package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sadsnake231/drawbridge/internal/config"
	"github.com/sadsnake231/drawbridge/internal/logging"
	"github.com/sadsnake231/drawbridge/internal/network"
)

func main() {
	configPath := flag.String("config", "/etc/drawbridge/config.yaml", "specify path to config")
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

	go func() {
		err := network.StartSniffing(cfg.Interface, cfg.Snaplen, cfg.Promisc, cfg.BPFFilter, sm)
		if err != nil {
			slog.Error("Sniffer failed", "error", err)
			os.Exit(1)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	slog.Info("Shutting down Drawbridge...")

	sm.Shutdown()
}
