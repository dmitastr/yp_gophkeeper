package main

import (
	"os"

	"gophkeep/internal/agent/agent"
	"gophkeep/internal/agent/commands"
)

func main() {
	root := commands.NewCmd(commands.RootDeps{
		NewAgent: agent.NewAgent,
	})

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
