package summarize

import "testing"

func TestMetaExtra(t *testing.T) {
	t.Parallel()
	meta := Meta{Limit: 10}
	meta.Put("table", "notes")
	if meta.String("table") != "notes" || meta.Limit != 10 {
		t.Fatalf("meta=%+v", meta)
	}
	meta.Put("page", 3)
	if meta.Int("page") != 3 {
		t.Fatalf("page=%d", meta.Int("page"))
	}
}
