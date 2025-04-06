package handler

import (
	"encoding/json"
	"net/http"

	"conciliacao-bancaria/internal/application/usecase"
	"conciliacao-bancaria/internal/infrastructure/http/dto/request"
	"conciliacao-bancaria/internal/infrastructure/http/dto/response"
)

// ReconciliationHandler gerencia as requisições HTTP relacionadas à conciliação
type ReconciliationHandler struct {
	reconciliationUseCase *usecase.ReconciliationUseCase
}

// NewReconciliationHandler cria uma nova instância do ReconciliationHandler
func NewReconciliationHandler(reconciliationUseCase *usecase.ReconciliationUseCase) *ReconciliationHandler {
	return &ReconciliationHandler{
		reconciliationUseCase: reconciliationUseCase,
	}
}

// RunReconciliation processa a requisição para executar o processo de conciliação
func (h *ReconciliationHandler) RunReconciliation(w http.ResponseWriter, r *http.Request) {
	var req request.ReconciliationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Erro ao decodificar requisição: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validar requisição
	if err := req.Validate(); err != nil {
		http.Error(w, "Dados inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Executar conciliação através do caso de uso
	result, err := h.reconciliationUseCase.RunReconciliation(r.Context(), req.ToReconciliationParams())
	if err != nil {
		handleError(w, err)
		return
	}

	// Converter resultado para a estrutura de resposta conforme requisito 3.a
	resp := response.ReconciliationResultResponse{
		BoletosConciliados:    make([]response.BilletReconciliationResponse, 0),
		BoletosNaoConciliados: make([]response.BilletResponse, 0),
	}

	// Preencher boletos conciliados
	for _, reconciled := range result.ReconciledBillets {
		resp.BoletosConciliados = append(resp.BoletosConciliados, response.FromBilletReconciliationDomain(reconciled))
	}

	// Preencher boletos não conciliados
	for _, notReconciled := range result.NotReconciledBillets {
		resp.BoletosNaoConciliados = append(resp.BoletosNaoConciliados, response.FromBilletDomain(notReconciled))
	}

	renderJSON(w, resp, http.StatusOK)
}

// GetReconciliationByID processa a requisição para obter detalhes de uma conciliação específica
func (h *ReconciliationHandler) GetReconciliationByID(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da conciliação da URL
	reconciliationID := extractPathParam(r, "id")
	if reconciliationID == "" {
		http.Error(w, "ID da conciliação é obrigatório", http.StatusBadRequest)
		return
	}

	// Buscar conciliação através do caso de uso
	reconciliation, err := h.reconciliationUseCase.GetReconciliationByID(r.Context(), reconciliationID)
	if err != nil {
		handleError(w, err)
		return
	}

	// Converter para resposta e retornar
	resp := response.FromReconciliationDomain(reconciliation)
	renderJSON(w, resp, http.StatusOK)
}

// ListReconciliations processa a requisição para listar todas as conciliações
func (h *ReconciliationHandler) ListReconciliations(w http.ResponseWriter, r *http.Request) {
	// Extrair parâmetros de paginação e filtros
	params := extractReconciliationQueryParams(r)

	// Buscar conciliações através do caso de uso
	reconciliations, err := h.reconciliationUseCase.ListReconciliations(r.Context(), params)
	if err != nil {
		handleError(w, err)
		return
	}

	// Converter para resposta e retornar
	var resp []response.ReconciliationSummaryResponse
	for _, reconciliation := range reconciliations {
		resp = append(resp, response.FromReconciliationSummaryDomain(reconciliation))
	}

	renderJSON(w, resp, http.StatusOK)
}

// GetBilletReconciliationStatus processa a requisição para obter o status de conciliação de um boleto específico
func (h *ReconciliationHandler) GetBilletReconciliationStatus(w http.ResponseWriter, r *http.Request) {
	// Extrair ID do boleto da URL
	billetID := extractPathParam(r, "billet_id")
	if billetID == "" {
		http.Error(w, "ID do boleto é obrigatório", http.StatusBadRequest)
		return
	}

	// Buscar status de conciliação através do caso de uso
	status, err := h.reconciliationUseCase.GetBilletReconciliationStatus(r.Context(), billetID)
	if err != nil {
		handleError(w, err)
		return
	}

	// Converter para resposta e retornar
	resp := response.BilletReconciliationStatusResponse{
		BilletID:           status.BilletID,
		ReconciliationID:   status.ReconciliationID,
		TransactionID:      status.TransactionID,
		Status:             status.Status,
		Strategy:           status.Strategy,
		AmountDiff:         status.AmountDiff,
		ReconciliationDate: status.ReconciliationDate,
	}

	renderJSON(w, resp, http.StatusOK)
}

// GetPaymentReconciliationStatus processa a requisição para obter o status de conciliação de um pagamento específico
func (h *ReconciliationHandler) GetPaymentReconciliationStatus(w http.ResponseWriter, r *http.Request) {
	// Extrair ID do pagamento da URL
	paymentID := extractPathParam(r, "transaction_id")
	if paymentID == "" {
		http.Error(w, "ID do pagamento é obrigatório", http.StatusBadRequest)
		return
	}

	// Buscar status de conciliação através do caso de uso
	status, err := h.reconciliationUseCase.GetPaymentReconciliationStatus(r.Context(), paymentID)
	if err != nil {
		handleError(w, err)
		return
	}

	// Converter para resposta e retornar
	resp := response.PaymentReconciliationStatusResponse{
		TransactionID:      status.TransactionID,
		ReconciliationID:   status.ReconciliationID,
		BilletID:           status.BilletID,
		Status:             status.Status,
		Strategy:           status.Strategy,
		AmountDiff:         status.AmountDiff,
		ReconciliationDate: status.ReconciliationDate,
	}

	renderJSON(w, resp, http.StatusOK)
}

// GetReconciliationStatistics processa a requisição para obter estatísticas de conciliação
func (h *ReconciliationHandler) GetReconciliationStatistics(w http.ResponseWriter, r *http.Request) {
	// Extrair parâmetros de filtro para as estatísticas
	params := extractReconciliationQueryParams(r)

	// Buscar estatísticas através do caso de uso
	stats, err := h.reconciliationUseCase.GetReconciliationStatistics(r.Context(), params)
	if err != nil {
		handleError(w, err)
		return
	}

	// Converter para resposta e retornar
	resp := response.ReconciliationStatisticsResponse{
		TotalBillets:                stats.TotalBillets,
		TotalPayments:               stats.TotalPayments,
		TotalReconciledBillets:      stats.TotalReconciledBillets,
		TotalNotReconciledBillets:   stats.TotalNotReconciledBillets,
		TotalMatchedByReferenceID:   stats.TotalMatchedByReferenceID,
		TotalMatchedByAccountAmount: stats.TotalMatchedByAccountAmount,
		TotalWithAmountDifference:   stats.TotalWithAmountDifference,
		AverageAmountDifference:     stats.AverageAmountDifference,
		ReconciliationRate:          stats.ReconciliationRate,
	}

	renderJSON(w, resp, http.StatusOK)
}

// extractReconciliationQueryParams extrai parâmetros de consulta específicos para conciliação
func extractReconciliationQueryParams(r *http.Request) map[string]string {
	params := make(map[string]string)

	// Extrair parâmetros comuns como paginação, ordenação, etc.
	query := r.URL.Query()

	// Parâmetros específicos para conciliação
	if limit := query.Get("limit"); limit != "" {
		params["limit"] = limit
	}

	if offset := query.Get("offset"); offset != "" {
		params["offset"] = offset
	}

	if startDate := query.Get("start_date"); startDate != "" {
		params["start_date"] = startDate
	}

	if endDate := query.Get("end_date"); endDate != "" {
		params["end_date"] = endDate
	}

	if bankAccount := query.Get("bank_account"); bankAccount != "" {
		params["bank_account"] = bankAccount
	}

	if status := query.Get("status"); status != "" {
		params["status"] = status
	}

	if strategy := query.Get("strategy"); strategy != "" {
		params["strategy"] = strategy
	}

	if tolerancePercentage := query.Get("tolerance_percentage"); tolerancePercentage != "" {
		params["tolerance_percentage"] = tolerancePercentage
	}

	return params
}
