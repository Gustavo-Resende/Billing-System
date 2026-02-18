package configuracao

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
	defer testDB.Close()

	code := m.Run()
	testDB.Close()
	os.Exit(code)
}

func TestConfiguracaoPostgres_Upsert(t *testing.T) {
	tx, cleanup := testutils.NewTestTx(t, testDB)
	defer cleanup()

	repo := NewConfiguracaoPostgres(tx)

	// 1. Create
	c1, _ := entity.NewConfiguracao("user1")
	c1.DiasAntesLembrete = 5

	err := repo.Save(c1)
	assert.NoError(t, err)

	// 2. Read
	found, err := repo.FindByUsuarioID("user1")
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, 5, found.DiasAntesLembrete)

	// 3. Update (Upsert)
	// Vamos criar uma nova entidade com o MESMO usuarioID para testar o ON CONFLICT
	c2, _ := entity.NewConfiguracao("user1")
	// NAO copiamos o ID. Deixamos gerar um novo.
	// O Upsert deve detectar conflito no usuario_id e atualizar o registro existente (que tem o ID do c1)

	c2.DiasAntesLembrete = 3
	c2.TemplateLembrete = "Novo template"

	err = repo.Save(c2)
	assert.NoError(t, err)

	found2, _ := repo.FindByUsuarioID("user1")
	assert.Equal(t, 3, found2.DiasAntesLembrete)
	assert.Equal(t, "Novo template", found2.TemplateLembrete)
	assert.Equal(t, c1.ID, found2.ID) // O ID deve ser o original (c1), não o do c2
}
