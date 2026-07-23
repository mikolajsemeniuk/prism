// Command runner executes the ablation experiment (roadmap step 5, pilot
// phase). It runs ONE condition (system prompt + optional grounding) of ONE
// scenario against ONE model, k times, in a controlled-but-real environment:
// each episode gets a fresh copy of the scenario's mini-repo in a temp dir,
// the model drives a real agent loop (list/read/write/bash tools that ACTUALLY
// execute on the copy), and success is checked by really running the
// scenario's success command. Nothing is faked — only the environment is
// controlled (reset per episode, a planted trap in the repo files).
//
// Composability mirrors the decomposition pipeline: the runner does one
// condition per invocation and appends to artefacts/episodes.jsonl; the
// Makefile drives the condition × model × k grid. This keeps the runner small
// and lets a scenario change be just a new folder under scenarios/.
//
// Metrics per episode: success (bool), steps, loop_repeats (identical tool
// calls repeated — the qwen-looping signal), triggered (did the planted trap
// event actually surface — the manipulation check that proves the scenario
// exercises the targeted component).
//
// Usage:
//
//	go run ./cmd/runner -scenario scenarios/selfrepair-01 -model qwen -label baseline -k 4
//	go run ./cmd/runner -scenario scenarios/selfrepair-01 -model qwen -label fragment \
//	    -system-file scenarios/selfrepair-01/stimuli/fragment.txt -k 4
//	go run ./cmd/runner -scenario scenarios/selfrepair-01 -scripted scenarios/selfrepair-01/scripted-good.json -label smoke -k 1
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type scenario struct {
	Name           string `json:"-"`
	Dir            string `json:"-"`
	Task            string `json:"task"`
	SuccessCmd      string `json:"success_cmd"`
	TriggerProbeCmd string `json:"trigger_probe_cmd"`
	TriggerPattern  string `json:"trigger_pattern"`
	MaxSteps        int    `json:"max_steps"`
	ToolTimeoutSec  int    `json:"tool_timeout_sec"`
}

type episodeResult struct {
	Scenario    string   `json:"scenario"`
	Model       string   `json:"model"`
	Condition   string   `json:"condition"`
	Episode     int      `json:"episode"`
	Success     bool     `json:"success"`
	Steps       int      `json:"steps"`
	LoopRepeats int      `json:"loop_repeats"`
	Triggered   bool     `json:"triggered"`
	SuccessStep int      `json:"success_step"` // step at which success first observed, 0 if never
	ToolSeq     []string `json:"tool_seq"`
	Error       string   `json:"error,omitempty"`
}

func main() {
	log.SetFlags(0)
	scenarioDir := flag.String("scenario", "", "scenario folder (required)")
	model := flag.String("model", "local", "model name (label + vLLM model id)")
	endpoint := flag.String("endpoint", "http://localhost:8000/v1", "OpenAI-compatible endpoint")
	label := flag.String("label", "baseline", "condition label for logging")
	systemFile := flag.String("system-file", "", "file with the system prompt (fragment/placebo); empty = baseline")
	grounding := flag.Bool("grounding", false, "inject an environment-context block (cwd + file listing)")
	scripted := flag.String("scripted", "", "replay tool calls from this file instead of calling a model (phase-0 smoke)")
	k := flag.Int("k", 4, "episodes to run")
	maxStepsFlag := flag.Int("max-steps", 0, "override scenario max_steps (0 = use scenario)")
	out := flag.String("out", "artefacts/episodes.jsonl", "episode log (appended)")
	flag.Parse()
	if *scenarioDir == "" {
		log.Fatal("-scenario is required")
	}

	sc := loadScenario(*scenarioDir)
	if *maxStepsFlag > 0 {
		sc.MaxSteps = *maxStepsFlag
	}

	var system string
	if *systemFile != "" {
		raw, err := os.ReadFile(*systemFile)
		if err != nil {
			log.Fatalf("system-file: %v", err)
		}
		system = strings.TrimSpace(string(raw))
	}

	var mdl modelClient
	if *scripted != "" {
		mdl = loadScriptedModel(*scripted)
	} else {
		mdl = &vllmClient{endpoint: *endpoint, model: *model}
	}

	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		log.Fatal(err)
	}
	logf, err := os.OpenFile(*out, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	defer logf.Close()
	enc := json.NewEncoder(logf)

	fmt.Printf("scenario=%s model=%s condition=%s grounding=%v k=%d max_steps=%d\n",
		sc.Name, *model, *label, *grounding, *k, sc.MaxSteps)

	var results []episodeResult
	for i := 1; i <= *k; i++ {
		r := runEpisode(sc, mdl, *model, *label, system, *grounding, i)
		results = append(results, r)
		if err := enc.Encode(r); err != nil {
			log.Fatal(err)
		}
		status := "fail"
		if r.Success {
			status = "OK"
		}
		fmt.Printf("  ep %d: %-4s steps=%-3d loops=%-3d trigger=%-5v %s\n",
			i, status, r.Steps, r.LoopRepeats, r.Triggered, r.Error)
	}
	printSummary(*label, results)
}

func loadScenario(dir string) scenario {
	raw, err := os.ReadFile(filepath.Join(dir, "scenario.json"))
	if err != nil {
		log.Fatalf("scenario: %v", err)
	}
	var sc scenario
	if err := json.Unmarshal(raw, &sc); err != nil {
		log.Fatalf("scenario.json: %v", err)
	}
	sc.Dir = dir
	sc.Name = filepath.Base(strings.TrimRight(dir, "/"))
	if sc.MaxSteps == 0 {
		sc.MaxSteps = 30
	}
	if sc.ToolTimeoutSec == 0 {
		sc.ToolTimeoutSec = 30
	}
	if _, err := os.Stat(filepath.Join(dir, "repo")); err != nil {
		log.Fatalf("scenario %s has no repo/ subfolder", sc.Name)
	}
	return sc
}

func printSummary(label string, rs []episodeResult) {
	n := len(rs)
	if n == 0 {
		return
	}
	succ, trig, sumSteps, sumLoops, nSucc := 0, 0, 0, 0, 0
	for _, r := range rs {
		if r.Success {
			succ++
			sumSteps += r.Steps
			nSucc++
		}
		if r.Triggered {
			trig++
		}
		sumLoops += r.LoopRepeats
	}
	meanSteps := 0.0
	if nSucc > 0 {
		meanSteps = float64(sumSteps) / float64(nSucc)
	}
	fmt.Printf("── %s: success %d/%d (%.0f%%)  trigger %d/%d  mean_steps(success)=%.1f  mean_loops=%.1f\n",
		label, succ, n, 100*float64(succ)/float64(n), trig, n, meanSteps, float64(sumLoops)/float64(n))
}

// groundingBlock builds the environment-context stimulus from the real repo
// copy: working directory + a listing of its files. This is source C of the
// stimulus roster (harness-generated, not a corpus fragment).
func groundingBlock(repoDir string) string {
	var files []string
	filepath.WalkDir(repoDir, func(p string, d os.DirEntry, e error) error {
		if e != nil || d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(repoDir, p)
		files = append(files, rel)
		return nil
	})
	sort.Strings(files)
	var b strings.Builder
	b.WriteString("Environment:\n")
	b.WriteString("Working directory: . (the project root)\n")
	b.WriteString("Files in the working directory:\n")
	for _, f := range files {
		fmt.Fprintf(&b, "  %s\n", f)
	}
	return b.String()
}
