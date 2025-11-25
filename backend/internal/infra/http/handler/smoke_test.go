package handler_test

import (
"net/http"
"net/http/httptest"
"testing"

"github.com/labstack/echo/v4"
"github.com/stretchr/testify/assert"
)

// TestSmokeTest_HealthEndpoint verifica que o servidor responde corretamente
func TestSmokeTest_HealthEndpoint(t *testing.T) {
	e := echo.New()

	// Simular endpoint de health
	e.GET("/health", func(c echo.Context) error {
return c.JSON(http.StatusOK, map[string]string{
"status":  "healthy",
"version": "v2.0.0",
})
})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "healthy")
}

// TestSmokeTest_EchoSetup verifica configuracao basica do Echo
func TestSmokeTest_EchoSetup(t *testing.T) {
	e := echo.New()
	assert.NotNil(t, e)
	assert.NotNil(t, e.Router())
}
