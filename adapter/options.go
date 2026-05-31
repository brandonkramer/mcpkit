package adapter

import (
	"fmt"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/brandonkramer/mcpkit/envelope"
)

// TextMode controls how much of the envelope appears in text content.
type TextMode int

const (
	// TextFull marshals the full envelope JSON into text content (default).
	TextFull TextMode = iota
	// TextSummary puts only the summary line in text content.
	TextSummary
	// TextNone omits text content; structured content still carries the envelope.
	TextNone
)

// Options configures adapter output.
type Options struct {
	Text TextMode
}

// DefaultOptions is the default adapter configuration.
var DefaultOptions = Options{Text: TextFull}

// FromEnvelopeWithOptions converts an envelope into an MCP CallToolResult.
func FromEnvelopeWithOptions(env *envelope.Envelope, opts Options) (*sdkmcp.CallToolResult, any, error) {
	text, err := formatText(env, opts)
	if err != nil {
		return Error(err)
	}
	result := &sdkmcp.CallToolResult{
		StructuredContent: *env,
		IsError:           !env.OK,
	}
	if text != "" {
		result.Content = []sdkmcp.Content{&sdkmcp.TextContent{Text: text}}
	}
	return result, *env, nil
}

func formatText(env *envelope.Envelope, opts Options) (string, error) {
	switch opts.Text {
	case TextNone:
		return "", nil
	case TextSummary:
		if env.Summary != "" {
			return env.Summary, nil
		}
		if !env.OK && env.Error != "" {
			return env.Error, nil
		}
		return env.Tool, nil
	default:
		return envelope.MarshalText(env)
	}
}

// OK builds a success CallToolResult from tool name, summary, data, and optional next hints.
func OK(tool, summary string, data any, next ...string) (*sdkmcp.CallToolResult, any, error) {
	return OKWithOptions(tool, summary, data, DefaultOptions, next...)
}

// OKWithOptions builds a success CallToolResult with adapter options.
func OKWithOptions(tool, summary string, data any, opts Options, next ...string) (*sdkmcp.CallToolResult, any, error) {
	env := envelope.Success(tool).
		Summary(summary).
		Data(data).
		NextHints(next...).
		Build()
	return FromEnvelopeWithOptions(&env, opts)
}

// ErrorCode builds an error CallToolResult with a stable code and message.
func ErrorCode(code string, err error) (*sdkmcp.CallToolResult, any, error) {
	return ErrorCodeWithOptions(code, err, DefaultOptions)
}

// ErrorCodeWithOptions builds an error CallToolResult with adapter options.
func ErrorCodeWithOptions(code string, err error, opts Options) (*sdkmcp.CallToolResult, any, error) {
	env := envelope.Error("request failed", err)
	if code != "" {
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		env.Diagnostics = append(env.Diagnostics, envelope.Diagnostic{
			Code:     code,
			Message:  msg,
			Severity: "error",
		})
	}
	text, encErr := formatText(&env, opts)
	if encErr != nil {
		text = fmt.Sprintf("error: %s", err.Error())
	}
	result := &sdkmcp.CallToolResult{
		StructuredContent: env,
		IsError:           true,
	}
	if text != "" {
		result.Content = []sdkmcp.Content{&sdkmcp.TextContent{Text: text}}
	}
	return result, env, nil
}
