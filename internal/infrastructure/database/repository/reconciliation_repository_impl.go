package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"conciliacao-bancaria/internal/domain/model"
	domainRepo "conciliacao-bancaria/internal/domain/repository"
)

// Garantir que ReconciliationRepositoryImpl implementa a interface ReconciliationRepository
var _ domainRepo.ReconciliationRepository = (*ReconciliationRepositoryImpl)(nil)

// ReconciliationRepositoryImpl implementa a interface de repositório para conciliações
type ReconciliationRepositoryImpl struct {
	db *sql.DB
}

// NewReconciliationRepository cria uma nova instância do repositório de conciliação
func NewReconciliationRepository(db *sql.DB) domainRepo.ReconciliationRepository {
	return &ReconciliationRepositoryImpl{
		db: db,
	}
}

// Create persiste uma nova conciliação no banco de dados
func (r *ReconciliationRepositoryImpl) Create(ctx context.Context, reconciliation *model.Reconciliation) error {
	query := `
		INSERT INTO reconciliation (
			id, billet_id, transaction_id, reconciliation_date, 
			conciliation_status, conciliation_strategy, amount_diff, reference_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Usar context com timeout para evitar operações longas em caso de problemas com o banco
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(
		ctxWithTimeout,
		query,
		reconciliation.ID,
		reconciliation.BilletID,
		reconciliation.TransactionID,
		reconciliation.ReconciliationDate,
		string(reconciliation.ConciliationStatus),
		string(reconciliation.ConciliationStrategy),
		reconciliation.AmountDiff,
		reconciliation.ReferenceID,
	)

	if err != nil {
		return fmt.Errorf("erro ao criar conciliação: %w", err)
	}

	return nil
}

// CreateMany persiste múltiplas conciliações no banco de dados
func (r *ReconciliationRepositoryImpl) CreateMany(ctx context.Context, reconciliations []*model.Reconciliation) error {
	if len(reconciliations) == 0 {
		return nil
	}

	// Iniciar uma transação para garantir a atomicidade da operação
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	// Defer para garantir que a transação será revertida em caso de erro
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		INSERT INTO reconciliation (
			id, billet_id, transaction_id, reconciliation_date, 
			conciliation_status, conciliation_strategy, amount_diff, reference_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("erro ao preparar statement: %w", err)
	}
	defer stmt.Close()

	for _, reconciliation := range reconciliations {
		_, err = stmt.ExecContext(
			ctx,
			reconciliation.ID,
			reconciliation.BilletID,
			reconciliation.TransactionID,
			reconciliation.ReconciliationDate,
			string(reconciliation.ConciliationStatus),
			string(reconciliation.ConciliationStrategy),
			reconciliation.AmountDiff,
			reconciliation.ReferenceID,
		)

		if err != nil {
			return fmt.Errorf("erro ao inserir conciliação %s: %w", reconciliation.ID, err)
		}
	}

	// Commit da transação
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	return nil
}

// GetByID recupera uma conciliação pelo seu ID
func (r *ReconciliationRepositoryImpl) GetByID(ctx context.Context, id string) (*model.Reconciliation, error) {
	query := `
		SELECT 
			id, billet_id, transaction_id, reconciliation_date, 
			conciliation_status, conciliation_strategy, amount_diff, reference_id
		FROM reconciliation
		WHERE id = ?
	`

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctxWithTimeout, query, id)

	reconciliation := &model.Reconciliation{}
	var conciliationStatus, conciliationStrategy string
	var referenceID sql.NullString

	err := row.Scan(
		&reconciliation.ID,
		&reconciliation.BilletID,
		&reconciliation.TransactionID,
		&reconciliation.ReconciliationDate,
		&conciliationStatus,
		&conciliationStrategy,
		&reconciliation.AmountDiff,
		&referenceID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("conciliação não encontrada: %w", err)
		}
		return nil, fmt.Errorf("erro ao buscar conciliação: %w", err)
	}

	// Converter os valores de string para os tipos de enum
	reconciliation.ConciliationStatus = model.ConciliationStatus(conciliationStatus)
	reconciliation.ConciliationStrategy = model.ConciliationStrategy(conciliationStrategy)

	// Tratar campo opcional
	if referenceID.Valid {
		reconciliation.ReferenceID = &referenceID.String
	}

	return reconciliation, nil
}

// GetAll recupera todas as conciliações
func (r *ReconciliationRepositoryImpl) GetAll(ctx context.Context) ([]*model.Reconciliation, error) {
	query := `
		SELECT 
			id, billet_id, transaction_id, reconciliation_date, 
			conciliation_status, conciliation_strategy, amount_diff, reference_id
		FROM reconciliation
		ORDER BY reconciliation_date DESC
	`

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctxWithTimeout, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar conciliações: %w", err)
	}
	defer rows.Close()

	reconciliations := []*model.Reconciliation{}

	for rows.Next() {
		reconciliation := &model.Reconciliation{}
		var conciliationStatus, conciliationStrategy string
		var referenceID sql.NullString

		err := rows.Scan(
			&reconciliation.ID,
			&reconciliation.BilletID,
			&reconciliation.TransactionID,
			&reconciliation.ReconciliationDate,
			&conciliationStatus,
			&conciliationStrategy,
			&reconciliation.AmountDiff,
			&referenceID,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao ler conciliação: %w", err)
		}

		// Converter os valores de string para os tipos de enum
		reconciliation.ConciliationStatus = model.ConciliationStatus(conciliationStatus)
		reconciliation.ConciliationStrategy = model.ConciliationStrategy(conciliationStrategy)

		// Tratar campo opcional
		if referenceID.Valid {
			reconciliation.ReferenceID = &referenceID.String
		}

		reconciliations = append(reconciliations, reconciliation)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao processar resultados: %w", err)
	}

	return reconciliations, nil
}

// GetByBilletID recupera conciliações por ID do boleto
func (r *ReconciliationRepositoryImpl) GetByBilletID(ctx context.Context, billetID string) ([]*model.Reconciliation, error) {
	query := `
		SELECT 
			id, billet_id, transaction_id, reconciliation_date, 
			conciliation_status, conciliation_strategy, amount_diff, reference_id
		FROM reconciliation
		WHERE billet_id = ?
		ORDER BY reconciliation_date DESC
	`

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctxWithTimeout, query, billetID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar conciliações por boleto: %w", err)
	}
	defer rows.Close()

	reconciliations := []*model.Reconciliation{}

	for rows.Next() {
		reconciliation := &model.Reconciliation{}
		var conciliationStatus, conciliationStrategy string
		var referenceID sql.NullString

		err := rows.Scan(
			&reconciliation.ID,
			&reconciliation.BilletID,
			&reconciliation.TransactionID,
			&reconciliation.ReconciliationDate,
			&conciliationStatus,
			&conciliationStrategy,
			&reconciliation.AmountDiff,
			&referenceID,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao ler conciliação: %w", err)
		}

		// Converter os valores de string para os tipos de enum
		reconciliation.ConciliationStatus = model.ConciliationStatus(conciliationStatus)
		reconciliation.ConciliationStrategy = model.ConciliationStrategy(conciliationStrategy)

		// Tratar campo opcional
		if referenceID.Valid {
			reconciliation.ReferenceID = &referenceID.String
		}

		reconciliations = append(reconciliations, reconciliation)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao processar resultados: %w", err)
	}

	return reconciliations, nil
}

// GetByTransactionID recupera conciliações por ID da transação
func (r *ReconciliationRepositoryImpl) GetByTransactionID(ctx context.Context, transactionID string) ([]*model.Reconciliation, error) {
	query := `
		SELECT 
			id, billet_id, transaction_id, reconciliation_date, 
			conciliation_status, conciliation_strategy, amount_diff, reference_id
		FROM reconciliation
		WHERE transaction_id = ?
		ORDER BY reconciliation_date DESC
	`

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctxWithTimeout, query, transactionID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar conciliações por transação: %w", err)
	}
	defer rows.Close()

	reconciliations := []*model.Reconciliation{}

	for rows.Next() {
		reconciliation := &model.Reconciliation{}
		var conciliationStatus, conciliationStrategy string
		var referenceID sql.NullString

		err := rows.Scan(
			&reconciliation.ID,
			&reconciliation.BilletID,
			&reconciliation.TransactionID,
			&reconciliation.ReconciliationDate,
			&conciliationStatus,
			&conciliationStrategy,
			&reconciliation.AmountDiff,
			&referenceID,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao ler conciliação: %w", err)
		}

		// Converter os valores de string para os tipos de enum
		reconciliation.ConciliationStatus = model.ConciliationStatus(conciliationStatus)
		reconciliation.ConciliationStrategy = model.ConciliationStrategy(conciliationStrategy)

		// Tratar campo opcional
		if referenceID.Valid {
			reconciliation.ReferenceID = &referenceID.String
		}

		reconciliations = append(reconciliations, reconciliation)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao processar resultados: %w", err)
	}

	return reconciliations, nil
}

// Update atualiza uma conciliação existente
func (r *ReconciliationRepositoryImpl) Update(ctx context.Context, reconciliation *model.Reconciliation) error {
	query := `
		UPDATE reconciliation 
		SET 
			billet_id = ?, 
			transaction_id = ?, 
			reconciliation_date = ?, 
			conciliation_status = ?, 
			conciliation_strategy = ?, 
			amount_diff = ?, 
			reference_id = ?
		WHERE id = ?
	`

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(
		ctxWithTimeout,
		query,
		reconciliation.BilletID,
		reconciliation.TransactionID,
		reconciliation.ReconciliationDate,
		string(reconciliation.ConciliationStatus),
		string(reconciliation.ConciliationStrategy),
		reconciliation.AmountDiff,
		reconciliation.ReferenceID,
		reconciliation.ID,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar conciliação: %w", err)
	}

	return nil
}

// Delete remove uma conciliação pelo ID
func (r *ReconciliationRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM reconciliation WHERE id = ?"

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctxWithTimeout, query, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir conciliação: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nenhuma conciliação encontrada com o ID: %s", id)
	}

	return nil
}

// GetReconciliationHistory recupera o histórico de conciliações para auditoria
func (r *ReconciliationRepositoryImpl) GetReconciliationHistory(ctx context.Context, billetID string) ([]*model.Reconciliation, error) {
	query := `
		SELECT 
			id, billet_id, transaction_id, reconciliation_date, 
			conciliation_status, conciliation_strategy, amount_diff, reference_id
		FROM reconciliation
		WHERE billet_id = ?
		ORDER BY reconciliation_date ASC
	`

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctxWithTimeout, query, billetID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar histórico de conciliações: %w", err)
	}
	defer rows.Close()

	reconciliations := []*model.Reconciliation{}

	for rows.Next() {
		reconciliation := &model.Reconciliation{}
		var conciliationStatus, conciliationStrategy string
		var referenceID sql.NullString

		err := rows.Scan(
			&reconciliation.ID,
			&reconciliation.BilletID,
			&reconciliation.TransactionID,
			&reconciliation.ReconciliationDate,
			&conciliationStatus,
			&conciliationStrategy,
			&reconciliation.AmountDiff,
			&referenceID,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao ler histórico de conciliação: %w", err)
		}

		// Converter os valores de string para os tipos de enum
		reconciliation.ConciliationStatus = model.ConciliationStatus(conciliationStatus)
		reconciliation.ConciliationStrategy = model.ConciliationStrategy(conciliationStrategy)

		// Tratar campo opcional
		if referenceID.Valid {
			reconciliation.ReferenceID = &referenceID.String
		}

		reconciliations = append(reconciliations, reconciliation)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao processar resultados do histórico: %w", err)
	}

	return reconciliations, nil
}
