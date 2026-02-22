package config

import (
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Interface    string        `yaml:"interface"`
	Sequence     []uint16      `yaml:"sequence"`
	KnockTimeout time.Duration `yaml:"knock-timeout"`
	SafePort     uint16        `yaml:"safe-port"`
	CloseTimeout time.Duration `yaml:"close-timeout"`
	LogFile      string        `yaml:"log-file"`
	Snaplen      int32         `yaml:"snaplen"`
	Promisc      bool          `yaml:"promisc"`
	BPFFilter    string        `yaml:"bpf-filter"`
}

func LoadConfig(path string) (*Config, error) {
	// Установка значений по умолчанию
	cfg := Config{
		LogFile:   "/var/log/drawbridge.log",
		Snaplen:   1024,
		Promisc:   false,
		BPFFilter: "tcp[tcpflags] & (tcp-syn) != 0",
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
