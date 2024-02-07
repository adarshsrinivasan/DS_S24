package sql

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"

	_ "github.com/lib/pq"
)

type connPool struct {
	mu                    sync.Mutex
	maxConns              int
	numberOfActiveClients int
}

// getClient retrieves a free client from the pool.
func (p *connPool) getClient(ctx context.Context, applicationName, schemaName string) (*clientObj, error) {
	for {
		if p.numberOfActiveClients < p.maxConns {
			client, err := newClient(ctx, applicationName, schemaName)
			if err != nil {
				err = fmt.Errorf("exception while creating new client: %v", err)
				logrus.Errorf("getClient: %v", err)
				return nil, err
			}
			p.mu.Lock()
			p.numberOfActiveClients++
			p.mu.Unlock()
			return client, nil
		} else {
			logrus.Errorf("DB connection pool max capacity. Waiting for free connection")
		}
	}
}

func (p *connPool) close(ctx context.Context, client *clientObj) error {
	p.mu.Lock()
	p.numberOfActiveClients--
	p.mu.Unlock()
	return client.bunClient.Close()
}
