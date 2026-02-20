package network

type FirewallManager interface {
	GrantAccess(ip string) error
	RevokeAccess(ip string) error
}
