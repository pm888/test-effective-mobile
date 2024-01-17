package migrations

import (
	"database/sql"

	"github.com/GuiaBolso/darwin"
)

var items = []darwin.Migration{
	{
		Version:     1,
		Description: `create table persons`,
		Script: `CREATE TABLE persons (
		id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        surname VARCHAR(255) NOT NULL,
        age INT,
        gender VARCHAR(10),
        nationality VARCHAR(50)
	)`,
	},
}

func Run(db *sql.DB) error {
	return darwin.New(darwin.NewGenericDriver(db, darwin.PostgresDialect{}), items, nil).Migrate()
}
