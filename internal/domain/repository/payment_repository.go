package repository

import (
	"context"

	"conciliacao-bancaria/internal/domain/model"
)

// PaymentRepository define as operações de repositório para pagamentos
type PaymentRepository interface {
	// Create persiste um novo pagamento no banco de dados
	Create(ctx context.Context, payment *model.Payment) error

	// CreateMany persiste múltiplos pagamentos no banco de dados
	CreateMany(ctx context.Context, payments []*model.Payment) error

	// GetByID recupera um pagamento pelo seu ID
	GetByID(ctx context.Context, id string) (*model.Payment, error)

	// GetAll recupera todos os pagamentos
	GetAll(ctx context.Context) ([]*model.Payment, error)

	// GetByBankAccount recupera pagamentos por conta bancária
	GetByBankAccount(ctx context.Context, bankAccount string) ([]*model.Payment, error)

	// GetByReferenceID recupera pagamentos por ID de referência
	GetByReferenceID(ctx context.Context, referenceID string) ([]*model.Payment, error)

	// Update atualiza um pagamento existente
	Update(ctx context.Context, payment *model.Payment) error

	// Delete remove um pagamento pelo ID
	Delete(ctx context.Context, id string) error

	// FindByBankAccountAndAmount encontra pagamentos por conta bancária e valor aproximado
	FindByBankAccountAndAmount(ctx context.Context, bankAccount string, amount float64, tolerance float64) ([]*model.Payment, error)
}
