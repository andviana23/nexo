package caixa

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GetCaixaAbertoUseCase retorna o caixa aberto do tenant
type GetCaixaAbertoUseCase struct {
	repo   port.CaixaDiarioRepository
	logger *zap.Logger
}

// NewGetCaixaAbertoUseCase cria nova instância do use case
func NewGetCaixaAbertoUseCase(repo port.CaixaDiarioRepository, logger *zap.Logger) *GetCaixaAbertoUseCase {
	return &GetCaixaAbertoUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute retorna o caixa aberto com operações
func (uc *GetCaixaAbertoUseCase) Execute(ctx context.Context, tenantID uuid.UUID) (*entity.CaixaDiario, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}

	caixa, err := uc.repo.FindAberto(ctx, tenantID)
	if err != nil {
		if err == domain.ErrCaixaNaoAberto {
			return nil, nil // Retorna nil sem erro = sem caixa aberto
		}
		return nil, fmt.Errorf("erro ao buscar caixa aberto: %w", err)
	}

	// Carregar operações
	operacoes, err := uc.repo.ListOperacoes(ctx, caixa.ID, tenantID)
	if err != nil {
		uc.logger.Warn("Erro ao carregar operações do caixa",
			zap.Error(err),
			zap.String("caixa_id", caixa.ID.String()),
		)
		// Continuar sem operações
	} else {
		caixa.Operacoes = operacoes
	}

	return caixa, nil
}

// GetCaixaByIDUseCase retorna um caixa específico por ID
type GetCaixaByIDUseCase struct {
	repo   port.CaixaDiarioRepository
	logger *zap.Logger
}

// NewGetCaixaByIDUseCase cria nova instância do use case
func NewGetCaixaByIDUseCase(repo port.CaixaDiarioRepository, logger *zap.Logger) *GetCaixaByIDUseCase {
	return &GetCaixaByIDUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute retorna um caixa por ID com operações
func (uc *GetCaixaByIDUseCase) Execute(ctx context.Context, caixaID, tenantID uuid.UUID) (*entity.CaixaDiario, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}
	if caixaID == uuid.Nil {
		return nil, fmt.Errorf("caixa_id é obrigatório")
	}

	caixa, err := uc.repo.FindByID(ctx, caixaID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar caixa: %w", err)
	}

	// Carregar operações
	operacoes, err := uc.repo.ListOperacoes(ctx, caixa.ID, tenantID)
	if err != nil {
		uc.logger.Warn("Erro ao carregar operações do caixa",
			zap.Error(err),
			zap.String("caixa_id", caixa.ID.String()),
		)
	} else {
		caixa.Operacoes = operacoes
	}

	return caixa, nil
}
