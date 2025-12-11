package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/firebase/genkit/go/plugins/server"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

var g *genkit.Genkit
var history map[string][]*ai.Message
var toolRefList []ai.ToolRef
var maxTokens int

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()
	history = make(map[string][]*ai.Message)

	//Initialize Genkit with the Google AI plugin
	g = genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{APIKey: os.Getenv("GEMINI_API_KEY")}),
		genkit.WithDefaultModel(os.Getenv("GEMINI_MODEL")),
	)
	//for OpenAI: g := genkit.Init(ctx, genkit.WithPlugins(&openai.OpenAI{APIKey: os.Getenv("OPENAI_API_KEY")}), genkit.WithDefaultModel(os.Getenv("OPENAI_API_MODEL")))

	toolRefList = []ai.ToolRef{
		getGeneralQuestionAnswerTool(g),
		getCurrentWeatherTool(g),
		guessAgeTool(g),
		guessGenderTool(g),
		getRecipeTool(g),
	}

	chatFlowAction := genkit.DefineFlow(g, "chatFlow", chatFlow)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /chat", genkit.Handler(chatFlowAction))
	//curl -X POST http://localhost:3400/chat -H 'Content-Type: application/json' -d '{"data": {"question": "hello", "sessionID": "abc"}}'

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "chat.html")
	})
	fmt.Println("Server starting, from browser visit http://127.0.0.1:3400")
	log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))

}

type ApiInput struct {
	Question  string `json:"question" jsonschema:"description=User question to be answered"`
	SessionID string `json:"sessionID,omitempty" jsonschema:"description=Session identifier for context tracking"`
}

func chatFlow(ctx context.Context, userInput *ApiInput) (string, error) {
	fmt.Println("Human:", userInput.Question)

	userHistory := []*ai.Message{}
	if userInput.SessionID != "" {
		userHistory = history[userInput.SessionID]
	}
	maxTokens, err := strconv.Atoi(os.Getenv("LLM_MAX_OUTPUT_TOKENS_INT"))
	if err != nil {
		maxTokens = 500 // default value
	}
	response, err := genkit.Generate(
		ctx,
		g,
		ai.WithSystem("You are a helpful assistant that provides accurate and concise information."),
		ai.WithMessages(userHistory...),
		ai.WithPrompt(userInput.Question),
		ai.WithTools(toolRefList...),
		ai.WithConfig(&genai.GenerateContentConfig{
			MaxOutputTokens: int32(maxTokens),
		}),
	)

	if err != nil {
		return "LLM generation failed", err
	}
	fmt.Println(string("AI: " + response.Text()))

	saveHistory(userInput.SessionID, userInput.Question, response.Text())
	return response.Text(), nil
}

func saveHistory(sessionID string, userMessage string, botMessage string) {
	if sessionID != "" && botMessage != "" {
		userHistory := history[sessionID]
		maxLength, err := strconv.Atoi(os.Getenv("USER_HISTORY_MAX_LENGTH"))
		if err != nil {
			maxLength = 10 // default value
		}
		if len(userHistory) > maxLength {
			userHistory = userHistory[len(userHistory)-maxLength:]
		}
		history[sessionID] = append(userHistory, &ai.Message{
			Role:    ai.RoleUser,
			Content: []*ai.Part{ai.NewTextPart(userMessage)},
		}, &ai.Message{
			Role:    ai.RoleModel,
			Content: []*ai.Part{ai.NewTextPart(botMessage)},
		})
	}
}
