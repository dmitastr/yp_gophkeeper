package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gophkeep/internal/agent/agent"
	"gophkeep/internal/agent/formatter"
	"gophkeep/internal/core/models"
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
		updateSecretCmd(deps),
		deleteSecretCmd(deps),
	)

	return cmd
}

func addSecretCmd(deps agent.RootDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use: "add",
		Long: `Add a secret by providing a JSON file or from command line argument
			
			Available secret types:
			gophkeep secrets add --file path/to/password.json
			
			// password.json
			{"login": "my_login", "password": "567890", "comment": "login for example.com" }
			-------
			gophkeep secrets add --file path/to/bank_card.json 
			
			// bank_card.json
			 {
				"number": 5120350100064537,
				"expiry_data": "01/30",
				"cardholder_name": "John Smith",
				"cvc": "000",
				"comment": "This is a credit card"
			}
			-------
			Accept any text or binary files
			
			gophkeep secrets add --file path/to/text_file -s text
			gophkeep secrets add --file path/to/binary_file -s binary`,

		Short: "Secrets add command",
		Example: `# Add a password 
		  gophkeep secrets add --file examples/password.json -s password
		
		  # Add a bank card 
		  gophkeep secrets add --file examples/bank_card.json -s bank_card
		
		  # Read from command line
		  gophkeep secrets add --input "text to add" --secret-type text`,
		Args: cobra.OnlyValidArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			iAgent := deps.NewAgent()

			secretType := viper.GetString("secret-type")
			input := viper.GetString("input")
			file := viper.GetString("file")

			address := viper.GetString("address")
			key := viper.GetString("key")
			token := viper.GetString("token")
			comment := viper.GetString("comment")

			secretRequest, err := iAgent.ParseInput(input, file, comment, models.SecretType(secretType))
			if err != nil {
				return fmt.Errorf("error parsing secret from command line arguments: %s", err)
			}

			deps.Logger.Info("Adding secret to app",
				zap.String("secret-type", string(secretRequest.SecretType)),
				zap.String("comment", secretRequest.Comment),
				zap.Int("body size", len(secretRequest.BodyString)),
			)

			connParams := &agent.ConnParams{Address: address, Key: key, Token: token}

			if err := iAgent.AddSecret(cmd.Context(), secretRequest, connParams); err != nil {
				deps.Logger.Error(err.Error())
				return fmt.Errorf("adding secret failed: %w", err)
			}
			cmd.Println("Secret added successfully")

			return nil
		},
	}

	var secretType models.SecretType
	cmd.Flags().StringP("input", "i", "", "secret input string")
	cmd.Flags().StringP("file", "f", "", "file path to read secret from")
	cmd.Flags().StringP("comment", "c", "", "optional comment for secret; max length 200")
	cmd.Flags().VarP(&models.SecretTypeValue{Value: &secretType}, "secret-type", "s", "secret type to add; (password|text|bank_card|binary)")

	cmd.MarkFlagsMutuallyExclusive("file", "input")
	cmd.MarkFlagsOneRequired("file", "input")

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

			secrets, err := iAgent.GetAllSecrets(cmd.Context(), connParams)
			if err != nil {
				return fmt.Errorf("error getting secrets: %w", err)
			}

			fmt.Printf("Fetched %d secrets\n", len(secrets))

			out := cmd.OutOrStdout()
			formatOutput := formatter.NewFormatter()
			formatOutput.FormatOutput(out, secrets)

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
			secretID := viper.GetInt("get-secret.id")
			output := viper.GetString("output")

			deps.Logger.Info("Get secret", zap.Int("secretID", secretID))

			connParams := &agent.ConnParams{Address: address, Key: key, Token: token}

			secret, err := iAgent.GetSecret(cmd.Context(), secretID, connParams)
			if err != nil {
				return fmt.Errorf("error getting secret: %w", err)
			}

			fmt.Printf("ID: %d\n", secret.ID)
			fmt.Printf("Created at: %s\n", secret.CreatedAt)
			fmt.Printf("Secret type: %s\n", secret.Type)
			fmt.Printf("Comment: %s\n", secret.Comment)

			if output == "" {
				fmt.Printf("Content: %s\n", string(secret.Content))
				return nil
			}

			if err := iAgent.WriteSecretToFile(secret.Content, output); err != nil {
				return fmt.Errorf("error writing secret to file: %w", err)
			}
			fmt.Printf("Content was written to file: %s\n", output)

			return nil
		},
	}
	cmd.Flags().IntP("id", "d", 0, "secret id")
	cmd.Flags().StringP("output", "o", "", "file to write secret to; default to stdout")

	cmd.MarkFlagRequired("id")

	_ = viper.BindPFlag("get-secret.id", cmd.Flags().Lookup("id"))
	_ = viper.BindPFlag("output", cmd.Flags().Lookup("output"))

	viper.SetEnvPrefix("GOPHKEEPER")
	_ = viper.BindEnv("id")
	_ = viper.BindEnv("output")

	return cmd
}

func updateSecretCmd(deps agent.RootDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "update",
		Long: `Update a secret by providing a JSON file or from command line argument. Check examples for secrets add command`,

		Short: "Secrets update command",
		RunE: func(cmd *cobra.Command, args []string) error {
			iAgent := deps.NewAgent()

			secretType := viper.GetString("update-secret.secret-type")
			input := viper.GetString("update-secret.input")
			file := viper.GetString("update-secret.file")
			comment := viper.GetString("update-secret.comment")
			secretID := viper.GetInt("update-secret.id")

			address := viper.GetString("address")
			key := viper.GetString("key")
			token := viper.GetString("token")

			secretRequest, err := iAgent.ParseInput(input, file, comment, models.SecretType(secretType))
			if err != nil {
				return fmt.Errorf("error parsing secret from command line arguments: %s", err)
			}

			secretRequest.ID = &secretID

			deps.Logger.Info("Update secret to app",
				zap.Int("secretID", secretID),
				zap.String("secret-type", string(secretRequest.SecretType)),
				zap.String("comment", secretRequest.Comment),
				zap.Int("body size", len(secretRequest.BodyString)),
			)

			connParams := &agent.ConnParams{Address: address, Key: key, Token: token}

			if err := iAgent.UpdateSecret(cmd.Context(), secretRequest, connParams); err != nil {
				deps.Logger.Error(err.Error())
				return fmt.Errorf("updating secret failed: %w", err)
			}

			return nil
		},
	}

	var secretType models.SecretType
	cmd.Flags().IntP("id", "d", 0, "secret id")
	cmd.Flags().StringP("input", "i", "", "secret input string")
	cmd.Flags().StringP("file", "f", "", "file path to read secret from")
	cmd.Flags().StringP("comment", "c", "", "optional comment for secret; max length 200")
	cmd.Flags().VarP(&models.SecretTypeValue{Value: &secretType}, "secret-type", "s", "secret type to add; (password|text|bank_card|binary)")

	cmd.MarkFlagsMutuallyExclusive("file", "input")
	cmd.MarkFlagsOneRequired("file", "input")
	cmd.MarkFlagRequired("id")

	_ = viper.BindPFlag("update-secret.secret-type", cmd.Flags().Lookup("secret-type"))
	_ = viper.BindPFlag("update-secret.input", cmd.Flags().Lookup("input"))
	_ = viper.BindPFlag("update-secret.file", cmd.Flags().Lookup("file"))
	_ = viper.BindPFlag("update-secret.comment", cmd.Flags().Lookup("comment"))
	_ = viper.BindPFlag("update-secret.id", cmd.Flags().Lookup("id"))

	viper.SetEnvPrefix("GOPHKEEPER")
	_ = viper.BindEnv("secret-type", "SECRET_TYPE")
	_ = viper.BindEnv("input")
	_ = viper.BindEnv("file")
	_ = viper.BindEnv("comment")

	return cmd
}

func deleteSecretCmd(deps agent.RootDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: `Delete a secret`,

		Short: "Secrets delete command",
		RunE: func(cmd *cobra.Command, args []string) error {
			iAgent := deps.NewAgent()

			address := viper.GetString("address")
			key := viper.GetString("key")
			token := viper.GetString("token")
			secretID := viper.GetInt("delete-secret.id")

			deps.Logger.Info("Delete secret",
				zap.Int("secretID", secretID),
			)

			connParams := &agent.ConnParams{Address: address, Key: key, Token: token}

			if err := iAgent.DeleteSecret(cmd.Context(), secretID, connParams); err != nil {
				deps.Logger.Error(err.Error())
				return fmt.Errorf("updating secret failed: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().IntP("id", "d", 0, "secret id")

	cmd.MarkFlagRequired("id")

	_ = viper.BindPFlag("delete-secret.id", cmd.Flags().Lookup("id"))

	return cmd
}
