package client

import (
	"context"
	"strings"
	"time"
)

// Result is a domain lookup result.
type Result struct {
	Expiry   time.Time
	Statuses []string
}

// Client is a DNS client impl.
type Client interface {
	Lookup(ctx context.Context, domain string, host string) (Result, error)
}

// nolint: gochecknoglobals
var eppStatuses = func() map[string]string {
	m := map[string]string{"active": "ok"} // RDAP "active" is EPP "ok" (RFC 8056)
	for _, s := range []string{
		"ok", "inactive",
		"addPeriod", "autoRenewPeriod", "renewPeriod", "transferPeriod", "redemptionPeriod",
		"pendingCreate", "pendingDelete", "pendingRenew", "pendingRestore", "pendingTransfer", "pendingUpdate",
		"clientHold", "clientDeleteProhibited", "clientRenewProhibited", "clientTransferProhibited", "clientUpdateProhibited",
		"serverHold", "serverDeleteProhibited", "serverRenewProhibited", "serverTransferProhibited", "serverUpdateProhibited",
	} {
		m[strings.ToLower(s)] = s
	}
	return m
}()

// EPPStatus normalizes a whois/rdap status ("clientHold", "client hold",
// "ok - normal state") to its EPP code. Unknown statuses return false.
func EPPStatus(s string) (string, bool) {
	s = strings.ToLower(strings.TrimSpace(s))
	if f := strings.Fields(s); len(f) > 0 {
		if code, ok := eppStatuses[strings.Trim(f[0], ".,;")]; ok {
			return code, true
		}
	}
	code, ok := eppStatuses[strings.ReplaceAll(s, " ", "")]
	return code, ok
}

// EPPStatuses normalizes and dedups statuses, dropping unknown ones.
func EPPStatuses(in []string) []string {
	var out []string
	seen := map[string]bool{}
	for _, s := range in {
		if code, ok := EPPStatus(s); ok && !seen[code] {
			seen[code] = true
			out = append(out, code)
		}
	}
	return out
}
