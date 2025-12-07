package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// commissionItemRepository implementa o repository de itens de comissão usando PostgreSQL
type commissionItemRepository struct {
	queries *db.Queries
}

// NewCommissionItemRepository cria uma nova instância do repository
func NewCommissionItemRepository(queries *db.Queries) repository.CommissionItemRepository {
	return &commissionItemRepository{
		queries: queries,
	}
}

// Create cria um novo item de comissão
func (r *commissionItemRepository) Create(ctx context.Context, item *entity.CommissionItem) (*entity.CommissionItem, error) {
	tenantID, err := uuid.Parse(item.TenantID)
	if err != nil {
		return nil, err
	}

	professionalID, err := uuid.Parse(item.ProfessionalID)
	if err != nil {
		return nil, err
	}

	var unitID pgtype.UUID
	if item.UnitID != nil {
		uid, err := uuid.Parse(*item.UnitID)
		if err != nil {
			return nil, err
		}
		unitID = pgtype.UUID{Bytes: uid, Valid: true}
	}

	var commandID pgtype.UUID
	if item.CommandID != nil {
		cid, err := uuid.Parse(*item.CommandID)
		if err != nil {
			return nil, err
		}
		commandID = pgtype.UUID{Bytes: cid, Valid: true}
	}

	var commandItemID pgtype.UUID
	if item.CommandItemID != nil {
		ciid, err := uuid.Parse(*item.CommandItemID)
		if err != nil {
			return nil, err
		}
		commandItemID = pgtype.UUID{Bytes: ciid, Valid: true}
	}

	var appointmentID pgtype.UUID
	if item.AppointmentID != nil {
		aid, err := uuid.Parse(*item.AppointmentID)
		if err != nil {
			return nil, err
		}
		appointmentID = pgtype.UUID{Bytes: aid, Valid: true}
	}

	var serviceID pgtype.UUID
	if item.ServiceID != nil {
		sid, err := uuid.Parse(*item.ServiceID)
		if err != nil {
			return nil, err
		}
		serviceID = pgtype.UUID{Bytes: sid, Valid: true}
	}

	var ruleID pgtype.UUID
	if item.RuleID != nil {
		rid, err := uuid.Parse(*item.RuleID)
		if err != nil {
			return nil, err
		}
		ruleID = pgtype.UUID{Bytes: rid, Valid: true}
	}

	result, err := r.queries.CreateCommissionItem(ctx, db.CreateCommissionItemParams{
		TenantID:         pgtype.UUID{Bytes: tenantID, Valid: true},
		UnitID:           unitID,
		ProfessionalID:   pgtype.UUID{Bytes: professionalID, Valid: true},
		CommandID:        commandID,
		CommandItemID:    commandItemID,
		AppointmentID:    appointmentID,
		ServiceID:        serviceID,
		ServiceName:      item.ServiceName,
		GrossValue:       item.GrossValue,
		CommissionRate:   item.CommissionRate,
		CommissionType:   item.CommissionType,
		CommissionValue:  item.CommissionValue,
		CommissionSource: item.CommissionSource,
		RuleID:           ruleID,
		ReferenceDate:    pgtype.Date{Time: item.ReferenceDate, Valid: true},
		Description:      item.Description,
		Status:           item.Status,
	})
	if err != nil {
		return nil, err
	}

	return commissionItemToDomain(result), nil
}

// CreateBatch cria múltiplos itens de comissão de uma vez
func (r *commissionItemRepository) CreateBatch(ctx context.Context, items []*entity.CommissionItem) ([]*entity.CommissionItem, error) {
	result := make([]*entity.CommissionItem, 0, len(items))
	for _, item := range items {
		created, err := r.Create(ctx, item)
		if err != nil {
			return nil, err
		}
		result = append(result, created)
	}
	return result, nil
}

// GetByID busca um item de comissão por ID
func (r *commissionItemRepository) GetByID(ctx context.Context, tenantID, id string) (*entity.CommissionItem, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	iid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.GetCommissionItemByID(ctx, db.GetCommissionItemByIDParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		ID:       pgtype.UUID{Bytes: iid, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("item de comissão não encontrado")
		}
		return nil, err
	}

	return commissionItemWithNamesToDomain(result), nil
}

// List lista itens de comissão com filtros
func (r *commissionItemRepository) List(ctx context.Context, tenantID string, professionalID *string, periodID *string, status *string, limit, offset int) ([]*entity.CommissionItem, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	// Se status for fornecido
	if status != nil && *status != "" {
		results, err := r.queries.ListCommissionItemsByStatus(ctx, db.ListCommissionItemsByStatusParams{
			TenantID: pgtype.UUID{Bytes: tid, Valid: true},
			Status:   *status,
			Limit:    int32(limit),
			Offset:   int32(offset),
		})
		if err != nil {
			return nil, err
		}

		items := make([]*entity.CommissionItem, 0, len(results))
		for _, result := range results {
			items = append(items, commissionItemByStatusRowToDomain(result))
		}
		return items, nil
	}

	// Se periodID for fornecido
	if periodID != nil && *periodID != "" {
		pid, err := uuid.Parse(*periodID)
		if err != nil {
			return nil, err
		}

		results, err := r.queries.ListCommissionItemsByPeriod(ctx, db.ListCommissionItemsByPeriodParams{
			TenantID: pgtype.UUID{Bytes: tid, Valid: true},
			PeriodID: pgtype.UUID{Bytes: pid, Valid: true},
		})
		if err != nil {
			return nil, err
		}

		items := make([]*entity.CommissionItem, 0, len(results))
		for _, result := range results {
			items = append(items, commissionItemByPeriodRowToDomain(result))
		}
		return items, nil
	}

	// Se professionalID for fornecido
	if professionalID != nil && *professionalID != "" {
		pid, err := uuid.Parse(*professionalID)
		if err != nil {
			return nil, err
		}

		results, err := r.queries.ListCommissionItemsByProfessional(ctx, db.ListCommissionItemsByProfessionalParams{
			TenantID:       pgtype.UUID{Bytes: tid, Valid: true},
			ProfessionalID: pgtype.UUID{Bytes: pid, Valid: true},
			Limit:          int32(limit),
			Offset:         int32(offset),
		})
		if err != nil {
			return nil, err
		}

		items := make([]*entity.CommissionItem, 0, len(results))
		for _, result := range results {
			items = append(items, commissionItemByProfessionalRowToDomain(result))
		}
		return items, nil
	}

	// Query padrão por tenant
	results, err := r.queries.ListCommissionItemsByTenant(ctx, db.ListCommissionItemsByTenantParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, err
	}

	items := make([]*entity.CommissionItem, 0, len(results))
	for _, result := range results {
		items = append(items, commissionItemByTenantRowToDomain(result))
	}
	return items, nil
}

// GetByProfessional busca itens de comissão por profissional
func (r *commissionItemRepository) GetByProfessional(ctx context.Context, tenantID, professionalID string) ([]*entity.CommissionItem, error) {
	return r.List(ctx, tenantID, &professionalID, nil, nil, 1000, 0)
}

// GetByPeriod busca itens de comissão por período
func (r *commissionItemRepository) GetByPeriod(ctx context.Context, tenantID, periodID string) ([]*entity.CommissionItem, error) {
	return r.List(ctx, tenantID, nil, &periodID, nil, 10000, 0)
}

// GetByCommandItem busca item de comissão por item de comanda
func (r *commissionItemRepository) GetByCommandItem(ctx context.Context, tenantID, commandItemID string) (*entity.CommissionItem, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	ciid, err := uuid.Parse(commandItemID)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.GetCommissionItemByCommandItem(ctx, db.GetCommissionItemByCommandItemParams{
		TenantID:      pgtype.UUID{Bytes: tid, Valid: true},
		CommandItemID: pgtype.UUID{Bytes: ciid, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Não encontrado retorna nil
		}
		return nil, err
	}

	return commissionItemToDomain(result), nil
}

// ListByCommand busca itens de comissão vinculados a uma comanda
// T-EST-003: Necessário para reverter comissões ao cancelar comanda
// Esta implementação busca por CommandItemID, já que não há query específica por CommandID
func (r *commissionItemRepository) ListByCommand(ctx context.Context, tenantID string, commandID uuid.UUID) ([]*entity.CommissionItem, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	// Buscar todos os itens pendentes e processados do tenant
	// e filtrar manualmente por CommandID
	// Nota: Uma query específica seria mais eficiente, mas usaremos o que temos

	// Buscar por status PENDENTE
	pendingItems, err := r.queries.ListCommissionItemsByStatus(ctx, db.ListCommissionItemsByStatusParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		Status:   "PENDENTE",
		Limit:    10000,
		Offset:   0,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// Buscar por status PROCESSADO
	processedItems, err := r.queries.ListCommissionItemsByStatus(ctx, db.ListCommissionItemsByStatusParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		Status:   "PROCESSADO",
		Limit:    10000,
		Offset:   0,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	var items []*entity.CommissionItem

	// Filtrar pendentes por CommandID
	for _, item := range pendingItems {
		if item.CommandID.Valid {
			var cmdID uuid.UUID
			copy(cmdID[:], item.CommandID.Bytes[:])
			if cmdID == commandID {
				items = append(items, commissionItemByStatusRowToDomain(item))
			}
		}
	}

	// Filtrar processados por CommandID
	for _, item := range processedItems {
		if item.CommandID.Valid {
			var cmdID uuid.UUID
			copy(cmdID[:], item.CommandID.Bytes[:])
			if cmdID == commandID {
				items = append(items, commissionItemByStatusRowToDomain(item))
			}
		}
	}

	return items, nil
}

// GetPendingByProfessional busca itens pendentes de um profissional
func (r *commissionItemRepository) GetPendingByProfessional(ctx context.Context, tenantID, professionalID string) ([]*entity.CommissionItem, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	pid, err := uuid.Parse(professionalID)
	if err != nil {
		return nil, err
	}

	results, err := r.queries.ListPendingCommissionItemsByProfessional(ctx, db.ListPendingCommissionItemsByProfessionalParams{
		TenantID:       pgtype.UUID{Bytes: tid, Valid: true},
		ProfessionalID: pgtype.UUID{Bytes: pid, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	items := make([]*entity.CommissionItem, 0, len(results))
	for _, result := range results {
		items = append(items, commissionItemPendingByProfRowToDomain(result))
	}
	return items, nil
}

// GetByDateRange busca itens de comissão por intervalo de datas
func (r *commissionItemRepository) GetByDateRange(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]*entity.CommissionItem, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	results, err := r.queries.ListCommissionItemsByDateRange(ctx, db.ListCommissionItemsByDateRangeParams{
		TenantID:        pgtype.UUID{Bytes: tid, Valid: true},
		ReferenceDate:   pgtype.Date{Time: startDate, Valid: true},
		ReferenceDate_2: pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	items := make([]*entity.CommissionItem, 0, len(results))
	for _, result := range results {
		items = append(items, commissionItemByDateRangeRowToDomain(result))
	}
	return items, nil
}

// GetTotalByPeriod retorna o total de comissão de um período
func (r *commissionItemRepository) GetTotalByPeriod(ctx context.Context, tenantID, periodID string) (float64, error) {
	// Busca itens do período e soma
	items, err := r.GetByPeriod(ctx, tenantID, periodID)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, item := range items {
		f, _ := item.CommissionValue.Float64()
		total += f
	}
	return total, nil
}

// GetTotalByProfessionalInRange retorna o total de comissão de um profissional em um intervalo
func (r *commissionItemRepository) GetTotalByProfessionalInRange(ctx context.Context, tenantID, professionalID string, startDate, endDate time.Time) (float64, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return 0, err
	}

	pid, err := uuid.Parse(professionalID)
	if err != nil {
		return 0, err
	}

	result, err := r.queries.SumCommissionsByProfessionalAndDateRange(ctx, db.SumCommissionsByProfessionalAndDateRangeParams{
		TenantID:        pgtype.UUID{Bytes: tid, Valid: true},
		ProfessionalID:  pgtype.UUID{Bytes: pid, Valid: true},
		ReferenceDate:   pgtype.Date{Time: startDate, Valid: true},
		ReferenceDate_2: pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		return 0, err
	}

	total, _ := result.TotalCommission.Float64()
	return total, nil
}

// GetSummaryByProfessional retorna resumo de comissões por profissional
func (r *commissionItemRepository) GetSummaryByProfessional(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]*entity.CommissionSummary, error) {
	// Implementação simplificada - agrupa dados manualmente
	items, err := r.GetByDateRange(ctx, tenantID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Agrupa por profissional
	byProfessional := make(map[string]*entity.CommissionSummary)
	for _, item := range items {
		if _, ok := byProfessional[item.ProfessionalID]; !ok {
			byProfessional[item.ProfessionalID] = &entity.CommissionSummary{
				ProfessionalID: item.ProfessionalID,
			}
		}
		summary := byProfessional[item.ProfessionalID]
		summary.TotalGross = summary.TotalGross.Add(item.GrossValue)
		summary.TotalCommission = summary.TotalCommission.Add(item.CommissionValue)
		summary.ItemsCount++
	}

	result := make([]*entity.CommissionSummary, 0, len(byProfessional))
	for _, summary := range byProfessional {
		result = append(result, summary)
	}
	return result, nil
}

// GetSummaryByService retorna resumo de comissões por serviço
func (r *commissionItemRepository) GetSummaryByService(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]*entity.CommissionByService, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	results, err := r.queries.GetCommissionSummaryByService(ctx, db.GetCommissionSummaryByServiceParams{
		TenantID:        pgtype.UUID{Bytes: tid, Valid: true},
		ReferenceDate:   pgtype.Date{Time: startDate, Valid: true},
		ReferenceDate_2: pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	summaries := make([]*entity.CommissionByService, 0, len(results))
	for _, result := range results {
		var serviceID string
		if result.ServiceID.Valid {
			serviceID = uuid.UUID(result.ServiceID.Bytes).String()
		}

		var serviceName string
		if result.ServiceName != nil {
			serviceName = *result.ServiceName
		}

		summaries = append(summaries, &entity.CommissionByService{
			ServiceID:       serviceID,
			ServiceName:     serviceName,
			TotalGross:      result.TotalGross,
			TotalCommission: result.TotalCommission,
			ItemsCount:      int(result.ItemsCount),
		})
	}
	return summaries, nil
}

// Process processa um item (vincula a um período)
func (r *commissionItemRepository) Process(ctx context.Context, tenantID, id, periodID string) (*entity.CommissionItem, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	iid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	pid, err := uuid.Parse(periodID)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.ProcessCommissionItem(ctx, db.ProcessCommissionItemParams{
		ID:       pgtype.UUID{Bytes: iid, Valid: true},
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		PeriodID: pgtype.UUID{Bytes: pid, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("item de comissão não encontrado ou não pode ser processado")
		}
		return nil, err
	}

	return commissionItemToDomain(result), nil
}

// AssignToPeriod vincula itens pendentes a um período
func (r *commissionItemRepository) AssignToPeriod(ctx context.Context, tenantID, professionalID, periodID string, startDate, endDate time.Time) (int64, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return 0, err
	}

	profID, err := uuid.Parse(professionalID)
	if err != nil {
		return 0, err
	}

	pid, err := uuid.Parse(periodID)
	if err != nil {
		return 0, err
	}

	err = r.queries.BulkProcessCommissionItems(ctx, db.BulkProcessCommissionItemsParams{
		TenantID:        pgtype.UUID{Bytes: tid, Valid: true},
		ProfessionalID:  pgtype.UUID{Bytes: profID, Valid: true},
		PeriodID:        pgtype.UUID{Bytes: pid, Valid: true},
		ReferenceDate:   pgtype.Date{Time: startDate, Valid: true},
		ReferenceDate_2: pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		return 0, err
	}

	// Retorna contagem estimada (não temos retorno real do bulk)
	return 0, nil
}

// Update atualiza um item de comissão
func (r *commissionItemRepository) Update(ctx context.Context, item *entity.CommissionItem) (*entity.CommissionItem, error) {
	id, err := uuid.Parse(item.ID)
	if err != nil {
		return nil, err
	}

	tenantID, err := uuid.Parse(item.TenantID)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.UpdateCommissionItem(ctx, db.UpdateCommissionItemParams{
		ID:               pgtype.UUID{Bytes: id, Valid: true},
		TenantID:         pgtype.UUID{Bytes: tenantID, Valid: true},
		CommissionRate:   item.CommissionRate,
		CommissionType:   item.CommissionType,
		CommissionValue:  item.CommissionValue,
		CommissionSource: item.CommissionSource,
		Description:      item.Description,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("item de comissão não encontrado")
		}
		return nil, err
	}

	return commissionItemToDomain(result), nil
}

// Delete remove um item de comissão (somente se PENDENTE)
func (r *commissionItemRepository) Delete(ctx context.Context, tenantID, id string) error {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return err
	}

	iid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteCommissionItem(ctx, db.DeleteCommissionItemParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		ID:       pgtype.UUID{Bytes: iid, Valid: true},
	})
}

// DeleteByCommandItem remove item de comissão por item de comanda
func (r *commissionItemRepository) DeleteByCommandItem(ctx context.Context, tenantID, commandItemID string) error {
	// Busca o item pelo command_item_id e deleta
	item, err := r.GetByCommandItem(ctx, tenantID, commandItemID)
	if err != nil {
		return err
	}
	if item == nil {
		return nil // Não existe, considera sucesso
	}

	return r.Delete(ctx, tenantID, item.ID)
}

// === Funções auxiliares de conversão ===

// commissionItemToDomain converte de db.CommissionItem para entity
func commissionItemToDomain(ci db.CommissionItem) *entity.CommissionItem {
	id := uuid.UUID(ci.ID.Bytes).String()
	tenantID := uuid.UUID(ci.TenantID.Bytes).String()
	professionalID := uuid.UUID(ci.ProfessionalID.Bytes).String()

	var unitID *string
	if ci.UnitID.Valid {
		uid := uuid.UUID(ci.UnitID.Bytes).String()
		unitID = &uid
	}

	var commandID *string
	if ci.CommandID.Valid {
		cid := uuid.UUID(ci.CommandID.Bytes).String()
		commandID = &cid
	}

	var commandItemID *string
	if ci.CommandItemID.Valid {
		ciid := uuid.UUID(ci.CommandItemID.Bytes).String()
		commandItemID = &ciid
	}

	var appointmentID *string
	if ci.AppointmentID.Valid {
		aid := uuid.UUID(ci.AppointmentID.Bytes).String()
		appointmentID = &aid
	}

	var serviceID *string
	if ci.ServiceID.Valid {
		sid := uuid.UUID(ci.ServiceID.Bytes).String()
		serviceID = &sid
	}

	var ruleID *string
	if ci.RuleID.Valid {
		rid := uuid.UUID(ci.RuleID.Bytes).String()
		ruleID = &rid
	}

	var periodID *string
	if ci.PeriodID.Valid {
		pid := uuid.UUID(ci.PeriodID.Bytes).String()
		periodID = &pid
	}

	var processedAt *time.Time
	if ci.ProcessedAt.Valid {
		processedAt = &ci.ProcessedAt.Time
	}

	return &entity.CommissionItem{
		ID:               id,
		TenantID:         tenantID,
		UnitID:           unitID,
		ProfessionalID:   professionalID,
		CommandID:        commandID,
		CommandItemID:    commandItemID,
		AppointmentID:    appointmentID,
		ServiceID:        serviceID,
		ServiceName:      ci.ServiceName,
		GrossValue:       ci.GrossValue,
		CommissionRate:   ci.CommissionRate,
		CommissionType:   ci.CommissionType,
		CommissionValue:  ci.CommissionValue,
		CommissionSource: ci.CommissionSource,
		RuleID:           ruleID,
		ReferenceDate:    ci.ReferenceDate.Time,
		Description:      ci.Description,
		Status:           ci.Status,
		PeriodID:         periodID,
		CreatedAt:        ci.CreatedAt.Time,
		UpdatedAt:        ci.UpdatedAt.Time,
		ProcessedAt:      processedAt,
	}
}

// commissionItemWithNamesToDomain converte GetCommissionItemByIDRow para entity
func commissionItemWithNamesToDomain(row db.GetCommissionItemByIDRow) *entity.CommissionItem {
	id := uuid.UUID(row.ID.Bytes).String()
	tenantID := uuid.UUID(row.TenantID.Bytes).String()
	professionalID := uuid.UUID(row.ProfessionalID.Bytes).String()

	var unitID *string
	if row.UnitID.Valid {
		uid := uuid.UUID(row.UnitID.Bytes).String()
		unitID = &uid
	}

	var commandID *string
	if row.CommandID.Valid {
		cid := uuid.UUID(row.CommandID.Bytes).String()
		commandID = &cid
	}

	var commandItemID *string
	if row.CommandItemID.Valid {
		ciid := uuid.UUID(row.CommandItemID.Bytes).String()
		commandItemID = &ciid
	}

	var appointmentID *string
	if row.AppointmentID.Valid {
		aid := uuid.UUID(row.AppointmentID.Bytes).String()
		appointmentID = &aid
	}

	var serviceID *string
	if row.ServiceID.Valid {
		sid := uuid.UUID(row.ServiceID.Bytes).String()
		serviceID = &sid
	}

	var ruleID *string
	if row.RuleID.Valid {
		rid := uuid.UUID(row.RuleID.Bytes).String()
		ruleID = &rid
	}

	var periodID *string
	if row.PeriodID.Valid {
		pid := uuid.UUID(row.PeriodID.Bytes).String()
		periodID = &pid
	}

	var processedAt *time.Time
	if row.ProcessedAt.Valid {
		processedAt = &row.ProcessedAt.Time
	}

	return &entity.CommissionItem{
		ID:               id,
		TenantID:         tenantID,
		UnitID:           unitID,
		ProfessionalID:   professionalID,
		CommandID:        commandID,
		CommandItemID:    commandItemID,
		AppointmentID:    appointmentID,
		ServiceID:        serviceID,
		ServiceName:      row.ServiceName,
		GrossValue:       row.GrossValue,
		CommissionRate:   row.CommissionRate,
		CommissionType:   row.CommissionType,
		CommissionValue:  row.CommissionValue,
		CommissionSource: row.CommissionSource,
		RuleID:           ruleID,
		ReferenceDate:    row.ReferenceDate.Time,
		Description:      row.Description,
		Status:           row.Status,
		PeriodID:         periodID,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
		ProcessedAt:      processedAt,
	}
}

// commissionItemByStatusRowToDomain converte ListCommissionItemsByStatusRow para entity
func commissionItemByStatusRowToDomain(row db.ListCommissionItemsByStatusRow) *entity.CommissionItem {
	id := uuid.UUID(row.ID.Bytes).String()
	tenantID := uuid.UUID(row.TenantID.Bytes).String()
	professionalID := uuid.UUID(row.ProfessionalID.Bytes).String()

	var unitID *string
	if row.UnitID.Valid {
		uid := uuid.UUID(row.UnitID.Bytes).String()
		unitID = &uid
	}

	var commandID *string
	if row.CommandID.Valid {
		cid := uuid.UUID(row.CommandID.Bytes).String()
		commandID = &cid
	}

	var commandItemID *string
	if row.CommandItemID.Valid {
		ciid := uuid.UUID(row.CommandItemID.Bytes).String()
		commandItemID = &ciid
	}

	var appointmentID *string
	if row.AppointmentID.Valid {
		aid := uuid.UUID(row.AppointmentID.Bytes).String()
		appointmentID = &aid
	}

	var serviceID *string
	if row.ServiceID.Valid {
		sid := uuid.UUID(row.ServiceID.Bytes).String()
		serviceID = &sid
	}

	var ruleID *string
	if row.RuleID.Valid {
		rid := uuid.UUID(row.RuleID.Bytes).String()
		ruleID = &rid
	}

	var periodID *string
	if row.PeriodID.Valid {
		pid := uuid.UUID(row.PeriodID.Bytes).String()
		periodID = &pid
	}

	var processedAt *time.Time
	if row.ProcessedAt.Valid {
		processedAt = &row.ProcessedAt.Time
	}

	return &entity.CommissionItem{
		ID:               id,
		TenantID:         tenantID,
		UnitID:           unitID,
		ProfessionalID:   professionalID,
		CommandID:        commandID,
		CommandItemID:    commandItemID,
		AppointmentID:    appointmentID,
		ServiceID:        serviceID,
		ServiceName:      row.ServiceName,
		GrossValue:       row.GrossValue,
		CommissionRate:   row.CommissionRate,
		CommissionType:   row.CommissionType,
		CommissionValue:  row.CommissionValue,
		CommissionSource: row.CommissionSource,
		RuleID:           ruleID,
		ReferenceDate:    row.ReferenceDate.Time,
		Description:      row.Description,
		Status:           row.Status,
		PeriodID:         periodID,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
		ProcessedAt:      processedAt,
	}
}

// commissionItemByPeriodRowToDomain converte ListCommissionItemsByPeriodRow para entity
func commissionItemByPeriodRowToDomain(row db.ListCommissionItemsByPeriodRow) *entity.CommissionItem {
	id := uuid.UUID(row.ID.Bytes).String()
	tenantID := uuid.UUID(row.TenantID.Bytes).String()
	professionalID := uuid.UUID(row.ProfessionalID.Bytes).String()

	var unitID *string
	if row.UnitID.Valid {
		uid := uuid.UUID(row.UnitID.Bytes).String()
		unitID = &uid
	}

	var commandID *string
	if row.CommandID.Valid {
		cid := uuid.UUID(row.CommandID.Bytes).String()
		commandID = &cid
	}

	var commandItemID *string
	if row.CommandItemID.Valid {
		ciid := uuid.UUID(row.CommandItemID.Bytes).String()
		commandItemID = &ciid
	}

	var appointmentID *string
	if row.AppointmentID.Valid {
		aid := uuid.UUID(row.AppointmentID.Bytes).String()
		appointmentID = &aid
	}

	var serviceID *string
	if row.ServiceID.Valid {
		sid := uuid.UUID(row.ServiceID.Bytes).String()
		serviceID = &sid
	}

	var ruleID *string
	if row.RuleID.Valid {
		rid := uuid.UUID(row.RuleID.Bytes).String()
		ruleID = &rid
	}

	var periodID *string
	if row.PeriodID.Valid {
		pid := uuid.UUID(row.PeriodID.Bytes).String()
		periodID = &pid
	}

	var processedAt *time.Time
	if row.ProcessedAt.Valid {
		processedAt = &row.ProcessedAt.Time
	}

	return &entity.CommissionItem{
		ID:               id,
		TenantID:         tenantID,
		UnitID:           unitID,
		ProfessionalID:   professionalID,
		CommandID:        commandID,
		CommandItemID:    commandItemID,
		AppointmentID:    appointmentID,
		ServiceID:        serviceID,
		ServiceName:      row.ServiceName,
		GrossValue:       row.GrossValue,
		CommissionRate:   row.CommissionRate,
		CommissionType:   row.CommissionType,
		CommissionValue:  row.CommissionValue,
		CommissionSource: row.CommissionSource,
		RuleID:           ruleID,
		ReferenceDate:    row.ReferenceDate.Time,
		Description:      row.Description,
		Status:           row.Status,
		PeriodID:         periodID,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
		ProcessedAt:      processedAt,
	}
}

// commissionItemByProfessionalRowToDomain converte ListCommissionItemsByProfessionalRow para entity
func commissionItemByProfessionalRowToDomain(row db.ListCommissionItemsByProfessionalRow) *entity.CommissionItem {
	id := uuid.UUID(row.ID.Bytes).String()
	tenantID := uuid.UUID(row.TenantID.Bytes).String()
	professionalID := uuid.UUID(row.ProfessionalID.Bytes).String()

	var unitID *string
	if row.UnitID.Valid {
		uid := uuid.UUID(row.UnitID.Bytes).String()
		unitID = &uid
	}

	var commandID *string
	if row.CommandID.Valid {
		cid := uuid.UUID(row.CommandID.Bytes).String()
		commandID = &cid
	}

	var commandItemID *string
	if row.CommandItemID.Valid {
		ciid := uuid.UUID(row.CommandItemID.Bytes).String()
		commandItemID = &ciid
	}

	var appointmentID *string
	if row.AppointmentID.Valid {
		aid := uuid.UUID(row.AppointmentID.Bytes).String()
		appointmentID = &aid
	}

	var serviceID *string
	if row.ServiceID.Valid {
		sid := uuid.UUID(row.ServiceID.Bytes).String()
		serviceID = &sid
	}

	var ruleID *string
	if row.RuleID.Valid {
		rid := uuid.UUID(row.RuleID.Bytes).String()
		ruleID = &rid
	}

	var periodID *string
	if row.PeriodID.Valid {
		pid := uuid.UUID(row.PeriodID.Bytes).String()
		periodID = &pid
	}

	var processedAt *time.Time
	if row.ProcessedAt.Valid {
		processedAt = &row.ProcessedAt.Time
	}

	return &entity.CommissionItem{
		ID:               id,
		TenantID:         tenantID,
		UnitID:           unitID,
		ProfessionalID:   professionalID,
		CommandID:        commandID,
		CommandItemID:    commandItemID,
		AppointmentID:    appointmentID,
		ServiceID:        serviceID,
		ServiceName:      row.ServiceName,
		GrossValue:       row.GrossValue,
		CommissionRate:   row.CommissionRate,
		CommissionType:   row.CommissionType,
		CommissionValue:  row.CommissionValue,
		CommissionSource: row.CommissionSource,
		RuleID:           ruleID,
		ReferenceDate:    row.ReferenceDate.Time,
		Description:      row.Description,
		Status:           row.Status,
		PeriodID:         periodID,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
		ProcessedAt:      processedAt,
	}
}

// commissionItemByTenantRowToDomain converte ListCommissionItemsByTenantRow para entity
func commissionItemByTenantRowToDomain(row db.ListCommissionItemsByTenantRow) *entity.CommissionItem {
	id := uuid.UUID(row.ID.Bytes).String()
	tenantID := uuid.UUID(row.TenantID.Bytes).String()
	professionalID := uuid.UUID(row.ProfessionalID.Bytes).String()

	var unitID *string
	if row.UnitID.Valid {
		uid := uuid.UUID(row.UnitID.Bytes).String()
		unitID = &uid
	}

	var commandID *string
	if row.CommandID.Valid {
		cid := uuid.UUID(row.CommandID.Bytes).String()
		commandID = &cid
	}

	var commandItemID *string
	if row.CommandItemID.Valid {
		ciid := uuid.UUID(row.CommandItemID.Bytes).String()
		commandItemID = &ciid
	}

	var appointmentID *string
	if row.AppointmentID.Valid {
		aid := uuid.UUID(row.AppointmentID.Bytes).String()
		appointmentID = &aid
	}

	var serviceID *string
	if row.ServiceID.Valid {
		sid := uuid.UUID(row.ServiceID.Bytes).String()
		serviceID = &sid
	}

	var ruleID *string
	if row.RuleID.Valid {
		rid := uuid.UUID(row.RuleID.Bytes).String()
		ruleID = &rid
	}

	var periodID *string
	if row.PeriodID.Valid {
		pid := uuid.UUID(row.PeriodID.Bytes).String()
		periodID = &pid
	}

	var processedAt *time.Time
	if row.ProcessedAt.Valid {
		processedAt = &row.ProcessedAt.Time
	}

	return &entity.CommissionItem{
		ID:               id,
		TenantID:         tenantID,
		UnitID:           unitID,
		ProfessionalID:   professionalID,
		CommandID:        commandID,
		CommandItemID:    commandItemID,
		AppointmentID:    appointmentID,
		ServiceID:        serviceID,
		ServiceName:      row.ServiceName,
		GrossValue:       row.GrossValue,
		CommissionRate:   row.CommissionRate,
		CommissionType:   row.CommissionType,
		CommissionValue:  row.CommissionValue,
		CommissionSource: row.CommissionSource,
		RuleID:           ruleID,
		ReferenceDate:    row.ReferenceDate.Time,
		Description:      row.Description,
		Status:           row.Status,
		PeriodID:         periodID,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
		ProcessedAt:      processedAt,
	}
}

// commissionItemPendingByProfRowToDomain converte ListPendingCommissionItemsByProfessionalRow para entity
func commissionItemPendingByProfRowToDomain(row db.ListPendingCommissionItemsByProfessionalRow) *entity.CommissionItem {
	id := uuid.UUID(row.ID.Bytes).String()
	tenantID := uuid.UUID(row.TenantID.Bytes).String()
	professionalID := uuid.UUID(row.ProfessionalID.Bytes).String()

	var unitID *string
	if row.UnitID.Valid {
		uid := uuid.UUID(row.UnitID.Bytes).String()
		unitID = &uid
	}

	var commandID *string
	if row.CommandID.Valid {
		cid := uuid.UUID(row.CommandID.Bytes).String()
		commandID = &cid
	}

	var commandItemID *string
	if row.CommandItemID.Valid {
		ciid := uuid.UUID(row.CommandItemID.Bytes).String()
		commandItemID = &ciid
	}

	var appointmentID *string
	if row.AppointmentID.Valid {
		aid := uuid.UUID(row.AppointmentID.Bytes).String()
		appointmentID = &aid
	}

	var serviceID *string
	if row.ServiceID.Valid {
		sid := uuid.UUID(row.ServiceID.Bytes).String()
		serviceID = &sid
	}

	var ruleID *string
	if row.RuleID.Valid {
		rid := uuid.UUID(row.RuleID.Bytes).String()
		ruleID = &rid
	}

	var periodID *string
	if row.PeriodID.Valid {
		pid := uuid.UUID(row.PeriodID.Bytes).String()
		periodID = &pid
	}

	var processedAt *time.Time
	if row.ProcessedAt.Valid {
		processedAt = &row.ProcessedAt.Time
	}

	return &entity.CommissionItem{
		ID:               id,
		TenantID:         tenantID,
		UnitID:           unitID,
		ProfessionalID:   professionalID,
		CommandID:        commandID,
		CommandItemID:    commandItemID,
		AppointmentID:    appointmentID,
		ServiceID:        serviceID,
		ServiceName:      row.ServiceName,
		GrossValue:       row.GrossValue,
		CommissionRate:   row.CommissionRate,
		CommissionType:   row.CommissionType,
		CommissionValue:  row.CommissionValue,
		CommissionSource: row.CommissionSource,
		RuleID:           ruleID,
		ReferenceDate:    row.ReferenceDate.Time,
		Description:      row.Description,
		Status:           row.Status,
		PeriodID:         periodID,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
		ProcessedAt:      processedAt,
	}
}

// commissionItemByDateRangeRowToDomain converte ListCommissionItemsByDateRangeRow para entity
func commissionItemByDateRangeRowToDomain(row db.ListCommissionItemsByDateRangeRow) *entity.CommissionItem {
	id := uuid.UUID(row.ID.Bytes).String()
	tenantID := uuid.UUID(row.TenantID.Bytes).String()
	professionalID := uuid.UUID(row.ProfessionalID.Bytes).String()

	var unitID *string
	if row.UnitID.Valid {
		uid := uuid.UUID(row.UnitID.Bytes).String()
		unitID = &uid
	}

	var commandID *string
	if row.CommandID.Valid {
		cid := uuid.UUID(row.CommandID.Bytes).String()
		commandID = &cid
	}

	var commandItemID *string
	if row.CommandItemID.Valid {
		ciid := uuid.UUID(row.CommandItemID.Bytes).String()
		commandItemID = &ciid
	}

	var appointmentID *string
	if row.AppointmentID.Valid {
		aid := uuid.UUID(row.AppointmentID.Bytes).String()
		appointmentID = &aid
	}

	var serviceID *string
	if row.ServiceID.Valid {
		sid := uuid.UUID(row.ServiceID.Bytes).String()
		serviceID = &sid
	}

	var ruleID *string
	if row.RuleID.Valid {
		rid := uuid.UUID(row.RuleID.Bytes).String()
		ruleID = &rid
	}

	var periodID *string
	if row.PeriodID.Valid {
		pid := uuid.UUID(row.PeriodID.Bytes).String()
		periodID = &pid
	}

	var processedAt *time.Time
	if row.ProcessedAt.Valid {
		processedAt = &row.ProcessedAt.Time
	}

	return &entity.CommissionItem{
		ID:               id,
		TenantID:         tenantID,
		UnitID:           unitID,
		ProfessionalID:   professionalID,
		CommandID:        commandID,
		CommandItemID:    commandItemID,
		AppointmentID:    appointmentID,
		ServiceID:        serviceID,
		ServiceName:      row.ServiceName,
		GrossValue:       row.GrossValue,
		CommissionRate:   row.CommissionRate,
		CommissionType:   row.CommissionType,
		CommissionValue:  row.CommissionValue,
		CommissionSource: row.CommissionSource,
		RuleID:           ruleID,
		ReferenceDate:    row.ReferenceDate.Time,
		Description:      row.Description,
		Status:           row.Status,
		PeriodID:         periodID,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
		ProcessedAt:      processedAt,
	}
}
