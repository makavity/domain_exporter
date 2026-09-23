package client

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type clifail int

func (clifail) Lookup(_ context.Context, domain string, host string) (Result, error) {
	return Result{}, errors.New("foo")
}

type clisuccess time.Time

func (c clisuccess) Lookup(_ context.Context, domain string, host string) (Result, error) {
	return Result{Expiry: time.Time(c)}, nil
}

func TestMulti(t *testing.T) {
	ctx := context.Background()
	t.Run("first client succeed", func(t *testing.T) {
		expected := time.Now()
		expire, err := NewMultiClient(clisuccess(expected), clifail(0)).Lookup(ctx, "a", "")
		require.NoError(t, err)
		require.Equal(t, expected, expire.Expiry)
	})
	t.Run("last client succeed", func(t *testing.T) {
		expected := time.Now()
		expire, err := NewMultiClient(clifail(0), clifail(0), clisuccess(expected)).Lookup(ctx, "a", "")
		require.NoError(t, err)
		require.Equal(t, expected, expire.Expiry)
	})
	t.Run("no client succeed", func(t *testing.T) {
		expire, err := NewMultiClient(clifail(0), clifail(0), clifail(0)).Lookup(ctx, "a", "")
		require.EqualError(t, err, "foo")
		require.Equal(t, expire.Expiry, time.Time{})
	})
}

func TestEPPStatuses(t *testing.T) {
	require.Equal(t, []string{"clientTransferProhibited", "ok", "serverHold", "redemptionPeriod"}, EPPStatuses([]string{
		"clientTransferProhibited https://icann.org/epp#clientTransferProhibited",
		"ok - normal state.",
		"active",
		"server hold",
		"redemption period",
		"REGISTERED, DELEGATED, VERIFIED",
		"connect",
	}))
	require.Nil(t, EPPStatuses(nil))
}
