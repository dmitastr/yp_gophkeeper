package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
	"gophkeep/internal/agent/client"
	"gophkeep/internal/agent/config"
	"gophkeep/internal/agent/cryptomanager"
	"gophkeep/internal/agent/filereader"
	"gophkeep/internal/core/models"
	"gophkeep/internal/core/requests"
	"gophkeep/internal/logger"
)

type ConnParams struct {
	Token   string
	Address string
	Key     string
}

type IAgent interface {
	Authenticate(username, password string, isNewUser bool, params *ConnParams) error
	Ping(params *ConnParams) error
	AddPassword(login, password string, params *ConnParams) error
	AddSecret(secret *requests.SecretRequest, params *ConnParams) error
	GetAllSecrets(params *ConnParams) ([]models.SecretInfo, error)
	GetSecret(secretID int, params *ConnParams) (*models.Secret, error)
	ReadSecretFromFile(file string) ([]byte, error)
	WriteSecretToFile(content []byte, file string) error
}

type Agent struct {
	connClient    client.IClient
	cryptoManager cryptomanager.CryptoManager
	fileReader    filereader.FileInputReader
	logger.ILogger
}

func NewAgent(log logger.ILogger) IAgent {
	a := &Agent{
		cryptoManager: cryptomanager.NewCryptoManager(),
		connClient:    client.NewClient(),
		fileReader:    filereader.NewFileInputReader(),
		ILogger:       log,
	}
	return a
}

func (a *Agent) Authenticate(username, password string, isNewUser bool, params *ConnParams) error {
	if username == "" || password == "" {
		return errors.New("username and Password cannot be empty")
	}

	req := &requests.AuthRequest{
		Username:  username,
		Password:  password,
		IsNewUser: isNewUser,
	}

	token, err := a.connClient.Authenticate(req, params.Address)
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

func (a *Agent) Ping(params *ConnParams) error {
	err := a.connClient.Ping(params.Token, params.Address)
	if err != nil {
		return fmt.Errorf("error pinging agent: %s", err)
	}
	return nil

}

func (a *Agent) AddPassword(login, password string, params *ConnParams) error {
	err := a.connClient.AddPassword(params.Token, params.Address, login, password)
	if err != nil {
		return fmt.Errorf("error pinging agent: %s", err)
	}
	return nil

}

func (a *Agent) AddSecret(secret *requests.SecretRequest, params *ConnParams) error {
	switch secret.SecretType {
	case models.PASSWORD:
		var passw models.PasswordRequestObject
		if err := json.Unmarshal(secret.Body, &passw); err != nil {
			return fmt.Errorf("error unmarshalling secret: %s", err)
		}
	}

	err := a.connClient.AddSecret(params.Token, params.Address, secret)
	if err != nil {
		return fmt.Errorf("error adding secret: %s", err)
	}
	return nil

}

func (a *Agent) GetAllSecrets(params *ConnParams) ([]models.SecretInfo, error) {
	secret, err := a.connClient.GetAllSecrets(params.Token, params.Address)
	if err != nil {
		a.Error("error getting all secrets: %s", zap.Error(err))
		return nil, err
	}
	return secret, nil
}

func (a *Agent) GetSecret(secretID int, params *ConnParams) (*models.Secret, error) {
	secret, err := a.connClient.GetSecret(params.Token, params.Address, secretID)
	if err != nil {
		a.Error("error getting secret: %s", zap.Error(err), zap.Int("secretID", secretID))

		return nil, err
	}
	return secret, nil
}

func (a *Agent) ReadSecretFromFile(file string) ([]byte, error) {
	return a.fileReader.FileRead(file)
}

func (a *Agent) WriteSecretToFile(content []byte, file string) error {
	return a.WriteSecretToFile(content, file)
}
