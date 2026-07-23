package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// runEpisode runs one full agent loop on a fresh copy of the scenario repo.
func runEpisode(sc scenario, mdl modelClient, model, label, system string, grounding bool, idx int) episodeResult {
	r := episodeResult{Scenario: sc.Name, Model: model, Condition: label, Episode: idx}

	root, err := os.MkdirTemp("", "prism-ep-")
	if err != nil {
		r.Error = "mktemp: " + err.Error()
		return r
	}
	defer os.RemoveAll(root)
	if err := copyDir(filepath.Join(sc.Dir, "repo"), root); err != nil {
		r.Error = "copy repo: " + err.Error()
		return r
	}
	timeout := time.Duration(sc.ToolTimeoutSec) * time.Second

	// sanity: the scenario must START in a failing state, otherwise there is
	// nothing to solve and the measurement is meaningless.
	if runSuccessCmd(root, sc.SuccessCmd, timeout) {
		r.Error = "scenario already passes at t=0 (broken scenario)"
		return r
	}

	// manipulation check: run the trap probe on the initial environment and
	// record whether the targeted trigger event actually surfaces. This is a
	// scenario property (independent of which commands the model chooses), so
	// the trigger rate is ~100% on a scenario that targets this component and
	// ~0% on a neutral one — exactly the pre-registered check.
	probe := sc.TriggerProbeCmd
	if probe == "" {
		probe = sc.SuccessCmd
	}
	if sc.TriggerPattern != "" {
		r.Triggered = strings.Contains(runCapture(root, probe, timeout), sc.TriggerPattern)
	}

	// build the system prompt for this condition
	var sysParts []string
	if grounding {
		sysParts = append(sysParts, groundingBlock(root))
	}
	if system != "" {
		sysParts = append(sysParts, system)
	}
	messages := []map[string]any{}
	if len(sysParts) > 0 {
		messages = append(messages, map[string]any{"role": "system", "content": strings.Join(sysParts, "\n\n")})
	}
	messages = append(messages, map[string]any{"role": "user", "content": sc.Task})

	tools := toolSchemas()
	seen := map[string]bool{} // for loop detection: name+args already issued

	for step := 1; step <= sc.MaxSteps; step++ {
		assistant, calls, err := mdl.next(messages, tools)
		if err != nil {
			r.Error = fmt.Sprintf("model at step %d: %v", step, err)
			break
		}
		messages = append(messages, assistant)
		r.Steps = step

		if len(calls) == 0 {
			// model declared done (no tool call): stop and let the final
			// success check below decide.
			break
		}
		for _, tc := range calls {
			key := tc.Name + "\x00" + tc.Args
			if seen[key] {
				r.LoopRepeats++
			}
			seen[key] = true
			r.ToolSeq = append(r.ToolSeq, tc.Name)

			result := execTool(root, timeout, tc.Name, tc.Args)
			messages = append(messages, map[string]any{
				"role": "tool", "tool_call_id": tc.ID, "content": result,
			})
		}
		if runSuccessCmd(root, sc.SuccessCmd, timeout) {
			r.Success = true
			r.SuccessStep = step
			break
		}
	}
	if !r.Success {
		// final check in case success came on the very last step's actions
		r.Success = runSuccessCmd(root, sc.SuccessCmd, timeout)
		if r.Success && r.SuccessStep == 0 {
			r.SuccessStep = r.Steps
		}
	}
	saveTranscript(sc.Name, model, label, idx, messages)
	return r
}

func saveTranscript(scenario, model, label string, idx int, messages []map[string]any) {
	dir := "artefacts/transcripts"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return
	}
	name := fmt.Sprintf("%s--%s--%s--%d.json", sanitize(scenario), sanitize(model), sanitize(label), idx)
	f, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		return
	}
	defer f.Close()
	writeJSON(f, messages)
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, p)
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(p, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func sanitize(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '/' || r == ' ' || r == ':' {
			return '-'
		}
		return r
	}, s)
}

func writeJSON(w io.Writer, v any) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
