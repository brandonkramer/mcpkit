package tool

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type bridgeArgs struct {
	ItemID string
}

type bridgeResult struct {
	ID string
}

func TestCallTyped(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("backend down")

	cases := []struct {
		name    string
		backend Backend
		method  string
		wantID  string
		wantErr string
		isSent  bool
	}{
		{
			name: "success",
			backend: FuncBackend(func(_ context.Context, method string, params any) (any, error) {
				args, ok := params.(bridgeArgs)
				if !ok || method != "items.get" {
					return nil, errors.New("unexpected call")
				}
				return bridgeResult{ID: args.ItemID}, nil
			}),
			method: "items.get",
			wantID: "01",
		},
		{
			name: "wraps backend error",
			backend: FuncBackend(func(context.Context, string, any) (any, error) {
				return nil, sentinel
			}),
			method:  "items.get",
			wantErr: "call items.get:",
			isSent:  true,
		},
		{
			name: "unexpected result type",
			backend: FuncBackend(func(context.Context, string, any) (any, error) {
				return "not-a-struct", nil
			}),
			method:  "items.get",
			wantErr: "unexpected result type",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			b := &Bridge{Backend: tc.backend}
			got, err := CallTyped[bridgeArgs, bridgeResult](b, context.Background(), tc.method, bridgeArgs{ItemID: "01"})
			if tc.wantErr != "" {
				if err == nil {
					t.Fatal("expected error")
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("err=%v", err)
				}
				if tc.isSent && !errors.Is(err, sentinel) {
					t.Fatalf("errors.Is failed: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("CallTyped: %v", err)
			}
			if got.ID != tc.wantID {
				t.Fatalf("got=%+v", got)
			}
		})
	}
}

func TestCallNoArgs(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("list failed")
	b := &Bridge{
		Backend: FuncBackend(func(_ context.Context, method string, _ any) (any, error) {
			if method != "items.list" {
				return nil, errors.New("wrong method")
			}
			return []bridgeResult{{ID: "01"}}, nil
		}),
	}

	got, err := CallNoArgs[[]bridgeResult](b, context.Background(), "items.list")
	if err != nil {
		t.Fatalf("CallNoArgs: %v", err)
	}
	if len(got) != 1 || got[0].ID != "01" {
		t.Fatalf("got=%+v", got)
	}

	b.Backend = FuncBackend(func(context.Context, string, any) (any, error) {
		return nil, sentinel
	})
	_, err = CallNoArgs[[]bridgeResult](b, context.Background(), "items.list")
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected wrapped sentinel: %v", err)
	}
}
