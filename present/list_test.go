package present

import "testing"

func TestListTruncates(t *testing.T) {
	t.Parallel()
	items := []string{"a", "b", "c"}
	env, meta := List(items, 2, "%d items")
	if env.Count != 2 || env.Total != 3 || !env.Truncated || meta.Limit != 2 {
		t.Fatalf("env=%+v meta=%+v", env, meta)
	}
}

func TestToolListItems(t *testing.T) {
	t.Parallel()
	raw := []string{"a", "b"}
	items, ok := ToolListItems[string](raw)
	if !ok || len(items) != 2 {
		t.Fatalf("items=%v ok=%v", items, ok)
	}
	env, _ := List(raw, 1, "")
	items, ok = ToolListItems[string](env)
	if !ok || len(items) != 1 || items[0] != "a" {
		t.Fatalf("items=%v ok=%v", items, ok)
	}
}

func TestToolListTotal(t *testing.T) {
	t.Parallel()
	if n, ok := ToolListTotal([]int{1, 2, 3}); !ok || n != 3 {
		t.Fatalf("n=%d ok=%v", n, ok)
	}
	env, _ := List([]int{1, 2, 3}, 2, "")
	if n, ok := ToolListTotal(env); !ok || n != 3 {
		t.Fatalf("n=%d ok=%v", n, ok)
	}
}
