package meiopagamento

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CreateMeioPagamentoUseCase implementa a criação de meio de pagamento
type CreateMeioPagamentoUseCase struct {
	repo port.MeioPagamentoRepository
}

// NewCreateMeioPagamentoUseCase cria uma nova instância
func NewCreateMeioPagamentoUseCase(repo port.MeioPagamentoRepository) *CreateMeioPagamentoUseCase {
	return &CreateMeioPagamentoUseCase{repo: repo}
}

// Execute cria um novo meio de pagamento
func (uc *CreateMeioPagamentoUseCase) Execute(ctx context.Context, tenantID string, req dto.CreateMeioPagamentoRequest) (*dto.MeioPagamentoResponse, error) {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant_id inválido: %w", err)
	}

	tipo := entity.TipoPagamento(req.Tipo)
	if !tipo.IsValid() {
		return nil, entity.ErrMeioPagamentoTipoInvalido
	}

	meio, err := entity.NewMeioPagamento(tenantUUID, req.Nome, tipo)
	if err != nil {
		return nil, err
	}

	// Set optional fields
	meio.Bandeira = req.Bandeira
	meio.Icone = req.Icone
	meio.Cor = req.Cor
	meio.OrdemExibicao = req.OrdemExibicao
	meio.Observacoes = req.Observacoes

	if req.Ativo != nil {
		meio.Ativo = *req.Ativo
	}

	// Parse and set taxa
	if req.Taxa != "" {
		taxa, err := decimal.NewFromString(req.Taxa)
		if err == nil {
			if err := meio.SetTaxa(taxa); err != nil {
				return nil, err
			}
		}
	}

	// Parse and set taxa fixa
	if req.TaxaFixa != "" {
		taxaFixa, err := decimal.NewFromString(req.TaxaFixa)
		if err == nil {
			if err := meio.SetTaxaFixa(taxaFixa); err != nil {
				return nil, err
			}
		}
	}

	// Set D+
	if req.DMais > 0 {
		if err := meio.SetDMais(req.DMais); err != nil {
			return nil, err
		}
	}

	if err := uc.repo.Create(ctx, meio); err != nil {
		return nil, fmt.Errorf("erro ao criar meio de pagamento: %w", err)
	}

	response := mapper.MeioPagamentoToResponse(meio)
	return &response, nil
}

// ListMeiosPagamentoUseCase implementa a listagem de meios de pagamento
type ListMeiosPagamentoUseCase struct {
	repo port.MeioPagamentoRepository
}

// NewListMeiosPagamentoUseCase cria uma nova instância
func NewListMeiosPagamentoUseCase(repo port.MeioPagamentoRepository) *ListMeiosPagamentoUseCase {
	return &ListMeiosPagamentoUseCase{repo: repo}
}

// Execute lista meios de pagamento com filtros
func (uc *ListMeiosPagamentoUseCase) Execute(ctx context.Context, tenantID string, filter dto.MeioPagamentoFilter) (*dto.ListMeiosPagamentoResponse, error) {
	var meios []*entity.MeioPagamento
	var err error

	if filter.ApenasAtivos {
		meios, err = uc.repo.ListAtivos(ctx, tenantID)
	} else if filter.Tipo != "" {
		tipo := entity.TipoPagamento(filter.Tipo)
		meios, err = uc.repo.ListByTipo(ctx, tenantID, tipo)
	} else {
		meios, err = uc.repo.List(ctx, tenantID)
	}

	if err != nil {
		return nil, err
	}

	// Filter by search if provided
	if filter.Search != "" {
		filtered := make([]*entity.MeioPagamento, 0)
		for _, m := range meios {
			if containsIgnoreCase(m.Nome, filter.Search) {
				filtered = append(filtered, m)
			}
		}
		meios = filtered
	}

	total, _ := uc.repo.Count(ctx, tenantID)
	totalAtivo, _ := uc.repo.CountAtivos(ctx, tenantID)

	response := mapper.MeiosPagamentoToListResponse(meios, total, totalAtivo)
	return &response, nil
}

// GetMeioPagamentoUseCase implementa a busca por ID
type GetMeioPagamentoUseCase struct {
	repo port.MeioPagamentoRepository
}

// NewGetMeioPagamentoUseCase cria uma nova instância
func NewGetMeioPagamentoUseCase(repo port.MeioPagamentoRepository) *GetMeioPagamentoUseCase {
	return &GetMeioPagamentoUseCase{repo: repo}
}

// Execute busca um meio de pagamento por ID
func (uc *GetMeioPagamentoUseCase) Execute(ctx context.Context, tenantID, id string) (*dto.MeioPagamentoResponse, error) {
	meio, err := uc.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	response := mapper.MeioPagamentoToResponse(meio)
	return &response, nil
}

// UpdateMeioPagamentoUseCase implementa a atualização
type UpdateMeioPagamentoUseCase struct {
	repo port.MeioPagamentoRepository
}

// NewUpdateMeioPagamentoUseCase cria uma nova instância
func NewUpdateMeioPagamentoUseCase(repo port.MeioPagamentoRepository) *UpdateMeioPagamentoUseCase {
	return &UpdateMeioPagamentoUseCase{repo: repo}
}

// Execute atualiza um meio de pagamento
func (uc *UpdateMeioPagamentoUseCase) Execute(ctx context.Context, tenantID, id string, req dto.UpdateMeioPagamentoRequest) (*dto.MeioPagamentoResponse, error) {
	meio, err := uc.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Nome != "" {
		meio.Nome = req.Nome
	}
	if req.Tipo != "" {
		tipo := entity.TipoPagamento(req.Tipo)
		if !tipo.IsValid() {
			return nil, entity.ErrMeioPagamentoTipoInvalido
		}
		meio.Tipo = tipo
	}
	if req.Bandeira != "" {
		meio.Bandeira = req.Bandeira
	}
	if req.Icone != "" {
		meio.Icone = req.Icone
	}
	if req.Cor != "" {
		meio.Cor = req.Cor
	}
	if req.Observacoes != "" {
		meio.Observacoes = req.Observacoes
	}
	if req.OrdemExibicao != nil {
		meio.OrdemExibicao = *req.OrdemExibicao
	}
	if req.DMais != nil {
		if err := meio.SetDMais(*req.DMais); err != nil {
			return nil, err
		}
	}
	if req.Ativo != nil {
		meio.Ativo = *req.Ativo
	}

	if req.Taxa != "" {
		taxa, err := decimal.NewFromString(req.Taxa)
		if err == nil {
			if err := meio.SetTaxa(taxa); err != nil {
				return nil, err
			}
		}
	}

	if req.TaxaFixa != "" {
		taxaFixa, err := decimal.NewFromString(req.TaxaFixa)
		if err == nil {
			if err := meio.SetTaxaFixa(taxaFixa); err != nil {
				return nil, err
			}
		}
	}

	if err := uc.repo.Update(ctx, meio); err != nil {
		return nil, fmt.Errorf("erro ao atualizar meio de pagamento: %w", err)
	}

	response := mapper.MeioPagamentoToResponse(meio)
	return &response, nil
}

// ToggleMeioPagamentoUseCase implementa toggle de status
type ToggleMeioPagamentoUseCase struct {
	repo port.MeioPagamentoRepository
}

// NewToggleMeioPagamentoUseCase cria uma nova instância
func NewToggleMeioPagamentoUseCase(repo port.MeioPagamentoRepository) *ToggleMeioPagamentoUseCase {
	return &ToggleMeioPagamentoUseCase{repo: repo}
}

// Execute alterna o status ativo/inativo
func (uc *ToggleMeioPagamentoUseCase) Execute(ctx context.Context, tenantID, id string) (*dto.MeioPagamentoResponse, error) {
	meio, err := uc.repo.Toggle(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	response := mapper.MeioPagamentoToResponse(meio)
	return &response, nil
}

// DeleteMeioPagamentoUseCase implementa exclusão
type DeleteMeioPagamentoUseCase struct {
	repo port.MeioPagamentoRepository
}

// NewDeleteMeioPagamentoUseCase cria uma nova instância
func NewDeleteMeioPagamentoUseCase(repo port.MeioPagamentoRepository) *DeleteMeioPagamentoUseCase {
	return &DeleteMeioPagamentoUseCase{repo: repo}
}

// Execute exclui um meio de pagamento
func (uc *DeleteMeioPagamentoUseCase) Execute(ctx context.Context, tenantID, id string) error {
	return uc.repo.Delete(ctx, tenantID, id)
}

// Helper function
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 &&
			(s[0] == substr[0] || s[0]+32 == substr[0] || s[0] == substr[0]+32) &&
			containsIgnoreCase(s[1:], substr[1:])))
}
