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

// SangriaInput define os dados de entrada para sangria
type SangriaInput struct {
	TenantID  uuid.UUID
	UsuarioID uuid.UUID
	Valor     decimal.Decimal
	Destino   string // DEPOSITO, PAGAMENTO, COFRE, OUTROS
	Descricao string
}

// SangriaUseCase implementa o registro de sangria
type SangriaUseCase struct {
	repo   port.CaixaDiarioRepository
	logger *zap.Logger
}

// NewSangriaUseCase cria nova instância do use case
func NewSangriaUseCase(repo port.CaixaDiarioRepository, logger *zap.Logger) *SangriaUseCase {
	return &SangriaUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute registra uma sangria no caixa
func (uc *SangriaUseCase) Execute(ctx context.Context, input SangriaInput) (*entity.OperacaoCaixa, error) {
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
	if input.Destino == "" {
		return nil, domain.ErrSangriaDestinoObrigatorio
	}

	// Buscar caixa aberto
	caixa, err := uc.repo.FindAberto(ctx, input.TenantID)
	if err != nil {
		if err == domain.ErrCaixaNaoAberto {
			return nil, domain.ErrCaixaNaoAberto
		}
		return nil, fmt.Errorf("erro ao buscar caixa aberto: %w", err)
	}

	// Criar operação de sangria
	operacao, err := entity.NewOperacaoSangria(
		caixa.ID,
		input.TenantID,
		input.UsuarioID,
		input.Valor,
		input.Destino,
		input.Descricao,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar operação de sangria: %w", err)
	}

	// Registrar na entidade caixa (atualiza totais)
	if err := caixa.RegistrarSangria(input.Valor); err != nil {
		return nil, fmt.Errorf("erro ao registrar sangria: %w", err)
	}

	// Persistir operação
	if err := uc.repo.CreateOperacao(ctx, operacao); err != nil {
		return nil, fmt.Errorf("erro ao salvar operação: %w", err)
	}

	// Atualizar totais do caixa
	if err := uc.repo.UpdateTotais(ctx, caixa.ID, input.TenantID, caixa.TotalSangrias, caixa.TotalReforcos, caixa.TotalEntradas); err != nil {
		return nil, fmt.Errorf("erro ao atualizar totais do caixa: %w", err)
	}

	uc.logger.Info("Sangria registrada com sucesso",
		zap.String("operacao_id", operacao.ID.String()),
		zap.String("caixa_id", caixa.ID.String()),
		zap.String("valor", input.Valor.String()),
		zap.String("destino", input.Destino),
	)

	return operacao, nil
}
