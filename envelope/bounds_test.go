package envelope

import "testing"

func TestCapSlice(t *testing.T) {
	t.Parallel()
	items := []int{1, 2, 3, 4, 5}
	capped, truncated := CapSlice(items, 3)
	if len(capped) != 3 || !truncated || capped[2] != 3 {
		t.Fatalf("capped=%v truncated=%v", capped, truncated)
	}
	unchanged, truncated := CapSlice(items, 10)
	if len(unchanged) != 5 || truncated {
		t.Fatalf("unchanged=%v truncated=%v", unchanged, truncated)
	}
	meta := CapSliceMeta(items, 3)
	if !meta.Truncated || meta.Total != 5 || meta.Limit != 3 {
		t.Fatalf("meta=%+v", meta)
	}
}

func TestPreviewField(t *testing.T) {
	t.Parallel()
	preview, truncated := PreviewField("hello", 10)
	if preview != "hello" || truncated {
		t.Fatalf("preview=%q truncated=%v", preview, truncated)
	}
	long := "0123456789abcdef"
	preview, truncated = PreviewField(long, 10)
	if preview != long[:10] || !truncated {
		t.Fatalf("preview=%q truncated=%v", preview, truncated)
	}
}
