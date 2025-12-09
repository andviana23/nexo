package metas_test

import (
	"context"
	"errors"
	"testing"

	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/metas"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// TenantID de teste fixo
var testTenantUUID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
var testTenantStr = testTenantUUID.String()

// MockMetaMensalRepository é um mock do repositório
type MockMetaMensalRepository struct {
	mock.Mock
}

func (m *MockMetaMensalRepository) Create(ctx context.Context, meta *entity.MetaMensal) error {
	args := m.Called(ctx, meta)
	return args.Error(0)
}

func (m *MockMetaMensalRepository) FindByID(ctx context.Context, tenantID, id string) (*entity.MetaMensal, error) {
	args := m.Called(ctx, tenantID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MetaMensal), args.Error(1)
}

func (m *MockMetaMensalRepository) FindByMesAno(ctx context.Context, tenantID string, mesAno valueobject.MesAno) (*entity.MetaMensal, error) {
	args := m.Called(ctx, tenantID, mesAno)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MetaMensal), args.Error(1)
}

func (m *MockMetaMensalRepository) ListAtivas(ctx context.Context, tenantID string) ([]*entity.MetaMensal, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.MetaMensal), args.Error(1)
}

func (m *MockMetaMensalRepository) ListByPeriod(ctx context.Context, tenantID string, inicio, fim valueobject.MesAno) ([]*entity.MetaMensal, error) {
	args := m.Called(ctx, tenantID, inicio, fim)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.MetaMensal), args.Error(1)
}

func (m *MockMetaMensalRepository) Update(ctx context.Context, meta *entity.MetaMensal) error {
	args := m.Called(ctx, meta)
	return args.Error(0)
}

func (m *MockMetaMensalRepository) Delete(ctx context.Context, tenantID, id string) error {
	args := m.Called(ctx, tenantID, id)
	return args.Error(0)
}

func TestSetMetaMensalUseCase_Execute_Success_CreateNew(t *testing.T) {
	// Arrange
	mockRepo := new(MockMetaMensalRepository)
	logger := zap.NewNop()
	useCase := metas.NewSetMetaMensalUseCase(mockRepo, logger)

	ctx := context.Background()
	tenantID := testTenantStr
	mesAno, _ := valueobject.NewMesAno("2025-01")
	metaFaturamento := valueobject.NewMoney(100000) // R$ 1.000,00

	input := metas.SetMetaMensalInput{
		TenantID:        testTenantStr,
		MesAno:          mesAno,
		MetaFaturamento: metaFaturamento,
		Origem:          valueobject.OrigemMetaManual,
	}

	// Mock: meta não existe ainda
	mockRepo.On("FindByMesAno", ctx, tenantID, mesAno).Return(nil, domain.ErrNotFound)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.MetaMensal")).Return(nil)

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testTenantUUID, result.TenantID)
	assert.Equal(t, mesAno, result.MesAno)
	assert.Equal(t, metaFaturamento, result.MetaFaturamento)
	mockRepo.AssertExpectations(t)
}

func TestSetMetaMensalUseCase_Execute_Success_UpdateExisting(t *testing.T) {
	// Arrange
	mockRepo := new(MockMetaMensalRepository)
	logger := zap.NewNop()
	useCase := metas.NewSetMetaMensalUseCase(mockRepo, logger)

	ctx := context.Background()
	tenantID := testTenantStr
	mesAno, _ := valueobject.NewMesAno("2025-01")
	oldMeta := valueobject.NewMoney(100000) // R$ 1.000,00
	newMeta := valueobject.NewMoney(150000) // R$ 1.500,00

	// Meta existente
	existingMeta := &entity.MetaMensal{
		ID:              "meta-123",
		TenantID:        testTenantUUID,
		MesAno:          mesAno,
		MetaFaturamento: oldMeta,
		Origem:          valueobject.OrigemMetaManual,
	}

	input := metas.SetMetaMensalInput{
		TenantID:        testTenantStr,
		MesAno:          mesAno,
		MetaFaturamento: newMeta,
		Origem:          valueobject.OrigemMetaManual,
	}

	// Mock: meta existe
	mockRepo.On("FindByMesAno", ctx, tenantID, mesAno).Return(existingMeta, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.MetaMensal")).Return(nil)

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testTenantUUID, result.TenantID)
	assert.Equal(t, mesAno, result.MesAno)
	assert.Equal(t, newMeta, result.MetaFaturamento)
	mockRepo.AssertExpectations(t)
}

func TestSetMetaMensalUseCase_Execute_Error_TenantIDRequired(t *testing.T) {
	// Arrange
	mockRepo := new(MockMetaMensalRepository)
	logger := zap.NewNop()
	useCase := metas.NewSetMetaMensalUseCase(mockRepo, logger)

	ctx := context.Background()
	mesAno, _ := valueobject.NewMesAno("2025-01")
	input := metas.SetMetaMensalInput{
		TenantID:        "", // vazio
		MesAno:          mesAno,
		MetaFaturamento: valueobject.NewMoney(100000),
		Origem:          valueobject.OrigemMetaManual,
	}

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrTenantIDRequired, err)
	assert.Nil(t, result)
	mockRepo.AssertNotCalled(t, "FindByMesAno")
	mockRepo.AssertNotCalled(t, "Create")
	mockRepo.AssertNotCalled(t, "Update")
}

func TestSetMetaMensalUseCase_Execute_Error_RepositoryCreate(t *testing.T) {
	// Arrange
	mockRepo := new(MockMetaMensalRepository)
	logger := zap.NewNop()
	useCase := metas.NewSetMetaMensalUseCase(mockRepo, logger)

	ctx := context.Background()
	tenantID := testTenantStr
	mesAno, _ := valueobject.NewMesAno("2025-01")
	metaFaturamento := valueobject.NewMoney(100000)

	input := metas.SetMetaMensalInput{
		TenantID:        testTenantStr,
		MesAno:          mesAno,
		MetaFaturamento: metaFaturamento,
		Origem:          valueobject.OrigemMetaManual,
	}

	// Mock: meta não existe, mas erro ao criar
	mockRepo.On("FindByMesAno", ctx, tenantID, mesAno).Return(nil, domain.ErrNotFound)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.MetaMensal")).Return(errors.New("database error"))

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
	mockRepo.AssertExpectations(t)
}

func TestSetMetaMensalUseCase_Execute_Error_RepositoryUpdate(t *testing.T) {
	// Arrange
	mockRepo := new(MockMetaMensalRepository)
	logger := zap.NewNop()
	useCase := metas.NewSetMetaMensalUseCase(mockRepo, logger)

	ctx := context.Background()
	tenantID := testTenantStr
	mesAno, _ := valueobject.NewMesAno("2025-01")
	oldMeta := valueobject.NewMoney(100000)
	newMeta := valueobject.NewMoney(150000)

	existingMeta := &entity.MetaMensal{
		ID:              "meta-123",
		TenantID:        testTenantUUID,
		MesAno:          mesAno,
		MetaFaturamento: oldMeta,
		Origem:          valueobject.OrigemMetaManual,
	}

	input := metas.SetMetaMensalInput{
		TenantID:        testTenantStr,
		MesAno:          mesAno,
		MetaFaturamento: newMeta,
		Origem:          valueobject.OrigemMetaManual,
	}

	// Mock: meta existe, mas erro ao atualizar
	mockRepo.On("FindByMesAno", ctx, tenantID, mesAno).Return(existingMeta, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.MetaMensal")).Return(errors.New("update error"))

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "update error")
	mockRepo.AssertExpectations(t)
}
