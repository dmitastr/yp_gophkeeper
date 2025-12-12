package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gophkeep/internal/agent/agent"
)

func NewPasswordsCmd(deps RootDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "passwords",
		Long:  "Add new password",
		Short: "Add new password",
		RunE: func(cmd *cobra.Command, args []string) error {
			login := viper.GetString("login")
			password := viper.GetString("password")
			address := viper.GetString("address")
			key := viper.GetString("key")
			token := viper.GetString("token")

			connParams := &agent.ConnParams{Address: address, Key: key, Token: token}

			if err := deps.NewAgent().AddPassword(login, password, connParams); err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringP("login", "l", "", "account's login to save")
	cmd.Flags().StringP("password", "p", "", "account's password to save")

	_ = viper.BindPFlag("login", cmd.Flags().Lookup("login"))
	_ = viper.BindPFlag("password", cmd.Flags().Lookup("password"))

	// Bind to environment variable AUTH_USERNAME
	viper.SetEnvPrefix("AUTH") // ENV prefix: AUTH_
	_ = viper.BindEnv("password")
	_ = viper.BindEnv("login")

	return cmd
}
