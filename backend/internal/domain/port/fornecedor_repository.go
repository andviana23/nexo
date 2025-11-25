package port

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/google/uuid"
)

// FornecedorRepository define as operações de persistência para fornecedores
type FornecedorRepository interface {
	// CRUD básico
	Create(ctx context.Context, fornecedor *entity.Fornecedor) error
	FindByID(ctx context.Context, tenantID, fornecedorID uuid.UUID) (*entity.Fornecedor, error)
	FindByCNPJ(ctx context.Context, cnpj string) (*entity.Fornecedor, error)
	Update(ctx context.Context, fornecedor *entity.Fornecedor) error
	Delete(ctx context.Context, tenantID, fornecedorID uuid.UUID) error

	// Listagens
	ListAll(ctx context.Context, tenantID uuid.UUID) ([]*entity.Fornecedor, error)
	ListAtivos(ctx context.Context, tenantID uuid.UUID) ([]*entity.Fornecedor, error)
}
