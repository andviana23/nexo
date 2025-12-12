package meiopagamento_test

import (
	"context"
	"errors"
	"testing"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/meiopagamento"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// =============================================================================
// MOCK REPOSITORY
// =============================================================================

type MockMeioPagamentoRepository struct {
	mock.Mock
}

func (m *MockMeioPagamentoRepository) Create(ctx context.Context, meio *entity.MeioPagamento) error {
	args := m.Called(ctx, meio)
	return args.Error(0)
}

func (m *MockMeioPagamentoRepository) FindByID(ctx context.Context, tenantID, id string) (*entity.MeioPagamento, error) {
	args := m.Called(ctx, tenantID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MeioPagamento), args.Error(1)
}

func (m *MockMeioPagamentoRepository) List(ctx context.Context, tenantID string) ([]*entity.MeioPagamento, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.MeioPagamento), args.Error(1)
}

func (m *MockMeioPagamentoRepository) ListAtivos(ctx context.Context, tenantID string) ([]*entity.MeioPagamento, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.MeioPagamento), args.Error(1)
}

func (m *MockMeioPagamentoRepository) ListByTipo(ctx context.Context, tenantID string, tipo entity.TipoPagamento) ([]*entity.MeioPagamento, error) {
	args := m.Called(ctx, tenantID, tipo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.MeioPagamento), args.Error(1)
}

func (m *MockMeioPagamentoRepository) Count(ctx context.Context, tenantID string) (int64, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMeioPagamentoRepository) CountAtivos(ctx context.Context, tenantID string) (int64, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMeioPagamentoRepository) Update(ctx context.Context, meio *entity.MeioPagamento) error {
	args := m.Called(ctx, meio)
	return args.Error(0)
}

func (m *MockMeioPagamentoRepository) Toggle(ctx context.Context, tenantID, id string) (*entity.MeioPagamento, error) {
	args := m.Called(ctx, tenantID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MeioPagamento), args.Error(1)
}

func (m *MockMeioPagamentoRepository) Delete(ctx context.Context, tenantID, id string) error {
	args := m.Called(ctx, tenantID, id)
	return args.Error(0)
}

func (m *MockMeioPagamentoRepository) ExistsByNome(ctx context.Context, tenantID, nome string) (bool, error) {
	args := m.Called(ctx, tenantID, nome)
	return args.Bool(0), args.Error(1)
}

// =============================================================================
// TEST HELPERS
// =============================================================================

func newTestMeioPagamento(tenantID uuid.UUID) *entity.MeioPagamento {
	meio, _ := entity.NewMeioPagamento(tenantID, "PIX Teste", entity.TipoPagamentoPIX)
	meio.Taxa = decimal.NewFromFloat(0)
	meio.TaxaFixa = decimal.NewFromFloat(0)
	meio.DMais = 0
	meio.Ativo = true
	return meio
}

// =============================================================================
// CREATE USE CASE TESTS
// =============================================================================

func TestCreateMeioPagamentoUseCase_Execute(t *testing.T) {
	t.Run("should create meio pagamento successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockMeioPagamentoRepository)
		useCase := meiopagamento.NewCreateMeioPagamentoUseCase(mockRepo)

		ctx := context.Background()
		tenantID := uuid.New().String()

		req := dto.CreateMeioPagamentoRequest{
			Nome:  "PIX Banco X",
			Tipo:  "PIX",
			Taxa:  "0",
			DMais: 0,
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.MeioPagamento")).Return(nil)

		// Act
		result, err := useCase.Execute(ctx, tenantID, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "PIX Banco X", result.Nome)
		assert.Equal(t, "PIX", result.Tipo)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail with invalid tenant_id", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockMeioPagamentoRepository)
		useCase := meiopagamento.NewCreateMeioPagamentoUseCase(mockRepo)

		ctx := context.Background()
		req := dto.CreateMeioPagamentoRequest{
			Nome: "PIX Teste",
			Tipo: "PIX",
		}

		// Act
		result, err := useCase.Execute(ctx, "invalid-uuid", req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "tenant_id inválido")
	})

	t.Run("should fail with invalid tipo", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockMeioPagamentoRepository)
		useCase := meiopagamento.NewCreateMeioPagamentoUseCase(mockRepo)

		ctx := context.Background()
		tenantID := uuid.New().String()

		req := dto.CreateMeioPagamentoRequest{
			Nome: "Teste",
			Tipo: "TIPO_INVALIDO",
		}

		// Act
		result, err := useCase.Execute(ctx, tenantID, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, entity.ErrMeioPagamentoTipoInvalido, err)
	})

	t.Run("should auto-generate nome when empty", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockMeioPagamentoRepository)
		useCase := meiopagamento.NewCreateMeioPagamentoUseCase(mockRepo)

		ctx := context.Background()
		tenantID := uuid.New().String()

		req := dto.CreateMeioPagamentoRequest{
			Nome: "",
			Tipo: "PIX",
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.MeioPagamento")).Return(nil)

		// Act
		result, err := useCase.Execute(ctx, tenantID, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "PIX", result.Nome)
		assert.Equal(t, "PIX", result.Tipo)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail with invalid taxa over 100", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockMeioPagamentoRepository)
		useCase := meiopagamento.NewCreateMeioPagamentoUseCase(mockRepo)

		ctx := context.Background()
		tenantID := uuid.New().String()

		req := dto.CreateMeioPagamentoRequest{
			Nome: "Crédito",
			Tipo: "CREDITO",
			Taxa: "150",
		}

		// Act
		result, err := useCase.Execute(ctx, tenantID, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, entity.ErrMeioPagamentoTaxaInvalida, err)
	})

	t.Run("should fail when repository returns error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockMeioPagamentoRepository)
		useCase := meiopagamento.NewCreateMeioPagamentoUseCase(mockRepo)

		ctx := context.Background()
		tenantID := uuid.New().String()

		req := dto.CreateMeioPagamentoRequest{
			Nome: "PIX",
			Tipo: "PIX",
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.MeioPagamento")).Return(errors.New("db error"))

		// Act
		result, err := useCase.Execute(ctx, tenantID, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "erro ao criar meio de pagamento")
	})

	t.Run("should create with all optional fields", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockMeioPagamentoRepository)
		useCase := meiopagamento.NewCreateMeioPagamentoUseCase(mockRepo)

		ctx := context.Background()
		tenantID := uuid.New().String()
		ativo := true

		req := dto.CreateMeioPagamentoRequest{
			Nome:          "Crédito Visa",
			Tipo:          "CREDITO",
			Bandeira:      "Visa",
			Taxa:          "2.99",
			TaxaFixa:      "0.50",
			DMais:         30,
			Icone:         "credit-card",
			Cor:           "#3b82f6",
			OrdemExibicao: 1,
			Observacoes:   "Cartão de crédito Visa",
			Ativo:         &ativo,
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.MeioPagamento")).Return(nil)

		// Act
		result, err := useCase.Execute(ctx, tenantID, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Crédito Visa", result.Nome)
		assert.Equal(t, "CREDITO", result.Tipo)
		assert.Equal(t, "Visa", result.Bandeira)
		assert.Equal(t, 30, result.DMais)
		mockRepo.AssertExpectations(t)
	})
}

// =============================================================================
// LIST USE CASE TESTS
// =============================================================================

func TestListMeiosPagamentoUseCase_Execute(t *testing.T) {
	t.Run("should list all meios pagamento", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockMeioPagamentoRepository)
		useCase := meiopagamento.NewListMeiosPagamentoUseCase(mockRepo)

		ctx := context.Background()
		tenantID := uuid.New().String()

		meios := []*entity.MeioPagamento{
			newTestMeioPagamento(uuid.MustParse(tenantID)),
			newTestMeioPagamento(uuid.MustParse(tenantID)),
		}

		mockRepo.On("List", ctx, tenantID).Return(meios, nil)
		mockRepo.On("Count", ctx, tenantID).Return(int64(2), nil)
		mockRepo.On("CountAtivos", ctx, tenantID).Return(int64(2), nil)

		filter := dto.MeioPagamentoFilter{}

		// Act
		result, err := useCase.Execute(ctx, tenantID, filter)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Data, 2)
		assert.Equal(t, int64(2), result.Total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should list only active meios", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockMeioPagamentoRepository)
		useCase := meiopagamento.NewListMeiosPagamentoUseCase(mockRepo)

		ctx := context.Background()
		tenantID := uuid.New().String()

		meios := []*entity.MeioPagamento{
			newTestMeioPagamento(uuid.MustParse(tenantID)),
		}

		mockRepo.On("ListAtivos", ctx, tenantID).Return(meios, nil)
		mockRepo.On("Count", ctx, tenantID).Return(int64(2), nil)
		mockRepo.On("CountAtivos", ctx, tenantID).Return(int64(1), nil)

		filter := dto.MeioPagamentoFilter{ApenasAtivos: true}

		// Act
		result, err := useCase.Execute(ctx, tenantID, filter)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Data, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should list by tipo", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockMeioPagamentoRepository)
		useCase := meiopagamento.NewListMeiosPagamentoUseCase(mockRepo)

		ctx := context.Background()
		tenantID := uuid.New().String()

		meios := []*entity.MeioPagamento{
			newTestMeioPagamento(uuid.MustParse(tenantID)),
		}

		mockRepo.On("ListByTipo", ctx, tenantID, entity.TipoPagamentoPIX).Return(meios, nil)
		mockRepo.On("Count", ctx, tenantID).Return(int64(5), nil)
		mockRepo.On("CountAtivos", ctx, tenantID).Return(int64(4), nil)

		filter := dto.MeioPagamentoFilter{Tipo: "PIX"}

		// Act
		result, err := useCase.Execute(ctx, tenantID, filter)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Data, 1)
		mockRepo.AssertExpectations(t)
	})
}

// =============================================================================
// GET USE CASE TESTS
// =============================================================================

func TestGetMeioPagamentoUseCase_Execute(t *testing.T) {
	t.Run("should get meio pagamento by id", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockMeioPagamentoRepository)
		useCase := meiopagamento.NewGetMeioPagamentoUseCase(mockRepo)

		ctx := context.Background()
		tenantID := uuid.New().String()
		meioPagamentoID := uuid.New().String()

		meio := newTestMeioPagamento(uuid.MustParse(tenantID))

		mockRepo.On("FindByID", ctx, tenantID, meioPagamentoID).Return(meio, nil)

		// Act
		result, err := useCase.Execute(ctx, tenantID, meioPagamentoID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when meio not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockMeioPagamentoRepository)
		useCase := meiopagamento.NewGetMeioPagamentoUseCase(mockRepo)

		ctx := context.Background()
		tenantID := uuid.New().String()
		meioPagamentoID := uuid.New().String()

		mockRepo.On("FindByID", ctx, tenantID, meioPagamentoID).Return(nil, errors.New("not found"))

		// Act
		result, err := useCase.Execute(ctx, tenantID, meioPagamentoID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

// =============================================================================
// UPDATE USE CASE TESTS
// =============================================================================

func TestUpdateMeioPagamentoUseCase_Execute(t *testing.T) {
	t.Run("should update meio pagamento successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockMeioPagamentoRepository)
		useCase := meiopagamento.NewUpdateMeioPagamentoUseCase(mockRepo)

		ctx := context.Background()
		tenantID := uuid.New().String()
		meioPagamentoID := uuid.New().String()

		meio := newTestMeioPagamento(uuid.MustParse(tenantID))

		mockRepo.On("FindByID", ctx, tenantID, meioPagamentoID).Return(meio, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.MeioPagamento")).Return(nil)

		req := dto.UpdateMeioPagamentoRequest{
			Nome: "PIX Atualizado",
			Taxa: "1.5",
		}

		// Act
		result, err := useCase.Execute(ctx, tenantID, meioPagamentoID, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "PIX Atualizado", result.Nome)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when meio not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockMeioPagamentoRepository)
		useCase := meiopagamento.NewUpdateMeioPagamentoUseCase(mockRepo)

		ctx := context.Background()
		tenantID := uuid.New().String()
		meioPagamentoID := uuid.New().String()

		mockRepo.On("FindByID", ctx, tenantID, meioPagamentoID).Return(nil, errors.New("not found"))

		req := dto.UpdateMeioPagamentoRequest{Nome: "Teste"}

		// Act
		result, err := useCase.Execute(ctx, tenantID, meioPagamentoID, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail with invalid tipo", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockMeioPagamentoRepository)
		useCase := meiopagamento.NewUpdateMeioPagamentoUseCase(mockRepo)

		ctx := context.Background()
		tenantID := uuid.New().String()
		meioPagamentoID := uuid.New().String()

		meio := newTestMeioPagamento(uuid.MustParse(tenantID))

		mockRepo.On("FindByID", ctx, tenantID, meioPagamentoID).Return(meio, nil)

		req := dto.UpdateMeioPagamentoRequest{
			Tipo: "TIPO_INVALIDO",
		}

		// Act
		result, err := useCase.Execute(ctx, tenantID, meioPagamentoID, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, entity.ErrMeioPagamentoTipoInvalido, err)
	})
}

// =============================================================================
// TOGGLE USE CASE TESTS
// =============================================================================

func TestToggleMeioPagamentoUseCase_Execute(t *testing.T) {
	t.Run("should toggle meio pagamento status", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockMeioPagamentoRepository)
		useCase := meiopagamento.NewToggleMeioPagamentoUseCase(mockRepo)

		ctx := context.Background()
		tenantID := uuid.New().String()
		meioPagamentoID := uuid.New().String()

		meio := newTestMeioPagamento(uuid.MustParse(tenantID))
		meio.Ativo = false // toggled to inactive

		mockRepo.On("Toggle", ctx, tenantID, meioPagamentoID).Return(meio, nil)

		// Act
		result, err := useCase.Execute(ctx, tenantID, meioPagamentoID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.Ativo)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when toggle fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockMeioPagamentoRepository)
		useCase := meiopagamento.NewToggleMeioPagamentoUseCase(mockRepo)

		ctx := context.Background()
		tenantID := uuid.New().String()
		meioPagamentoID := uuid.New().String()

		mockRepo.On("Toggle", ctx, tenantID, meioPagamentoID).Return(nil, errors.New("toggle failed"))

		// Act
		result, err := useCase.Execute(ctx, tenantID, meioPagamentoID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

// =============================================================================
// DELETE USE CASE TESTS
// =============================================================================

func TestDeleteMeioPagamentoUseCase_Execute(t *testing.T) {
	t.Run("should delete meio pagamento successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockMeioPagamentoRepository)
		useCase := meiopagamento.NewDeleteMeioPagamentoUseCase(mockRepo)

		ctx := context.Background()
		tenantID := uuid.New().String()
		meioPagamentoID := uuid.New().String()

		mockRepo.On("Delete", ctx, tenantID, meioPagamentoID).Return(nil)

		// Act
		err := useCase.Execute(ctx, tenantID, meioPagamentoID)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when delete fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockMeioPagamentoRepository)
		useCase := meiopagamento.NewDeleteMeioPagamentoUseCase(mockRepo)

		ctx := context.Background()
		tenantID := uuid.New().String()
		meioPagamentoID := uuid.New().String()

		mockRepo.On("Delete", ctx, tenantID, meioPagamentoID).Return(errors.New("delete failed"))

		// Act
		err := useCase.Execute(ctx, tenantID, meioPagamentoID)

		// Assert
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
