package pricing

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SimularPrecoInput define os dados de entrada para simulação
type SimularPrecoInput struct {
	TenantID       string
	ItemID         string
	TipoItem       string // SERVICO ou PRODUTO
	CustoMateriais valueobject.Money
	CustoMaoDeObra valueobject.Money
	PrecoAtual     valueobject.Money
	Parametros     *entity.ParametrosSimulacao
}

// SimularPrecoUseCase simula um preço baseado nos custos e configuração
type SimularPrecoUseCase struct {
	configRepo port.PrecificacaoConfigRepository
	simRepo    port.PrecificacaoSimulacaoRepository
	logger     *zap.Logger
}

// NewSimularPrecoUseCase cria nova instância
func NewSimularPrecoUseCase(
	configRepo port.PrecificacaoConfigRepository,
	simRepo port.PrecificacaoSimulacaoRepository,
	logger *zap.Logger,
) *SimularPrecoUseCase {
	return &SimularPrecoUseCase{
		configRepo: configRepo,
		simRepo:    simRepo,
		logger:     logger,
	}
}

// Execute executa a simulação de preço
func (uc *SimularPrecoUseCase) Execute(ctx context.Context, input SimularPrecoInput) (*entity.PrecificacaoSimulacao, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	// Buscar configuração do tenant
	config, err := uc.configRepo.FindByTenantID(ctx, input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("configuração de precificação não encontrada: %w", err)
	}

	// Converter tenant_id de string para uuid.UUID
	tenantUUID, err := uuid.Parse(input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant_id inválido: %w", err)
	}

	// Criar simulação usando a configuração
	simulacao, err := entity.NewPrecificacaoSimulacao(
		tenantUUID,
		input.ItemID,
		input.TipoItem,
		input.CustoMateriais,
		input.CustoMaoDeObra,
		config.MargemDesejada,
		config.ComissaoPercentualDefault,
		config.ImpostoPercentual,
		input.PrecoAtual,
		input.Parametros,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar simulação: %w", err)
	}

	// Calcular preço sugerido
	simulacao.CalcularPrecoSugerido()

	uc.logger.Info("Simulação de preço executada",
		zap.String("tenant_id", input.TenantID),
		zap.String("item_id", input.ItemID),
		zap.String("preco_sugerido", simulacao.PrecoSugerido.String()),
	)

	return simulacao, nil
}
