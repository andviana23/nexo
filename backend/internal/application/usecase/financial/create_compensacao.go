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

// CreateCompensacaoInput define os dados de entrada para criar compensação
type CreateCompensacaoInput struct {
	TenantID        string
	ReceitaID       string
	MeioPagamentoID string
	DataTransacao   time.Time
	ValorBruto      valueobject.Money
	TaxaPercentual  valueobject.Percentage
	TaxaFixa        valueobject.Money
	DMais           valueobject.DMais
}

// CreateCompensacaoUseCase implementa a criação de compensação bancária
type CreateCompensacaoUseCase struct {
	repo   port.CompensacaoBancariaRepository
	logger *zap.Logger
}

// NewCreateCompensacaoUseCase cria nova instância do use case
func NewCreateCompensacaoUseCase(
	repo port.CompensacaoBancariaRepository,
	logger *zap.Logger,
) *CreateCompensacaoUseCase {
	return &CreateCompensacaoUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute cria uma nova compensação bancária
func (uc *CreateCompensacaoUseCase) Execute(ctx context.Context, input CreateCompensacaoInput) (*entity.CompensacaoBancaria, error) {
	// Validações de entrada
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	if input.ReceitaID == "" {
		return nil, fmt.Errorf("ID da receita é obrigatório")
	}

	if input.MeioPagamentoID == "" {
		return nil, fmt.Errorf("ID do meio de pagamento é obrigatório")
	}

	if input.DataTransacao.IsZero() {
		return nil, fmt.Errorf("data da transação é obrigatória")
	}

	// Converter TenantID para UUID
	tenantUUID, err := uuid.Parse(input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant_id inválido: %w", err)
	}

	// Criar entidade de domínio
	comp, err := entity.NewCompensacaoBancaria(
		tenantUUID,
		input.ReceitaID,
		input.MeioPagamentoID,
		input.DataTransacao,
		input.ValorBruto,
		input.TaxaPercentual,
		input.TaxaFixa,
		input.DMais,
	)
	if err != nil {
		uc.logger.Error("Erro ao criar entidade CompensacaoBancaria",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("receita_id", input.ReceitaID),
		)
		return nil, fmt.Errorf("erro ao criar compensação bancária: %w", err)
	}

	// Persistir no repositório
	if err := uc.repo.Create(ctx, comp); err != nil {
		uc.logger.Error("Erro ao persistir compensação bancária",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("comp_id", comp.ID),
		)
		return nil, fmt.Errorf("erro ao salvar compensação: %w", err)
	}

	uc.logger.Info("Compensação bancária criada",
		zap.String("tenant_id", input.TenantID),
		zap.String("comp_id", comp.ID),
		zap.String("receita_id", comp.ReceitaID),
		zap.String("data_compensacao", comp.DataCompensacao.Format("2006-01-02")),
	)

	return comp, nil
}
