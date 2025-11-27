package mapper

import (
	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
)

// =============================================================================
// Entity → Response DTO
// =============================================================================

// CategoriaServicoToResponse converte entidade para DTO de resposta
func CategoriaServicoToResponse(c *entity.CategoriaServico) *dto.CategoriaServicoResponse {
	if c == nil {
		return nil
	}

	return &dto.CategoriaServicoResponse{
		ID:           c.ID.String(),
		TenantID:     c.TenantID.String(),
		Nome:         c.Nome,
		Descricao:    c.Descricao,
		Cor:          c.Cor,
		Icone:        c.Icone,
		Ativa:        c.Ativa,
		CriadoEm:     c.CriadoEm.Format("2006-01-02T15:04:05Z07:00"),
		AtualizadoEm: c.AtualizadoEm.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// CategoriasServicosToResponse converte lista de entidades para lista de DTOs
func CategoriasServicosToResponse(categorias []*entity.CategoriaServico) []*dto.CategoriaServicoResponse {
	if categorias == nil {
		return []*dto.CategoriaServicoResponse{}
	}

	result := make([]*dto.CategoriaServicoResponse, 0, len(categorias))
	for _, c := range categorias {
		result = append(result, CategoriaServicoToResponse(c))
	}
	return result
}

// CategoriaServicoToResponseWithCount converte entidade para DTO com contagem
func CategoriaServicoToResponseWithCount(c *entity.CategoriaServico, totalServicos int64) *dto.CategoriaServicoWithCountResponse {
	if c == nil {
		return nil
	}

	return &dto.CategoriaServicoWithCountResponse{
		CategoriaServicoResponse: *CategoriaServicoToResponse(c),
		TotalServicos:            totalServicos,
	}
}
