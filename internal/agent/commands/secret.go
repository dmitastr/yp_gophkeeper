package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gophkeep/internal/agent/agent"
	"gophkeep/internal/core/models"
	"gophkeep/internal/core/requests"
)

func NewSecretsCmd(deps agent.RootDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Long:  "Root command for handling secrets",
		Short: "Secrets root command",
	}

	cmd.AddCommand(
		addSecretCmd(deps),
		getAllSecretsCmd(deps),
		getSecretCmd(deps),
	)

	return cmd
}

func addSecretCmd(deps agent.RootDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Long:  "Add secret command for adding secrets",
		Short: "Secrets add command",
		RunE: func(cmd *cobra.Command, args []string) error {
			iAgent := deps.NewAgent()

			secretType := viper.GetString("secret-type")
			input := viper.GetString("input")
			file := viper.GetString("file")

			address := viper.GetString("address")
			key := viper.GetString("key")
			token := viper.GetString("token")
			comment := viper.GetString("comment")

			var body []byte
			if file != "" {
				var err error
				body, err = iAgent.ReadSecretFromFile(file)
				if err != nil {
					return fmt.Errorf("error reading secret from file: %s", err)
				}
			} else {
				body = []byte(input)
			}

			deps.Logger.Info("Adding secret to app", zap.String("secret-type", secretType), zap.String("comment", comment))

			connParams := &agent.ConnParams{Address: address, Key: key, Token: token}

			secret := &requests.SecretRequest{Body: body, SecretType: models.SecretType(secretType), Comment: comment}
			if err := iAgent.AddSecret(secret, connParams); err != nil {
				deps.Logger.Error(err.Error())
				return fmt.Errorf("authentication failed: %w", err)
			}

			return nil
		},
	}
	cmd.Flags().StringP("secret-type", "s", "", "secret type to add")
	cmd.Flags().StringP("input", "i", "", "secret input string")
	cmd.Flags().StringP("file", "f", "", "file path to read secret from")
	cmd.Flags().StringP("comment", "c", "", "optional comment for secret; max length 200")

	_ = viper.BindPFlag("secret-type", cmd.Flags().Lookup("secret-type"))
	_ = viper.BindPFlag("input", cmd.Flags().Lookup("input"))
	_ = viper.BindPFlag("file", cmd.Flags().Lookup("file"))
	_ = viper.BindPFlag("comment", cmd.Flags().Lookup("comment"))

	viper.SetEnvPrefix("GOPHKEEPER")
	_ = viper.BindEnv("secret-type", "SECRET_TYPE")
	_ = viper.BindEnv("input")
	_ = viper.BindEnv("file")
	_ = viper.BindEnv("comment")

	return cmd
}

func getAllSecretsCmd(deps agent.RootDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Long:  "List all secrets for user",
		Short: "List secrets",
		RunE: func(cmd *cobra.Command, args []string) error {
			iAgent := deps.NewAgent()

			address := viper.GetString("address")
			key := viper.GetString("key")
			token := viper.GetString("token")

			connParams := &agent.ConnParams{Address: address, Key: key, Token: token}

			secrets, err := iAgent.GetAllSecrets(connParams)
			if err != nil {
				return fmt.Errorf("error getting secrets: %w", err)
			}

			fmt.Printf("Fetched %d secrets\n", len(secrets))

			w := tabwriter.NewWriter(os.Stdout, 10, 1, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tTYPE\tCREATED_AT\t")

			for _, secret := range secrets {
				fmt.Fprintf(w, "%d\t%s\t%s\t\n", secret.ID, secret.Type, secret.CreatedAt)
			}

			w.Flush()

			return nil
		},
	}

	return cmd
}

func getSecretCmd(deps agent.RootDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Long:  "Get secret by id",
		Short: "Get secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			iAgent := deps.NewAgent()

			address := viper.GetString("address")
			key := viper.GetString("key")
			token := viper.GetString("token")
			secretID := viper.GetInt("id")
			output := viper.GetString("output")

			connParams := &agent.ConnParams{Address: address, Key: key, Token: token}

			secret, err := iAgent.GetSecret(secretID, connParams)
			if err != nil {
				return fmt.Errorf("error getting secret: %w", err)
			}

			if output == "" {
				fmt.Printf("ID: %d\n", secret.ID)
				fmt.Printf("Created at: %s\n", secret.CreatedAt)
				fmt.Printf("Secret type: %s\n", secret.Type)
				fmt.Printf("Content: %s\n", string(secret.Content))
				fmt.Printf("Comment: %s\n", secret.Comment)
				return nil
			}

			if err := iAgent.WriteSecretToFile(secret.Content, output); err != nil {
				return fmt.Errorf("error writing secret to file: %w", err)
			}

			return nil
		},
	}
	cmd.Flags().IntP("id", "d", -1, "secret id")
	cmd.Flags().StringP("output", "o", "", "file to write secret to; default to stdout")

	_ = viper.BindPFlag("id", cmd.Flags().Lookup("id"))
	_ = viper.BindPFlag("output", cmd.Flags().Lookup("output"))

	viper.SetEnvPrefix("GOPHKEEPER")
	_ = viper.BindEnv("id")
	_ = viper.BindEnv("output")

	return cmd
}
