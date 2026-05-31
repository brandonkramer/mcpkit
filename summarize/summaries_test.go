package summarize

import (
	"testing"

	"github.com/brandonkramer/mcpkit/envelope"
)

func TestToolSummaries(t *testing.T) {
	t.Parallel()
	briefs := ToolSummaries{
		"tasks.get": {
			Summarize: func(data any, _ Meta) (string, []string, bool) {
				id, ok := data.(string)
				if !ok {
					return "", nil, false
				}
				return "task " + id, []string{"tasks.list"}, true
			},
			Subjects: func(data any, _ Meta) ([]envelope.Subject, bool) {
				id, ok := data.(string)
				if !ok {
					return nil, false
				}
				return []envelope.Subject{{Kind: "task", ID: id}}, true
			},
		},
	}
	reg := RegistryFromSummaries(briefs)
	summary, next := reg.Summarize("tasks.get", "01TASK", Meta{})
	if summary != "task 01TASK" || len(next) != 1 {
		t.Fatalf("summary=%q next=%v", summary, next)
	}
	subjects := reg.Subjects("tasks.get", "01TASK", Meta{})
	if len(subjects) != 1 || subjects[0].ID != "01TASK" {
		t.Fatalf("subjects=%+v", subjects)
	}
}
