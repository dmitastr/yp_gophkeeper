package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gophkeep/internal/logger"
)

type ConfigProvider interface {
	Logger() logger.ILogger
	GetConfig() Config
}

type configProvider struct {
	config Config
	logger logger.ILogger
}

type Config struct {
	Address   string `json:"address" env:"ADDRESS" mapstructure:"address"`
	DBConnStr string `json:"db_conn_str" env:"DBCONNSTR" mapstructure:"db-conn-str"`
	Key       string `json:"key" env:"KEY" mapstructure:"key"`
	MaxSize   int64  `json:"max_size" env:"MAX_SIZE" mapstructure:"max-size"`
}

func NewConfig() (ConfigProvider, error) {
	log := logger.NewLogger()
	cp := configProvider{logger: log}

	maxSizeDefault := int64(100) << 20

	flagSet := pflag.NewFlagSet("GophKeeper", pflag.ExitOnError)
	flagSet.StringP("address", "a", "localhost:8080", "set app host and port")
	flagSet.StringP("db-conn-str", "d", "", "postgres connection url")
	flagSet.StringP("key", "k", "", "key for encrypting data")
	flagSet.Int64P("max-size", "m", maxSizeDefault, "key for encrypting data")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("error parsing flags: %w", err)
	}

	_ = viper.BindPFlags(flagSet)

	viper.AutomaticEnv()

	// Bind environment variables
	_ = viper.BindEnv("address", "ADDRESS")
	_ = viper.BindEnv("db-conn-str", "DBCONNSTR")
	_ = viper.BindEnv("key", "KEY")
	_ = viper.BindEnv("max-size", "MAX_SIZE")

	if cfgPath := viper.GetString("configProvider"); cfgPath != "" {
		viper.SetConfigFile(cfgPath)
		if err := viper.ReadInConfig(); err != nil {
			cp.Logger().Error("error reading configProvider file, %s\n", zap.Error(err))
		}
	}

	var cfg Config
	// Unmarshal the configuration into the Config struct
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %w", err)
	}

	cp.config = cfg

	return &cp, nil

}

func (c configProvider) Logger() logger.ILogger {
	return c.logger
}

func (c configProvider) GetConfig() Config {
	return c.config
}
