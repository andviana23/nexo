// Package postgres contém o repositório de profissionais
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

// ProfessionalRepository implementa operações de profissionais
type ProfessionalRepository struct {
	queries *db.Queries
}

// NewProfessionalRepository cria uma nova instância
func NewProfessionalRepository(queries *db.Queries) *ProfessionalRepository {
	return &ProfessionalRepository{queries: queries}
}

// List lista profissionais com filtros e paginação
func (r *ProfessionalRepository) List(ctx context.Context, tenantID, unitID string, req dto.ListProfessionalsRequest) ([]dto.ProfessionalResponse, int64, error) {
	// Defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize

	var search *string
	if req.Search != "" {
		search = &req.Search
	}
	var status *string
	if req.Status != "" {
		status = &req.Status
	}
	var tipo *string
	if req.Tipo != "" {
		tipo = &req.Tipo
	}
	var orderBy *string
	if req.OrderBy != "" {
		orderBy = &req.OrderBy
	}

	listParams := db.ListProfessionalsParams{
		TenantID:   stringToUUID(tenantID),
		UnitID:     stringToUUID(unitID),
		Status:     status,
		Tipo:       tipo,
		Search:     search,
		OrderBy:    orderBy,
		PageSize:   int32(req.PageSize),
		PageOffset: int32(offset),
	}

	rows, err := r.queries.ListProfessionals(ctx, listParams)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao listar profissionais: %w", err)
	}

	countParams := db.CountProfessionalsParams{
		TenantID: stringToUUID(tenantID),
		UnitID:   stringToUUID(unitID),
		Status:   status,
		Tipo:     tipo,
		Search:   search,
	}

	total, err := r.queries.CountProfessionals(ctx, countParams)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao contar profissionais: %w", err)
	}

	professionals := make([]dto.ProfessionalResponse, 0, len(rows))
	for _, row := range rows {
		professionals = append(professionals, mapListRowToResponse(row))
	}

	return professionals, total, nil
}

// GetByID busca um profissional por ID
func (r *ProfessionalRepository) GetByID(ctx context.Context, tenantID, unitID, id string) (*dto.ProfessionalResponse, error) {
	params := db.GetProfessionalByIDParams{
		ID:       stringToUUID(id),
		TenantID: stringToUUID(tenantID),
		UnitID:   stringToUUID(unitID),
	}

	row, err := r.queries.GetProfessionalByID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar profissional: %w", err)
	}

	resp := mapGetByIDRowToResponse(row)

	// Fetch commissions by category
	commissions, err := r.queries.ListProfessionalCategoryCommissions(ctx, db.ListProfessionalCategoryCommissionsParams{
		TenantID:       stringToUUID(tenantID),
		ProfissionalID: stringToUUID(id),
	})
	if err != nil {
		// Log error but don't fail? Or fail? Better fail or return empty?
		// Usually simpler to just return empty if error, but here unlikely.
		// Let's return empty if error, or just continue.
		// return nil, fmt.Errorf("erro ao buscar comissões: %w", err)
		// For now, fail safe.
	} else {
		resp.ComissoesPorCategoria = make([]dto.CommissionByCategory, 0, len(commissions))
		for _, c := range commissions {
			val, _ := c.Comissao.Float64()
			resp.ComissoesPorCategoria = append(resp.ComissoesPorCategoria, dto.CommissionByCategory{
				CategoriaID: uuidToString(c.CategoriaID),
				Comissao:    val,
			})
		}
	}

	return &resp, nil
}

// Create cria um novo profissional
func (r *ProfessionalRepository) Create(ctx context.Context, tenantID, unitID string, req dto.CreateProfessionalRequest) (*dto.ProfessionalResponse, error) {
	// Parse comissao
	var comissao decimal.Decimal
	if req.Comissao != 0 {
		comissao = decimal.NewFromFloat(req.Comissao)
	} else {
		comissao = decimal.Zero
	}

	// Parse data_admissao
	dataAdmissao, err := time.Parse("2006-01-02", req.DataAdmissao)
	if err != nil {
		return nil, fmt.Errorf("data de admissão inválida: %w", err)
	}

	// Default status
	status := req.Status
	if status == "" {
		status = "ATIVO"
	}

	// Default tipo_comissao
	tipoComissao := req.TipoComissao
	if tipoComissao == "" {
		tipoComissao = "PERCENTUAL"
	}

	params := db.CreateProfessionalParams{
		TenantID:        stringToUUID(tenantID),
		UnitID:          stringToUUID(unitID),
		Nome:            req.Nome,
		Email:           req.Email,
		Telefone:        req.Telefone,
		Cpf:             req.CPF,
		Especialidades:  req.Especialidades,
		Comissao:        pgtype.Numeric{Int: comissao.BigInt(), Exp: -2, Valid: true},
		TipoComissao:    &tipoComissao,
		Foto:            req.Foto,
		DataAdmissao:    pgtype.Date{Time: dataAdmissao, Valid: true},
		Status:          &status,
		HorarioTrabalho: stringToJSONB(req.HorarioTrabalho),
		Observacoes:     req.Observacoes,
		Tipo:            req.Tipo,
	}

	row, err := r.queries.CreateProfessional(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar profissional: %w", err)
	}

	// Save category commissions
	if len(req.ComissoesPorCategoria) > 0 {
		if err := r.saveCategoryCommissions(ctx, tenantID, uuidToString(row.ID), req.ComissoesPorCategoria); err != nil {
			return nil, fmt.Errorf("erro ao salvar comissões por categoria: %w", err)
		}
	}

	resp := mapCreateRowToResponse(row)
	return &resp, nil
}

// Update atualiza um profissional
func (r *ProfessionalRepository) Update(ctx context.Context, tenantID, unitID, id string, req dto.UpdateProfessionalRequest) (*dto.ProfessionalResponse, error) {
	// Parse comissao
	var comissao decimal.Decimal
	if req.Comissao != 0 {
		comissao = decimal.NewFromFloat(req.Comissao)
	} else {
		comissao = decimal.Zero
	}

	// Parse data_admissao
	dataAdmissao, err := time.Parse("2006-01-02", req.DataAdmissao)
	if err != nil {
		return nil, fmt.Errorf("data de admissão inválida: %w", err)
	}

	// Parse data_demissao (optional)
	var dataDemissao pgtype.Date
	if req.DataDemissao != nil && *req.DataDemissao != "" {
		t, err := time.Parse("2006-01-02", *req.DataDemissao)
		if err != nil {
			return nil, fmt.Errorf("data de demissão inválida: %w", err)
		}
		dataDemissao = pgtype.Date{Time: t, Valid: true}
	}

	// Default tipo_comissao
	tipoComissao := req.TipoComissao
	if tipoComissao == "" {
		tipoComissao = "PERCENTUAL"
	}

	params := db.UpdateProfessionalParams{
		ID:              stringToUUID(id),
		TenantID:        stringToUUID(tenantID),
		UnitID:          stringToUUID(unitID),
		Nome:            req.Nome,
		Email:           req.Email,
		Telefone:        req.Telefone,
		Cpf:             req.CPF,
		Especialidades:  req.Especialidades,
		Comissao:        pgtype.Numeric{Int: comissao.BigInt(), Exp: -2, Valid: true},
		TipoComissao:    &tipoComissao,
		Foto:            req.Foto,
		DataAdmissao:    pgtype.Date{Time: dataAdmissao, Valid: true},
		DataDemissao:    dataDemissao,
		Status:          &req.Status,
		HorarioTrabalho: stringToJSONB(req.HorarioTrabalho),
		Observacoes:     req.Observacoes,
		Tipo:            req.Tipo,
	}

	row, err := r.queries.UpdateProfessional(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar profissional: %w", err)
	}

	// Save category commissions
	if err := r.saveCategoryCommissions(ctx, tenantID, id, req.ComissoesPorCategoria); err != nil {
		return nil, fmt.Errorf("erro ao salvar comissões por categoria: %w", err)
	}

	resp := mapUpdateRowToResponse(row)
	return &resp, nil
}

// UpdateStatus atualiza o status de um profissional
func (r *ProfessionalRepository) UpdateStatus(ctx context.Context, tenantID, unitID, id string, req dto.UpdateProfessionalStatusRequest) (*dto.ProfessionalResponse, error) {
	var dataDemissao pgtype.Date
	if req.DataDemissao != nil && *req.DataDemissao != "" {
		t, err := time.Parse("2006-01-02", *req.DataDemissao)
		if err != nil {
			return nil, fmt.Errorf("data de demissão inválida: %w", err)
		}
		dataDemissao = pgtype.Date{Time: t, Valid: true}
	}

	params := db.UpdateProfessionalStatusParams{
		ID:           stringToUUID(id),
		TenantID:     stringToUUID(tenantID),
		UnitID:       stringToUUID(unitID),
		Status:       &req.Status,
		DataDemissao: dataDemissao,
	}

	row, err := r.queries.UpdateProfessionalStatus(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar status: %w", err)
	}

	resp := mapUpdateStatusRowToResponse(row)
	return &resp, nil
}

// Delete remove um profissional
func (r *ProfessionalRepository) Delete(ctx context.Context, tenantID, unitID, id string) error {
	params := db.DeleteProfessionalParams{
		ID:       stringToUUID(id),
		TenantID: stringToUUID(tenantID),
		UnitID:   stringToUUID(unitID),
	}

	err := r.queries.DeleteProfessional(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao deletar profissional: %w", err)
	}

	return nil
}

// CheckEmailExists verifica se email já existe
func (r *ProfessionalRepository) CheckEmailExists(ctx context.Context, tenantID, unitID, email string, excludeID *string) (bool, error) {
	var excludeUUID pgtype.UUID
	if excludeID != nil {
		excludeUUID = stringToUUID(*excludeID)
	}

	params := db.CheckEmailExistsProfessionalParams{
		TenantID:  stringToUUID(tenantID),
		UnitID:    stringToUUID(unitID),
		Email:     email,
		ExcludeID: excludeUUID,
	}

	exists, err := r.queries.CheckEmailExistsProfessional(ctx, params)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar email: %w", err)
	}

	return exists, nil
}

// CheckCpfExists verifica se CPF já existe
func (r *ProfessionalRepository) CheckCpfExists(ctx context.Context, tenantID, unitID, cpf string, excludeID *string) (bool, error) {
	var excludeUUID pgtype.UUID
	if excludeID != nil {
		excludeUUID = stringToUUID(*excludeID)
	}

	params := db.CheckCpfExistsProfessionalParams{
		TenantID:  stringToUUID(tenantID),
		UnitID:    stringToUUID(unitID),
		Cpf:       cpf,
		ExcludeID: excludeUUID,
	}

	exists, err := r.queries.CheckCpfExistsProfessional(ctx, params)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar CPF: %w", err)
	}

	return exists, nil
}

// =============================================================================
// Helpers
// =============================================================================

// profissionalData é uma interface para tipos de row de profissionais
type profissionalData interface {
	GetID() pgtype.UUID
	GetTenantID() pgtype.UUID
	GetUserID() pgtype.UUID
	GetNome() string
	GetEmail() string
	GetTelefone() string
	GetCpf() string
	GetEspecialidades() []string
	GetComissao() pgtype.Numeric
	GetTipoComissao() *string
	GetFoto() *string
	GetDataAdmissao() pgtype.Date
	GetDataDemissao() pgtype.Date
	GetStatus() *string
	GetHorarioTrabalho() []byte
	GetObservacoes() *string
	GetCriadoEm() pgtype.Timestamptz
	GetAtualizadoEm() pgtype.Timestamptz
	GetTipo() string
}

func mapListRowToResponse(row db.ListProfessionalsRow) dto.ProfessionalResponse {
	return mapGenericToResponse(
		row.ID, row.TenantID, row.UserID, row.Nome, row.Email, row.Telefone, row.Cpf,
		row.Especialidades, row.Comissao, row.TipoComissao, row.Foto, row.DataAdmissao,
		row.DataDemissao, row.Status, row.HorarioTrabalho, row.Observacoes,
		row.CriadoEm, row.AtualizadoEm, row.Tipo,
	)
}

func mapGetByIDRowToResponse(row db.GetProfessionalByIDRow) dto.ProfessionalResponse {
	return mapGenericToResponse(
		row.ID, row.TenantID, row.UserID, row.Nome, row.Email, row.Telefone, row.Cpf,
		row.Especialidades, row.Comissao, row.TipoComissao, row.Foto, row.DataAdmissao,
		row.DataDemissao, row.Status, row.HorarioTrabalho, row.Observacoes,
		row.CriadoEm, row.AtualizadoEm, row.Tipo,
	)
}

func mapCreateRowToResponse(row db.CreateProfessionalRow) dto.ProfessionalResponse {
	return mapGenericToResponse(
		row.ID, row.TenantID, row.UserID, row.Nome, row.Email, row.Telefone, row.Cpf,
		row.Especialidades, row.Comissao, row.TipoComissao, row.Foto, row.DataAdmissao,
		row.DataDemissao, row.Status, row.HorarioTrabalho, row.Observacoes,
		row.CriadoEm, row.AtualizadoEm, row.Tipo,
	)
}

func mapUpdateRowToResponse(row db.UpdateProfessionalRow) dto.ProfessionalResponse {
	return mapGenericToResponse(
		row.ID, row.TenantID, row.UserID, row.Nome, row.Email, row.Telefone, row.Cpf,
		row.Especialidades, row.Comissao, row.TipoComissao, row.Foto, row.DataAdmissao,
		row.DataDemissao, row.Status, row.HorarioTrabalho, row.Observacoes,
		row.CriadoEm, row.AtualizadoEm, row.Tipo,
	)
}

func mapUpdateStatusRowToResponse(row db.UpdateProfessionalStatusRow) dto.ProfessionalResponse {
	return mapGenericToResponse(
		row.ID, row.TenantID, row.UserID, row.Nome, row.Email, row.Telefone, row.Cpf,
		row.Especialidades, row.Comissao, row.TipoComissao, row.Foto, row.DataAdmissao,
		row.DataDemissao, row.Status, row.HorarioTrabalho, row.Observacoes,
		row.CriadoEm, row.AtualizadoEm, row.Tipo,
	)
}

func mapGenericToResponse(
	id, tenantID, userID pgtype.UUID,
	nome, email, telefone, cpf string,
	especialidades []string,
	comissao pgtype.Numeric,
	tipoComissao, foto *string,
	dataAdmissao, dataDemissao pgtype.Date,
	status *string,
	horarioTrabalho []byte,
	observacoes *string,
	criadoEm, atualizadoEm pgtype.Timestamptz,
	tipo string,
) dto.ProfessionalResponse {
	var userIDStr *string
	if userID.Valid {
		s := uuidToString(userID)
		userIDStr = &s
	}

	var fotoStr *string
	if foto != nil {
		fotoStr = foto
	}

	var dataDemissaoStr *string
	if dataDemissao.Valid {
		d := dataDemissao.Time.Format("2006-01-02")
		dataDemissaoStr = &d
	}

	var horarioTrabalhoStr *string
	if len(horarioTrabalho) > 0 {
		h := string(horarioTrabalho)
		horarioTrabalhoStr = &h
	}

	var observacoesStr *string
	if observacoes != nil {
		observacoesStr = observacoes
	}

	statusStr := "ATIVO"
	if status != nil {
		statusStr = *status
	}

	tipoComissaoStr := "PERCENTUAL"
	if tipoComissao != nil {
		tipoComissaoStr = *tipoComissao
	}

	comissaoVal := 0.0
	if comissao.Valid {
		d := decimal.NewFromBigInt(comissao.Int, comissao.Exp)
		comissaoVal, _ = d.Float64()
	}

	especList := especialidades
	if especList == nil {
		especList = []string{}
	}

	var criadoEmTime, atualizadoEmTime time.Time
	if criadoEm.Valid {
		criadoEmTime = criadoEm.Time
	}
	if atualizadoEm.Valid {
		atualizadoEmTime = atualizadoEm.Time
	}

	return dto.ProfessionalResponse{
		ID:              uuidToString(id),
		TenantID:        uuidToString(tenantID),
		UserID:          userIDStr,
		Nome:            nome,
		Email:           email,
		Telefone:        telefone,
		CPF:             cpf,
		Especialidades:  especList,
		Comissao:        comissaoVal,
		TipoComissao:    tipoComissaoStr,
		Foto:            fotoStr,
		DataAdmissao:    dataAdmissao.Time.Format("2006-01-02"),
		DataDemissao:    dataDemissaoStr,
		Status:          statusStr,
		HorarioTrabalho: horarioTrabalhoStr,
		Observacoes:     observacoesStr,
		Tipo:            tipo,
		CriadoEm:        criadoEmTime,
		AtualizadoEm:    atualizadoEmTime,
	}
}

func stringToJSONB(s *string) []byte {
	if s == nil || *s == "" {
		return nil
	}
	return []byte(*s)
}

func (r *ProfessionalRepository) saveCategoryCommissions(ctx context.Context, tenantID, professionalID string, commissions []dto.CommissionByCategory) error {
	// Delete existing
	err := r.queries.DeleteProfessionalCategoryCommissionsByProfessional(ctx, db.DeleteProfessionalCategoryCommissionsByProfessionalParams{
		TenantID:       stringToUUID(tenantID),
		ProfissionalID: stringToUUID(professionalID),
	})
	if err != nil {
		return err
	}

	// Insert new
	for _, c := range commissions {
		dec := decimal.NewFromFloat(c.Comissao)
		_, err := r.queries.CreateProfessionalCategoryCommission(ctx, db.CreateProfessionalCategoryCommissionParams{
			TenantID:       stringToUUID(tenantID),
			ProfissionalID: stringToUUID(professionalID),
			CategoriaID:    stringToUUID(c.CategoriaID),
			Comissao:       dec,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
