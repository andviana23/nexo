package handler

import (
	"encoding/json"
	"io"
	"net"
	"net/http"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	subUC "github.com/andviana23/barber-analytics-backend/internal/application/usecase/subscription"
	"github.com/andviana23/barber-analytics-backend/internal/infra/gateway/asaas"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// WebhookHandler handles incoming webhooks from Asaas
// Reference: FLUXO_ASSINATURA.md — Seção 6.6 (Fluxo Processar Webhook)
type WebhookHandler struct {
	processUC    *subUC.ProcessWebhookUseCase
	processUCV2  *subUC.ProcessWebhookUseCaseV2 // Use case v2 com suporte a log
	useV2        bool
	webhookToken string
	logger       *zap.Logger
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(
	processUC *subUC.ProcessWebhookUseCase,
	webhookToken string,
	logger *zap.Logger,
) *WebhookHandler {
	return &WebhookHandler{
		processUC:    processUC,
		webhookToken: webhookToken,
		logger:       logger,
		useV2:        false,
	}
}

// NewWebhookHandlerV2 creates a new webhook handler using V2 use case
func NewWebhookHandlerV2(
	processUCV2 *subUC.ProcessWebhookUseCaseV2,
	webhookToken string,
	logger *zap.Logger,
) *WebhookHandler {
	return &WebhookHandler{
		processUCV2:  processUCV2,
		webhookToken: webhookToken,
		logger:       logger,
		useV2:        true,
	}
}

// HandleAsaasWebhook processes incoming webhooks from Asaas
// POST /webhooks/asaas
// Reference: REGRA AS-005 - Webhook deve responder 200 em até 5s
func (h *WebhookHandler) HandleAsaasWebhook(c echo.Context) error {
	// Step 1: Validate webhook signature/token (AS-009)
	// Asaas sends token in header "asaas-access-token"
	token := c.Request().Header.Get("asaas-access-token")
	if token == "" {
		// Also check X-Asaas-Signature for newer API versions
		token = c.Request().Header.Get("X-Asaas-Signature")
	}

	if h.webhookToken != "" && token != h.webhookToken {
		h.logger.Warn("invalid webhook token",
			zap.String("received_token", maskToken(token)),
		)
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Invalid webhook token",
		})
	}

	// Step 2: Read request body
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		h.logger.Error("failed to read webhook body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Failed to read request body",
		})
	}

	// Step 3: Parse webhook event
	var event asaas.WebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		h.logger.Error("failed to parse webhook event",
			zap.Error(err),
			zap.String("body", string(body)),
		)
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Failed to parse webhook event",
		})
	}

	h.logger.Info("received Asaas webhook",
		zap.String("event", event.Event),
		zap.String("payment_id", getPaymentID(event)),
		zap.String("subscription_id", getSubscriptionID(event)),
	)

	// Step 4: Process webhook (AS-010)
	// Note: We respond 200 immediately and process async if needed
	// For now, we process synchronously but quickly
	ctx := c.Request().Context()

	if h.useV2 {
		// Extrair IP do cliente
		clientIP := net.ParseIP(c.RealIP())
		if err := h.processUCV2.Execute(ctx, event, body, clientIP); err != nil {
			h.logger.Error("failed to process webhook (v2)",
				zap.String("event", event.Event),
				zap.Error(err),
			)
			// Still return 200 to avoid Asaas retrying
		}
	} else {
		if err := h.processUC.Execute(ctx, event); err != nil {
			h.logger.Error("failed to process webhook",
				zap.String("event", event.Event),
				zap.Error(err),
			)
			// Still return 200 to avoid Asaas retrying
			// The error is logged for investigation
		}
	}

	// Step 5: Return 200 OK immediately (REGRA AS-005)
	return c.JSON(http.StatusOK, map[string]string{
		"status": "received",
	})
}

// maskToken masks the token for logging (shows first 4 and last 4 chars)
func maskToken(token string) string {
	if len(token) <= 8 {
		return "***"
	}
	return token[:4] + "****" + token[len(token)-4:]
}

// getPaymentID extracts payment ID from webhook event
func getPaymentID(event asaas.WebhookEvent) string {
	if event.Payment != nil {
		return event.Payment.ID
	}
	return ""
}

// getSubscriptionID extracts subscription ID from webhook event
func getSubscriptionID(event asaas.WebhookEvent) string {
	if event.Payment != nil {
		return event.Payment.Subscription
	}
	return ""
}
