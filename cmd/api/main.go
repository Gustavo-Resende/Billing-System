package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/teusf/billing-system/config"
	"github.com/teusf/billing-system/internal/infrastructure/database"
	"github.com/teusf/billing-system/internal/infrastructure/logger"
)

func main() {
	// 1. Carrega Configurações
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// 2. Configura Logger
	isDebug := cfg.AppEnv == "development"
	log := logger.NewLogger(isDebug)
	defer log.Sync()

	log.Info("Starting Billing System API",
		zap.String("env", cfg.AppEnv),
		zap.String("port", cfg.AppPort),
	)

	// 3. Conecta ao banco de dados
	db, err := database.NewPostgresConnection(cfg, log)
	if err != nil {
		log.Fatal("Could not connect to database", zap.Error(err))
	}
	defer db.Close()

	// 4. Executa as migrations
	// Nota: em produção, este caminho deve ser absoluto ou relativo corretamente ao binário
	if err := database.RunMigrations(db, "internal/infrastructure/database/migrations", log); err != nil {
		log.Fatal("Failed to run migrations", zap.Error(err))
	}

	// 5. Configura Router
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health Check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 6. Inicia o servidor
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Info("Server listening", zap.String("addr", addr))

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Server failed", zap.Error(err))
	}
}
