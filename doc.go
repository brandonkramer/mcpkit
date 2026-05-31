// Module mcpkit provides agent-facing MCP ergonomics on top of github.com/modelcontextprotocol/go-sdk.
//
// Integration levels:
//   - minimal: envelope + adapter only (hand-written tool handlers)
//   - proxy: tool.AddProxy pipeline with Backend + summarize.Registry
//   - full: present list shaping + server.ServeStdio or go-sdk HTTP transports
//
// Subpackages:
//   - envelope: domain-neutral tool result JSON envelopes
//   - adapter: envelope to CallToolResult mapping
//   - summarize: pluggable summary/subject hooks and responder
//   - present: ToolListData, previews, capped list views
//   - tool: proxy-tool registration pipeline
//   - server: stdio server bootstrap (HTTP via go-sdk)
//
// mcpkit intentionally does not reimplement MCP transport or protocol handling.
package mcpkit
