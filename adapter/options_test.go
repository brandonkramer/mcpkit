package adapter

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/brandonkramer/mcpkit/envelope"
)

func TestFromEnvelopeWithOptionsTextModes(t *testing.T) {
	t.Parallel()

	env := envelope.Success("ping").Summary("pong").Build()

	cases := []struct {
		name        string
		mode        TextMode
		wantContent bool
		wantSummary string
		wantJSON    bool
	}{
		{name: "full json", mode: TextFull, wantContent: true, wantJSON: true},
		{name: "summary only", mode: TextSummary, wantContent: true, wantSummary: "pong"},
		{name: "none", mode: TextNone, wantContent: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, _, err := FromEnvelopeWithOptions(&env, Options{Text: tc.mode})
			if err != nil {
				t.Fatal(err)
			}
			if tc.wantContent {
				if len(result.Content) == 0 {
					t.Fatal("expected text content")
				}
				text := result.Content[0].(*sdkmcp.TextContent).Text
				if tc.wantSummary != "" && text != tc.wantSummary {
					t.Fatalf("text=%q", text)
				}
				if tc.wantJSON && !strings.Contains(text, `"summary": "pong"`) {
					t.Fatalf("text=%q", text)
				}
			} else if len(result.Content) != 0 {
				t.Fatalf("content=%v", result.Content)
			}
			if result.StructuredContent == nil {
				t.Fatal("expected structured content")
			}
		})
	}
}

func TestErrorCodeWithOptions(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		mode        TextMode
		wantContent bool
	}{
		{name: "summary text", mode: TextSummary, wantContent: true},
		{name: "no text", mode: TextNone, wantContent: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, envAny, err := ErrorCodeWithOptions("not_found", errors.New("missing"), Options{Text: tc.mode})
			if err != nil {
				t.Fatal(err)
			}
			env := envAny.(envelope.Envelope)
			if env.Diagnostics[0].Code != "not_found" || !result.IsError {
				t.Fatalf("env=%+v isError=%v", env, result.IsError)
			}
			if tc.wantContent {
				if len(result.Content) == 0 {
					t.Fatal("expected text content")
				}
			} else if len(result.Content) != 0 {
				t.Fatalf("content=%v", result.Content)
			}
		})
	}
}

func TestOKWithOptionsStructuredContent(t *testing.T) {
	t.Parallel()
	result, envAny, err := OKWithOptions("demo", "ok", map[string]int{"n": 1}, Options{Text: TextNone})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"ok":true`) {
		t.Fatalf("structured=%s", raw)
	}
	env := envAny.(envelope.Envelope)
	if !env.OK || env.Summary != "ok" {
		t.Fatalf("env=%+v", env)
	}
}
