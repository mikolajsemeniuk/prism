package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const maxToolOutput = 4000

// toolSchemas is the tool list advertised to the model (OpenAI function
// format). The tools really execute on the per-episode repo copy.
func toolSchemas() []map[string]any {
	strParam := func(desc string) map[string]any {
		return map[string]any{"type": "string", "description": desc}
	}
	fn := func(name, desc string, props map[string]any, required []string) map[string]any {
		// Coerce nil → empty slice: strict tool-call parsers (Mistral,
		// Llama, Granite) reject "required": null and require an array.
		if required == nil {
			required = []string{}
		}
		return map[string]any{
			"type": "function",
			"function": map[string]any{
				"name": name, "description": desc,
				"parameters": map[string]any{"type": "object", "properties": props, "required": required},
			},
		}
	}
	return []map[string]any{
		fn("ls", "List files in a directory of the project (relative path; default '.').",
			map[string]any{"path": strParam("relative directory path")}, nil),
		fn("read_file", "Read a text file from the project.",
			map[string]any{"path": strParam("relative file path")}, []string{"path"}),
		fn("write_file", "Create or overwrite a text file in the project.",
			map[string]any{"path": strParam("relative file path"), "content": strParam("full file content")},
			[]string{"path", "content"}),
		fn("bash", "Run a shell command in the project root and return its combined stdout+stderr. Use this to run tests, scripts, installers, etc.",
			map[string]any{"command": strParam("the shell command")}, []string{"command"}),
	}
}

// execTool runs one tool call against the repo copy at root and returns the
// text result the model will see.
func execTool(root string, timeout time.Duration, name, argsJSON string) string {
	var args map[string]string
	_ = json.Unmarshal([]byte(argsJSON), &args)
	switch name {
	case "ls":
		return toolLS(root, args["path"])
	case "read_file":
		return toolRead(root, args["path"])
	case "write_file":
		return toolWrite(root, args["path"], args["content"])
	case "bash":
		return toolBash(root, timeout, args["command"])
	default:
		return "error: unknown tool " + name
	}
}

// safeJoin keeps tool paths inside the repo copy (defense against the model
// wandering outside its sandbox via ../).
func safeJoin(root, rel string) (string, bool) {
	if rel == "" {
		rel = "."
	}
	p := filepath.Clean(filepath.Join(root, rel))
	if p != root && !strings.HasPrefix(p, root+string(os.PathSeparator)) {
		return "", false
	}
	return p, true
}

func toolLS(root, rel string) string {
	p, ok := safeJoin(root, rel)
	if !ok {
		return "error: path escapes the project"
	}
	entries, err := os.ReadDir(p)
	if err != nil {
		return "error: " + err.Error()
	}
	var b strings.Builder
	for _, e := range entries {
		suffix := ""
		if e.IsDir() {
			suffix = "/"
		}
		fmt.Fprintf(&b, "%s%s\n", e.Name(), suffix)
	}
	if b.Len() == 0 {
		return "(empty directory)"
	}
	return b.String()
}

func toolRead(root, rel string) string {
	p, ok := safeJoin(root, rel)
	if !ok {
		return "error: path escapes the project"
	}
	raw, err := os.ReadFile(p)
	if err != nil {
		return "error: " + err.Error()
	}
	return truncate(string(raw))
}

func toolWrite(root, rel, content string) string {
	p, ok := safeJoin(root, rel)
	if !ok {
		return "error: path escapes the project"
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return "error: " + err.Error()
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		return "error: " + err.Error()
	}
	return "wrote " + rel
}

func toolBash(root string, timeout time.Duration, command string) string {
	if command == "" {
		return "error: empty command"
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "-lc", command)
	cmd.Dir = root
	// isolate: no inherited network-y env surprises; keep PATH for python etc.
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	outBytes, err := cmd.CombinedOutput()
	out := string(outBytes)
	if ctx.Err() == context.DeadlineExceeded {
		out += "\n[command timed out]"
	} else if err != nil {
		out += fmt.Sprintf("\n[exit: %v]", err)
	}
	if strings.TrimSpace(out) == "" {
		out = "(no output)"
	}
	return truncate(out)
}

// runSuccessCmd runs the scenario's success command in the repo copy and
// reports whether it exited 0 (task solved).
func runSuccessCmd(root, successCmd string, timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "-lc", successCmd)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	return cmd.Run() == nil
}

// runCapture runs a command in the repo copy and returns its combined output
// (used by the manipulation-check probe, which needs the full text).
func runCapture(root, command string, timeout time.Duration) string {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "-lc", command)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	out, _ := cmd.CombinedOutput()
	return string(out)
}

func truncate(s string) string {
	if len(s) <= maxToolOutput {
		return s
	}
	return s[:maxToolOutput] + "\n[...truncated...]"
}
