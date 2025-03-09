package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCPort int `yaml:"grpc_port" env:"GRPC_PORT" env-default:"50051"`
}

func New() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig("./config/config.env", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
