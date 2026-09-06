package omp

import (
	"context"
	"os/exec"
	"strings"

	"github.com/entireio/external-agents/agents/entire-agent-omp/internal/protocol"
	"github.com/entireio/external-agents/llm"
)

type Agent struct {
	Generator llm.Generator
}

func New() *Agent {
	generator, _, err := llm.FromEnvironment()
	if err != nil {
		return &Agent{}
	}
	return NewWithGenerator(generator)
}

// NewWithGenerator creates an agent with an optional text generator.
func NewWithGenerator(generator llm.Generator) *Agent { return &Agent{Generator: generator} }

func (a *Agent) Info() protocol.InfoResponse {
	return protocol.InfoResponse{
		ProtocolVersion: protocol.ProtocolVersion,
		Name:            "omp",
		Type:            "Oh My Pi",
		Description:     "Oh My Pi terminal coding agent",
		IsPreview:       true,
		ProtectedDirs:   []string{".omp"},
		ProtectedFiles:  []string{extensionRelFile},
		HookNames: []string{
			HookNameSessionStart,
			HookNameAgentStart,
			HookNameAgentEnd,
			HookNameSessionShutdown,
		},
		Capabilities: protocol.DeclaredCapabilities{
			Hooks:              true,
			TranscriptAnalyzer: true,
			CompactTranscript:  true,
			TextGenerator:      a.Generator != nil,
			UsesTerminal:       true,
		},
	}
}

func (a *Agent) GenerateText(prompt, model string) (string, error) {
	if a.Generator == nil {
		return "", protocol.ErrTextGeneratorUnavailable
	}
	return a.Generator.Generate(context.Background(), llm.Request{Prompt: prompt, Model: model})
}

func (a *Agent) Detect() protocol.DetectResponse {
	_, err := exec.LookPath("omp")
	return protocol.DetectResponse{Present: err == nil}
}

func (a *Agent) GetSessionID(input *protocol.HookInputJSON) string {
	if input == nil {
		return ""
	}
	return input.SessionID
}

func (a *Agent) FormatResumeCommand(sessionID string) string {
	if sessionID == "" {
		return "omp --continue"
	}
	return "omp --resume " + shellQuote(sessionID)
}

func shellQuote(value string) string {
	if strings.IndexFunc(value, func(r rune) bool { return !safeShellRune(r) }) == -1 {
		return value
	}
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func safeShellRune(r rune) bool {
	return r == '-' || r == '_' || r == '.' || r == ':' ||
		(r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
