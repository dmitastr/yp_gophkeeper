package config

type Config struct {
	Token   string `mapstructure:"token" json:"token"`
	Address string `mapstructure:"address" json:"address"`
}
