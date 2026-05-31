package present_test

import (
	"testing"

	"github.com/brandonkramer/mcpkit/present"
)

func TestWrapToolListExport(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		items     []string
		limit     int
		wantCount int
		truncated bool
	}{
		{name: "no cap", items: []string{"a", "b"}, limit: 0, wantCount: 2},
		{name: "truncated", items: []string{"a", "b", "c"}, limit: 2, wantCount: 2, truncated: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, meta := present.WrapToolList(tc.items, tc.limit, "%d total")
			if data.Count != tc.wantCount || data.Truncated != tc.truncated {
				t.Fatalf("data=%+v", data)
			}
			if tc.truncated && meta.Limit != tc.limit {
				t.Fatalf("meta=%+v", meta)
			}
		})
	}
}

func TestRedactTextExport(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		text        string
		limit       int
		wantPrev    string
		truncated   bool
		wantCleared string
	}{
		{name: "short", text: "hello", limit: 10, wantPrev: "hello", wantCleared: "hello"},
		{name: "long", text: "0123456789abcdef", limit: 10, wantPrev: "0123456789", truncated: true, wantCleared: ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			prev, cleared := present.RedactText(tc.text, tc.limit)
			if prev.Preview != tc.wantPrev || prev.Truncated != tc.truncated || cleared != tc.wantCleared {
				t.Fatalf("prev=%+v cleared=%q", prev, cleared)
			}
		})
	}
}
