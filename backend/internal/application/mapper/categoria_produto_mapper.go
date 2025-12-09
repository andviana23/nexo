package mapper

import (
	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
)

// CategoriaProdutoToResponse converte entity para response DTO
func CategoriaProdutoToResponse(e *entity.CategoriaProdutoEntity) dto.CategoriaProdutoResponse {
	return dto.CategoriaProdutoResponse{
		ID:           e.ID.String(),
		Nome:         e.Nome,
		Descricao:    e.Descricao,
		Cor:          e.Cor,
		Icone:        e.Icone,
		CentroCusto:  string(e.CentroCusto),
		Ativa:        e.Ativa,
		CriadoEm:     e.CriadoEm,
		AtualizadoEm: e.AtualizadoEm,
	}
}

// CategoriasProdutosToListResponse converte lista de entities para response DTO
func CategoriasProdutosToListResponse(entities []*entity.CategoriaProdutoEntity) dto.ListCategoriaProdutoResponse {
	categorias := make([]dto.CategoriaProdutoResponse, len(entities))
	for i, e := range entities {
		categorias[i] = CategoriaProdutoToResponse(e)
	}
	return dto.ListCategoriaProdutoResponse{
		Categorias: categorias,
		Total:      len(categorias),
	}
}
