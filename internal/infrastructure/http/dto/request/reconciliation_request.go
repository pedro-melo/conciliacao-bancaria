package request

import "time"

// ReconciliationRequest representa a estrutura de dados para solicitar uma conciliação
type ReconciliationRequest struct {
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	FilterAccounts []string  `json:"filter_accounts,omitempty"`
	Tolerance      *float64  `json:"tolerance,omitempty"` // Tolerância para conciliação com valor diferente (padrão 5%)
}

// ReconciliationByIDsRequest representa a solicitação de conciliação para conjuntos específicos de boletos e pagamentos
type ReconciliationByIDsRequest struct {
	BilletIDs      []string `json:"billet_ids"`
	TransactionIDs []string `json:"transaction_ids"`
	Tolerance      *float64 `json:"tolerance,omitempty"` // Tolerância para conciliação com valor diferente (padrão 5%)
}
