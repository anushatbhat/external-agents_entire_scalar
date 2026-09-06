package llm

import (
	"context"
	"reflect"
	"testing"
)

type recordingRunner struct {
	command string
	args    []string
	input   []byte
}

func (r *recordingRunner) Run(_ context.Context, command string, args []string, input []byte) ([]byte, error) {
	r.command, r.args, r.input = command, append([]string(nil), args...), append([]byte(nil), input...)
	return []byte("  generated summary\n"), nil
}

func TestCommandGeneratorForwardsPromptAndModel(t *testing.T) {
	runner := &recordingRunner{}
	got, err := (CommandGenerator{
		Command: "provider",
		Args:    []string{"generate", "--model={model}"},
		Runner:  runner,
	}).Generate(context.Background(), Request{Prompt: "summarize this", Model: "model-1"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "generated summary" {
		t.Fatalf("generated text = %q", got)
	}
	if runner.command != "provider" || !reflect.DeepEqual(runner.args, []string{"generate", "--model=model-1"}) || string(runner.input) != "summarize this" {
		t.Fatalf("command = %q, args = %#v, input = %q", runner.command, runner.args, runner.input)
	}
}
