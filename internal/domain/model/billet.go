package model

import (
	"time"
)

// Billet representa um boleto emitido no sistema
type Billet struct {
	ID           string    `json:"billet_id"`
	BankAccount  string    `json:"bank_account"`
	Amount       float64   `json:"amount"`
	IssuanceDate time.Time `json:"issuance_date"`
	ReferenceID  *string   `json:"reference_id,omitempty"`

	// Campos adicionais para controle interno
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewBillet cria uma nova instância de Billet
func NewBillet(id, bankAccount string, amount float64, issuanceDate time.Time, referenceID *string) *Billet {
	now := time.Now()

	return &Billet{
		ID:           id,
		BankAccount:  bankAccount,
		Amount:       amount,
		IssuanceDate: issuanceDate,
		ReferenceID:  referenceID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
