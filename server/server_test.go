package server

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/brandonkramer/mcpkit/adapter"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func connectTestServer(t *testing.T, server *sdkmcp.Server) *sdkmcp.ClientSession {
	t.Helper()
	ct, st := sdkmcp.NewInMemoryTransports()
	if _, err := server.Connect(context.Background(), st, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "test-client", Version: "test"}, nil)
	cs, err := client.Connect(context.Background(), ct, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	return cs
}

func TestNewRegistersTools(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	var gotWorkDir string

	server, absWorkDir, err := New(Config{
		Name:    "test-server",
		Version: "0.0.1",
		WorkDir: workDir,
		Register: func(s *sdkmcp.Server, wd string) {
			gotWorkDir = wd
			sdkmcp.AddTool(s, &sdkmcp.Tool{
				Name: "ping", Description: "Health check",
			}, func(_ context.Context, _ *sdkmcp.CallToolRequest, _ struct{}) (*sdkmcp.CallToolResult, any, error) {
				return adapter.OK("ping", "pong", map[string]string{"status": "ok"})
			})
		},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	wantWorkDir, err := filepath.Abs(workDir)
	if err != nil {
		t.Fatal(err)
	}
	if absWorkDir != wantWorkDir || gotWorkDir != wantWorkDir {
		t.Fatalf("workDir got=%q abs=%q register=%q want=%q", gotWorkDir, absWorkDir, gotWorkDir, wantWorkDir)
	}

	cs := connectTestServer(t, server)
	tools, err := cs.ListTools(context.Background(), &sdkmcp.ListToolsParams{})
	if err != nil {
		t.Fatalf("ListTools: %v", err)
	}
	if len(tools.Tools) != 1 || tools.Tools[0].Name != "ping" {
		t.Fatalf("tools=%+v", tools.Tools)
	}

	result, err := cs.CallTool(context.Background(), &sdkmcp.CallToolParams{Name: "ping"})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if result.IsError {
		t.Fatalf("result=%+v", result)
	}
}
