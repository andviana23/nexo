package appointment

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Tenant de teste fixo para manter consistência
var testTenantUUID = uuid.MustParse("11111111-1111-1111-1111-111111111111")

func TestCancelAppointmentUseCase_Execute(t *testing.T) {
	logger := zap.NewNop()

	t.Run("should cancel appointment successfully", func(t *testing.T) {
		// Create a valid appointment in CREATED status
		services := []entity.AppointmentService{
			{ServiceID: "svc-1", ServiceName: "Corte", PriceAtBooking: valueobject.NewMoneyFromFloat(50.0), DurationAtBooking: 30},
		}
		appointment, _ := entity.NewAppointment(testTenantUUID, "prof-123", "cust-123", time.Now().Add(24*time.Hour), services)

		mockRepo := &MockAppointmentRepository{
			FindByIDFn: func(ctx context.Context, tenantID, id string) (*entity.Appointment, error) {
				return appointment, nil
			},
		}

		uc := NewCancelAppointmentUseCase(mockRepo, logger)

		input := CancelAppointmentInput{
			TenantID:      testTenantUUID.String(),
			AppointmentID: appointment.ID,
			Reason:        "Cliente solicitou cancelamento",
		}

		result, err := uc.Execute(context.Background(), input)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.Status != valueobject.AppointmentStatusCanceled {
			t.Errorf("expected status CANCELED, got %s", result.Status)
		}
		if mockRepo.UpdateCalls != 1 {
			t.Errorf("expected Update to be called once, got %d", mockRepo.UpdateCalls)
		}
	})

	t.Run("should fail without tenant_id", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{}
		uc := NewCancelAppointmentUseCase(mockRepo, logger)

		input := CancelAppointmentInput{
			TenantID:      "",
			AppointmentID: "appt-123",
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrTenantIDRequired) {
			t.Errorf("expected ErrTenantIDRequired, got %v", err)
		}
	})

	t.Run("should fail without appointment_id", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{}
		uc := NewCancelAppointmentUseCase(mockRepo, logger)

		input := CancelAppointmentInput{
			TenantID:      testTenantUUID.String(),
			AppointmentID: "",
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidID) {
			t.Errorf("expected ErrInvalidID, got %v", err)
		}
	})

	t.Run("should fail when appointment not found", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{
			FindByIDFn: func(ctx context.Context, tenantID, id string) (*entity.Appointment, error) {
				return nil, domain.ErrAppointmentNotFound
			},
		}

		uc := NewCancelAppointmentUseCase(mockRepo, logger)

		input := CancelAppointmentInput{
			TenantID:      testTenantUUID.String(),
			AppointmentID: "appt-not-exists",
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("should fail when appointment already completed", func(t *testing.T) {
		services := []entity.AppointmentService{
			{ServiceID: "svc-1", ServiceName: "Corte", PriceAtBooking: valueobject.NewMoneyFromFloat(50.0), DurationAtBooking: 30},
		}
		appointment, _ := entity.NewAppointment(testTenantUUID, "prof-123", "cust-123", time.Now().Add(24*time.Hour), services)
		appointment.Confirm()
		appointment.StartService()
		appointment.Complete()

		mockRepo := &MockAppointmentRepository{
			FindByIDFn: func(ctx context.Context, tenantID, id string) (*entity.Appointment, error) {
				return appointment, nil
			},
		}

		uc := NewCancelAppointmentUseCase(mockRepo, logger)

		input := CancelAppointmentInput{
			TenantID:      testTenantUUID.String(),
			AppointmentID: appointment.ID,
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error when canceling completed appointment, got nil")
		}
	})
}

func TestRescheduleAppointmentUseCase_Execute(t *testing.T) {
	logger := zap.NewNop()

	t.Run("should reschedule appointment successfully", func(t *testing.T) {
		services := []entity.AppointmentService{
			{ServiceID: "svc-1", ServiceName: "Corte", PriceAtBooking: valueobject.NewMoneyFromFloat(50.0), DurationAtBooking: 30},
		}
		originalTime := time.Now().Add(24 * time.Hour)
		appointment, _ := entity.NewAppointment(testTenantUUID, "prof-123", "cust-123", originalTime, services)

		mockRepo := &MockAppointmentRepository{
			FindByIDFn: func(ctx context.Context, tenantID, id string) (*entity.Appointment, error) {
				return appointment, nil
			},
			CheckConflictFn: func(ctx context.Context, tenantID, professionalID string, startTime, endTime time.Time, excludeAppointmentID string) (bool, error) {
				return false, nil // No conflict
			},
		}

		uc := NewRescheduleAppointmentUseCase(mockRepo, &MockProfessionalReader{}, logger)

		newTime := time.Now().Add(48 * time.Hour)
		input := RescheduleAppointmentInput{
			TenantID:      testTenantUUID.String(),
			AppointmentID: appointment.ID,
			NewStartTime:  newTime,
		}

		result, err := uc.Execute(context.Background(), input)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !result.StartTime.Equal(newTime) {
			t.Errorf("expected start_time %v, got %v", newTime, result.StartTime)
		}
		if mockRepo.UpdateCalls != 1 {
			t.Errorf("expected Update to be called once, got %d", mockRepo.UpdateCalls)
		}
	})

	t.Run("should fail on time conflict", func(t *testing.T) {
		services := []entity.AppointmentService{
			{ServiceID: "svc-1", ServiceName: "Corte", PriceAtBooking: valueobject.NewMoneyFromFloat(50.0), DurationAtBooking: 30},
		}
		appointment, _ := entity.NewAppointment(testTenantUUID, "prof-123", "cust-123", time.Now().Add(24*time.Hour), services)

		mockRepo := &MockAppointmentRepository{
			FindByIDFn: func(ctx context.Context, tenantID, id string) (*entity.Appointment, error) {
				return appointment, nil
			},
			CheckConflictFn: func(ctx context.Context, tenantID, professionalID string, startTime, endTime time.Time, excludeAppointmentID string) (bool, error) {
				return true, nil // Has conflict
			},
		}

		uc := NewRescheduleAppointmentUseCase(mockRepo, &MockProfessionalReader{}, logger)

		input := RescheduleAppointmentInput{
			TenantID:      testTenantUUID.String(),
			AppointmentID: appointment.ID,
			NewStartTime:  time.Now().Add(48 * time.Hour),
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected conflict error, got nil")
		}
		if !errors.Is(err, domain.ErrAppointmentConflict) {
			t.Errorf("expected ErrAppointmentConflict, got %v", err)
		}
	})

	t.Run("should fail without tenant_id", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{}
		uc := NewRescheduleAppointmentUseCase(mockRepo, &MockProfessionalReader{}, logger)

		input := RescheduleAppointmentInput{
			TenantID:      "",
			AppointmentID: "appt-123",
			NewStartTime:  time.Now().Add(24 * time.Hour),
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrTenantIDRequired) {
			t.Errorf("expected ErrTenantIDRequired, got %v", err)
		}
	})
}

func TestUpdateAppointmentStatusUseCase_Execute(t *testing.T) {
	logger := zap.NewNop()
	mockCommandRepo := &MockCommandRepository{}

	t.Run("should confirm appointment", func(t *testing.T) {
		services := []entity.AppointmentService{
			{ServiceID: "svc-1", ServiceName: "Corte", PriceAtBooking: valueobject.NewMoneyFromFloat(50.0), DurationAtBooking: 30},
		}
		appointment, _ := entity.NewAppointment(testTenantUUID, "prof-123", "cust-123", time.Now().Add(24*time.Hour), services)

		mockRepo := &MockAppointmentRepository{
			FindByIDFn: func(ctx context.Context, tenantID, id string) (*entity.Appointment, error) {
				return appointment, nil
			},
		}

		uc := NewUpdateAppointmentStatusUseCase(mockRepo, mockCommandRepo, logger)

		input := UpdateAppointmentStatusInput{
			TenantID:      testTenantUUID.String(),
			AppointmentID: appointment.ID,
			NewStatus:     valueobject.AppointmentStatusConfirmed,
		}

		result, err := uc.Execute(context.Background(), input)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.Status != valueobject.AppointmentStatusConfirmed {
			t.Errorf("expected status CONFIRMED, got %s", result.Status)
		}
	})

	t.Run("should start service from confirmed", func(t *testing.T) {
		services := []entity.AppointmentService{
			{ServiceID: "svc-1", ServiceName: "Corte", PriceAtBooking: valueobject.NewMoneyFromFloat(50.0), DurationAtBooking: 30},
		}
		appointment, _ := entity.NewAppointment(testTenantUUID, "prof-123", "cust-123", time.Now().Add(24*time.Hour), services)
		appointment.Confirm()

		mockRepo := &MockAppointmentRepository{
			FindByIDFn: func(ctx context.Context, tenantID, id string) (*entity.Appointment, error) {
				return appointment, nil
			},
		}

		uc := NewUpdateAppointmentStatusUseCase(mockRepo, mockCommandRepo, logger)

		input := UpdateAppointmentStatusInput{
			TenantID:      testTenantUUID.String(),
			AppointmentID: appointment.ID,
			NewStatus:     valueobject.AppointmentStatusInService,
		}

		result, err := uc.Execute(context.Background(), input)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.Status != valueobject.AppointmentStatusInService {
			t.Errorf("expected status IN_SERVICE, got %s", result.Status)
		}
	})

	t.Run("should finish service to awaiting_payment", func(t *testing.T) {
		services := []entity.AppointmentService{
			{ServiceID: "svc-1", ServiceName: "Corte", PriceAtBooking: valueobject.NewMoneyFromFloat(50.0), DurationAtBooking: 30},
		}
		appointment, _ := entity.NewAppointment(testTenantUUID, "prof-123", "cust-123", time.Now().Add(24*time.Hour), services)
		appointment.Confirm()
		appointment.StartService()

		mockRepo := &MockAppointmentRepository{
			FindByIDFn: func(ctx context.Context, tenantID, id string) (*entity.Appointment, error) {
				return appointment, nil
			},
		}

		uc := NewUpdateAppointmentStatusUseCase(mockRepo, mockCommandRepo, logger)

		input := UpdateAppointmentStatusInput{
			TenantID:      testTenantUUID.String(),
			AppointmentID: appointment.ID,
			NewStatus:     valueobject.AppointmentStatusAwaitingPayment,
		}

		result, err := uc.Execute(context.Background(), input)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.Status != valueobject.AppointmentStatusAwaitingPayment {
			t.Errorf("expected status AWAITING_PAYMENT, got %s", result.Status)
		}
	})

	t.Run("should complete with closed command", func(t *testing.T) {
		commandID := uuid.New()
		services := []entity.AppointmentService{
			{ServiceID: "svc-1", ServiceName: "Corte", PriceAtBooking: valueobject.NewMoneyFromFloat(50.0), DurationAtBooking: 30},
		}
		appointment, _ := entity.NewAppointment(testTenantUUID, "prof-123", "cust-123", time.Now().Add(24*time.Hour), services)
		appointment.Confirm()
		appointment.StartService()
		appointment.FinishService()
		appointment.CommandID = commandID.String()

		mockRepo := &MockAppointmentRepository{
			FindByIDFn: func(ctx context.Context, tenantID, id string) (*entity.Appointment, error) {
				return appointment, nil
			},
		}

		// Mock comanda fechada
		localMockCommandRepo := &MockCommandRepository{
			FindByIDFn: func(ctx context.Context, cmdID, tenantID uuid.UUID) (*entity.Command, error) {
				return &entity.Command{
					ID:       cmdID,
					TenantID: tenantID,
					Status:   entity.CommandStatusClosed,
				}, nil
			},
		}

		uc := NewUpdateAppointmentStatusUseCase(mockRepo, localMockCommandRepo, logger)

		input := UpdateAppointmentStatusInput{
			TenantID:      testTenantUUID.String(),
			AppointmentID: appointment.ID,
			NewStatus:     valueobject.AppointmentStatusDone,
		}

		result, err := uc.Execute(context.Background(), input)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.Status != valueobject.AppointmentStatusDone {
			t.Errorf("expected status DONE, got %s", result.Status)
		}
	})

	t.Run("should fail invalid status transition", func(t *testing.T) {
		services := []entity.AppointmentService{
			{ServiceID: "svc-1", ServiceName: "Corte", PriceAtBooking: valueobject.NewMoneyFromFloat(50.0), DurationAtBooking: 30},
		}
		appointment, _ := entity.NewAppointment(testTenantUUID, "prof-123", "cust-123", time.Now().Add(24*time.Hour), services)
		// Appointment is in CREATED status

		mockRepo := &MockAppointmentRepository{
			FindByIDFn: func(ctx context.Context, tenantID, id string) (*entity.Appointment, error) {
				return appointment, nil
			},
		}

		uc := NewUpdateAppointmentStatusUseCase(mockRepo, mockCommandRepo, logger)

		input := UpdateAppointmentStatusInput{
			TenantID:      testTenantUUID.String(),
			AppointmentID: appointment.ID,
			NewStatus:     valueobject.AppointmentStatusDone, // Can't go directly to DONE from CREATED
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error for invalid status transition, got nil")
		}
	})

	t.Run("should mark no_show", func(t *testing.T) {
		services := []entity.AppointmentService{
			{ServiceID: "svc-1", ServiceName: "Corte", PriceAtBooking: valueobject.NewMoneyFromFloat(50.0), DurationAtBooking: 30},
		}
		appointment, _ := entity.NewAppointment(testTenantUUID, "prof-123", "cust-123", time.Now().Add(24*time.Hour), services)
		appointment.Confirm()

		mockRepo := &MockAppointmentRepository{
			FindByIDFn: func(ctx context.Context, tenantID, id string) (*entity.Appointment, error) {
				return appointment, nil
			},
		}

		uc := NewUpdateAppointmentStatusUseCase(mockRepo, mockCommandRepo, logger)

		input := UpdateAppointmentStatusInput{
			TenantID:      testTenantUUID.String(),
			AppointmentID: appointment.ID,
			NewStatus:     valueobject.AppointmentStatusNoShow,
		}

		result, err := uc.Execute(context.Background(), input)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.Status != valueobject.AppointmentStatusNoShow {
			t.Errorf("expected status NO_SHOW, got %s", result.Status)
		}
	})
}

func TestGetAppointmentUseCase_Execute(t *testing.T) {
	logger := zap.NewNop()

	t.Run("should get appointment successfully", func(t *testing.T) {
		services := []entity.AppointmentService{
			{ServiceID: "svc-1", ServiceName: "Corte", PriceAtBooking: valueobject.NewMoneyFromFloat(50.0), DurationAtBooking: 30},
		}
		appointment, _ := entity.NewAppointment(testTenantUUID, "prof-123", "cust-123", time.Now().Add(24*time.Hour), services)

		mockRepo := &MockAppointmentRepository{
			FindByIDFn: func(ctx context.Context, tenantID, id string) (*entity.Appointment, error) {
				return appointment, nil
			},
		}

		uc := NewGetAppointmentUseCase(mockRepo, logger)

		input := GetAppointmentInput{
			TenantID:      testTenantUUID.String(),
			AppointmentID: appointment.ID,
		}

		result, err := uc.Execute(context.Background(), input)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.ID != appointment.ID {
			t.Errorf("expected appointment ID %s, got %s", appointment.ID, result.ID)
		}
	})

	t.Run("should fail without tenant_id", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{}
		uc := NewGetAppointmentUseCase(mockRepo, logger)

		input := GetAppointmentInput{
			TenantID:      "",
			AppointmentID: "appt-123",
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrTenantIDRequired) {
			t.Errorf("expected ErrTenantIDRequired, got %v", err)
		}
	})

	t.Run("should fail without appointment_id", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{}
		uc := NewGetAppointmentUseCase(mockRepo, logger)

		input := GetAppointmentInput{
			TenantID:      testTenantUUID.String(),
			AppointmentID: "",
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidID) {
			t.Errorf("expected ErrInvalidID, got %v", err)
		}
	})

	t.Run("should fail when appointment not found", func(t *testing.T) {
		mockRepo := &MockAppointmentRepository{
			FindByIDFn: func(ctx context.Context, tenantID, id string) (*entity.Appointment, error) {
				return nil, domain.ErrAppointmentNotFound
			},
		}

		uc := NewGetAppointmentUseCase(mockRepo, logger)

		input := GetAppointmentInput{
			TenantID:      testTenantUUID.String(),
			AppointmentID: "appt-not-exists",
		}

		_, err := uc.Execute(context.Background(), input)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
