package summarize

import (
	"testing"

	"github.com/brandonkramer/mcpkit/envelope"
)

func TestResponderUsesRegistry(t *testing.T) {
	t.Parallel()
	reg := &Registry{
		Summarize: func(tool string, _ any, _ Meta) (string, []string) {
			return "done " + tool, []string{"next"}
		},
		Subjects: func(_ string, _ any, _ Meta) []envelope.Subject {
			return []envelope.Subject{{Kind: "item", ID: "1"}}
		},
	}
	_, envAny, err := NewResponder(reg).OK("demo", map[string]any{"x": 1}, Meta{})
	if err != nil {
		t.Fatal(err)
	}
	env := envAny.(envelope.Envelope)
	if env.Summary != "done demo" || len(env.Next) != 1 || len(env.Subjects) != 1 {
		t.Fatalf("env=%+v", env)
	}
}
