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

2. **Tools**  
   Tools are model-controlled functions that AI models can call to execute specific tasks. They are defined with clear input and output schemas, allowing AI models to understand how to interact with them. For example, a tool might allow an AI model to fetch weather data or create a calendar event.

3. **Resources**  
   Resources are application-controlled data sources that AI models can access to retrieve information. Unlike tools, resources do not perform actions but provide data that can inform the AI model's responses. Examples include accessing a user's contact list or retrieving documents from a database.

4. **Prompts**  
   Prompts are user-controlled templates that help standardize interactions between AI models and external systems. They can be used to format queries, provide context, or guide the AI model's behavior in specific scenarios.

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

## Step 2: Setting Up Your Environment

In this step, you'll prepare your development environment to build and run an MCP server using Go. We'll also configure supporting services such as PostgreSQL and an LLM provider.

---

### 🧰 Prerequisites

Before proceeding, make sure you have the following:

- Go SDK (1.21 or later recommended)
- IDE (Visual Studio Code or your preferred Go-compatible editor)
- Docker & Docker Compose
- PostgreSQL (we'll use Docker to simplify this)
- An LLM API key (such as Gemini or OpenAI-compatible endpoint)
- Dive (AI chat client that supports MCP): [https://github.com/OpenAgentPlatform/Dive](https://github.com/OpenAgentPlatform/Dive)

---

### 🛠️ Installing Required Dependencies

Install Go:

- macOS: `brew install go`
- Ubuntu: `sudo apt install golang-go`
- Windows: Download from [https://go.dev/dl/](https://go.dev/dl/)

Install Docker:

- macOS: [https://docs.docker.com/desktop/mac/install/](https://docs.docker.com/desktop/mac/install/)
- Ubuntu: `sudo apt install docker.io docker-compose`
- Windows: [https://docs.docker.com/desktop/windows/install/](https://docs.docker.com/desktop/windows/install/)

Install VS Code (optional but recommended):

- [https://code.visualstudio.com/Download](https://code.visualstudio.com/Download)

---

### 🐘 Setting Up PostgreSQL with Docker
docker-compose file is available on `compose` directory


##### Download and install Dive AI Chat client
https://github.com/OpenAgentPlatform/Dive/releases


### Step 3: Creating Your First MCP Server

- Initializing an MCP server
- Understanding server configuration options
- Starting the server

### Step 4: Implementing Your First Tool - Echo

- Creating a simple echo tool
- Defining tool parameters
- Implementing the handler function
- Registering the tool with the MCP server

### Step 5: Testing Your Echo Tool

- Making requests to your echo tool
- Understanding the response format
- Debugging common issues

### Step 6: Database Integration

- Setting up database connections
- Creating repositories for data access
- Implementing database operations

### Step 7: Creating a Database-Backed Tool

- Designing a tool that interacts with the database
- Implementing the handler with database operations
- Error handling and response formatting

### Step 8: Advanced MCP Features

- Tool dependencies
- Context propagation
- Authentication and authorization
- Rate limiting and throttling

### Step 9: Best Practices

- Structuring your MCP application
- Error handling strategies
- Testing MCP tools
- Documentation

### Step 10: Building a Complete Application

- Putting it all together
- Creating a multi-tool MCP server
- Deploying your MCP application

## Hands-on Exercises

Throughout this workshop, you'll complete the following exercises:

1. Create a basic MCP server with an echo tool
2. Implement a tool that retrieves data from a database
3. Create a tool that performs CRUD operations
4. Build a complex tool that aggregates data from multiple sources
5. Deploy your MCP server as a standalone application

## Resources

- [MCP Documentation](https://github.com/FreePeak/cortex)
- [Go Documentation](https://golang.org/doc/)
- [GORM Documentation](https://gorm.io/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## Next Steps

After completing this workshop, you can:
- Contribute to the MCP project
- Build more complex tools
- Integrate MCP with other services and protocols
- Share your tools with the community

Happy coding!
