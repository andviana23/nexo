package port

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CaixaDiarioRepository define operações de persistência para caixa diário
type CaixaDiarioRepository interface {
	// Create cria um novo caixa diário (abertura)
	Create(ctx context.Context, caixa *entity.CaixaDiario) error

	// FindByID busca um caixa por ID
	FindByID(ctx context.Context, caixaID, tenantID uuid.UUID) (*entity.CaixaDiario, error)

	// FindAberto busca o caixa aberto do tenant (deve existir apenas 1)
	FindAberto(ctx context.Context, tenantID uuid.UUID) (*entity.CaixaDiario, error)

	// Update atualiza um caixa existente
	Update(ctx context.Context, caixa *entity.CaixaDiario) error

	// UpdateTotais atualiza apenas os totais do caixa (após sangria, reforço, venda)
	UpdateTotais(ctx context.Context, caixaID, tenantID uuid.UUID, sangrias, reforcos, entradas decimal.Decimal) error

	// Fechar fecha o caixa com os valores de fechamento
	Fechar(ctx context.Context, caixa *entity.CaixaDiario) error

	// ListHistorico lista caixas fechados com paginação
	ListHistorico(ctx context.Context, tenantID uuid.UUID, filters CaixaFilters) ([]*entity.CaixaDiario, error)

	// CountHistorico conta o total de caixas fechados
	CountHistorico(ctx context.Context, tenantID uuid.UUID, filters CaixaFilters) (int64, error)

	// ======== Operações ========

	// CreateOperacao registra uma operação no caixa
	CreateOperacao(ctx context.Context, op *entity.OperacaoCaixa) error

	// ListOperacoes lista todas as operações de um caixa
	ListOperacoes(ctx context.Context, caixaID, tenantID uuid.UUID) ([]entity.OperacaoCaixa, error)

	// ListOperacoesByTipo lista operações filtradas por tipo
	ListOperacoesByTipo(ctx context.Context, caixaID, tenantID uuid.UUID, tipo entity.TipoOperacaoCaixa) ([]entity.OperacaoCaixa, error)

	// SumOperacoesByTipo soma os valores de operações por tipo
	SumOperacoesByTipo(ctx context.Context, caixaID, tenantID uuid.UUID) (map[entity.TipoOperacaoCaixa]decimal.Decimal, error)
}

// CaixaFilters representa filtros para busca de caixas
type CaixaFilters struct {
	DataInicio *time.Time // Filtrar por data de abertura >= DataInicio
	DataFim    *time.Time // Filtrar por data de abertura <= DataFim
	UsuarioID  *uuid.UUID // Filtrar por usuário de abertura ou fechamento
	Limit      int
	Offset     int
}

// DefaultCaixaFilters retorna filtros padrão
func DefaultCaixaFilters() CaixaFilters {
	return CaixaFilters{
		Limit:  20,
		Offset: 0,
	}
}
