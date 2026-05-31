package adapter

import (
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/brandonkramer/mcpkit/envelope"
)

func TestFromEnvelopeSuccess(t *testing.T) {
	t.Parallel()
	env := envelope.Success("demo").Summary("ok").Build()
	result, envAny, err := FromEnvelope(&env)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := envAny.(envelope.Envelope)
	if !ok || !got.OK {
		t.Fatalf("env=%v ok=%v", envAny, ok)
	}
	if result.IsError {
		t.Fatal("expected success result")
	}
}

func TestErrorResult(t *testing.T) {
	t.Parallel()
	result, envAny, err := Error(errTest("boom"))
	if err != nil {
		t.Fatal(err)
	}
	env := envAny.(envelope.Envelope)
	if env.OK || !result.IsError {
		t.Fatalf("env=%+v isError=%v", env, result.IsError)
	}
}

type errTest string

func (e errTest) Error() string { return string(e) }

func TestOK(t *testing.T) {
	t.Parallel()
	result, envAny, err := OK("notes.list", "2 note(s)", []string{"a", "b"}, "notes.create")
	if err != nil {
		t.Fatal(err)
	}
	env := envAny.(envelope.Envelope)
	if env.Summary != "2 note(s)" || len(env.Next) != 1 || result.IsError {
		t.Fatalf("env=%+v isError=%v", env, result.IsError)
	}
}

func TestErrorCode(t *testing.T) {
	t.Parallel()
	result, envAny, err := ErrorCode("not_found", errTest("missing"))
	if err != nil {
		t.Fatal(err)
	}
	env := envAny.(envelope.Envelope)
	if len(env.Diagnostics) != 1 || env.Diagnostics[0].Code != "not_found" || !result.IsError {
		t.Fatalf("env=%+v", env)
	}
}

func TestTextSummaryMode(t *testing.T) {
	t.Parallel()
	env := envelope.Success("ping").Summary("pong").Build()
	result, _, err := FromEnvelopeWithOptions(&env, Options{Text: TextSummary})
	if err != nil {
		t.Fatal(err)
	}
	text := result.Content[0].(*sdkmcp.TextContent).Text
	if text != "pong" {
		t.Fatalf("text=%q", text)
	}
}
