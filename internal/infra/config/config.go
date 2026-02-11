package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type (
	Config struct {
		Logger Logger `yaml:"logger"`
		Psql   Psql   `yaml:"postgresql"`
	}

	Psql struct {
		Host            string        `yaml:"host"`
		Port            int           `yaml:"port"`
		Username        string        `yaml:"username"`
		Password        string        `yaml:"password"`
		Database        string        `yaml:"database"`
		SSLMode         string        `yaml:"sslmode"`
		MaxConnections  int32         `yaml:"max_connections"`
		MinConnections  int32         `yaml:"min_connections"`
		MaxConnLifetime time.Duration `yaml:"max_connection_lifetime"`
		MaxConnIdleTime time.Duration `yaml:"max_connection_idle_time"`
	}

	Logger struct {
		Level  string `yaml:"level"`
		Pretty bool   `yaml:"pretty" env-default:"false"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	_, err := os.Stat(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot find config file path")
	}

	cfg := &Config{}
	err = cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read config file")
	}

	return cfg, nil
}
