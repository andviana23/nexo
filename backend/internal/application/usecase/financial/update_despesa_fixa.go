package financial

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"go.uber.org/zap"
)

// UpdateDespesaFixaInput define os dados de entrada para atualizar despesa fixa
type UpdateDespesaFixaInput struct {
	TenantID      string
	ID            string
	Descricao     string
	CategoriaID   string
	UnidadeID     string
	Fornecedor    string
	Valor         valueobject.Money
	DiaVencimento int
	Observacoes   string
}

// UpdateDespesaFixaUseCase implementa a atualização de despesa fixa
type UpdateDespesaFixaUseCase struct {
	repo   port.DespesaFixaRepository
	logger *zap.Logger
}

// NewUpdateDespesaFixaUseCase cria nova instância do use case
func NewUpdateDespesaFixaUseCase(
	repo port.DespesaFixaRepository,
	logger *zap.Logger,
) *UpdateDespesaFixaUseCase {
	return &UpdateDespesaFixaUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute atualiza uma despesa fixa existente
func (uc *UpdateDespesaFixaUseCase) Execute(ctx context.Context, input UpdateDespesaFixaInput) (*entity.DespesaFixa, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	if input.ID == "" {
		return nil, domain.ErrInvalidID
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

	// Buscar despesa existente
	despesa, err := uc.repo.FindByID(ctx, input.TenantID, input.ID)
	if err != nil {
		uc.logger.Error("Erro ao buscar despesa fixa para atualização",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("id", input.ID),
		)
		return nil, err
	}

	// Verificar duplicidade de descrição (excluindo o próprio registro)
	if input.Descricao != despesa.Descricao {
		exists, err := uc.repo.ExistsByDescricao(ctx, input.TenantID, input.Descricao, &input.ID)
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
	}

	// Atualizar campos
	err = despesa.Update(
		input.Descricao,
		input.Valor,
		input.DiaVencimento,
		input.CategoriaID,
		input.Fornecedor,
		input.UnidadeID,
		input.Observacoes,
	)
	if err != nil {
		return nil, err
	}

	// Persistir
	if err := uc.repo.Update(ctx, despesa); err != nil {
		uc.logger.Error("Erro ao atualizar despesa fixa",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("id", input.ID),
		)
		return nil, err
	}

	uc.logger.Info("Despesa fixa atualizada com sucesso",
		zap.String("id", despesa.ID),
		zap.String("tenant_id", despesa.TenantID.String()),
		zap.String("descricao", despesa.Descricao),
	)

	return despesa, nil
}
