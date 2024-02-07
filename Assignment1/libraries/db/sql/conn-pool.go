package sql

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"

	_ "github.com/lib/pq"
)

type connPool struct {
	mu       sync.Mutex
	clients  chan *clientObj
	maxConns int
}

func (p *connPool) initialize(ctx context.Context, applicationName, schemaName string, maxConns int) error {
	p.clients = make(chan *clientObj, maxConns)
	p.maxConns = maxConns

	for i := 0; i < maxConns; i++ {
		client, err := newClient(ctx, applicationName, schemaName)
		if err != nil {
			err = fmt.Errorf("exception while creating new client: %v", err)
			logrus.Errorf("initialize: %v", err)
			return err
		}
		p.clients <- client
	}
	return nil
}

// getClient retrieves a free client from the pool.
func (p *connPool) getClient(ctx context.Context) *clientObj {
	for {
		select {
		case client := <-p.clients:
			return client
		}
	}
}

func (p *connPool) close(ctx context.Context, client *clientObj) {
	p.clients <- client
}
