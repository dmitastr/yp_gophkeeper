package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
	"gophkeep/internal/agent/agent"
	"gophkeep/internal/agent/client"
)

type LoginCmd struct {
	Cmd *cobra.Command
	c   *client.Client
}

func NewLoginCmd(deps agent.RootDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Long:  "Login existing user",
		Short: "login user",
		RunE: func(cmd *cobra.Command, args []string) error {
			username := viper.GetString("username")
			address := viper.GetString("address")
			key := viper.GetString("key")

			fmt.Print("Enter password: ")
			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			fmt.Println() // newline after password input

			password := string(bytePassword)

			connParams := &agent.ConnParams{Address: address, Key: key}
			if err := deps.NewAgent().Authenticate(username, password, false, connParams); err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringP("username", "u", "", "username for auth call")

	_ = viper.BindPFlag("username", cmd.Flags().Lookup("username"))

	viper.SetEnvPrefix("GOPHKEEPER")
	_ = viper.BindEnv("username")

	return cmd
}
