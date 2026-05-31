package present

//
// ────────────────────────────────────────
// tool list payloads.
//

import (
	"fmt"

	"github.com/brandonkramer/mcpkit/summarize"
)

const defaultTruncationHint = "%d items available. Add filters to narrow results, or request a specific item by ID."

// ToolListData is the structured list payload placed in tool result data.
type ToolListData struct {
	Items     any    `json:"items"`
	Count     int    `json:"count"`
	Total     int    `json:"total,omitempty"`
	Truncated bool   `json:"truncated,omitempty"`
	Hint      string `json:"hint,omitempty"`
}

// WrapToolList caps items and attaches list metadata for tool result data.
func WrapToolList[T any](items []T, limit int, hintTemplate string) (ToolListData, summarize.Meta) {
	total := len(items)
	count := total
	if limit > 0 && limit < total {
		count = limit
	}
	truncated := count < total

	var hint string
	if truncated {
		switch {
		case hintTemplate != "":
			hint = fmt.Sprintf(hintTemplate, total)
		default:
			hint = fmt.Sprintf(defaultTruncationHint, total)
		}
	}

	data := ToolListData{
		Items:     items[:count],
		Count:     count,
		Total:     total,
		Truncated: truncated,
		Hint:      hint,
	}
	meta := summarize.Meta{}
	if truncated && limit > 0 {
		meta.Limit = limit
	}
	return data, meta
}
