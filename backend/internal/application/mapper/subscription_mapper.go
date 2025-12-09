package mapper

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
)

// SubscriptionToResponse converte entidade para DTO de resposta
func SubscriptionToResponse(s *entity.Subscription) *dto.SubscriptionResponse {
	if s == nil {
		return nil
	}

	var dataAtivacao, dataVencimento *string
	if s.DataAtivacao != nil {
		v := s.DataAtivacao.Format(time.RFC3339)
		dataAtivacao = &v
	}
	if s.DataVencimento != nil {
		v := s.DataVencimento.Format(time.RFC3339)
		dataVencimento = &v
	}

	return &dto.SubscriptionResponse{
		ID:                 s.ID.String(),
		ClienteID:          s.ClienteID.String(),
		ClienteNome:        s.ClienteNome,
		ClienteTelefone:    s.ClienteTelefone,
		PlanoID:            s.PlanoID.String(),
		PlanoNome:          s.PlanoNome,
		FormaPagamento:     string(s.FormaPagamento),
		Status:             string(s.Status),
		Valor:              s.Valor.StringFixed(2),
		LinkPagamento:      s.LinkPagamento,
		DataAtivacao:       dataAtivacao,
		DataVencimento:     dataVencimento,
		ServicosUtilizados: s.ServicosUtilizados,
		CreatedAt:          s.CreatedAt.Format(time.RFC3339),
	}
}

// SubscriptionsToResponse converte lista de entidades
func SubscriptionsToResponse(list []*entity.Subscription) []*dto.SubscriptionResponse {
	if list == nil {
		return []*dto.SubscriptionResponse{}
	}
	out := make([]*dto.SubscriptionResponse, 0, len(list))
	for _, s := range list {
		out = append(out, SubscriptionToResponse(s))
	}
	return out
}

// SubscriptionMetricsToResponse converte métricas de domínio para DTO
func SubscriptionMetricsToResponse(m *entity.SubscriptionMetrics) *dto.SubscriptionMetricsResponse {
	if m == nil {
		return nil
	}
	receitaMensal, _ := m.ReceitaMensal.Float64()
	return &dto.SubscriptionMetricsResponse{
		TotalAssinantesAtivos:   m.TotalAtivas,
		TotalInativas:           m.TotalInativas,
		TotalInadimplentes:      m.TotalInadimplentes,
		TotalPlanosAtivos:       m.TotalPlanosAtivos,
		ReceitaMensal:           receitaMensal,
		TaxaRenovacao:           m.TaxaRenovacao,
		RenovacoesProximos7Dias: m.RenovacoesProximos7Dias,
	}
}
