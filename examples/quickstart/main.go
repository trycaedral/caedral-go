package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/caedral/caedral-go"
)

func loadRootEnv() {
	candidates := []string{
		filepath.Join("..", "..", ".env"),
		filepath.Join("..", "..", "..", ".env"),
	}
	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
		return
	}
}

func main() {
	loadRootEnv()

	apiKey := os.Getenv("CAEDRAL_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("CAEDRAL_TEST_API_KEY")
	}
	if apiKey == "" {
		log.Fatal("set CAEDRAL_API_KEY or CAEDRAL_TEST_API_KEY")
	}

	baseURL := os.Getenv("CAEDRAL_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:5001"
	}

	client, err := caedral.NewClient(apiKey, caedral.WithBaseURL(baseURL))
	if err != nil {
		log.Fatalf("new client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	models, err := client.Models.List(ctx)
	if err != nil {
		log.Fatalf("models.list: %v", err)
	}

	ids := make([]string, 0, len(models.Data))
	for _, model := range models.Data {
		ids = append(ids, model.ID)
	}
	fmt.Println("Models:", ids)

	content := "Say hello in one short sentence."
	completion, err := client.Chat.Completions.Create(ctx, caedral.ChatCompletionRequest{
		Model: "caedral-base",
		Messages: []caedral.ChatMessage{
			{Role: "user", Content: &content},
		},
	})
	if err != nil {
		log.Fatalf("chat completion: %v", err)
	}

	assistant, _ := completion.Choices[0].Message["content"].(string)
	fmt.Println("Assistant:", assistant)
}
