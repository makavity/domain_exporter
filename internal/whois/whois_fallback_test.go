package whois

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/caarlos0/domain_exporter/internal/client"
	"github.com/stretchr/testify/require"
)

func TestFallbackToRDAPSucceedsWhenWhoisFails(t *testing.T) {
	previous := rdapFallbackLookup
	t.Cleanup(func() {
		rdapFallbackLookup = previous
	})

	expected := time.Date(2026, 7, 17, 0, 0, 0, 0, time.UTC)
	rdapFallbackLookup = func(_ context.Context, domain string) (client.Result, error) {
		require.Equal(t, "tanukifamily.ru", domain)
		return client.Result{Expiry: expected, Statuses: []string{"ok"}}, nil
	}

	expiration, err := fallbackToRDAP(context.Background(), "tanukifamily.ru", "", errors.New("whois failed"))
	require.NoError(t, err)
	require.Equal(t, expected, expiration.Expiry)
	require.Equal(t, []string{"ok"}, expiration.Statuses)
}

func TestFallbackToRDAPPreservesErrorForExplicitHost(t *testing.T) {
	previous := rdapFallbackLookup
	t.Cleanup(func() {
		rdapFallbackLookup = previous
	})

	rdapCalled := false
	rdapFallbackLookup = func(_ context.Context, _ string) (client.Result, error) {
		rdapCalled = true
		return client.Result{}, nil
	}

	before := time.Now()
	expiration, err := fallbackToRDAP(context.Background(), "google.com", "whois.dot.ph", errors.New("whois failed"))
	require.ErrorContains(t, err, "whois failed")
	require.WithinDuration(t, before, expiration.Expiry, time.Second)
	require.False(t, rdapCalled)
}

func TestFallbackToRDAPCombinesErrorsWhenRdapFails(t *testing.T) {
	previous := rdapFallbackLookup
	t.Cleanup(func() {
		rdapFallbackLookup = previous
	})

	rdapFallbackLookup = func(_ context.Context, _ string) (client.Result, error) {
		return client.Result{}, errors.New("rdap failed")
	}

	_, err := fallbackToRDAP(context.Background(), "missing.ru", "", errors.New("whois failed"))
	require.ErrorContains(t, err, "whois failed")
	require.ErrorContains(t, err, "rdap fallback failed: rdap failed")
}
