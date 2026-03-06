package main

import (
	"github.com/cylixlee/cortex/pkg/agent"
)

func main() {
	llmClient := newLLMClient()
	sessionStore := newSessionStore()
	consoleUI := newConsoleUI()

	a := agent.NewAgent(llmClient, sessionStore, consoleUI)
	consoleUI.Run(a)
}
