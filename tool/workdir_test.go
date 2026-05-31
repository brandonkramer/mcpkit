package tool

import (
	"testing"
)

type workDirArgs struct {
	WorkDir string `json:"work_dir"`
}

func TestResolveWorkDir(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		defaultDir string
		arg        string
		want       string
	}{
		{name: "uses default when arg empty", defaultDir: "/srv", arg: "", want: "/srv"},
		{name: "prefers arg", defaultDir: "/srv", arg: "/tmp", want: "/tmp"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := ResolveWorkDir(tc.defaultDir, tc.arg); got != tc.want {
				t.Fatalf("ResolveWorkDir(%q, %q) = %q; want %q", tc.defaultDir, tc.arg, got, tc.want)
			}
		})
	}
}

func TestPrepareWorkDir(t *testing.T) {
	t.Parallel()

	args := workDirArgs{}
	got, meta, err := PrepareWorkDir("/default", args,
		func(a *workDirArgs, wd string) { a.WorkDir = wd },
		func(a workDirArgs) string { return a.WorkDir },
	)
	if err != nil {
		t.Fatalf("PrepareWorkDir: %v", err)
	}
	if got.WorkDir != "/default" {
		t.Fatalf("WorkDir=%q", got.WorkDir)
	}
	if gotMeta := meta; gotMeta.Limit != 0 || gotMeta.Extra != nil {
		t.Fatalf("meta=%+v", meta)
	}

	got, _, err = PrepareWorkDir("/default", workDirArgs{WorkDir: "/override"},
		func(a *workDirArgs, wd string) { a.WorkDir = wd },
		func(a workDirArgs) string { return a.WorkDir },
	)
	if err != nil {
		t.Fatalf("PrepareWorkDir override: %v", err)
	}
	if got.WorkDir != "/override" {
		t.Fatalf("WorkDir=%q", got.WorkDir)
	}
}
