package handler

import (
	"context"
	"github.com/FreePeak/cortex/pkg/server"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleEcho_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	testMessage := "Hello, world!"
	request := server.ToolCallRequest{
		Name: "echo",
		Parameters: map[string]interface{}{
			"message": testMessage,
		},
	}

	// Act
	response, err := HandleEcho(ctx, request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Check response structure
	responseMap, ok := response.(map[string]interface{})
	assert.True(t, ok, "Response should be a map")

	content, ok := responseMap["content"].([]map[string]interface{})
	assert.True(t, ok, "Response should have content array")
	assert.Len(t, content, 1, "Content should have one item")

	item := content[0]
	assert.Equal(t, "text", item["type"], "Content item should have type 'text'")

	// Check that the response text contains the timestamp and original message
	responseText, ok := item["text"].(string)
	assert.True(t, ok, "Content item should have text field")

	// The timestamp is dynamic, so we just check that the message is included
	assert.Contains(t, responseText, testMessage, "Response should contain the original message")

	// Verify timestamp format: [timestamp] message
	assert.Regexp(t, `^\[\d+\] Hello, world!$`, responseText, "Response should have timestamp format")
}

func TestHandleEcho_MissingMessage(t *testing.T) {
	// Arrange
	ctx := context.Background()
	request := server.ToolCallRequest{
		Name:       "echo",
		Parameters: map[string]interface{}{},
	}

	// Act
	response, err := HandleEcho(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "missing or invalid 'message' parameter")
}

func TestHandleEcho_InvalidMessageType(t *testing.T) {
	// Arrange
	ctx := context.Background()
	request := server.ToolCallRequest{
		Name: "echo",
		Parameters: map[string]interface{}{
			"message": 123, // Not a string
		},
	}

	// Act
	response, err := HandleEcho(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "missing or invalid 'message' parameter")
}
