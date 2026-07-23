package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// modelClient abstracts the decision maker: either a real vLLM model or a
// scripted replay (phase-0 smoke, zero API cost).
type modelClient interface {
	// next returns the assistant message (raw, to append verbatim so
	// tool_calls round-trip) and the parsed tool calls it requested.
	next(messages []map[string]any, tools []map[string]any) (assistant map[string]any, calls []toolCall, err error)
}

type toolCall struct {
	ID   string
	Name string
	Args string
}

// ----- real vLLM (OpenAI-compatible /chat/completions) -----

type vllmClient struct {
	endpoint string
	model    string
}

func (c *vllmClient) next(messages, tools []map[string]any) (map[string]any, []toolCall, error) {
	body, _ := json.Marshal(map[string]any{
		"model":       c.model,
		"messages":    messages,
		"tools":       tools,
		"tool_choice": "auto",
		"temperature": 0.0,
	})
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Post(c.endpoint+"/chat/completions", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("%s: %.300s", resp.Status, string(raw))
	}
	var parsed struct {
		Choices []struct {
			Message map[string]any `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, nil, err
	}
	if len(parsed.Choices) == 0 {
		return nil, nil, fmt.Errorf("no choices in response")
	}
	msg := parsed.Choices[0].Message
	return msg, parseToolCalls(msg), nil
}

func parseToolCalls(msg map[string]any) []toolCall {
	raw, ok := msg["tool_calls"].([]any)
	if !ok {
		return nil
	}
	var calls []toolCall
	for _, tc := range raw {
		m, ok := tc.(map[string]any)
		if !ok {
			continue
		}
		id, _ := m["id"].(string)
		fn, ok := m["function"].(map[string]any)
		if !ok {
			continue
		}
		name, _ := fn["name"].(string)
		args, _ := fn["arguments"].(string)
		calls = append(calls, toolCall{ID: id, Name: name, Args: args})
	}
	return calls
}

// ----- scripted replay (phase-0 smoke) -----

type scriptedModel struct {
	steps []scriptedStep
	i     int
}

type scriptedStep struct {
	Tool string            `json:"tool"`
	Args map[string]string `json:"args"`
}

func loadScriptedModel(path string) *scriptedModel {
	raw, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var steps []scriptedStep
	if err := json.Unmarshal(raw, &steps); err != nil {
		panic(err)
	}
	return &scriptedModel{steps: steps}
}

func (m *scriptedModel) next(_, _ []map[string]any) (map[string]any, []toolCall, error) {
	if m.i >= len(m.steps) {
		// no more scripted actions → signal "done" (no tool calls)
		return map[string]any{"role": "assistant", "content": "done"}, nil, nil
	}
	s := m.steps[m.i]
	m.i++
	argsJSON, _ := json.Marshal(s.Args)
	id := fmt.Sprintf("call_%d", m.i)
	assistant := map[string]any{
		"role":    "assistant",
		"content": nil,
		"tool_calls": []any{map[string]any{
			"id": id, "type": "function",
			"function": map[string]any{"name": s.Tool, "arguments": string(argsJSON)},
		}},
	}
	return assistant, []toolCall{{ID: id, Name: s.Tool, Args: string(argsJSON)}}, nil
}
