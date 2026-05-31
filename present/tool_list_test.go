package present

import "testing"

func TestWrapToolList(t *testing.T) {
	t.Parallel()
	data, meta := WrapToolList([]string{"a", "b", "c"}, 2, "%d items")
	if data.Count != 2 || data.Total != 3 || !data.Truncated || meta.Limit != 2 {
		t.Fatalf("data=%+v meta=%+v", data, meta)
	}
	items, ok := data.Items.([]string)
	if !ok || len(items) != 2 || items[0] != "a" {
		t.Fatalf("items=%v ok=%v", items, ok)
	}
}

func TestWrapToolListDefaultHint(t *testing.T) {
	t.Parallel()
	data, _ := WrapToolList([]int{1, 2, 3}, 1, "")
	if data.Hint == "" {
		t.Fatal("expected default hint")
	}
}
