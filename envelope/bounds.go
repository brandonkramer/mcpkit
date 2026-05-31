package envelope

//
// ────────────────────────────────────────
// bounds and previews.
//

import "fmt"

// BoundList is a capped slice with truncation metadata.
type BoundList[T any] struct {
	Items     []T
	Truncated bool
	Total     int
	Limit     int
}

// CapSlice returns at most limit items and reports whether the input was truncated.
func CapSlice[T any](items []T, limit int) (out []T, truncated bool) {
	if limit <= 0 || len(items) <= limit {
		return items, false
	}
	return items[:limit], true
}

// CapSliceMeta returns a capped slice with truncation metadata.
func CapSliceMeta[T any](items []T, limit int) BoundList[T] {
	out, truncated := CapSlice(items, limit)
	return BoundList[T]{
		Items:     out,
		Truncated: truncated,
		Total:     len(items),
		Limit:     limit,
	}
}

// TruncateText returns a preview capped at limit bytes/chars and whether truncation occurred.
func TruncateText(text string, limit int) (preview string, truncated bool) {
	if limit <= 0 || len(text) <= limit {
		return text, false
	}
	return text[:limit], true
}

// PreviewField is an alias for TruncateText.
func PreviewField(text string, limit int) (preview string, truncated bool) {
	return TruncateText(text, limit)
}

// AppendLimitHint appends a paging hint when a list response hit its limit.
func AppendLimitHint(next []string, atLimit bool, hint string) []string {
	if atLimit && hint != "" {
		return append(next, hint)
	}
	return next
}

// SummarizeCount formats a count summary and optional next hints.
func SummarizeCount(count int, unit string, next []string) (summary string, out []string) {
	return fmt.Sprintf("%d %s", count, unit), append([]string(nil), next...)
}
