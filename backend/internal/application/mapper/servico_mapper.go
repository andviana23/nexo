package mapper

import (
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// =============================================================================
// Entity → Response DTO
// =============================================================================

// ServicoToResponse converte entidade para DTO de resposta
func ServicoToResponse(s *entity.Servico) *dto.ServicoResponse {
	if s == nil {
		return nil
	}

	// Converter UUIDs de profissionais para strings
	profissionaisIDs := make([]string, 0, len(s.ProfissionaisIDs))
	for _, id := range s.ProfissionaisIDs {
		profissionaisIDs = append(profissionaisIDs, id.String())
	}

	// Converter preco para centavos
	precoCentavos := s.Preco.Mul(decimal.NewFromInt(100)).IntPart()

	resp := &dto.ServicoResponse{
		ID:               s.ID.String(),
		TenantID:         s.TenantID.String(),
		Nome:             s.Nome,
		Preco:            s.Preco.StringFixed(2),
		PrecoCentavos:    precoCentavos,
		Duracao:          s.Duracao,
		DuracaoFormatada: s.FormatarDuracao(),
		Comissao:         s.Comissao.StringFixed(2),
		ProfissionaisIDs: profissionaisIDs,
		Tags:             s.Tags,
		Ativo:            s.Ativo,
		CriadoEm:         s.CriadoEm.Format("2006-01-02T15:04:05Z07:00"),
		AtualizadoEm:     s.AtualizadoEm.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Campos opcionais
	if s.CategoriaID != uuid.Nil {
		categoriaID := s.CategoriaID.String()
		resp.CategoriaID = &categoriaID
	}
	if s.Descricao != "" {
		resp.Descricao = &s.Descricao
	}
	if s.Cor != "" {
		resp.Cor = &s.Cor
	}
	if s.Imagem != "" {
		resp.Imagem = &s.Imagem
	}
	if s.Observacoes != "" {
		resp.Observacoes = &s.Observacoes
	}
	if s.CategoriaNome != "" {
		resp.CategoriaNome = &s.CategoriaNome
	}
	if s.CategoriaCor != "" {
		resp.CategoriaCor = &s.CategoriaCor
	}

	return resp
}

// ServicosToResponse converte lista de entidades para lista de DTOs
func ServicosToResponse(servicos []*entity.Servico) []*dto.ServicoResponse {
	if servicos == nil {
		return []*dto.ServicoResponse{}
	}

	result := make([]*dto.ServicoResponse, 0, len(servicos))
	for _, s := range servicos {
		result = append(result, ServicoToResponse(s))
	}
	return result
}

// ServicoToSimplificadoResponse converte entidade para DTO simplificado
func ServicoToSimplificadoResponse(s *entity.Servico) *dto.ServicoSimplificadoResponse {
	if s == nil {
		return nil
	}

	resp := &dto.ServicoSimplificadoResponse{
		ID:      s.ID.String(),
		Nome:    s.Nome,
		Preco:   s.Preco.StringFixed(2),
		Duracao: s.Duracao,
		Ativo:   s.Ativo,
	}

	if s.CategoriaNome != "" {
		resp.CategoriaNome = &s.CategoriaNome
	}
	if s.Cor != "" {
		resp.Cor = &s.Cor
	}

	return resp
}

// ServicoStatsToResponse converte stats do repository para DTO
func ServicoStatsToResponse(stats *port.ServicoStats) *dto.ServicoStatsResponse {
	if stats == nil {
		return nil
	}

	return &dto.ServicoStatsResponse{
		TotalServicos:    stats.TotalServicos,
		ServicosAtivos:   stats.ServicosAtivos,
		ServicosInativos: stats.ServicosInativos,
		PrecoMedio:       fmt.Sprintf("%.2f", stats.PrecoMedio),
		DuracaoMedia:     stats.DuracaoMedia,
		ComissaoMedia:    fmt.Sprintf("%.2f", stats.ComissaoMedia),
	}
}

// =============================================================================
// Request DTO → Domain Entity
// =============================================================================

// CreateServicoRequestToEntity converte DTO de criação para entidade
func CreateServicoRequestToEntity(req *dto.CreateServicoRequest, tenantID uuid.UUID) (*entity.Servico, error) {
	// Parse preco
	preco, err := decimal.NewFromString(req.Preco)
	if err != nil {
		return nil, fmt.Errorf("preço inválido: %w", err)
	}

	// Parse unit_id
	var unitID uuid.UUID
	if req.UnitID != nil {
		unitID, err = uuid.Parse(*req.UnitID)
		if err != nil {
			return nil, fmt.Errorf("unit_id inválido: %w", err)
		}
	}

	// Criar entidade base
	servico, err := entity.NewServico(tenantID, unitID, req.Nome, preco, req.Duracao)
	if err != nil {
		return nil, err
	}

	// Campos opcionais
	if req.CategoriaID != nil {
		categoriaUUID, err := uuid.Parse(*req.CategoriaID)
		if err != nil {
			return nil, fmt.Errorf("categoria_id inválido: %w", err)
		}
		servico.SetCategoria(categoriaUUID)
	}

	if req.Descricao != nil {
		servico.SetDescricao(*req.Descricao)
	}

	if req.Comissao != nil {
		comissao, err := decimal.NewFromString(*req.Comissao)
		if err != nil {
			return nil, fmt.Errorf("comissão inválida: %w", err)
		}
		if err := servico.SetComissao(comissao); err != nil {
			return nil, err
		}
	}

	if req.Cor != nil {
		if err := servico.SetCor(*req.Cor); err != nil {
			return nil, err
		}
	}

	if req.Imagem != nil {
		servico.SetImagem(*req.Imagem)
	}

	if req.Observacoes != nil {
		servico.SetObservacoes(*req.Observacoes)
	}

	if len(req.Tags) > 0 {
		servico.SetTags(req.Tags)
	}

	if len(req.ProfissionaisIDs) > 0 {
		profissionaisUUIDs := make([]uuid.UUID, 0, len(req.ProfissionaisIDs))
		for _, id := range req.ProfissionaisIDs {
			uid, err := uuid.Parse(id)
			if err != nil {
				return nil, fmt.Errorf("profissional_id inválido (%s): %w", id, err)
			}
			profissionaisUUIDs = append(profissionaisUUIDs, uid)
		}
		servico.SetProfissionais(profissionaisUUIDs)
	}

	return servico, nil
}

// UpdateServicoRequestToEntity aplica dados de update em entidade existente
func UpdateServicoRequestToEntity(req *dto.UpdateServicoRequest, servico *entity.Servico) error {
	// Parse preco
	preco, err := decimal.NewFromString(req.Preco)
	if err != nil {
		return fmt.Errorf("preço inválido: %w", err)
	}

	// Update campos principais
	if err := servico.Update(req.Nome, preco, req.Duracao); err != nil {
		return err
	}

	// Categoria (apenas atualizar se for enviada)
	if req.CategoriaID != nil {
		categoriaUUID, err := uuid.Parse(*req.CategoriaID)
		if err != nil {
			return fmt.Errorf("categoria_id inválido: %w", err)
		}
		servico.SetCategoria(categoriaUUID)
	}
	// Se não enviada, mantém a categoria existente

	// Descrição (apenas atualizar se enviada)
	if req.Descricao != nil {
		servico.SetDescricao(*req.Descricao)
	}

	// Comissão (apenas atualizar se enviada)
	if req.Comissao != nil {
		comissao, err := decimal.NewFromString(*req.Comissao)
		if err != nil {
			return fmt.Errorf("comissão inválida: %w", err)
		}
		if err := servico.SetComissao(comissao); err != nil {
			return err
		}
	}

	// Cor (apenas atualizar se enviada)
	if req.Cor != nil {
		if err := servico.SetCor(*req.Cor); err != nil {
			return err
		}
	}

	// Imagem (apenas atualizar se enviada)
	if req.Imagem != nil {
		servico.SetImagem(*req.Imagem)
	}

	// Observações (apenas atualizar se enviadas)
	if req.Observacoes != nil {
		servico.SetObservacoes(*req.Observacoes)
	}

	// Tags (apenas atualizar se enviadas)
	if len(req.Tags) > 0 {
		servico.SetTags(req.Tags)
	}

	// Profissionais (apenas atualizar se enviados)
	if len(req.ProfissionaisIDs) > 0 {
		profissionaisUUIDs := make([]uuid.UUID, 0, len(req.ProfissionaisIDs))
		for _, id := range req.ProfissionaisIDs {
			uid, err := uuid.Parse(id)
			if err != nil {
				return fmt.Errorf("profissional_id inválido (%s): %w", id, err)
			}
			profissionaisUUIDs = append(profissionaisUUIDs, uid)
		}
		servico.SetProfissionais(profissionaisUUIDs)
	}

	return nil
}
