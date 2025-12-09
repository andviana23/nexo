package postgres

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/google/uuid"
)

// MeioPagamentoRepository implementa operações de persistência para MeioPagamento
type MeioPagamentoRepository struct {
	queries *db.Queries
}

// NewMeioPagamentoRepository cria uma nova instância do repositório
func NewMeioPagamentoRepository(queries *db.Queries) *MeioPagamentoRepository {
	return &MeioPagamentoRepository{queries: queries}
}

// =============================================================================
// CREATE
// =============================================================================

// Create persiste um novo meio de pagamento
func (r *MeioPagamentoRepository) Create(ctx context.Context, meio *entity.MeioPagamento) error {
	dMais := int32(meio.DMais)
	ordemExibicao := int32(meio.OrdemExibicao)

	params := db.CreateMeioPagamentoParams{
		ID:            stringToUUID(meio.ID.String()),
		TenantID:      stringToUUID(meio.TenantID.String()),
		Nome:          meio.Nome,
		Tipo:          string(meio.Tipo),
		Bandeira:      strPtrToPgText(meio.Bandeira),
		Taxa:          decimalToNumeric(meio.Taxa),
		TaxaFixa:      decimalToNumeric(meio.TaxaFixa),
		DMais:         &dMais,
		Icone:         strPtrToPgText(meio.Icone),
		Cor:           strPtrToPgText(meio.Cor),
		OrdemExibicao: &ordemExibicao,
		Observacoes:   strPtrToPgText(meio.Observacoes),
		Ativo:         &meio.Ativo,
	}

	_, err := r.queries.CreateMeioPagamento(ctx, params)
	return err
}

// =============================================================================
// READ
// =============================================================================

// FindByID busca meio de pagamento por ID
func (r *MeioPagamentoRepository) FindByID(ctx context.Context, tenantID, id string) (*entity.MeioPagamento, error) {
	params := db.GetMeioPagamentoByIDParams{
		ID:       stringToUUID(id),
		TenantID: stringToUUID(tenantID),
	}

	row, err := r.queries.GetMeioPagamentoByID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("meio de pagamento não encontrado: %w", err)
	}

	return mapDBToMeioPagamento(row), nil
}

// List lista todos os meios de pagamento do tenant
func (r *MeioPagamentoRepository) List(ctx context.Context, tenantID string) ([]*entity.MeioPagamento, error) {
	rows, err := r.queries.ListMeiosPagamento(ctx, stringToUUID(tenantID))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar meios de pagamento: %w", err)
	}

	meios := make([]*entity.MeioPagamento, 0, len(rows))
	for _, row := range rows {
		meios = append(meios, mapDBToMeioPagamento(row))
	}

	return meios, nil
}

// ListAtivos lista apenas os meios de pagamento ativos
func (r *MeioPagamentoRepository) ListAtivos(ctx context.Context, tenantID string) ([]*entity.MeioPagamento, error) {
	rows, err := r.queries.ListMeiosPagamentoAtivos(ctx, stringToUUID(tenantID))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar meios de pagamento ativos: %w", err)
	}

	meios := make([]*entity.MeioPagamento, 0, len(rows))
	for _, row := range rows {
		meios = append(meios, mapDBToMeioPagamento(row))
	}

	return meios, nil
}

// ListByTipo lista meios de pagamento por tipo
func (r *MeioPagamentoRepository) ListByTipo(ctx context.Context, tenantID string, tipo entity.TipoPagamento) ([]*entity.MeioPagamento, error) {
	params := db.ListMeiosPagamentoPorTipoParams{
		TenantID: stringToUUID(tenantID),
		Tipo:     string(tipo),
	}

	rows, err := r.queries.ListMeiosPagamentoPorTipo(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar meios de pagamento por tipo: %w", err)
	}

	meios := make([]*entity.MeioPagamento, 0, len(rows))
	for _, row := range rows {
		meios = append(meios, mapDBToMeioPagamento(row))
	}

	return meios, nil
}

// Count retorna a quantidade total de meios de pagamento
func (r *MeioPagamentoRepository) Count(ctx context.Context, tenantID string) (int64, error) {
	return r.queries.CountMeiosPagamento(ctx, stringToUUID(tenantID))
}

// CountAtivos retorna a quantidade de meios de pagamento ativos
func (r *MeioPagamentoRepository) CountAtivos(ctx context.Context, tenantID string) (int64, error) {
	return r.queries.CountMeiosPagamentoAtivos(ctx, stringToUUID(tenantID))
}

// =============================================================================
// UPDATE
// =============================================================================

// Update atualiza um meio de pagamento existente
func (r *MeioPagamentoRepository) Update(ctx context.Context, meio *entity.MeioPagamento) error {
	dMais := int32(meio.DMais)
	ordemExibicao := int32(meio.OrdemExibicao)

	params := db.UpdateMeioPagamentoParams{
		ID:            stringToUUID(meio.ID.String()),
		TenantID:      stringToUUID(meio.TenantID.String()),
		Nome:          meio.Nome,
		Tipo:          string(meio.Tipo),
		Bandeira:      strPtrToPgText(meio.Bandeira),
		Taxa:          decimalToNumeric(meio.Taxa),
		TaxaFixa:      decimalToNumeric(meio.TaxaFixa),
		DMais:         &dMais,
		Icone:         strPtrToPgText(meio.Icone),
		Cor:           strPtrToPgText(meio.Cor),
		OrdemExibicao: &ordemExibicao,
		Observacoes:   strPtrToPgText(meio.Observacoes),
		Ativo:         &meio.Ativo,
	}

	_, err := r.queries.UpdateMeioPagamento(ctx, params)
	return err
}

// Toggle alterna o status ativo/inativo
func (r *MeioPagamentoRepository) Toggle(ctx context.Context, tenantID, id string) (*entity.MeioPagamento, error) {
	params := db.ToggleMeioPagamentoAtivoParams{
		ID:       stringToUUID(id),
		TenantID: stringToUUID(tenantID),
	}

	row, err := r.queries.ToggleMeioPagamentoAtivo(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao alternar status: %w", err)
	}

	return mapDBToMeioPagamento(row), nil
}

// =============================================================================
// DELETE
// =============================================================================

// Delete remove um meio de pagamento
func (r *MeioPagamentoRepository) Delete(ctx context.Context, tenantID, id string) error {
	params := db.DeleteMeioPagamentoParams{
		ID:       stringToUUID(id),
		TenantID: stringToUUID(tenantID),
	}

	return r.queries.DeleteMeioPagamento(ctx, params)
}

// ExistsByNome verifica se existe meio de pagamento com o nome
func (r *MeioPagamentoRepository) ExistsByNome(ctx context.Context, tenantID, nome string) (bool, error) {
	params := db.ExistsMeioPagamentoByNomeParams{
		TenantID: stringToUUID(tenantID),
		Lower:    nome,
	}

	return r.queries.ExistsMeioPagamentoByNome(ctx, params)
}

// =============================================================================
// MAPPER
// =============================================================================

func mapDBToMeioPagamento(row db.MeiosPagamento) *entity.MeioPagamento {
	id, _ := uuid.Parse(uuidToString(row.ID))
	tenantID, _ := uuid.Parse(uuidToString(row.TenantID))

	var icone, cor, observacoes, bandeira string
	if row.Icone != nil {
		icone = *row.Icone
	}
	if row.Cor != nil {
		cor = *row.Cor
	}
	if row.Observacoes != nil {
		observacoes = *row.Observacoes
	}
	if row.Bandeira != nil {
		bandeira = *row.Bandeira
	}

	var ordemExibicao, dMais int
	if row.OrdemExibicao != nil {
		ordemExibicao = int(*row.OrdemExibicao)
	}
	if row.DMais != nil {
		dMais = int(*row.DMais)
	}

	var ativo bool
	if row.Ativo != nil {
		ativo = *row.Ativo
	}

	return &entity.MeioPagamento{
		ID:            id,
		TenantID:      tenantID,
		Nome:          row.Nome,
		Tipo:          entity.TipoPagamento(row.Tipo),
		Bandeira:      bandeira,
		Taxa:          numericToDecimal(row.Taxa),
		TaxaFixa:      numericToDecimal(row.TaxaFixa),
		DMais:         dMais,
		Icone:         icone,
		Cor:           cor,
		OrdemExibicao: ordemExibicao,
		Observacoes:   observacoes,
		Ativo:         ativo,
		CriadoEm:      row.CriadoEm.Time,
		AtualizadoEm:  row.AtualizadoEm.Time,
	}
}
