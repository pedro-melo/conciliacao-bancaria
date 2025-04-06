package response

import "time"

// PaymentResponse representa a estrutura de dados para a resposta de um pagamento
type PaymentResponse struct {
	TransactionID string    `json:"transaction_id"`
	BankAccount   string    `json:"bank_account"`
	Amount        float64   `json:"amount"`
	PaymentDate   time.Time `json:"payment_date"`
	ReferenceID   *string   `json:"reference_id,omitempty"`
	Status        string    `json:"status"`              // Status atual do pagamento (recebido, conciliado, estornado, etc.)
	BilletID      *string   `json:"billet_id,omitempty"` // ID do boleto relacionado, se conciliado
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PaymentListResponse representa uma lista paginada de pagamentos para resposta
type PaymentListResponse struct {
	Payments    []PaymentResponse `json:"payments"`
	TotalCount  int64             `json:"total_count"`
	PageSize    int               `json:"page_size"`
	CurrentPage int               `json:"current_page"`
	TotalPages  int               `json:"total_pages"`
}
