package tool

import "context"

// Backend executes tool calls against an application service.
type Backend interface {
	Call(ctx context.Context, method string, params any) (any, error)
}

// FuncBackend adapts a function to Backend.
type FuncBackend func(ctx context.Context, method string, params any) (any, error)

// Call implements Backend.
func (f FuncBackend) Call(ctx context.Context, method string, params any) (any, error) {
	return f(ctx, method, params)
}
