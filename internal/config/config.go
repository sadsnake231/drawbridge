package config

import (
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Interface    string        `yaml:"interface"`
	Sequence     []int         `yaml:"sequence"`
	KnockTimeout time.Duration `yaml:"knock-timeout"`
	SafePort     int           `yaml:"safe-port"`
	CloseTimeout time.Duration `yaml:"close-timeout"`
	LogFile      string        `yaml:"log-file" default:"/var/log/drawbridge.log"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil

}
