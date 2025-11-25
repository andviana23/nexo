package domain

import "errors"

// Domain errors
var (
	// Erros gerais
	ErrTenantIDRequired = errors.New("tenant_id é obrigatório")
	ErrInvalidID        = errors.New("ID inválido")
	ErrInvalidTenantID  = errors.New("tenant_id inválido")
	ErrNotFound         = errors.New("recurso não encontrado")
	ErrAlreadyExists    = errors.New("recurso já existe")

	// Erros de MesAno
	ErrMesAnoRequired = errors.New("mes_ano é obrigatório")
	ErrMesAnoInvalido = errors.New("mes_ano inválido (formato esperado: YYYY-MM)")

	// Erros de valores monetários
	ErrValorInvalido = errors.New("valor inválido")
	ErrValorNegativo = errors.New("valor não pode ser negativo")
	ErrValorZero     = errors.New("valor não pode ser zero")

	// Erros de status
	ErrStatusInvalido = errors.New("status inválido")

	// Erros de compensação
	ErrCompensacaoJaCompensada = errors.New("compensação já está marcada como compensada")
	ErrDataCompensacaoInvalida = errors.New("data de compensação inválida")

	// Erros de contas
	ErrContaJaPaga            = errors.New("conta já está paga")
	ErrContaCancelada         = errors.New("conta está cancelada")
	ErrDataVencimentoInvalida = errors.New("data de vencimento inválida")

	// Erros de metas
	ErrMetaInvalida = errors.New("meta inválida")
	ErrMetaNegativa = errors.New("meta não pode ser negativa")

	// Erros de precificação
	ErrMargemInvalida = errors.New("margem inválida (deve estar entre 5-100%)")
	ErrMarkupInvalido = errors.New("markup inválido (deve ser >= 1)")

	// Erros de autenticação
	ErrEmailNaoEncontrado   = errors.New("Email não encontrado")
	ErrSenhaIncorreta       = errors.New("Senha incorreta")
	ErrContaDesativada      = errors.New("Conta desativada")
	ErrRefreshTokenInvalido = errors.New("Refresh token inválido ou expirado")
	ErrUsuarioNaoEncontrado = errors.New("Usuário não encontrado")
	ErrTokenInvalido        = errors.New("Token inválido")
)
