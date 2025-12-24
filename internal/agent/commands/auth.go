package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
	"gophkeep/internal/agent/agent"
)

func NewAuthCmd(deps agent.RootDeps) *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Long:  "RegisterUser via login and password",
		Short: "get auth data",
		RunE: func(cmd *cobra.Command, args []string) error {
			username := viper.GetString("username")
			address := viper.GetString("address")
			key := viper.GetString("key")

			fmt.Print("Enter password: ")
			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			fmt.Println()

			password := string(bytePassword)

			connParams := &agent.ConnParams{Address: address, Key: key}
			if err := deps.NewAgent("").Authenticate(cmd.Context(), username, password, true, connParams); err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}
			return nil
		},
	}
	authCmd.Flags().StringP("username", "u", "", "username for auth call")

	_ = viper.BindPFlag("username", authCmd.Flags().Lookup("username"))

	viper.SetEnvPrefix("GOPHKEEPER")
	_ = viper.BindEnv("username")

	return authCmd
}
