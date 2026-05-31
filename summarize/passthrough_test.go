package summarize

import "testing"

func TestPassthrough(t *testing.T) {
	t.Parallel()
	reg := Passthrough()
	summary, next := reg.Summarize("items.list", nil, Meta{})
	if summary != "items.list" || len(next) != 0 {
		t.Fatalf("summary=%q next=%v", summary, next)
	}
	if len(reg.Subjects("items.list", nil, Meta{})) != 0 {
		t.Fatal("expected no subjects")
	}
}
