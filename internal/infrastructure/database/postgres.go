package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/teusf/billing-system/config"
	"go.uber.org/zap"
)

func NewPostgresConnection(cfg *config.Config, logger *zap.Logger) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	// Máscara de senha para log
	safeConnStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=*** dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBSSLMode,
	)
	logger.Debug("Tentando conectar ao banco", zap.String("dsn", safeConnStr))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexao com banco: %w", err)
	}

	// Configuração do Pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verificar conexão com retry
	for i := 0; i < 5; i++ {
		err = db.Ping()
		if err == nil {
			break
		}

		// Se o erro for "banco de dados não existe" (código 3D000 ou mensagem)
		if strings.Contains(err.Error(), "3D000") || strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "nao existe") {
			logger.Warn("Banco de dados nao existe (detectado 3D000). Tentando criar...", zap.String("db", cfg.DBName))
			if createErr := createDatabase(cfg, logger); createErr != nil {
				return nil, fmt.Errorf("erro ao criar banco de dados: %w", createErr)
			}
			// Tenta pingar novamente na próxima iteração
			continue
		}

		logger.Warn("Falha ao pingar banco, retentando em 2s...", zap.Error(err))
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco apos tentativas: %w", err)
	}

	logger.Info("Conectado ao PostgreSQL com sucesso")

	return db, nil
}

func createDatabase(cfg *config.Config, logger *zap.Logger) error {
	// Conecta ao banco 'postgres' (padrão) para criar o nosso
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("erro ao conectar ao banco postgres para criar db: %w", err)
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	if err != nil {
		return fmt.Errorf("erro ao executar CREATE DATABASE: %w", err)
	}

	logger.Info("Banco de dados criado com sucesso", zap.String("db", cfg.DBName))
	return nil
}

// RunMigrations executa os arquivos .sql na ordem correta
func RunMigrations(db *sql.DB, migrationsPath string, logger *zap.Logger) error {
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("erro ao ler diretorio de migrations: %w", err)
	}

	// Ordenar por nome (001, 002, etc)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		logger.Info("Verificando migration", zap.String("file", file.Name()))

		content, err := os.ReadFile(filepath.Join(migrationsPath, file.Name()))
		if err != nil {
			return fmt.Errorf("erro ao ler arquivo %s: %w", file.Name(), err)
		}

		// Divide o arquivo em comandos separados se necessário, mas para este projeto
		// assumimos que cada arquivo pode rodar inteiro (separado por ;)
		// ou é um comando só. Executar tudo de uma vez costuma funcionar para DDL simples.
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("erro ao executar migration %s: %w", file.Name(), err)
		}
	}

	logger.Info("Todas as migrations verificadas/executadas com sucesso")
	return nil
}
