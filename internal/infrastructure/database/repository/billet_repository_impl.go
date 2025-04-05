package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"conciliacao-bancaria/internal/domain/model"
	"conciliacao-bancaria/internal/domain/repository"
)

// billetRepositoryImpl implementa a interface BilletRepository
type billetRepositoryImpl struct {
	db *sql.DB
}

// NewBilletRepository cria uma nova instância de BilletRepository
func NewBilletRepository(db *sql.DB) repository.BilletRepository {
	return &billetRepositoryImpl{db: db}
}

// Create persiste um novo boleto no banco de dados
func (r *billetRepositoryImpl) Create(ctx context.Context, billet *model.Billet) error {
	query := `
		INSERT INTO bank_reconciliation.billets 
		(id, bank_account, amount, issuance_date, reference_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	now := time.Now()
	var referenceID *string
	if billet.ReferenceID != nil {
		referenceID = billet.ReferenceID
	}

	_, err := r.db.ExecContext(ctx, query,
		billet.ID,
		billet.BankAccount,
		billet.Amount,
		billet.IssuanceDate,
		referenceID,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("erro ao criar boleto: %w", err)
	}

	return nil
}

// CreateMany persiste múltiplos boletos no banco de dados
func (r *billetRepositoryImpl) CreateMany(ctx context.Context, billets []*model.Billet) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	query := `
		INSERT INTO bank_reconciliation.billets 
		(id, bank_account, amount, issuance_date, reference_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao preparar statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now()

	for _, billet := range billets {
		var referenceID *string
		if billet.ReferenceID != nil {
			referenceID = billet.ReferenceID
		}

		_, err := stmt.ExecContext(ctx,
			billet.ID,
			billet.BankAccount,
			billet.Amount,
			billet.IssuanceDate,
			referenceID,
			now,
			now,
		)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("erro ao criar boleto no batch: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao fazer commit da transação: %w", err)
	}

	return nil
}

// GetByID recupera um boleto pelo seu ID
func (r *billetRepositoryImpl) GetByID(ctx context.Context, id string) (*model.Billet, error) {
	query := `
		SELECT id, bank_account, amount, issuance_date, reference_id, created_at, updated_at
		FROM bank_reconciliation.billets
		WHERE id = $1
	`

	var billet model.Billet
	var referenceID sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&billet.ID,
		&billet.BankAccount,
		&billet.Amount,
		&billet.IssuanceDate,
		&referenceID,
		&billet.CreatedAt,
		&billet.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("boleto não encontrado: %s", id)
		}
		return nil, fmt.Errorf("erro ao buscar boleto: %w", err)
	}

	if referenceID.Valid {
		refID := referenceID.String
		billet.ReferenceID = &refID
	}

	return &billet, nil
}

// GetAll recupera todos os boletos
func (r *billetRepositoryImpl) GetAll(ctx context.Context) ([]*model.Billet, error) {
	query := `
		SELECT id, bank_account, amount, issuance_date, reference_id, created_at, updated_at
		FROM bank_reconciliation.billets
		ORDER BY issuance_date
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar boletos: %w", err)
	}
	defer rows.Close()

	var billets []*model.Billet

	for rows.Next() {
		var billet model.Billet
		var referenceID sql.NullString

		err := rows.Scan(
			&billet.ID,
			&billet.BankAccount,
			&billet.Amount,
			&billet.IssuanceDate,
			&referenceID,
			&billet.CreatedAt,
			&billet.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao ler boleto: %w", err)
		}

		if referenceID.Valid {
			refID := referenceID.String
			billet.ReferenceID = &refID
		}

		billets = append(billets, &billet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre boletos: %w", err)
	}

	return billets, nil
}

// GetByBankAccount recupera boletos por conta bancária
func (r *billetRepositoryImpl) GetByBankAccount(ctx context.Context, bankAccount string) ([]*model.Billet, error) {
	query := `
		SELECT id, bank_account, amount, issuance_date, reference_id, created_at, updated_at
		FROM bank_reconciliation.billets
		WHERE bank_account = $1
		ORDER BY issuance_date
	`

	rows, err := r.db.QueryContext(ctx, query, bankAccount)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar boletos por conta bancária: %w", err)
	}
	defer rows.Close()

	var billets []*model.Billet

	for rows.Next() {
		var billet model.Billet
		var referenceID sql.NullString

		err := rows.Scan(
			&billet.ID,
			&billet.BankAccount,
			&billet.Amount,
			&billet.IssuanceDate,
			&referenceID,
			&billet.CreatedAt,
			&billet.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao ler boleto: %w", err)
		}

		if referenceID.Valid {
			refID := referenceID.String
			billet.ReferenceID = &refID
		}

		billets = append(billets, &billet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre boletos: %w", err)
	}

	return billets, nil
}

// GetByReferenceID recupera boletos por ID de referência
func (r *billetRepositoryImpl) GetByReferenceID(ctx context.Context, referenceID string) ([]*model.Billet, error) {
	query := `
		SELECT id, bank_account, amount, issuance_date, reference_id, created_at, updated_at
		FROM bank_reconciliation.billets
		WHERE reference_id = $1
		ORDER BY issuance_date
	`

	rows, err := r.db.QueryContext(ctx, query, referenceID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar boletos por ID de referência: %w", err)
	}
	defer rows.Close()

	var billets []*model.Billet

	for rows.Next() {
		var billet model.Billet
		var refID sql.NullString

		err := rows.Scan(
			&billet.ID,
			&billet.BankAccount,
			&billet.Amount,
			&billet.IssuanceDate,
			&refID,
			&billet.CreatedAt,
			&billet.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao ler boleto: %w", err)
		}

		if refID.Valid {
			id := refID.String
			billet.ReferenceID = &id
		}

		billets = append(billets, &billet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre boletos: %w", err)
	}

	return billets, nil
}

// Update atualiza um boleto existente
func (r *billetRepositoryImpl) Update(ctx context.Context, billet *model.Billet) error {
	query := `
		UPDATE bank_reconciliation.billets
		SET bank_account = $1, amount = $2, issuance_date = $3, reference_id = $4
		WHERE id = $5
	`

	var referenceID *string
	if billet.ReferenceID != nil {
		referenceID = billet.ReferenceID
	}

	result, err := r.db.ExecContext(ctx, query,
		billet.BankAccount,
		billet.Amount,
		billet.IssuanceDate,
		referenceID,
		billet.ID,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar boleto: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("boleto não encontrado: %s", billet.ID)
	}

	return nil
}

// Delete remove um boleto pelo ID
func (r *billetRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM bank_reconciliation.billets
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir boleto: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("boleto não encontrado: %s", id)
	}

	return nil
}

// FindNonReconciled encontra boletos que ainda não foram conciliados
func (r *billetRepositoryImpl) FindNonReconciled(ctx context.Context) ([]*model.Billet, error) {
	query := `
		SELECT b.id, b.bank_account, b.amount, b.issuance_date, b.reference_id, b.created_at, b.updated_at
		FROM bank_reconciliation.billets b
		LEFT JOIN bank_reconciliation.reconciliations r ON b.id = r.billet_id
		WHERE r.id IS NULL
		ORDER BY b.issuance_date
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar boletos não conciliados: %w", err)
	}
	defer rows.Close()

	var billets []*model.Billet

	for rows.Next() {
		var billet model.Billet
		var referenceID sql.NullString

		err := rows.Scan(
			&billet.ID,
			&billet.BankAccount,
			&billet.Amount,
			&billet.IssuanceDate,
			&referenceID,
			&billet.CreatedAt,
			&billet.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao ler boleto não conciliado: %w", err)
		}

		if referenceID.Valid {
			refID := referenceID.String
			billet.ReferenceID = &refID
		}

		billets = append(billets, &billet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre boletos não conciliados: %w", err)
	}

	return billets, nil
}
