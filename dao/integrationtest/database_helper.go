package integrationtest

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func openDB(ctx context.Context, mysqlContainer *mysql.MySQLContainer) (*sql.DB, error) {

	connStr, err := mysqlContainer.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}
	log.Println("MySQL connection string: ", connStr)

	db, err := sql.Open("mysql", connStr+"?parseTime=true&loc=UTC")
	return db, err
}

func startMySQLContainer() (*mysql.MySQLContainer, error) {
	ctx := context.Background()

	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0.36",
		mysql.WithDatabase("example"),
		mysql.WithUsername("appuser"),
		mysql.WithPassword("passme"),
		mysql.WithScripts(filepath.Join("testdata", "schema.sql")),
	)

	if err != nil {
		log.Printf("failed to start container: %s\n", err)
		return nil, err
	}

	return mysqlContainer, nil
}

func cleanUp(db *sql.DB, mysqlContainer *mysql.MySQLContainer) {
	log.Println("closing DB handler")
	if err := db.Close(); err != nil {
		log.Println("failed to close database", err)
	}

	log.Println("shutting down MySQL container")
	if err := testcontainers.TerminateContainer(mysqlContainer); err != nil {
		log.Printf("failed to terminate container: %s\n", err)
	}
}
