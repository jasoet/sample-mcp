package main

import (
	"context"
	"github.com/FreePeak/cortex/pkg/server"
	"github.com/FreePeak/cortex/pkg/tools"
	"log"
	"os"
	//"sample-mcp/config"
	//"sample-mcp/db"
	"sample-mcp/handler"
	//"sample-mcp/ops"
)

func main() {
	logger := log.New(os.Stderr, "[cortex-stdio] ", log.LstdFlags)

	//cfg, err := config.LoadConfig()
	//if err != nil {
	//	logger.Fatalf("Failed to load configuration: %v", err)
	//}
	//logger.Printf("Database configuration loaded: Type=%s, Host=%s, Port=%d, Database=%s",
	//	cfg.Database.DbType, cfg.Database.Host, cfg.Database.Port, cfg.Database.DbName)
	//
	//dbConfig := cfg.Database
	//pool, err := dbConfig.Pool()
	//if err != nil {
	//	logger.Fatalf("Failed to load database: %v", err)
	//}
	//
	//err = db.RunMigrations(pool)
	//if err != nil {
	//	logger.Fatalf("Failed to run migration: %v", err)
	//}
	//
	//_, err = ops.NewQueryOps(ops.WithGormDB(pool))
	//if err != nil {
	//	logger.Fatalf("Failed to initiate query ops: %v", err)
	//}

	mcpServer := server.NewMCPServer("Cortex Stdio Server", "1.0.0", logger)

	echoTool := tools.NewTool("echo",
		tools.WithDescription("Echoes back the input message"),
		tools.WithString("message",
			tools.Description("The message to echo back"),
			tools.Required(),
		),
	)

	var err error

	ctx := context.Background()
	err = mcpServer.AddTool(ctx, echoTool, handler.HandleEcho)
	if err != nil {
		logger.Fatalf("Error adding echo tool: %v", err)
	}

	logger.Printf("Server ready. The following tools are available:\n")
	logger.Printf("- echo\n")

	if err := mcpServer.ServeStdio(); err != nil {
		logger.Printf("Error serving stdio: %v\n", err)
		os.Exit(1)
	}
}
