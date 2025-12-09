// Package postgres implementa os readers para consultas auxiliares do módulo de agendamentos.
package postgres

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ============================================================================
// ProfessionalReader Implementation
// ============================================================================

// ProfessionalReaderPG implementa port.ProfessionalReader usando sqlc.
type ProfessionalReaderPG struct {
	queries *db.Queries
}

// NewProfessionalReader cria uma nova instância do reader de profissionais.
func NewProfessionalReader(queries *db.Queries) *ProfessionalReaderPG {
	return &ProfessionalReaderPG{queries: queries}
}

// Exists verifica se um profissional existe e está ativo.
func (r *ProfessionalReaderPG) Exists(ctx context.Context, tenantID, professionalID string) (bool, error) {
	params := db.ProfessionalExistsParams{
		ID:       uuidStringToPgtype(professionalID),
		TenantID: uuidStringToPgtype(tenantID),
	}

	exists, err := r.queries.ProfessionalExists(ctx, params)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar profissional: %w", err)
	}

	return exists, nil
}

// FindByID busca dados básicos do profissional.
func (r *ProfessionalReaderPG) FindByID(ctx context.Context, tenantID, professionalID string) (*port.ProfessionalInfo, error) {
	params := db.GetProfessionalInfoParams{
		ID:       uuidStringToPgtype(professionalID),
		TenantID: uuidStringToPgtype(tenantID),
	}

	row, err := r.queries.GetProfessionalInfo(ctx, params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar profissional: %w", err)
	}

	status := ""
	if row.Status != nil {
		status = *row.Status
	}
	color := ""
	if row.Cor != nil {
		color = *row.Cor
	}

	// Converter string para *string (sqlc retorna string, mas pode ser vazio)
	var comissao *string
	if row.Comissao != "" {
		comissao = &row.Comissao
	}

	return &port.ProfessionalInfo{
		ID:           pgUUIDToString(row.ID),
		Name:         row.Nome,
		Status:       status,
		Color:        color,
		Comissao:     comissao,
		TipoComissao: row.TipoComissao,
	}, nil
}

// ListActive lista profissionais ativos.
func (r *ProfessionalReaderPG) ListActive(ctx context.Context, tenantID string) ([]*port.ProfessionalInfo, error) {
	rows, err := r.queries.ListActiveProfessionals(ctx, uuidStringToPgtype(tenantID))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar profissionais ativos: %w", err)
	}

	professionals := make([]*port.ProfessionalInfo, 0, len(rows))
	for _, row := range rows {
		status := ""
		if row.Status != nil {
			status = *row.Status
		}
		color := ""
		if row.Cor != nil {
			color = *row.Cor
		}
		professionals = append(professionals, &port.ProfessionalInfo{
			ID:     pgUUIDToString(row.ID),
			Name:   row.Nome,
			Status: status,
			Color:  color,
		})
	}

	return professionals, nil
}

// ============================================================================
// CustomerReader Implementation
// ============================================================================

// CustomerReaderPG implementa port.CustomerReader usando sqlc.
type CustomerReaderPG struct {
	queries *db.Queries
}

// NewCustomerReader cria uma nova instância do reader de clientes.
func NewCustomerReader(queries *db.Queries) *CustomerReaderPG {
	return &CustomerReaderPG{queries: queries}
}

// Exists verifica se um cliente existe e está ativo.
func (r *CustomerReaderPG) Exists(ctx context.Context, tenantID, customerID string) (bool, error) {
	params := db.CustomerExistsParams{
		ID:       uuidStringToPgtype(customerID),
		TenantID: uuidStringToPgtype(tenantID),
	}

	exists, err := r.queries.CustomerExists(ctx, params)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar cliente: %w", err)
	}

	return exists, nil
}

// FindByID busca dados básicos do cliente.
func (r *CustomerReaderPG) FindByID(ctx context.Context, tenantID, customerID string) (*port.CustomerInfo, error) {
	params := db.GetCustomerInfoParams{
		ID:       uuidStringToPgtype(customerID),
		TenantID: uuidStringToPgtype(tenantID),
	}

	row, err := r.queries.GetCustomerInfo(ctx, params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar cliente: %w", err)
	}

	email := ""
	if row.Email != nil {
		email = *row.Email
	}

	return &port.CustomerInfo{
		ID:    pgUUIDToString(row.ID),
		Name:  row.Nome,
		Phone: row.Telefone,
		Email: email,
	}, nil
}

// ============================================================================
// ServiceReader Implementation
// ============================================================================

// ServiceReaderPG implementa port.ServiceReader usando sqlc.
type ServiceReaderPG struct {
	queries *db.Queries
}

// NewServiceReader cria uma nova instância do reader de serviços.
func NewServiceReader(queries *db.Queries) *ServiceReaderPG {
	return &ServiceReaderPG{queries: queries}
}

// Exists verifica se um serviço existe e está ativo.
func (r *ServiceReaderPG) Exists(ctx context.Context, tenantID, serviceID string) (bool, error) {
	params := db.ServiceExistsParams{
		ID:       uuidStringToPgtype(serviceID),
		TenantID: uuidStringToPgtype(tenantID),
	}

	exists, err := r.queries.ServiceExists(ctx, params)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar serviço: %w", err)
	}

	return exists, nil
}

// FindByID busca dados do serviço.
func (r *ServiceReaderPG) FindByID(ctx context.Context, tenantID, serviceID string) (*port.ServiceInfo, error) {
	params := db.GetServiceInfoParams{
		ID:       uuidStringToPgtype(serviceID),
		TenantID: uuidStringToPgtype(tenantID),
	}

	row, err := r.queries.GetServiceInfo(ctx, params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar serviço: %w", err)
	}

	active := true
	if row.Ativo != nil {
		active = *row.Ativo
	}

	// Converter string para *string
	var comissao *string
	if row.Comissao != "" {
		comissao = &row.Comissao
	}

	return &port.ServiceInfo{
		ID:       pgUUIDToString(row.ID),
		Name:     row.Nome,
		Price:    valueobject.NewMoneyFromDecimal(row.Preco),
		Duration: int(row.Duracao),
		Active:   active,
		Comissao: comissao,
	}, nil
}

// FindByIDs busca múltiplos serviços.
func (r *ServiceReaderPG) FindByIDs(ctx context.Context, tenantID string, serviceIDs []string) ([]*port.ServiceInfo, error) {
	// Converter []string para []pgtype.UUID
	pgUUIDs := make([]pgtype.UUID, 0, len(serviceIDs))
	for _, id := range serviceIDs {
		pgUUIDs = append(pgUUIDs, uuidStringToPgtype(id))
	}

	params := db.GetServicesByIDsParams{
		TenantID: uuidStringToPgtype(tenantID),
		Column2:  pgUUIDs,
	}

	rows, err := r.queries.GetServicesByIDs(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar serviços: %w", err)
	}

	services := make([]*port.ServiceInfo, 0, len(rows))
	for _, row := range rows {
		active := true
		if row.Ativo != nil {
			active = *row.Ativo
		}
		// Converter string para *string
		var comissao *string
		if row.Comissao != "" {
			comissao = &row.Comissao
		}
		services = append(services, &port.ServiceInfo{
			ID:       pgUUIDToString(row.ID),
			Name:     row.Nome,
			Price:    valueobject.NewMoneyFromDecimal(row.Preco),
			Duration: int(row.Duracao),
			Active:   active,
			Comissao: comissao,
		})
	}

	return services, nil
}
