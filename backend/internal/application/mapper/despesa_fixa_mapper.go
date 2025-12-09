package mapper

import (
	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/financial"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/shopspring/decimal"
)

// DespesaFixaMapper converte entre DTOs e entidades de despesa fixa
type DespesaFixaMapper struct{}

// NewDespesaFixaMapper cria nova instância do mapper
func NewDespesaFixaMapper() *DespesaFixaMapper {
	return &DespesaFixaMapper{}
}

// ToCreateInput converte CreateDespesaFixaRequest para CreateDespesaFixaInput
func (m *DespesaFixaMapper) ToCreateInput(req dto.CreateDespesaFixaRequest, tenantID string) (*financial.CreateDespesaFixaInput, error) {
	valor, err := decimal.NewFromString(req.Valor)
	if err != nil {
		return nil, err
	}

	return &financial.CreateDespesaFixaInput{
		TenantID:      tenantID,
		Descricao:     req.Descricao,
		CategoriaID:   req.CategoriaID,
		Fornecedor:    req.Fornecedor,
		Valor:         valueobject.NewMoneyFromDecimal(valor),
		DiaVencimento: req.DiaVencimento,
		UnidadeID:     req.UnidadeID,
		Observacoes:   req.Observacoes,
	}, nil
}

// ToUpdateInput converte UpdateDespesaFixaRequest para UpdateDespesaFixaInput
func (m *DespesaFixaMapper) ToUpdateInput(req dto.UpdateDespesaFixaRequest, tenantID, id string) (*financial.UpdateDespesaFixaInput, error) {
	valor, err := decimal.NewFromString(req.Valor)
	if err != nil {
		return nil, err
	}

	return &financial.UpdateDespesaFixaInput{
		TenantID:      tenantID,
		ID:            id,
		Descricao:     req.Descricao,
		CategoriaID:   req.CategoriaID,
		Fornecedor:    req.Fornecedor,
		Valor:         valueobject.NewMoneyFromDecimal(valor),
		DiaVencimento: req.DiaVencimento,
		UnidadeID:     req.UnidadeID,
		Observacoes:   req.Observacoes,
	}, nil
}

// ToResponse converte entity.DespesaFixa para DespesaFixaResponse
func (m *DespesaFixaMapper) ToResponse(despesa *entity.DespesaFixa) dto.DespesaFixaResponse {
	return dto.DespesaFixaResponse{
		ID:            despesa.ID,
		UnidadeID:     despesa.UnidadeID,
		Descricao:     despesa.Descricao,
		CategoriaID:   despesa.CategoriaID,
		Fornecedor:    despesa.Fornecedor,
		Valor:         despesa.Valor.Raw(),
		DiaVencimento: despesa.DiaVencimento,
		Ativo:         despesa.Ativo,
		Observacoes:   despesa.Observacoes,
		CriadoEm:      despesa.CriadoEm,
		AtualizadoEm:  despesa.AtualizadoEm,
	}
}

// ToListResponse converte slice de entities para DespesasFixasListResponse
func (m *DespesaFixaMapper) ToListResponse(despesas []*entity.DespesaFixa, total int64, page, pageSize int) dto.DespesasFixasListResponse {
	data := make([]dto.DespesaFixaResponse, 0, len(despesas))
	for _, d := range despesas {
		data = append(data, m.ToResponse(d))
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return dto.DespesasFixasListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// ToGerarContasResponse converte output do use case para DTO
func (m *DespesaFixaMapper) ToGerarContasResponse(output *financial.GerarContasFromDespesasFixasOutput) dto.GerarContasResponse {
	return dto.GerarContasResponse{
		TotalDespesas:   output.TotalDespesas,
		ContasCriadas:   output.ContasCriadas,
		Erros:           output.Erros,
		DetalhesErros:   output.DetalhesErros,
		TempoExecucaoMs: output.TempoExecucaoMs,
	}
}

// ToSummaryResponse cria resposta de resumo
func (m *DespesaFixaMapper) ToSummaryResponse(total, totalAtivas int64, valorTotal valueobject.Money) dto.DespesasFixasSummaryResponse {
	return dto.DespesasFixasSummaryResponse{
		Total:       total,
		TotalAtivas: totalAtivas,
		ValorTotal:  valorTotal.Raw(),
	}
}
