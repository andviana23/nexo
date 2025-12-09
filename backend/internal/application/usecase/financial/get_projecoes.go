package financial

import (
	"context"

	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// ProjecoesInput representa os parâmetros de entrada para projeções
type ProjecoesInput struct {
	TenantID   string
	MesesAhead int // Quantos meses projetar (default: 3, max: 12)
}

// ProjecaoMensal representa a projeção de um mês específico
type ProjecaoMensal struct {
	Ano                int
	Mes                int
	NomeMes            string
	ReceitaProjetada   valueobject.Money
	DespesasProjetadas valueobject.Money
	DespesasFixas      valueobject.Money
	LucroProjetado     valueobject.Money
	DiasUteis          int
	MetaDiaria         valueobject.Money
	Confianca          string // "Alta", "Média", "Baixa"
}

// ProjecoesOutput representa o resultado das projeções financeiras
type ProjecoesOutput struct {
	Projecoes           []ProjecaoMensal
	MediaReceita3Meses  valueobject.Money
	MediaDespesas3Meses valueobject.Money
	TendenciaReceita    string // "Crescente", "Estável", "Decrescente"
	DataGeracao         time.Time
}

// GetProjecoesUseCase calcula projeções financeiras baseadas no histórico
type GetProjecoesUseCase struct {
	contaPagarRepo   port.ContaPagarRepository
	contaReceberRepo port.ContaReceberRepository
	despesaFixaRepo  port.DespesaFixaRepository
	logger           *zap.Logger
}

// NewGetProjecoesUseCase cria uma nova instância do use case
func NewGetProjecoesUseCase(
	contaPagarRepo port.ContaPagarRepository,
	contaReceberRepo port.ContaReceberRepository,
	despesaFixaRepo port.DespesaFixaRepository,
	logger *zap.Logger,
) *GetProjecoesUseCase {
	return &GetProjecoesUseCase{
		contaPagarRepo:   contaPagarRepo,
		contaReceberRepo: contaReceberRepo,
		despesaFixaRepo:  despesaFixaRepo,
		logger:           logger,
	}
}

// Execute executa a lógica de projeção financeira
func (uc *GetProjecoesUseCase) Execute(ctx context.Context, input ProjecoesInput) (*ProjecoesOutput, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	// Validar e normalizar meses
	if input.MesesAhead < 1 {
		input.MesesAhead = 3
	}
	if input.MesesAhead > 12 {
		input.MesesAhead = 12
	}

	hoje := time.Now()
	statusPago := valueobject.StatusContaPago

	// 1. Buscar histórico dos últimos 3 meses
	var historicoReceitas []valueobject.Money
	var historicoDespesas []valueobject.Money

	for i := 1; i <= 3; i++ {
		mesRef := hoje.AddDate(0, -i, 0)
		inicio := time.Date(mesRef.Year(), mesRef.Month(), 1, 0, 0, 0, 0, time.UTC)
		fim := inicio.AddDate(0, 1, 0).Add(-time.Nanosecond)

		receita, err := uc.contaReceberRepo.SumByPeriod(ctx, input.TenantID, inicio, fim, &statusPago)
		if err != nil {
			uc.logger.Warn("erro ao buscar receita do histórico", zap.Error(err), zap.Int("mes", i))
			receita = valueobject.Zero()
		}

		despesa, err := uc.contaPagarRepo.SumByPeriod(ctx, input.TenantID, inicio, fim, &statusPago)
		if err != nil {
			uc.logger.Warn("erro ao buscar despesa do histórico", zap.Error(err), zap.Int("mes", i))
			despesa = valueobject.Zero()
		}

		historicoReceitas = append(historicoReceitas, receita)
		historicoDespesas = append(historicoDespesas, despesa)
	}

	// 2. Calcular médias
	mediaReceita := calcularMediaMoney(historicoReceitas)
	mediaDespesa := calcularMediaMoney(historicoDespesas)

	// 3. Detectar tendência de receitas
	tendencia := detectarTendenciaMoney(historicoReceitas)

	// 4. Buscar despesas fixas (garantidas)
	despesasFixas, err := uc.despesaFixaRepo.SumAtivas(ctx, input.TenantID)
	if err != nil {
		uc.logger.Warn("erro ao buscar despesas fixas", zap.Error(err))
		despesasFixas = valueobject.Zero()
	}

	// 5. Gerar projeções
	var projecoes []ProjecaoMensal

	for i := 1; i <= input.MesesAhead; i++ {
		mesProjetado := hoje.AddDate(0, i, 0)
		ano := mesProjetado.Year()
		mes := int(mesProjetado.Month())

		// Aplicar fator de tendência
		fatorTendencia := calcularFatorTendenciaMoney(tendencia, i)
		receitaProjetada := mediaReceita.Mul(fatorTendencia)

		// Despesas variáveis estimadas = média despesas - despesas fixas
		despesasVariaveis := mediaDespesa.Sub(despesasFixas)
		if despesasVariaveis.IsNegative() {
			despesasVariaveis = valueobject.Zero()
		}

		// Total despesas projetadas = fixas + variáveis
		despesasProjetadas := despesasFixas.Add(despesasVariaveis)

		// Lucro projetado
		lucroProjetado := receitaProjetada.Sub(despesasProjetadas)

		// Dias úteis do mês
		diasUteis := calcularDiasUteisMes(ano, mes)

		// Meta diária para atingir a projeção
		var metaDiaria valueobject.Money
		if diasUteis > 0 {
			metaDiaria = receitaProjetada.Div(decimal.NewFromInt(int64(diasUteis)))
		} else {
			metaDiaria = valueobject.Zero()
		}

		// Nível de confiança diminui com a distância
		confianca := calcularConfiancaProjecao(i, tendencia)

		projecoes = append(projecoes, ProjecaoMensal{
			Ano:                ano,
			Mes:                mes,
			NomeMes:            nomeMes(mes),
			ReceitaProjetada:   receitaProjetada,
			DespesasProjetadas: despesasProjetadas,
			DespesasFixas:      despesasFixas,
			LucroProjetado:     lucroProjetado,
			DiasUteis:          diasUteis,
			MetaDiaria:         metaDiaria,
			Confianca:          confianca,
		})
	}

	return &ProjecoesOutput{
		Projecoes:           projecoes,
		MediaReceita3Meses:  mediaReceita,
		MediaDespesas3Meses: mediaDespesa,
		TendenciaReceita:    tendencia,
		DataGeracao:         hoje,
	}, nil
}

// calcularMediaMoney calcula a média de uma lista de Money
func calcularMediaMoney(valores []valueobject.Money) valueobject.Money {
	if len(valores) == 0 {
		return valueobject.Zero()
	}

	soma := valueobject.Zero()
	for _, v := range valores {
		soma = soma.Add(v)
	}

	return soma.Div(decimal.NewFromInt(int64(len(valores))))
}

// detectarTendenciaMoney analisa a tendência de valores Money
func detectarTendenciaMoney(valores []valueobject.Money) string {
	if len(valores) < 2 {
		return "Estável"
	}

	// valores[0] = mês mais recente, valores[len-1] = mês mais antigo
	maisRecente := valores[0].Value()
	maisAntigo := valores[len(valores)-1].Value()

	if maisAntigo.IsZero() {
		return "Estável"
	}

	diff := maisRecente.Sub(maisAntigo)
	percentual := diff.Div(maisAntigo).Mul(decimal.NewFromInt(100))

	if percentual.GreaterThan(decimal.NewFromInt(5)) {
		return "Crescente"
	} else if percentual.LessThan(decimal.NewFromInt(-5)) {
		return "Decrescente"
	}
	return "Estável"
}

// calcularFatorTendenciaMoney retorna o fator multiplicador baseado na tendência
func calcularFatorTendenciaMoney(tendencia string, mesAhead int) decimal.Decimal {
	base := decimal.NewFromInt(1)
	fator := decimal.NewFromFloat(0.02) // 2% por mês

	switch tendencia {
	case "Crescente":
		return base.Add(fator.Mul(decimal.NewFromInt(int64(mesAhead))))
	case "Decrescente":
		resultado := base.Sub(fator.Mul(decimal.NewFromInt(int64(mesAhead))))
		// Não deixar o fator ficar negativo ou muito baixo
		if resultado.LessThan(decimal.NewFromFloat(0.5)) {
			return decimal.NewFromFloat(0.5)
		}
		return resultado
	default:
		return base
	}
}

// calcularConfiancaProjecao determina o nível de confiança da projeção
func calcularConfiancaProjecao(mesAhead int, tendencia string) string {
	// Confiança diminui com a distância
	if mesAhead == 1 {
		return "Alta"
	} else if mesAhead <= 3 {
		return "Média"
	}
	return "Baixa"
}

// calcularDiasUteisMes retorna uma estimativa de dias úteis
// TODO: Implementar cálculo real considerando feriados brasileiros
func calcularDiasUteisMes(ano, mes int) int {
	// Primeira abordagem: contar dias de semana
	inicio := time.Date(ano, time.Month(mes), 1, 0, 0, 0, 0, time.UTC)
	fim := inicio.AddDate(0, 1, -1) // Último dia do mês

	diasUteis := 0
	for d := inicio; !d.After(fim); d = d.AddDate(0, 0, 1) {
		weekday := d.Weekday()
		if weekday != time.Saturday && weekday != time.Sunday {
			diasUteis++
		}
	}

	// Subtrair uma estimativa de feriados (média brasileira: ~1 por mês)
	diasUteis -= 1
	if diasUteis < 15 {
		diasUteis = 15 // Mínimo razoável
	}

	return diasUteis
}
