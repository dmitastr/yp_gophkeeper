package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"
	"gophkeep/internal/agent/client"
	"gophkeep/internal/agent/config"
	"gophkeep/internal/agent/filereader"
	"gophkeep/internal/core/models"
	"gophkeep/internal/core/requests"
	"gophkeep/internal/core/validation"
	"gophkeep/internal/logger"
)

type ConnParams struct {
	Token   string
	Address string
	Key     string
}

type IAgent interface {
	Authenticate(ctx context.Context, username, password string, isNewUser bool, params *ConnParams) error
	Ping(ctx context.Context, params *ConnParams) error
	AddSecret(ctx context.Context, secretRequest *requests.SecretRequest, params *ConnParams) error
	UpdateSecret(ctx context.Context, secretRequest *requests.SecretRequest, params *ConnParams) error
	DeleteSecret(ctx context.Context, secretID int, params *ConnParams) error
	GetAllSecrets(ctx context.Context, params *ConnParams) ([]models.SecretInfo, error)
	GetSecret(ctx context.Context, secretID int, params *ConnParams) (*models.Secret, error)
	ReadSecretFromFile(file string) (*requests.SecretRequest, error)
	WriteSecretToFile(content []byte, file string) error
	ParseInput(input, fileName, comment string, secretType models.SecretType) (*requests.SecretRequest, error)
}

type Agent struct {
	connClient client.IClient
	fileReader filereader.FileInputReader
	logger.ILogger
	secretValidator validation.IValidator
}

func NewAgent(log logger.ILogger, bearerToken string) IAgent {
	connClient := client.NewRestyClient(log, bearerToken)
	a := &Agent{
		connClient:      connClient,
		fileReader:      filereader.NewFileInputReader(),
		ILogger:         log,
		secretValidator: validation.NewValidator(nil),
	}
	return a
}

func (a *Agent) Authenticate(ctx context.Context, username, password string, isNewUser bool, params *ConnParams) error {
	if username == "" || password == "" {
		return ErrorEmptyAuthData
	}

	req := &requests.AuthRequest{
		Username:  username,
		Password:  password,
		IsNewUser: isNewUser,
	}

	token, err := a.connClient.Authenticate(ctx, req, params.Address)
	if err != nil {
		return err
	}

	a.Info("Receive auth data", zap.Stringp("token", token))

	file, err := os.Create("agent_config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	cfg := &config.Config{Token: *token, Address: params.Address}
	if err := json.NewEncoder(file).Encode(cfg); err != nil {
		return err
	}

	return nil
}

func (a *Agent) Ping(ctx context.Context, params *ConnParams) error {
	err := a.connClient.Ping(ctx, params.Address)
	if err != nil {
		return fmt.Errorf("error pinging agent: %s", err)
	}
	return nil

}

func (a *Agent) AddSecret(ctx context.Context, secretRequest *requests.SecretRequest, params *ConnParams) error {
	if err := a.connClient.AddSecret(ctx, params.Address, secretRequest); err != nil {
		return fmt.Errorf("error adding secretRequest: %s", err)
	}
	return nil
}

func (a *Agent) GetAllSecrets(ctx context.Context, params *ConnParams) ([]models.SecretInfo, error) {
	secret, err := a.connClient.GetAllSecrets(ctx, params.Address)
	if err != nil {
		a.Error("error getting all secrets: %s", zap.Error(err))
		return nil, err
	}
	return secret, nil
}

func (a *Agent) GetSecret(ctx context.Context, secretID int, params *ConnParams) (*models.Secret, error) {
	secret, err := a.connClient.GetSecret(ctx, params.Address, secretID)
	if err != nil {
		a.Error("error getting secret: %s", zap.Error(err), zap.Int("secretID", secretID))

		return nil, err
	}
	return secret, nil
}

func (a *Agent) UpdateSecret(ctx context.Context, secretRequest *requests.SecretRequest, params *ConnParams) error {
	if err := a.connClient.AddSecret(ctx, params.Address, secretRequest); err != nil {
		return fmt.Errorf("error adding secret: %s", err)
	}
	return nil
}

func (a *Agent) DeleteSecret(ctx context.Context, secretID int, params *ConnParams) error {

	if err := a.connClient.DeleteSecret(ctx, params.Address, secretID); err != nil {
		return fmt.Errorf("error deleting secret: %s", err)
	}
	return nil
}

func (a *Agent) ReadSecretFromFile(file string) (*requests.SecretRequest, error) {
	fileContent, err := a.fileReader.FileRead(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	var r requests.SecretRequest
	if err := json.Unmarshal(fileContent, &r); err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %w", err)
	}
	return &r, nil
}

func (a *Agent) WriteSecretToFile(content []byte, file string) error {
	return a.fileReader.FileWrite(content, file)
}

func (a *Agent) ParseInput(input, fileName, comment string, secretType models.SecretType) (*requests.SecretRequest, error) {
	return a.fileReader.ParseInput(input, fileName, comment, secretType)
}
