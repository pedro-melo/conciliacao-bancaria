package usecase

import (
	"context"
	"fmt"
	"time"

	"conciliacao-bancaria/internal/domain/model"
	"conciliacao-bancaria/internal/domain/repository"
	"conciliacao-bancaria/pkg/errors"
)

// BilletUseCase implementa os casos de uso relacionados a boletos
type BilletUseCase struct {
	billetRepository repository.BilletRepository
}

// NewBilletUseCase cria uma nova instância do BilletUseCase
func NewBilletUseCase(billetRepo repository.BilletRepository) *BilletUseCase {
	return &BilletUseCase{
		billetRepository: billetRepo,
	}
}

// ImportResult representa o resultado de uma operação de importação em lote
type ImportResult struct {
	Imported int      `json:"imported"`
	Errors   []string `json:"errors,omitempty"`
}

// CreateBillet cria um novo boleto
func (uc *BilletUseCase) CreateBillet(ctx context.Context, billet *model.Billet) (*model.Billet, error) {
	// Validar dados do boleto
	if err := validateBillet(billet); err != nil {
		return nil, err
	}

	// Verificar se já existe um boleto com o mesmo ID
	existingBillet, err := uc.billetRepository.GetByID(ctx, billet.BilletID)
	if err != nil && !errors.IsNotFoundError(err) {
		return nil, errors.NewDatabaseError("verificar existência", err)
	}

	if existingBillet != nil {
		return nil, errors.NewConflictError("boleto", billet.BilletID, "boleto com este ID já existe")
	}

	// Criar boleto no repositório
	createdBillet, err := uc.billetRepository.Create(ctx, billet)
	if err != nil {
		return nil, errors.NewDatabaseError("criar", err)
	}

	return createdBillet, nil
}

// GetBilletByID busca um boleto pelo ID
func (uc *BilletUseCase) GetBilletByID(ctx context.Context, billetID string) (*model.Billet, error) {
	if billetID == "" {
		return nil, errors.NewValidationError("billet_id", "ID do boleto não pode ser vazio")
	}

	billet, err := uc.billetRepository.GetByID(ctx, billetID)
	if err != nil {
		return nil, err
	}

	return billet, nil
}

// ListBillets lista boletos com base em parâmetros de filtro
func (uc *BilletUseCase) ListBillets(ctx context.Context, params map[string]string) ([]*model.Billet, error) {
	// Criar filtro com base nos parâmetros
	filter := createBilletFilter(params)

	// Buscar boletos no repositório
	billets, err := uc.billetRepository.List(ctx, filter)
	if err != nil {
		return nil, errors.NewDatabaseError("listar", err)
	}

	return billets, nil
}

// ImportBillets importa uma lista de boletos
func (uc *BilletUseCase) ImportBillets(ctx context.Context, billetsData []interface{}) (*ImportResult, error) {
	result := &ImportResult{
		Imported: 0,
		Errors:   []string{},
	}

	// Converter e validar cada boleto
	billets := make([]*model.Billet, 0, len(billetsData))
	for i, data := range billetsData {
		billet, ok := data.(*model.Billet)
		if !ok {
			result.Errors = append(result.Errors,
				"erro na conversão do item "+string(i)+": formato inválido")
			continue
		}

		if err := validateBillet(billet); err != nil {
			result.Errors = append(result.Errors,
				"erro na validação do boleto "+billet.BilletID+": "+err.Error())
			continue
		}

		billets = append(billets, billet)
	}

	// Salvar boletos válidos no repositório
	for _, billet := range billets {
		_, err := uc.billetRepository.Create(ctx, billet)
		if err != nil {
			if errors.IsConflictError(err) {
				// Caso já exista, apenas ignoramos ou atualizamos
				// Neste caso, estamos decidindo por ignorar boletos duplicados
				result.Errors = append(result.Errors,
					"boleto "+billet.BilletID+" já existe e foi ignorado")
			} else {
				result.Errors = append(result.Errors,
					"erro ao salvar boleto "+billet.BilletID+": "+err.Error())
			}
			continue
		}

		result.Imported++
	}

	return result, nil
}

// UpdateBillet atualiza um boleto existente
func (uc *BilletUseCase) UpdateBillet(ctx context.Context, billet *model.Billet) (*model.Billet, error) {
	// Validar dados do boleto
	if err := validateBillet(billet); err != nil {
		return nil, err
	}

	// Verificar se o boleto existe
	existingBillet, err := uc.billetRepository.GetByID(ctx, billet.ID)
	if err != nil {
		return nil, err
	}

	// Se o boleto já estiver conciliado, não pode ser alterado
	if existingBillet.ReconciliationID != "" {
		return nil, errors.NewValidationError("", "boleto já conciliado não pode ser alterado")
	}

	// Atualizar boleto no repositório
	updatedBillet, err := uc.billetRepository.Update(ctx, billet)
	if err != nil {
		return nil, errors.NewDatabaseError("atualizar", err)
	}

	return updatedBillet, nil
}

// DeleteBillet remove um boleto pelo ID
func (uc *BilletUseCase) DeleteBillet(ctx context.Context, billetID string) error {
	if billetID == "" {
		return errors.NewValidationError("billet_id", "ID do boleto não pode ser vazio")
	}

	// Verificar se o boleto existe
	billet, err := uc.billetRepository.GetByID(ctx, billetID)
	if err != nil {
		return err
	}

	// Se o boleto já estiver conciliado, não pode ser excluído
	if billet.ReconciliationID != "" {
		return errors.NewValidationError("", "boleto conciliado não pode ser excluído")
	}

	// Excluir boleto do repositório
	if err := uc.billetRepository.Delete(ctx, billetID); err != nil {
		return errors.NewDatabaseError("excluir", err)
	}

	return nil
}

// validateBillet valida os dados de um boleto
func validateBillet(billet *model.Billet) error {
	if billet == nil {
		return errors.NewValidationError("", "boleto não pode ser nulo")
	}

	if billet.BilletID == "" {
		return errors.NewValidationError("billet_id", "ID do boleto é obrigatório")
	}

	if billet.BankAccount == "" {
		return errors.NewValidationError("bank_account", "conta bancária é obrigatória")
	}

	if billet.Amount <= 0 {
		return errors.NewValidationError("amount", "valor deve ser maior que zero")
	}

	// Verificar se a data de emissão é válida (não nula e não futura)
	if billet.IssuanceDate.IsZero() {
		return errors.NewValidationError("issuance_date", "data de emissão é obrigatória")
	}

	// Não permitir datas futuras
	if billet.IssuanceDate.After(time.Now()) {
		return errors.NewValidationError("issuance_date", "data de emissão não pode ser futura")
	}

	return nil
}

// createBilletFilter cria um filtro para busca de boletos com base nos parâmetros
func createBilletFilter(params map[string]string) *model.BilletFilter {
	filter := &model.BilletFilter{}

	// Aplicar filtros de parâmetros
	if bankAccount, ok := params["bank_account"]; ok {
		filter.BankAccount = bankAccount
	}

	if referenceID, ok := params["reference_id"]; ok {
		filter.ReferenceID = referenceID
	}

	// Filtros de data
	if startDateStr, ok := params["start_date"]; ok {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			filter.StartDate = &startDate
		}
	}

	if endDateStr, ok := params["end_date"]; ok {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			filter.EndDate = &endDate
		}
	}

	// Filtros de valor
	if minAmountStr, ok := params["min_amount"]; ok {
		var minAmount float64
		if _, err := fmt.Sscanf(minAmountStr, "%f", &minAmount); err == nil {
			filter.MinAmount = &minAmount
		}
	}

	if maxAmountStr, ok := params["max_amount"]; ok {
		var maxAmount float64
		if _, err := fmt.Sscanf(maxAmountStr, "%f", &maxAmount); err == nil {
			filter.MaxAmount = &maxAmount
		}
	}

	// Filtros de paginação
	if limitStr, ok := params["limit"]; ok {
		var limit int64
		if _, err := fmt.Sscanf(limitStr, "%d", &limit); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if offsetStr, ok := params["offset"]; ok {
		var offset int64
		if _, err := fmt.Sscanf(offsetStr, "%d", &offset); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	return filter
}
