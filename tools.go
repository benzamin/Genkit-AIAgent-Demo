package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"google.golang.org/genai"
)

func getCurrentWeatherTool(g *genkit.Genkit) ai.Tool {
	return genkit.DefineTool(g, "getWeather", "Gets the current weather in a given location",
		func(ctx *ai.ToolContext, input struct {
			Location string `json:"location"`
		}) (string, error) {
			fmt.Println("TOOL CALLED: Weather tool with location:", input.Location)
			return fmt.Sprintf("The current weather in %s is 63°F and sunny.", input.Location), nil
		},
	)
}

func getGeneralQuestionAnswerTool(g *genkit.Genkit) ai.Tool {
	return genkit.DefineTool(g, "generalQuestionAnswer", "Answers any general questions",
		func(ctx *ai.ToolContext, input struct {
			Question string `json:"question"`
		}) (string, error) {
			fmt.Println("TOOL CALLED: General Question Answer tool with question:", input.Question)
			resp, err := genkit.Generate(ctx, g,
				ai.WithPrompt(input.Question),
				ai.WithConfig(&genai.GenerateContentConfig{
					MaxOutputTokens: int32(maxTokens),
				}),
			)
			return resp.Text(), err
		},
	)
}

type ageToolResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

func guessAgeTool(g *genkit.Genkit) ai.Tool {
	return genkit.DefineTool(g, "guessAge", "Guesses the age of a person based on their name",
		func(ctx *ai.ToolContext, input struct {
			Name string `json:"name"`
		}) (int, error) {
			fmt.Println("TOOL CALLED: Age tool  with name:", input.Name)
			resp, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", url.QueryEscape(input.Name)))
			if err != nil {
				return 0, fmt.Errorf("failed to call age API: %w", err)
			}
			defer resp.Body.Close()

			var result ageToolResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return 0, fmt.Errorf("failed to decode response: %w", err)
			}
			age := result.Age
			return age, nil
		},
	)
}

type genderToolResponse struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

func guessGenderTool(g *genkit.Genkit) ai.Tool {
	return genkit.DefineTool(g, "guessGender", "Guesses the gender of a person based on their name",
		func(ctx *ai.ToolContext, input struct {
			Name string `json:"name"`
		}) (string, error) {
			fmt.Println("TOOL CALLED: Gender tool with name:", input.Name)
			resp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", url.QueryEscape(input.Name)))
			if err != nil {
				return "", fmt.Errorf("failed to call gender API: %w", err)
			}
			defer resp.Body.Close()
			var result genderToolResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return "", fmt.Errorf("failed to decode response: %w", err)
			}
			gender := result.Gender + fmt.Sprintf(" (with probability %.2f)", result.Probability)
			return gender, nil
		},
	)
}

// Define input schema
type RecipeInput struct {
	Ingredient          string `json:"ingredient" jsonschema:"description=Main ingredient or cuisine type"`
	DietaryRestrictions string `json:"dietaryRestrictions,omitempty" jsonschema:"description=Any dietary restrictions"`
}

// Define output schema
type Recipe struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	PrepTime     string   `json:"prepTime"`
	CookTime     string   `json:"cookTime"`
	Servings     int      `json:"servings"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	Tips         []string `json:"tips,omitempty"`
}

func getRecipeTool(g *genkit.Genkit) ai.Tool {
	return genkit.DefineTool(g, "getRecipe", "Generates a recipe based on main ingredient and dietary restrictions",
		func(ctx *ai.ToolContext, input RecipeInput) (*Recipe, error) {
			fmt.Println("TOOL CALLED: Recipe tool with ingredient:", input.Ingredient, "and dietary restrictions:", input.DietaryRestrictions)
			recipe, _, err := genkit.GenerateData[Recipe](
				context.Background(), g,
				ai.WithPrompt(fmt.Sprintf(`Create a recipe with the following requirements: Main ingredient: %s Dietary restrictions: %s`, input.Ingredient, input.DietaryRestrictions)),
				ai.WithConfig(&genai.GenerateContentConfig{
					MaxOutputTokens: int32(maxTokens),
				}),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to generate recipe: %w", err)
			}
			return recipe, nil
		},
	)
}
