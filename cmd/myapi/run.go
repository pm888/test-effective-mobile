package main

import (
	"fmt"
	"log"

	"test_effective_mobile/test-effective-mobile/internal/config"
	"test_effective_mobile/test-effective-mobile/internal/database"
	migrations "test_effective_mobile/test-effective-mobile/internal/database/migration"
)

const migrationMsgOK = "migration ... ok"

func Run() error {
	c, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	db, err := database.ConnectPostgresql(*c)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	defer db.Close()

	err = migrations.Run(db)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	} else {
		log.Println(migrationMsgOK)
	}
	return nil
}
