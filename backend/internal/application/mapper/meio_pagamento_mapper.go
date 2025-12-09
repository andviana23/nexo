package mapper

import (
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
)

// MeioPagamentoToResponse converte entidade para DTO de resposta
func MeioPagamentoToResponse(meio *entity.MeioPagamento) dto.MeioPagamentoResponse {
	return dto.MeioPagamentoResponse{
		ID:            meio.ID.String(),
		TenantID:      meio.TenantID.String(),
		Nome:          meio.Nome,
		Tipo:          string(meio.Tipo),
		TipoLabel:     meio.Tipo.DisplayName(),
		Bandeira:      meio.Bandeira,
		Taxa:          meio.Taxa.StringFixed(2),
		TaxaFixa:      meio.TaxaFixa.StringFixed(2),
		DMais:         meio.DMais,
		DMaisLabel:    fmt.Sprintf("D+%d", meio.DMais),
		Icone:         meio.Icone,
		Cor:           meio.Cor,
		OrdemExibicao: meio.OrdemExibicao,
		Observacoes:   meio.Observacoes,
		Ativo:         meio.Ativo,
		CriadoEm:      meio.CriadoEm,
		AtualizadoEm:  meio.AtualizadoEm,
	}
}

// MeiosPagamentoToListResponse converte lista de entidades para response
func MeiosPagamentoToListResponse(meios []*entity.MeioPagamento, total, totalAtivo int64) dto.ListMeiosPagamentoResponse {
	data := make([]dto.MeioPagamentoResponse, 0, len(meios))
	for _, meio := range meios {
		data = append(data, MeioPagamentoToResponse(meio))
	}

	return dto.ListMeiosPagamentoResponse{
		Data:       data,
		Total:      total,
		TotalAtivo: totalAtivo,
	}
}
