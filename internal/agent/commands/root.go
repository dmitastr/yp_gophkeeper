package commands

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gophkeep/internal/agent/agent"
	"gophkeep/internal/agent/config"
)

type RootDeps struct {
	NewAgent func() agent.IAgent
}

func NewCmd(deps RootDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cobra",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}
	cmd.PersistentFlags().StringP("key", "k", "", "key for encrypting token")
	cmd.PersistentFlags().StringP("address", "a", "", "address for calling gophkeeper")

	viper.SetEnvPrefix("AUTH") // ENV prefix: AUTH_
	_ = viper.BindEnv("username")
	_ = viper.BindEnv("address")
	// Bind flag to Viper
	// Bind to environment variable AUTH_USERNAME
	_ = viper.BindPFlag("address", cmd.PersistentFlags().Lookup("address"))
	_ = viper.BindPFlag("key", cmd.PersistentFlags().Lookup("key"))

	// Bind to environment variable: APP_GLOBAL_OPT
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	_ = viper.BindEnv("KEY")
	_ = viper.BindEnv("ADDRESS")

	viper.SetConfigFile("agent_config.json")
	err := viper.ReadInConfig()

	if err != nil && !errors.Is(err, viper.ConfigFileNotFoundError{}) {
		panic(err)
	}

	var cfg config.Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	cmd.AddCommand(NewPingCmd(deps))
	cmd.AddCommand(NewAuthCmd(deps))
	cmd.AddCommand(NewPasswordsCmd(deps))

	return cmd
}
