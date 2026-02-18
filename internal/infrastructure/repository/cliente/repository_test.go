package cliente

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teusf/billing-system/internal/domain/entity"
	"github.com/teusf/billing-system/internal/infrastructure/repository/testutils"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = testutils.SetupTestDB()
	if err != nil {
		log.Fatalf("Falha ao configurar banco de teste: %v", err)
	}
	// defer testDB.Close() // Removed defer

	if err := testutils.ResetAndMigrate(testDB, "../../database/migrations"); err != nil {
		log.Fatalf("Falha nas migrações: %v", err)
	}

	code := m.Run()
	testDB.Close() // Moved Close here
	os.Exit(code)
}

func TestClientePostgres_CRUD(t *testing.T) {
	// Inicia transação que será revertida no final (Rollback)
	tx, cleanup := testutils.NewTestTx(t, testDB)
	defer cleanup()

	// Repositório usa a transação, não o banco direto
	repo := NewClientePostgres(tx)

	// 1. Create - Usa nome válido (>3 chars) para não dar erro
	client, err := entity.NewCliente("Cliente 1", "5511999998888", "c1@test.com")
	assert.NoError(t, err)
	if client == nil {
		t.Fatal("Cliente nil")
	}

	err = repo.Save(client)
	assert.NoError(t, err)

	// 2. Read
	found, err := repo.FindByID(client.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, client.Nome, found.Nome)
	assert.Equal(t, client.WhatsApp, found.WhatsApp)

	foundZap, err := repo.FindByWhatsApp(client.WhatsApp)
	assert.NoError(t, err)
	assert.NotNil(t, foundZap)
	assert.Equal(t, client.ID, foundZap.ID)

	// 3. Update
	client.Nome = "Jane Doe"
	client.Ativo = false
	err = repo.Update(client)
	assert.NoError(t, err)

	found2, _ := repo.FindByID(client.ID)
	assert.Equal(t, "Jane Doe", found2.Nome)
	assert.False(t, found2.Ativo)

	// 4. FindAll
	all, err := repo.FindAll()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(all), 1)

	// 5. Delete
	err = repo.Delete(client.ID)
	assert.NoError(t, err)

	found3, err := repo.FindByID(client.ID)
	assert.NoError(t, err)
	assert.Nil(t, found3)
}
