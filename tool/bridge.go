package tool

import (
	"context"
	"fmt"

	"github.com/brandonkramer/mcpkit/summarize"
)

// Bridge carries shared proxy-tool state.
type Bridge struct {
	Backend   Backend
	Responder *summarize.Responder
}

// CallTyped executes a backend call with typed result decoding.
func CallTyped[Args, Result any](b *Bridge, ctx context.Context, method string, args Args) (Result, error) {
	raw, err := b.Backend.Call(ctx, method, args)
	if err != nil {
		var zero Result
		return zero, err
	}
	out, ok := raw.(Result)
	if !ok {
		var zero Result
		return zero, fmt.Errorf("unexpected result type for %s", method)
	}
	return out, nil
}

// CallNoArgs executes a backend call with no params.
func CallNoArgs[Result any](b *Bridge, ctx context.Context, method string) (Result, error) {
	raw, err := b.Backend.Call(ctx, method, nil)
	if err != nil {
		var zero Result
		return zero, err
	}
	out, ok := raw.(Result)
	if !ok {
		var zero Result
		return zero, fmt.Errorf("unexpected result type for %s", method)
	}
	return out, nil
}
