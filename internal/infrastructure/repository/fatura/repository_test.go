package fatura

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/teusf/billing-system/internal/domain/entity"
	"github.com/teusf/billing-system/internal/infrastructure/repository/cliente"
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

func TestFaturaPostgres_CRUD(t *testing.T) {
	tx, cleanup := testutils.NewTestTx(t, testDB)
	defer cleanup()

	// Setup dependency
	cRepo := cliente.NewClientePostgres(tx)
	client, _ := entity.NewCliente("Cliente 1", "5511999998888", "c1@test.com")
	cRepo.Save(client)

	repo := NewFaturaPostgres(tx)

	// 1. Create
	vencimento := time.Now().AddDate(0, 0, 5)
	f, _ := entity.NewFatura(client.ID, 150.00, vencimento, "Consultoria")
	err := repo.Save(f)
	assert.NoError(t, err)

	// 2. Read
	found, err := repo.FindByID(f.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, f.Numero, found.Numero)
	assert.Equal(t, f.Valor, found.Valor)

	// 3. Update (Pagar)
	f.MarcarComoPaga()
	err = repo.Update(f)
	assert.NoError(t, err)

	found2, _ := repo.FindByID(f.ID)
	assert.Equal(t, entity.StatusPaga, found2.Status)
	assert.NotNil(t, found2.DataPagamento)

	// 4. FindByClienteID
	list, err := repo.FindByClienteID(client.ID)
	assert.NoError(t, err)
	assert.Len(t, list, 1) // Deve ter 1 fatura
}

func TestFaturaPostgres_Filtros(t *testing.T) {
	tx, cleanup := testutils.NewTestTx(t, testDB)
	defer cleanup()

	cRepo := cliente.NewClientePostgres(tx)
	client, err := entity.NewCliente("Cliente 1", "5511888888888", "c1@test.com")
	assert.NoError(t, err)
	cRepo.Save(client)

	repo := NewFaturaPostgres(tx)

	// Fatura 1: Vence hoje (pendente)
	// Criamos com data futura para passar na validação do NewFatura, depois forçamos para "Agora"
	// para garantir que o teste de "Vencendo Hoje" funcione mesmo se rodar às 23:59
	f1, err := entity.NewFatura(client.ID, 100, time.Now().Add(24*time.Hour), "Hoje")
	assert.NoError(t, err)
	if f1 == nil {
		t.Fatal("Falha ao criar f1: nil")
	}
	f1.DataVencimento = time.Now() // Força ser "Hoje" para o teste
	repo.Save(f1)

	// Fatura 2: Vence em 3 dias (pendente)
	f2, err := entity.NewFatura(client.ID, 200, time.Now().AddDate(0, 0, 3), "Futuro")
	assert.NoError(t, err)
	if f2 == nil {
		t.Fatal("Falha ao criar f2: nil")
	}
	repo.Save(f2)

	// Fatura 3: Paga
	f3, err := entity.NewFatura(client.ID, 300, time.Now().AddDate(0, 0, 5), "Paga")
	assert.NoError(t, err)
	if f3 == nil {
		t.Fatal("Falha ao criar f3: nil")
	}
	f3.MarcarComoPaga()
	repo.Save(f3)

	// Test FindPendentes
	pendentes, err := repo.FindPendentes()
	assert.NoError(t, err)
	// Como o banco de teste é limpo por transação, deve ter 2 (f1, f2)
	// Se tiver sujeira de outros testes, pode falhar. Mas rollback garante limpeza.
	assert.Len(t, pendentes, 2)

	// Test FindVencendoEm(0) -> Hoje
	hoje, err := repo.FindVencendoEm(0)
	assert.NoError(t, err)
	assert.Len(t, hoje, 1)
	if len(hoje) > 0 {
		assert.Equal(t, f1.ID, hoje[0].ID)
	}

	// Test FindVencendoEm(3)
	tresDias, err := repo.FindVencendoEm(3)
	assert.NoError(t, err)
	assert.Len(t, tresDias, 1)
	if len(tresDias) > 0 {
		assert.Equal(t, f2.ID, tresDias[0].ID)
	}
}
