// Package caixa contém os use cases do módulo Caixa Diário
package caixa

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// AbrirCaixaInput define os dados de entrada para abrir caixa
type AbrirCaixaInput struct {
	TenantID     uuid.UUID
	UsuarioID    uuid.UUID
	SaldoInicial decimal.Decimal
}

// AbrirCaixaUseCase implementa a abertura de caixa
type AbrirCaixaUseCase struct {
	repo   port.CaixaDiarioRepository
	logger *zap.Logger
}

// NewAbrirCaixaUseCase cria nova instância do use case
func NewAbrirCaixaUseCase(repo port.CaixaDiarioRepository, logger *zap.Logger) *AbrirCaixaUseCase {
	return &AbrirCaixaUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute abre um novo caixa diário
func (uc *AbrirCaixaUseCase) Execute(ctx context.Context, input AbrirCaixaInput) (*entity.CaixaDiario, error) {
	// Validar tenant_id
	if input.TenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}

	// Validar usuario_id
	if input.UsuarioID == uuid.Nil {
		return nil, fmt.Errorf("usuario_id é obrigatório")
	}

	// Verificar se já existe caixa aberto (RN-CAI-001)
	caixaAberto, err := uc.repo.FindAberto(ctx, input.TenantID)
	if err != nil && err != domain.ErrCaixaNaoAberto {
		uc.logger.Error("Erro ao verificar caixa aberto",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID.String()),
		)
		return nil, fmt.Errorf("erro ao verificar caixa aberto: %w", err)
	}

	if caixaAberto != nil {
		uc.logger.Warn("Tentativa de abrir caixa com outro já aberto",
			zap.String("tenant_id", input.TenantID.String()),
			zap.String("caixa_aberto_id", caixaAberto.ID.String()),
		)
		return nil, domain.ErrCaixaJaAberto
	}

	// Criar novo caixa
	caixa, err := entity.NewCaixaDiario(input.TenantID, input.UsuarioID, input.SaldoInicial)
	if err != nil {
		uc.logger.Error("Erro ao criar entidade CaixaDiario",
			zap.Error(err),
		)
		return nil, fmt.Errorf("erro ao criar caixa: %w", err)
	}

	// Persistir
	if err := uc.repo.Create(ctx, caixa); err != nil {
		uc.logger.Error("Erro ao persistir caixa",
			zap.Error(err),
			zap.String("caixa_id", caixa.ID.String()),
		)
		return nil, fmt.Errorf("erro ao salvar caixa: %w", err)
	}

	uc.logger.Info("Caixa aberto com sucesso",
		zap.String("caixa_id", caixa.ID.String()),
		zap.String("tenant_id", input.TenantID.String()),
		zap.String("usuario_id", input.UsuarioID.String()),
		zap.String("saldo_inicial", input.SaldoInicial.String()),
	)

	return caixa, nil
}
