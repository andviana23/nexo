package postgres

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

// ServicoRepository implementa port.ServicoRepository usando PostgreSQL
type ServicoRepository struct {
	queries *db.Queries
}

// NewServicoRepository cria uma nova instância
func NewServicoRepository(queries *db.Queries) *ServicoRepository {
	return &ServicoRepository{queries: queries}
}

// =============================================================================
// CREATE
// =============================================================================

// Create persiste um novo serviço
func (r *ServicoRepository) Create(ctx context.Context, servico *entity.Servico) error {
	params := db.CreateServicoParams{
		ID:               stringToUUID(servico.ID.String()),
		TenantID:         stringToUUID(servico.TenantID.String()),
		UnitID:           stringToUUIDNullable(servico.UnitID.String()),
		CategoriaID:      stringToUUIDNullable(servico.CategoriaID.String()),
		Nome:             servico.Nome,
		Descricao:        strPtrToPgText(servico.Descricao),
		Preco:            servico.Preco,
		Duracao:          int32(servico.Duracao),
		Comissao:         decimalToNumeric(servico.Comissao),
		Cor:              strPtrToPgText(servico.Cor),
		Imagem:           strPtrToPgText(servico.Imagem),
		ProfissionaisIds: uuidSliceToDBSlice(servico.ProfissionaisIDs),
		Observacoes:      strPtrToPgText(servico.Observacoes),
		Tags:             servico.Tags,
		Ativo:            &servico.Ativo,
	}

	_, err := r.queries.CreateServico(ctx, params)
	return err
}

// =============================================================================
// READ
// =============================================================================

// FindByID busca serviço por ID
func (r *ServicoRepository) FindByID(ctx context.Context, tenantID, unitID, id string) (*entity.Servico, error) {
	params := db.GetServicoByIDParams{
		ID:       stringToUUID(id),
		TenantID: stringToUUID(tenantID),
		UnitID:   stringToUUIDNullable(unitID),
	}

	row, err := r.queries.GetServicoByID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("serviço não encontrado: %w", err)
	}

	return mapDBRowToServico(row), nil
}

// List lista todos os serviços com filtros
func (r *ServicoRepository) List(ctx context.Context, tenantID, unitID string, filter port.ServicoFilter) ([]*entity.Servico, error) {
	tenantUUID := stringToUUID(tenantID)
	unitUUID := stringToUUIDNullable(unitID)

	// Seleção baseada nos filtros
	if filter.Search != "" {
		params := db.SearchServicosParams{
			TenantID:   tenantUUID,
			UnitID:     unitUUID,
			SearchTerm: filter.Search,
		}
		rows, err := r.queries.SearchServicos(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar serviços: %w", err)
		}
		return mapSearchRowsToServicos(rows), nil
	}

	if filter.CategoriaID != "" {
		params := db.ListServicosByCategoriaParams{
			TenantID:    tenantUUID,
			UnitID:      unitUUID,
			CategoriaID: stringToUUID(filter.CategoriaID),
		}
		rows, err := r.queries.ListServicosByCategoria(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("erro ao listar serviços por categoria: %w", err)
		}
		return mapCategoriaRowsToServicos(rows), nil
	}

	if filter.ProfissionalID != "" {
		params := db.ListServicosByProfissionalParams{
			TenantID:       tenantUUID,
			UnitID:         unitUUID,
			ProfissionalID: stringToUUID(filter.ProfissionalID),
		}
		rows, err := r.queries.ListServicosByProfissional(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("erro ao listar serviços por profissional: %w", err)
		}
		return mapProfissionalRowsToServicos(rows), nil
	}

	if filter.ApenasAtivos {
		params := db.ListServicosAtivosParams{
			TenantID: tenantUUID,
			UnitID:   unitUUID,
		}
		rows, err := r.queries.ListServicosAtivos(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("erro ao listar serviços ativos: %w", err)
		}
		return mapAtivosRowsToServicos(rows), nil
	}

	// Listagem padrão
	params := db.ListServicosParams{
		TenantID: tenantUUID,
		UnitID:   unitUUID,
	}
	rows, err := r.queries.ListServicos(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar serviços: %w", err)
	}
	return mapListRowsToServicos(rows), nil
}

// ListByCategoria lista serviços de uma categoria específica
func (r *ServicoRepository) ListByCategoria(ctx context.Context, tenantID, unitID, categoriaID string) ([]*entity.Servico, error) {
	params := db.ListServicosByCategoriaParams{
		TenantID:    stringToUUID(tenantID),
		UnitID:      stringToUUIDNullable(unitID),
		CategoriaID: stringToUUID(categoriaID),
	}

	rows, err := r.queries.ListServicosByCategoria(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar serviços por categoria: %w", err)
	}

	return mapCategoriaRowsToServicos(rows), nil
}

// ListByProfissional lista serviços que um profissional pode realizar
func (r *ServicoRepository) ListByProfissional(ctx context.Context, tenantID, unitID, profissionalID string) ([]*entity.Servico, error) {
	params := db.ListServicosByProfissionalParams{
		TenantID:       stringToUUID(tenantID),
		UnitID:         stringToUUIDNullable(unitID),
		ProfissionalID: stringToUUID(profissionalID),
	}

	rows, err := r.queries.ListServicosByProfissional(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar serviços por profissional: %w", err)
	}

	return mapProfissionalRowsToServicos(rows), nil
}

// FindByIDs busca múltiplos serviços por IDs
func (r *ServicoRepository) FindByIDs(ctx context.Context, tenantID, unitID string, ids []string) ([]*entity.Servico, error) {
	idsUUID := make([]pgtype.UUID, 0, len(ids))
	for _, id := range ids {
		idsUUID = append(idsUUID, stringToUUID(id))
	}

	params := db.GetServicosByIDsParams{
		TenantID:   stringToUUID(tenantID),
		ServicoIds: idsUUID,
	}

	rows, err := r.queries.GetServicosByIDs(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar serviços por IDs: %w", err)
	}

	return mapIDsRowsToServicos(rows), nil
}

// =============================================================================
// UPDATE
// =============================================================================

// Update atualiza um serviço existente
func (r *ServicoRepository) Update(ctx context.Context, servico *entity.Servico) error {
	params := db.UpdateServicoParams{
		ID:               stringToUUID(servico.ID.String()),
		TenantID:         stringToUUID(servico.TenantID.String()),
		CategoriaID:      stringToUUIDNullable(servico.CategoriaID.String()),
		Nome:             servico.Nome,
		Descricao:        strPtrToPgText(servico.Descricao),
		Preco:            servico.Preco,
		Duracao:          int32(servico.Duracao),
		Comissao:         decimalToNumeric(servico.Comissao),
		Cor:              strPtrToPgText(servico.Cor),
		Imagem:           strPtrToPgText(servico.Imagem),
		ProfissionaisIds: uuidSliceToDBSlice(servico.ProfissionaisIDs),
		Observacoes:      strPtrToPgText(servico.Observacoes),
		Tags:             servico.Tags,
		UnitID:           stringToUUIDNullable(servico.UnitID.String()),
	}

	_, err := r.queries.UpdateServico(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao atualizar serviço: %w", err)
	}

	return nil
}

// ToggleStatus ativa/desativa um serviço
func (r *ServicoRepository) ToggleStatus(ctx context.Context, tenantID, unitID, id string, ativo bool) error {
	params := db.ToggleServicoStatusParams{
		ID:       stringToUUID(id),
		TenantID: stringToUUID(tenantID),
		Ativo:    &ativo,
	}

	_, err := r.queries.ToggleServicoStatus(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao alterar status do serviço: %w", err)
	}

	return nil
}

// UpdateCategoria atualiza a categoria de um serviço
func (r *ServicoRepository) UpdateCategoria(ctx context.Context, tenantID, unitID, id, categoriaID string) error {
	params := db.UpdateServicoCategoriaParams{
		ID:          stringToUUID(id),
		TenantID:    stringToUUID(tenantID),
		CategoriaID: stringToUUID(categoriaID),
	}

	_, err := r.queries.UpdateServicoCategoria(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao atualizar categoria do serviço: %w", err)
	}

	return nil
}

// UpdateProfissionais atualiza a lista de profissionais de um serviço
func (r *ServicoRepository) UpdateProfissionais(ctx context.Context, tenantID, unitID, id string, profissionaisIDs []string) error {
	uuids := make([]uuid.UUID, 0, len(profissionaisIDs))
	for _, pid := range profissionaisIDs {
		if uid, err := uuid.Parse(pid); err == nil {
			uuids = append(uuids, uid)
		}
	}

	params := db.UpdateServicoProfissionaisParams{
		ID:               stringToUUID(id),
		TenantID:         stringToUUID(tenantID),
		ProfissionaisIds: uuidSliceToDBSlice(uuids),
	}

	_, err := r.queries.UpdateServicoProfissionais(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao atualizar profissionais do serviço: %w", err)
	}

	return nil
}

// =============================================================================
// DELETE
// =============================================================================

// Delete deleta um serviço
func (r *ServicoRepository) Delete(ctx context.Context, tenantID, unitID, id string) error {
	params := db.DeleteServicoParams{
		ID:       stringToUUID(id),
		TenantID: stringToUUID(tenantID),
	}

	err := r.queries.DeleteServico(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao deletar serviço: %w", err)
	}

	return nil
}

// DeleteByCategoria deleta todos os serviços de uma categoria
func (r *ServicoRepository) DeleteByCategoria(ctx context.Context, tenantID, unitID, categoriaID string) error {
	params := db.DeleteServicosByCategoriaParams{
		CategoriaID: stringToUUID(categoriaID),
		TenantID:    stringToUUID(tenantID),
	}

	err := r.queries.DeleteServicosByCategoria(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao deletar serviços por categoria: %w", err)
	}

	return nil
}

// =============================================================================
// QUERIES AUXILIARES
// =============================================================================

// CheckNomeExists verifica se já existe serviço com o mesmo nome
func (r *ServicoRepository) CheckNomeExists(ctx context.Context, tenantID, unitID, nome, excludeID string) (bool, error) {
	var idUUID pgtype.UUID

	if excludeID != "" {
		idUUID = stringToUUID(excludeID)
	} else {
		idUUID = pgtype.UUID{Valid: false}
	}

	params := db.CheckServicoNomeExistsParams{
		TenantID: stringToUUID(tenantID),
		UnitID:   stringToUUIDNullable(unitID),
		Lower:    nome,
		ID:       idUUID,
	}

	exists, err := r.queries.CheckServicoNomeExists(ctx, params)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar nome: %w", err)
	}

	return exists, nil
}

// GetStats retorna estatísticas dos serviços
func (r *ServicoRepository) GetStats(ctx context.Context, tenantID, unitID string) (*port.ServicoStats, error) {
	params := db.GetServicosStatsParams{
		TenantID: stringToUUID(tenantID),
		UnitID:   stringToUUIDNullable(unitID),
	}
	row, err := r.queries.GetServicosStats(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar estatísticas: %w", err)
	}

	return &port.ServicoStats{
		TotalServicos:    row.TotalServicos,
		ServicosAtivos:   row.ServicosAtivos,
		ServicosInativos: row.ServicosInativos,
		PrecoMedio:       interfaceToFloat64(row.PrecoMedio),
		DuracaoMedia:     interfaceToFloat64(row.DuracaoMedia),
		ComissaoMedia:    interfaceToFloat64(row.ComissaoMedia),
	}, nil
}

// Count conta total de serviços do tenant
func (r *ServicoRepository) Count(ctx context.Context, tenantID, unitID string) (int64, error) {
	params := db.CountServicosByTenantParams{
		TenantID: stringToUUID(tenantID),
		UnitID:   stringToUUIDNullable(unitID),
	}
	count, err := r.queries.CountServicosByTenant(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("erro ao contar serviços: %w", err)
	}
	return count, nil
}

// CountAtivos conta total de serviços ativos do tenant
func (r *ServicoRepository) CountAtivos(ctx context.Context, tenantID, unitID string) (int64, error) {
	params := db.CountServicosAtivosByTenantParams{
		TenantID: stringToUUID(tenantID),
		UnitID:   stringToUUIDNullable(unitID),
	}
	count, err := r.queries.CountServicosAtivosByTenant(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("erro ao contar serviços ativos: %w", err)
	}
	return count, nil
}

// =============================================================================
// HELPERS
// =============================================================================

// stringToUUIDNullable converte string para pgtype.UUID nullable
func stringToUUIDNullable(id string) pgtype.UUID {
	if id == "" {
		return pgtype.UUID{Valid: false}
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: uid, Valid: true}
}

// uuidSliceToDBSlice converte []uuid.UUID para []pgtype.UUID
func uuidSliceToDBSlice(uuids []uuid.UUID) []pgtype.UUID {
	result := make([]pgtype.UUID, 0, len(uuids))
	for _, uid := range uuids {
		result = append(result, pgtype.UUID{Bytes: uid, Valid: true})
	}
	return result
}

// dbSliceToUUIDSlice converte []pgtype.UUID para []uuid.UUID
func dbSliceToUUIDSlice(dbUUIDs []pgtype.UUID) []uuid.UUID {
	result := make([]uuid.UUID, 0, len(dbUUIDs))
	for _, dbUUID := range dbUUIDs {
		if dbUUID.Valid {
			result = append(result, uuid.UUID(dbUUID.Bytes))
		}
	}
	return result
}

// interfaceToFloat64 converte interface{} para float64
func interfaceToFloat64(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int64:
		return float64(val)
	case int32:
		return float64(val)
	case int:
		return float64(val)
	case decimal.Decimal:
		f, _ := val.Float64()
		return f
	case string:
		d, _ := decimal.NewFromString(val)
		f, _ := d.Float64()
		return f
	default:
		return 0
	}
}

// =============================================================================
// MAPPERS
// =============================================================================

// mapDBRowToServico converte GetServicoByIDRow para entity.Servico
func mapDBRowToServico(row db.GetServicoByIDRow) *entity.Servico {
	servico := &entity.Servico{
		ID:               uuid.UUID(row.ID.Bytes),
		TenantID:         uuid.UUID(row.TenantID.Bytes),
		Nome:             row.Nome,
		Preco:            row.Preco,
		Duracao:          int(row.Duracao),
		Comissao:         numericToDecimal(row.Comissao),
		ProfissionaisIDs: dbSliceToUUIDSlice(row.ProfissionaisIds),
		Tags:             row.Tags,
		Ativo:            row.Ativo != nil && *row.Ativo,
		CriadoEm:         row.CriadoEm.Time,
		AtualizadoEm:     row.AtualizadoEm.Time,
	}

	// Campos opcionais
	if row.CategoriaID.Valid {
		servico.CategoriaID = uuid.UUID(row.CategoriaID.Bytes)
	}
	if row.Descricao != nil {
		servico.Descricao = *row.Descricao
	}
	if row.Cor != nil {
		servico.Cor = *row.Cor
	}
	if row.Imagem != nil {
		servico.Imagem = *row.Imagem
	}
	if row.Observacoes != nil {
		servico.Observacoes = *row.Observacoes
	}
	if row.CategoriaNome != nil {
		servico.CategoriaNome = *row.CategoriaNome
	}
	if row.CategoriaCor != nil {
		servico.CategoriaCor = *row.CategoriaCor
	}

	return servico
}

// mapListRowsToServicos converte []ListServicosRow para []*entity.Servico
func mapListRowsToServicos(rows []db.ListServicosRow) []*entity.Servico {
	servicos := make([]*entity.Servico, 0, len(rows))
	for _, row := range rows {
		servico := &entity.Servico{
			ID:               uuid.UUID(row.ID.Bytes),
			TenantID:         uuid.UUID(row.TenantID.Bytes),
			Nome:             row.Nome,
			Preco:            row.Preco,
			Duracao:          int(row.Duracao),
			Comissao:         numericToDecimal(row.Comissao),
			ProfissionaisIDs: dbSliceToUUIDSlice(row.ProfissionaisIds),
			Tags:             row.Tags,
			Ativo:            row.Ativo != nil && *row.Ativo,
			CriadoEm:         row.CriadoEm.Time,
			AtualizadoEm:     row.AtualizadoEm.Time,
		}

		if row.CategoriaID.Valid {
			servico.CategoriaID = uuid.UUID(row.CategoriaID.Bytes)
		}
		if row.Descricao != nil {
			servico.Descricao = *row.Descricao
		}
		if row.Cor != nil {
			servico.Cor = *row.Cor
		}
		if row.Imagem != nil {
			servico.Imagem = *row.Imagem
		}
		if row.Observacoes != nil {
			servico.Observacoes = *row.Observacoes
		}
		if row.CategoriaNome != nil {
			servico.CategoriaNome = *row.CategoriaNome
		}
		if row.CategoriaCor != nil {
			servico.CategoriaCor = *row.CategoriaCor
		}

		servicos = append(servicos, servico)
	}
	return servicos
}

// mapAtivosRowsToServicos converte []ListServicosAtivosRow para []*entity.Servico
func mapAtivosRowsToServicos(rows []db.ListServicosAtivosRow) []*entity.Servico {
	servicos := make([]*entity.Servico, 0, len(rows))
	for _, row := range rows {
		servico := &entity.Servico{
			ID:               uuid.UUID(row.ID.Bytes),
			TenantID:         uuid.UUID(row.TenantID.Bytes),
			Nome:             row.Nome,
			Preco:            row.Preco,
			Duracao:          int(row.Duracao),
			Comissao:         numericToDecimal(row.Comissao),
			ProfissionaisIDs: dbSliceToUUIDSlice(row.ProfissionaisIds),
			Tags:             row.Tags,
			Ativo:            row.Ativo != nil && *row.Ativo,
			CriadoEm:         row.CriadoEm.Time,
			AtualizadoEm:     row.AtualizadoEm.Time,
		}

		if row.CategoriaID.Valid {
			servico.CategoriaID = uuid.UUID(row.CategoriaID.Bytes)
		}
		if row.Descricao != nil {
			servico.Descricao = *row.Descricao
		}
		if row.Cor != nil {
			servico.Cor = *row.Cor
		}
		if row.Imagem != nil {
			servico.Imagem = *row.Imagem
		}
		if row.Observacoes != nil {
			servico.Observacoes = *row.Observacoes
		}
		if row.CategoriaNome != nil {
			servico.CategoriaNome = *row.CategoriaNome
		}
		if row.CategoriaCor != nil {
			servico.CategoriaCor = *row.CategoriaCor
		}

		servicos = append(servicos, servico)
	}
	return servicos
}

// mapCategoriaRowsToServicos converte []ListServicosByCategoriaRow para []*entity.Servico
func mapCategoriaRowsToServicos(rows []db.ListServicosByCategoriaRow) []*entity.Servico {
	servicos := make([]*entity.Servico, 0, len(rows))
	for _, row := range rows {
		servico := &entity.Servico{
			ID:               uuid.UUID(row.ID.Bytes),
			TenantID:         uuid.UUID(row.TenantID.Bytes),
			Nome:             row.Nome,
			Preco:            row.Preco,
			Duracao:          int(row.Duracao),
			Comissao:         numericToDecimal(row.Comissao),
			ProfissionaisIDs: dbSliceToUUIDSlice(row.ProfissionaisIds),
			Tags:             row.Tags,
			Ativo:            row.Ativo != nil && *row.Ativo,
			CriadoEm:         row.CriadoEm.Time,
			AtualizadoEm:     row.AtualizadoEm.Time,
		}

		if row.CategoriaID.Valid {
			servico.CategoriaID = uuid.UUID(row.CategoriaID.Bytes)
		}
		if row.Descricao != nil {
			servico.Descricao = *row.Descricao
		}
		if row.Cor != nil {
			servico.Cor = *row.Cor
		}
		if row.Imagem != nil {
			servico.Imagem = *row.Imagem
		}
		if row.Observacoes != nil {
			servico.Observacoes = *row.Observacoes
		}
		if row.CategoriaNome != nil {
			servico.CategoriaNome = *row.CategoriaNome
		}
		if row.CategoriaCor != nil {
			servico.CategoriaCor = *row.CategoriaCor
		}

		servicos = append(servicos, servico)
	}
	return servicos
}

// mapProfissionalRowsToServicos converte []ListServicosByProfissionalRow para []*entity.Servico
func mapProfissionalRowsToServicos(rows []db.ListServicosByProfissionalRow) []*entity.Servico {
	servicos := make([]*entity.Servico, 0, len(rows))
	for _, row := range rows {
		servico := &entity.Servico{
			ID:               uuid.UUID(row.ID.Bytes),
			TenantID:         uuid.UUID(row.TenantID.Bytes),
			Nome:             row.Nome,
			Preco:            row.Preco,
			Duracao:          int(row.Duracao),
			Comissao:         numericToDecimal(row.Comissao),
			ProfissionaisIDs: dbSliceToUUIDSlice(row.ProfissionaisIds),
			Tags:             row.Tags,
			Ativo:            row.Ativo != nil && *row.Ativo,
			CriadoEm:         row.CriadoEm.Time,
			AtualizadoEm:     row.AtualizadoEm.Time,
		}

		if row.CategoriaID.Valid {
			servico.CategoriaID = uuid.UUID(row.CategoriaID.Bytes)
		}
		if row.Descricao != nil {
			servico.Descricao = *row.Descricao
		}
		if row.Cor != nil {
			servico.Cor = *row.Cor
		}
		if row.Imagem != nil {
			servico.Imagem = *row.Imagem
		}
		if row.Observacoes != nil {
			servico.Observacoes = *row.Observacoes
		}
		if row.CategoriaNome != nil {
			servico.CategoriaNome = *row.CategoriaNome
		}
		if row.CategoriaCor != nil {
			servico.CategoriaCor = *row.CategoriaCor
		}

		servicos = append(servicos, servico)
	}
	return servicos
}

// mapSearchRowsToServicos converte []SearchServicosRow para []*entity.Servico
func mapSearchRowsToServicos(rows []db.SearchServicosRow) []*entity.Servico {
	servicos := make([]*entity.Servico, 0, len(rows))
	for _, row := range rows {
		servico := &entity.Servico{
			ID:               uuid.UUID(row.ID.Bytes),
			TenantID:         uuid.UUID(row.TenantID.Bytes),
			Nome:             row.Nome,
			Preco:            row.Preco,
			Duracao:          int(row.Duracao),
			Comissao:         numericToDecimal(row.Comissao),
			ProfissionaisIDs: dbSliceToUUIDSlice(row.ProfissionaisIds),
			Tags:             row.Tags,
			Ativo:            row.Ativo != nil && *row.Ativo,
			CriadoEm:         row.CriadoEm.Time,
			AtualizadoEm:     row.AtualizadoEm.Time,
		}

		if row.CategoriaID.Valid {
			servico.CategoriaID = uuid.UUID(row.CategoriaID.Bytes)
		}
		if row.Descricao != nil {
			servico.Descricao = *row.Descricao
		}
		if row.Cor != nil {
			servico.Cor = *row.Cor
		}
		if row.Imagem != nil {
			servico.Imagem = *row.Imagem
		}
		if row.Observacoes != nil {
			servico.Observacoes = *row.Observacoes
		}
		if row.CategoriaNome != nil {
			servico.CategoriaNome = *row.CategoriaNome
		}
		if row.CategoriaCor != nil {
			servico.CategoriaCor = *row.CategoriaCor
		}

		servicos = append(servicos, servico)
	}
	return servicos
}

// mapIDsRowsToServicos converte []GetServicosByIDsRow para []*entity.Servico
func mapIDsRowsToServicos(rows []db.GetServicosByIDsRow) []*entity.Servico {
	servicos := make([]*entity.Servico, 0, len(rows))
	for _, row := range rows {
		servico := &entity.Servico{
			ID:               uuid.UUID(row.ID.Bytes),
			TenantID:         uuid.UUID(row.TenantID.Bytes),
			Nome:             row.Nome,
			Preco:            row.Preco,
			Duracao:          int(row.Duracao),
			Comissao:         numericToDecimal(row.Comissao),
			ProfissionaisIDs: dbSliceToUUIDSlice(row.ProfissionaisIds),
			Tags:             row.Tags,
			Ativo:            row.Ativo != nil && *row.Ativo,
			CriadoEm:         row.CriadoEm.Time,
			AtualizadoEm:     row.AtualizadoEm.Time,
		}

		if row.CategoriaID.Valid {
			servico.CategoriaID = uuid.UUID(row.CategoriaID.Bytes)
		}
		if row.Descricao != nil {
			servico.Descricao = *row.Descricao
		}
		if row.Cor != nil {
			servico.Cor = *row.Cor
		}
		if row.Imagem != nil {
			servico.Imagem = *row.Imagem
		}
		if row.Observacoes != nil {
			servico.Observacoes = *row.Observacoes
		}
		if row.CategoriaNome != nil {
			servico.CategoriaNome = *row.CategoriaNome
		}
		if row.CategoriaCor != nil {
			servico.CategoriaCor = *row.CategoriaCor
		}

		servicos = append(servicos, servico)
	}
	return servicos
}
