package client

import (
	"context"
)

type multiClient []Client

func (clients multiClient) Lookup(ctx context.Context, domain string, host string) (Result, error) {
	var res Result
	var err error
	for _, client := range clients {
		res, err = client.Lookup(ctx, domain, host)
		if err == nil {
			break
		}
	}
	return res, err
}

// NewMultiClient returns a client that wraps multiple clients.
// It returns the first success, or, if all clients fail, the latest failure.
func NewMultiClient(clients ...Client) Client {
	return multiClient(clients)
}
