package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/teusf/billing-system/config"
	"github.com/teusf/billing-system/internal/infrastructure/database"
	"go.uber.org/zap"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	// 1. Setup global do banco de testes
	var err error
	testDB, err = setupTestDB()
	if err != nil {
		log.Printf("SKIP: Banco de dados de teste não disponível: %v", err)
		os.Exit(0)
	}
	defer testDB.Close()

	// 2. Limpar banco de dados (Reset Schema)
	if _, err := testDB.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO postgres; GRANT ALL ON SCHEMA public TO public;"); err != nil {
		log.Fatalf("Erro ao resetar schema do banco de teste: %v", err)
	}

	// 3. Rodar Migrations para garantir schema atualizado
	logger := zap.NewNop() // Logger desabilitado para testes
	if err := database.RunMigrations(testDB, "../database/migrations", logger); err != nil {
		log.Fatalf("Erro ao rodar migrations no banco de teste: %v", err)
	}

	// 3. Rodar os testes
	code := m.Run()

	os.Exit(code)
}

func setupTestDB() (*sql.DB, error) {
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

// newTestTx inicia uma transação e retorna uma função de limpeza (rollback)
func newTestTx(t *testing.T) (*sql.Tx, func()) {
	t.Helper()

	if testDB == nil {
		t.Skip("Banco de dados de teste não inicializado")
	}

	tx, err := testDB.Begin()
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
