package network

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

type flag = string

const (
	check  flag = "-C"
	add    flag = "-A"
	delete flag = "-D"
)

type IPTablesExecutor struct {
	safePort string
	timeout  time.Duration
}

func NewIPTablesExecutor(safePort int, timeout time.Duration) *IPTablesExecutor {
	ipte := &IPTablesExecutor{safePort: strconv.Itoa(safePort), timeout: timeout}

	return ipte
}

func (ipte *IPTablesExecutor) GrantAccess(ip string) error {
	err := ipte.ruleAction(ip, check)
	if err != nil {
		err = ipte.ruleAction(ip, add)
		if err != nil {
			fmt.Printf("couldn't add iptables rule: %v\n", err.Error())
			return err
		}
		fmt.Printf("access granted for %v\n", ip)

		time.AfterFunc(ipte.timeout, func() {
			ipte.RevokeAccess(ip)
		})
	}

	return nil
}

func (ipte *IPTablesExecutor) RevokeAccess(ip string) error {
	err := ipte.ruleAction(ip, delete)
	if err != nil {
		fmt.Printf("couldn't delete iptables rule: %v", err.Error())
		return err
	}

	fmt.Printf("access revoked for %v", ip)
	return nil
}

func (ipte *IPTablesExecutor) ruleAction(ip, key string) error {
	cmd := exec.Command("iptables", key, "INPUT", "-s", ip, "-p", "tcp", "--dport", ipte.safePort, "-j", "ACCEPT")
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
