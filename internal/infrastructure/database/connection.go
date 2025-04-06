package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

// DBConfig representa a configuração de conexão com o banco de dados
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Connection representa uma conexão com o banco de dados
type Connection struct {
	DB *sql.DB
}

// NewConnection cria uma nova conexão com o banco de dados
func NewConnection() (*Connection, error) {
	config := DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "conciliacao"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir conexão com o banco de dados: %w", err)
	}

	// Configurar o pool de conexões
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verificar se a conexão está funcionando
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("falha ao conectar no banco de dados: %w", err)
	}

	log.Println("Conexão com o banco de dados estabelecida com sucesso")
	return &Connection{DB: db}, nil
}

// Close fecha a conexão com o banco de dados
func (c *Connection) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

// getEnv retorna o valor da variável de ambiente ou um valor padrão
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
