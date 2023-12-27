package main

import (
	"context"
	"fmt"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

func (m *Module) ArangoInit() error {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{m.config.ArangoHost},
	})
	if err != nil {
		return err
	}
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(m.config.ArangoUser, m.config.ArangoPassword),
	})
	if err != nil {
		return fmt.Errorf(`error creating arango client: %v`, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	db, err := client.Database(ctx, m.config.ArangoDBName)
	if err != nil {
		return fmt.Errorf(`cannot connect to communication db: %v`, err)
	}
	m.db = db

	return nil
}
