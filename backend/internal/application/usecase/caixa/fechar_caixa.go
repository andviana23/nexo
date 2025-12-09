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

// FecharCaixaInput define os dados de entrada para fechamento
type FecharCaixaInput struct {
	TenantID      uuid.UUID
	UsuarioID     uuid.UUID
	SaldoReal     decimal.Decimal
	Justificativa *string
}

// FecharCaixaUseCase implementa o fechamento de caixa
type FecharCaixaUseCase struct {
	repo   port.CaixaDiarioRepository
	logger *zap.Logger
}

// NewFecharCaixaUseCase cria nova instância do use case
func NewFecharCaixaUseCase(repo port.CaixaDiarioRepository, logger *zap.Logger) *FecharCaixaUseCase {
	return &FecharCaixaUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute fecha o caixa diário atual
func (uc *FecharCaixaUseCase) Execute(ctx context.Context, input FecharCaixaInput) (*entity.CaixaDiario, error) {
	// Validações
	if input.TenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}
	if input.UsuarioID == uuid.Nil {
		return nil, fmt.Errorf("usuario_id é obrigatório")
	}
	if input.SaldoReal.IsNegative() {
		return nil, fmt.Errorf("saldo real não pode ser negativo")
	}

	// Buscar caixa aberto
	caixa, err := uc.repo.FindAberto(ctx, input.TenantID)
	if err != nil {
		if err == domain.ErrCaixaNaoAberto {
			return nil, domain.ErrCaixaNaoAberto
		}
		return nil, fmt.Errorf("erro ao buscar caixa aberto: %w", err)
	}

	// Carregar operações para recálculo se necessário
	operacoes, err := uc.repo.ListOperacoes(ctx, caixa.ID, input.TenantID)
	if err != nil {
		uc.logger.Warn("Erro ao carregar operações para recálculo",
			zap.Error(err),
		)
		// Continuar sem operações - usar valores do caixa
	}
	caixa.Operacoes = operacoes

	// Fechar caixa (calcula divergência automaticamente)
	err = caixa.Fechar(input.UsuarioID, input.SaldoReal, input.Justificativa)
	if err != nil {
		// Verificar se é erro de justificativa obrigatória
		if err.Error() == "justificativa obrigatória para divergência maior que R$ 5,00" {
			return nil, domain.ErrCaixaJustificativaObrigatoria
		}
		return nil, fmt.Errorf("erro ao fechar caixa: %w", err)
	}

	// Persistir fechamento
	if err := uc.repo.Fechar(ctx, caixa); err != nil {
		return nil, fmt.Errorf("erro ao salvar fechamento: %w", err)
	}

	// Log com informações de divergência
	logFields := []zap.Field{
		zap.String("caixa_id", caixa.ID.String()),
		zap.String("tenant_id", input.TenantID.String()),
		zap.String("saldo_esperado", caixa.SaldoEsperado.String()),
		zap.String("saldo_real", input.SaldoReal.String()),
	}

	if caixa.Divergencia != nil {
		logFields = append(logFields, zap.String("divergencia", caixa.Divergencia.String()))
	}

	if caixa.TemDivergencia() {
		uc.logger.Warn("Caixa fechado com divergência", logFields...)
	} else {
		uc.logger.Info("Caixa fechado com sucesso", logFields...)
	}

	return caixa, nil
}
