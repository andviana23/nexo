package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/financial"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/andviana23/barber-analytics-backend/internal/infra/http/handler"
	"github.com/andviana23/barber-analytics-backend/internal/infra/repository/postgres"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupFinancialHandler() (*handler.FinancialHandler, *echo.Echo) {
	queries := db.New(testDBPool)

	// Repositories
	contaPagarRepo := postgres.NewContaPagarRepository(queries)
	contaReceberRepo := postgres.NewContaReceberRepository(queries)
	compensacaoRepo := postgres.NewCompensacaoBancariaRepository(queries)
	fluxoCaixaRepo := postgres.NewFluxoCaixaDiarioRepository(queries)
	dreRepo := postgres.NewDREMensalRepository(queries)

	// Use cases - ContaPagar
	createContaPagarUC := financial.NewCreateContaPagarUseCase(contaPagarRepo, testLogger)
	getContaPagarUC := financial.NewGetContaPagarUseCase(contaPagarRepo, testLogger)
	listContasPagarUC := financial.NewListContasPagarUseCase(contaPagarRepo, testLogger)
	updateContaPagarUC := financial.NewUpdateContaPagarUseCase(contaPagarRepo, testLogger)
	deleteContaPagarUC := financial.NewDeleteContaPagarUseCase(contaPagarRepo, testLogger)
	marcarPagamentoUC := financial.NewMarcarPagamentoUseCase(contaPagarRepo, testLogger)

	// Use cases - ContaReceber
	createContaReceberUC := financial.NewCreateContaReceberUseCase(contaReceberRepo, testLogger)
	getContaReceberUC := financial.NewGetContaReceberUseCase(contaReceberRepo, testLogger)
	listContasReceberUC := financial.NewListContasReceberUseCase(contaReceberRepo, testLogger)
	updateContaReceberUC := financial.NewUpdateContaReceberUseCase(contaReceberRepo, testLogger)
	deleteContaReceberUC := financial.NewDeleteContaReceberUseCase(contaReceberRepo, testLogger)
	marcarRecebimentoUC := financial.NewMarcarRecebimentoUseCase(contaReceberRepo, testLogger)

	// Use cases - Compensação
	createCompensacaoUC := financial.NewCreateCompensacaoUseCase(compensacaoRepo, testLogger)
	getCompensacaoUC := financial.NewGetCompensacaoUseCase(compensacaoRepo, testLogger)
	listCompensacoesUC := financial.NewListCompensacoesUseCase(compensacaoRepo, testLogger)
	deleteCompensacaoUC := financial.NewDeleteCompensacaoUseCase(compensacaoRepo, testLogger)
	marcarCompensacaoUC := financial.NewMarcarCompensacaoUseCase(compensacaoRepo, contaReceberRepo, testLogger)

	// Use cases - FluxoCaixa e DRE
	generateFluxoDiarioUC := financial.NewGenerateFluxoDiarioUseCase(fluxoCaixaRepo, contaPagarRepo, contaReceberRepo, compensacaoRepo, testLogger)
	getFluxoCaixaUC := financial.NewGetFluxoCaixaUseCase(fluxoCaixaRepo, testLogger)
	listFluxoCaixaUC := financial.NewListFluxoCaixaUseCase(fluxoCaixaRepo, testLogger)
	// DRE com commissionItemRepo nil (teste não precisa de comissões)
	generateDREUC := financial.NewGenerateDREUseCase(dreRepo, contaPagarRepo, contaReceberRepo, nil, testLogger)
	getDREUC := financial.NewGetDREUseCase(dreRepo, testLogger)
	listDREUC := financial.NewListDREUseCase(dreRepo, testLogger)

	// Handler
	financialHandler := handler.NewFinancialHandler(
		createContaPagarUC,
		getContaPagarUC,
		listContasPagarUC,
		updateContaPagarUC,
		deleteContaPagarUC,
		marcarPagamentoUC,
		createContaReceberUC,
		getContaReceberUC,
		listContasReceberUC,
		updateContaReceberUC,
		deleteContaReceberUC,
		marcarRecebimentoUC,
		createCompensacaoUC,
		getCompensacaoUC,
		listCompensacoesUC,
		deleteCompensacaoUC,
		marcarCompensacaoUC,
		generateFluxoDiarioUC,
		getFluxoCaixaUC,
		listFluxoCaixaUC,
		generateDREUC,
		getDREUC,
		listDREUC,
		nil, // getPainelMensalUC - not needed for these tests
		nil, // getProjecoesUC - not needed for these tests
		testLogger,
	)

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	return financialHandler, e
}

func TestFinancialHandler_CreateContaPagar_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste de integração")
	}

	financialHandler, e := setupFinancialHandler()

	payload := map[string]interface{}{
		"descricao":       "Aluguel Janeiro",
		"categoria_id":    "00000000-0000-0000-0000-000000000002",
		"fornecedor":      "Imobiliária XYZ",
		"valor":           "200000",
		"tipo":            "FIXO",
		"data_vencimento": "2026-01-10",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/financial/payables", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("X-Tenant-ID", "e2e00000-0000-0000-0000-000000000001")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("tenant_id", "e2e00000-0000-0000-0000-000000000001")

	err := financialHandler.CreateContaPagar(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotEmpty(t, response["id"])
	assert.Equal(t, "Aluguel Janeiro", response["descricao"])
	assert.Equal(t, "200000", response["valor"])
}

func TestFinancialHandler_ListContasPagar_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste de integração")
	}

	financialHandler, e := setupFinancialHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/financial/payables?page=1&page_size=10", nil)
	req.Header.Set("X-Tenant-ID", "e2e00000-0000-0000-0000-000000000001")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("tenant_id", "e2e00000-0000-0000-0000-000000000001")

	err := financialHandler.ListContasPagar(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotNil(t, response["data"])
}
