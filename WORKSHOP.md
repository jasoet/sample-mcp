# Model Context Protocol (MCP) Workshop

This workshop will guide you through understanding and implementing the Model Context Protocol (MCP) in a Go application. By the end of this workshop, you'll be able to create your own MCP server with custom tools that can interact with databases and other services.

## What is Model Context Protocol (MCP)?

Model Context Protocol (MCP) is a protocol for creating tools that can be exposed through a server. Each tool has specific inputs, parameters, and handlers. MCP allows you to create a standardized interface for your services, making them easily accessible and usable.

## Workshop Outline

### Step 1: Understanding MCP Basics

The Model Context Protocol (MCP) is an open standard developed by Anthropic to standardize how AI applications interact with external tools, data sources, and systems. It provides a structured way for AI models to access and utilize external functionalities and data, enhancing their capabilities and contextual understanding.

#### Core Components of MCP

1. **MCP Server**  
   The MCP Server acts as a bridge between AI models and external systems. It exposes various capabilities to clients, including:
   - **Tools**: Functions that AI models can invoke to perform specific actions, such as calling APIs or executing commands.
   - **Resources**: Data sources that provide information to AI models without performing actions, akin to RESTful GET endpoints.
   - **Prompts**: Predefined templates that guide AI models in generating responses or performing tasks.

   MCP Servers can be implemented in various programming languages and communicate with clients through different transports, such as standard input/output (stdio) or HTTP with Server-Sent Events (SSE).

   **Example**: An MCP server might expose a weather tool that accepts a city name and returns current weather conditions, or a database tool that can query customer information.

2. **Tools**  
   Tools are model-controlled functions that AI models can call to execute specific tasks. They are defined with clear input and output schemas, allowing AI models to understand how to interact with them. Each tool has:
   - A unique identifier
   - A description of its purpose
   - A schema defining required and optional parameters
   - A handler function that executes the tool's logic

   **Example**: A calendar tool might have parameters for `date`, `time`, and `description`, and return a confirmation when an event is created.

   ```json
   {
     "name": "create_calendar_event",
     "description": "Creates a new calendar event",
     "parameters": {
       "date": "string (YYYY-MM-DD)",
       "time": "string (HH:MM)",
       "description": "string"
     }
   }
   ```

3. **Resources**  
   Resources are application-controlled data sources that AI models can access to retrieve information. Unlike tools, resources do not perform actions but provide data that can inform the AI model's responses. Examples include accessing a user's contact list or retrieving documents from a database.

   **Example**: A contacts resource might provide access to a user's address book, allowing the AI to reference contact information when composing emails.

4. **Prompts**  
   Prompts are user-controlled templates that help standardize interactions between AI models and external systems. They can be used to format queries, provide context, or guide the AI model's behavior in specific scenarios.

   **Example**: A customer service prompt template might include placeholders for customer information and previous interaction history to help the AI generate appropriate responses.

#### MCP Communication Workflow

The interaction between AI models and external systems via MCP follows a structured workflow:

1. **Initialization**  
   An MCP client establishes a connection with an MCP server. This involves a handshake process where the client and server negotiate capabilities and protocol versions.

2. **Capability Discovery**  
   The client queries the server to discover available tools, resources, and prompts. The server responds with a list of its capabilities, allowing the client to understand what functionalities are accessible.

3. **Context Provisioning**  
   Based on the discovered capabilities, the client can present resources and prompts to the AI model, enriching its context and guiding its interactions.

4. **Invocation and Execution**  
   When the AI model determines that a specific tool is needed to fulfill a task, it instructs the client to invoke the corresponding tool on the server. The server executes the tool and returns the result to the client.

5. **Response Integration**  
   The client integrates the server's response into the AI model's context, enabling it to generate informed and context-aware outputs.

By understanding these core components and workflows, participants in the Model Context Protocol Workshop will gain a solid foundation for building AI applications that can effectively interact with a wide range of external tools and data sources.

A short presentation about MCP is available here: [MCP Presentation Slides](https://docs.google.com/presentation/d/1Fuq5TYie_VHCAiuKFy3OaGzoVRzxCXa1riX8dWoZwEY/edit?usp=sharing)

### Step 2: Setting Up Your Environment

In this step, you'll prepare your development environment to build and run an MCP server using Go. We'll also configure supporting services such as PostgreSQL and an LLM provider.

#### Prerequisites

Before proceeding, make sure you have the following:

- Go SDK (1.21 or later recommended)
- IDE (Visual Studio Code or your preferred Go-compatible editor)
- Docker & Docker Compose
- PostgreSQL (we'll use Docker to simplify this)
- An LLM API key (such as Gemini or OpenAI-compatible endpoint)
- Dive (AI chat client that supports MCP): [https://github.com/OpenAgentPlatform/Dive](https://github.com/OpenAgentPlatform/Dive)

#### Installing Required Dependencies

Install Go:

```bash
# macOS
brew install go

# Ubuntu
sudo apt install golang-go

# Windows
# Download from https://go.dev/dl/
```

Install Docker:

```bash
# macOS
# Download from https://docs.docker.com/desktop/mac/install/

# Ubuntu
sudo apt install docker.io docker-compose

# Windows
# Download from https://docs.docker.com/desktop/windows/install/
```

Install VS Code (optional but recommended):
- [https://code.visualstudio.com/Download](https://code.visualstudio.com/Download)

#### Setting Up PostgreSQL with Docker

The docker-compose file is available in the `compose` directory:

```bash
# Start PostgreSQL using Docker Compose
cd compose
docker-compose up -d
```

This will start a PostgreSQL instance with the following configuration:
- Host: localhost
- Port: 5432
- Username: postgres
- Password: postgres
- Database: postgres

You can verify the database is running with:
```bash
docker ps
# or
docker-compose ps
```

#### Configuring LLM API Keys

To use an LLM provider with your MCP server, you'll need to obtain and configure API keys:

1. **For OpenAI API**:
   - Sign up at [OpenAI Platform](https://platform.openai.com/)
   - Create an API key in your account dashboard
   - Store your API key, we will use it later

2. **For Google Gemini API**:
   - Sign up for Google AI Studio at [Google AI Studio](https://makersuite.google.com/)
   - Create an API key in your account
   - Store your API key, we will use it later

3. **For Anthropic Claude API**:
   - Sign up at [Anthropic Console](https://console.anthropic.com/)
   - Create an API key in your account settings
   - Store your API key, we will use it later

You can add these environment variables to your shell profile (`.bashrc`, `.zshrc`, etc.) for persistence.

#### Download and Install Dive AI Chat Client

Download and install the Dive AI Chat client, which supports MCP:
- [Dive Releases](https://github.com/OpenAgentPlatform/Dive/releases)

After installation, configure Dive to use your LLM API key in the settings panel.

### Step 3: Technical Deep Dive into MCP Server

In this step, we'll explore the technical aspects of the Model Context Protocol Server and understand how it communicates with clients and LLMs.

#### MCP Server Architecture

The MCP Server is the core component that handles communication between AI clients and tools. It processes JSON-based requests and responses, manages tool registration, and executes tool handlers when invoked.

#### Communication Flow

The MCP protocol follows a specific communication pattern:

1. **Tool Discovery Phase**:
   ```
   ┌────────────┐                                      ┌────────────┐
   │            │  JSON-RPC (tool list request)        │            │
   │ MCP Client │ ───────────────────────────────────> | MCP Server │
   │            │                                      │            │
   │            │  JSON-RPC (tool list)                │            │
   │            │ <─────────────────────────────────── │            │
   └────────────┘                                      └────────────┘
   ```

2. **Tool Execution Phase**:
   ```
   ┌────────────┐                                      ┌────────────┐
   │            │  Send Prompt + List of Tools         │            │
   │ MCP Client │ ───────────────────────────────────> |   LLM      │
   │            │                                      │            │
   │            │  tools execution suggestion          │            │
   │            │ <─────────────────────────────────── │            │
   └────────────┘                                      └────────────┘

   ┌────────────┐                                      ┌────────────┐
   │            │  JSON-RPC  (tool execution)          │            │
   │ MCP Client │ ───────────────────────────────────> | MCP Server │
   │            │                                      │            │
   │            │  JSON-RPC  (execution result)        │            │
   │            │ <─────────────────────────────────── │            │
   └────────────┘                                      └────────────┘
   ┌────────────┐                                      ┌────────────┐
   │            │  execution result                    │            │
   │ MCP Client │ ───────────────────────────────────> |   LLM      │
   │            │                                      │            │
   │            │  prompt result                       │            │
   │            │ <─────────────────────────────────── │            │
   └────────────┘                                      └────────────┘
   ```
#### Example MCP Request/Response JSON

1. **Tool Discovery Request**:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/list"
}
```

2. **Tool Discovery Response**:
```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": {
    "tools": [
      {
        "description": "Echoes back the input message",
        "inputSchema": {
          "properties": {
            "message": {
              "description": "The message to echo back",
              "type": "string"
            }
          },
          "required": ["message"],
          "type": "object"
        },
        "name": "echo"
      }
    ]
  }
}
```

3. **Tool Execution Request**:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "echo",
    "parameters": {
      "message": "Hello, MCP Server!"
    }
  }
}
```

4. **Tool Execution Response**:
```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "text": "[1746772565] Hello, MCP Server!",
        "type": "text"
      }
    ]
  }
}
```

#### Starting the MCP Server

Using cortex library (github.com/FreePeak/cortex/) to setup MCP server.

```go
package main

import (
    "context"
    "log"
    "os"
    "github.com/FreePeak/cortex/pkg/server"
    "github.com/FreePeak/cortex/pkg/tools"
)

func main() {
    // Create a logger
    logger := log.New(os.Stderr, "[cortex-stdio] ", log.LstdFlags)

    // Create a new MCP server
    mcpServer := server.NewMCPServer("My MCP Server", "1.0.0", logger)

    // Register tools (will be covered in Step 4)

    // Start the server using stdio transport
    if err := mcpServer.ServeStdio(); err != nil {
        logger.Printf("Error serving stdio: %v\n", err)
        os.Exit(1)
    }
}
```

This technical foundation will help you understand how the MCP server processes requests and communicates with clients and LLMs, preparing you for implementing your own tools in the next steps.

### Step 4: Implementing Your First Tool - Echo

In this step, we'll implement a simple "echo" tool that takes a message as input and returns it with a timestamp. This will demonstrate the basic structure of an MCP tool and how to integrate it with the MCP server.

#### 1. Creating a Simple Echo Tool

First, let's define our echo tool in the main.go file. We'll use the `tools` package from the cortex library to create a new tool with a name and description:

```go
echoTool := tools.NewTool("echo",
    tools.WithDescription("Echoes back the input message"),
    tools.WithString("message",
        tools.Description("The message to echo back"),
        tools.Required(),
    ),
)
```

This code creates a new tool named "echo" with a description and a single parameter:
- The tool name is "echo"
- The description explains what the tool does: "Echoes back the input message"
- It has one parameter called "message" which is a string and is required

#### 2. Implementing the Handler Function

Next, we need to implement the handler function that will process the tool call. Create a new file called `handler/echo.go` with the following content:

```go
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

    // Extract and validate the message parameter
    message, ok := request.Parameters["message"].(string)
    if !ok {
        return nil, fmt.Errorf("missing or invalid 'message' parameter")
    }

    // Add a timestamp to the message
    timestamp := fmt.Sprintf("%d", time.Now().Unix())
    responseMessage := fmt.Sprintf("[%s] %s", timestamp, message)

    // Return the response in the format expected by the MCP server
    return map[string]interface{}{
        "content": []map[string]interface{}{
            {
                "type": "text",
                "text": responseMessage,
            },
        },
    }, nil
}
```

This handler function:
1. Extracts the "message" parameter from the request
2. Validates that the parameter exists and is a string
3. Adds a Unix timestamp to the message
4. Returns the response in the format expected by the MCP server

#### 3. Registering the Tool with the MCP Server

Finally, we need to register our tool with the MCP server in the main.go file:

```go
ctx := context.Background()
err = mcpServer.AddTool(ctx, echoTool, handler.HandleEcho)
if err != nil {
    logger.Fatalf("Error adding echo tool: %v", err)
}

logger.Printf("Server ready. The following tools are available:\n")
logger.Printf("- echo\n")
```

This code:
1. Creates a background context
2. Adds the echo tool to the MCP server, associating it with the HandleEcho handler function
3. Logs an error if the tool couldn't be added
4. Logs a message indicating that the server is ready and which tools are available

#### 4. Complete Implementation

Here's the complete implementation of the main.go file with the echo tool:

```go
package main

import (
    "context"
    "github.com/FreePeak/cortex/pkg/server"
    "github.com/FreePeak/cortex/pkg/tools"
    "log"
    "os"
    "sample-mcp/handler"
)

func main() {
    // Create a logger
    logger := log.New(os.Stderr, "[cortex-stdio] ", log.LstdFlags)

    // Create a new MCP server
    mcpServer := server.NewMCPServer("Cortex Stdio Server", "1.0.0", logger)

    // Create the echo tool
    echoTool := tools.NewTool("echo",
        tools.WithDescription("Echoes back the input message"),
        tools.WithString("message",
            tools.Description("The message to echo back"),
            tools.Required(),
        ),
    )

    var err error

    // Register the echo tool with the MCP server
    ctx := context.Background()
    err = mcpServer.AddTool(ctx, echoTool, handler.HandleEcho)
    if err != nil {
        logger.Fatalf("Error adding echo tool: %v", err)
    }

    logger.Printf("Server ready. The following tools are available:\n")
    logger.Printf("- echo\n")

    // Start the server using stdio transport
    if err := mcpServer.ServeStdio(); err != nil {
        logger.Printf("Error serving stdio: %v\n", err)
        os.Exit(1)
    }
}
```

With this implementation, your MCP server now has a functional echo tool that can receive messages and echo them back with a timestamp. In the next step, we'll test this tool to ensure it works as expected.

### Step 5: Testing Your Echo Tool

In this step, we'll build the MCP server application and test the echo tool we created in Step 4. We'll use command-line tools to send JSON-RPC requests to the server and examine the responses.

#### Building the Application

First, you need to build the application. You have two options:

1. Using standard Go build:

```bash
# For Unix/Linux/macOS
go build -o mcp-server ./main.go

# For Windows
go build -o mcp-server.exe ./main.go
```

2. Using Mage build tools (if you have Mage installed):

```bash
mage compileBuild
```

This will create an executable file named `mcp-server` (or `mcp-server.exe` on Windows) in your project directory.

#### Testing the Echo Tool

Once you've built the application, you can test it by sending JSON-RPC requests through standard input. We'll use the `echo` command to pipe JSON requests to our server.

##### 1. Listing Available Tools

To list all available tools, send a tools/list request:

```bash
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list"}' | ./mcp-server
```

You should receive a response similar to:

```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": {
    "tools": [
      {
        "description": "Echoes back the input message",
        "inputSchema": {
          "properties": {
            "message": {
              "description": "The message to echo back",
              "type": "string"
            }
          },
          "required": ["message"],
          "type": "object"
        },
        "name": "echo"
      }
    ]
  }
}
```

This response shows that our server has one tool available: the `echo` tool we created in Step 4.

##### 2. Executing the Echo Tool

To execute the echo tool, send a tools/call request with the tool name and parameters:

```bash
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "echo", "parameters": {"message": "Hello, MCP Server!"}}}' | ./mcp-server
```

You should receive a response similar to:

```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "text": "[1746772565] Hello, MCP Server!",
        "type": "text"
      }
    ]
  }
}
```

The response includes our message with a timestamp prefix, exactly as we implemented in the `HandleEcho` function.

#### Understanding the Response Format

The response from the MCP server follows the JSON-RPC 2.0 specification:

1. **id**: Matches the id from the request, allowing you to correlate requests and responses.
2. **jsonrpc**: Always "2.0", indicating the JSON-RPC protocol version.
3. **result**: Contains the actual response data from the tool for successful responses.
   - For the `tools/list` method, it contains a list of available tools with their descriptions and input schemas.
   - For the `tools/call` method, it contains the output from the tool, which in our case is a content array with a text element.
4. **error**: Present only in error responses, contains information about the error.
   - Contains **code** (a number that indicates the error type)
   - Contains **message** (a short description of the error)
   - May contain **data** (additional information about the error)

Example of an error response:
```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "error": {
    "code": -32602,
    "message": "Invalid params",
    "data": "The 'message' parameter is required"
  }
}
```

Standard JSON-RPC 2.0 error codes:
- **-32700**: Parse error - Invalid JSON was received
- **-32600**: Invalid Request - The JSON sent is not a valid Request object
- **-32601**: Method not found - The method does not exist / is not available
- **-32602**: Invalid params - Invalid method parameter(s)
- **-32603**: Internal error - Internal JSON-RPC error
- **-32000 to -32099**: Server error - Reserved for implementation-defined server errors

#### Batch Requests

JSON-RPC 2.0 also supports batch requests, allowing multiple requests to be sent in a single HTTP request. This can improve performance by reducing network overhead. Here's an example of a batch request:

```json
[
  {"jsonrpc": "2.0", "id": 1, "method": "tools/list"},
  {"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "echo", "parameters": {"message": "Hello"}}}
]
```

The server would respond with an array of responses:

```json
[
  {
    "id": 1,
    "jsonrpc": "2.0",
    "result": {
      "tools": [
        {
          "description": "Echoes back the input message",
          "inputSchema": {
            "properties": {
              "message": {
                "description": "The message to echo back",
                "type": "string"
              }
            },
            "required": ["message"],
            "type": "object"
          },
          "name": "echo"
        }
      ]
    }
  },
  {
    "id": 2,
    "jsonrpc": "2.0",
    "result": {
      "content": [
        {
          "text": "[1746772565] Hello",
          "type": "text"
        }
      ]
    }
  }
]
```

Note: The MCP server implementation in this workshop may not support batch requests, but it's a standard feature of JSON-RPC 2.0 that you might want to implement in a production environment.

#### Debugging Common Issues

If you encounter issues when testing your MCP server, here are some common problems and solutions:

1. **Server not starting**: Make sure you've built the application correctly and the executable file exists.

2. **Permission denied**: If you're on Unix/Linux/macOS, you might need to make the executable file executable:
   ```bash
   chmod +x ./mcp-server
   ```

3. **Invalid JSON**: Make sure your JSON requests are properly formatted. A common issue is missing quotes around property names or string values.

4. **Missing parameters**: If you get an error about missing or invalid parameters, check that you're providing all required parameters with the correct types.

5. **Server crashes**: Check the error messages in the console. The server logs errors to stderr, which should help you identify the issue.

#### Next Steps

Now that you've successfully built and tested your echo tool, you're ready to move on to more complex tools that interact with databases and other services. In the next step, we'll set up database integration to create tools that can store and retrieve data.

### Step 6: Test MCP Server using Dive 

In this step, we'll test our MCP server using the Dive chat client. Dive is an AI chat client that supports the Model Context Protocol, allowing us to interact with our MCP server through a user-friendly interface.

#### Building the MCP Server

Before we can use our MCP server with Dive, we need to build it:

```bash
# For Unix/Linux/macOS
go build -o mcp-server ./main.go

# For Windows
go build -o mcp-server.exe ./main.go

# Or using Mage (if installed)
mage compileBuild
```

Make sure the executable has the proper permissions:

```bash
chmod +x ./mcp-server
```

#### Registering the MCP Server in Dive

1. Open the Dive chat client that you installed in Step 2.

2. Go to Tools Management(MCP).

3. Click on "Add" and use the following configuration:

```json
{
  "mcpServers": {
    "echoAccount": {
      "transport": "stdio",
      "enabled": true,
      "command": "<binary location>",
      "args": [],
      "env": {},
      "url": null
    }
  }
}
```

Replace `<binary location>` with the full path to your `mcp-server` executable. For example:
- On macOS/Linux: `/Users/username/path/to/sample-mcp/mcp-server`
- On Windows: `C:\Users\username\path\to\sample-mcp\mcp-server.exe`

4. Click "Save" to register the MCP server.

#### Validating the Configuration

To validate that the MCP server is correctly registered and working:

1. In the Dive chat interface, start a new conversation.

2. Type a message like "I want to use the echo tool."

3. Dive should detect that you want to use a tool and suggest the echo tool from your MCP server.

4. When prompted, provide a message to echo, such as "Hello, MCP Server!"

5. The MCP server should process your request and return the message, which will be displayed in the chat.

If everything is working correctly, you should see a response similar to:

```
Hello, MCP Server!
```

#### Troubleshooting

If you encounter issues:

1. **MCP Server not found**: Make sure the path to the binary is correct and the file exists.

2. **Permission denied**: Ensure the binary has execution permissions.

3. **Tool not showing up**: Check that the MCP server is enabled in Dive settings.

4. **Error in tool execution**: Look at the Dive logs for error messages. You might also see error messages in the terminal if you're running the MCP server manually.

5. **Connection issues**: Make sure the transport type is set to "stdio" and that the command path is correct.

By successfully testing your MCP server with Dive, you've verified that your implementation works with a real AI chat client. This is an important step before moving on to more complex tools and integrations.


### Step 7: Database Integration

In this step, we'll implement a database integration tool that allows users to query the database for accounts, categories, and transactions. We'll use the existing `QueryOps` code to connect to the database and retrieve data.

#### 1. Understanding the Database Structure

Our application uses three main entities:

1. **Account** - Represents a financial account
   - Fields: AccountID, Name, AccountType, CreatedAt, UpdatedAt

2. **Category** - Represents a transaction category
   - Fields: CategoryID, Name, CategoryType, CreatedAt, UpdatedAt

3. **Transaction** - Represents a financial transaction
   - Fields: TransactionID, AccountID, CategoryID, Amount, TransactionDate, Description, CreatedAt, UpdatedAt

The `QueryOps` class provides methods to query these entities from the database.

#### 2. Implementing the DatabaseHandler

Now, let's implement a `DatabaseHandler` similar to the `EchoHandler` but using the existing `QueryOps` code to connect to the database. Create a new file called `handler/database.go` with the following content:

```go
package handler

import (
	"context"
	"fmt"
	"github.com/FreePeak/cortex/pkg/server"
	"log"
	"sample-mcp/ops"
	"strconv"
)

// HandleGetAllAccounts handles requests to get all accounts
func HandleGetAllAccounts(ctx context.Context, request server.ToolCallRequest) (interface{}, error) {
	log.Printf("Handling get all accounts tool call with name: %s", request.Name)

	// Create a new QueryOps instance
	// In a real application, you would inject this dependency
	queryOps, err := ops.NewQueryOps(ops.WithGormDB(nil)) // Replace nil with actual DB connection
	if err != nil {
		return nil, fmt.Errorf("failed to create query ops: %v", err)
	}

	accounts, err := queryOps.GetAllAccounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %v", err)
	}

	// Return the response in the format expected by the MCP server
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Query result: %v", accounts),
			},
		},
	}, nil
}

// HandleGetAccountById handles requests to get an account by ID
func HandleGetAccountById(ctx context.Context, request server.ToolCallRequest) (interface{}, error) {
	log.Printf("Handling get account by ID tool call with name: %s", request.Name)

	// Extract account ID parameter
	accountIDStr, ok := request.Parameters["account_id"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'account_id' parameter")
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID format: %v", err)
	}

	// Create a new QueryOps instance
	// In a real application, you would inject this dependency
	queryOps, err := ops.NewQueryOps(ops.WithGormDB(nil)) // Replace nil with actual DB connection
	if err != nil {
		return nil, fmt.Errorf("failed to create query ops: %v", err)
	}

	account, err := queryOps.GetAccountByID(ctx, uint(accountID))
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %v", err)
	}

	// Return the response in the format expected by the MCP server
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Query result: %v", account),
			},
		},
	}, nil
}

// HandleGetAllCategories handles requests to get all categories
func HandleGetAllCategories(ctx context.Context, request server.ToolCallRequest) (interface{}, error) {
	log.Printf("Handling get all categories tool call with name: %s", request.Name)

	// Create a new QueryOps instance
	// In a real application, you would inject this dependency
	queryOps, err := ops.NewQueryOps(ops.WithGormDB(nil)) // Replace nil with actual DB connection
	if err != nil {
		return nil, fmt.Errorf("failed to create query ops: %v", err)
	}

	categories, err := queryOps.GetAllCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %v", err)
	}

	// Return the response in the format expected by the MCP server
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Query result: %v", categories),
			},
		},
	}, nil
}

// HandleGetTransactionsByAccount handles requests to get transactions by account ID
func HandleGetTransactionsByAccount(ctx context.Context, request server.ToolCallRequest) (interface{}, error) {
	log.Printf("Handling get transactions by account tool call with name: %s", request.Name)

	// Extract account ID parameter
	accountIDStr, ok := request.Parameters["account_id"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'account_id' parameter")
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID format: %v", err)
	}

	// Create a new QueryOps instance
	// In a real application, you would inject this dependency
	queryOps, err := ops.NewQueryOps(ops.WithGormDB(nil)) // Replace nil with actual DB connection
	if err != nil {
		return nil, fmt.Errorf("failed to create query ops: %v", err)
	}

	transactions, err := queryOps.GetTransactionsByAccountID(ctx, uint(accountID))
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %v", err)
	}

	// Return the response in the format expected by the MCP server
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Query result: %v", transactions),
			},
		},
	}, nil
}
```

#### 3. Registering the Database Tool with the MCP Server

Next, we need to register our database tool with the MCP server in the main.go file:

```go
// Create the get all accounts tool
getAllAccountsTool := tools.NewTool("get_all_accounts",
	tools.WithDescription("Retrieves all accounts from the database"),
)

// Register the get all accounts tool with the MCP server
err = mcpServer.AddTool(ctx, getAllAccountsTool, handler.HandleGetAllAccounts)
if err != nil {
	logger.Fatalf("Error adding get all accounts tool: %v", err)
}

// Create the get account by ID tool
getAccountByIdTool := tools.NewTool("get_account_by_id",
	tools.WithDescription("Retrieves an account by its ID"),
	tools.WithString("account_id",
		tools.Description("The ID of the account to retrieve"),
		tools.Required(),
	),
)

// Register the get account by ID tool with the MCP server
err = mcpServer.AddTool(ctx, getAccountByIdTool, handler.HandleGetAccountById)
if err != nil {
	logger.Fatalf("Error adding get account by ID tool: %v", err)
}

// Create the get all categories tool
getAllCategoriesTool := tools.NewTool("get_all_categories",
	tools.WithDescription("Retrieves all categories from the database"),
)

// Register the get all categories tool with the MCP server
err = mcpServer.AddTool(ctx, getAllCategoriesTool, handler.HandleGetAllCategories)
if err != nil {
	logger.Fatalf("Error adding get all categories tool: %v", err)
}

// Create the get transactions by account tool
getTransactionsByAccountTool := tools.NewTool("get_transactions_by_account",
	tools.WithDescription("Retrieves all transactions for a specific account"),
	tools.WithString("account_id",
		tools.Description("The ID of the account to retrieve transactions for"),
		tools.Required(),
	),
)

// Register the get transactions by account tool with the MCP server
err = mcpServer.AddTool(ctx, getTransactionsByAccountTool, handler.HandleGetTransactionsByAccount)
if err != nil {
	logger.Fatalf("Error adding get transactions by account tool: %v", err)
}

logger.Printf("- get_all_accounts\n")
logger.Printf("- get_account_by_id\n")
logger.Printf("- get_all_categories\n")
logger.Printf("- get_transactions_by_account\n")
```

#### 4. Using the Database Tools

Once implemented, you can use the database tools to query financial data from the database. Here are some example queries:

1. **Get All Accounts**:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "get_all_accounts",
    "parameters": {}
  }
}
```

2. **Get Account by ID**:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "get_account_by_id",
    "parameters": {
      "account_id": "1"
    }
  }
}
```

3. **Get All Categories**:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "get_all_categories",
    "parameters": {}
  }
}
```

4. **Get Transactions by Account**:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "get_transactions_by_account",
    "parameters": {
      "account_id": "1"
    }
  }
}
```

This implementation demonstrates how to create a database integration tool that uses the existing `QueryOps` code to connect to the database and retrieve data. In a real application, you would need to properly initialize the database connection and handle errors more robustly.


## Resources

- [Model Context Protocol Documentation](https://modelcontextprotocol.io/introduction)
- [Magefile - Build Tool](https://magefile.org/)
- [Cortex MCP Documentation](https://github.com/FreePeak/cortex)
- [Go Documentation](https://golang.org/doc/)
- [GORM Documentation](https://gorm.io/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## Next Steps

After completing this workshop, you can:
- Build more complex tools
- Integrate MCP with other services and protocols
- Share your tools with the community

Happy coding!
