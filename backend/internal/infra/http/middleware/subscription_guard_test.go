package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// =============================================================================
// T-ASAAS-003: Testes do Middleware de Verificação de Assinatura
// =============================================================================

// MockSubscriptionChecker é um mock para o SubscriptionChecker
type MockSubscriptionChecker struct {
	mock.Mock
}

func (m *MockSubscriptionChecker) GetTenantSubscriptionStatus(ctx context.Context, tenantID uuid.UUID) (*entity.Subscription, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Subscription), args.Error(1)
}

func (m *MockSubscriptionChecker) HasActiveSubscription(ctx context.Context, tenantID uuid.UUID, gracePeriodDays int) (bool, error) {
	args := m.Called(ctx, tenantID, gracePeriodDays)
	return args.Bool(0), args.Error(1)
}

// setupTest cria um contexto de teste com Echo
func setupTest(tenantID string) (*echo.Echo, *httptest.ResponseRecorder) {
	e := echo.New()
	rec := httptest.NewRecorder()
	return e, rec
}

// TestSubscriptionGuard_AllowsActiveSubscription testa que assinaturas ativas passam
func TestSubscriptionGuard_AllowsActiveSubscription(t *testing.T) {
	e, rec := setupTest("")
	mockChecker := new(MockSubscriptionChecker)
	logger, _ := zap.NewDevelopment()
	tenantID := uuid.New()

	mockChecker.On("HasActiveSubscription", mock.Anything, tenantID, 5).Return(true, nil)

	middleware := RequireActiveSubscription(mockChecker, logger)

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	c := e.NewContext(req, rec)
	c.Set("tenant_id", tenantID.String())

	err := middleware(handler)(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockChecker.AssertExpectations(t)
}

// TestSubscriptionGuard_BlocksInactiveSubscription testa que assinaturas inativas são bloqueadas
func TestSubscriptionGuard_BlocksInactiveSubscription(t *testing.T) {
	e, rec := setupTest("")
	mockChecker := new(MockSubscriptionChecker)
	logger, _ := zap.NewDevelopment()
	tenantID := uuid.New()

	mockChecker.On("HasActiveSubscription", mock.Anything, tenantID, 5).Return(false, nil)

	middleware := RequireActiveSubscription(mockChecker, logger)

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	c := e.NewContext(req, rec)
	c.Set("tenant_id", tenantID.String())

	err := middleware(handler)(c)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusPaymentRequired, httpErr.Code)
	mockChecker.AssertExpectations(t)
}

// TestSubscriptionGuard_AllowsWithoutTenantID testa que requisições sem tenant_id passam
func TestSubscriptionGuard_AllowsWithoutTenantID(t *testing.T) {
	e, rec := setupTest("")
	mockChecker := new(MockSubscriptionChecker)
	logger, _ := zap.NewDevelopment()

	middleware := RequireActiveSubscription(mockChecker, logger)

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	c := e.NewContext(req, rec)
	// Sem tenant_id no contexto

	err := middleware(handler)(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	// Não deve chamar o checker sem tenant_id
	mockChecker.AssertNotCalled(t, "HasActiveSubscription")
}

// =============================================================================
// Testes do DefaultSubscriptionChecker
// =============================================================================

// MockTenantSubscriptionRepo é um mock para o TenantSubscriptionRepo
type MockTenantSubscriptionRepo struct {
	mock.Mock
}

func (m *MockTenantSubscriptionRepo) ListByStatus(ctx context.Context, tenantID uuid.UUID, status entity.SubscriptionStatus) ([]*entity.Subscription, error) {
	args := m.Called(ctx, tenantID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Subscription), args.Error(1)
}

func (m *MockTenantSubscriptionRepo) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.Subscription, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Subscription), args.Error(1)
}

// TestDefaultChecker_ActiveSubscription testa assinatura ativa
func TestDefaultChecker_ActiveSubscription(t *testing.T) {
	mockRepo := new(MockTenantSubscriptionRepo)
	logger, _ := zap.NewDevelopment()
	checker := NewDefaultSubscriptionChecker(mockRepo, logger)

	tenantID := uuid.New()
	activeSub := &entity.Subscription{
		ID:       uuid.New(),
		TenantID: tenantID,
		Status:   entity.StatusAtivo,
	}

	mockRepo.On("ListByStatus", mock.Anything, tenantID, entity.StatusAtivo).Return([]*entity.Subscription{activeSub}, nil)

	hasActive, err := checker.HasActiveSubscription(context.Background(), tenantID, 5)

	assert.NoError(t, err)
	assert.True(t, hasActive)
	mockRepo.AssertExpectations(t)
}

// TestDefaultChecker_ExpiredSubscription testa assinatura vencida
func TestDefaultChecker_ExpiredSubscription(t *testing.T) {
	mockRepo := new(MockTenantSubscriptionRepo)
	logger, _ := zap.NewDevelopment()
	checker := NewDefaultSubscriptionChecker(mockRepo, logger)

	tenantID := uuid.New()
	expiredDate := time.Now().AddDate(0, 0, -10) // Vencida há 10 dias
	expiredSub := &entity.Subscription{
		ID:             uuid.New(),
		TenantID:       tenantID,
		Status:         entity.StatusAtivo,
		DataVencimento: &expiredDate,
	}

	mockRepo.On("ListByStatus", mock.Anything, tenantID, entity.StatusAtivo).Return([]*entity.Subscription{expiredSub}, nil)

	hasActive, err := checker.HasActiveSubscription(context.Background(), tenantID, 5)

	assert.NoError(t, err)
	assert.False(t, hasActive) // Vencida há 10 dias > 5 dias de tolerância
	mockRepo.AssertExpectations(t)
}

// TestDefaultChecker_WithinGracePeriod testa assinatura dentro do período de tolerância
func TestDefaultChecker_WithinGracePeriod(t *testing.T) {
	mockRepo := new(MockTenantSubscriptionRepo)
	logger, _ := zap.NewDevelopment()
	checker := NewDefaultSubscriptionChecker(mockRepo, logger)

	tenantID := uuid.New()
	recentExpiry := time.Now().AddDate(0, 0, -3) // Vencida há 3 dias
	sub := &entity.Subscription{
		ID:             uuid.New(),
		TenantID:       tenantID,
		Status:         entity.StatusAtivo,
		DataVencimento: &recentExpiry,
	}

	mockRepo.On("ListByStatus", mock.Anything, tenantID, entity.StatusAtivo).Return([]*entity.Subscription{sub}, nil)

	hasActive, err := checker.HasActiveSubscription(context.Background(), tenantID, 5)

	assert.NoError(t, err)
	assert.True(t, hasActive) // Vencida há 3 dias < 5 dias de tolerância
	mockRepo.AssertExpectations(t)
}

// TestDefaultChecker_InadimplenteStatus testa status INADIMPLENTE
func TestDefaultChecker_InadimplenteStatus(t *testing.T) {
	mockRepo := new(MockTenantSubscriptionRepo)
	logger, _ := zap.NewDevelopment()
	checker := NewDefaultSubscriptionChecker(mockRepo, logger)

	tenantID := uuid.New()
	inadimplenteSub := &entity.Subscription{
		ID:       uuid.New(),
		TenantID: tenantID,
		Status:   entity.StatusInadimplente,
	}

	// Primeiro busca ativas (nenhuma)
	mockRepo.On("ListByStatus", mock.Anything, tenantID, entity.StatusAtivo).Return([]*entity.Subscription{}, nil)
	// Depois busca todas
	mockRepo.On("ListByTenant", mock.Anything, tenantID).Return([]*entity.Subscription{inadimplenteSub}, nil)

	hasActive, err := checker.HasActiveSubscription(context.Background(), tenantID, 5)

	assert.NoError(t, err)
	assert.False(t, hasActive) // Status INADIMPLENTE = sem acesso
	mockRepo.AssertExpectations(t)
}

// TestDefaultChecker_NoSubscription testa tenant sem assinatura
func TestDefaultChecker_NoSubscription(t *testing.T) {
	mockRepo := new(MockTenantSubscriptionRepo)
	logger, _ := zap.NewDevelopment()
	checker := NewDefaultSubscriptionChecker(mockRepo, logger)

	tenantID := uuid.New()

	mockRepo.On("ListByStatus", mock.Anything, tenantID, entity.StatusAtivo).Return([]*entity.Subscription{}, nil)
	mockRepo.On("ListByTenant", mock.Anything, tenantID).Return([]*entity.Subscription{}, nil)

	hasActive, err := checker.HasActiveSubscription(context.Background(), tenantID, 5)

	assert.NoError(t, err)
	assert.False(t, hasActive) // Sem assinatura = sem acesso
	mockRepo.AssertExpectations(t)
}

// TestDefaultChecker_AguardandoPagamento_WithinCarencia testa novo tenant em carência
func TestDefaultChecker_AguardandoPagamento_WithinCarencia(t *testing.T) {
	mockRepo := new(MockTenantSubscriptionRepo)
	logger, _ := zap.NewDevelopment()
	checker := NewDefaultSubscriptionChecker(mockRepo, logger)

	tenantID := uuid.New()
	recentCreation := time.Now().AddDate(0, 0, -3) // Criada há 3 dias
	newSub := &entity.Subscription{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Status:    entity.StatusAguardandoPagamento,
		CreatedAt: recentCreation,
	}

	mockRepo.On("ListByStatus", mock.Anything, tenantID, entity.StatusAtivo).Return([]*entity.Subscription{}, nil)
	mockRepo.On("ListByTenant", mock.Anything, tenantID).Return([]*entity.Subscription{newSub}, nil)

	hasActive, err := checker.HasActiveSubscription(context.Background(), tenantID, 5)

	assert.NoError(t, err)
	assert.True(t, hasActive) // Aguardando pagamento há 3 dias < 7 dias de carência
	mockRepo.AssertExpectations(t)
}

// TestDefaultChecker_AguardandoPagamento_ExpiredCarencia testa tenant que passou da carência
func TestDefaultChecker_AguardandoPagamento_ExpiredCarencia(t *testing.T) {
	mockRepo := new(MockTenantSubscriptionRepo)
	logger, _ := zap.NewDevelopment()
	checker := NewDefaultSubscriptionChecker(mockRepo, logger)

	tenantID := uuid.New()
	oldCreation := time.Now().AddDate(0, 0, -10) // Criada há 10 dias
	oldSub := &entity.Subscription{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Status:    entity.StatusAguardandoPagamento,
		CreatedAt: oldCreation,
	}

	mockRepo.On("ListByStatus", mock.Anything, tenantID, entity.StatusAtivo).Return([]*entity.Subscription{}, nil)
	mockRepo.On("ListByTenant", mock.Anything, tenantID).Return([]*entity.Subscription{oldSub}, nil)

	hasActive, err := checker.HasActiveSubscription(context.Background(), tenantID, 5)

	assert.NoError(t, err)
	assert.False(t, hasActive) // Aguardando pagamento há 10 dias > 7 dias de carência
	mockRepo.AssertExpectations(t)
}
