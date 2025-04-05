package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

// DatabaseConfig contém as configurações para conexão com o banco de dados
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Connection gerencia a conexão com o banco de dados
type Connection struct {
	db *sql.DB
}

// NewConnection cria uma nova instância de conexão com o banco de dados
func NewConnection(config DatabaseConfig) (*Connection, error) {
	// Construir a string de conexão
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	// Abrir conexão com o banco de dados
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}

	// Configurar pool de conexões
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verificar se a conexão está funcionando
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}

	log.Println("Conexão com o banco de dados estabelecida com sucesso")
	return &Connection{db: db}, nil
}

// GetDB retorna a instância do banco de dados
func (c *Connection) GetDB() *sql.DB {
	return c.db
}

// Close fecha a conexão com o banco de dados
func (c *Connection) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// ExecuteTransaction executa uma função dentro de uma transação
func (c *Connection) ExecuteTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			// Garantir que a transação seja revertida em caso de pânico
			_ = tx.Rollback()
			panic(p) // Re-lançar o pânico após o rollback
		}
	}()

	// Executar a função passada dentro da transação
	if err := fn(tx); err != nil {
		// Reverter em caso de erro
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("erro na transação: %v, rollback falhou: %w", err, rbErr)
		}
		return err
	}

	// Commit da transação se tudo estiver OK
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao fazer commit da transação: %w", err)
	}

	return nil
}
