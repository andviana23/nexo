package mapper

import (
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PlanToResponse converte entidade para DTO de resposta
func PlanToResponse(p *entity.Plan) *dto.PlanResponse {
	if p == nil {
		return nil
	}

	return &dto.PlanResponse{
		ID:              p.ID.String(),
		Nome:            p.Nome,
		Descricao:       p.Descricao,
		Valor:           p.Valor.StringFixed(2),
		Periodicidade:   p.Periodicidade,
		QtdServicos:     p.QtdServicos,
		LimiteUsoMensal: p.LimiteUsoMensal,
		Ativo:           p.Ativo,
		CreatedAt:       p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       p.UpdatedAt.Format(time.RFC3339),
	}
}

// PlansToResponse converte slice de entidades para slice de DTOs
func PlansToResponse(plans []*entity.Plan) []*dto.PlanResponse {
	if plans == nil {
		return []*dto.PlanResponse{}
	}
	out := make([]*dto.PlanResponse, 0, len(plans))
	for _, p := range plans {
		out = append(out, PlanToResponse(p))
	}
	return out
}

// CreatePlanRequestToEntity cria entidade a partir do payload
func CreatePlanRequestToEntity(req *dto.CreatePlanRequest, tenantID uuid.UUID) (*entity.Plan, error) {
	valor, err := decimal.NewFromString(req.Valor)
	if err != nil {
		return nil, fmt.Errorf("valor inválido: %w", err)
	}

	plan, err := entity.NewPlan(tenantID, req.Nome, valor)
	if err != nil {
		return nil, err
	}

	plan.Descricao = req.Descricao
	plan.QtdServicos = req.QtdServicos
	plan.LimiteUsoMensal = req.LimiteUsoMensal
	return plan, nil
}

// UpdatePlanRequestToEntity aplica payload em entidade existente
func UpdatePlanRequestToEntity(plan *entity.Plan, req *dto.UpdatePlanRequest) error {
	valor, err := decimal.NewFromString(req.Valor)
	if err != nil {
		return fmt.Errorf("valor inválido: %w", err)
	}
	return plan.UpdateBasicFields(req.Nome, req.Descricao, valor, req.QtdServicos, req.LimiteUsoMensal, req.Ativo)
}
