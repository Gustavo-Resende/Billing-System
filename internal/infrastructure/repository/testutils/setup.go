package testutils

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"github.com/teusf/billing-system/config"
	"github.com/teusf/billing-system/internal/infrastructure/database"
	"go.uber.org/zap"
)

// SetupTestDB connects to the test database and returns the connection.
// It assumes a running Postgres instance configured for testing.
func SetupTestDB() (*sql.DB, error) {
	// Configuração hardcoded para o container de teste (postgres-test)
	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "5433", // Porta do container de teste
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "billing_test",
		DBSSLMode:  "disable",
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// ResetAndMigrate cleans the public schema and runs all migrations.
func ResetAndMigrate(db *sql.DB, migrationsPath string) error {
	// 1. Limpar banco de dados (Reset Schema)
	if _, err := db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO postgres; GRANT ALL ON SCHEMA public TO public;"); err != nil {
		return fmt.Errorf("erro ao resetar schema: %v", err)
	}

	// 2. Rodar Migrations
	logger := zap.NewNop()
	if err := database.RunMigrations(db, migrationsPath, logger); err != nil {
		return fmt.Errorf("erro ao rodar migrations: %v", err)
	}

	return nil
}

// NewTestTx starts a transaction and returns a rollback function.
// It requires an active database connection.
func NewTestTx(t *testing.T, db *sql.DB) (*sql.Tx, func()) {
	t.Helper()

	if db == nil {
		t.Skip("Banco de dados de teste não inicializado")
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Falha ao iniciar transação: %v", err)
	}

	cleanup := func() {
		if err := tx.Rollback(); err != nil {
			t.Logf("Erro no rollback (ignorar se a transação já foi comitada/abortada manualmente): %v", err)
		}
	}

	return tx, cleanup
}
