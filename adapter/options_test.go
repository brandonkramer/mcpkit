package adapter

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/brandonkramer/mcpkit/envelope"
)

func TestTextFullMode(t *testing.T) {
	t.Parallel()
	env := envelope.Success("ping").Summary("pong").Build()
	result, _, err := FromEnvelopeWithOptions(&env, Options{Text: TextFull})
	if err != nil {
		t.Fatal(err)
	}
	text := result.Content[0].(*sdkmcp.TextContent).Text
	if !strings.Contains(text, `"summary": "pong"`) {
		t.Fatalf("text=%q", text)
	}
}

func TestTextNoneMode(t *testing.T) {
	t.Parallel()
	env := envelope.Success("ping").Summary("pong").Build()
	result, _, err := FromEnvelopeWithOptions(&env, Options{Text: TextNone})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Content) != 0 {
		t.Fatalf("content=%v", result.Content)
	}
	if result.StructuredContent == nil {
		t.Fatal("expected structured content")
	}
}

func TestErrorCodeWithOptionsTextSummary(t *testing.T) {
	t.Parallel()
	errTest := errors.New("missing")
	result, envAny, err := ErrorCodeWithOptions("not_found", errTest, Options{Text: TextSummary})
	if err != nil {
		t.Fatal(err)
	}
	env := envAny.(envelope.Envelope)
	if env.Diagnostics[0].Code != "not_found" || !result.IsError {
		t.Fatalf("env=%+v isError=%v", env, result.IsError)
	}
	text := result.Content[0].(*sdkmcp.TextContent).Text
	if text == "" {
		t.Fatal("expected summary text")
	}
}

func TestErrorCodeWithOptionsTextNone(t *testing.T) {
	t.Parallel()
	result, _, err := ErrorCodeWithOptions("bad", errors.New("nope"), Options{Text: TextNone})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Content) != 0 || !result.IsError {
		t.Fatalf("content=%v isError=%v", result.Content, result.IsError)
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
