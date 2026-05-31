package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/brandonkramer/mcpkit/adapter"
	"github.com/brandonkramer/mcpkit/envelope"
	"github.com/brandonkramer/mcpkit/present"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	os.Exit(run())
}

func run() int {
	store := NewStore(Note{ID: "01", Title: "Welcome", Body: "Hello from mcpkit examples/notes."})
	server := sdkmcp.NewServer(&sdkmcp.Implementation{Name: "notes", Version: "0.1.0"}, &sdkmcp.ServerOptions{
		Instructions: "Minimal notes MCP using mcpkit envelope + adapter.",
	})
	registerNotesTools(server, store)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := server.Run(ctx, &sdkmcp.StdioTransport{}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func registerNotesTools(server *sdkmcp.Server, store *Store) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name: "notes.list", Description: "List notes",
	}, func(_ context.Context, _ *sdkmcp.CallToolRequest, _ struct{}) (*sdkmcp.CallToolResult, any, error) {
		notes := store.List()
		data, _ := present.WrapToolList(notes, 20, "%d notes total")
		return adapter.OK("notes.list", fmt.Sprintf("%d note(s)", len(notes)), data, "notes.create title=\"...\"")
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name: "notes.get", Description: "Get one note by id",
	}, func(_ context.Context, _ *sdkmcp.CallToolRequest, args struct {
		ID string `json:"id"`
	}) (*sdkmcp.CallToolResult, any, error) {
		note, err := store.Get(args.ID)
		if err != nil {
			return adapter.ErrorCode("not_found", err)
		}
		prev, body := present.RedactText(note.Body, 500)
		return adapter.OKWithOptions("notes.get", "note "+note.ID, map[string]any{
			"id": note.ID, "title": note.Title,
			"body_preview": prev.Preview, "truncated": prev.Truncated, "body": body,
		}, adapter.Options{Text: adapter.TextSummary},
			envelope.Hint("notes.update", envelope.HintParam("id", note.ID)),
		)
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name: "notes.create", Description: "Create a note",
	}, func(_ context.Context, _ *sdkmcp.CallToolRequest, args struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}) (*sdkmcp.CallToolResult, any, error) {
		if args.Title == "" {
			return adapter.ErrorCode("invalid_argument", fmt.Errorf("title is required"))
		}
		note := store.Create(args.Title, args.Body)
		return adapter.OK("notes.create", "created note "+note.ID, note,
			envelope.Hint("notes.get", envelope.HintParam("id", note.ID)),
		)
	})
}
