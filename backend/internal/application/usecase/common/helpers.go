package common

import (
	"github.com/google/uuid"
)

// ParseTenantID converte uma string de tenant ID para uuid.UUID
// Retorna erro se a string não for um UUID válido
func ParseTenantID(tenantID string) (uuid.UUID, error) {
	return uuid.Parse(tenantID)
}

// MustParseTenantID converte uma string de tenant ID para uuid.UUID
// Panic se a string não for um UUID válido (usar apenas quando já validado)
func MustParseTenantID(tenantID string) uuid.UUID {
	return uuid.MustParse(tenantID)
}
