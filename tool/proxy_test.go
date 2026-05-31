package tool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/brandonkramer/mcpkit/envelope"
	"github.com/brandonkramer/mcpkit/summarize"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type proxyArgs struct {
	Name string `json:"name"`
}

type proxyResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

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

func decodeEnvelope(t *testing.T, result *sdkmcp.CallToolResult) envelope.Envelope {
	t.Helper()
	raw, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatalf("marshal structured: %v", err)
	}
	var env envelope.Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	return env
}

func TestAddProxyBackendCall(t *testing.T) {
	t.Parallel()

	server := sdkmcp.NewServer(&sdkmcp.Implementation{Name: "test", Version: "test"}, nil)
	bridge := &Bridge{
		Backend: FuncBackend(func(_ context.Context, method string, params any) (any, error) {
			if method != "items.create" {
				return nil, fmt.Errorf("unexpected method %q", method)
			}
			args, ok := params.(proxyArgs)
			if !ok {
				return nil, errors.New("bad params")
			}
			return proxyResult{ID: "01", Name: args.Name}, nil
		}),
		Responder: summarize.NewResponder(summarize.Passthrough()),
	}

	AddProxy(server, bridge, ProxySpec[proxyArgs, proxyResult]{
		Name:        "items.create",
		Description: "Create item",
	})

	cs := connectTestServer(t, server)
	result, err := cs.CallTool(context.Background(), &sdkmcp.CallToolParams{
		Name:      "items.create",
		Arguments: map[string]any{"name": "alpha"},
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	env := decodeEnvelope(t, result)
	if !env.OK || env.Tool != "items.create" {
		t.Fatalf("env=%+v", env)
	}
}

func TestAddProxyPipelineHooks(t *testing.T) {
	t.Parallel()

	server := sdkmcp.NewServer(&sdkmcp.Implementation{Name: "test", Version: "test"}, nil)
	bridge := &Bridge{
		Responder: summarize.NewResponder(summarize.Passthrough()),
	}

	AddProxy(server, bridge, ProxySpec[proxyArgs, proxyResult]{
		Name:        "items.create",
		Description: "Create item",
		Normalize: func(args proxyArgs) (proxyArgs, error) {
			if args.Name == "" {
				return args, errors.New("name required")
			}
			args.Name = "N:" + args.Name
			return args, nil
		},
		Validate: func(args proxyArgs) error {
			if args.Name == "N:bad" {
				return errors.New("invalid name")
			}
			return nil
		},
		Prepare: func(_ *Bridge, args proxyArgs) (proxyArgs, summarize.Meta, error) {
			return args, summarize.Meta{Limit: 10}, nil
		},
		Execute: func(_ context.Context, _ *Bridge, args proxyArgs) (proxyResult, error) {
			return proxyResult{ID: "99", Name: args.Name}, nil
		},
		Present: func(_ proxyArgs, out proxyResult) (any, summarize.Meta) {
			return map[string]string{"id": out.ID, "name": out.Name}, summarize.Meta{Limit: 10}
		},
	})

	cs := connectTestServer(t, server)

	t.Run("success", func(t *testing.T) {
		result, err := cs.CallTool(context.Background(), &sdkmcp.CallToolParams{
			Name:      "items.create",
			Arguments: map[string]any{"name": "alpha"},
		})
		if err != nil {
			t.Fatalf("CallTool: %v", err)
		}
		env := decodeEnvelope(t, result)
		if !env.OK {
			t.Fatalf("env=%+v", env)
		}
	})

	t.Run("normalize error", func(t *testing.T) {
		result, err := cs.CallTool(context.Background(), &sdkmcp.CallToolParams{
			Name:      "items.create",
			Arguments: map[string]any{"name": ""},
		})
		if err != nil {
			t.Fatalf("CallTool: %v", err)
		}
		if !result.IsError {
			t.Fatal("expected tool error result")
		}
	})

	t.Run("validate error", func(t *testing.T) {
		result, err := cs.CallTool(context.Background(), &sdkmcp.CallToolParams{
			Name:      "items.create",
			Arguments: map[string]any{"name": "bad"},
		})
		if err != nil {
			t.Fatalf("CallTool: %v", err)
		}
		if !result.IsError {
			t.Fatal("expected tool error result")
		}
	})
}

func TestAddProxyNoArgs(t *testing.T) {
	t.Parallel()

	server := sdkmcp.NewServer(&sdkmcp.Implementation{Name: "test", Version: "test"}, nil)
	bridge := &Bridge{
		Backend: FuncBackend(func(_ context.Context, method string, _ any) (any, error) {
			if method != "items.list" {
				return nil, fmt.Errorf("unexpected method %q", method)
			}
			return []proxyResult{{ID: "01", Name: "alpha"}}, nil
		}),
		Responder: summarize.NewResponder(summarize.Passthrough()),
	}

	AddProxyNoArgs[[]proxyResult](server, bridge, "items.list", "List items", nil)

	cs := connectTestServer(t, server)
	result, err := cs.CallTool(context.Background(), &sdkmcp.CallToolParams{Name: "items.list"})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	env := decodeEnvelope(t, result)
	if !env.OK || env.Tool != "items.list" {
		t.Fatalf("env=%+v", env)
	}
}
