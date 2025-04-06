package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"conciliacao-bancaria/internal/domain/model"
	"conciliacao-bancaria/internal/domain/repository"
)

// SQLPaymentRepository implementa a interface PaymentRepository usando SQL
type SQLPaymentRepository struct {
	db *sql.DB
}

// NewPaymentRepository cria uma nova instância de SQLPaymentRepository
func NewPaymentRepository(db *sql.DB) repository.PaymentRepository {
	return &SQLPaymentRepository{db: db}
}

// Create persiste um novo pagamento no banco de dados
func (r *SQLPaymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	query := `
		INSERT INTO payments (
			id, bank_account, amount, payment_date, reference_id, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`

	now := time.Now()
	_, err := r.db.ExecContext(
		ctx,
		query,
		payment.ID,
		payment.BankAccount,
		payment.Amount,
		payment.PaymentDate,
		payment.ReferenceID,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("falha ao criar pagamento: %w", err)
	}

	return nil
}

// CreateMany persiste múltiplos pagamentos no banco de dados
func (r *SQLPaymentRepository) CreateMany(ctx context.Context, payments []*model.Payment) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("falha ao iniciar transação: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		INSERT INTO payments (
			id, bank_account, amount, payment_date, reference_id, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("falha ao preparar declaração: %w", err)
	}
	defer stmt.Close()

	now := time.Now()
	for _, payment := range payments {
		_, err = stmt.ExecContext(
			ctx,
			payment.ID,
			payment.BankAccount,
			payment.Amount,
			payment.PaymentDate,
			payment.ReferenceID,
			now,
			now,
		)

		if err != nil {
			return fmt.Errorf("falha ao inserir pagamento %s: %w", payment.ID, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("falha ao confirmar transação: %w", err)
	}

	return nil
}

// GetByID recupera um pagamento pelo seu ID
func (r *SQLPaymentRepository) GetByID(ctx context.Context, id string) (*model.Payment, error) {
	query := `
		SELECT 
			id, bank_account, amount, payment_date, reference_id, created_at, updated_at
		FROM 
			payments 
		WHERE 
			id = $1
	`

	var payment model.Payment
	var referenceID sql.NullString
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID,
		&payment.BankAccount,
		&payment.Amount,
		&payment.PaymentDate,
		&referenceID,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Não encontrado
		}
		return nil, fmt.Errorf("falha ao recuperar pagamento: %w", err)
	}

	if referenceID.Valid {
		refID := referenceID.String
		payment.ReferenceID = &refID
	}

	return &payment, nil
}

// GetAll recupera todos os pagamentos
func (r *SQLPaymentRepository) GetAll(ctx context.Context) ([]*model.Payment, error) {
	query := `
		SELECT 
			id, bank_account, amount, payment_date, reference_id, created_at, updated_at
		FROM 
			payments
		ORDER BY
			payment_date
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar pagamentos: %w", err)
	}
	defer rows.Close()

	var payments []*model.Payment
	for rows.Next() {
		var payment model.Payment
		var referenceID sql.NullString
		var createdAt, updatedAt time.Time

		if err := rows.Scan(
			&payment.ID,
			&payment.BankAccount,
			&payment.Amount,
			&payment.PaymentDate,
			&referenceID,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("falha ao ler pagamento: %w", err)
		}

		if referenceID.Valid {
			refID := referenceID.String
			payment.ReferenceID = &refID
		}

		payments = append(payments, &payment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados: %w", err)
	}

	return payments, nil
}

// GetByBankAccount recupera pagamentos por conta bancária
func (r *SQLPaymentRepository) GetByBankAccount(ctx context.Context, bankAccount string) ([]*model.Payment, error) {
	query := `
		SELECT 
			id, bank_account, amount, payment_date, reference_id, created_at, updated_at
		FROM 
			payments
		WHERE
			bank_account = $1
		ORDER BY
			payment_date
	`

	rows, err := r.db.QueryContext(ctx, query, bankAccount)
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar pagamentos por conta bancária: %w", err)
	}
	defer rows.Close()

	var payments []*model.Payment
	for rows.Next() {
		var payment model.Payment
		var referenceID sql.NullString
		var createdAt, updatedAt time.Time

		if err := rows.Scan(
			&payment.ID,
			&payment.BankAccount,
			&payment.Amount,
			&payment.PaymentDate,
			&referenceID,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("falha ao ler pagamento: %w", err)
		}

		if referenceID.Valid {
			refID := referenceID.String
			payment.ReferenceID = &refID
		}

		payments = append(payments, &payment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados: %w", err)
	}

	return payments, nil
}

// GetByReferenceID recupera pagamentos por ID de referência
func (r *SQLPaymentRepository) GetByReferenceID(ctx context.Context, referenceID string) ([]*model.Payment, error) {
	query := `
		SELECT 
			id, bank_account, amount, payment_date, reference_id, created_at, updated_at
		FROM 
			payments
		WHERE
			reference_id = $1
		ORDER BY
			payment_date
	`

	rows, err := r.db.QueryContext(ctx, query, referenceID)
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar pagamentos por ID de referência: %w", err)
	}
	defer rows.Close()

	var payments []*model.Payment
	for rows.Next() {
		var payment model.Payment
		var refID sql.NullString
		var createdAt, updatedAt time.Time

		if err := rows.Scan(
			&payment.ID,
			&payment.BankAccount,
			&payment.Amount,
			&payment.PaymentDate,
			&refID,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("falha ao ler pagamento: %w", err)
		}

		if refID.Valid {
			refIDStr := refID.String
			payment.ReferenceID = &refIDStr
		}

		payments = append(payments, &payment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados: %w", err)
	}

	return payments, nil
}

// Update atualiza um pagamento existente
func (r *SQLPaymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	query := `
		UPDATE payments
		SET
			bank_account = $1,
			amount = $2,
			payment_date = $3,
			reference_id = $4,
			updated_at = $5
		WHERE
			id = $6
	`

	now := time.Now()
	result, err := r.db.ExecContext(
		ctx,
		query,
		payment.BankAccount,
		payment.Amount,
		payment.PaymentDate,
		payment.ReferenceID,
		now,
		payment.ID,
	)

	if err != nil {
		return fmt.Errorf("falha ao atualizar pagamento: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("falha ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nenhum pagamento atualizado com o ID: %s", payment.ID)
	}

	return nil
}

// Delete remove um pagamento pelo ID
func (r *SQLPaymentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM payments WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("falha ao excluir pagamento: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("falha ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nenhum pagamento excluído com o ID: %s", id)
	}

	return nil
}

// FindByBankAccountAndAmount encontra pagamentos por conta bancária e valor aproximado
func (r *SQLPaymentRepository) FindByBankAccountAndAmount(ctx context.Context, bankAccount string, amount float64, tolerance float64) ([]*model.Payment, error) {
	// Calculando o intervalo de tolerância
	minAmount := amount - (amount * tolerance / 100)
	maxAmount := amount + (amount * tolerance / 100)

	query := `
		SELECT 
			id, bank_account, amount, payment_date, reference_id, created_at, updated_at
		FROM 
			payments
		WHERE
			bank_account = $1
			AND amount BETWEEN $2 AND $3
		ORDER BY
			payment_date
	`

	rows, err := r.db.QueryContext(ctx, query, bankAccount, minAmount, maxAmount)
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar pagamentos por conta e valor: %w", err)
	}
	defer rows.Close()

	var payments []*model.Payment
	for rows.Next() {
		var payment model.Payment
		var referenceID sql.NullString
		var createdAt, updatedAt time.Time

		if err := rows.Scan(
			&payment.ID,
			&payment.BankAccount,
			&payment.Amount,
			&payment.PaymentDate,
			&referenceID,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("falha ao ler pagamento: %w", err)
		}

		if referenceID.Valid {
			refID := referenceID.String
			payment.ReferenceID = &refID
		}

		payments = append(payments, &payment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados: %w", err)
	}

	return payments, nil
}
