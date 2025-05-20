package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"mori/pkg/models"

	_ "github.com/lib/pq" // PostgreSQL driver
	migrate "github.com/rubenv/sql-migrate"

	"github.com/joho/godotenv"
)

// InitDB initializes the PostgreSQL database connection.
func InitDB() *sql.DB {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get database connection details from .env
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Run migrations
	err = Migrations(db)
	if err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	return db
}

// InitRepositories initializes all repositories with the database connection.
func InitRepositories(db *sql.DB) *models.Repositories {
	return &models.Repositories{
		UserRepo:     &UserRepository{DB: db},
		SessionRepo:  &SessionRepository{DB: db},
		GroupRepo:    &GroupRepository{DB: db},
		NotifRepo:    &NotifRepository{DB: db},
		MsgRepo:      &MsgRepository{DB: db},
		LLMConvoRepo: &LLMConvoRepository{DB: db},
	}
}

// Migrations applies database migrations.
func Migrations(db *sql.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "pkg/db/migration/PostgreSql",
	}

	// Execute migrations for PostgreSQL
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	fmt.Printf("Applied %d migrations to PostgreSQL database!\n", n)
	return nil
}
