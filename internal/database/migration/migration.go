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
	{
		Version:     2,
		Description: `create table nationality`,
		Script: `CREATE TABLE nationalities (
    id SERIAL PRIMARY KEY,
    person_id INT REFERENCES persons(id),
    country_id VARCHAR(10),
    probability FLOAT
    )`,
	},
	{
		Version:     3,
		Description: "rename and change type of nationality column",
		Script: `
			ALTER TABLE persons ADD COLUMN nationality_id INT;
			ALTER TABLE persons DROP COLUMN nationality;
		`,
	},
	{
		Version:     4,
		Description: "remove nationality_id",
		Script: `
			ALTER TABLE persons DROP COLUMN nationality_id;
		`,
	},
}

func Run(db *sql.DB) error {
	return darwin.New(darwin.NewGenericDriver(db, darwin.PostgresDialect{}), items, nil).Migrate()
}
