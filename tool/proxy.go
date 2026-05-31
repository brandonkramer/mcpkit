package tool

import (
	"context"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/brandonkramer/mcpkit/present"
	"github.com/brandonkramer/mcpkit/summarize"
)

// ProxySpec configures one proxied MCP tool.
type ProxySpec[Args, Result any] struct {
	Name        string
	Description string
	Normalize   func(Args) (Args, error)
	Validate    func(Args) error
	Prepare     func(*Bridge, Args) (Args, summarize.Meta, error)
	Execute     func(context.Context, *Bridge, Args) (Result, error)
	Present     func(Args, Result) (any, summarize.Meta)
}

// AddProxy registers a typed proxy tool on a go-sdk server.
func AddProxy[Args, Result any](server *sdkmcp.Server, b *Bridge, spec ProxySpec[Args, Result]) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        spec.Name,
		Description: spec.Description,
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, args Args) (*sdkmcp.CallToolResult, any, error) {
		var err error
		if spec.Normalize != nil {
			args, err = spec.Normalize(args)
			if err != nil {
				return b.Responder.Error(err)
			}
		}
		if spec.Validate != nil {
			if err := spec.Validate(args); err != nil {
				return b.Responder.Error(err)
			}
		}
		meta := summarize.Meta{}
		if spec.Prepare != nil {
			args, meta, err = spec.Prepare(b, args)
			if err != nil {
				return b.Responder.Error(err)
			}
		}
		var out Result
		switch {
		case spec.Execute != nil:
			out, err = spec.Execute(ctx, b, args)
		default:
			out, err = CallTyped[Args, Result](b, ctx, spec.Name, args)
		}
		if err != nil {
			return b.Responder.Error(err)
		}
		data := any(out)
		if spec.Present != nil {
			data, meta = spec.Present(args, out)
		}
		return b.Responder.OK(spec.Name, data, meta)
	})
}

// AddProxyNoArgs registers a proxy tool with no arguments.
func AddProxyNoArgs[Result any](server *sdkmcp.Server, b *Bridge, name, description string, execute func(context.Context, *Bridge) (Result, error)) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        name,
		Description: description,
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, _ struct{}) (*sdkmcp.CallToolResult, any, error) {
		var out Result
		var err error
		if execute != nil {
			out, err = execute(ctx, b)
		} else {
			out, err = CallNoArgs[Result](b, ctx, name)
		}
		if err != nil {
			return b.Responder.Error(err)
		}
		return b.Responder.OK(name, out, summarize.Meta{})
	})
}

// PresentCappedList wraps a slice in structured tool list data.
func PresentCappedList[T any](items []T, limit int, hintTemplate string) (any, summarize.Meta) {
	data, meta := present.WrapToolList(items, limit, hintTemplate)
	return data, meta
}
