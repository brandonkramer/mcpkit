// Package server provides stdio MCP server bootstrap helpers built on go-sdk.
//
// HTTP, SSE, and streamable transports are provided by go-sdk directly.
// Build a server with [New], register tools, then serve with sdkmcp.NewStreamableHTTPHandler
// or other go-sdk transports. Mcpkit does not wrap HTTP.
package server
