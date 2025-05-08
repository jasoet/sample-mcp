package main

import (
	"context"
	"fmt"
	"github.com/FreePeak/cortex/pkg/server"
	"github.com/FreePeak/cortex/pkg/tools"
	"log"
	"os"
	"sample-mcp/config"
	"sample-mcp/handler"
)

func main() {
	logger := log.New(os.Stderr, "[cortex-stdio] ", log.LstdFlags)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}
	logger.Printf("Database configuration loaded: Type=%s, Host=%s, Port=%d, Database=%s",
		cfg.Database.DbType, cfg.Database.Host, cfg.Database.Port, cfg.Database.DbName)

	mcpServer := server.NewMCPServer("Cortex Stdio Server", "1.0.0", logger)

	echoTool := tools.NewTool("echo",
		tools.WithDescription("Echoes back the input message"),
		tools.WithString("message",
			tools.Description("The message to echo back"),
			tools.Required(),
		),
	)

	ctx := context.Background()
	err = mcpServer.AddTool(ctx, echoTool, handler.HandleEcho)
	if err != nil {
		logger.Fatalf("Error adding echo tool: %v", err)
	}

	_, _ = fmt.Fprintf(os.Stderr, "Server ready. The following tools are available:\n")
	_, _ = fmt.Fprintf(os.Stderr, "- echo\n")

	if err := mcpServer.ServeStdio(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error serving stdio: %v\n", err)
		os.Exit(1)
	}
}
