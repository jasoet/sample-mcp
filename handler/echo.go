package handler

import (
	"context"
	"fmt"
	"github.com/FreePeak/cortex/pkg/server"
	"log"
	"time"
)

func HandleEcho(ctx context.Context, request server.ToolCallRequest) (interface{}, error) {
	log.Printf("Handling echo tool call with name: %s", request.Name)

	message, ok := request.Parameters["message"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'message' parameter")
	}

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	responseMessage := fmt.Sprintf("[%s] %s", timestamp, message)

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": responseMessage,
			},
		},
	}, nil
}
