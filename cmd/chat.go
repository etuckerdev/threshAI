package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"threshAI/internal/core/memory"

	"github.com/spf13/cobra"
)

var (
	model       string
	interactive bool
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start an interactive chat session",
	Long: `Start an interactive chat session with the AI.
Supports conversation history and context management.`,
	GroupID: "core",
	Run: func(cmd *cobra.Command, args []string) {
		if interactive {
			startInteractiveChat()
		} else if len(args) > 0 {
			handleSingleMessage(strings.Join(args, " "))
		} else {
			fmt.Println("Error: Please provide a message or use --interactive for chat mode")
		}
	},
}

func startInteractiveChat() {
	fmt.Println("Starting interactive chat session (type 'exit' to quit)")
	fmt.Println("----------------------------------------------------")

	scanner := bufio.NewScanner(os.Stdin)
	mem := memory.LoadMemory()
	defer mem.Save()

	for {
		fmt.Print("\nUser > ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if strings.ToLower(strings.TrimSpace(input)) == "exit" {
			fmt.Println("Ending chat session...")
			break
		}

		if input == "" {
			continue
		}

		// Handle the message
		handleMessage(input, mem)
	}
}

func handleSingleMessage(message string) {
	mem := memory.LoadMemory()
	defer mem.Save()
	handleMessage(message, mem)
}

func handleMessage(input string, mem *memory.Memory) {
	// Get relevant context from memory
	context := mem.RetrieveRelevantContext(input)

	// Generate response based on input and context
	response := generateResponse(input, context)

	// Display the response
	fmt.Printf("\nAI > %s\n", response)

	// Store the interaction
	mem.AddInteraction(input, response)
}

func generateResponse(input string, context []memory.Interaction) string {
	// Simple response generation - this can be enhanced with actual LLM integration
	if len(context) > 0 {
		return fmt.Sprintf("I remember our previous conversation about %s. Regarding your current question: %s",
			context[0].UserInput, simpleResponse(input))
	}
	return simpleResponse(input)
}

func simpleResponse(input string) string {
	// Basic response templates - can be extended or replaced with actual LLM
	if strings.Contains(strings.ToLower(input), "hello") {
		return "Hello! How can I assist you today?"
	} else if strings.Contains(strings.ToLower(input), "help") {
		return "I can help you with various tasks. What would you like to know?"
	} else {
		return "I understand you're asking about " + input + ". Could you elaborate?"
	}
}

func init() {
	chatCmd.Flags().StringVarP(&model, "model", "m", "default", "Model to use for chat")
	chatCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Start interactive chat session")

	chatCmd.GroupID = "core"
	rootCmd.AddCommand(chatCmd)
}
