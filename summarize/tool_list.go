package summarize

import "github.com/brandonkramer/mcpkit/envelope"

// ToolListSummary configures summary and next-hint generation for a list tool result.
type ToolListSummary[T any] struct {
	CountLabel string
	Next       []string
	PageHint   string
	SampleSize int
	NextFor    func(T) []string
}

// ToolListItems builds envelope summary and next hints from a counted tool list payload.
func ToolListItems[T any](items []T, total int, meta Meta, spec ToolListSummary[T]) (summary string, next []string) {
	next = append([]string(nil), spec.Next...)
	sample := spec.SampleSize
	if sample <= 0 {
		sample = 3
	}
	if spec.NextFor != nil {
		for i, item := range items {
			if i >= sample {
				break
			}
			next = append(next, spec.NextFor(item)...)
		}
	}
	next = envelope.AppendLimitHint(next, meta.Limit > 0 && total >= meta.Limit, spec.PageHint)
	return envelope.SummarizeCount(total, spec.CountLabel, next)
}
