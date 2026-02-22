package network

import (
	"fmt"
	"os/exec"
)

type flag = string

const (
	check  flag = "-C"
	add    flag = "-A"
	delete flag = "-D"
)

type IPTablesExecutor struct {
	safePort string
}

func NewIPTablesExecutor(safePort uint16) *IPTablesExecutor {
	return &IPTablesExecutor{safePort: fmt.Sprintf("%d", safePort)}
}

func (ipte *IPTablesExecutor) GrantAccess(ip string) error {
	if err := ipte.ruleAction(ip, check); err != nil {
		if err = ipte.ruleAction(ip, add); err != nil {
			return fmt.Errorf("couldn't add iptables rule: %w", err)
		}
	}
	return nil
}

func (ipte *IPTablesExecutor) RevokeAccess(ip string) error {
	if err := ipte.ruleAction(ip, delete); err != nil {
		return fmt.Errorf("couldn't delete iptables rule: %w", err)
	}
	return nil
}

func (ipte *IPTablesExecutor) ruleAction(ip, key string) error {
	cmd := exec.Command("iptables", "-w", key, "INPUT", "-s", ip, "-p", "tcp", "--dport", ipte.safePort, "-j", "ACCEPT")
	return cmd.Run()
}
