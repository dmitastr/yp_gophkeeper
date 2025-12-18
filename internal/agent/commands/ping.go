package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gophkeep/internal/agent/agent"
)

type PingCmd struct {
	Cmd *cobra.Command
}

func NewPingCmd(deps agent.RootDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ping",
		Long:  "Ping service",
		Short: "ping service",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := &agent.ConnParams{
				Token:   viper.GetString("token"),
				Address: viper.GetString("address"),
				Key:     viper.GetString("key"),
			}

			if err := deps.NewAgent().Ping(params); err != nil {
				return fmt.Errorf("ping service error: %w", err)
			}
			return nil
		},
	}
	return cmd
}
