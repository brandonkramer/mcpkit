package summarize

import "github.com/brandonkramer/mcpkit/envelope"

// SummaryFunc returns a short summary and continuation hints for a tool result.
type SummaryFunc func(tool string, data any, meta Meta) (summary string, next []string)

// SubjectsFunc returns entity references extracted from a tool result.
type SubjectsFunc func(tool string, data any, meta Meta) []envelope.Subject

// Registry holds pluggable summarization hooks for a tool server.
type Registry struct {
	Summarize SummaryFunc
	Subjects  SubjectsFunc
}

// Default returns a registry that uses the tool name as summary.
func Default() *Registry {
	return &Registry{
		Summarize: func(tool string, _ any, _ Meta) (string, []string) { return tool, nil },
		Subjects:  func(string, any, Meta) []envelope.Subject { return nil },
	}
}
