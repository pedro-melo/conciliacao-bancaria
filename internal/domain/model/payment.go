package model

import (
	"time"
)

// Payment representa um pagamento bancário recebido no sistema
type Payment struct {
	ID          string    `json:"transaction_id"`
	BankAccount string    `json:"bank_account"`
	Amount      float64   `json:"amount"`
	PaymentDate time.Time `json:"payment_date"`
	ReferenceID *string   `json:"reference_id,omitempty"`

	// Campos adicionais para controle interno
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewPayment cria uma nova instância de Payment
func NewPayment(id, bankAccount string, amount float64, paymentDate time.Time, referenceID *string) *Payment {
	now := time.Now()

	return &Payment{
		ID:          id,
		BankAccount: bankAccount,
		Amount:      amount,
		PaymentDate: paymentDate,
		ReferenceID: referenceID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
