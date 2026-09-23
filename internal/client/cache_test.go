package client

import (
	"context"
	"fmt"
	"testing"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/require"
)

type testClient struct {
	result *time.Time
}

func (f testClient) Lookup(_ context.Context, _ string, _ string) (Result, error) {
	return Result{Expiry: *f.result}, nil
}

type errTestClient struct{}

func (f errTestClient) Lookup(_ context.Context, _ string, _ string) (Result, error) {
	return Result{Expiry: time.Now()}, fmt.Errorf("failed to get domain info blah")
}

func TestCachedClient(t *testing.T) {
	ctx := context.Background()
	cache := cache.New(1*time.Minute, 1*time.Minute)
	expected := time.Now()
	domain := "foo.bar"
	host := ""

	cli := NewCachedClient(testClient{result: &expected}, cache)

	// test getting from out fake client
	t.Run("get fresh", func(t *testing.T) {
		res, err := cli.Lookup(ctx, domain, host)
		require.NoError(t, err)
		require.Equal(t, expected, res.Expiry)
	})

	// here we change the inner fake client result, but the result
	// should be the cached one
	t.Run("get from cache", func(t *testing.T) {
		oldExpected := expected
		expected = time.Now()
		res, err := cli.Lookup(ctx, domain, host)
		require.NoError(t, err)
		require.Equal(t, oldExpected, res.Expiry)
	})

	// here we flush the cache and verify that the result is the one
	// from the fake client
	t.Run("flush cache", func(t *testing.T) {
		cache.Flush()
		res, err := cli.Lookup(ctx, domain, host)
		require.NoError(t, err)
		require.Equal(t, expected, res.Expiry)
	})

	t.Run("do not cache errors", func(t *testing.T) {
		cache.Flush()

		cli := NewCachedClient(errTestClient{}, cache)
		_, err := cli.Lookup(ctx, domain, host)
		require.Error(t, err)

		_, err = cli.Lookup(ctx, domain, host)
		require.Error(t, err)

		cached, got := cache.Get(domain)
		require.Nil(t, cached)
		require.False(t, got)
	})
}
