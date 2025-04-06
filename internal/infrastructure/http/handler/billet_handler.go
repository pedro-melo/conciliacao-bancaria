package handler

import (
	"encoding/json"
	"net/http"

	"conciliacao-bancaria/internal/application/usecase"
	"conciliacao-bancaria/internal/infrastructure/http/dto/request"
	"conciliacao-bancaria/internal/infrastructure/http/dto/response"
	"conciliacao-bancaria/pkg/errors"
)

// BilletHandler gerencia as requisições HTTP relacionadas a boletos
type BilletHandler struct {
	billetUseCase *usecase.BilletUseCase
}

// NewBilletHandler cria uma nova instância do BilletHandler
func NewBilletHandler(billetUseCase *usecase.BilletUseCase) *BilletHandler {
	return &BilletHandler{
		billetUseCase: billetUseCase,
	}
}

// CreateBillet processa a requisição para criar um novo boleto
func (h *BilletHandler) CreateBillet(w http.ResponseWriter, r *http.Request) {
	var req request.BilletRequest
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

	// Criar boleto através do caso de uso
	billet, err := h.billetUseCase.CreateBillet(r.Context(), req.ToBilletDomain())
	if err != nil {
		handleError(w, err)
		return
	}

	// Converter para resposta e retornar
	resp := response.FromBilletDomain(billet)
	renderJSON(w, resp, http.StatusCreated)
}

// GetBilletByID processa a requisição para buscar um boleto por ID
func (h *BilletHandler) GetBilletByID(w http.ResponseWriter, r *http.Request) {
	// Extrair ID do boleto da URL
	billetID := extractPathParam(r, "id")
	if billetID == "" {
		http.Error(w, "ID do boleto é obrigatório", http.StatusBadRequest)
		return
	}

	// Buscar boleto através do caso de uso
	billet, err := h.billetUseCase.GetBilletByID(r.Context(), billetID)
	if err != nil {
		handleError(w, err)
		return
	}

	// Converter para resposta e retornar
	resp := response.FromBilletDomain(billet)
	renderJSON(w, resp, http.StatusOK)
}

// ListBillets processa a requisição para listar todos os boletos
func (h *BilletHandler) ListBillets(w http.ResponseWriter, r *http.Request) {
	// Extrair parâmetros de paginação e filtros (se necessário)
	params := extractQueryParams(r)

	// Buscar boletos através do caso de uso
	billets, err := h.billetUseCase.ListBillets(r.Context(), params)
	if err != nil {
		handleError(w, err)
		return
	}

	// Converter para resposta e retornar
	var resp []response.BilletResponse
	for _, billet := range billets {
		resp = append(resp, response.FromBilletDomain(billet))
	}

	renderJSON(w, resp, http.StatusOK)
}

// ImportBillets processa a requisição para importar uma lista de boletos
func (h *BilletHandler) ImportBillets(w http.ResponseWriter, r *http.Request) {
	var req []request.BilletRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Erro ao decodificar requisição: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validar cada boleto na requisição
	for i, billetReq := range req {
		if err := billetReq.Validate(); err != nil {
			http.Error(w, "Dados inválidos no boleto "+string(i)+": "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Converter requisições para domínio
	domainBillets := make([]interface{}, len(req))
	for i, billetReq := range req {
		domainBillets[i] = billetReq.ToBilletDomain()
	}

	// Importar boletos através do caso de uso
	results, err := h.billetUseCase.ImportBillets(r.Context(), domainBillets)
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

// DeleteBillet processa a requisição para excluir um boleto
func (h *BilletHandler) DeleteBillet(w http.ResponseWriter, r *http.Request) {
	// Extrair ID do boleto da URL
	billetID := extractPathParam(r, "id")
	if billetID == "" {
		http.Error(w, "ID do boleto é obrigatório", http.StatusBadRequest)
		return
	}

	// Excluir boleto através do caso de uso
	err := h.billetUseCase.DeleteBillet(r.Context(), billetID)
	if err != nil {
		handleError(w, err)
		return
	}

	// Retornar sucesso sem conteúdo
	w.WriteHeader(http.StatusNoContent)
}

// handleError trata os diversos tipos de erro e define o status HTTP adequado
func handleError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *errors.NotFoundError:
		http.Error(w, e.Error(), http.StatusNotFound)
	case *errors.ValidationError:
		http.Error(w, e.Error(), http.StatusBadRequest)
	case *errors.ConflictError:
		http.Error(w, e.Error(), http.StatusConflict)
	default:
		http.Error(w, "Erro interno do servidor: "+err.Error(), http.StatusInternalServerError)
	}
}

// renderJSON serializa uma resposta para JSON e escreve no ResponseWriter
func renderJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Erro ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// extractPathParam extrai um parâmetro da URL
func extractPathParam(r *http.Request, param string) string {
	// Esta função depende da implementação do router usado (como gorilla/mux, chi, etc.)
	// Portanto, é uma função básica que deve ser adaptada conforme o router escolhido
	return r.PathValue(param) // Usando PathValue do net/http a partir do Go 1.22
}

// extractQueryParams extrai parâmetros de consulta da URL
func extractQueryParams(r *http.Request) map[string]string {
	params := make(map[string]string)

	// Extrair parâmetros comuns como paginação, ordenação, etc.
	query := r.URL.Query()

	// Exemplo de parâmetros que podem ser úteis para listagem
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

	return params
}
