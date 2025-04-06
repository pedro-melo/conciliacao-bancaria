package response

import "time"

// BilletResponse representa a estrutura de dados para a resposta de um boleto
type BilletResponse struct {
	BilletID      string    `json:"billet_id"`
	BankAccount   string    `json:"bank_account"`
	Amount        float64   `json:"amount"`
	IssuanceDate  time.Time `json:"issuance_date"`
	ReferenceID   *string   `json:"reference_id,omitempty"`
	Status        string    `json:"status"`                   // Status atual do boleto (emitido, conciliado, cancelado, etc.)
	TransactionID *string   `json:"transaction_id,omitempty"` // ID da transação relacionada, se conciliado
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// BilletListResponse representa uma lista paginada de boletos para resposta
type BilletListResponse struct {
	Billets     []BilletResponse `json:"billets"`
	TotalCount  int64            `json:"total_count"`
	PageSize    int              `json:"page_size"`
	CurrentPage int              `json:"current_page"`
	TotalPages  int              `json:"total_pages"`
}
