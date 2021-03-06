// +build mem

package data

import (
	"github.com/dstpierre/gosaas/data/mem"
	"github.com/dstpierre/gosaas/model"
)

// Open initializes in-memory services.
func (db *DB) Open(driverName, dataSourceName string) error {
	conn, err := model.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}

	db.Users = &mem.Users{}
	db.Webhooks = &mem.Webhooks{}

	// we use this to populate test data for unit test
	db.CopySession = true

	db.Connection = &conn
	return nil
}
