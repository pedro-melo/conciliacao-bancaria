package request

import "time"

// BilletRequest representa a estrutura de dados para a requisição de criação ou atualização de um boleto
type BilletRequest struct {
	BilletID     string    `json:"billet_id"`
	BankAccount  string    `json:"bank_account"`
	Amount       float64   `json:"amount"`
	IssuanceDate time.Time `json:"issuance_date"`
	ReferenceID  *string   `json:"reference_id,omitempty"`
}

// BilletBatchRequest representa uma lista de boletos para processamento em lote
type BilletBatchRequest struct {
	Billets []BilletRequest `json:"billets"`
}
