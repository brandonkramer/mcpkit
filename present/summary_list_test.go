package present

import (
	"testing"

	"github.com/brandonkramer/mcpkit/summarize"
)

func TestSummarizeToolList(t *testing.T) {
	t.Parallel()
	data, _ := WrapToolList([]string{"a", "b", "c"}, 2, "")
	summary, next := SummarizeToolList(data, summarize.Meta{Limit: 2}, summarize.ToolListSummary[string]{
		CountLabel: "item(s)",
		Next:       []string{"items.list"},
		NextFor: func(item string) []string {
			return []string{"items.get name=" + item}
		},
	})
	if summary != "3 item(s)" {
		t.Fatalf("summary=%q", summary)
	}
	if len(next) < 3 {
		t.Fatalf("next=%v", next)
	}
}
