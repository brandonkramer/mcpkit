package present

import (
	"reflect"

	"github.com/brandonkramer/mcpkit/envelope"
	"github.com/brandonkramer/mcpkit/summarize"
)

// List wraps items for tool result data. Prefer WrapToolList for new code.
func List[T any](items []T, limit int, hintTemplate string) (ToolListData, summarize.Meta) {
	return WrapToolList(items, limit, hintTemplate)
}

// CapSlice caps a slice without list metadata wrapping.
func CapSlice[T any](items []T, limit int) ([]T, summarize.Meta) {
	out, truncated := envelope.CapSlice(items, limit)
	meta := summarize.Meta{}
	if truncated {
		meta.Limit = limit
	}
	return out, meta
}

// PreviewField returns a bounded text preview.
func PreviewField(text string, limit int) (preview string, truncated bool) {
	return envelope.PreviewField(text, limit)
}

func listDataTotal(data any) (int, bool) {
	list, ok := data.(ToolListData)
	if !ok {
		return 0, false
	}
	return list.Total, true
}

func listDataItems[T any](data any) ([]T, bool) {
	list, ok := data.(ToolListData)
	if !ok {
		return nil, false
	}
	items, ok := list.Items.([]T)
	if !ok {
		return nil, false
	}
	return items, true
}

// ToolListItems returns typed items from tool result data (structured list or bare slice).
func ToolListItems[T any](data any) ([]T, bool) {
	if items, ok := listDataItems[T](data); ok {
		return items, true
	}
	items, ok := data.([]T)
	if !ok {
		return nil, false
	}
	return items, true
}

// ToolListTotal returns the item total for tool result list data.
func ToolListTotal(data any) (int, bool) {
	if n, ok := listDataTotal(data); ok {
		return n, true
	}
	if data == nil {
		return 0, false
	}
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return 0, false
	}
	return v.Len(), true
}
