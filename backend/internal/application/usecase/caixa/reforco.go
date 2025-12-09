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

// ReforcoInput define os dados de entrada para reforço
type ReforcoInput struct {
	TenantID  uuid.UUID
	UsuarioID uuid.UUID
	Valor     decimal.Decimal
	Origem    string // TROCO, CAPITAL_GIRO, TRANSFERENCIA, OUTROS
	Descricao string
}

// ReforcoUseCase implementa o registro de reforço
type ReforcoUseCase struct {
	repo   port.CaixaDiarioRepository
	logger *zap.Logger
}

// NewReforcoUseCase cria nova instância do use case
func NewReforcoUseCase(repo port.CaixaDiarioRepository, logger *zap.Logger) *ReforcoUseCase {
	return &ReforcoUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute registra um reforço no caixa
func (uc *ReforcoUseCase) Execute(ctx context.Context, input ReforcoInput) (*entity.OperacaoCaixa, error) {
	// Validações
	if input.TenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}
	if input.UsuarioID == uuid.Nil {
		return nil, fmt.Errorf("usuario_id é obrigatório")
	}
	if input.Valor.IsNegative() || input.Valor.IsZero() {
		return nil, domain.ErrValorInvalido
	}
	if input.Origem == "" {
		return nil, domain.ErrReforcoOrigemObrigatoria
	}

	// Buscar caixa aberto
	caixa, err := uc.repo.FindAberto(ctx, input.TenantID)
	if err != nil {
		if err == domain.ErrCaixaNaoAberto {
			return nil, domain.ErrCaixaNaoAberto
		}
		return nil, fmt.Errorf("erro ao buscar caixa aberto: %w", err)
	}

	// Criar operação de reforço
	operacao, err := entity.NewOperacaoReforco(
		caixa.ID,
		input.TenantID,
		input.UsuarioID,
		input.Valor,
		input.Origem,
		input.Descricao,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar operação de reforço: %w", err)
	}

	// Registrar na entidade caixa (atualiza totais)
	if err := caixa.RegistrarReforco(input.Valor); err != nil {
		return nil, fmt.Errorf("erro ao registrar reforço: %w", err)
	}

	// Persistir operação
	if err := uc.repo.CreateOperacao(ctx, operacao); err != nil {
		return nil, fmt.Errorf("erro ao salvar operação: %w", err)
	}

	// Atualizar totais do caixa
	if err := uc.repo.UpdateTotais(ctx, caixa.ID, input.TenantID, caixa.TotalSangrias, caixa.TotalReforcos, caixa.TotalEntradas); err != nil {
		return nil, fmt.Errorf("erro ao atualizar totais do caixa: %w", err)
	}

	uc.logger.Info("Reforço registrado com sucesso",
		zap.String("operacao_id", operacao.ID.String()),
		zap.String("caixa_id", caixa.ID.String()),
		zap.String("valor", input.Valor.String()),
		zap.String("origem", input.Origem),
	)

	return operacao, nil
}
