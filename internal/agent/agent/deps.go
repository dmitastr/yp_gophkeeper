package agent

import "gophkeep/internal/logger"

type RootDeps struct {
	Logger logger.ILogger
}

func (r *RootDeps) NewAgent(bearerToken string) IAgent {
	return NewAgent(r.Logger, bearerToken)
}
