package financial

import (
	"context"

	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// PainelMensalInput representa os parâmetros de entrada para o painel mensal
type PainelMensalInput struct {
	TenantID  string
	UnidadeID *string
	Ano       int
	Mes       int
}

// PainelMensalOutput representa os dados consolidados do painel mensal
type PainelMensalOutput struct {
	// Período
	Ano     int
	Mes     int
	NomeMes string

	// Receitas
	ReceitaRealizada valueobject.Money
	ReceitaPendente  valueobject.Money
	ReceitaTotal     valueobject.Money

	// Despesas
	DespesasFixas     valueobject.Money
	DespesasVariaveis valueobject.Money
	DespesasPagas     valueobject.Money
	DespesasPendentes valueobject.Money
	DespesasTotal     valueobject.Money

	// Resultados
	LucroBruto    valueobject.Money
	LucroLiquido  valueobject.Money
	MargemLiquida decimal.Decimal // Em percentual

	// Metas
	MetaMensal     valueobject.Money
	PercentualMeta decimal.Decimal // Em percentual
	DiferencaMeta  valueobject.Money
	StatusMeta     string // "Atingida", "Em andamento", "Abaixo"

	// Caixa
	SaldoCaixaAtual valueobject.Money

	// Comparativo
	VariacaoMesAnterior decimal.Decimal // Em percentual
	TendenciaVariacao   string          // "up", "down", "stable"
}

// GetPainelMensalUseCase busca e agrega dados financeiros do mês
type GetPainelMensalUseCase struct {
	contaPagarRepo   port.ContaPagarRepository
	contaReceberRepo port.ContaReceberRepository
	despesaFixaRepo  port.DespesaFixaRepository
	metaMensalRepo   port.MetaMensalRepository
	fluxoCaixaRepo   port.FluxoCaixaDiarioRepository
	logger           *zap.Logger
}

// NewGetPainelMensalUseCase cria uma nova instância do use case
func NewGetPainelMensalUseCase(
	contaPagarRepo port.ContaPagarRepository,
	contaReceberRepo port.ContaReceberRepository,
	despesaFixaRepo port.DespesaFixaRepository,
	metaMensalRepo port.MetaMensalRepository,
	fluxoCaixaRepo port.FluxoCaixaDiarioRepository,
	logger *zap.Logger,
) *GetPainelMensalUseCase {
	return &GetPainelMensalUseCase{
		contaPagarRepo:   contaPagarRepo,
		contaReceberRepo: contaReceberRepo,
		despesaFixaRepo:  despesaFixaRepo,
		metaMensalRepo:   metaMensalRepo,
		fluxoCaixaRepo:   fluxoCaixaRepo,
		logger:           logger,
	}
}

// Execute executa a lógica de agregação do painel mensal
func (uc *GetPainelMensalUseCase) Execute(ctx context.Context, input PainelMensalInput) (*PainelMensalOutput, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	// Validar ano e mês
	if input.Ano < 2000 || input.Ano > 2100 {
		return nil, fmt.Errorf("ano inválido: %d", input.Ano)
	}
	if input.Mes < 1 || input.Mes > 12 {
		return nil, fmt.Errorf("mês inválido: %d", input.Mes)
	}

	// 1. Definir período do mês
	inicio := time.Date(input.Ano, time.Month(input.Mes), 1, 0, 0, 0, 0, time.UTC)
	fim := inicio.AddDate(0, 1, 0).Add(-time.Nanosecond) // Último instante do mês

	// 2. Buscar receitas realizadas (RECEBIDO)
	statusRecebido := valueobject.StatusContaRecebido
	receitaRealizada, err := uc.contaReceberRepo.SumByPeriod(ctx, input.TenantID, inicio, fim, &statusRecebido)
	if err != nil {
		uc.logger.Error("erro ao buscar receita realizada", zap.Error(err))
		receitaRealizada = valueobject.Zero()
	}

	// 3. Buscar receitas pendentes
	statusPendente := valueobject.StatusContaPendente
	receitaPendente, err := uc.contaReceberRepo.SumByPeriod(ctx, input.TenantID, inicio, fim, &statusPendente)
	if err != nil {
		uc.logger.Error("erro ao buscar receita pendente", zap.Error(err))
		receitaPendente = valueobject.Zero()
	}

	// 4. Buscar despesas pagas
	statusPago := valueobject.StatusContaPago
	despesasPagas, err := uc.contaPagarRepo.SumByPeriod(ctx, input.TenantID, inicio, fim, &statusPago)
	if err != nil {
		uc.logger.Error("erro ao buscar despesas pagas", zap.Error(err))
		despesasPagas = valueobject.Zero()
	}

	// 5. Buscar despesas pendentes
	despesasPendentes, err := uc.contaPagarRepo.SumByPeriod(ctx, input.TenantID, inicio, fim, &statusPendente)
	if err != nil {
		uc.logger.Error("erro ao buscar despesas pendentes", zap.Error(err))
		despesasPendentes = valueobject.Zero()
	}

	// 6. Buscar despesas fixas ativas
	despesasFixas, err := uc.despesaFixaRepo.SumAtivas(ctx, input.TenantID)
	if err != nil {
		uc.logger.Error("erro ao buscar despesas fixas", zap.Error(err))
		despesasFixas = valueobject.Zero()
	}

	// 7. Buscar meta do mês
	mesAno, _ := valueobject.NewMesAnoFromInts(input.Ano, input.Mes)
	meta, err := uc.metaMensalRepo.FindByMesAno(ctx, input.TenantID, mesAno)
	var metaMensal valueobject.Money
	if err != nil || meta == nil {
		metaMensal = valueobject.Zero()
	} else {
		metaMensal = meta.MetaFaturamento
	}

	// 8. Buscar saldo do caixa mais recente
	hoje := time.Now()
	caixa, err := uc.fluxoCaixaRepo.FindByData(ctx, input.TenantID, hoje)
	var saldoCaixaAtual valueobject.Money
	if err != nil || caixa == nil {
		saldoCaixaAtual = valueobject.Zero()
	} else {
		saldoCaixaAtual = caixa.SaldoFinal
	}

	// 9. Calcular totais
	receitaTotal := receitaRealizada.Add(receitaPendente)
	despesasTotal := despesasPagas.Add(despesasPendentes)

	// Despesas variáveis = total pagas - fixas
	despesasVariaveis := despesasPagas.Sub(despesasFixas)
	if despesasVariaveis.IsNegative() {
		despesasVariaveis = valueobject.Zero()
	}

	// 10. Calcular lucro
	lucroBruto := receitaRealizada.Sub(despesasPagas)
	lucroLiquido := lucroBruto // Simplificado (sem IR por enquanto)

	// 11. Calcular margem líquida
	var margemLiquida decimal.Decimal
	if !receitaRealizada.IsZero() {
		margemLiquida = lucroLiquido.Value().Div(receitaRealizada.Value()).Mul(decimal.NewFromInt(100))
	}

	// 12. Calcular percentual e diferença da meta
	var percentualMeta decimal.Decimal
	var diferencaMeta valueobject.Money
	var statusMeta string

	if !metaMensal.IsZero() {
		percentualMeta = receitaRealizada.Value().Div(metaMensal.Value()).Mul(decimal.NewFromInt(100))
		diferencaMeta = receitaRealizada.Sub(metaMensal)

		if percentualMeta.GreaterThanOrEqual(decimal.NewFromInt(100)) {
			statusMeta = "Atingida"
		} else if percentualMeta.GreaterThanOrEqual(decimal.NewFromInt(80)) {
			statusMeta = "Em andamento"
		} else {
			statusMeta = "Abaixo"
		}
	} else {
		diferencaMeta = valueobject.Zero()
		statusMeta = "Sem meta"
	}

	// 13. Calcular variação vs mês anterior
	mesAnteriorInicio := inicio.AddDate(0, -1, 0)
	mesAnteriorFim := inicio.Add(-time.Nanosecond)

	receitaMesAnterior, err := uc.contaReceberRepo.SumByPeriod(ctx, input.TenantID, mesAnteriorInicio, mesAnteriorFim, &statusPago)
	if err != nil {
		receitaMesAnterior = valueobject.Zero()
	}

	var variacaoMesAnterior decimal.Decimal
	var tendenciaVariacao string

	if !receitaMesAnterior.IsZero() {
		diff := receitaRealizada.Sub(receitaMesAnterior)
		variacaoMesAnterior = diff.Value().Div(receitaMesAnterior.Value()).Mul(decimal.NewFromInt(100))

		if variacaoMesAnterior.GreaterThan(decimal.NewFromFloat(2)) {
			tendenciaVariacao = "up"
		} else if variacaoMesAnterior.LessThan(decimal.NewFromFloat(-2)) {
			tendenciaVariacao = "down"
		} else {
			tendenciaVariacao = "stable"
		}
	} else {
		tendenciaVariacao = "stable"
	}

	return &PainelMensalOutput{
		Ano:                 input.Ano,
		Mes:                 input.Mes,
		NomeMes:             nomeMes(input.Mes),
		ReceitaRealizada:    receitaRealizada,
		ReceitaPendente:     receitaPendente,
		ReceitaTotal:        receitaTotal,
		DespesasFixas:       despesasFixas,
		DespesasVariaveis:   despesasVariaveis,
		DespesasPagas:       despesasPagas,
		DespesasPendentes:   despesasPendentes,
		DespesasTotal:       despesasTotal,
		LucroBruto:          lucroBruto,
		LucroLiquido:        lucroLiquido,
		MargemLiquida:       margemLiquida,
		MetaMensal:          metaMensal,
		PercentualMeta:      percentualMeta,
		DiferencaMeta:       diferencaMeta,
		StatusMeta:          statusMeta,
		SaldoCaixaAtual:     saldoCaixaAtual,
		VariacaoMesAnterior: variacaoMesAnterior,
		TendenciaVariacao:   tendenciaVariacao,
	}, nil
}

// nomeMes retorna o nome do mês em português
func nomeMes(mes int) string {
	nomes := []string{
		"", // índice 0 não usado
		"Janeiro",
		"Fevereiro",
		"Março",
		"Abril",
		"Maio",
		"Junho",
		"Julho",
		"Agosto",
		"Setembro",
		"Outubro",
		"Novembro",
		"Dezembro",
	}
	if mes >= 1 && mes <= 12 {
		return nomes[mes]
	}
	return ""
}
