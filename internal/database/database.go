package database

import (
	"database/sql"
	"fmt"
	"log"

	"test_effective_mobile/test-effective-mobile/internal/config"
	"test_effective_mobile/test-effective-mobile/internal/models"

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
func AddPerson(person *models.Person, s *sql.DB) error {
	query := `
		INSERT INTO persons (name, surname, age, gender)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, surname, age, gender
	`
	createdPerson := new(models.Person)
	row := s.QueryRow(query, person.Name, person.Surname, person.Age, person.Gender)
	err := row.Scan(
		&createdPerson.ID,
		&createdPerson.Name,
		&createdPerson.Surname,
		&createdPerson.Age,
		&createdPerson.Gender,
	)
	if err != nil {
		return err
	}
	person.ID = createdPerson.ID
	err = AddNationalitiesForUser(person, s)
	if err != nil {
		return err
	}

	return nil
}

func AddNationalitiesForUser(createPerson *models.Person, s *sql.DB) error {
	var natID int
	queryNationalities := `
		INSERT INTO nationalities (person_id,country_id, probability)
		VALUES ($1, $2,$3)
		RETURNING id
	`
	for _, country := range createPerson.Country {
		err := s.QueryRow(queryNationalities, createPerson.ID, country.CountryID, country.Probability).Scan(&natID)
		if err != nil {
			return err
		}

	}
	return nil
}
