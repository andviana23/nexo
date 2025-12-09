// Package mapper contém funções para conversão entre entidades de domínio e DTOs HTTP.
// Responsável por transformar dados entre camadas da aplicação seguindo Clean Architecture.
package mapper

import (
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/shopspring/decimal"
)

// ToContaPagarResponse converte entidade ContaPagar para DTO Response
func ToContaPagarResponse(conta *entity.ContaPagar) dto.ContaPagarResponse {
	var dataPagamento *string
	if conta.DataPagamento != nil {
		dp := conta.DataPagamento.Format("2006-01-02")
		dataPagamento = &dp
	}

	return dto.ContaPagarResponse{
		ID:             conta.ID,
		Descricao:      conta.Descricao,
		CategoriaID:    conta.CategoriaID,
		Fornecedor:     conta.Fornecedor,
		Valor:          conta.Valor.Raw(),
		Tipo:           string(conta.Tipo),
		Recorrente:     conta.Recorrente,
		Periodicidade:  conta.Periodicidade,
		DataVencimento: conta.DataVencimento.Format("2006-01-02"),
		DataPagamento:  dataPagamento,
		Status:         string(conta.Status),
		ComprovanteURL: conta.ComprovanteURL,
		PixCode:        conta.PixCode,
		Observacoes:    conta.Observacoes,
		CriadoEm:       conta.CriadoEm.Format(time.RFC3339),
		AtualizadoEm:   conta.AtualizadoEm.Format(time.RFC3339),
	}
}

// ToContaReceberResponse converte entidade ContaReceber para DTO Response
func ToContaReceberResponse(conta *entity.ContaReceber) dto.ContaReceberResponse {
	var dataRecebimento *string
	if conta.DataRecebimento != nil {
		dr := conta.DataRecebimento.Format("2006-01-02")
		dataRecebimento = &dr
	}

	return dto.ContaReceberResponse{
		ID:              conta.ID,
		Origem:          conta.Origem,
		AssinaturaID:    conta.AssinaturaID,
		DescricaoOrigem: conta.DescricaoOrigem,
		Valor:           conta.Valor.Raw(),
		ValorPago:       conta.ValorPago.Raw(),
		ValorAberto:     conta.ValorAberto.Raw(),
		DataVencimento:  conta.DataVencimento.Format("2006-01-02"),
		DataRecebimento: dataRecebimento,
		Status:          string(conta.Status),
		Observacoes:     conta.Observacoes,
		CriadoEm:        conta.CriadoEm.Format(time.RFC3339),
		AtualizadoEm:    conta.AtualizadoEm.Format(time.RFC3339),
	}
}

// ToFluxoCaixaDiarioResponse converte entidade para DTO Response
func ToFluxoCaixaDiarioResponse(fluxo *entity.FluxoCaixaDiario) dto.FluxoCaixaDiarioResponse {
	return dto.FluxoCaixaDiarioResponse{
		ID:                  fluxo.ID,
		Data:                fluxo.Data.Format("2006-01-02"),
		SaldoInicial:        fluxo.SaldoInicial.Raw(),
		EntradasConfirmadas: fluxo.EntradasConfirmadas.Raw(),
		EntradasPrevistas:   fluxo.EntradasPrevistas.Raw(),
		SaidasPagas:         fluxo.SaidasPagas.Raw(),
		SaidasPrevistas:     fluxo.SaidasPrevistas.Raw(),
		SaldoFinal:          fluxo.SaldoFinal.Raw(),
		ProcessadoEm:        fluxo.ProcessadoEm.Format(time.RFC3339),
	}
}

// ToCompensacaoBancariaResponse converte entidade para DTO Response
func ToCompensacaoBancariaResponse(comp *entity.CompensacaoBancaria) dto.CompensacaoBancariaResponse {
	var dataCompensado *string
	if comp.DataCompensado != nil {
		dc := comp.DataCompensado.Format("2006-01-02")
		dataCompensado = &dc
	}

	return dto.CompensacaoBancariaResponse{
		ID:              comp.ID,
		ReceitaID:       comp.ReceitaID,
		DataTransacao:   comp.DataTransacao.Format("2006-01-02"),
		DataCompensacao: comp.DataCompensacao.Format("2006-01-02"),
		DataCompensado:  dataCompensado,
		ValorBruto:      comp.ValorBruto.Raw(),
		TaxaPercentual:  comp.TaxaPercentual.String(),
		TaxaFixa:        comp.TaxaFixa.Raw(),
		ValorLiquido:    comp.ValorLiquido.Raw(),
		MeioPagamentoID: comp.MeioPagamentoID,
		DMais:           comp.DMais.Dias(),
		Status:          string(comp.Status),
	}
}

// ToDREMensalResponse converte entidade DRE para DTO Response
func ToDREMensalResponse(dre *entity.DREMensal) dto.DREMensalResponse {
	return dto.DREMensalResponse{
		ID:                   dre.ID,
		MesAno:               dre.MesAno.String(),
		ReceitaServicos:      dre.ReceitaServicos.Raw(),
		ReceitaProdutos:      dre.ReceitaProdutos.Raw(),
		ReceitaPlanos:        dre.ReceitaPlanos.Raw(),
		ReceitaTotal:         dre.ReceitaTotal.Raw(),
		CustoComissoes:       dre.CustoComissoes.Raw(),
		CustoInsumos:         dre.CustoInsumos.Raw(),
		CustoVariavelTotal:   dre.CustoVariavelTotal.Raw(),
		DespesaFixa:          dre.DespesaFixa.Raw(),
		DespesaVariavel:      dre.DespesaVariavel.Raw(),
		DespesaTotal:         dre.DespesaTotal.Raw(),
		ResultadoBruto:       dre.ResultadoBruto.Raw(),
		ResultadoOperacional: dre.ResultadoOperacional.Raw(),
		MargemBruta:          dre.MargemBruta.String(),
		MargemOperacional:    dre.MargemOperacional.String(),
		LucroLiquido:         dre.LucroLiquido.Raw(),
		ProcessadoEm:         dre.ProcessadoEm.Format(time.RFC3339),
	}
}

// FromCreateContaPagarRequest converte DTO Request para parâmetros do use case
func FromCreateContaPagarRequest(req dto.CreateContaPagarRequest) (
	valor valueobject.Money,
	tipo valueobject.TipoCusto,
	dataVencimento time.Time,
	err error,
) {
	// Parse valor
	valorDecimal, err := decimal.NewFromString(req.Valor)
	if err != nil {
		return valueobject.Money{}, valueobject.TipoCusto(""), time.Time{}, fmt.Errorf("valor inválido: %w", err)
	}
	valor = valueobject.NewMoneyFromDecimal(valorDecimal)

	// Parse tipo
	tipo = valueobject.TipoCusto(req.Tipo)
	if !tipo.IsValid() {
		return valueobject.Money{}, valueobject.TipoCusto(""), time.Time{}, fmt.Errorf("tipo inválido")
	}

	// Parse data
	dataVencimento, err = time.Parse("2006-01-02", req.DataVencimento)
	if err != nil {
		return valueobject.Money{}, valueobject.TipoCusto(""), time.Time{}, fmt.Errorf("data de vencimento inválida: %w", err)
	}

	return valor, tipo, dataVencimento, nil
}

// FromCreateContaReceberRequest converte DTO Request para parâmetros do use case
func FromCreateContaReceberRequest(req dto.CreateContaReceberRequest) (
	valor valueobject.Money,
	dataVencimento time.Time,
	err error,
) {
	// Parse valor
	valorDecimal, err := decimal.NewFromString(req.Valor)
	if err != nil {
		return valueobject.Money{}, time.Time{}, fmt.Errorf("valor inválido: %w", err)
	}
	valor = valueobject.NewMoneyFromDecimal(valorDecimal)

	// Parse data
	dataVencimento, err = time.Parse("2006-01-02", req.DataVencimento)
	if err != nil {
		return valueobject.Money{}, time.Time{}, fmt.Errorf("data de vencimento inválida: %w", err)
	}

	return valor, dataVencimento, nil
}
