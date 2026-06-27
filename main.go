package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()

	// 1. Initialize the client. It automatically picks up the GEMINI_API_KEY env var.
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to create GenAI client: %v", err)
	}

	// 2. Start a persistent chat session
	chat, err := client.Chats.Create(ctx, "gemini-2.5-flash", nil, nil)
	if err != nil {
		log.Fatalf("Failed to create chat session: %v", err)
	}

	fmt.Println("🤖 Gemini Go CLI Initialized. (Type 'exit' or 'quit' to stop)")
	fmt.Println("-----------------------------------------------------------")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\nYou > ")
		if !scanner.Scan() {
			break
		}

		userInput := strings.TrimSpace(scanner.Text())
		if userInput == "" {
			continue
		}

		if userInput == "exit" || userInput == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		fmt.Print("Gemini > ")

		// 3. Create the text part payload
		contentPart := genai.Part{Text: userInput}

		// 4. Send message stream. It returns a single streaming function.
		stream := chat.SendMessageStream(ctx, contentPart)

		// 5. Execute the streaming function by passing a callback
		var streamErr error
		stream(func(resp *genai.GenerateContentResponse, err error) bool {
			if err != nil {
				streamErr = err
				return false // stop the stream on error
			}

			// Safely print text chunk using the .Text() helper
			if resp != nil {
				fmt.Print(resp.Text())
			}
			return true // continue listening for chunks
		})

		if streamErr != nil {
			fmt.Printf("\nError during streaming: %v\n", streamErr)
		}

		fmt.Println() // Add a trailing newline after the response finishes streaming
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading standard input: %v", err)
	}
}
