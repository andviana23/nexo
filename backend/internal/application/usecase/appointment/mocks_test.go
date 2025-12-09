package appointment

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

// ============================================================================
// Mock AppointmentRepository
// ============================================================================

// MockAppointmentRepository implementa port.AppointmentRepository para testes.
type MockAppointmentRepository struct {
	CreateFn                         func(ctx context.Context, appointment *entity.Appointment) error
	FindByIDFn                       func(ctx context.Context, tenantID, id string) (*entity.Appointment, error)
	UpdateFn                         func(ctx context.Context, appointment *entity.Appointment) error
	DeleteFn                         func(ctx context.Context, tenantID, id string) error
	ListFn                           func(ctx context.Context, tenantID string, filter port.AppointmentFilter) ([]*entity.Appointment, int64, error)
	ListByProfessionalAndDateRangeFn func(ctx context.Context, tenantID, professionalID string, startDate, endDate time.Time) ([]*entity.Appointment, error)
	ListByCustomerFn                 func(ctx context.Context, tenantID, customerID string) ([]*entity.Appointment, error)
	CheckConflictFn                  func(ctx context.Context, tenantID, professionalID string, startTime, endTime time.Time, excludeAppointmentID string) (bool, error)
	CheckBlockedTimeConflictFn       func(ctx context.Context, tenantID, professionalID string, startTime, endTime time.Time) (bool, error)
	CheckMinimumIntervalConflictFn   func(ctx context.Context, tenantID, professionalID string, startTime, endTime time.Time, excludeAppointmentID string, intervalMinutes int) (bool, error)
	CountByStatusFn                  func(ctx context.Context, tenantID string, status valueobject.AppointmentStatus) (int64, error)
	GetDailyStatsFn                  func(ctx context.Context, tenantID string, date time.Time) (*port.AppointmentDailyStats, error)

	// Tracking calls
	CreateCalls   int
	UpdateCalls   int
	FindByIDCalls int
}

func (m *MockAppointmentRepository) Create(ctx context.Context, appointment *entity.Appointment) error {
	m.CreateCalls++
	if m.CreateFn != nil {
		return m.CreateFn(ctx, appointment)
	}
	return nil
}

func (m *MockAppointmentRepository) FindByID(ctx context.Context, tenantID, id string) (*entity.Appointment, error) {
	m.FindByIDCalls++
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, tenantID, id)
	}
	return nil, nil
}

func (m *MockAppointmentRepository) Update(ctx context.Context, appointment *entity.Appointment) error {
	m.UpdateCalls++
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, appointment)
	}
	return nil
}

func (m *MockAppointmentRepository) Delete(ctx context.Context, tenantID, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, tenantID, id)
	}
	return nil
}

func (m *MockAppointmentRepository) List(ctx context.Context, tenantID string, filter port.AppointmentFilter) ([]*entity.Appointment, int64, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, tenantID, filter)
	}
	return nil, 0, nil
}

func (m *MockAppointmentRepository) ListByProfessionalAndDateRange(ctx context.Context, tenantID, professionalID string, startDate, endDate time.Time) ([]*entity.Appointment, error) {
	if m.ListByProfessionalAndDateRangeFn != nil {
		return m.ListByProfessionalAndDateRangeFn(ctx, tenantID, professionalID, startDate, endDate)
	}
	return nil, nil
}

func (m *MockAppointmentRepository) ListByCustomer(ctx context.Context, tenantID, customerID string) ([]*entity.Appointment, error) {
	if m.ListByCustomerFn != nil {
		return m.ListByCustomerFn(ctx, tenantID, customerID)
	}
	return nil, nil
}

func (m *MockAppointmentRepository) CheckConflict(ctx context.Context, tenantID, professionalID string, startTime, endTime time.Time, excludeAppointmentID string) (bool, error) {
	if m.CheckConflictFn != nil {
		return m.CheckConflictFn(ctx, tenantID, professionalID, startTime, endTime, excludeAppointmentID)
	}
	return false, nil
}

func (m *MockAppointmentRepository) CheckBlockedTimeConflict(ctx context.Context, tenantID, professionalID string, startTime, endTime time.Time) (bool, error) {
	if m.CheckBlockedTimeConflictFn != nil {
		return m.CheckBlockedTimeConflictFn(ctx, tenantID, professionalID, startTime, endTime)
	}
	return false, nil
}

func (m *MockAppointmentRepository) CheckMinimumIntervalConflict(ctx context.Context, tenantID, professionalID string, startTime, endTime time.Time, excludeAppointmentID string, intervalMinutes int) (bool, error) {
	if m.CheckMinimumIntervalConflictFn != nil {
		return m.CheckMinimumIntervalConflictFn(ctx, tenantID, professionalID, startTime, endTime, excludeAppointmentID, intervalMinutes)
	}
	return false, nil
}

func (m *MockAppointmentRepository) CountByStatus(ctx context.Context, tenantID string, status valueobject.AppointmentStatus) (int64, error) {
	if m.CountByStatusFn != nil {
		return m.CountByStatusFn(ctx, tenantID, status)
	}
	return 0, nil
}

func (m *MockAppointmentRepository) GetDailyStats(ctx context.Context, tenantID string, date time.Time) (*port.AppointmentDailyStats, error) {
	if m.GetDailyStatsFn != nil {
		return m.GetDailyStatsFn(ctx, tenantID, date)
	}
	return nil, nil
}

// ============================================================================
// Mock ProfessionalReader
// ============================================================================

// MockProfessionalReader implementa port.ProfessionalReader para testes.
type MockProfessionalReader struct {
	ExistsFn     func(ctx context.Context, tenantID, professionalID string) (bool, error)
	FindByIDFn   func(ctx context.Context, tenantID, professionalID string) (*port.ProfessionalInfo, error)
	ListActiveFn func(ctx context.Context, tenantID string) ([]*port.ProfessionalInfo, error)
}

func (m *MockProfessionalReader) Exists(ctx context.Context, tenantID, professionalID string) (bool, error) {
	if m.ExistsFn != nil {
		return m.ExistsFn(ctx, tenantID, professionalID)
	}
	return true, nil
}

func (m *MockProfessionalReader) FindByID(ctx context.Context, tenantID, professionalID string) (*port.ProfessionalInfo, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, tenantID, professionalID)
	}
	return &port.ProfessionalInfo{ID: professionalID, Name: "Test Professional"}, nil
}

func (m *MockProfessionalReader) ListActive(ctx context.Context, tenantID string) ([]*port.ProfessionalInfo, error) {
	if m.ListActiveFn != nil {
		return m.ListActiveFn(ctx, tenantID)
	}
	return nil, nil
}

// ============================================================================
// Mock CustomerReader
// ============================================================================

// MockCustomerReader implementa port.CustomerReader para testes.
type MockCustomerReader struct {
	ExistsFn   func(ctx context.Context, tenantID, customerID string) (bool, error)
	FindByIDFn func(ctx context.Context, tenantID, customerID string) (*port.CustomerInfo, error)
}

func (m *MockCustomerReader) Exists(ctx context.Context, tenantID, customerID string) (bool, error) {
	if m.ExistsFn != nil {
		return m.ExistsFn(ctx, tenantID, customerID)
	}
	return true, nil
}

func (m *MockCustomerReader) FindByID(ctx context.Context, tenantID, customerID string) (*port.CustomerInfo, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, tenantID, customerID)
	}
	return &port.CustomerInfo{ID: customerID, Name: "Test Customer"}, nil
}

// ============================================================================
// Mock ServiceReader
// ============================================================================

// MockServiceReader implementa port.ServiceReader para testes.
type MockServiceReader struct {
	ExistsFn    func(ctx context.Context, tenantID, serviceID string) (bool, error)
	FindByIDFn  func(ctx context.Context, tenantID, serviceID string) (*port.ServiceInfo, error)
	FindByIDsFn func(ctx context.Context, tenantID string, serviceIDs []string) ([]*port.ServiceInfo, error)
}

func (m *MockServiceReader) Exists(ctx context.Context, tenantID, serviceID string) (bool, error) {
	if m.ExistsFn != nil {
		return m.ExistsFn(ctx, tenantID, serviceID)
	}
	return true, nil
}

func (m *MockServiceReader) FindByID(ctx context.Context, tenantID, serviceID string) (*port.ServiceInfo, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, tenantID, serviceID)
	}
	return &port.ServiceInfo{
		ID:       serviceID,
		Name:     "Test Service",
		Price:    valueobject.NewMoneyFromFloat(50.0),
		Duration: 30,
		Active:   true,
	}, nil
}

func (m *MockServiceReader) FindByIDs(ctx context.Context, tenantID string, serviceIDs []string) ([]*port.ServiceInfo, error) {
	if m.FindByIDsFn != nil {
		return m.FindByIDsFn(ctx, tenantID, serviceIDs)
	}
	services := make([]*port.ServiceInfo, 0, len(serviceIDs))
	for _, id := range serviceIDs {
		services = append(services, &port.ServiceInfo{
			ID:       id,
			Name:     "Service " + id,
			Price:    valueobject.NewMoneyFromFloat(50.0),
			Duration: 30,
			Active:   true,
		})
	}
	return services, nil
}

// ============================================================================
// Mock CommandRepository
// ============================================================================

// MockCommandRepository implementa port.CommandRepository para testes.
type MockCommandRepository struct {
	CreateFn              func(ctx context.Context, command *entity.Command) error
	FindByIDFn            func(ctx context.Context, commandID, tenantID uuid.UUID) (*entity.Command, error)
	FindByAppointmentIDFn func(ctx context.Context, appointmentID, tenantID uuid.UUID) (*entity.Command, error)
	UpdateFn              func(ctx context.Context, command *entity.Command) error
	DeleteFn              func(ctx context.Context, commandID, tenantID uuid.UUID) error
	ListFn                func(ctx context.Context, tenantID uuid.UUID, filters port.CommandFilters) ([]*entity.Command, error)
	AddItemFn             func(ctx context.Context, item *entity.CommandItem) error
	UpdateItemFn          func(ctx context.Context, item *entity.CommandItem) error
	RemoveItemFn          func(ctx context.Context, itemID, tenantID uuid.UUID) error
	GetItemsFn            func(ctx context.Context, commandID, tenantID uuid.UUID) ([]entity.CommandItem, error)
	AddPaymentFn          func(ctx context.Context, payment *entity.CommandPayment) error
	RemovePaymentFn       func(ctx context.Context, paymentID, tenantID uuid.UUID) error
	GetPaymentsFn         func(ctx context.Context, commandID, tenantID uuid.UUID) ([]entity.CommandPayment, error)
}

func (m *MockCommandRepository) Create(ctx context.Context, command *entity.Command) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, command)
	}
	return nil
}

func (m *MockCommandRepository) FindByID(ctx context.Context, commandID, tenantID uuid.UUID) (*entity.Command, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, commandID, tenantID)
	}
	return nil, nil
}

func (m *MockCommandRepository) FindByAppointmentID(ctx context.Context, appointmentID, tenantID uuid.UUID) (*entity.Command, error) {
	if m.FindByAppointmentIDFn != nil {
		return m.FindByAppointmentIDFn(ctx, appointmentID, tenantID)
	}
	return nil, nil
}

func (m *MockCommandRepository) Update(ctx context.Context, command *entity.Command) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, command)
	}
	return nil
}

func (m *MockCommandRepository) Delete(ctx context.Context, commandID, tenantID uuid.UUID) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, commandID, tenantID)
	}
	return nil
}

func (m *MockCommandRepository) List(ctx context.Context, tenantID uuid.UUID, filters port.CommandFilters) ([]*entity.Command, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, tenantID, filters)
	}
	return nil, nil
}

func (m *MockCommandRepository) AddItem(ctx context.Context, item *entity.CommandItem) error {
	if m.AddItemFn != nil {
		return m.AddItemFn(ctx, item)
	}
	return nil
}

func (m *MockCommandRepository) UpdateItem(ctx context.Context, item *entity.CommandItem) error {
	if m.UpdateItemFn != nil {
		return m.UpdateItemFn(ctx, item)
	}
	return nil
}

func (m *MockCommandRepository) RemoveItem(ctx context.Context, itemID, tenantID uuid.UUID) error {
	if m.RemoveItemFn != nil {
		return m.RemoveItemFn(ctx, itemID, tenantID)
	}
	return nil
}

func (m *MockCommandRepository) GetItems(ctx context.Context, commandID, tenantID uuid.UUID) ([]entity.CommandItem, error) {
	if m.GetItemsFn != nil {
		return m.GetItemsFn(ctx, commandID, tenantID)
	}
	return nil, nil
}

func (m *MockCommandRepository) AddPayment(ctx context.Context, payment *entity.CommandPayment) error {
	if m.AddPaymentFn != nil {
		return m.AddPaymentFn(ctx, payment)
	}
	return nil
}

func (m *MockCommandRepository) RemovePayment(ctx context.Context, paymentID, tenantID uuid.UUID) error {
	if m.RemovePaymentFn != nil {
		return m.RemovePaymentFn(ctx, paymentID, tenantID)
	}
	return nil
}

func (m *MockCommandRepository) GetPayments(ctx context.Context, commandID, tenantID uuid.UUID) ([]entity.CommandPayment, error) {
	if m.GetPaymentsFn != nil {
		return m.GetPaymentsFn(ctx, commandID, tenantID)
	}
	return nil, nil
}
