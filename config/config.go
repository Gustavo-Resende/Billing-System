package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	// App
	AppEnv   string `mapstructure:"APP_ENV"`
	AppPort  string `mapstructure:"APP_PORT"`
	LogLevel string `mapstructure:"LOG_LEVEL"`

	// Database
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSL_MODE"`

	// RabbitMQ
	RabbitMQHost     string `mapstructure:"RABBITMQ_HOST"`
	RabbitMQPort     string `mapstructure:"RABBITMQ_PORT"`
	RabbitMQUser     string `mapstructure:"RABBITMQ_USER"`
	RabbitMQPassword string `mapstructure:"RABBITMQ_PASSWORD"`

	// Evolution API
	EvolutionAPIURL   string `mapstructure:"EVOLUTION_API_URL"`
	EvolutionAPIKey   string `mapstructure:"EVOLUTION_API_KEY"`
	EvolutionInstance string `mapstructure:"EVOLUTION_INSTANCE"`

	// Business Rules
	LembreteDiasAntes  int    `mapstructure:"LEMBRETE_DIAS_ANTES"`
	HorarioInicioEnvio string `mapstructure:"HORARIO_INICIO_ENVIO"`
	HorarioFimEnvio    string `mapstructure:"HORARIO_FIM_ENVIO"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("env")

	// Default values
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LEMBRETE_DIAS_ANTES", 3)

	viper.AutomaticEnv() // Read from env variables

	if err := viper.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist, we fallback to env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}
