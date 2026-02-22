package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/sadsnake231/drawbridge/internal/config"
	"github.com/sadsnake231/drawbridge/internal/network"
)

func main() {
	configPath := flag.String("config", "./drawbridge.yaml", "specify path to config")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	exec := network.NewIPTablesExecutor(cfg.SafePort, cfg.CloseTimeout)

	sm := network.NewStateManager(cfg.Sequence, cfg.KnockTimeout, exec)

	fmt.Printf("Starting Drawbridge on interface: %s\n", cfg.Interface)
	network.StartSniffing(cfg.Interface, sm)
}
