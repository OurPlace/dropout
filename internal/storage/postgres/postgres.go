package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ourplace/dropout/internal/settings"
	"github.com/ourplace/dropout/internal/storage"

	// lib/pq is our postgres driver
	_ "github.com/lib/pq"
)

var _ storage.Datastore = (*Database)(nil)

// Database ...
type Database struct {
	*sql.DB
	config settings.Database

	table  string
	query  string
	driver string
}

// ReportDropout reports that we dropped out.
func (db *Database) ReportDropout(ctx context.Context) error {
	// Report drop out to our database.
	_, err := db.ExecContext(ctx, db.query)
	if err != nil {
		return err
	}

	return nil
}

// Open creates a new database connection and pings it.
func Open(ctx context.Context, config settings.Database) (*Database, error) {
	var (
		// Set defaults.
		database = &Database{
			config: config,
			table:  config.Table,
			driver: config.Driver,
		}
		err error
	)

	// By default, our table will be located in `network.dropout`.
	// But this location can be overridden.
	if database.table == "" {
		database.table = "network.dropout"
	}

	// The ability to change the driver is so that way we can easily test this
	// code against go-sqlmock. In practice, we normally won't ever need to do
	// this, but we may use the `config.Driver` option to specify the type of
	// database interface we plan on using so we can swap easily.
	if database.driver == "" {
		database.driver = "postgres"
	}

	// Create the query once.
	database.query = fmt.Sprintf(
		"INSERT INTO %s (at) VALUES (current_timestamp);",
		database.table,
	)

	// Connect to our database.
	database.DB, err = sql.Open(database.driver, config.DSN)
	if err != nil {
		return nil, err
	}

	// Ensure that it's actually there.
	if err = database.DB.PingContext(ctx); err != nil {
		return nil, err
	}

	// Finally done.
	return database, nil
}
