// Package mapper contém funções para conversão entre entidades de domínio e DTOs HTTP.
package mapper

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
)

// ============================================================
// CAIXA DIÁRIO - Mappers
// ============================================================

// ToCaixaDiarioResponse converte entity.CaixaDiario para dto.CaixaDiarioResponse
func ToCaixaDiarioResponse(caixa *entity.CaixaDiario) dto.CaixaDiarioResponse {
	resp := dto.CaixaDiarioResponse{
		ID:                       caixa.ID.String(),
		UsuarioAberturaID:        caixa.UsuarioAberturaID.String(),
		UsuarioAberturaNome:      caixa.UsuarioAberturaNome,
		DataAbertura:             caixa.DataAbertura.Format(time.RFC3339),
		SaldoInicial:             caixa.SaldoInicial.String(),
		TotalEntradas:            caixa.TotalEntradas.String(),
		TotalSaidas:              caixa.TotalSaidas.String(),
		TotalSangrias:            caixa.TotalSangrias.String(),
		TotalReforcos:            caixa.TotalReforcos.String(),
		SaldoEsperado:            caixa.SaldoEsperado.String(),
		Status:                   string(caixa.Status),
		JustificativaDivergencia: caixa.JustificativaDivergencia,
		CreatedAt:                caixa.CreatedAt.Format(time.RFC3339),
		UpdatedAt:                caixa.UpdatedAt.Format(time.RFC3339),
		UsuarioFechamentoNome:    caixa.UsuarioFechamentoNome,
	}

	// Campos opcionais
	if caixa.UsuarioFechamentoID != nil {
		ufID := caixa.UsuarioFechamentoID.String()
		resp.UsuarioFechamentoID = &ufID
	}

	if caixa.DataFechamento != nil {
		df := caixa.DataFechamento.Format(time.RFC3339)
		resp.DataFechamento = &df
	}

	if caixa.SaldoReal != nil {
		sr := caixa.SaldoReal.String()
		resp.SaldoReal = &sr
	}

	if caixa.Divergencia != nil {
		div := caixa.Divergencia.String()
		resp.Divergencia = &div
	}

	// Mapear operações se existirem
	if len(caixa.Operacoes) > 0 {
		resp.Operacoes = make([]dto.OperacaoCaixaResponse, len(caixa.Operacoes))
		for i, op := range caixa.Operacoes {
			resp.Operacoes[i] = ToOperacaoCaixaResponse(&op)
		}
	}

	return resp
}

// ToCaixaDiarioResumoResponse converte entity.CaixaDiario para dto.CaixaDiarioResumoResponse
func ToCaixaDiarioResumoResponse(caixa *entity.CaixaDiario) dto.CaixaDiarioResumoResponse {
	resp := dto.CaixaDiarioResumoResponse{
		ID:                    caixa.ID.String(),
		UsuarioAberturaNome:   caixa.UsuarioAberturaNome,
		UsuarioFechamentoNome: caixa.UsuarioFechamentoNome,
		DataAbertura:          caixa.DataAbertura.Format(time.RFC3339),
		SaldoInicial:          caixa.SaldoInicial.String(),
		SaldoEsperado:         caixa.SaldoEsperado.String(),
		Status:                string(caixa.Status),
		TemDivergencia:        caixa.TemDivergencia(),
	}

	if caixa.DataFechamento != nil {
		df := caixa.DataFechamento.Format(time.RFC3339)
		resp.DataFechamento = &df
	}

	if caixa.SaldoReal != nil {
		sr := caixa.SaldoReal.String()
		resp.SaldoReal = &sr
	}

	if caixa.Divergencia != nil {
		div := caixa.Divergencia.String()
		resp.Divergencia = &div
	}

	return resp
}

// ToOperacaoCaixaResponse converte entity.OperacaoCaixa para dto.OperacaoCaixaResponse
func ToOperacaoCaixaResponse(op *entity.OperacaoCaixa) dto.OperacaoCaixaResponse {
	return dto.OperacaoCaixaResponse{
		ID:          op.ID.String(),
		Tipo:        string(op.Tipo),
		Valor:       op.Valor.String(),
		Descricao:   op.Descricao,
		Destino:     op.Destino,
		Origem:      op.Origem,
		UsuarioID:   op.UsuarioID.String(),
		UsuarioNome: op.UsuarioNome,
		CreatedAt:   op.CreatedAt.Format(time.RFC3339),
	}
}

// ToCaixaStatusResponse cria resposta de status do caixa
func ToCaixaStatusResponse(caixa *entity.CaixaDiario, ultimoFechamento *time.Time) dto.CaixaStatusResponse {
	resp := dto.CaixaStatusResponse{
		Aberto: caixa != nil && caixa.Status == entity.StatusCaixaAberto,
	}

	if caixa != nil && caixa.Status == entity.StatusCaixaAberto {
		caixaResp := ToCaixaDiarioResponse(caixa)
		resp.CaixaAtual = &caixaResp
	}

	if ultimoFechamento != nil {
		uf := ultimoFechamento.Format(time.RFC3339)
		resp.UltimoFechamento = &uf
	}

	return resp
}

// ToListCaixaHistoricoResponse converte lista de caixas para resposta paginada
func ToListCaixaHistoricoResponse(caixas []*entity.CaixaDiario, total int64, page, pageSize int) dto.ListCaixaHistoricoResponse {
	items := make([]dto.CaixaDiarioResumoResponse, len(caixas))
	for i, c := range caixas {
		items[i] = ToCaixaDiarioResumoResponse(c)
	}

	// Proteção contra divisão por zero
	totalPages := 0
	if pageSize > 0 {
		totalPages = int(total) / pageSize
		if int(total)%pageSize > 0 {
			totalPages++
		}
	}

	return dto.ListCaixaHistoricoResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
