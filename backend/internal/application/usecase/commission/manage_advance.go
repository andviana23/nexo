package commission

import (
	"context"


	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
)

// ApproveAdvanceInput representa a entrada para aprovar um adiantamento
type ApproveAdvanceInput struct {
	TenantID   string
	AdvanceID  string
	ApprovedBy string
}

// ApproveAdvanceOutput representa a saída da aprovação
type ApproveAdvanceOutput struct {
	Advance *entity.Advance
}

// ApproveAdvanceUseCase aprova um adiantamento
type ApproveAdvanceUseCase struct {
	advanceRepo repository.AdvanceRepository
}

// NewApproveAdvanceUseCase cria uma nova instância do use case
func NewApproveAdvanceUseCase(advanceRepo repository.AdvanceRepository) *ApproveAdvanceUseCase {
	return &ApproveAdvanceUseCase{
		advanceRepo: advanceRepo,
	}
}

// Execute executa o use case
func (uc *ApproveAdvanceUseCase) Execute(ctx context.Context, input ApproveAdvanceInput) (*ApproveAdvanceOutput, error) {
	// Verifica se existe
	advance, err := uc.advanceRepo.GetByID(ctx, input.TenantID, input.AdvanceID)
	if err != nil {
		return nil, err
	}

	if advance == nil {
		return nil, domain.ErrAdvanceNotFound
	}

	// Verifica se pode aprovar
	if !advance.CanApprove() {
		return nil, domain.ErrAdiantamentoNaoPodeAprovar
	}

	// Aprova
	approved, err := uc.advanceRepo.Approve(ctx, input.TenantID, input.AdvanceID, input.ApprovedBy)
	if err != nil {
		return nil, err
	}

	return &ApproveAdvanceOutput{Advance: approved}, nil
}

// RejectAdvanceInput representa a entrada para rejeitar um adiantamento
type RejectAdvanceInput struct {
	TenantID        string
	AdvanceID       string
	RejectedBy      string
	RejectionReason string
}

// RejectAdvanceOutput representa a saída da rejeição
type RejectAdvanceOutput struct {
	Advance *entity.Advance
}

// RejectAdvanceUseCase rejeita um adiantamento
type RejectAdvanceUseCase struct {
	advanceRepo repository.AdvanceRepository
}

// NewRejectAdvanceUseCase cria uma nova instância do use case
func NewRejectAdvanceUseCase(advanceRepo repository.AdvanceRepository) *RejectAdvanceUseCase {
	return &RejectAdvanceUseCase{
		advanceRepo: advanceRepo,
	}
}

// Execute executa o use case
func (uc *RejectAdvanceUseCase) Execute(ctx context.Context, input RejectAdvanceInput) (*RejectAdvanceOutput, error) {
	// Verifica se existe
	advance, err := uc.advanceRepo.GetByID(ctx, input.TenantID, input.AdvanceID)
	if err != nil {
		return nil, err
	}

	if advance == nil {
		return nil, domain.ErrAdvanceNotFound
	}

	// Verifica se pode rejeitar
	if !advance.CanReject() {
		return nil, domain.ErrAdiantamentoNaoPodeRejeitar
	}

	// Valida motivo
	if input.RejectionReason == "" {
		return nil, domain.ErrMotivoRejeicaoObrigatorio
	}

	// Rejeita
	rejected, err := uc.advanceRepo.Reject(ctx, input.TenantID, input.AdvanceID, input.RejectedBy, input.RejectionReason)
	if err != nil {
		return nil, err
	}

	return &RejectAdvanceOutput{Advance: rejected}, nil
}

// MarkAdvanceDeductedInput representa a entrada para marcar um adiantamento como deduzido
type MarkAdvanceDeductedInput struct {
	TenantID  string
	AdvanceID string
	PeriodID  string
}

// MarkAdvanceDeductedOutput representa a saída da marcação
type MarkAdvanceDeductedOutput struct {
	Advance *entity.Advance
}

// MarkAdvanceDeductedUseCase marca um adiantamento como deduzido
type MarkAdvanceDeductedUseCase struct {
	advanceRepo repository.AdvanceRepository
}

// NewMarkAdvanceDeductedUseCase cria uma nova instância do use case
func NewMarkAdvanceDeductedUseCase(advanceRepo repository.AdvanceRepository) *MarkAdvanceDeductedUseCase {
	return &MarkAdvanceDeductedUseCase{
		advanceRepo: advanceRepo,
	}
}

// Execute executa o use case
func (uc *MarkAdvanceDeductedUseCase) Execute(ctx context.Context, input MarkAdvanceDeductedInput) (*MarkAdvanceDeductedOutput, error) {
	// Verifica se existe
	advance, err := uc.advanceRepo.GetByID(ctx, input.TenantID, input.AdvanceID)
	if err != nil {
		return nil, err
	}

	if advance == nil {
		return nil, domain.ErrAdvanceNotFound
	}

	// Verifica se pode deduzir
	if !advance.CanDeduct() {
		return nil, domain.ErrAdiantamentoNaoPodeDeduzir
	}

	// Marca como deduzido
	deducted, err := uc.advanceRepo.MarkDeducted(ctx, input.TenantID, input.AdvanceID, input.PeriodID)
	if err != nil {
		return nil, err
	}

	return &MarkAdvanceDeductedOutput{Advance: deducted}, nil
}

// CancelAdvanceInput representa a entrada para cancelar um adiantamento
type CancelAdvanceInput struct {
	TenantID  string
	AdvanceID string
}

// CancelAdvanceOutput representa a saída do cancelamento
type CancelAdvanceOutput struct {
	Advance *entity.Advance
}

// CancelAdvanceUseCase cancela um adiantamento
type CancelAdvanceUseCase struct {
	advanceRepo repository.AdvanceRepository
}

// NewCancelAdvanceUseCase cria uma nova instância do use case
func NewCancelAdvanceUseCase(advanceRepo repository.AdvanceRepository) *CancelAdvanceUseCase {
	return &CancelAdvanceUseCase{
		advanceRepo: advanceRepo,
	}
}

// Execute executa o use case
func (uc *CancelAdvanceUseCase) Execute(ctx context.Context, input CancelAdvanceInput) (*CancelAdvanceOutput, error) {
	// Verifica se existe
	advance, err := uc.advanceRepo.GetByID(ctx, input.TenantID, input.AdvanceID)
	if err != nil {
		return nil, err
	}

	if advance == nil {
		return nil, domain.ErrAdvanceNotFound
	}

	// Verifica se pode cancelar
	if !advance.CanCancel() {
		return nil, domain.ErrAdiantamentoNaoPodeCancelar
	}

	// Cancela
	cancelled, err := uc.advanceRepo.Cancel(ctx, input.TenantID, input.AdvanceID)
	if err != nil {
		return nil, err
	}

	return &CancelAdvanceOutput{Advance: cancelled}, nil
}

// DeleteAdvanceInput representa a entrada para deletar um adiantamento
type DeleteAdvanceInput struct {
	TenantID  string
	AdvanceID string
}

// DeleteAdvanceOutput representa a saída da exclusão
type DeleteAdvanceOutput struct {
	Success bool
}

// DeleteAdvanceUseCase deleta um adiantamento
type DeleteAdvanceUseCase struct {
	advanceRepo repository.AdvanceRepository
}

// NewDeleteAdvanceUseCase cria uma nova instância do use case
func NewDeleteAdvanceUseCase(advanceRepo repository.AdvanceRepository) *DeleteAdvanceUseCase {
	return &DeleteAdvanceUseCase{
		advanceRepo: advanceRepo,
	}
}

// Execute executa o use case
func (uc *DeleteAdvanceUseCase) Execute(ctx context.Context, input DeleteAdvanceInput) (*DeleteAdvanceOutput, error) {
	// Verifica se existe
	advance, err := uc.advanceRepo.GetByID(ctx, input.TenantID, input.AdvanceID)
	if err != nil {
		return nil, err
	}

	if advance == nil {
		return nil, domain.ErrAdvanceNotFound
	}

	// Só pode deletar se estiver pendente
	if advance.Status != "PENDING" {
		return nil, domain.ErrAdiantamentoNaoPodeCancelar
	}

	// Deleta
	err = uc.advanceRepo.Delete(ctx, input.TenantID, input.AdvanceID)
	if err != nil {
		return nil, err
	}

	return &DeleteAdvanceOutput{Success: true}, nil
}
