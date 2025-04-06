package handler

import (
	"encoding/json"
	"net/http"

	"conciliacao-bancaria/internal/application/usecase"
	"conciliacao-bancaria/internal/infrastructure/http/dto/request"
	"conciliacao-bancaria/internal/infrastructure/http/dto/response"
)

// PaymentHandler gerencia as requisições HTTP relacionadas a pagamentos
type PaymentHandler struct {
	paymentUseCase *usecase.PaymentUseCase
}

// NewPaymentHandler cria uma nova instância do PaymentHandler
func NewPaymentHandler(paymentUseCase *usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{
		paymentUseCase: paymentUseCase,
	}
}

// CreatePayment processa a requisição para criar um novo pagamento
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req request.PaymentRequest
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

	// Criar pagamento através do caso de uso
	payment, err := h.paymentUseCase.CreatePayment(r.Context(), req.ToPaymentDomain())
	if err != nil {
		handleError(w, err)
		return
	}

	// Converter para resposta e retornar
	resp := response.FromPaymentDomain(payment)
	renderJSON(w, resp, http.StatusCreated)
}

// GetPaymentByID processa a requisição para buscar um pagamento por ID
func (h *PaymentHandler) GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	// Extrair ID do pagamento da URL
	paymentID := extractPathParam(r, "id")
	if paymentID == "" {
		http.Error(w, "ID do pagamento é obrigatório", http.StatusBadRequest)
		return
	}

	// Buscar pagamento através do caso de uso
	payment, err := h.paymentUseCase.GetPaymentByID(r.Context(), paymentID)
	if err != nil {
		handleError(w, err)
		return
	}

	// Converter para resposta e retornar
	resp := response.FromPaymentDomain(payment)
	renderJSON(w, resp, http.StatusOK)
}

// ListPayments processa a requisição para listar todos os pagamentos
func (h *PaymentHandler) ListPayments(w http.ResponseWriter, r *http.Request) {
	// Extrair parâmetros de paginação e filtros
	params := extractPaymentQueryParams(r)

	// Buscar pagamentos através do caso de uso
	payments, err := h.paymentUseCase.ListPayments(r.Context(), params)
	if err != nil {
		handleError(w, err)
		return
	}

	// Converter para resposta e retornar
	var resp []response.PaymentResponse
	for _, payment := range payments {
		resp = append(resp, response.FromPaymentDomain(payment))
	}

	renderJSON(w, resp, http.StatusOK)
}

// ImportPayments processa a requisição para importar uma lista de pagamentos
func (h *PaymentHandler) ImportPayments(w http.ResponseWriter, r *http.Request) {
	var req []request.PaymentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Erro ao decodificar requisição: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validar cada pagamento na requisição
	for i, paymentReq := range req {
		if err := paymentReq.Validate(); err != nil {
			http.Error(w, "Dados inválidos no pagamento "+string(i)+": "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Converter requisições para domínio
	domainPayments := make([]interface{}, len(req))
	for i, paymentReq := range req {
		domainPayments[i] = paymentReq.ToPaymentDomain()
	}

	// Importar pagamentos através do caso de uso
	results, err := h.paymentUseCase.ImportPayments(r.Context(), domainPayments)
	if err != nil {
		handleError(w, err)
		return
	}

	// Converter para resposta e retornar
	var resp struct {
		Imported int      `json:"imported"`
		Errors   []string `json:"errors,omitempty"`
	}
	resp.Imported = results.Imported
	resp.Errors = results.Errors

	renderJSON(w, resp, http.StatusOK)
}

// GetPaymentsByBankAccount processa a requisição para buscar pagamentos por conta bancária
func (h *PaymentHandler) GetPaymentsByBankAccount(w http.ResponseWriter, r *http.Request) {
	// Extrair conta bancária da URL
	bankAccount := extractPathParam(r, "bank_account")
	if bankAccount == "" {
		http.Error(w, "Conta bancária é obrigatória", http.StatusBadRequest)
		return
	}

	// Buscar pagamentos através do caso de uso
	payments, err := h.paymentUseCase.GetPaymentsByBankAccount(r.Context(), bankAccount)
	if err != nil {
		handleError(w, err)
		return
	}

	// Converter para resposta e retornar
	var resp []response.PaymentResponse
	for _, payment := range payments {
		resp = append(resp, response.FromPaymentDomain(payment))
	}

	renderJSON(w, resp, http.StatusOK)
}

// GetPaymentsByReferenceID processa a requisição para buscar pagamentos por referenceID
func (h *PaymentHandler) GetPaymentsByReferenceID(w http.ResponseWriter, r *http.Request) {
	// Extrair referenceID da URL
	referenceID := extractPathParam(r, "reference_id")
	if referenceID == "" {
		http.Error(w, "ID de referência é obrigatório", http.StatusBadRequest)
		return
	}

	// Buscar pagamentos através do caso de uso
	payments, err := h.paymentUseCase.GetPaymentsByReferenceID(r.Context(), referenceID)
	if err != nil {
		handleError(w, err)
		return
	}

	// Converter para resposta e retornar
	var resp []response.PaymentResponse
	for _, payment := range payments {
		resp = append(resp, response.FromPaymentDomain(payment))
	}

	renderJSON(w, resp, http.StatusOK)
}

// DeletePayment processa a requisição para excluir um pagamento
func (h *PaymentHandler) DeletePayment(w http.ResponseWriter, r *http.Request) {
	// Extrair ID do pagamento da URL
	paymentID := extractPathParam(r, "id")
	if paymentID == "" {
		http.Error(w, "ID do pagamento é obrigatório", http.StatusBadRequest)
		return
	}

	// Excluir pagamento através do caso de uso
	err := h.paymentUseCase.DeletePayment(r.Context(), paymentID)
	if err != nil {
		handleError(w, err)
		return
	}

	// Retornar sucesso sem conteúdo
	w.WriteHeader(http.StatusNoContent)
}

// extractPaymentQueryParams extrai parâmetros de consulta específicos para pagamentos
func extractPaymentQueryParams(r *http.Request) map[string]string {
	params := make(map[string]string)

	// Extrair parâmetros comuns como paginação, ordenação, etc.
	query := r.URL.Query()

	// Exemplo de parâmetros que podem ser úteis para listagem de pagamentos
	if limit := query.Get("limit"); limit != "" {
		params["limit"] = limit
	}

	if offset := query.Get("offset"); offset != "" {
		params["offset"] = offset
	}

	if bankAccount := query.Get("bank_account"); bankAccount != "" {
		params["bank_account"] = bankAccount
	}

	if minAmount := query.Get("min_amount"); minAmount != "" {
		params["min_amount"] = minAmount
	}

	if maxAmount := query.Get("max_amount"); maxAmount != "" {
		params["max_amount"] = maxAmount
	}

	if startDate := query.Get("start_date"); startDate != "" {
		params["start_date"] = startDate
	}

	if endDate := query.Get("end_date"); endDate != "" {
		params["end_date"] = endDate
	}

	if referenceID := query.Get("reference_id"); referenceID != "" {
		params["reference_id"] = referenceID
	}

	if transactionID := query.Get("transaction_id"); transactionID != "" {
		params["transaction_id"] = transactionID
	}

	return params
}
