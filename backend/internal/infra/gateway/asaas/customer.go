// Package asaas - Customer API methods
// Reference: https://docs.asaas.com/reference/criar-novo-cliente
// REGRAS: AS-001, AS-002, AS-003, RN-CLI-002
package asaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"go.uber.org/zap"
)

// ============================================================================
// CUSTOMER METHODS
// ============================================================================

// FindCustomerByNameAndPhone searches for a customer by name and mobile phone
// Reference: REGRA AS-001 - Busca de cliente no Asaas é SEMPRE por Nome + Telefone
func (c *Client) FindCustomerByNameAndPhone(ctx context.Context, name, phone string) (*CustomerResponse, error) {
	// Clean phone number (remove non-digits)
	cleanPhone := cleanPhoneNumber(phone)

	// Build query params
	params := url.Values{}
	params.Set("name", name)
	params.Set("mobilePhone", cleanPhone)

	path := "/customers?" + params.Encode()

	c.logger.Debug("searching customer by name and phone",
		zap.String("name", name),
		zap.String("phone", cleanPhone),
	)

	body, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("find customer by name and phone: %w", err)
	}

	var listResp CustomerListResponse
	if err := json.Unmarshal(body, &listResp); err != nil {
		return nil, fmt.Errorf("unmarshal customer list: %w", err)
	}

	if listResp.TotalCount == 0 || len(listResp.Data) == 0 {
		return nil, nil // Not found
	}

	// Return first match
	return &listResp.Data[0], nil
}

// CreateCustomer creates a new customer in Asaas
// Reference: REGRA AS-002, AS-003 - CPF is NOT required
func (c *Client) CreateCustomer(ctx context.Context, req CustomerRequest) (*CustomerResponse, error) {
	c.logger.Debug("creating customer in Asaas",
		zap.String("name", req.Name),
		zap.String("externalReference", req.ExternalReference),
	)

	body, err := c.doRequest(ctx, "POST", "/customers", req)
	if err != nil {
		return nil, fmt.Errorf("create customer: %w", err)
	}

	var resp CustomerResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal customer response: %w", err)
	}

	c.logger.Info("customer created in Asaas",
		zap.String("asaas_customer_id", resp.ID),
		zap.String("name", resp.Name),
	)

	return &resp, nil
}

// GetCustomer retrieves a customer by ID
func (c *Client) GetCustomer(ctx context.Context, customerID string) (*CustomerResponse, error) {
	path := fmt.Sprintf("/customers/%s", customerID)

	body, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("get customer: %w", err)
	}

	var resp CustomerResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal customer: %w", err)
	}

	return &resp, nil
}

// FindOrCreateCustomer finds a customer by name+phone or creates if not exists
// This implements the logic from FLUXO_ASSINATURA.md Section 6.1:
// 1. Check if local client already has asaas_customer_id (caller should handle)
// 2. Search in Asaas by name + phone (REGRA AS-001)
// 3. If found: return existing ID (avoid duplication - RN-CLI-002)
// 4. If not found: create new customer (without CPF - AS-002, AS-003)
func (c *Client) FindOrCreateCustomer(ctx context.Context, req CustomerRequest) (*CustomerResponse, bool, error) {
	// Try to find existing customer
	existing, err := c.FindCustomerByNameAndPhone(ctx, req.Name, req.MobilePhone)
	if err != nil {
		c.logger.Warn("error searching customer, will try to create",
			zap.String("name", req.Name),
			zap.Error(err),
		)
	}

	if existing != nil {
		c.logger.Info("customer already exists in Asaas",
			zap.String("asaas_customer_id", existing.ID),
			zap.String("name", existing.Name),
		)
		return existing, false, nil // false = was not created (already existed)
	}

	// Create new customer
	created, err := c.CreateCustomer(ctx, req)
	if err != nil {
		return nil, false, fmt.Errorf("create customer: %w", err)
	}

	return created, true, nil // true = was created
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// cleanPhoneNumber removes all non-digit characters from phone
func cleanPhoneNumber(phone string) string {
	var result []byte
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			result = append(result, byte(ch))
		}
	}
	return string(result)
}
