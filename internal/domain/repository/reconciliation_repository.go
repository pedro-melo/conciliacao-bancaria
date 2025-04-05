package repository

import (
	"context"

	"conciliacao-bancaria/internal/domain/model"
)

// ReconciliationRepository define as operações de repositório para conciliações
type ReconciliationRepository interface {
	// Create persiste uma nova conciliação no banco de dados
	Create(ctx context.Context, reconciliation *model.Reconciliation) error

	// CreateMany persiste múltiplas conciliações no banco de dados
	CreateMany(ctx context.Context, reconciliations []*model.Reconciliation) error

	// GetByID recupera uma conciliação pelo seu ID
	GetByID(ctx context.Context, id string) (*model.Reconciliation, error)

	// GetAll recupera todas as conciliações
	GetAll(ctx context.Context) ([]*model.Reconciliation, error)

	// GetByBilletID recupera conciliações por ID do boleto
	GetByBilletID(ctx context.Context, billetID string) ([]*model.Reconciliation, error)

	// GetByTransactionID recupera conciliações por ID da transação
	GetByTransactionID(ctx context.Context, transactionID string) ([]*model.Reconciliation, error)

	// Update atualiza uma conciliação existente
	Update(ctx context.Context, reconciliation *model.Reconciliation) error

	// Delete remove uma conciliação pelo ID
	Delete(ctx context.Context, id string) error

	// GetReconciliationHistory recupera o histórico de conciliações para auditoria
	GetReconciliationHistory(ctx context.Context, billetID string) ([]*model.Reconciliation, error)
}
