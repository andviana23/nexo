package appointment

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
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
