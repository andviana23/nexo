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

// CreateContaPagarInput define os dados de entrada para criar conta a pagar
type CreateContaPagarInput struct {
	TenantID       string
	Descricao      string
	CategoriaID    string
	Fornecedor     string
	Valor          valueobject.Money
	Tipo           valueobject.TipoCusto
	DataVencimento time.Time
	Recorrente     bool
	Periodicidade  string
	PixCode        string
	Observacoes    string
}

// CreateContaPagarUseCase implementa a criação de conta a pagar
type CreateContaPagarUseCase struct {
	repo   port.ContaPagarRepository
	logger *zap.Logger
}

// NewCreateContaPagarUseCase cria nova instância do use case
func NewCreateContaPagarUseCase(
	repo port.ContaPagarRepository,
	logger *zap.Logger,
) *CreateContaPagarUseCase {
	return &CreateContaPagarUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute cria uma nova conta a pagar
func (uc *CreateContaPagarUseCase) Execute(ctx context.Context, input CreateContaPagarInput) (*entity.ContaPagar, error) {
	// Validações de entrada
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	if input.Descricao == "" {
		return nil, fmt.Errorf("descrição é obrigatória")
	}

	if input.CategoriaID == "" {
		return nil, fmt.Errorf("categoria é obrigatória")
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
	conta, err := entity.NewContaPagar(
		tenantUUID,
		input.Descricao,
		input.CategoriaID,
		input.Fornecedor,
		input.Valor,
		input.Tipo,
		input.DataVencimento,
		input.Recorrente,
		input.Periodicidade,
	)
	if err != nil {
		uc.logger.Error("Erro ao criar entidade ContaPagar",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("descricao", input.Descricao),
		)
		return nil, fmt.Errorf("erro ao criar conta a pagar: %w", err)
	}

	// Atribuir campos opcionais
	conta.PixCode = input.PixCode
	conta.Observacoes = input.Observacoes

	// Persistir no repositório
	if err := uc.repo.Create(ctx, conta); err != nil {
		uc.logger.Error("Erro ao persistir conta a pagar",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("conta_id", conta.ID),
		)
		return nil, fmt.Errorf("erro ao salvar conta a pagar: %w", err)
	}

	uc.logger.Info("Conta a pagar criada com sucesso",
		zap.String("tenant_id", input.TenantID),
		zap.String("conta_id", conta.ID),
		zap.String("descricao", conta.Descricao),
	)

	return conta, nil
}
