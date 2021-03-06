// +build !mem

package data

import (
	"github.com/dstpierre/gosaas/data/mongo"
	"github.com/dstpierre/gosaas/model"
)

// Open creates the database connection and initialize the MongoDB services.
func (db *DB) Open(driverName, dataSourceName string) error {
	conn, err := model.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}

	//  for mongo, we need to copy the connection session at each requests
	// this is done in our api's ServeHTTP
	db.Users = &mongo.Users{}
	db.Webhooks = &mongo.Webhooks{}

	db.Connection = conn

	db.DatabaseName = "gosaas"
	db.CopySession = true
	return nil
}
