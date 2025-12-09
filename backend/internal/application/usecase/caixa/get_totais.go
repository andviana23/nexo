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

// TotaisCaixa representa os totais do caixa por tipo
type TotaisCaixa struct {
	TotalVendas   decimal.Decimal
	TotalSangrias decimal.Decimal
	TotalReforcos decimal.Decimal
	TotalDespesas decimal.Decimal
	SaldoAtual    decimal.Decimal
}

// GetTotaisCaixaUseCase retorna os totais do caixa aberto
type GetTotaisCaixaUseCase struct {
	repo   port.CaixaDiarioRepository
	logger *zap.Logger
}

// NewGetTotaisCaixaUseCase cria nova instância do use case
func NewGetTotaisCaixaUseCase(repo port.CaixaDiarioRepository, logger *zap.Logger) *GetTotaisCaixaUseCase {
	return &GetTotaisCaixaUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute retorna os totais do caixa aberto
func (uc *GetTotaisCaixaUseCase) Execute(ctx context.Context, tenantID uuid.UUID) (*TotaisCaixa, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}

	// Buscar caixa aberto
	caixa, err := uc.repo.FindAberto(ctx, tenantID)
	if err != nil {
		if err == domain.ErrCaixaNaoAberto {
			return nil, domain.ErrCaixaNaoAberto
		}
		return nil, fmt.Errorf("erro ao buscar caixa aberto: %w", err)
	}

	// Buscar somas por tipo
	sums, err := uc.repo.SumOperacoesByTipo(ctx, caixa.ID, tenantID)
	if err != nil {
		uc.logger.Warn("Erro ao somar operações, usando valores do caixa",
			zap.Error(err),
		)
		// Usar valores do próprio caixa
		return &TotaisCaixa{
			TotalVendas:   caixa.TotalEntradas,
			TotalSangrias: caixa.TotalSangrias,
			TotalReforcos: caixa.TotalReforcos,
			TotalDespesas: decimal.Zero,
			SaldoAtual:    caixa.SaldoEsperado,
		}, nil
	}

	// Montar totais a partir das somas
	totais := &TotaisCaixa{
		TotalVendas:   sums[entity.TipoOperacaoVenda],
		TotalSangrias: sums[entity.TipoOperacaoSangria],
		TotalReforcos: sums[entity.TipoOperacaoReforco],
		TotalDespesas: sums[entity.TipoOperacaoDespesa],
	}

	// Calcular saldo atual
	// Saldo = Inicial + Vendas + Reforços - Sangrias - Despesas
	totais.SaldoAtual = caixa.SaldoInicial.
		Add(totais.TotalVendas).
		Add(totais.TotalReforcos).
		Sub(totais.TotalSangrias).
		Sub(totais.TotalDespesas)

	return totais, nil
}
