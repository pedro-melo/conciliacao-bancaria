package repository

import (
	"context"

	"conciliacao-bancaria/internal/domain/model"
)

// BilletRepository define as operações de repositório para boletos
type BilletRepository interface {
	// Create persiste um novo boleto no banco de dados
	Create(ctx context.Context, billet *model.Billet) error

	// CreateMany persiste múltiplos boletos no banco de dados
	CreateMany(ctx context.Context, billets []*model.Billet) error

	// GetByID recupera um boleto pelo seu ID
	GetByID(ctx context.Context, id string) (*model.Billet, error)

	// GetAll recupera todos os boletos
	GetAll(ctx context.Context) ([]*model.Billet, error)

	// GetByBankAccount recupera boletos por conta bancária
	GetByBankAccount(ctx context.Context, bankAccount string) ([]*model.Billet, error)

	// GetByReferenceID recupera boletos por ID de referência
	GetByReferenceID(ctx context.Context, referenceID string) ([]*model.Billet, error)

	// Update atualiza um boleto existente
	Update(ctx context.Context, billet *model.Billet) error

	// Delete remove um boleto pelo ID
	Delete(ctx context.Context, id string) error

	// FindNonReconciled encontra boletos que ainda não foram conciliados
	FindNonReconciled(ctx context.Context) ([]*model.Billet, error)
}
