package agent

import "gophkeep/internal/logger"

type RootDeps struct {
	Logger logger.ILogger
}

func (r *RootDeps) NewAgent() IAgent {
	return NewAgent(r.Logger)
}
