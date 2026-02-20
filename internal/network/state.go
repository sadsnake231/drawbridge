package network

import (
	"fmt"
	"sync"
	"time"
)

type StateManager struct {
	secretSequence []int
	timeout        time.Duration
	states         sync.Map
	executor       FirewallManager
}

func NewStateManager(seq []int, t time.Duration, exec FirewallManager) *StateManager {
	sm := &StateManager{secretSequence: seq, timeout: t, executor: exec}
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

	if !state.IsRightTime(sm.timeout) {
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
				fmt.Printf("Error granting access: %v\n", err)
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

func (ips *ipState) IsRightTime(timeout time.Duration) bool {
	return time.Since(ips.startTime) <= timeout
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
				if !ips.IsRightTime(sm.timeout) {
					sm.deleteRecord(k.(string))
				}
				return true
			})
		}
	}
}
