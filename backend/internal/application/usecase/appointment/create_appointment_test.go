package appointment

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Tenant de teste fixo para manter consistência
var testTenantID = "11111111-1111-1111-1111-111111111111"

func TestCreateAppointmentUseCase_Execute(t *testing.T) {
	logger := zap.NewNop()

	t.Run("should create appointment successfully", func(t *testing.T) {
		// Arrange
		mockRepo := &MockAppointmentRepository{}
		mockCommandRepo := &MockCommandRepository{}
		mockProfReader := &MockProfessionalReader{}
		mockCustReader := &MockCustomerReader{}
		mockSvcReader := &MockServiceReader{
			FindByIDsFn: func(ctx context.Context, tenantID string, serviceIDs []string) ([]*port.ServiceInfo, error) {
				return []*port.ServiceInfo{
					{ID: "svc-1", Name: "Corte", Price: valueobject.NewMoneyFromFloat(50.0), Duration: 30, Active: true},
				}, nil
			},
		}

		uc := NewCreateAppointmentUseCase(mockRepo, mockCommandRepo, mockSvcReader, mockProfReader, mockCustReader, logger)

		input := CreateAppointmentInput{
			TenantID:       testTenantID,
			ProfessionalID: "prof-123",
			CustomerID:     "cust-123",
			StartTime:      time.Now().Add(24 * time.Hour),
			ServiceIDs:     []string{"svc-1"},
			Notes:          "Test appointment",
		}

		// Act
		result, err := uc.Execute(context.Background(), input)

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
		expectedTenantUUID, _ := uuid.Parse(input.TenantID)
		if result.TenantID != expectedTenantUUID {
			t.Errorf("expected tenant_id %s, got %s", input.TenantID, result.TenantID.String())
		}
		if result.ProfessionalID != input.ProfessionalID {
			t.Errorf("expected professional_id %s, got %s", input.ProfessionalID, result.ProfessionalID)
		}
		if result.CustomerID != input.CustomerID {
			t.Errorf("expected customer_id %s, got %s", input.CustomerID, result.CustomerID)
		}
		if mockRepo.CreateCalls != 1 {
			t.Errorf("expected Create to be called once, got %d", mockRepo.CreateCalls)
		}
	})

	t.Run("should fail without tenant_id", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{}
		mockCommandRepo := &MockCommandRepository{}
		mockProfReader := &MockProfessionalReader{}
		mockCustReader := &MockCustomerReader{}
		mockSvcReader := &MockServiceReader{}

		uc := NewCreateAppointmentUseCase(mockRepo, mockCommandRepo, mockSvcReader, mockProfReader, mockCustReader, logger)

		input := CreateAppointmentInput{
			TenantID:       "",
			ProfessionalID: "prof-123",
			CustomerID:     "cust-123",
			StartTime:      time.Now().Add(24 * time.Hour),
			ServiceIDs:     []string{"svc-1"},
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrTenantIDRequired) {
			t.Errorf("expected ErrTenantIDRequired, got %v", err)
		}
	})

	t.Run("should fail without professional_id", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{}
		mockCommandRepo := &MockCommandRepository{}
		mockProfReader := &MockProfessionalReader{}
		mockCustReader := &MockCustomerReader{}
		mockSvcReader := &MockServiceReader{}

		uc := NewCreateAppointmentUseCase(mockRepo, mockCommandRepo, mockSvcReader, mockProfReader, mockCustReader, logger)

		input := CreateAppointmentInput{
			TenantID:       testTenantID,
			ProfessionalID: "",
			CustomerID:     "cust-123",
			StartTime:      time.Now().Add(24 * time.Hour),
			ServiceIDs:     []string{"svc-1"},
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAppointmentProfessionalRequired) {
			t.Errorf("expected ErrAppointmentProfessionalRequired, got %v", err)
		}
	})

	t.Run("should fail without customer_id", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{}
		mockCommandRepo := &MockCommandRepository{}
		mockProfReader := &MockProfessionalReader{}
		mockCustReader := &MockCustomerReader{}
		mockSvcReader := &MockServiceReader{}

		uc := NewCreateAppointmentUseCase(mockRepo, mockCommandRepo, mockSvcReader, mockProfReader, mockCustReader, logger)

		input := CreateAppointmentInput{
			TenantID:       testTenantID,
			ProfessionalID: "prof-123",
			CustomerID:     "",
			StartTime:      time.Now().Add(24 * time.Hour),
			ServiceIDs:     []string{"svc-1"},
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAppointmentCustomerRequired) {
			t.Errorf("expected ErrAppointmentCustomerRequired, got %v", err)
		}
	})

	t.Run("should fail without start_time", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{}
		mockCommandRepo := &MockCommandRepository{}
		mockProfReader := &MockProfessionalReader{}
		mockCustReader := &MockCustomerReader{}
		mockSvcReader := &MockServiceReader{}

		uc := NewCreateAppointmentUseCase(mockRepo, mockCommandRepo, mockSvcReader, mockProfReader, mockCustReader, logger)

		input := CreateAppointmentInput{
			TenantID:       testTenantID,
			ProfessionalID: "prof-123",
			CustomerID:     "cust-123",
			StartTime:      time.Time{}, // zero
			ServiceIDs:     []string{"svc-1"},
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAppointmentStartTimeRequired) {
			t.Errorf("expected ErrAppointmentStartTimeRequired, got %v", err)
		}
	})

	t.Run("should fail without services", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{}
		mockCommandRepo := &MockCommandRepository{}
		mockProfReader := &MockProfessionalReader{}
		mockCustReader := &MockCustomerReader{}
		mockSvcReader := &MockServiceReader{}

		uc := NewCreateAppointmentUseCase(mockRepo, mockCommandRepo, mockSvcReader, mockProfReader, mockCustReader, logger)

		input := CreateAppointmentInput{
			TenantID:       testTenantID,
			ProfessionalID: "prof-123",
			CustomerID:     "cust-123",
			StartTime:      time.Now().Add(24 * time.Hour),
			ServiceIDs:     []string{},
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAppointmentServicesRequired) {
			t.Errorf("expected ErrAppointmentServicesRequired, got %v", err)
		}
	})

	t.Run("should fail when professional not found", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{}
		mockCommandRepo := &MockCommandRepository{}
		mockProfReader := &MockProfessionalReader{
			ExistsFn: func(ctx context.Context, tenantID, professionalID string) (bool, error) {
				return false, nil
			},
		}
		mockCustReader := &MockCustomerReader{}
		mockSvcReader := &MockServiceReader{}

		uc := NewCreateAppointmentUseCase(mockRepo, mockCommandRepo, mockSvcReader, mockProfReader, mockCustReader, logger)

		input := CreateAppointmentInput{
			TenantID:       testTenantID,
			ProfessionalID: "prof-not-exists",
			CustomerID:     "cust-123",
			StartTime:      time.Now().Add(24 * time.Hour),
			ServiceIDs:     []string{"svc-1"},
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAppointmentProfessionalNotFound) {
			t.Errorf("expected ErrAppointmentProfessionalNotFound, got %v", err)
		}
	})

	t.Run("should fail when customer not found", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{}
		mockCommandRepo := &MockCommandRepository{}
		mockProfReader := &MockProfessionalReader{}
		mockCustReader := &MockCustomerReader{
			ExistsFn: func(ctx context.Context, tenantID, customerID string) (bool, error) {
				return false, nil
			},
		}
		mockSvcReader := &MockServiceReader{}

		uc := NewCreateAppointmentUseCase(mockRepo, mockCommandRepo, mockSvcReader, mockProfReader, mockCustReader, logger)

		input := CreateAppointmentInput{
			TenantID:       testTenantID,
			ProfessionalID: "prof-123",
			CustomerID:     "cust-not-exists",
			StartTime:      time.Now().Add(24 * time.Hour),
			ServiceIDs:     []string{"svc-1"},
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAppointmentCustomerNotFound) {
			t.Errorf("expected ErrAppointmentCustomerNotFound, got %v", err)
		}
	})

	t.Run("should fail when service not found", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{}
		mockCommandRepo := &MockCommandRepository{}
		mockProfReader := &MockProfessionalReader{}
		mockCustReader := &MockCustomerReader{}
		mockSvcReader := &MockServiceReader{
			FindByIDsFn: func(ctx context.Context, tenantID string, serviceIDs []string) ([]*port.ServiceInfo, error) {
				// Return fewer services than requested
				return []*port.ServiceInfo{}, nil
			},
		}

		uc := NewCreateAppointmentUseCase(mockRepo, mockCommandRepo, mockSvcReader, mockProfReader, mockCustReader, logger)

		input := CreateAppointmentInput{
			TenantID:       testTenantID,
			ProfessionalID: "prof-123",
			CustomerID:     "cust-123",
			StartTime:      time.Now().Add(24 * time.Hour),
			ServiceIDs:     []string{"svc-not-exists"},
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAppointmentServiceNotFound) {
			t.Errorf("expected ErrAppointmentServiceNotFound, got %v", err)
		}
	})

	t.Run("should fail on time conflict", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{
			CheckConflictFn: func(ctx context.Context, tenantID, professionalID string, startTime, endTime time.Time, excludeAppointmentID string) (bool, error) {
				return true, nil // Has conflict
			},
		}
		mockCommandRepo := &MockCommandRepository{}
		mockProfReader := &MockProfessionalReader{}
		mockCustReader := &MockCustomerReader{}
		mockSvcReader := &MockServiceReader{
			FindByIDsFn: func(ctx context.Context, tenantID string, serviceIDs []string) ([]*port.ServiceInfo, error) {
				return []*port.ServiceInfo{
					{ID: "svc-1", Name: "Corte", Price: valueobject.NewMoneyFromFloat(50.0), Duration: 30, Active: true},
				}, nil
			},
		}

		uc := NewCreateAppointmentUseCase(mockRepo, mockCommandRepo, mockSvcReader, mockProfReader, mockCustReader, logger)

		input := CreateAppointmentInput{
			TenantID:       testTenantID,
			ProfessionalID: "prof-123",
			CustomerID:     "cust-123",
			StartTime:      time.Now().Add(24 * time.Hour),
			ServiceIDs:     []string{"svc-1"},
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAppointmentConflict) {
			t.Errorf("expected ErrAppointmentConflict, got %v", err)
		}
	})

	t.Run("should calculate end_time and total_price correctly", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{}
		mockCommandRepo := &MockCommandRepository{}
		mockProfReader := &MockProfessionalReader{}
		mockCustReader := &MockCustomerReader{}
		mockSvcReader := &MockServiceReader{
			FindByIDsFn: func(ctx context.Context, tenantID string, serviceIDs []string) ([]*port.ServiceInfo, error) {
				return []*port.ServiceInfo{
					{ID: "svc-1", Name: "Corte", Price: valueobject.NewMoneyFromFloat(50.0), Duration: 30, Active: true},
					{ID: "svc-2", Name: "Barba", Price: valueobject.NewMoneyFromFloat(30.0), Duration: 20, Active: true},
				}, nil
			},
		}

		uc := NewCreateAppointmentUseCase(mockRepo, mockCommandRepo, mockSvcReader, mockProfReader, mockCustReader, logger)

		startTime := time.Now().Add(24 * time.Hour)
		input := CreateAppointmentInput{
			TenantID:       testTenantID,
			ProfessionalID: "prof-123",
			CustomerID:     "cust-123",
			StartTime:      startTime,
			ServiceIDs:     []string{"svc-1", "svc-2"},
		}

		result, err := uc.Execute(context.Background(), input)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Check end_time = start_time + 30min + 20min = start_time + 50min
		expectedEndTime := startTime.Add(50 * time.Minute)
		if !result.EndTime.Equal(expectedEndTime) {
			t.Errorf("expected end_time %v, got %v", expectedEndTime, result.EndTime)
		}

		// Check total_price = 50 + 30 = 80
		expectedPrice := valueobject.NewMoneyFromFloat(80.0)
		if !result.TotalPrice.Equals(expectedPrice) {
			t.Errorf("expected total_price %s, got %s", expectedPrice.String(), result.TotalPrice.String())
		}
	})

	t.Run("should fail when service is inactive", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{}
		mockCommandRepo := &MockCommandRepository{}
		mockProfReader := &MockProfessionalReader{}
		mockCustReader := &MockCustomerReader{}
		mockSvcReader := &MockServiceReader{
			FindByIDsFn: func(ctx context.Context, tenantID string, serviceIDs []string) ([]*port.ServiceInfo, error) {
				return []*port.ServiceInfo{
					{ID: "svc-1", Name: "Corte Inativo", Price: valueobject.NewMoneyFromFloat(50.0), Duration: 30, Active: false},
				}, nil
			},
		}

		uc := NewCreateAppointmentUseCase(mockRepo, mockCommandRepo, mockSvcReader, mockProfReader, mockCustReader, logger)

		input := CreateAppointmentInput{
			TenantID:       testTenantID,
			ProfessionalID: "prof-123",
			CustomerID:     "cust-123",
			StartTime:      time.Now().Add(24 * time.Hour),
			ServiceIDs:     []string{"svc-1"},
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error for inactive service, got nil")
		}
	})
}

func TestListAppointmentsUseCase_Execute(t *testing.T) {
	logger := zap.NewNop()

	t.Run("should list appointments with default pagination", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{
			ListFn: func(ctx context.Context, tenantID string, filter port.AppointmentFilter) ([]*entity.Appointment, int64, error) {
				return []*entity.Appointment{}, 0, nil
			},
		}

		uc := NewListAppointmentsUseCase(mockRepo, logger)

		input := ListAppointmentsInput{
			TenantID: "tenant-123",
		}

		result, err := uc.Execute(context.Background(), input)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.Page != 1 {
			t.Errorf("expected page 1, got %d", result.Page)
		}
		if result.PageSize != 20 {
			t.Errorf("expected page_size 20, got %d", result.PageSize)
		}
	})

	t.Run("should fail without tenant_id", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{}
		uc := NewListAppointmentsUseCase(mockRepo, logger)

		input := ListAppointmentsInput{
			TenantID: "",
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrTenantIDRequired) {
			t.Errorf("expected ErrTenantIDRequired, got %v", err)
		}
	})

	t.Run("should respect page size limit", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{
			ListFn: func(ctx context.Context, tenantID string, filter port.AppointmentFilter) ([]*entity.Appointment, int64, error) {
				if filter.PageSize > 100 {
					t.Error("page_size should be capped at 100")
				}
				return []*entity.Appointment{}, 0, nil
			},
		}

		uc := NewListAppointmentsUseCase(mockRepo, logger)

		input := ListAppointmentsInput{
			TenantID: "tenant-123",
			PageSize: 500, // Requesting more than allowed
		}

		result, err := uc.Execute(context.Background(), input)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.PageSize != 20 { // Should default to 20 since 500 > 100
			t.Errorf("expected page_size to be capped, got %d", result.PageSize)
		}
	})

	t.Run("should filter by date range", func(t *testing.T) {
		startDate := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2025, 12, 1, 23, 59, 59, 0, time.UTC)

		mockRepo := &MockAppointmentRepository{
			ListFn: func(ctx context.Context, tenantID string, filter port.AppointmentFilter) ([]*entity.Appointment, int64, error) {
				if filter.StartDate.IsZero() {
					t.Error("expected start_date to be set")
				}
				if filter.EndDate.IsZero() {
					t.Error("expected end_date to be set")
				}
				if !filter.StartDate.Equal(startDate) {
					t.Errorf("expected start_date %v, got %v", startDate, filter.StartDate)
				}
				if !filter.EndDate.Equal(endDate) {
					t.Errorf("expected end_date %v, got %v", endDate, filter.EndDate)
				}
				return []*entity.Appointment{}, 0, nil
			},
		}

		uc := NewListAppointmentsUseCase(mockRepo, logger)

		input := ListAppointmentsInput{
			TenantID:  "tenant-123",
			StartDate: startDate,
			EndDate:   endDate,
		}

		_, err := uc.Execute(context.Background(), input)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("should filter by multiple statuses", func(t *testing.T) {
		statuses := []valueobject.AppointmentStatus{
			valueobject.AppointmentStatusCreated,
			valueobject.AppointmentStatusConfirmed,
			valueobject.AppointmentStatusAwaitingPayment,
		}

		mockRepo := &MockAppointmentRepository{
			ListFn: func(ctx context.Context, tenantID string, filter port.AppointmentFilter) ([]*entity.Appointment, int64, error) {
				if len(filter.Statuses) != 3 {
					t.Errorf("expected 3 statuses, got %d", len(filter.Statuses))
				}
				// Verificar se todos os status esperados estão presentes
				for _, expected := range statuses {
					found := false
					for _, actual := range filter.Statuses {
						if actual == expected {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("expected status %s to be in filter", expected)
					}
				}
				return []*entity.Appointment{}, 0, nil
			},
		}

		uc := NewListAppointmentsUseCase(mockRepo, logger)

		input := ListAppointmentsInput{
			TenantID: "tenant-123",
			Statuses: statuses,
		}

		_, err := uc.Execute(context.Background(), input)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("should filter by single status as array", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{
			ListFn: func(ctx context.Context, tenantID string, filter port.AppointmentFilter) ([]*entity.Appointment, int64, error) {
				if len(filter.Statuses) != 1 {
					t.Errorf("expected 1 status, got %d", len(filter.Statuses))
				}
				if filter.Statuses[0] != valueobject.AppointmentStatusAwaitingPayment {
					t.Errorf("expected status AWAITING_PAYMENT, got %s", filter.Statuses[0])
				}
				return []*entity.Appointment{}, 0, nil
			},
		}

		uc := NewListAppointmentsUseCase(mockRepo, logger)

		input := ListAppointmentsInput{
			TenantID: "tenant-123",
			Statuses: []valueobject.AppointmentStatus{valueobject.AppointmentStatusAwaitingPayment},
		}

		_, err := uc.Execute(context.Background(), input)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}
