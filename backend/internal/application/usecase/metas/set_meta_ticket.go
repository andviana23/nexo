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

// SetMetaTicketInput define os dados de entrada
type SetMetaTicketInput struct {
	TenantID   string
	MesAno     valueobject.MesAno
	Tipo       valueobject.TipoMetaTicket
	BarbeiroID *string
	MetaValor  valueobject.Money
}

// SetMetaTicketUseCase implementa a definição de meta de ticket médio
type SetMetaTicketUseCase struct {
	repo   port.MetaTicketMedioRepository
	logger *zap.Logger
}

// NewSetMetaTicketUseCase cria nova instância do use case
func NewSetMetaTicketUseCase(
	repo port.MetaTicketMedioRepository,
	logger *zap.Logger,
) *SetMetaTicketUseCase {
	return &SetMetaTicketUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute cria ou atualiza uma meta de ticket médio
func (uc *SetMetaTicketUseCase) Execute(ctx context.Context, input SetMetaTicketInput) (*entity.MetaTicketMedio, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	// Validar tipo e barbeiroID
	if input.Tipo == valueobject.TipoMetaTicketBarbeiro && (input.BarbeiroID == nil || *input.BarbeiroID == "") {
		return nil, fmt.Errorf("barbeiro ID é obrigatório para meta de ticket de barbeiro")
	}

	// Buscar meta existente
	var existing *entity.MetaTicketMedio
	var err error

	if input.Tipo == valueobject.TipoMetaTicketBarbeiro {
		existing, err = uc.repo.FindBarbeiroByMesAno(ctx, input.TenantID, *input.BarbeiroID, input.MesAno)
	} else {
		existing, err = uc.repo.FindGeralByMesAno(ctx, input.TenantID, input.MesAno)
	}

	if err == nil && existing != nil {
		// Atualizar valor da meta (atribuição direta)
		existing.MetaValor = input.MetaValor
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
	meta, err := entity.NewMetaTicketMedio(
		tenantUUID,
		input.MesAno,
		input.Tipo,
		input.MetaValor,
		input.BarbeiroID,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar meta de ticket médio: %w", err)
	}

	if err := uc.repo.Create(ctx, meta); err != nil {
		return nil, fmt.Errorf("erro ao salvar meta de ticket médio: %w", err)
	}

	uc.logger.Info("Meta de ticket médio definida",
		zap.String("tenant_id", input.TenantID),
		zap.String("tipo", string(input.Tipo)),
		zap.String("mes_ano", input.MesAno.String()),
	)

	return meta, nil
}
