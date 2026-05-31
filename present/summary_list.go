package present

//
// ────────────────────────────────────────
// list summarization.
//

import "github.com/brandonkramer/mcpkit/summarize"

// SummarizeToolList builds envelope summary and next hints from tool result list data.
func SummarizeToolList[T any](data any, meta summarize.Meta, spec summarize.ToolListSummary[T]) (summary string, next []string) {
	total, ok := ToolListTotal(data)
	if !ok {
		return "", nil
	}
	items, _ := ToolListItems[T](data)
	return summarize.ToolListItems(items, total, meta, spec)
}
