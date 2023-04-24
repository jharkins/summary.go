package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Initiate a chat conversation with GPT-3.",
	Long:  `Initiate a chat conversation with GPT-3 with streaming and turn-based chatting.`,
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

func printChatHelp() {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("To continue (say, if there's a truncation), just press %s.\n", green("return"))
	fmt.Printf("To end the chat, type %s or %s.\n", red("exit"), red("quit"))
	fmt.Printf("You will be offered to %s your transcript when you quit.\n", yellow("save"))

}

func chat(client *openai.Client) {
	printChatHelp()

	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)

	messages := []openai.ChatCompletionMessage{}

	// Yep - lets see about this.
	messages = append(messages, openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleSystem,
		Content: `
		You are summary.go, a chatbot that uses the OpenAI API to create engaging human/AI experiences through text.

		Aid us in our task, be verbose, creative, and kind in your responses.
		`,
	})

	chatTranscript := ""

	for {
		color.Magenta("\nYou: ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)

		if userInput == "exit" || userInput == "quit" {
			break
		}

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
		color.Cyan("\nOpenAI: ")
		fullResponse := ""
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				color.Green("\n\nTransmission Over.\n\n")
				break
			}

			if err != nil {
				color.Red("\nError: ")
				fmt.Printf("\nTransmission error: %v\n", err)
				break
			}

			fmt.Print(response.Choices[0].Delta.Content)
			fullResponse += response.Choices[0].Delta.Content
		}

		chatTranscript += "You:\n " + userInput + "\n\n"
		chatTranscript += "Assistant:\n " + fullResponse + "\n\n\n"

		// Add the assistant's response to the list of messages
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: fullResponse,
		})

		stream.Close()
	}

	saveTranscript(chatTranscript)
}

func saveTranscript(chatTranscript string) {
	fmt.Print("\n\nWould you like to save the chat transcript? [Y/n]: ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)

	if response == "Y" || response == "y" || response == "" {
		usr, err := user.Current()
		if err != nil {
			fmt.Println("Error retrieving user's home directory:", err)
			return
		}

		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := fmt.Sprintf("summary_go_chat_%s.txt", timestamp)
		dir := filepath.Join(usr.HomeDir, "summary_go_chats")
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("Error creating chat directory:", err)
			return
		}

		filePath := filepath.Join(dir, filename)
		err = ioutil.WriteFile(filePath, []byte(chatTranscript), 0644)
		if err != nil {
			fmt.Println("Error saving chat transcript:", err)
			return
		}

		fmt.Printf("Chat transcript saved to: %s\n", filePath)
	} else {
		fmt.Println("Chat transcript not saved.")
	}
}
