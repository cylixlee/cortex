package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cylixlee/cortex/pkg/agent"
	"github.com/cylixlee/cortex/pkg/llm"
)

type consoleUI struct{}

func newConsoleUI() *consoleUI {
	return &consoleUI{}
}

func (c *consoleUI) OnResponse(resp llm.ChatResponse) {
	fmt.Print(resp.Delta)
}

func (c *consoleUI) OnError(err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
}

func (c *consoleUI) Run(agent *agent.Agent) {
	fmt.Println("Welcome to Cortex CLI! Type 'exit' or 'quit' to exit.")
	fmt.Println()

	ctx := context.Background()
	sessionID := "default"

	for {
		fmt.Print("> ")
		
		var input string
		fmt.Scanln(&input)

		if input == "exit" || input == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		if input == "" {
			continue
		}

		agent.Chat(ctx, sessionID, input)
		fmt.Println()
	}
}
