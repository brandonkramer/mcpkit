package summarize

import (
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/brandonkramer/mcpkit/adapter"
	"github.com/brandonkramer/mcpkit/envelope"
)

// Responder builds MCP tool results from envelopes using a registry.
type Responder struct {
	Registry *Registry
	Adapter  adapter.Options
}

// NewResponder returns a responder with the given registry or Default when nil.
func NewResponder(reg *Registry) *Responder {
	if reg == nil {
		reg = Default()
	}
	return &Responder{Registry: reg, Adapter: adapter.DefaultOptions}
}

// OK builds a success CallToolResult for tool name and payload data.
func (r *Responder) OK(tool string, data any, meta Meta) (*sdkmcp.CallToolResult, any, error) {
	summary, next := r.Registry.Summarize(tool, data, meta)
	subjects := r.Registry.Subjects(tool, data, meta)
	env := envelope.Success(tool).
		Summary(summary).
		Subjects(subjects...).
		Data(data).
		NextHints(next...).
		Build()
	return adapter.FromEnvelopeWithOptions(&env, r.Adapter)
}

// Error builds an error CallToolResult.
func (r *Responder) Error(err error) (*sdkmcp.CallToolResult, any, error) {
	return adapter.ErrorCodeWithOptions("", err, r.Adapter)
}
