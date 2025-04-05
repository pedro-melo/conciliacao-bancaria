package model

import (
	"time"
)

// ConciliationStatus define os possíveis status de uma conciliação
type ConciliationStatus string

const (
	StatusSuccessful     ConciliationStatus = "conciliado_com_sucesso"
	StatusDifferentValue ConciliationStatus = "valor_diferente"
	StatusNotReconciled  ConciliationStatus = "nao_conciliado"
)

// ConciliationStrategy define as estratégias possíveis de conciliação
type ConciliationStrategy string

const (
	StrategyReferenceID       ConciliationStrategy = "reference_id"
	StrategyAccountAmountDate ConciliationStrategy = "conta_valor_data"
)

// Reconciliation representa o resultado da conciliação entre boleto e pagamento
type Reconciliation struct {
	ID                   string               `json:"id"`
	BilletID             string               `json:"billet_id"`
	TransactionID        *string              `json:"transaction_id,omitempty"`
	BankAccount          string               `json:"bank_account"`
	ConciliationStatus   ConciliationStatus   `json:"conciliation_status"`
	ConciliationStrategy ConciliationStrategy `json:"conciliation_strategy"`
	AmountDiff           float64              `json:"amount_diff"`
	ReferenceID          *string              `json:"reference_id,omitempty"`

	// Campos adicionais
	ReconciliationDate time.Time `json:"reconciliation_date"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// NewReconciliation cria uma nova instância de Reconciliation
func NewReconciliation(
	billetID string,
	transactionID *string,
	bankAccount string,
	status ConciliationStatus,
	strategy ConciliationStrategy,
	amountDiff float64,
	referenceID *string,
) *Reconciliation {
	now := time.Now()

	return &Reconciliation{
		ID:                   generateUUID(),
		BilletID:             billetID,
		TransactionID:        transactionID,
		BankAccount:          bankAccount,
		ConciliationStatus:   status,
		ConciliationStrategy: strategy,
		AmountDiff:           amountDiff,
		ReferenceID:          referenceID,
		ReconciliationDate:   now,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
}

// generateUUID é uma função auxiliar para gerar um UUID
// Em uma implementação real, você usaria uma biblioteca para gerar UUIDs
func generateUUID() string {
	// Implementação simplificada para exemplo
	return "rec-" + time.Now().Format("20060102150405")
}

// Definindo o modelo para resposta de reconciliação
type ReconciliationResult struct {
	ReconciledBillets    []ReconciledBillet `json:"boletos_conciliados"`
	NonReconciledBillets []Billet           `json:"boletos_nao_conciliados"`
}

// ReconciledBillet representa um boleto que foi conciliado com um pagamento
type ReconciledBillet struct {
	BilletID             string               `json:"billet_id"`
	BankAccount          string               `json:"bank_account"`
	TransactionID        string               `json:"transaction_id"`
	ConciliationStatus   ConciliationStatus   `json:"conciliation_status"`
	ConciliationStrategy ConciliationStrategy `json:"conciliation_strategy"`
	ReferenceID          *string              `json:"reference_id,omitempty"`
	AmountDiff           float64              `json:"amount_diff"`
}
