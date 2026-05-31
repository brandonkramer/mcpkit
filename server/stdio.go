package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Config describes a stdio MCP server instance.
type Config struct {
	Name         string
	Version      string
	Instructions string
	WorkDir      string
	Register     func(server *sdkmcp.Server, workDir string)
}

// New creates a go-sdk MCP server from config.
func New(cfg Config) (*sdkmcp.Server, string, error) {
	workDir := cfg.WorkDir
	if workDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, "", err
		}
		workDir = cwd
	}
	abs, err := filepath.Abs(workDir)
	if err != nil {
		return nil, "", fmt.Errorf("work_dir: %w", err)
	}
	server := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    cfg.Name,
		Version: cfg.Version,
	}, &sdkmcp.ServerOptions{Instructions: cfg.Instructions})
	if cfg.Register != nil {
		cfg.Register(server, abs)
	}
	return server, abs, nil
}

// ServeStdio runs the MCP server over stdin/stdout until interrupted.
func ServeStdio(cfg Config) error {
	server, _, err := New(cfg)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return server.Run(ctx, &sdkmcp.StdioTransport{})
}
