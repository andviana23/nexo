// Package postgres contém implementações de repositórios usando PostgreSQL.
package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// CustomerRepository implementa CustomerRepository usando PostgreSQL
type CustomerRepository struct {
	queries *db.Queries
}

// NewCustomerRepository cria uma nova instância
func NewCustomerRepository(queries *db.Queries) *CustomerRepository {
	return &CustomerRepository{queries: queries}
}

// =============================================================================
// CREATE
// =============================================================================

// Create persiste um novo cliente
func (r *CustomerRepository) Create(ctx context.Context, customer *entity.Customer) error {
	params := db.CreateCustomerParams{
		ID:                  stringToUUID(customer.ID),
		TenantID:            uuidToPgUUID(customer.TenantID),
		Nome:                customer.Nome,
		Telefone:            customer.Telefone,
		Email:               customer.Email,
		Cpf:                 customer.CPF,
		DataNascimento:      timeToPgDate(customer.DataNascimento),
		Genero:              customer.Genero,
		EnderecoLogradouro:  customer.EnderecoLogradouro,
		EnderecoNumero:      customer.EnderecoNumero,
		EnderecoComplemento: customer.EnderecoComplemento,
		EnderecoBairro:      customer.EnderecoBairro,
		EnderecoCidade:      customer.EnderecoCidade,
		EnderecoEstado:      customer.EnderecoEstado,
		EnderecoCep:         customer.EnderecoCEP,
		Observacoes:         customer.Observacoes,
		Tags:                customer.Tags,
	}

	_, err := r.queries.CreateCustomer(ctx, params)
	return err
}

// =============================================================================
// READ
// =============================================================================

// FindByID busca cliente por ID
func (r *CustomerRepository) FindByID(ctx context.Context, tenantID, id string) (*entity.Customer, error) {
	params := db.GetCustomerByIDParams{
		ID:       stringToUUID(id),
		TenantID: stringToUUID(tenantID),
	}

	row, err := r.queries.GetCustomerByID(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapClienteToCustomer(row), nil
}

// FindByPhone busca cliente por telefone
func (r *CustomerRepository) FindByPhone(ctx context.Context, tenantID, phone string) (*entity.Customer, error) {
	params := db.GetCustomerByPhoneParams{
		TenantID: stringToUUID(tenantID),
		Telefone: phone,
	}

	row, err := r.queries.GetCustomerByPhone(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapClienteToCustomer(row), nil
}

// FindByCPF busca cliente por CPF
func (r *CustomerRepository) FindByCPF(ctx context.Context, tenantID, cpf string) (*entity.Customer, error) {
	params := db.GetCustomerByCPFParams{
		TenantID: stringToUUID(tenantID),
		Cpf:      &cpf,
	}

	row, err := r.queries.GetCustomerByCPF(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapClienteToCustomer(row), nil
}

// =============================================================================
// UPDATE
// =============================================================================

// Update atualiza um cliente existente
func (r *CustomerRepository) Update(ctx context.Context, customer *entity.Customer) error {
	params := db.UpdateCustomerParams{
		ID:                  stringToUUID(customer.ID),
		TenantID:            uuidToPgUUID(customer.TenantID),
		Nome:                customer.Nome,
		Telefone:            customer.Telefone,
		Email:               customer.Email,
		Cpf:                 customer.CPF,
		DataNascimento:      timeToPgDate(customer.DataNascimento),
		Genero:              customer.Genero,
		EnderecoLogradouro:  customer.EnderecoLogradouro,
		EnderecoNumero:      customer.EnderecoNumero,
		EnderecoComplemento: customer.EnderecoComplemento,
		EnderecoBairro:      customer.EnderecoBairro,
		EnderecoCidade:      customer.EnderecoCidade,
		EnderecoEstado:      customer.EnderecoEstado,
		EnderecoCep:         customer.EnderecoCEP,
		Observacoes:         customer.Observacoes,
		Tags:                customer.Tags,
	}

	_, err := r.queries.UpdateCustomer(ctx, params)
	return err
}

// UpdateTags atualiza apenas as tags do cliente
func (r *CustomerRepository) UpdateTags(ctx context.Context, tenantID, id string, tags []string) error {
	params := db.UpdateCustomerTagsParams{
		ID:       stringToUUID(id),
		TenantID: stringToUUID(tenantID),
		Tags:     tags,
	}

	_, err := r.queries.UpdateCustomerTags(ctx, params)
	return err
}

// =============================================================================
// LIST
// =============================================================================

// List lista clientes com filtros e paginação
func (r *CustomerRepository) List(ctx context.Context, tenantID string, filter port.CustomerFilter) ([]*entity.Customer, int64, error) {
	// Defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	offset := (filter.Page - 1) * filter.PageSize

	// Define valores para nullable params
	var search *string
	if filter.Search != "" {
		search = &filter.Search
	}

	listParams := db.ListCustomersParams{
		TenantID:   stringToUUID(tenantID),
		Ativo:      filter.Ativo, // Já é *bool, nil significa "todos"
		Search:     search,
		Tags:       filter.Tags,
		OrderBy:    filter.OrderBy,
		PageSize:   int32(filter.PageSize),
		PageOffset: int32(offset),
	}

	rows, err := r.queries.ListCustomers(ctx, listParams)
	if err != nil {
		return nil, 0, err
	}

	countParams := db.CountCustomersParams{
		TenantID: stringToUUID(tenantID),
		Ativo:    filter.Ativo,
		Search:   search,
		Tags:     filter.Tags,
	}

	total, err := r.queries.CountCustomers(ctx, countParams)
	if err != nil {
		return nil, 0, err
	}

	customers := make([]*entity.Customer, 0, len(rows))
	for _, row := range rows {
		customers = append(customers, mapClienteToCustomer(row))
	}

	return customers, total, nil
}

// ListActive lista clientes ativos para selects
func (r *CustomerRepository) ListActive(ctx context.Context, tenantID string) ([]*port.CustomerSummary, error) {
	rows, err := r.queries.ListActiveCustomers(ctx, stringToUUID(tenantID))
	if err != nil {
		return nil, err
	}

	summaries := make([]*port.CustomerSummary, 0, len(rows))
	for _, row := range rows {
		summaries = append(summaries, &port.CustomerSummary{
			ID:       uuidToString(row.ID),
			Nome:     row.Nome,
			Telefone: row.Telefone,
			Email:    row.Email,
			Tags:     row.Tags,
		})
	}

	return summaries, nil
}

// Search busca rápida de clientes
func (r *CustomerRepository) Search(ctx context.Context, tenantID, query string) ([]*port.CustomerSummary, error) {
	params := db.SearchCustomersParams{
		TenantID: stringToUUID(tenantID),
		Column2:  &query,
	}

	rows, err := r.queries.SearchCustomers(ctx, params)
	if err != nil {
		return nil, err
	}

	summaries := make([]*port.CustomerSummary, 0, len(rows))
	for _, row := range rows {
		summaries = append(summaries, &port.CustomerSummary{
			ID:       uuidToString(row.ID),
			Nome:     row.Nome,
			Telefone: row.Telefone,
			Email:    row.Email,
			Tags:     row.Tags,
		})
	}

	return summaries, nil
}

// =============================================================================
// DELETE (Soft Delete)
// =============================================================================

// Inactivate inativa um cliente (soft delete)
func (r *CustomerRepository) Inactivate(ctx context.Context, tenantID, id string) error {
	params := db.InactivateCustomerParams{
		ID:       stringToUUID(id),
		TenantID: stringToUUID(tenantID),
	}

	return r.queries.InactivateCustomer(ctx, params)
}

// Reactivate reativa um cliente
func (r *CustomerRepository) Reactivate(ctx context.Context, tenantID, id string) error {
	params := db.ReactivateCustomerParams{
		ID:       stringToUUID(id),
		TenantID: stringToUUID(tenantID),
	}

	return r.queries.ReactivateCustomer(ctx, params)
}

// =============================================================================
// VALIDAÇÕES
// =============================================================================

// CheckPhoneExists verifica se telefone já existe
func (r *CustomerRepository) CheckPhoneExists(ctx context.Context, tenantID, phone string, excludeID *string) (bool, error) {
	params := db.CheckPhoneExistsParams{
		TenantID: stringToUUID(tenantID),
		Telefone: phone,
	}

	if excludeID != nil {
		params.Column3 = stringToUUID(*excludeID)
	}

	return r.queries.CheckPhoneExists(ctx, params)
}

// CheckCPFExists verifica se CPF já existe
func (r *CustomerRepository) CheckCPFExists(ctx context.Context, tenantID, cpf string, excludeID *string) (bool, error) {
	params := db.CheckCPFExistsParams{
		TenantID: stringToUUID(tenantID),
		Cpf:      &cpf,
	}

	if excludeID != nil {
		params.Column3 = stringToUUID(*excludeID)
	}

	return r.queries.CheckCPFExists(ctx, params)
}

// CheckEmailExists verifica se email já existe
func (r *CustomerRepository) CheckEmailExists(ctx context.Context, tenantID, email string, excludeID *string) (bool, error) {
	params := db.CheckEmailExistsParams{
		TenantID: stringToUUID(tenantID),
		Email:    &email,
	}

	if excludeID != nil {
		params.Column3 = stringToUUID(*excludeID)
	}

	return r.queries.CheckEmailExists(ctx, params)
}

// =============================================================================
// ESTATÍSTICAS E MÉTRICAS
// =============================================================================

// GetStats retorna estatísticas de clientes
func (r *CustomerRepository) GetStats(ctx context.Context, tenantID string) (*port.CustomerStats, error) {
	row, err := r.queries.GetCustomerStats(ctx, stringToUUID(tenantID))
	if err != nil {
		return nil, err
	}

	return &port.CustomerStats{
		TotalAtivos:        row.TotalAtivos,
		TotalInativos:      row.TotalInativos,
		NovosUltimos30Dias: row.NovosUltimos30Dias,
		TotalGeral:         row.TotalGeral,
	}, nil
}

// GetWithHistory busca cliente com histórico de atendimentos
func (r *CustomerRepository) GetWithHistory(ctx context.Context, tenantID, id string) (*port.CustomerWithHistory, error) {
	params := db.GetCustomerWithHistoryParams{
		ID:       stringToUUID(id),
		TenantID: stringToUUID(tenantID),
	}

	row, err := r.queries.GetCustomerWithHistory(ctx, params)
	if err != nil {
		return nil, err
	}

	customer := mapHistoryRowToCustomer(row)

	result := &port.CustomerWithHistory{
		Customer:          customer,
		TotalAtendimentos: row.TotalAtendimentos,
		TotalGasto:        fmt.Sprintf("%.2f", float64(row.TotalGasto)/100),
	}

	// Calcula ticket médio
	if row.TotalAtendimentos > 0 {
		ticketMedio := float64(row.TotalGasto) / float64(row.TotalAtendimentos) / 100
		result.TicketMedio = fmt.Sprintf("%.2f", ticketMedio)
	} else {
		result.TicketMedio = "0.00"
	}

	// Último atendimento
	if row.UltimoAtendimento != nil {
		if t, ok := row.UltimoAtendimento.(time.Time); ok {
			result.UltimoAtendimento = &t
		}
	}

	return result, nil
}

// GetDataForExport busca dados completos para exportação LGPD
func (r *CustomerRepository) GetDataForExport(ctx context.Context, tenantID, id string) (*port.CustomerExport, error) {
	params := db.GetCustomerDataForExportParams{
		ID:       stringToUUID(id),
		TenantID: stringToUUID(tenantID),
	}

	row, err := r.queries.GetCustomerDataForExport(ctx, params)
	if err != nil {
		return nil, err
	}

	customer := mapExportRowToCustomer(row)

	result := &port.CustomerExport{
		Customer:              customer,
		HistoricoAtendimentos: []port.CustomerAppointmentHistory{},
	}

	// Parse histórico JSON
	if row.HistoricoAtendimentos != nil {
		var rawHistorico []map[string]interface{}

		// Tenta fazer type assertion para []byte
		if jsonBytes, ok := row.HistoricoAtendimentos.([]byte); ok {
			if err := json.Unmarshal(jsonBytes, &rawHistorico); err == nil {
				for _, h := range rawHistorico {
					history := port.CustomerAppointmentHistory{}

					if data, ok := h["data"].(string); ok {
						if t, err := time.Parse(time.RFC3339, data); err == nil {
							history.Data = t
						}
					}
					if status, ok := h["status"].(string); ok {
						history.Status = status
					}
					if prof, ok := h["profissional"].(string); ok {
						history.Profissional = prof
					}
					if valor, ok := h["valor_total"].(float64); ok {
						history.ValorTotal = fmt.Sprintf("%.2f", valor/100)
					}

					result.HistoricoAtendimentos = append(result.HistoricoAtendimentos, history)
				}
			}
		}
	}

	// Calcular métricas
	result.TotalVisitas = int64(len(result.HistoricoAtendimentos))

	var totalGasto float64
	for _, h := range result.HistoricoAtendimentos {
		if v, err := parseFloat(h.ValorTotal); err == nil {
			totalGasto += v
		}
	}

	result.TotalGasto = fmt.Sprintf("%.2f", totalGasto)

	if result.TotalVisitas > 0 {
		result.TicketMedio = fmt.Sprintf("%.2f", totalGasto/float64(result.TotalVisitas))
	} else {
		result.TicketMedio = "0.00"
	}

	return result, nil
}

// =============================================================================
// Helpers de conversão
// =============================================================================

func stringToUUID(s string) pgtype.UUID {
	var pgUUID pgtype.UUID
	if uid, err := uuid.Parse(s); err == nil {
		pgUUID.Bytes = uid
		pgUUID.Valid = true
	}
	return pgUUID
}

func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	uid := uuid.UUID(u.Bytes)
	return uid.String()
}

func timeToPgDate(t *time.Time) pgtype.Date {
	var pgDate pgtype.Date
	if t != nil {
		pgDate.Time = *t
		pgDate.Valid = true
	}
	return pgDate
}

func pgDateToTime(d pgtype.Date) *time.Time {
	if !d.Valid {
		return nil
	}
	return &d.Time
}

func pgTimestamptzToTime(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// =============================================================================
// Mappers
// =============================================================================

func mapClienteToCustomer(row db.Cliente) *entity.Customer {
	var ativo bool
	if row.Ativo != nil {
		ativo = *row.Ativo
	}

	return &entity.Customer{
		ID:                  uuidToString(row.ID),
		TenantID: pgtypeToEntityUUID(row.TenantID),
		Nome:                row.Nome,
		Telefone:            row.Telefone,
		Email:               row.Email,
		CPF:                 row.Cpf,
		DataNascimento:      pgDateToTime(row.DataNascimento),
		Genero:              row.Genero,
		EnderecoLogradouro:  row.EnderecoLogradouro,
		EnderecoNumero:      row.EnderecoNumero,
		EnderecoComplemento: row.EnderecoComplemento,
		EnderecoBairro:      row.EnderecoBairro,
		EnderecoCidade:      row.EnderecoCidade,
		EnderecoEstado:      row.EnderecoEstado,
		EnderecoCEP:         row.EnderecoCep,
		Observacoes:         row.Observacoes,
		Tags:                row.Tags,
		Ativo:               ativo,
		CreatedAt:           pgTimestamptzToTime(row.CriadoEm),
		UpdatedAt:           pgTimestamptzToTime(row.AtualizadoEm),
	}
}

func mapHistoryRowToCustomer(row db.GetCustomerWithHistoryRow) *entity.Customer {
	var ativo bool
	if row.Ativo != nil {
		ativo = *row.Ativo
	}

	return &entity.Customer{
		ID:                  uuidToString(row.ID),
		TenantID: pgtypeToEntityUUID(row.TenantID),
		Nome:                row.Nome,
		Telefone:            row.Telefone,
		Email:               row.Email,
		CPF:                 row.Cpf,
		DataNascimento:      pgDateToTime(row.DataNascimento),
		Genero:              row.Genero,
		EnderecoLogradouro:  row.EnderecoLogradouro,
		EnderecoNumero:      row.EnderecoNumero,
		EnderecoComplemento: row.EnderecoComplemento,
		EnderecoBairro:      row.EnderecoBairro,
		EnderecoCidade:      row.EnderecoCidade,
		EnderecoEstado:      row.EnderecoEstado,
		EnderecoCEP:         row.EnderecoCep,
		Observacoes:         row.Observacoes,
		Tags:                row.Tags,
		Ativo:               ativo,
		CreatedAt:           pgTimestamptzToTime(row.CriadoEm),
		UpdatedAt:           pgTimestamptzToTime(row.AtualizadoEm),
	}
}

func mapExportRowToCustomer(row db.GetCustomerDataForExportRow) *entity.Customer {
	var ativo bool
	if row.Ativo != nil {
		ativo = *row.Ativo
	}

	return &entity.Customer{
		ID:                  uuidToString(row.ID),
		TenantID: pgtypeToEntityUUID(row.TenantID),
		Nome:                row.Nome,
		Telefone:            row.Telefone,
		Email:               row.Email,
		CPF:                 row.Cpf,
		DataNascimento:      pgDateToTime(row.DataNascimento),
		Genero:              row.Genero,
		EnderecoLogradouro:  row.EnderecoLogradouro,
		EnderecoNumero:      row.EnderecoNumero,
		EnderecoComplemento: row.EnderecoComplemento,
		EnderecoBairro:      row.EnderecoBairro,
		EnderecoCidade:      row.EnderecoCidade,
		EnderecoEstado:      row.EnderecoEstado,
		EnderecoCEP:         row.EnderecoCep,
		Observacoes:         row.Observacoes,
		Tags:                row.Tags,
		Ativo:               ativo,
		CreatedAt:           pgTimestamptzToTime(row.CriadoEm),
		UpdatedAt:           pgTimestamptzToTime(row.AtualizadoEm),
	}
}
