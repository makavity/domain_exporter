package refresher

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/caarlos0/domain_exporter/internal/client"
	"github.com/caarlos0/domain_exporter/internal/safeconfig"
)

type fakeOk struct{}

func (fakeOk) Lookup(ctx context.Context, domain string, host string) (client.Result, error) {
	return client.Result{}, nil
}

type fakeFail struct{}

func (fakeFail) Lookup(ctx context.Context, domain string, host string) (client.Result, error) {
	return client.Result{}, errors.New("foo")
}

func Test_refresher_Refresh(t *testing.T) {
	tests := []struct {
		name      string
		refresher Refresher
	}{
		{
			name:      "refresh is ok",
			refresher: New(time.Second, fakeOk{}, time.Second, safeconfig.Domain{Name: "foo.com", Host: ""}),
		},
		{
			name:      "refresh is failed",
			refresher: New(time.Second, fakeFail{}, time.Second, safeconfig.Domain{Name: "foo.com", Host: ""}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			tt.refresher.Refresh(context.Background())
			defer tt.refresher.Stop()
		})
	}
}
