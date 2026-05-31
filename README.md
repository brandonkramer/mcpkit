# mcpkit

Agent-facing MCP ergonomics built on [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk).

It adds structured agent tool results, optional proxy-tool plumbing, and list/preview helpers.

**Single runtime dependency:** `modelcontextprotocol/go-sdk`

## Install

From [pkg.go.dev](https://pkg.go.dev/github.com/brandonkramer/mcpkit):

```bash
go get github.com/brandonkramer/mcpkit
```

## Development

Lefthook and golangci-lint are pinned in `go.mod` as **tools** (dev-only). Install git hooks once per clone:

```bash
make install-hooks
```

Hooks and `make lint` use `go tool` binaries from `go.mod`. Pre-commit lints staged `.go` files; pre-push runs `./scripts/check.sh`. CI runs the same checks.

```bash
make check    # full local CI script
make test
make lint
```

---

| Level | Packages | When |
| --- | --- | --- |
| **Minimal** | `envelope`, `adapter` | Hand-written tool handlers; you build the envelope yourself |
| **Proxy** | + `summarize`, `tool` | Tools call a backend (RPC, DB, HTTP); shared normalize→present→envelope pipeline |
| **Full** | + `present`, `server` | Capped list payloads, previews, stdio bootstrap |

---

## Minimal: envelope + adapter

Register tools directly with go-sdk 

```go
sdkmcp.AddTool(server, &sdkmcp.Tool{
    Name: "ping", Description: "Health check",
}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, _ struct{}) (*sdkmcp.CallToolResult, any, error) {
    env := envelope.Success("ping").
        Summary("pong").
        Data(map[string]string{"status": "ok"}).
        Build()
    return adapter.FromEnvelope(&env)
})
```

Shorthand helpers (same envelope shape, less boilerplate):

```go
return adapter.OK("notes.list", "2 note(s)", data, "notes.create title=\"...\"")
return adapter.ErrorCode("not_found", err)
```

On error without a code:

```go
return adapter.Error(err)
```

**Text mode** for clients that only read `content` text (structured envelope still in `structuredContent`):

```go
adapter.OKWithOptions("notes.get", "note 01", data, adapter.Options{Text: adapter.TextSummary})
// TextFull (default): full envelope JSON in text
// TextSummary: summary line only
// TextNone: no text content
```

No `summarize.Registry`, no `tool.Bridge`, no `server` helper required.

---

## Proxy: backend + registry

When tools proxy to a service and share summarization logic:

```go
reg := summarize.Passthrough() // or your own Registry; Passthrough() for day-one proxy wiring

bridge := &tool.Bridge{
    Backend:   tool.FuncBackend(myBackend.Call),
    Responder: summarize.NewResponder(reg),
}

tool.AddProxy(server, bridge, tool.ProxySpec[CreateArgs, Item]{
    Name: "items.create", Description: "Create an item",
    Prepare: func(_ *tool.Bridge, args CreateArgs) (CreateArgs, summarize.Meta, error) {
        return tool.PrepareWorkDir(defaultWorkDir, args, setWorkDir, getWorkDir)
    },
})
```

`WorkDir` is app-specific — use `tool.ResolveWorkDir` / `tool.PrepareWorkDir` in your `prepare` hooks, not on `Bridge`.

List tool summarization:

```go
summary, next := present.SummarizeToolList(data, meta, summarize.ToolListSummary[Item]{
    CountLabel: "item(s)",
    Next:       []string{"items.create title=\"...\""},
    PageHint:   "narrow filter=... or use CLI for full list",
    NextFor: func(item Item) []string {
        return []string{envelope.Hint("items.get", envelope.HintParam("id", item.ID))}
    },
})
```

---

## `summarize.Meta`

`Meta` passes context from `prepare` / `present` hooks to summarizers.

| Field | Purpose |
| --- | --- |
| `Limit` | List cap applied during present (used by `SummarizeToolList`) |
| `SinceSeq` | Optional paging cursor (ignore if unused) |
| `RunID` | Optional entity id context (ignore if unused) |
| `Extra` | Open bag: `meta.Put("table", "notes")`, read with `meta.String("table")` |

Only use the fields your app needs. A replay/poll proxy might use `Limit`, `SinceSeq`, and `RunID`; a notes MCP might only set `Limit` or use `Extra` alone.

---

## Transport

| Transport | Use |
| --- | --- |
| **Stdio** | `server.ServeStdio(server.Config{...})` |
| **HTTP / streamable** | go-sdk — mcpkit does not wrap HTTP |

Build once with `server.New`, then serve over any go-sdk transport:

```go
server, _, err := mcpserver.New(mcpserver.Config{
    Name: "notes", Version: "1.0.0",
    Register: func(s *sdkmcp.Server, _ string) { registerNotesTools(s, db) },
})

// Stdio
_ = server.Run(ctx, &sdkmcp.StdioTransport{})

// HTTP (go-sdk) — see modelcontextprotocol/go-sdk StreamableHTTPHandler
handler := sdkmcp.NewStreamableHTTPHandler(func(_ *http.Request) *sdkmcp.Server {
    return server
}, nil)
_ = http.ListenAndServe(":8080", handler)
```

---

## Packages

| Package | Purpose |
| --- | --- |
| `envelope` | `ok`, `summary`, `subjects`, `data`, `next_actions`, diagnostics |
| `adapter` | `OK`, `ErrorCode`, text modes; map envelopes to `CallToolResult` |
| `summarize` | `Passthrough`, `ToolSummaries`, `ToolListItems`, `Meta` |
| `present` | `ToolListData`, `WrapToolList`, `ToolListItems`, `MapCapList`, `RedactText` |
| `tool` | Proxy pipeline, `ResolveWorkDir`, `PrepareWorkDir`, `NormalizeField`, `NormalizeID` |
| `server` | Stdio bootstrap via `ServeStdio` / `New` |
