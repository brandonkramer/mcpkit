package present

import "testing"

func TestMapCapList(t *testing.T) {
	t.Parallel()
	type in struct{ N int }
	type out struct{ Twice int }
	views, meta := MapCapList([]in{{1}, {2}, {3}}, 2, func(v *in) out {
		return out{Twice: v.N * 2}
	})
	if len(views) != 2 || views[0].Twice != 2 || views[1].Twice != 4 || meta.Limit != 2 {
		t.Fatalf("views=%+v meta=%+v", views, meta)
	}
}
