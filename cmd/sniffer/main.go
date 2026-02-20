package main

import (
	"time"

	"github.com/sadsnake231/drawbridge/internal/network"
)

func main() {
	exec := &network.LogExecutor{}

	seq := []int{1111, 2222, 3333}
	sm := network.NewStateManager(seq, 10*time.Second, exec)

	network.StartSniffing("lo", sm)
}
