package mensagem

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/teusf/billing-system/internal/domain/entity"
	"github.com/teusf/billing-system/internal/infrastructure/repository/cliente"
	"github.com/teusf/billing-system/internal/infrastructure/repository/fatura"
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

func TestMensagemPostgres_CRUD(t *testing.T) {
	tx, cleanup := testutils.NewTestTx(t, testDB)
	defer cleanup()

	// Setup deps - Repositórios precisam usar a transação (tx)
	cRepo := cliente.NewClientePostgres(tx)
	client, _ := entity.NewCliente("Cliente 1", "5511999998888", "c1@test.com")
	if err := cRepo.Save(client); err != nil {
		t.Fatalf("Failed to save client: %v", err)
	}

	fRepo := fatura.NewFaturaPostgres(tx)
	// Usa data fixa para teste consistente
	vencimento := time.Now().AddDate(0, 0, 5)
	fatura, _ := entity.NewFatura(client.ID, 100, vencimento, "F1")
	if err := fRepo.Save(fatura); err != nil {
		t.Fatalf("Failed to save fatura: %v", err)
	}

	repo := NewMensagemPostgres(tx)

	// 1. Create
	msg, err := entity.NewMensagem(fatura.ID, client.ID, client.WhatsApp, "Ola", entity.TipoMensagemLembrete)
	assert.NoError(t, err)

	err = repo.Save(msg)
	assert.NoError(t, err)

	// 2. Read
	found, err := repo.FindByID(msg.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, msg.Conteudo, found.Conteudo)

	// 3. Update (Enviada)
	msg.MarcarComoEnviada()
	err = repo.Update(msg)
	assert.NoError(t, err)

	found2, _ := repo.FindByID(msg.ID)
	assert.Equal(t, entity.StatusMensagemEnviada, found2.Status)
	assert.NotNil(t, found2.EnviadoEm)

	// 4. FindByStatus
	list, err := repo.FindByStatus(entity.StatusMensagemEnviada)
	assert.NoError(t, err)
	// Como limpamos o banco, deve ser 1. Se tiver sujeira, >= 1
	assert.GreaterOrEqual(t, len(list), 1)

	// 5. FindParaDLQ (Fail simulation)
	msgFalha, _ := entity.NewMensagem(fatura.ID, client.ID, client.WhatsApp, "Fail", entity.TipoMensagemCobranca)
	msgFalha.MarcarComoFalha("Erro 1")
	msgFalha.MarcarComoFalha("Erro 2")
	msgFalha.MarcarComoFalha("Erro 3")
	msgFalha.MarcarComoFalha("Erro 4")
	msgFalha.MarcarComoFalha("Erro 5") // Total 5 falhas

	repo.Save(msgFalha)

	dlq, err := repo.FindParaDLQ()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(dlq), 1) // deve achar msgFalha
}
