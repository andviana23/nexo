package blockedtime_test

import (
	"context"
	"testing"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/blockedtime"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBlockedTimeRepository é um mock do repositório
type MockBlockedTimeRepository struct {
	mock.Mock
}

func (m *MockBlockedTimeRepository) Create(ctx context.Context, bt *entity.BlockedTime) (*entity.BlockedTime, error) {
	args := m.Called(ctx, bt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BlockedTime), args.Error(1)
}

func (m *MockBlockedTimeRepository) GetByID(ctx context.Context, tenantID, id string) (*entity.BlockedTime, error) {
	args := m.Called(ctx, tenantID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BlockedTime), args.Error(1)
}

func (m *MockBlockedTimeRepository) List(ctx context.Context, tenantID string, professionalID *string, startDate, endDate *time.Time) ([]*entity.BlockedTime, error) {
	args := m.Called(ctx, tenantID, professionalID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.BlockedTime), args.Error(1)
}

func (m *MockBlockedTimeRepository) GetInRange(ctx context.Context, tenantID, professionalID string, startTime, endTime time.Time) ([]*entity.BlockedTime, error) {
	args := m.Called(ctx, tenantID, professionalID, startTime, endTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.BlockedTime), args.Error(1)
}

func (m *MockBlockedTimeRepository) CheckConflict(ctx context.Context, tenantID, professionalID string, startTime, endTime time.Time, excludeID *string) (bool, error) {
	args := m.Called(ctx, tenantID, professionalID, startTime, endTime, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockBlockedTimeRepository) Update(ctx context.Context, bt *entity.BlockedTime) (*entity.BlockedTime, error) {
	args := m.Called(ctx, bt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BlockedTime), args.Error(1)
}

func (m *MockBlockedTimeRepository) Delete(ctx context.Context, tenantID, id string) error {
	args := m.Called(ctx, tenantID, id)
	return args.Error(0)
}

// TestCreateBlockedTime testa a criação de bloqueio sem conflito
func TestCreateBlockedTime_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockBlockedTimeRepository)
	useCase := blockedtime.NewCreateBlockedTimeUseCase(mockRepo)

	ctx := context.Background()
	tenantID := uuid.New().String()
	tenantUUID := uuid.MustParse(tenantID)
	professionalID := uuid.New().String()
	startTime := time.Now().Add(1 * time.Hour)
	endTime := startTime.Add(1 * time.Hour)
	reason := "Almoço"

	// Mock: Sem conflito
	mockRepo.On("CheckConflict", ctx, tenantID, professionalID, startTime, endTime, (*string)(nil)).Return(false, nil)

	// Mock: Criação bem-sucedida
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.BlockedTime")).Return(&entity.BlockedTime{
		ID:             "blocked-789",
		TenantID:       tenantUUID,
		ProfessionalID: professionalID,
		StartTime:      startTime,
		EndTime:        endTime,
		Reason:         reason,
	}, nil)

	// Act
	input := blockedtime.CreateBlockedTimeInput{
		TenantID:       tenantID,
		ProfessionalID: professionalID,
		StartTime:      startTime,
		EndTime:        endTime,
		Reason:         reason,
	}

	output, err := useCase.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotNil(t, output.BlockedTime)
	assert.Equal(t, tenantUUID, output.BlockedTime.TenantID)
	assert.Equal(t, professionalID, output.BlockedTime.ProfessionalID)
	assert.Equal(t, reason, output.BlockedTime.Reason)

	mockRepo.AssertExpectations(t)
}

// TestCreateBlockedTime_Conflict testa criação de bloqueio com conflito
func TestCreateBlockedTime_Conflict(t *testing.T) {
	// Arrange
	mockRepo := new(MockBlockedTimeRepository)
	useCase := blockedtime.NewCreateBlockedTimeUseCase(mockRepo)

	ctx := context.Background()
	tenantID := uuid.New().String()
	professionalID := uuid.New().String()
	startTime := time.Now().Add(1 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)

	// Mock: Há conflito
	mockRepo.On("CheckConflict", ctx, tenantID, professionalID, startTime, endTime, (*string)(nil)).Return(true, nil)

	// Act
	input := blockedtime.CreateBlockedTimeInput{
		TenantID:       tenantID,
		ProfessionalID: professionalID,
		StartTime:      startTime,
		EndTime:        endTime,
		Reason:         "Reunião",
	}

	output, err := useCase.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, entity.ErrTimeRangeOverlap, err)

	mockRepo.AssertExpectations(t)
}

// TestListBlockedTimes testa a listagem de bloqueios
func TestListBlockedTimes_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockBlockedTimeRepository)
	useCase := blockedtime.NewListBlockedTimesUseCase(mockRepo)

	ctx := context.Background()
	tenantID := uuid.New().String()
	tenantUUID := uuid.MustParse(tenantID)
	professionalID := uuid.New().String()

	blockedID1 := uuid.New().String()
	blockedID2 := uuid.New().String()

	expectedList := []*entity.BlockedTime{
		{
			ID:             blockedID1,
			TenantID:       tenantUUID,
			ProfessionalID: professionalID,
			Reason:         "Almoço",
		},
		{
			ID:             blockedID2,
			TenantID:       tenantUUID,
			ProfessionalID: professionalID,
			Reason:         "Reunião",
		},
	}

	mockRepo.On("List", ctx, tenantID, &professionalID, (*time.Time)(nil), (*time.Time)(nil)).Return(expectedList, nil)

	// Act
	input := blockedtime.ListBlockedTimesInput{
		TenantID:       tenantID,
		ProfessionalID: &professionalID,
	}

	output, err := useCase.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Len(t, output.BlockedTimes, 2)
	assert.Equal(t, blockedID1, output.BlockedTimes[0].ID)
	assert.Equal(t, blockedID2, output.BlockedTimes[1].ID)

	mockRepo.AssertExpectations(t)
}

// TestDeleteBlockedTime testa a exclusão
func TestDeleteBlockedTime_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockBlockedTimeRepository)
	useCase := blockedtime.NewDeleteBlockedTimeUseCase(mockRepo)

	ctx := context.Background()
	tenantID := uuid.New().String()
	tenantUUID := uuid.MustParse(tenantID)
	blockedID := uuid.New().String()

	existingBlocked := &entity.BlockedTime{
		ID:       blockedID,
		TenantID: tenantUUID,
		Reason:   "Almoço",
	}

	mockRepo.On("GetByID", ctx, tenantID, blockedID).Return(existingBlocked, nil)
	mockRepo.On("Delete", ctx, tenantID, blockedID).Return(nil)

	// Act
	input := blockedtime.DeleteBlockedTimeInput{
		TenantID: tenantID,
		ID:       blockedID,
	}

	err := useCase.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
