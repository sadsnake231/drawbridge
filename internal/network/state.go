package network

import (
	"log/slog"
	"sync"
	"time"
)

type StateManager struct {
	secretSequence []int
	knockTimeout   time.Duration
	states         sync.Map
	executor       FirewallManager
	closeTimeout   time.Duration
}

func NewStateManager(seq []int, knockTimeout time.Duration, closeTimeout time.Duration, exec FirewallManager) *StateManager {
	sm := &StateManager{
		secretSequence: seq,
		knockTimeout:   knockTimeout,
		executor:       exec,
		closeTimeout:   closeTimeout,
	}
	go sm.runCleanUp()

	return sm
}

type ipState struct {
	startTime time.Time
	stage     int
}

func (sm *StateManager) HandlePacket(ip string, port int) {
	val, exists := sm.states.Load(ip)
	if !exists {
		if sm.IsRightPort(port, 0) {
			sm.states.Store(ip, &ipState{startTime: time.Now(), stage: 1})
		}
		return
	}

	state := val.(*ipState)

	if !state.IsRightTime(sm.knockTimeout) {
		sm.deleteRecord(ip)

		if sm.IsRightPort(port, 0) {
			sm.states.Store(ip, &ipState{startTime: time.Now(), stage: 1})
		}
		return
	}

	if sm.IsRightPort(port, state.stage) {
		state.stage++
		if state.stage == len(sm.secretSequence) {
			if err := sm.executor.GrantAccess(ip); err != nil {
				slog.Error("failed to grant access", "error", err, "ip", ip)
			} else {
				slog.Info("access granted", "ip", ip)

				time.AfterFunc(sm.closeTimeout, func() {
					if err := sm.executor.RevokeAccess(ip); err != nil {
						slog.Error("failed to revoke access on timeout", "error", err, "ip", ip)
					} else {
						slog.Info("access revoked", "ip", ip)
					}
				})
			}

			sm.deleteRecord(ip)
		}
	} else {
		sm.deleteRecord(ip)
	}

}

func (sm *StateManager) IsRightPort(port, stage int) bool {
	if stage >= len(sm.secretSequence) {
		return false
	}

	return port == sm.secretSequence[stage]
}

func (ips *ipState) IsRightTime(knockTimeout time.Duration) bool {
	return time.Since(ips.startTime) <= knockTimeout
}

func (sm *StateManager) deleteRecord(ip string) {
	sm.states.Delete(ip)
}

func (sm *StateManager) runCleanUp() {
	ticker := time.NewTicker(10 * time.Minute)
	for {
		select {
		case <-ticker.C:
			sm.states.Range(func(k, v any) bool {
				ips := v.(*ipState)
				if !ips.IsRightTime(sm.knockTimeout) {
					sm.deleteRecord(k.(string))
				}
				return true
			})
		}
	}
}
