# Genkit AI Agent Demo (Go)

This repository demonstrates how to build an intelligent chat agent using the **Genkit Go** framework. It showcases how to create a conversational AI that maintains context (history) and leverages **tools** to perform specific tasks or retrieve external information, going beyond simple text generation.

## Features

- **Contextual Chat**: Maintains conversation history per session ID, allowing for multi-turn conversations.
- **Tool Integration**: The agent is equipped with several tools to demonstrate function calling capabilities:
    - **Weather Tool**: Retrieves current weather information (Mocked for demo).
    - **Age Guesser**: Predicts age based on a name using `api.agify.io`.
    - **Gender Guesser**: Predicts gender based on a name using `api.genderize.io`.
    - **Recipe Generator**: Generates structured recipe data based on ingredients and dietary restrictions.
    - **General QA**: Handles general queries using the LLM's knowledge.
- **Structured Output**: Demonstrates generating structured data (JSON) for the recipe tool.
- **HTTP Server**: Exposes the chat flow via a REST API.

## Project Structure

- **`main.go`**: The entry point. It initializes Genkit, configures the Google AI plugin, defines the chat flow, manages conversation history, and starts the HTTP server.
- **`tools.go`**: Contains the definitions of the tools available to the agent.
- **`.env`**: Configuration file for API keys and settings.

## Prerequisites

- [Go](https://go.dev/dl/) 1.22 or later.
- A Google AI Studio API Key (for Gemini models). Get one [here](https://aistudio.google.com/).

## Installation & Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/benzamin/Genkit-AIAgent-Demo.git
    cd Genkit-AIAgent-Demo
    ```

2.  **Install dependencies:**
    ```bash
    go mod download
    ```

3.  **Configure Environment Variables:**
    Create a `.env` file in the root directory (or rename a sample if provided) and add your API keys:

    ```env
    # .env
    GEMINI_MODEL=googleai/gemini-2.5-flash
    GOOGLE_API_KEY=your_google_api_key_here
    
    # Optional Configuration
    USER_HISTORY_MAX_LENGTH=10
    LLM_MAX_OUTPUT_TOKENS=500
    ```

## Running the Application

Start the server:

```bash
go run .
```

The server will start on [http://127.0.0.1:3400](http://127.0.0.1:3400).

## Usage

Once the server is running, you can visit [Simple Chat Page](http://127.0.0.1:3400) in your browser to use the simple chat interface.

Alternatively, you can interact with the agent using `curl` or any API client (like Postman).

### Basic Chat
```bash
curl -X POST http://localhost:3400/chat \
     -H 'Content-Type: application/json' \
     -d '{"data": {"question": "Hello, who are you?", "sessionID": "session-1"}}'
```

### Using Tools (Age/Gender Guessing)
The agent will automatically decide to use the `guessAge` or `guessGender` tool based on your question.

```bash
curl -X POST http://localhost:3400/chat \
     -H 'Content-Type: application/json' \
     -d '{"data": {"question": "Guess the age of Benzamin", "sessionID": "session-1"}}'
```

### Using Tools (Recipe Generation)
Ask for a recipe to see structured data generation in action.

```bash
curl -X POST http://localhost:3400/chat \
     -H 'Content-Type: application/json' \
     -d '{"data": {"question": "Give me a recipe for chicken with keto diet", "sessionID": "session-1"}}'
```

## How It Works

1.  **Initialization**: `main.go` initializes the Genkit instance with the Google AI plugin.
2.  **Tool Definition**: Tools are defined in `tools.go` using `genkit.DefineTool`. Some tools make external HTTP requests (like Agify), while others use the LLM itself to generate structured data.
3.  **Flow Definition**: A `chatFlow` is defined to handle incoming requests.
4.  **Processing**:
    - The flow retrieves the conversation history for the given `sessionID`.
    - It calls `genkit.Generate` with the user's prompt, history, and the list of available tools.
    - The LLM decides whether to answer directly or call a tool.
    - If a tool is called, Genkit executes the Go function defined in `tools.go` and feeds the result back to the LLM.
    - The final response is sent back to the user and saved to history.

## Future Improvements / Wishlist

- **Persistence**: Currently, conversation history is stored in an in-memory variable (`history` map). This means all context is lost when the server restarts. Future versions should implement persistence using a database (like PostgreSQL) or a cache (like Redis).
- **Project Structure**: This project uses a flat structure for simplicity. As the project grows, it should be refactored to follow standard Go project layout (e.g., `cmd/`, `internal/`, `pkg/`) and design patterns.
- **Error Handling**: Enhanced error handling and logging.
- **Authentication**: Add authentication for the API endpoints.

## Contributing

Feel free to fork this repository and add more tools! To add a new tool:
1.  Define a new function in `tools.go` that returns `ai.Tool`.
2.  Add the new tool to the `toolRefList` in `main.go`.

## License

[MIT License](LICENSE)
