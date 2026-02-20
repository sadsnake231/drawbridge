package network

import (
	"fmt"
	"time"
)

type LogExecutor struct{}

func (le *LogExecutor) GrantAccess(ip string) error {
	fmt.Printf("Access Granted for %v\n", ip)

	time.AfterFunc(1*time.Hour, func() {
		le.RevokeAccess(ip)
	})

	return nil
}

func (le *LogExecutor) RevokeAccess(ip string) error {
	fmt.Printf("Access Revoked for %v\n", ip)
	return nil
}
