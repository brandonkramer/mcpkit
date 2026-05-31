package envelope_test

import (
	"testing"

	"github.com/brandonkramer/mcpkit/envelope"
)

func TestSuccessExport(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		tool    string
		summary string
	}{
		{name: "sets ok and summary", tool: "ping", summary: "pong"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			env := envelope.Success(tc.tool).Summary(tc.summary).Build()
			if !env.OK || env.Tool != tc.tool || env.Summary != tc.summary {
				t.Fatalf("env=%+v", env)
			}
		})
	}
}

func TestHintExport(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		tool string
		part string
		want string
	}{
		{name: "tool only", tool: "notes.get", want: "notes.get"},
		{name: "with param", tool: "notes.get", part: envelope.HintParam("id", "01"), want: `notes.get id="01"`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var got string
			if tc.part == "" {
				got = envelope.Hint(tc.tool)
			} else {
				got = envelope.Hint(tc.tool, tc.part)
			}
			if got != tc.want {
				t.Fatalf("got=%q want=%q", got, tc.want)
			}
		})
	}
}
