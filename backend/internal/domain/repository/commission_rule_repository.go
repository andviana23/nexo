package repository

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
)

// CommissionRuleRepository define as operações de repositório para regras de comissão
type CommissionRuleRepository interface {
	// Create cria uma nova regra de comissão
	Create(ctx context.Context, rule *entity.CommissionRule) (*entity.CommissionRule, error)

	// GetByID busca uma regra de comissão por ID
	GetByID(ctx context.Context, tenantID, id string) (*entity.CommissionRule, error)

	// List lista todas as regras de comissão de um tenant
	List(ctx context.Context, tenantID string) ([]*entity.CommissionRule, error)

	// ListActive lista regras de comissão ativas de um tenant
	ListActive(ctx context.Context, tenantID string) ([]*entity.CommissionRule, error)

	// GetEffective busca regra de comissão vigente em uma data
	GetEffective(ctx context.Context, tenantID string, date time.Time) ([]*entity.CommissionRule, error)

	// GetEffectiveByUnit busca regra vigente específica de uma unidade
	GetEffectiveByUnit(ctx context.Context, tenantID, unitID string, date time.Time) (*entity.CommissionRule, error)

	// GetEffectiveGlobal busca regra vigente global do tenant (sem unidade)
	GetEffectiveGlobal(ctx context.Context, tenantID string, date time.Time) (*entity.CommissionRule, error)

	// Update atualiza uma regra de comissão
	Update(ctx context.Context, rule *entity.CommissionRule) (*entity.CommissionRule, error)

	// Delete remove uma regra de comissão
	Delete(ctx context.Context, tenantID, id string) error

	// Deactivate desativa uma regra de comissão
	Deactivate(ctx context.Context, tenantID, id string) error
}
