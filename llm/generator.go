// Package llm provides small, provider-neutral text generation primitives for
// external agents. Providers are supplied by callers rather than embedded here.
package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	CommandEnv = "ENTIRE_LLM_COMMAND"
	ArgsEnv    = "ENTIRE_LLM_ARGS"
)

// Request is a text-generation request from an external agent.
type Request struct {
	Prompt string
	Model  string
}

// Generator generates text for a request. Implementations may call any LLM
// provider, or be a test double.
type Generator interface {
	Generate(context.Context, Request) (string, error)
}

// GeneratorFunc adapts a function into a Generator, making test doubles small.
type GeneratorFunc func(context.Context, Request) (string, error)

func (f GeneratorFunc) Generate(ctx context.Context, request Request) (string, error) {
	return f(ctx, request)
}

// CommandRunner executes a configured provider command.
type CommandRunner interface {
	Run(context.Context, string, []string, []byte) ([]byte, error)
}

// CommandGenerator sends a prompt to a command on standard input. Every
// occurrence of {model} in Args is replaced with the requested model.
type CommandGenerator struct {
	Command string
	Args    []string
	Runner  CommandRunner
}

func (g CommandGenerator) Generate(ctx context.Context, request Request) (string, error) {
	if strings.TrimSpace(g.Command) == "" {
		return "", errors.New("LLM command is required")
	}
	args := make([]string, len(g.Args))
	for i, arg := range g.Args {
		args[i] = strings.ReplaceAll(arg, "{model}", request.Model)
	}
	runner := g.Runner
	if runner == nil {
		runner = execRunner{}
	}
	output, err := runner.Run(ctx, g.Command, args, []byte(request.Prompt))
	if err != nil {
		return "", fmt.Errorf("run LLM command: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// FromEnvironment returns an optional command-backed generator. It is absent
// when ENTIRE_LLM_COMMAND is unset; ENTIRE_LLM_ARGS must be a JSON string array.
func FromEnvironment() (Generator, bool, error) {
	command := strings.TrimSpace(os.Getenv(CommandEnv))
	if command == "" {
		return nil, false, nil
	}
	var args []string
	if rawArgs := strings.TrimSpace(os.Getenv(ArgsEnv)); rawArgs != "" {
		if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
			return nil, false, fmt.Errorf("parse %s: %w", ArgsEnv, err)
		}
	}
	return CommandGenerator{Command: command, Args: args}, true, nil
}

type execRunner struct{}

func (execRunner) Run(ctx context.Context, command string, args []string, input []byte) ([]byte, error) {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdin = strings.NewReader(string(input))
	return cmd.Output()
}
