package main

import (
	"gophkeep/internal/agent/agent"
	"gophkeep/internal/agent/commands"
	"gophkeep/internal/logger"
)

func main() {
	log := logger.NewLogger()
	root := commands.NewCmd(agent.RootDeps{
		Logger: log,
	})

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
