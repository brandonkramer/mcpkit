package adapter_test

import (
	"errors"
	"testing"

	"github.com/brandonkramer/mcpkit/adapter"
	"github.com/brandonkramer/mcpkit/envelope"
)

func TestOKExport(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		tool    string
		summary string
		wantErr bool
	}{
		{name: "success", tool: "notes.list", summary: "2 note(s)"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, envAny, err := adapter.OK(tc.tool, tc.summary, []string{"a"}, "notes.create")
			if err != nil {
				t.Fatal(err)
			}
			env := envAny.(envelope.Envelope)
			if env.Tool != tc.tool || env.Summary != tc.summary || result.IsError {
				t.Fatalf("env=%+v isError=%v", env, result.IsError)
			}
		})
	}
}

func TestErrorCodeExport(t *testing.T) {
	t.Parallel()

	result, envAny, err := adapter.ErrorCode("not_found", errors.New("missing"))
	if err != nil {
		t.Fatal(err)
	}
	env := envAny.(envelope.Envelope)
	if env.Diagnostics[0].Code != "not_found" || !result.IsError {
		t.Fatalf("env=%+v", env)
	}
}
