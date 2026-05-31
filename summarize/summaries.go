package summarize

import "github.com/brandonkramer/mcpkit/envelope"

// ToolSummary summarizes one MCP tool result and extracts envelope subjects.
type ToolSummary struct {
	Summarize func(data any, meta Meta) (summary string, next []string, ok bool)
	Subjects  func(data any, meta Meta) ([]envelope.Subject, bool)
}

// ToolSummaries maps MCP tool names to result summarizers.
type ToolSummaries map[string]ToolSummary

// Summarize dispatches to a registered tool summarizer or returns the tool name.
func (b ToolSummaries) Summarize(tool string, data any, meta Meta) (summary string, next []string) {
	if fn, ok := b[tool]; ok && fn.Summarize != nil {
		if summary, next, ok := fn.Summarize(data, meta); ok {
			return summary, next
		}
	}
	return tool, nil
}

// Subjects dispatches to a registered tool summarizer.
func (b ToolSummaries) Subjects(tool string, data any, meta Meta) []envelope.Subject {
	if fn, ok := b[tool]; ok && fn.Subjects != nil {
		if subjects, ok := fn.Subjects(data, meta); ok {
			return subjects
		}
	}
	return nil
}

// RegistryFromSummaries builds a Registry backed by per-tool summarizers.
func RegistryFromSummaries(summaries ToolSummaries) *Registry {
	return &Registry{Summarize: summaries.Summarize, Subjects: summaries.Subjects}
}
