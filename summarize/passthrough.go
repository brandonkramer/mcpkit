package summarize

import "github.com/brandonkramer/mcpkit/envelope"

// Passthrough returns a registry that uses the tool name as summary and omits subjects.
// Use with summarize.NewResponder for proxy tools before domain summarizers are ready.
func Passthrough() *Registry {
	return &Registry{
		Summarize: func(tool string, _ any, _ Meta) (string, []string) { return tool, nil },
		Subjects:  func(string, any, Meta) []envelope.Subject { return nil },
	}
}
