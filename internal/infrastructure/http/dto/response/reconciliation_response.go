package response

import "time"

// ReconciliationItemResponse representa um item conciliado na resposta da API
type ReconciliationItemResponse struct {
	BilletID             string    `json:"billet_id"`
	TransactionID        string    `json:"transaction_id"`
	BankAccount          string    `json:"bank_account"`
	ConciliationStatus   string    `json:"conciliation_status"`    // conciliado_com_sucesso, valor_diferente
	ConciliationStrategy string    `json:"conciliation_strategy"`  // reference_id, conta_valor_data
	AmountDiff           float64   `json:"amount_diff"`            // Diferença de valor (se houver)
	ReferenceID          *string   `json:"reference_id,omitempty"` // Quando utilizado na conciliação
	ReconciliationDate   time.Time `json:"reconciliation_date"`    // Data da conciliação
}

// NonReconciledBilletResponse representa um boleto não conciliado na resposta da API
type NonReconciledBilletResponse struct {
	BilletID     string    `json:"billet_id"`
	BankAccount  string    `json:"bank_account"`
	Amount       float64   `json:"amount"`
	IssuanceDate time.Time `json:"issuance_date"`
	ReferenceID  *string   `json:"reference_id,omitempty"`
}

// ReconciliationResponse representa a estrutura de dados para a resposta de uma conciliação
type ReconciliationResponse struct {
	ReconciliationID      string                        `json:"reconciliation_id"`
	ReconciliationDate    time.Time                     `json:"reconciliation_date"`
	BoletosConciliados    []ReconciliationItemResponse  `json:"boletos_conciliados"`
	BoletosNaoConciliados []NonReconciledBilletResponse `json:"boletos_nao_conciliados"`
	TotalConciliados      int                           `json:"total_conciliados"`
	TotalNaoConciliados   int                           `json:"total_nao_conciliados"`
	Tolerance             float64                       `json:"tolerance"`
}

// ReconciliationHistoryResponse representa o histórico de conciliações para um boleto ou pagamento específico
type ReconciliationHistoryResponse struct {
	EntityID              string                      `json:"entity_id"`   // Pode ser billet_id ou transaction_id
	EntityType            string                      `json:"entity_type"` // "boleto" ou "pagamento"
	CurrentStatus         string                      `json:"current_status"`
	ReconciliationHistory []ReconciliationHistoryItem `json:"reconciliation_history"`
}

// ReconciliationHistoryItem representa um item do histórico de conciliação
type ReconciliationHistoryItem struct {
	ReconciliationID     string    `json:"reconciliation_id"`
	ReconciliationDate   time.Time `json:"reconciliation_date"`
	Status               string    `json:"status"`
	PairedWith           string    `json:"paired_with,omitempty"` // ID do boleto ou transação com o qual foi pareado
	ConciliationStrategy string    `json:"conciliation_strategy,omitempty"`
	AmountDiff           float64   `json:"amount_diff,omitempty"`
}

// ReconciliationListResponse representa uma lista paginada de conciliações para resposta
type ReconciliationListResponse struct {
	Reconciliations []ReconciliationSummary `json:"reconciliations"`
	TotalCount      int64                   `json:"total_count"`
	PageSize        int                     `json:"page_size"`
	CurrentPage     int                     `json:"current_page"`
	TotalPages      int                     `json:"total_pages"`
}

// ReconciliationSummary representa um resumo de uma conciliação para listagem
type ReconciliationSummary struct {
	ReconciliationID    string    `json:"reconciliation_id"`
	ReconciliationDate  time.Time `json:"reconciliation_date"`
	TotalProcessed      int       `json:"total_processed"`
	TotalConciliados    int       `json:"total_conciliados"`
	TotalNaoConciliados int       `json:"total_nao_conciliados"`
	Tolerance           float64   `json:"tolerance"`
}
