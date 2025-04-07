package errors

import (
	"errors"
	"fmt"
)

// Erros básicos para reutilização
var (
	ErrNotFound         = errors.New("recurso não encontrado")
	ErrInvalidInput     = errors.New("dados de entrada inválidos")
	ErrAlreadyExists    = errors.New("recurso já existe")
	ErrDatabaseError    = errors.New("erro na operação com banco de dados")
	ErrUnauthorized     = errors.New("não autorizado")
	ErrInternalError    = errors.New("erro interno do servidor")
	ErrInvalidOperation = errors.New("operação inválida")
)

// NotFoundError representa erro de recurso não encontrado
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s com ID %s não encontrado", e.Resource, e.ID)
}

// ValidationError representa erro de validação de dados
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("erro de validação no campo %s: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("erro de validação: %s", e.Message)
}

// ConflictError representa erro de conflito (recurso já existe, etc)
type ConflictError struct {
	Resource string
	ID       string
	Reason   string
}

func (e *ConflictError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("conflito com %s (ID: %s): %s", e.Resource, e.ID, e.Reason)
	}
	return fmt.Sprintf("conflito com %s (ID: %s)", e.Resource, e.ID)
}

// DatabaseError representa erro de operação com banco de dados
type DatabaseError struct {
	Operation string
	Err       error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("erro na operação '%s' do banco de dados: %v", e.Operation, e.Err)
}

// NewNotFoundError cria um novo erro de recurso não encontrado
func NewNotFoundError(resource, id string) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

// NewValidationError cria um novo erro de validação
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// NewConflictError cria um novo erro de conflito
func NewConflictError(resource, id, reason string) *ConflictError {
	return &ConflictError{
		Resource: resource,
		ID:       id,
		Reason:   reason,
	}
}

// NewDatabaseError cria um novo erro de banco de dados
func NewDatabaseError(operation string, err error) *DatabaseError {
	return &DatabaseError{
		Operation: operation,
		Err:       err,
	}
}

// IsNotFoundError verifica se um erro é do tipo NotFoundError
func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

// IsValidationError verifica se um erro é do tipo ValidationError
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// IsConflictError verifica se um erro é do tipo ConflictError
func IsConflictError(err error) bool {
	_, ok := err.(*ConflictError)
	return ok
}

// IsDatabaseError verifica se um erro é do tipo DatabaseError
func IsDatabaseError(err error) bool {
	_, ok := err.(*DatabaseError)
	return ok
}

// Wrap adiciona contexto a um erro existente
func Wrap(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}
