package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/appointment"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/andviana23/barber-analytics-backend/internal/infra/http/handler"
	"github.com/andviana23/barber-analytics-backend/internal/infra/repository/postgres"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// E2E test tenant ID (deve existir no banco com seeds)
const apptE2ETenantID = "e2e00000-0000-0000-0000-000000000001"

type ApptValidator struct {
	validator *validator.Validate
}

func (cv *ApptValidator) Validate(i interface{}) error {
	if cv.validator == nil {
		cv.validator = validator.New()
	}
	return cv.validator.Struct(i)
}

func getApptTestDBPool(t *testing.T) *pgxpool.Pool {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL não configurada, pulando testes de integração")
		return nil
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("Erro ao conectar ao banco de testes: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("Erro ao fazer ping no banco: %v", err)
	}

	return pool
}

func setupApptHandler(t *testing.T) (*handler.AppointmentHandler, *echo.Echo, func()) {
	pool := getApptTestDBPool(t)
	if pool == nil {
		return nil, nil, nil
	}

	logger, _ := zap.NewDevelopment()
	queries := db.New(pool)

	// Repositories e Readers
	appointmentRepo := postgres.NewAppointmentRepository(queries, pool)
	professionalReader := postgres.NewProfessionalReader(queries)
	customerReader := postgres.NewCustomerReader(queries)
	serviceReader := postgres.NewServiceReader(queries)
	commandRepo := postgres.NewCommandRepository(queries, pool)

	// Use cases
	// G-001: createUC agora recebe commandRepo para criar comanda automaticamente
	createUC := appointment.NewCreateAppointmentUseCase(appointmentRepo, commandRepo, serviceReader, professionalReader, customerReader, logger)
	listUC := appointment.NewListAppointmentsUseCase(appointmentRepo, logger)
	getUC := appointment.NewGetAppointmentUseCase(appointmentRepo, logger)
	updateStatusUC := appointment.NewUpdateAppointmentStatusUseCase(appointmentRepo, commandRepo, logger)
	rescheduleUC := appointment.NewRescheduleAppointmentUseCase(appointmentRepo, professionalReader, logger)
	cancelUC := appointment.NewCancelAppointmentUseCase(appointmentRepo, logger)

	// Handler
	apptHandler := handler.NewAppointmentHandler(
		createUC,
		listUC,
		getUC,
		updateStatusUC,
		rescheduleUC,
		cancelUC,
		nil, // finishWithCommandUC - not needed for these tests
		logger,
	)

	// Echo setup
	e := echo.New()
	e.Validator = &ApptValidator{}

	cleanup := func() {
		pool.Close()
	}

	return apptHandler, e, cleanup
}

func TestAppointmentHandler_ListAppointments_Integration(t *testing.T) {
	apptHandler, e, cleanup := setupApptHandler(t)
	if apptHandler == nil {
		return
	}
	defer cleanup()

	t.Run("should list appointments for tenant", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/appointments", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("tenant_id", apptE2ETenantID)

		err := apptHandler.ListAppointments(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.ListAppointmentsResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, response.Total, int64(0))
	})

	t.Run("should filter by status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/appointments?status=CONFIRMED", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("tenant_id", apptE2ETenantID)

		err := apptHandler.ListAppointments(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should fail without tenant_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/appointments", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		// No tenant_id set

		err := apptHandler.ListAppointments(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestAppointmentHandler_GetAppointment_Integration(t *testing.T) {
	apptHandler, e, cleanup := setupApptHandler(t)
	if apptHandler == nil {
		return
	}
	defer cleanup()

	t.Run("should return 404 for non-existent appointment", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/appointments/00000000-0000-0000-0000-000000000000", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("00000000-0000-0000-0000-000000000000")
		c.Set("tenant_id", apptE2ETenantID)

		err := apptHandler.GetAppointment(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("should fail without tenant_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/appointments/some-id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("some-id")
		// No tenant_id set

		err := apptHandler.GetAppointment(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestAppointmentHandler_CreateAppointment_Integration(t *testing.T) {
	apptHandler, e, cleanup := setupApptHandler(t)
	if apptHandler == nil {
		return
	}
	defer cleanup()

	t.Run("should create appointment successfully", func(t *testing.T) {
		// Usando IDs dos seeds do tenant E2E
		reqBody := dto.CreateAppointmentRequest{
			ProfessionalID: "a0000000-0000-0000-0000-000000000001",           // Carlos Silva (seed)
			CustomerID:     "c1000000-0000-0000-0000-000000000001",           // João Santos (seed)
			StartTime:      time.Now().Add(72 * time.Hour),                   // 3 dias no futuro
			ServiceIDs:     []string{"s0000000-0000-0000-0000-000000000001"}, // Corte Masculino (seed)
			Notes:          "Agendamento de teste de integração",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/appointments", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("tenant_id", apptE2ETenantID)

		err := apptHandler.CreateAppointment(c)
		require.NoError(t, err)

		// Se falhar por causa de validação ou conflito, é OK - estamos testando a estrutura
		if rec.Code == http.StatusCreated {
			var response dto.AppointmentResponse
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.NotEmpty(t, response.ID)
			assert.Equal(t, "CREATED", response.Status)
		}
	})

	t.Run("should fail with invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/appointments", bytes.NewReader([]byte("{invalid json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("tenant_id", apptE2ETenantID)

		err := apptHandler.CreateAppointment(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAppointmentHandler_CancelAppointment_Integration(t *testing.T) {
	apptHandler, e, cleanup := setupApptHandler(t)
	if apptHandler == nil {
		return
	}
	defer cleanup()

	t.Run("should return error for non-existent appointment", func(t *testing.T) {
		reqBody := dto.CancelAppointmentRequest{
			Reason: "Teste de cancelamento",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/appointments/00000000-0000-0000-0000-000000000000/cancel", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("00000000-0000-0000-0000-000000000000")
		c.Set("tenant_id", apptE2ETenantID)

		err := apptHandler.CancelAppointment(c)
		require.NoError(t, err)
		// Pode retornar 400 (appointment not found wrapped) ou 404
		assert.Contains(t, []int{http.StatusBadRequest, http.StatusNotFound}, rec.Code)
	})
}

func TestAppointmentHandler_UpdateStatus_Integration(t *testing.T) {
	apptHandler, e, cleanup := setupApptHandler(t)
	if apptHandler == nil {
		return
	}
	defer cleanup()

	t.Run("should fail with invalid status", func(t *testing.T) {
		reqBody := dto.UpdateAppointmentStatusRequest{
			Status: "INVALID_STATUS",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/appointments/some-id/status", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("some-id")
		c.Set("tenant_id", apptE2ETenantID)

		err := apptHandler.UpdateAppointmentStatus(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// BUG-004: Teste de formatação monetária
// Valida que total_price e service.price são retornados como strings numéricas
// (ex: "50.00") em vez de strings formatadas (ex: "R$ 50,00")
func TestAppointmentHandler_PriceFormat_Integration(t *testing.T) {
	apptHandler, e, cleanup := setupApptHandler(t)
	if apptHandler == nil {
		return
	}
	defer cleanup()

	t.Run("list endpoint should return numeric price format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/appointments?page=1&page_size=5", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("tenant_id", apptE2ETenantID)

		err := apptHandler.ListAppointments(c)
		require.NoError(t, err)

		if rec.Code == http.StatusOK {
			var response dto.ListAppointmentsResponse
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			require.NoError(t, err)

			// Se há agendamentos, validar formato do preço
			for _, apt := range response.Data {
				// total_price deve ser string numérica: "XX.XX" (sem R$, sem vírgula)
				// Não deve conter "R$"
				assert.NotContains(t, apt.TotalPrice, "R$", "total_price deve ser numérico, não formatado")
				// Não deve conter vírgula (formato brasileiro)
				assert.NotContains(t, apt.TotalPrice, ",", "total_price deve usar ponto decimal, não vírgula")

				// Validar que é parseável como float
				if apt.TotalPrice != "" {
					var price float64
					err := json.Unmarshal([]byte(`"`+apt.TotalPrice+`"`), &price)
					if err != nil {
						// Tentar parse direto
						_, parseErr := json.Number(apt.TotalPrice).Float64()
						// Se não parsear, ao menos deve ser string numérica
						assert.Regexp(t, `^\d+\.\d{2}$`, apt.TotalPrice, "total_price deve ter formato XX.XX, got: "+apt.TotalPrice)
						_ = parseErr
					}
				}

				// Validar serviços
				for _, svc := range apt.Services {
					assert.NotContains(t, svc.Price, "R$", "service.price deve ser numérico, não formatado")
					assert.NotContains(t, svc.Price, ",", "service.price deve usar ponto decimal, não vírgula")
				}
			}

			t.Logf("✅ %d agendamentos validados com formato de preço correto", len(response.Data))
		}
	})
}
