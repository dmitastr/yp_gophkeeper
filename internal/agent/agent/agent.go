package agent

import (
	"encoding/json"
	"fmt"
	"os"

	"gophkeep/internal/agent/client"
	"gophkeep/internal/agent/config"
	"gophkeep/internal/agent/cryptomanager"
)

type ConnParams struct {
	Token   string
	Address string
	Key     string
}

type IAgent interface {
	Authenticate(username, password string, params *ConnParams) error
	Ping(params *ConnParams) error
	AddPassword(login, password string, params *ConnParams) error
}

type ProfileConfig struct {
	Token string `json:"Token"`
}

type Agent struct {
	connClient    *client.Client
	cryptoManager cryptomanager.CryptoManager
}

func NewAgent() IAgent {
	a := &Agent{cryptoManager: cryptomanager.NewCryptoManager(), connClient: client.NewClient()}
	return a
}

func (a *Agent) Authenticate(username, password string, params *ConnParams) error {
	req := &client.AuthRequest{
		Username: username,
		Password: password,
	}

	token, err := a.connClient.Authenticate(req, params.Address)
	if err != nil {
		return err
	}

	fmt.Printf("Token is %s\n", *token)
	//
	// tokenEncrypted, err := a.cryptoManager.Encrypt(token, []byte(params.Key))
	// if err != nil {
	// 	return err
	// }

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
	// tokenDecrypted, err := a.cryptoManager.Decrypt(&params.Token, []byte(params.Key))
	// if err != nil {
	// 	return fmt.Errorf("error decrypting Token: %s", err)
	// }

	err := a.connClient.Ping(params.Token, params.Address)
	if err != nil {
		return fmt.Errorf("error pinging agent: %s", err)
	}
	return nil

}

func (a *Agent) AddPassword(login, password string, params *ConnParams) error {
	// tokenDecrypted, err := a.cryptoManager.Decrypt(&params.Token, []byte(params.Key))
	//
	// if err != nil {
	// 	return fmt.Errorf("error decrypting Token: %s", err)
	// }

	err := a.connClient.AddPassword(params.Token, params.Address, login, password)
	if err != nil {
		return fmt.Errorf("error pinging agent: %s", err)
	}
	return nil

}
