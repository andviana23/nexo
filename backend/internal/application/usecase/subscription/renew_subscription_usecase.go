package subscription

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
)

// RenewSubscriptionUseCase renova assinaturas manuais (PIX/Dinheiro)
type RenewSubscriptionUseCase struct {
	subRepo port.SubscriptionRepository
}

// NewRenewSubscriptionUseCase cria instância
func NewRenewSubscriptionUseCase(subRepo port.SubscriptionRepository) *RenewSubscriptionUseCase {
	return &RenewSubscriptionUseCase{subRepo: subRepo}
}

// Execute renova assinatura manual
func (uc *RenewSubscriptionUseCase) Execute(ctx context.Context, tenantID, subscriptionID string, req dto.RenewSubscriptionRequest) (*dto.SubscriptionResponse, error) {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, domain.ErrInvalidTenantID
	}
	subUUID, err := uuid.Parse(subscriptionID)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	sub, err := uc.subRepo.GetByID(ctx, subUUID, tenantUUID)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, domain.ErrSubscriptionNotFound
	}
	if sub.Status == entity.StatusCancelado {
		return nil, domain.ErrSubscriptionCannotReactivate
	}

	forma := entity.PaymentMethod(req.FormaPagamento)
	if forma != entity.PaymentMethodPix && forma != entity.PaymentMethodDinheiro {
		return nil, domain.ErrSubscriptionPaymentMethodInvalid
	}

	now := time.Now()
	venc := now.AddDate(0, 0, 30)

	// Atualizar status e datas
	if err := uc.subRepo.Activate(ctx, subUUID, tenantUUID, now, venc); err != nil {
		return nil, err
	}
	_ = uc.subRepo.ResetServicosUtilizados(ctx, subUUID, tenantUUID)

	// Registrar pagamento
	pagamento := &entity.SubscriptionPayment{
		ID:              uuid.New(),
		TenantID:        tenantUUID,
		SubscriptionID:  subUUID,
		Valor:           sub.Valor,
		FormaPagamento:  forma,
		Status:          entity.PaymentStatusConfirmado,
		DataPagamento:   &now,
		CodigoTransacao: req.CodigoTransacao,
	}
	_ = uc.subRepo.CreatePayment(ctx, pagamento)

	// Garantir flag de assinante
	_ = uc.subRepo.SetClienteAsSubscriber(ctx, sub.ClienteID, tenantUUID, true)

	// Atualizar struct para resposta
	sub.Status = entity.StatusAtivo
	sub.FormaPagamento = forma
	sub.DataAtivacao = &now
	sub.DataVencimento = &venc
	sub.ServicosUtilizados = 0

	return mapper.SubscriptionToResponse(sub), nil
}
