package omp

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/entireio/external-agents/agents/entire-agent-omp/internal/protocol"
	"github.com/entireio/external-agents/llm"
)

func TestDetectUsesPath(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("PATH", dir)
	if New().Detect().Present {
		t.Fatal("Detect().Present = true without omp on PATH")
	}
	if err := os.WriteFile(filepath.Join(dir, "omp"), []byte("#!/bin/sh\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	if !New().Detect().Present {
		t.Fatal("Detect().Present = false with omp on PATH")
	}
}

func TestGenerateTextUsesInjectedGenerator(t *testing.T) {
	var request llm.Request
	agent := NewWithGenerator(llm.GeneratorFunc(func(_ context.Context, got llm.Request) (string, error) {
		request = got
		return "short summary", nil
	}))
	if !agent.Info().Capabilities.TextGenerator {
		t.Fatal("text_generator = false with configured generator")
	}
	text, err := agent.GenerateText("summarize this", "model-1")
	if err != nil {
		t.Fatal(err)
	}
	if text != "short summary" || request != (llm.Request{Prompt: "summarize this", Model: "model-1"}) {
		t.Fatalf("text = %q, request = %#v", text, request)
	}
}

func TestGetSessionID(t *testing.T) {
	agent := New()
	if got := agent.GetSessionID(nil); got != "" {
		t.Fatalf("GetSessionID(nil) = %q", got)
	}
	if got := agent.GetSessionID(&protocol.HookInputJSON{SessionID: "session-1"}); got != "session-1" {
		t.Fatalf("GetSessionID() = %q", got)
	}
}

func TestFormatResumeCommand(t *testing.T) {
	tests := []struct{ id, want string }{
		{"", "omp --continue"},
		{"019abc-safe", "omp --resume 019abc-safe"},
		{"two words", "omp --resume 'two words'"},
		{"a'b", "omp --resume 'a'\\''b'"},
		{"$(touch nope)", "omp --resume '$(touch nope)'"},
	}
	for _, test := range tests {
		if got := New().FormatResumeCommand(test.id); got != test.want {
			t.Errorf("FormatResumeCommand(%q) = %q, want %q", test.id, got, test.want)
		}
	}
}
