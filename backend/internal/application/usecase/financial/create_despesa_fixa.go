package financial

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateDespesaFixaInput define os dados de entrada para criar despesa fixa
type CreateDespesaFixaInput struct {
	TenantID      string
	Descricao     string
	CategoriaID   string
	UnidadeID     string
	Fornecedor    string
	Valor         valueobject.Money
	DiaVencimento int
	Observacoes   string
}

// CreateDespesaFixaUseCase implementa a criação de despesa fixa
type CreateDespesaFixaUseCase struct {
	repo   port.DespesaFixaRepository
	logger *zap.Logger
}

// NewCreateDespesaFixaUseCase cria nova instância do use case
func NewCreateDespesaFixaUseCase(
	repo port.DespesaFixaRepository,
	logger *zap.Logger,
) *CreateDespesaFixaUseCase {
	return &CreateDespesaFixaUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute cria uma nova despesa fixa
func (uc *CreateDespesaFixaUseCase) Execute(ctx context.Context, input CreateDespesaFixaInput) (*entity.DespesaFixa, error) {
	// Validações de entrada
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	if input.Descricao == "" {
		return nil, domain.ErrDescricaoRequired
	}

	if input.DiaVencimento < 1 || input.DiaVencimento > 31 {
		return nil, domain.ErrDiaVencimentoInvalido
	}

	if input.Valor.IsNegative() || input.Valor.IsZero() {
		return nil, domain.ErrValorNegativo
	}

	// Verificar duplicidade de descrição
	exists, err := uc.repo.ExistsByDescricao(ctx, input.TenantID, input.Descricao, nil)
	if err != nil {
		uc.logger.Error("Erro ao verificar duplicidade de despesa fixa",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
		)
		return nil, err
	}
	if exists {
		return nil, domain.ErrDespesaFixaJaExiste
	}

	// Converter TenantID para UUID
	tenantUUID, err := uuid.Parse(input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant_id inválido: %w", err)
	}

	// Criar entidade de domínio
	despesa, err := entity.NewDespesaFixa(
		tenantUUID,
		input.Descricao,
		input.Valor,
		input.DiaVencimento,
	)
	if err != nil {
		uc.logger.Error("Erro ao criar entidade DespesaFixa",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
		)
		return nil, err
	}

	// Atribuir campos opcionais
	if input.CategoriaID != "" {
		despesa.CategoriaID = input.CategoriaID
	}
	if input.UnidadeID != "" {
		despesa.UnidadeID = input.UnidadeID
	}
	if input.Fornecedor != "" {
		despesa.Fornecedor = input.Fornecedor
	}
	if input.Observacoes != "" {
		despesa.Observacoes = input.Observacoes
	}

	// Persistir
	if err := uc.repo.Create(ctx, despesa); err != nil {
		uc.logger.Error("Erro ao criar despesa fixa",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
		)
		return nil, err
	}

	uc.logger.Info("Despesa fixa criada com sucesso",
		zap.String("id", despesa.ID),
		zap.String("tenant_id", despesa.TenantID.String()),
		zap.String("descricao", despesa.Descricao),
		zap.Int("dia_vencimento", despesa.DiaVencimento),
	)

	return despesa, nil
}
