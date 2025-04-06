package request

import "time"

// PaymentRequest representa a estrutura de dados para a requisição de criação ou atualização de um pagamento
type PaymentRequest struct {
	TransactionID string    `json:"transaction_id"`
	BankAccount   string    `json:"bank_account"`
	Amount        float64   `json:"amount"`
	PaymentDate   time.Time `json:"payment_date"`
	ReferenceID   *string   `json:"reference_id,omitempty"`
}

// PaymentBatchRequest representa uma lista de pagamentos para processamento em lote
type PaymentBatchRequest struct {
	Payments []PaymentRequest `json:"payments"`
}
