/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Set up OpenAI client
		client := openai.NewClient(viper.GetString("openai-key"))

		// Start a chat conversation
		fmt.Println("Starting a chat conversation...")
		chat(client)

	},
}

func init() {
	rootCmd.AddCommand(chatCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// func chat(client *openai.Client) {
// 	ctx := context.Background()

// 	req := openai.ChatCompletionRequest{
// 		Model:     openai.GPT3Dot5Turbo,
// 		MaxTokens: 200,
// 		Messages: []openai.ChatCompletionMessage{
// 			{
// 				Role:    openai.ChatMessageRoleUser,
// 				Content: "Good morning! How are you?",
// 			},
// 		},
// 		Stream: true,
// 	}
// 	stream, err := client.CreateChatCompletionStream(ctx, req)
// 	if err != nil {
// 		fmt.Printf("ChatCompletionStream error: %v\n", err)
// 		return
// 	}
// 	defer stream.Close()

// 	fmt.Printf("Stream response: ")
// 	for {
// 		response, err := stream.Recv()
// 		if errors.Is(err, io.EOF) {
// 			fmt.Println("\nStream finished")
// 			return
// 		}

// 		if err != nil {
// 			fmt.Printf("\nStream error: %v\n", err)
// 			return
// 		}

// 		fmt.Printf(response.Choices[0].Delta.Content)
// 	}
// }

func chat(client *openai.Client) {
	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)

	messages := []openai.ChatCompletionMessage{}

	for {
		fmt.Print("You: ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)

		// Add user message to the list of messages
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: userInput,
		})

		// Call the OpenAI API with the user's message
		req := openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo,
			MaxTokens: 200,
			Messages:  messages,
			Stream:    true,
		}
		stream, err := client.CreateChatCompletionStream(ctx, req)
		if err != nil {
			fmt.Printf("ChatCompletionStream error: %v\n", err)
			return
		}

		// Receive responses from the OpenAI API
		fmt.Print("Assistant: ")
		fullResponse := ""
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				fmt.Println("\nStream finished")
				break
			}

			if err != nil {
				fmt.Printf("\nStream error: %v\n", err)
				break
			}

			fmt.Print(response.Choices[0].Delta.Content)
			fullResponse += response.Choices[0].Delta.Content
		}

		// Add the assistant's response to the list of messages
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: fullResponse,
		})

		stream.Close()
	}
}
