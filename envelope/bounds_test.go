package envelope

import "testing"

func TestCapSlice(t *testing.T) {
	t.Parallel()

	items := []int{1, 2, 3, 4, 5}

	cases := []struct {
		name      string
		limit     int
		wantLen   int
		truncated bool
	}{
		{name: "caps when over limit", limit: 3, wantLen: 3, truncated: true},
		{name: "unchanged when under limit", limit: 10, wantLen: 5, truncated: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			capped, truncated := CapSlice(items, tc.limit)
			if len(capped) != tc.wantLen || truncated != tc.truncated {
				t.Fatalf("capped=%v truncated=%v", capped, truncated)
			}
		})
	}

	meta := CapSliceMeta(items, 3)
	if !meta.Truncated || meta.Total != 5 || meta.Limit != 3 {
		t.Fatalf("meta=%+v", meta)
	}
}

func TestPreviewField(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		text      string
		limit     int
		want      string
		truncated bool
	}{
		{name: "short text", text: "hello", limit: 10, want: "hello"},
		{name: "long text", text: "0123456789abcdef", limit: 10, want: "0123456789", truncated: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			preview, truncated := PreviewField(tc.text, tc.limit)
			if preview != tc.want || truncated != tc.truncated {
				t.Fatalf("preview=%q truncated=%v", preview, truncated)
			}
		})
	}
}
