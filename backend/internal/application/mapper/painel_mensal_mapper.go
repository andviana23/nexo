package mapper

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/financial"
)

// ToPainelMensalResponse converte o output do use case para DTO de resposta
func ToPainelMensalResponse(output *financial.PainelMensalOutput) *dto.PainelMensalResponse {
	if output == nil {
		return nil
	}

	return &dto.PainelMensalResponse{
		Ano:                 output.Ano,
		Mes:                 output.Mes,
		NomeMes:             output.NomeMes,
		ReceitaRealizada:    output.ReceitaRealizada.Value().StringFixed(2),
		ReceitaPendente:     output.ReceitaPendente.Value().StringFixed(2),
		ReceitaTotal:        output.ReceitaTotal.Value().StringFixed(2),
		DespesasFixas:       output.DespesasFixas.Value().StringFixed(2),
		DespesasVariaveis:   output.DespesasVariaveis.Value().StringFixed(2),
		DespesasPagas:       output.DespesasPagas.Value().StringFixed(2),
		DespesasPendentes:   output.DespesasPendentes.Value().StringFixed(2),
		DespesasTotal:       output.DespesasTotal.Value().StringFixed(2),
		LucroBruto:          output.LucroBruto.Value().StringFixed(2),
		LucroLiquido:        output.LucroLiquido.Value().StringFixed(2),
		MargemLiquida:       output.MargemLiquida.StringFixed(2),
		MetaMensal:          output.MetaMensal.Value().StringFixed(2),
		PercentualMeta:      output.PercentualMeta.StringFixed(2),
		DiferencaMeta:       output.DiferencaMeta.Value().StringFixed(2),
		StatusMeta:          output.StatusMeta,
		SaldoCaixaAtual:     output.SaldoCaixaAtual.Value().StringFixed(2),
		VariacaoMesAnterior: output.VariacaoMesAnterior.StringFixed(2),
		TendenciaVariacao:   output.TendenciaVariacao,
	}
}

// ToProjecoesResponse converte o output do use case para DTO de resposta
func ToProjecoesResponse(output *financial.ProjecoesOutput) *dto.ProjecoesResponse {
	if output == nil {
		return nil
	}

	projecoes := make([]dto.ProjecaoMensalResponse, 0, len(output.Projecoes))
	for _, p := range output.Projecoes {
		projecoes = append(projecoes, dto.ProjecaoMensalResponse{
			Ano:                p.Ano,
			Mes:                p.Mes,
			NomeMes:            p.NomeMes,
			ReceitaProjetada:   p.ReceitaProjetada.Value().StringFixed(2),
			DespesasProjetadas: p.DespesasProjetadas.Value().StringFixed(2),
			DespesasFixas:      p.DespesasFixas.Value().StringFixed(2),
			LucroProjetado:     p.LucroProjetado.Value().StringFixed(2),
			DiasUteis:          p.DiasUteis,
			MetaDiaria:         p.MetaDiaria.Value().StringFixed(2),
			Confianca:          p.Confianca,
		})
	}

	return &dto.ProjecoesResponse{
		Projecoes:           projecoes,
		MediaReceita3Meses:  output.MediaReceita3Meses.Value().StringFixed(2),
		MediaDespesas3Meses: output.MediaDespesas3Meses.Value().StringFixed(2),
		TendenciaReceita:    output.TendenciaReceita,
		DataGeracao:         output.DataGeracao.Format(time.RFC3339),
	}
}
