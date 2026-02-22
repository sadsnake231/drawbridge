package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"

	"github.com/sadsnake231/drawbridge/internal/config"
	"github.com/sadsnake231/drawbridge/internal/logging"
	"github.com/sadsnake231/drawbridge/internal/network"
)

func main() {
	configPath := flag.String("config", "./drawbridge.yaml", "specify path to config")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logFile, err := logging.InitLogger(cfg.LogFile)
	if err != nil {
		log.Fatalf("Failed to initialize log file: %v", err)
	}
	defer logFile.Close()

	slog.Info("Drawbridge started", "interface", cfg.Interface)

	exec := network.NewIPTablesExecutor(cfg.SafePort)

	sm := network.NewStateManager(cfg.Sequence, cfg.KnockTimeout, cfg.CloseTimeout, exec)

	fmt.Printf("Starting Drawbridge on interface: %s\n", cfg.Interface)
	network.StartSniffing(cfg.Interface, sm)
}
