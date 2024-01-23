package main

import (
	"fmt"
	"log"
	"net/http"

	"test_effective_mobile/test-effective-mobile/internal/api"
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
	a := api.New(db)
	router := api.SetupRoutes(a)
	port := ":8080"
	log.Printf("Server ... running on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, router))
	return nil
}
