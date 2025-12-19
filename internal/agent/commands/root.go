package commands

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gophkeep/internal/agent/agent"
	"gophkeep/internal/agent/config"
)

func NewCmd(deps agent.RootDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cobra",
		Short: "Root cobra command",
		Long: `Gophkeeper is a CLI library for Go for interaction with Gophkeeper application.
			It allows to store and retrieve secrets.`,
	}
	cmd.PersistentFlags().StringP("key", "k", "", "key for encrypting token")
	cmd.PersistentFlags().StringP("address", "a", "", "address for calling gophkeeper")

	viper.AutomaticEnv()

	viper.SetEnvPrefix("GOPHKEEPER")
	_ = viper.BindEnv("address")
	_ = viper.BindEnv("KEY")

	_ = viper.BindPFlag("address", cmd.PersistentFlags().Lookup("address"))
	_ = viper.BindPFlag("key", cmd.PersistentFlags().Lookup("key"))

	viper.SetConfigName("agent_config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			panic(err)
		}
	}

	var cfg config.Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	cmd.AddCommand(
		NewPingCmd(deps),
		NewLoginCmd(deps),
		NewAuthCmd(deps),
		NewPasswordsCmd(deps),
		NewSecretsCmd(deps),
	)

	return cmd
}
