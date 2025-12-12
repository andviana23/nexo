// Package postgres implementa os repositórios usando PostgreSQL e sqlc.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

// AppointmentRepository implementa port.AppointmentRepository usando sqlc.
type AppointmentRepository struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

// NewAppointmentRepository cria uma nova instância do repositório.
func NewAppointmentRepository(queries *db.Queries, pool *pgxpool.Pool) *AppointmentRepository {
	return &AppointmentRepository{
		queries: queries,
		pool:    pool,
	}
}

// Create cria um novo agendamento com seus serviços (transação).
func (r *AppointmentRepository) Create(ctx context.Context, appointment *entity.Appointment) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// 1. Criar o agendamento
	params := db.CreateAppointmentParams{
		ID:                    uuidStringToPgtype(appointment.ID),
		TenantID:              entityUUIDToPgtype(appointment.TenantID),
		UnitID:                entityUUIDToPgtype(appointment.UnitID),
		ProfessionalID:        uuidStringToPgtype(appointment.ProfessionalID),
		CustomerID:            uuidStringToPgtype(appointment.CustomerID),
		StartTime:             timestampToTimestamptz(appointment.StartTime),
		EndTime:               timestampToTimestamptz(appointment.EndTime),
		Status:                appointment.Status.String(),
		TotalPrice:            appointment.TotalPrice.Value(),
		Notes:                 strPtrToPgText(appointment.Notes),
		CanceledReason:        strPtrToPgText(appointment.CanceledReason),
		GoogleCalendarEventID: strPtrToPgText(appointment.GoogleCalendarEventID),
		CommandID:             uuidStrPtrToPgtype(appointment.CommandID),
	}

	result, err := qtx.CreateAppointment(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar agendamento: %w", err)
	}

	// 2. Criar os serviços do agendamento
	for _, svc := range appointment.Services {
		svcParams := db.CreateAppointmentServiceParams{
			AppointmentID:     uuidStringToPgtype(appointment.ID),
			ServiceID:         uuidStringToPgtype(svc.ServiceID),
			PriceAtBooking:    svc.PriceAtBooking.Value(),
			DurationAtBooking: int32(svc.DurationAtBooking),
		}
		if err := qtx.CreateAppointmentService(ctx, svcParams); err != nil {
			return fmt.Errorf("erro ao criar serviço do agendamento: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	// Atualizar timestamps da entidade
	appointment.CreatedAt = timestamptzToTime(result.CreatedAt)
	appointment.UpdatedAt = timestamptzToTime(result.UpdatedAt)

	return nil
}

// FindByID busca um agendamento por ID.
func (r *AppointmentRepository) FindByID(ctx context.Context, tenantID, unitID, id string) (*entity.Appointment, error) {
	params := db.GetAppointmentByIDParams{
		ID:       uuidStringToPgtype(id),
		TenantID: uuidStringToPgtype(tenantID),
		UnitID:   uuidStringToPgtype(unitID),
	}

	row, err := r.queries.GetAppointmentByID(ctx, params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrAppointmentNotFound
		}
		return nil, fmt.Errorf("erro ao buscar agendamento: %w", err)
	}

	// Buscar serviços
	services, err := r.queries.GetAppointmentServices(ctx, uuidStringToPgtype(id))
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar serviços do agendamento: %w", err)
	}

	return r.rowToDomain(&row, services), nil
}

// Update atualiza um agendamento existente.
func (r *AppointmentRepository) Update(ctx context.Context, appointment *entity.Appointment) error {
	params := db.UpdateAppointmentParams{
		ID:                    uuidStringToPgtype(appointment.ID),
		TenantID:              entityUUIDToPgtype(appointment.TenantID),
		ProfessionalID:        uuidStringToPgtype(appointment.ProfessionalID),
		StartTime:             timestampToTimestamptz(appointment.StartTime),
		EndTime:               timestampToTimestamptz(appointment.EndTime),
		Status:                appointment.Status.String(),
		TotalPrice:            appointment.TotalPrice.Value(),
		Notes:                 strPtrToPgText(appointment.Notes),
		CanceledReason:        strPtrToPgText(appointment.CanceledReason),
		GoogleCalendarEventID: strPtrToPgText(appointment.GoogleCalendarEventID),
		CheckedInAt:           timePtrToPgTimestamptz(appointment.CheckedInAt),
		StartedAt:             timePtrToPgTimestamptz(appointment.StartedAt),
		FinishedAt:            timePtrToPgTimestamptz(appointment.FinishedAt),
		CommandID:             uuidStrPtrToPgtype(appointment.CommandID),
	}

	result, err := r.queries.UpdateAppointment(ctx, params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrAppointmentNotFound
		}
		return fmt.Errorf("erro ao atualizar agendamento: %w", err)
	}

	appointment.UpdatedAt = timestamptzToTime(result.UpdatedAt)
	return nil
}

// Delete remove um agendamento (soft delete via status CANCELED).
func (r *AppointmentRepository) Delete(ctx context.Context, tenantID, unitID, id string) error {
	params := db.DeleteAppointmentParams{
		ID:       uuidStringToPgtype(id),
		TenantID: uuidStringToPgtype(tenantID),
		UnitID:   uuidStringToPgtype(unitID),
	}
	return r.queries.DeleteAppointment(ctx, params)
}

// List lista agendamentos com filtros.
func (r *AppointmentRepository) List(ctx context.Context, tenantID string, filter port.AppointmentFilter) ([]*entity.Appointment, int64, error) {
	// Preparar parâmetros
	params := db.ListAppointmentsParams{
		TenantID: uuidStringToPgtype(tenantID),
		UnitID:   uuidStringToPgtype(filter.UnitID),
		Limit:    int32(filter.PageSize),
		Offset:   int32((filter.Page - 1) * filter.PageSize),
	}

	// Filtros opcionais
	if filter.ProfessionalID != "" {
		params.Column2 = uuidStringToPgtype(filter.ProfessionalID)
	}
	if filter.CustomerID != "" {
		params.Column3 = uuidStringToPgtype(filter.CustomerID)
	}
	// Converter array de status para []string
	if len(filter.Statuses) > 0 {
		statusStrings := make([]string, len(filter.Statuses))
		for i, s := range filter.Statuses {
			statusStrings[i] = s.String()
		}
		params.Column4 = statusStrings
	}
	if !filter.StartDate.IsZero() {
		params.Column5 = timestampToTimestamptz(filter.StartDate)
	}
	if !filter.EndDate.IsZero() {
		params.Column6 = timestampToTimestamptz(filter.EndDate)
	}

	// Buscar lista
	rows, err := r.queries.ListAppointments(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao listar agendamentos: %w", err)
	}

	// Contar total
	countParams := db.CountAppointmentsParams{
		TenantID: params.TenantID,
		Column2:  params.Column2,
		Column3:  params.Column3,
		Column4:  params.Column4,
		Column5:  params.Column5,
		Column6:  params.Column6,
	}
	total, err := r.queries.CountAppointments(ctx, countParams)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao contar agendamentos: %w", err)
	}

	// Converter para entidades
	appointments := make([]*entity.Appointment, 0, len(rows))
	appointmentIDs := make([]pgtype.UUID, 0, len(rows))
	for _, row := range rows {
		appointments = append(appointments, r.listRowToDomain(&row))
		appointmentIDs = append(appointmentIDs, row.ID)
	}

	// Carregar serviços para todos os agendamentos de uma vez (evita N+1)
	if len(appointmentIDs) > 0 {
		services, err := r.queries.GetServicesForAppointments(ctx, appointmentIDs)
		if err != nil {
			return nil, 0, fmt.Errorf("erro ao buscar serviços dos agendamentos: %w", err)
		}

		// Mapear serviços por appointment_id
		servicesByAppointment := make(map[string][]entity.AppointmentService)
		for _, svc := range services {
			apptID := pgUUIDToString(svc.AppointmentID)
			servicesByAppointment[apptID] = append(servicesByAppointment[apptID], entity.AppointmentService{
				AppointmentID:     apptID,
				ServiceID:         pgUUIDToString(svc.ServiceID),
				PriceAtBooking:    valueobject.NewMoneyFromDecimal(svc.PriceAtBooking),
				DurationAtBooking: int(svc.DurationAtBooking),
				CreatedAt:         timestamptzToTime(svc.CreatedAt),
				ServiceName:       svc.ServiceName,
			})
		}

		// Atribuir serviços aos agendamentos
		for _, appt := range appointments {
			if svcs, ok := servicesByAppointment[appt.ID]; ok {
				appt.Services = svcs
			}
		}
	}

	return appointments, total, nil
}

// ListByProfessionalAndDateRange lista agendamentos de um profissional em um período.
func (r *AppointmentRepository) ListByProfessionalAndDateRange(
	ctx context.Context,
	tenantID string,
	unitID string,
	professionalID string,
	startDate, endDate time.Time,
) ([]*entity.Appointment, error) {
	params := db.ListAppointmentsByProfessionalAndDateRangeParams{
		TenantID:       uuidStringToPgtype(tenantID),
		UnitID:         uuidStringToPgtype(unitID),
		ProfessionalID: uuidStringToPgtype(professionalID),
		StartTime:      timestampToTimestamptz(startDate),
		StartTime_2:    timestampToTimestamptz(endDate),
	}

	rows, err := r.queries.ListAppointmentsByProfessionalAndDateRange(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar agendamentos por profissional: %w", err)
	}

	appointments := make([]*entity.Appointment, 0, len(rows))
	appointmentIDs := make([]pgtype.UUID, 0, len(rows))
	for _, row := range rows {
		appointments = append(appointments, r.professionalRangeRowToDomain(&row))
		appointmentIDs = append(appointmentIDs, row.ID)
	}

	// Carregar serviços para todos os agendamentos de uma vez (evita N+1)
	if len(appointmentIDs) > 0 {
		services, err := r.queries.GetServicesForAppointments(ctx, appointmentIDs)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar serviços dos agendamentos: %w", err)
		}

		servicesByAppointment := make(map[string][]entity.AppointmentService)
		for _, svc := range services {
			apptID := pgUUIDToString(svc.AppointmentID)
			servicesByAppointment[apptID] = append(servicesByAppointment[apptID], entity.AppointmentService{
				AppointmentID:     apptID,
				ServiceID:         pgUUIDToString(svc.ServiceID),
				PriceAtBooking:    valueobject.NewMoneyFromDecimal(svc.PriceAtBooking),
				DurationAtBooking: int(svc.DurationAtBooking),
				CreatedAt:         timestamptzToTime(svc.CreatedAt),
				ServiceName:       svc.ServiceName,
			})
		}

		for _, appt := range appointments {
			if svcs, ok := servicesByAppointment[appt.ID]; ok {
				appt.Services = svcs
			}
		}
	}

	return appointments, nil
}

// ListByCustomer lista agendamentos de um cliente.
func (r *AppointmentRepository) ListByCustomer(ctx context.Context, tenantID, unitID, customerID string) ([]*entity.Appointment, error) {
	params := db.ListAppointmentsByCustomerParams{
		TenantID:   uuidStringToPgtype(tenantID),
		UnitID:     uuidStringToPgtype(unitID),
		CustomerID: uuidStringToPgtype(customerID),
	}

	rows, err := r.queries.ListAppointmentsByCustomer(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar agendamentos por cliente: %w", err)
	}

	appointments := make([]*entity.Appointment, 0, len(rows))
	appointmentIDs := make([]pgtype.UUID, 0, len(rows))
	for _, row := range rows {
		appointments = append(appointments, r.customerRowToDomain(&row))
		appointmentIDs = append(appointmentIDs, row.ID)
	}

	// Carregar serviços para todos os agendamentos de uma vez (evita N+1)
	if len(appointmentIDs) > 0 {
		services, err := r.queries.GetServicesForAppointments(ctx, appointmentIDs)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar serviços dos agendamentos: %w", err)
		}

		servicesByAppointment := make(map[string][]entity.AppointmentService)
		for _, svc := range services {
			apptID := pgUUIDToString(svc.AppointmentID)
			servicesByAppointment[apptID] = append(servicesByAppointment[apptID], entity.AppointmentService{
				AppointmentID:     apptID,
				ServiceID:         pgUUIDToString(svc.ServiceID),
				PriceAtBooking:    valueobject.NewMoneyFromDecimal(svc.PriceAtBooking),
				DurationAtBooking: int(svc.DurationAtBooking),
				CreatedAt:         timestamptzToTime(svc.CreatedAt),
				ServiceName:       svc.ServiceName,
			})
		}

		for _, appt := range appointments {
			if svcs, ok := servicesByAppointment[appt.ID]; ok {
				appt.Services = svcs
			}
		}
	}

	return appointments, nil
}

// CheckConflict verifica se há conflito de horário.
func (r *AppointmentRepository) CheckConflict(
	ctx context.Context,
	tenantID string,
	unitID string,
	professionalID string,
	startTime, endTime time.Time,
	excludeAppointmentID string,
) (bool, error) {
	excludeID := pgtype.UUID{}
	if excludeAppointmentID != "" {
		excludeID = uuidStringToPgtype(excludeAppointmentID)
	}

	params := db.CheckAppointmentConflictParams{
		TenantID:       uuidStringToPgtype(tenantID),
		UnitID:         uuidStringToPgtype(unitID),
		ProfessionalID: uuidStringToPgtype(professionalID),
		ID:             excludeID,
		StartTime:      timestampToTimestamptz(startTime),
		EndTime:        timestampToTimestamptz(endTime),
	}

	hasConflict, err := r.queries.CheckAppointmentConflict(ctx, params)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar conflito: %w", err)
	}

	return hasConflict, nil
}

// CheckBlockedTimeConflict verifica se há conflito com horários bloqueados.
func (r *AppointmentRepository) CheckBlockedTimeConflict(
	ctx context.Context,
	tenantID string,
	unitID string,
	professionalID string,
	startTime, endTime time.Time,
) (bool, error) {
	params := db.CheckBlockedTimeConflictForAppointmentParams{
		TenantID:       uuidStringToPgtype(tenantID),
		UnitID:         uuidStringToPgtype(unitID),
		ProfessionalID: uuidStringToPgtype(professionalID),
		StartTime:      timestampToTimestamptz(startTime),
		EndTime:        timestampToTimestamptz(endTime),
	}

	hasConflict, err := r.queries.CheckBlockedTimeConflictForAppointment(ctx, params)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar conflito com bloqueio: %w", err)
	}

	return hasConflict, nil
}

// CheckMinimumIntervalConflict verifica se há conflito de intervalo mínimo.
func (r *AppointmentRepository) CheckMinimumIntervalConflict(
	ctx context.Context,
	tenantID string,
	unitID string,
	professionalID string,
	startTime, endTime time.Time,
	excludeAppointmentID string,
	intervalMinutes int,
) (bool, error) {
	excludeID := pgtype.UUID{}
	if excludeAppointmentID != "" {
		excludeID = uuidStringToPgtype(excludeAppointmentID)
	}

	params := db.CheckMinimumIntervalConflictParams{
		TenantID:        uuidStringToPgtype(tenantID),
		UnitID:          uuidStringToPgtype(unitID),
		ProfessionalID:  uuidStringToPgtype(professionalID),
		ExcludeID:       excludeID,
		StartTime:       timestampToTimestamptz(startTime),
		EndTime:         timestampToTimestamptz(endTime),
		IntervalMinutes: int32(intervalMinutes),
	}

	hasConflict, err := r.queries.CheckMinimumIntervalConflict(ctx, params)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar intervalo mínimo: %w", err)
	}

	return hasConflict, nil
}

// CountByStatus conta agendamentos por status.
func (r *AppointmentRepository) CountByStatus(ctx context.Context, tenantID, unitID string, status valueobject.AppointmentStatus) (int64, error) {
	params := db.CountAppointmentsByStatusParams{
		TenantID: uuidStringToPgtype(tenantID),
		UnitID:   uuidStringToPgtype(unitID),
		Status:   status.String(),
	}

	count, err := r.queries.CountAppointmentsByStatus(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("erro ao contar agendamentos: %w", err)
	}

	return count, nil
}

// GetDailyStats retorna estatísticas diárias.
func (r *AppointmentRepository) GetDailyStats(ctx context.Context, tenantID, unitID string, date time.Time) (*port.AppointmentDailyStats, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	params := db.GetDailyAppointmentStatsParams{
		TenantID:    uuidStringToPgtype(tenantID),
		UnitID:      uuidStringToPgtype(unitID),
		StartTime:   timestampToTimestamptz(startOfDay),
		StartTime_2: timestampToTimestamptz(endOfDay),
	}

	row, err := r.queries.GetDailyAppointmentStats(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter estatísticas diárias: %w", err)
	}

	// Converter TotalRevenue para Money
	var totalRevenue valueobject.Money
	if row.TotalRevenue != nil {
		switch v := row.TotalRevenue.(type) {
		case decimal.Decimal:
			totalRevenue = valueobject.NewMoneyFromDecimal(v)
		case float64:
			totalRevenue = valueobject.NewMoneyFromFloat(v)
		default:
			totalRevenue = valueobject.Zero()
		}
	}

	return &port.AppointmentDailyStats{
		TotalAppointments: row.TotalAppointments,
		CompletedCount:    row.CompletedCount,
		CanceledCount:     row.CanceledCount,
		NoShowCount:       row.NoShowCount,
		TotalRevenue:      totalRevenue,
	}, nil
}

// === Métodos de conversão ===

func (r *AppointmentRepository) rowToDomain(row *db.GetAppointmentByIDRow, services []db.GetAppointmentServicesRow) *entity.Appointment {
	// Converter serviços
	domainServices := make([]entity.AppointmentService, 0, len(services))
	for _, svc := range services {
		domainServices = append(domainServices, entity.AppointmentService{
			AppointmentID:     pgUUIDToString(svc.AppointmentID),
			ServiceID:         pgUUIDToString(svc.ServiceID),
			PriceAtBooking:    valueobject.NewMoneyFromDecimal(svc.PriceAtBooking),
			DurationAtBooking: int(svc.DurationAtBooking),
			CreatedAt:         timestamptzToTime(svc.CreatedAt),
			ServiceName:       svc.ServiceName,
		})
	}

	status, _ := valueobject.ParseAppointmentStatus(row.Status)

	return &entity.Appointment{
		ID:                    pgUUIDToString(row.ID),
		TenantID:              pgtypeToEntityUUID(row.TenantID),
		ProfessionalID:        pgUUIDToString(row.ProfessionalID),
		CustomerID:            pgUUIDToString(row.CustomerID),
		StartTime:             timestamptzToTime(row.StartTime),
		EndTime:               timestamptzToTime(row.EndTime),
		CheckedInAt:           pgTimestamptzToTimePtr(row.CheckedInAt),
		StartedAt:             pgTimestamptzToTimePtr(row.StartedAt),
		FinishedAt:            pgTimestamptzToTimePtr(row.FinishedAt),
		Status:                status,
		TotalPrice:            valueobject.NewMoneyFromDecimal(row.TotalPrice),
		Notes:                 pgTextToStr(row.Notes),
		CanceledReason:        pgTextToStr(row.CanceledReason),
		GoogleCalendarEventID: pgTextToStr(row.GoogleCalendarEventID),
		CommandID:             pgUUIDPtrToString(row.CommandID),
		Services:              domainServices,
		ProfessionalName:      row.ProfessionalName,
		CustomerName:          row.CustomerName,
		CustomerPhone:         row.CustomerPhone,
		CreatedAt:             timestamptzToTime(row.CreatedAt),
		UpdatedAt:             timestamptzToTime(row.UpdatedAt),
	}
}

func (r *AppointmentRepository) listRowToDomain(row *db.ListAppointmentsRow) *entity.Appointment {
	status, _ := valueobject.ParseAppointmentStatus(row.Status)

	return &entity.Appointment{
		ID:                    pgUUIDToString(row.ID),
		TenantID:              pgtypeToEntityUUID(row.TenantID),
		ProfessionalID:        pgUUIDToString(row.ProfessionalID),
		CustomerID:            pgUUIDToString(row.CustomerID),
		StartTime:             timestamptzToTime(row.StartTime),
		EndTime:               timestamptzToTime(row.EndTime),
		CheckedInAt:           pgTimestamptzToTimePtr(row.CheckedInAt),
		StartedAt:             pgTimestamptzToTimePtr(row.StartedAt),
		FinishedAt:            pgTimestamptzToTimePtr(row.FinishedAt),
		Status:                status,
		TotalPrice:            valueobject.NewMoneyFromDecimal(row.TotalPrice),
		Notes:                 pgTextToStr(row.Notes),
		CanceledReason:        pgTextToStr(row.CanceledReason),
		GoogleCalendarEventID: pgTextToStr(row.GoogleCalendarEventID),
		CommandID:             pgUUIDPtrToString(row.CommandID),
		ProfessionalName:      row.ProfessionalName,
		CustomerName:          row.CustomerName,
		CustomerPhone:         row.CustomerPhone,
		CreatedAt:             timestamptzToTime(row.CreatedAt),
		UpdatedAt:             timestamptzToTime(row.UpdatedAt),
	}
}

func (r *AppointmentRepository) professionalRangeRowToDomain(row *db.ListAppointmentsByProfessionalAndDateRangeRow) *entity.Appointment {
	status, _ := valueobject.ParseAppointmentStatus(row.Status)

	return &entity.Appointment{
		ID:                    pgUUIDToString(row.ID),
		TenantID:              pgtypeToEntityUUID(row.TenantID),
		ProfessionalID:        pgUUIDToString(row.ProfessionalID),
		CustomerID:            pgUUIDToString(row.CustomerID),
		StartTime:             timestamptzToTime(row.StartTime),
		EndTime:               timestamptzToTime(row.EndTime),
		CheckedInAt:           pgTimestamptzToTimePtr(row.CheckedInAt),
		StartedAt:             pgTimestamptzToTimePtr(row.StartedAt),
		FinishedAt:            pgTimestamptzToTimePtr(row.FinishedAt),
		Status:                status,
		TotalPrice:            valueobject.NewMoneyFromDecimal(row.TotalPrice),
		Notes:                 pgTextToStr(row.Notes),
		CanceledReason:        pgTextToStr(row.CanceledReason),
		GoogleCalendarEventID: pgTextToStr(row.GoogleCalendarEventID),
		CommandID:             pgUUIDPtrToString(row.CommandID),
		ProfessionalName:      row.ProfessionalName,
		CustomerName:          row.CustomerName,
		CustomerPhone:         row.CustomerPhone,
		CreatedAt:             timestamptzToTime(row.CreatedAt),
		UpdatedAt:             timestamptzToTime(row.UpdatedAt),
	}
}

func (r *AppointmentRepository) customerRowToDomain(row *db.ListAppointmentsByCustomerRow) *entity.Appointment {
	status, _ := valueobject.ParseAppointmentStatus(row.Status)

	return &entity.Appointment{
		ID:                    pgUUIDToString(row.ID),
		TenantID:              pgtypeToEntityUUID(row.TenantID),
		ProfessionalID:        pgUUIDToString(row.ProfessionalID),
		CustomerID:            pgUUIDToString(row.CustomerID),
		StartTime:             timestamptzToTime(row.StartTime),
		EndTime:               timestamptzToTime(row.EndTime),
		CheckedInAt:           pgTimestamptzToTimePtr(row.CheckedInAt),
		StartedAt:             pgTimestamptzToTimePtr(row.StartedAt),
		FinishedAt:            pgTimestamptzToTimePtr(row.FinishedAt),
		Status:                status,
		TotalPrice:            valueobject.NewMoneyFromDecimal(row.TotalPrice),
		Notes:                 pgTextToStr(row.Notes),
		CanceledReason:        pgTextToStr(row.CanceledReason),
		GoogleCalendarEventID: pgTextToStr(row.GoogleCalendarEventID),
		CommandID:             pgUUIDPtrToString(row.CommandID),
		ProfessionalName:      row.ProfessionalName,
		CustomerName:          row.CustomerName,
		CustomerPhone:         row.CustomerPhone,
		CreatedAt:             timestamptzToTime(row.CreatedAt),
		UpdatedAt:             timestamptzToTime(row.UpdatedAt),
	}
}
