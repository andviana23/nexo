package command

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

// AddCommandItemUseCase implementa a adição de item à comanda
type AddCommandItemUseCase struct {
	repo        port.CommandRepository
	produtoRepo port.ProdutoRepository
	mapper      *mapper.CommandMapper
}

// NewAddCommandItemUseCase cria uma nova instância do use case
func NewAddCommandItemUseCase(
	repo port.CommandRepository,
	produtoRepo port.ProdutoRepository,
	mapper *mapper.CommandMapper,
) *AddCommandItemUseCase {
	return &AddCommandItemUseCase{
		repo:        repo,
		produtoRepo: produtoRepo,
		mapper:      mapper,
	}
}

// Execute adiciona um item à comanda e recalcula totais
// T-EST-001: Para produtos, valida disponibilidade de estoque antes de adicionar
func (uc *AddCommandItemUseCase) Execute(ctx context.Context, commandID, tenantID, userID uuid.UUID, req *dto.AddCommandItemRequest) (*dto.CommandResponse, error) {
	// Buscar comanda existente
	command, err := uc.repo.FindByID(ctx, commandID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get command: %w", err)
	}

	if command == nil {
		return nil, fmt.Errorf("command not found")
	}

	// T-EST-001: Validar estoque para itens do tipo PRODUTO
	if req.Tipo == string(entity.CommandItemTypeProduto) && uc.produtoRepo != nil {
		itemUUID, err := uuid.Parse(req.ItemID)
		if err != nil {
			return nil, fmt.Errorf("invalid item_id: %w", err)
		}

		produto, err := uc.produtoRepo.FindByID(ctx, tenantID, itemUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product: %w", err)
		}
		if produto == nil {
			return nil, fmt.Errorf("produto não encontrado")
		}

		// Verificar se produto está ativo
		if !produto.Ativo {
			return nil, fmt.Errorf("produto '%s' está inativo", produto.Nome)
		}

		// Verificar disponibilidade de estoque
		quantidadeSolicitada := decimal.NewFromInt(int64(req.Quantidade))
		if produto.QuantidadeAtual.LessThan(quantidadeSolicitada) {
			return nil, fmt.Errorf(
				"estoque insuficiente para '%s': disponível %.2f, solicitado %d",
				produto.Nome,
				produto.QuantidadeAtual.InexactFloat64(),
				req.Quantidade,
			)
		}
	}

	// Converter request para entity
	item, err := uc.mapper.FromAddCommandItemRequest(req, commandID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to map item: %w", err)
	}

	// Adicionar item via domain logic
	if err := command.AddItem(*item); err != nil {
		return nil, fmt.Errorf("failed to add item to command: %w", err)
	}

	// Persistir item
	if err := uc.repo.AddItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to persist item: %w", err)
	}

	// Atualizar comanda (totais recalculados)
	if err := uc.repo.Update(ctx, command); err != nil {
		return nil, fmt.Errorf("failed to update command: %w", err)
	}

	// Buscar comanda atualizada
	updated, err := uc.repo.FindByID(ctx, commandID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated command: %w", err)
	}

	// Converter para DTO
	response := uc.mapper.ToCommandResponse(updated)
	return response, nil
}
