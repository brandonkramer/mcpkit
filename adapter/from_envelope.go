package adapter

import (
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/brandonkramer/mcpkit/envelope"
)

// FromEnvelope converts an envelope into an MCP CallToolResult.
func FromEnvelope(env *envelope.Envelope) (*sdkmcp.CallToolResult, any, error) {
	return FromEnvelopeWithOptions(env, DefaultOptions)
}

// Error converts an error into an MCP error tool result.
func Error(err error) (*sdkmcp.CallToolResult, any, error) {
	return ErrorCode("", err)
}
