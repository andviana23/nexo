package caixa

import (
	"context"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ListHistoricoInput define os filtros para listagem do histórico
type ListHistoricoInput struct {
	TenantID   uuid.UUID
	DataInicio *time.Time
	DataFim    *time.Time
	UsuarioID  *uuid.UUID
	Page       int
	PageSize   int
}

// ListHistoricoOutput representa a saída paginada
type ListHistoricoOutput struct {
	Caixas []*entity.CaixaDiario
	Total  int64
}

// ListHistoricoUseCase lista o histórico de caixas fechados
type ListHistoricoUseCase struct {
	repo   port.CaixaDiarioRepository
	logger *zap.Logger
}

// NewListHistoricoUseCase cria nova instância do use case
func NewListHistoricoUseCase(repo port.CaixaDiarioRepository, logger *zap.Logger) *ListHistoricoUseCase {
	return &ListHistoricoUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute retorna o histórico de caixas fechados
func (uc *ListHistoricoUseCase) Execute(ctx context.Context, input ListHistoricoInput) (*ListHistoricoOutput, error) {
	if input.TenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}

	// Aplicar defaults
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}
	if input.PageSize > 100 {
		input.PageSize = 100
	}

	filters := port.CaixaFilters{
		DataInicio: input.DataInicio,
		DataFim:    input.DataFim,
		UsuarioID:  input.UsuarioID,
		Limit:      input.PageSize,
		Offset:     (input.Page - 1) * input.PageSize,
	}

	// Buscar caixas
	caixas, err := uc.repo.ListHistorico(ctx, input.TenantID, filters)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar histórico: %w", err)
	}

	// Contar total
	total, err := uc.repo.CountHistorico(ctx, input.TenantID, filters)
	if err != nil {
		uc.logger.Warn("Erro ao contar histórico, usando len(caixas)",
			zap.Error(err),
		)
		total = int64(len(caixas))
	}

	return &ListHistoricoOutput{
		Caixas: caixas,
		Total:  total,
	}, nil
}
