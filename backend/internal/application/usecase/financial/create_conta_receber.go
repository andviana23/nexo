package financial

import (
	"context"

	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateContaReceberInput define os dados de entrada para criar conta a receber
type CreateContaReceberInput struct {
	TenantID        string
	Descricao       string
	Origem          string
	Subtipo         valueobject.SubtipoReceita
	AssinaturaID    *string
	Valor           valueobject.Money
	DataVencimento  time.Time
	MetodoPagamento string
	Observacoes     string
}

// CreateContaReceberUseCase implementa a criação de conta a receber
type CreateContaReceberUseCase struct {
	repo   port.ContaReceberRepository
	logger *zap.Logger
}

// NewCreateContaReceberUseCase cria nova instância do use case
func NewCreateContaReceberUseCase(
	repo port.ContaReceberRepository,
	logger *zap.Logger,
) *CreateContaReceberUseCase {
	return &CreateContaReceberUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute cria uma nova conta a receber
func (uc *CreateContaReceberUseCase) Execute(ctx context.Context, input CreateContaReceberInput) (*entity.ContaReceber, error) {
	// Validações de entrada
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	if input.Descricao == "" {
		return nil, fmt.Errorf("descrição é obrigatória")
	}

	if input.Origem == "" {
		return nil, fmt.Errorf("origem é obrigatória")
	}

	if input.DataVencimento.Before(time.Now().Add(-24 * time.Hour)) {
		return nil, fmt.Errorf("data de vencimento não pode ser no passado")
	}

	// Converter TenantID para UUID
	tenantUUID, err := uuid.Parse(input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant_id inválido: %w", err)
	}

	// Criar entidade de domínio
	conta, err := entity.NewContaReceber(
		tenantUUID,
		input.Origem,
		input.AssinaturaID,
		input.Descricao,
		input.Valor,
		input.DataVencimento,
	)
	if err != nil {
		uc.logger.Error("Erro ao criar entidade ContaReceber",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("descricao", input.Descricao),
		)
		return nil, fmt.Errorf("erro ao criar conta a receber: %w", err)
	}

	// Atribuir campos opcionais
	conta.Observacoes = input.Observacoes

	// Persistir no repositório
	if err := uc.repo.Create(ctx, conta); err != nil {
		uc.logger.Error("Erro ao persistir conta a receber",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("conta_id", conta.ID),
		)
		return nil, fmt.Errorf("erro ao salvar conta a receber: %w", err)
	}

	uc.logger.Info("Conta a receber criada com sucesso",
		zap.String("tenant_id", input.TenantID),
		zap.String("conta_id", conta.ID),
		zap.String("descricao", conta.DescricaoOrigem),
	)

	return conta, nil
}
