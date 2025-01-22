package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"threshAI/internal/personality"
)

func GenerateResponse(userInput string) string {
	// Use absolute path to thresh binary
	cmd := exec.Command("/mnt/c/src/threshAI/bin/thresh", "generate", userInput)
	cmd.Stderr = os.Stderr // Redirect stderr to the terminal for debugging
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error running thresh generate: %v\n", err)
		return "Oops, something went wrong. Let's try again!"
	}

	// Print raw output for debugging
	fmt.Printf("Raw output: %s\n", string(output))

	// Extract the generated text (skip debug logs)
	generatedText := string(output)
	if split := strings.Split(generatedText, "Generated: "); len(split) > 1 {
		generatedText = strings.TrimSpace(split[1])
		fmt.Printf("Extracted text: %s\n", generatedText)
	} else {
		fmt.Printf("No 'Generated:' prefix found in output\n")
	}
	return generatedText
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userInput := r.FormValue("input")
	if userInput == "" {
		http.Error(w, "Missing 'input' parameter", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received user input: %s\n", userInput)

	// Generate response
	response := GenerateResponse(userInput)
	fmt.Printf("Generated response: %s\n", response)

	// Create response JSON with proper escaping
	responseData := struct {
		Response string `json:"response"`
	}{
		Response: response,
	}

	// Encode as JSON with proper escaping
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func greetingHandler(w http.ResponseWriter, r *http.Request) {
	greeting := personality.GetGreeting()

	responseData := struct {
		Greeting string `json:"greeting"`
	}{
		Greeting: greeting,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		fmt.Printf("Error encoding greeting JSON: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Serve static files (HTML, CSS, JS)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// Handle chat requests
	http.HandleFunc("/chat", chatHandler)
	http.HandleFunc("/greeting", greetingHandler)

	// Start server
	fmt.Println("Eidos web server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
