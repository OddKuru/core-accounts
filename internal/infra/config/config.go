package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type (
	Config struct {
		App    App    `yaml:"app"`
		Logger Logger `yaml:"logger"`
		Psql   Psql   `yaml:"postgresql"`
	}

	App struct {
		StartTimeout    time.Duration `yaml:"start_timeout" env-default:"5s"`
		GracefulTimeout time.Duration `yaml:"graceful_timeout" env-default:"10s"`
		GRPCAddress     string        `yaml:"grpc_address"`
		GRPCTimeout     time.Duration `yaml:"grpc_timeout" env-default:"10s"`
	}

	Psql struct {
		Host                string        `yaml:"host"`
		Port                int           `yaml:"port"`
		Username            string        `yaml:"username"`
		Password            string        `yaml:"password" json:"-"`
		Database            string        `yaml:"database"`
		SSLMode             string        `yaml:"sslmode"`
		MaxConnections      int32         `yaml:"max_connections" env-default:"20"`
		MinConnections      int32         `yaml:"min_connections" env-default:"2"`
		MaxConnLifetime     time.Duration `yaml:"max_connection_lifetime" env-default:"10m"`
		MaxConnIdleTime     time.Duration `yaml:"max_connection_idle_time" env-default:"10m"`
		HealthCheckInterval time.Duration `yaml:"health_check_interval" env-default:"10m"`
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
