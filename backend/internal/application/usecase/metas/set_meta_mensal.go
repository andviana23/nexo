package metas

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

// SetMetaMensalInput define os dados de entrada
type SetMetaMensalInput struct {
	TenantID        string
	MesAno          valueobject.MesAno
	MetaFaturamento valueobject.Money
	Origem          valueobject.OrigemMeta
}

// SetMetaMensalUseCase implementa a definição de meta mensal
type SetMetaMensalUseCase struct {
	repo   port.MetaMensalRepository
	logger *zap.Logger
}

// NewSetMetaMensalUseCase cria nova instância do use case
func NewSetMetaMensalUseCase(
	repo port.MetaMensalRepository,
	logger *zap.Logger,
) *SetMetaMensalUseCase {
	return &SetMetaMensalUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute cria ou atualiza uma meta mensal
func (uc *SetMetaMensalUseCase) Execute(ctx context.Context, input SetMetaMensalInput) (*entity.MetaMensal, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	// Buscar meta existente para o mês
	existing, err := uc.repo.FindByMesAno(ctx, input.TenantID, input.MesAno)
	if err == nil && existing != nil {
		// Atualizar meta existente
		if err := existing.AtualizarMeta(input.MetaFaturamento); err != nil {
			return nil, fmt.Errorf("erro ao atualizar meta: %w", err)
		}
		if err := uc.repo.Update(ctx, existing); err != nil {
			return nil, fmt.Errorf("erro ao salvar meta atualizada: %w", err)
		}
		return existing, nil
	}

	// Converter tenant_id de string para uuid.UUID
	tenantUUID, err := uuid.Parse(input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant_id inválido: %w", err)
	}

	// Criar nova meta
	meta, err := entity.NewMetaMensal(tenantUUID, input.MesAno, input.MetaFaturamento, input.Origem)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar meta mensal: %w", err)
	}

	if err := uc.repo.Create(ctx, meta); err != nil {
		return nil, fmt.Errorf("erro ao salvar meta mensal: %w", err)
	}

	uc.logger.Info("Meta mensal definida",
		zap.String("tenant_id", input.TenantID),
		zap.String("mes_ano", input.MesAno.String()),
		zap.String("meta", input.MetaFaturamento.String()),
	)

	return meta, nil
}
