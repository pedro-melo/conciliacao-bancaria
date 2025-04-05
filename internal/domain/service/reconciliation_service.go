package service

import (
	"context"
	"math"
	"time"

	"conciliacao-bancaria/internal/domain/model"
)

// TolerancePercentage define a tolerância percentual para diferença de valores (5%)
const TolerancePercentage = 5.0

// ReconciliationService define as operações de serviço para conciliação
type ReconciliationService interface {
	// ReconcileBilletsWithPayments realiza a conciliação entre boletos e pagamentos
	ReconcileBilletsWithPayments(ctx context.Context, billets []*model.Billet, payments []*model.Payment) (*model.ReconciliationResult, error)

	// GetReconciliationStatus recupera o status de conciliação de um boleto
	GetReconciliationStatus(ctx context.Context, billetID string) (*model.Reconciliation, error)
}

// DefaultReconciliationService implementa ReconciliationService
type DefaultReconciliationService struct {
	// Dependências podem ser adicionadas aqui
}

// NewReconciliationService cria uma nova instância de DefaultReconciliationService
func NewReconciliationService() ReconciliationService {
	return &DefaultReconciliationService{}
}

// ReconcileBilletsWithPayments realiza a conciliação entre boletos e pagamentos
func (s *DefaultReconciliationService) ReconcileBilletsWithPayments(
	ctx context.Context,
	billets []*model.Billet,
	payments []*model.Payment,
) (*model.ReconciliationResult, error) {
	// Mapa para acompanhar boletos já conciliados
	reconciledBilletsMap := make(map[string]bool)

	// Mapa para acompanhar pagamentos já utilizados
	usedPaymentsMap := make(map[string]bool)

	result := &model.ReconciliationResult{
		ReconciledBillets:    []model.ReconciledBillet{},
		NonReconciledBillets: []model.Billet{},
	}

	// 1ª Estratégia: Conciliação por reference_id
	s.reconcileByReferenceID(billets, payments, reconciledBilletsMap, usedPaymentsMap, &result.ReconciledBillets)

	// 2ª Estratégia: Conciliação por conta, valor e data
	s.reconcileByAccountValueDate(billets, payments, reconciledBilletsMap, usedPaymentsMap, &result.ReconciledBillets)

	// Adicionar boletos não conciliados
	for _, billet := range billets {
		if !reconciledBilletsMap[billet.ID] {
			result.NonReconciledBillets = append(result.NonReconciledBillets, *billet)
		}
	}

	return result, nil
}

// GetReconciliationStatus recupera o status de conciliação de um boleto
func (s *DefaultReconciliationService) GetReconciliationStatus(ctx context.Context, billetID string) (*model.Reconciliation, error) {
	// Implementação completa seria feita na camada de aplicação com acesso ao repositório
	return nil, nil
}

// reconcileByReferenceID implementa a 1ª estratégia de conciliação
func (s *DefaultReconciliationService) reconcileByReferenceID(
	billets []*model.Billet,
	payments []*model.Payment,
	reconciledBilletsMap map[string]bool,
	usedPaymentsMap map[string]bool,
	reconciledBillets *[]model.ReconciledBillet,
) {
	// Mapear pagamentos por referenceID para acesso rápido
	paymentsByReferenceID := make(map[string]*model.Payment)
	for _, payment := range payments {
		if payment.ReferenceID != nil && *payment.ReferenceID != "" && !usedPaymentsMap[payment.ID] {
			paymentsByReferenceID[*payment.ReferenceID] = payment
		}
	}

	// Tentar conciliar boletos pelo referenceID
	for _, billet := range billets {
		// Pular boletos já conciliados
		if reconciledBilletsMap[billet.ID] {
			continue
		}

		// Verificar se o boleto tem referenceID válido
		if billet.ReferenceID == nil || *billet.ReferenceID == "" {
			continue
		}

		// Verificar se existe um pagamento com o mesmo referenceID
		payment, found := paymentsByReferenceID[*billet.ReferenceID]
		if !found {
			continue
		}

		// Calcular diferença de valor
		amountDiff := math.Abs(payment.Amount - billet.Amount)
		amountDiffPercentage := (amountDiff / billet.Amount) * 100

		// Determinar status de conciliação
		var status model.ConciliationStatus
		if amountDiff == 0 {
			status = model.StatusSuccessful
		} else if amountDiffPercentage <= TolerancePercentage {
			status = model.StatusDifferentValue
		} else {
			// Se a diferença de valor for muito grande, não concilia por referenceID
			continue
		}

		// Adicionar à lista de boletos conciliados
		*reconciledBillets = append(*reconciledBillets, model.ReconciledBillet{
			BilletID:             billet.ID,
			BankAccount:          billet.BankAccount,
			TransactionID:        payment.ID,
			ConciliationStatus:   status,
			ConciliationStrategy: model.StrategyReferenceID,
			ReferenceID:          billet.ReferenceID,
			AmountDiff:           amountDiff,
		})

		// Marcar boleto e pagamento como utilizados
		reconciledBilletsMap[billet.ID] = true
		usedPaymentsMap[payment.ID] = true
	}
}

// reconcileByAccountValueDate implementa a 2ª estratégia de conciliação
func (s *DefaultReconciliationService) reconcileByAccountValueDate(
	billets []*model.Billet,
	payments []*model.Payment,
	reconciledBilletsMap map[string]bool,
	usedPaymentsMap map[string]bool,
	reconciledBillets *[]model.ReconciledBillet,
) {
	// Para cada pagamento não utilizado
	for _, payment := range payments {
		if usedPaymentsMap[payment.ID] {
			continue
		}

		var bestBillet *model.Billet
		var minDateDiff time.Duration = time.Duration(math.MaxInt64)
		var bestAmountDiff float64 = math.MaxFloat64

		// Procurar o melhor boleto para este pagamento
		for _, billet := range billets {
			// Pular boletos já conciliados
			if reconciledBilletsMap[billet.ID] {
				continue
			}

			// Verificar se conta bancária corresponde
			if billet.BankAccount != payment.BankAccount {
				continue
			}

			// Calcular diferença de valor
			amountDiff := math.Abs(payment.Amount - billet.Amount)
			amountDiffPercentage := (amountDiff / billet.Amount) * 100

			// Verificar se está dentro da tolerância
			if amountDiffPercentage > TolerancePercentage {
				continue
			}

			// Calcular diferença de data
			dateDiff := payment.PaymentDate.Sub(billet.IssuanceDate)
			if dateDiff < 0 {
				dateDiff = -dateDiff
			}

			// Critérios para escolher o melhor boleto:
			// 1. Priorizar a menor diferença de data
			// 2. Em caso de empate, priorizar a menor diferença de valor
			// 3. Em caso de empate, priorizar o boleto mais antigo
			isBetter := false

			if bestBillet == nil {
				isBetter = true
			} else if dateDiff < minDateDiff {
				isBetter = true
			} else if dateDiff == minDateDiff && amountDiff < bestAmountDiff {
				isBetter = true
			} else if dateDiff == minDateDiff && amountDiff == bestAmountDiff && billet.IssuanceDate.Before(bestBillet.IssuanceDate) {
				isBetter = true
			}

			if isBetter {
				bestBillet = billet
				minDateDiff = dateDiff
				bestAmountDiff = amountDiff
			}
		}

		// Se encontrou um boleto para conciliar
		if bestBillet != nil {
			// Determinar status de conciliação
			var status model.ConciliationStatus
			if bestAmountDiff == 0 {
				status = model.StatusSuccessful
			} else {
				status = model.StatusDifferentValue
			}

			// Adicionar à lista de boletos conciliados
			*reconciledBillets = append(*reconciledBillets, model.ReconciledBillet{
				BilletID:             bestBillet.ID,
				BankAccount:          bestBillet.BankAccount,
				TransactionID:        payment.ID,
				ConciliationStatus:   status,
				ConciliationStrategy: model.StrategyAccountAmountDate,
				ReferenceID:          bestBillet.ReferenceID,
				AmountDiff:           bestAmountDiff,
			})

			// Marcar boleto e pagamento como utilizados
			reconciledBilletsMap[bestBillet.ID] = true
			usedPaymentsMap[payment.ID] = true
		}
	}
}
