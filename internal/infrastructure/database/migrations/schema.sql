-- Criação de Schema
CREATE SCHEMA IF NOT EXISTS bank_reconciliation;

-- Definição de tabelas

-- Tabela de Boletos
CREATE TABLE IF NOT EXISTS bank_reconciliation.billets (
    id VARCHAR(50) PRIMARY KEY,
    bank_account VARCHAR(50) NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    issuance_date TIMESTAMP NOT NULL,
    reference_id VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Pagamentos
CREATE TABLE IF NOT EXISTS bank_reconciliation.payments (
    id VARCHAR(50) PRIMARY KEY,
    bank_account VARCHAR(50) NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    payment_date TIMESTAMP NOT NULL,
    reference_id VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Conciliações
CREATE TABLE IF NOT EXISTS bank_reconciliation.reconciliations (
    id VARCHAR(50) PRIMARY KEY,
    billet_id VARCHAR(50) NOT NULL,
    transaction_id VARCHAR(50),
    bank_account VARCHAR(50) NOT NULL,
    conciliation_status VARCHAR(30) NOT NULL,
    conciliation_strategy VARCHAR(30) NOT NULL,
    amount_diff DECIMAL(15, 2) NOT NULL,
    reference_id VARCHAR(50),
    reconciliation_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_billet_id FOREIGN KEY (billet_id) REFERENCES bank_reconciliation.billets(id),
    CONSTRAINT fk_transaction_id FOREIGN KEY (transaction_id) REFERENCES bank_reconciliation.payments(id)
);

-- Índices para melhorar performance de consultas

-- Índices para tabela de boletos
CREATE INDEX IF NOT EXISTS idx_billets_bank_account ON bank_reconciliation.billets(bank_account);
CREATE INDEX IF NOT EXISTS idx_billets_reference_id ON bank_reconciliation.billets(reference_id);
CREATE INDEX IF NOT EXISTS idx_billets_issuance_date ON bank_reconciliation.billets(issuance_date);
CREATE INDEX IF NOT EXISTS idx_billets_amount ON bank_reconciliation.billets(amount);

-- Índices para tabela de pagamentos
CREATE INDEX IF NOT EXISTS idx_payments_bank_account ON bank_reconciliation.payments(bank_account);
CREATE INDEX IF NOT EXISTS idx_payments_reference_id ON bank_reconciliation.payments(reference_id);
CREATE INDEX IF NOT EXISTS idx_payments_payment_date ON bank_reconciliation.payments(payment_date);
CREATE INDEX IF NOT EXISTS idx_payments_amount ON bank_reconciliation.payments(amount);

-- Índices para tabela de conciliações
CREATE INDEX IF NOT EXISTS idx_reconciliations_billet_id ON bank_reconciliation.reconciliations(billet_id);
CREATE INDEX IF NOT EXISTS idx_reconciliations_transaction_id ON bank_reconciliation.reconciliations(transaction_id);
CREATE INDEX IF NOT EXISTS idx_reconciliations_status ON bank_reconciliation.reconciliations(conciliation_status);
CREATE INDEX IF NOT EXISTS idx_reconciliations_date ON bank_reconciliation.reconciliations(reconciliation_date);

-- Função para atualizar o updated_at automaticamente
CREATE OR REPLACE FUNCTION bank_reconciliation.update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers para atualizar automaticamente o updated_at
CREATE TRIGGER update_billets_modtime
BEFORE UPDATE ON bank_reconciliation.billets
FOR EACH ROW
EXECUTE FUNCTION bank_reconciliation.update_modified_column();

CREATE TRIGGER update_payments_modtime
BEFORE UPDATE ON bank_reconciliation.payments
FOR EACH ROW
EXECUTE FUNCTION bank_reconciliation.update_modified_column();

CREATE TRIGGER update_reconciliations_modtime
BEFORE UPDATE ON bank_reconciliation.reconciliations
FOR EACH ROW
EXECUTE FUNCTION bank_reconciliation.update_modified_column();