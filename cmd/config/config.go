package config

import (
	"errors"
	"os"
)

type Config struct {
	PublicConfig  PublicConfig
	PrivateConfig PrivateConfig
}

type PublicConfig struct {
	PublicKey        string
	PostgresUser     string
	PostgresHost     string
	PostgresPort     string
	PostgresDatabase string
	PostgresPassword string
}

type PrivateConfig struct {
	SecretKey        string
	PostgresPassword string
}

var (
	envs = map[string]PublicConfig{
		"local": {
			PublicKey:        "localAdmin",
			PostgresUser:     "c4g9",
			PostgresPort:     "5432",
			PostgresHost:     "localhost",
			PostgresDatabase: "test-database",
			PostgresPassword: "pass_c4g9",
		},
		"dev": {
			PublicKey: "devAdmin",
		},
		"prod": {
			PublicKey: "prodAdmin",
		},
	}
)

func NewConfig(env string) (Config, error) {

	publicConfig, exists := envs[env]
	if !exists {
		return Config{}, errors.New("env doest not exists")
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return Config{}, errors.New("secret key doest not exists in env")
	}

	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	if postgresPassword == "" {
		return Config{}, errors.New("postgres password doest not exists in env")
	}

	return Config{
		PublicConfig: publicConfig,
		PrivateConfig: PrivateConfig{
			SecretKey:        secretKey,
			PostgresPassword: postgresPassword,
		},
	}, nil

}
