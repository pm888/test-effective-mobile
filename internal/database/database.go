package database

import (
	"database/sql"
	"fmt"
	"log"

	"test_effective_mobile/test-effective-mobile/internal/config"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	pgGoodMsg = "postgres  ...  starts"
	pgBadMsg  = "error connect BD"
)

func ConnectPostgresql(c config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable", c.DbUser, c.DbPassword, c.DbHost, c.DbPort, c.DbName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.Wrap(err, pgBadMsg)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, pgBadMsg)
	}
	log.Println(pgGoodMsg)
	return db, nil
}
func DeletePersonByID(id int, s *sql.DB) error {
	_, err := s.Exec("DELETE FROM persons WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
