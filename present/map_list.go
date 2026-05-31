package present

import "github.com/brandonkramer/mcpkit/summarize"

// MapCapList caps items and maps each element into a presentation view.
func MapCapList[T, V any](items []T, limit int, mapFn func(*T) V) ([]V, summarize.Meta) {
	capped, meta := CapSlice(items, limit)
	out := make([]V, len(capped))
	for i := range capped {
		out[i] = mapFn(&capped[i])
	}
	return out, meta
}
