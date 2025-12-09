package metas

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

// SetMetaBarbeiroInput define os dados de entrada
type SetMetaBarbeiroInput struct {
	TenantID           string
	BarbeiroID         string
	MesAno             valueobject.MesAno
	MetaServicosGerais valueobject.Money
	MetaServicosExtras valueobject.Money
	MetaProdutos       valueobject.Money
}

// SetMetaBarbeiroUseCase implementa a definição de meta por barbeiro
type SetMetaBarbeiroUseCase struct {
	repo   port.MetaBarbeiroRepository
	logger *zap.Logger
}

// NewSetMetaBarbeiroUseCase cria nova instância do use case
func NewSetMetaBarbeiroUseCase(
	repo port.MetaBarbeiroRepository,
	logger *zap.Logger,
) *SetMetaBarbeiroUseCase {
	return &SetMetaBarbeiroUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute cria ou atualiza uma meta de barbeiro
func (uc *SetMetaBarbeiroUseCase) Execute(ctx context.Context, input SetMetaBarbeiroInput) (*entity.MetaBarbeiro, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if input.BarbeiroID == "" {
		return nil, fmt.Errorf("barbeiro ID é obrigatório")
	}

	// Buscar meta existente
	existing, err := uc.repo.FindByBarbeiroMesAno(ctx, input.TenantID, input.BarbeiroID, input.MesAno)
	if err == nil && existing != nil {
		// Atualizar metas (atribuição direta)
		existing.MetaServicosGerais = input.MetaServicosGerais
		existing.MetaServicosExtras = input.MetaServicosExtras
		existing.MetaProdutos = input.MetaProdutos
		existing.AtualizadoEm = time.Now()

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
	meta, err := entity.NewMetaBarbeiro(
		tenantUUID,
		input.BarbeiroID,
		input.MesAno,
		input.MetaServicosGerais,
		input.MetaServicosExtras,
		input.MetaProdutos,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar meta de barbeiro: %w", err)
	}

	if err := uc.repo.Create(ctx, meta); err != nil {
		return nil, fmt.Errorf("erro ao salvar meta de barbeiro: %w", err)
	}

	uc.logger.Info("Meta de barbeiro definida",
		zap.String("tenant_id", input.TenantID),
		zap.String("barbeiro_id", input.BarbeiroID),
		zap.String("mes_ano", input.MesAno.String()),
	)

	return meta, nil
}
