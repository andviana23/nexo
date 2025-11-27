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

	// Erros de agendamento
	ErrAppointmentProfessionalRequired    = errors.New("profissional é obrigatório")
	ErrAppointmentCustomerRequired        = errors.New("cliente é obrigatório")
	ErrAppointmentStartTimeRequired       = errors.New("horário de início é obrigatório")
	ErrAppointmentServicesRequired        = errors.New("pelo menos um serviço é obrigatório")
	ErrAppointmentInvalidStatus           = errors.New("status de agendamento inválido")
	ErrAppointmentInvalidStatusTransition = errors.New("transição de status inválida")
	ErrAppointmentInvalidTimeRange        = errors.New("intervalo de horário inválido")
	ErrAppointmentConflict                = errors.New("conflito de horário com outro agendamento")
	ErrAppointmentCannotReschedule        = errors.New("não é possível reagendar este agendamento")
	ErrAppointmentNotFound                = errors.New("agendamento não encontrado")
	ErrAppointmentProfessionalNotFound    = errors.New("profissional não encontrado")
	ErrAppointmentCustomerNotFound        = errors.New("cliente não encontrado")
	ErrAppointmentServiceNotFound         = errors.New("serviço não encontrado")

	// Erros de cliente
	ErrCustomerNameRequired        = errors.New("nome do cliente é obrigatório")
	ErrCustomerNameTooShort        = errors.New("nome do cliente deve ter pelo menos 3 caracteres")
	ErrCustomerPhoneRequired       = errors.New("telefone do cliente é obrigatório")
	ErrCustomerPhoneInvalid        = errors.New("telefone inválido (deve ter 10 ou 11 dígitos)")
	ErrCustomerPhoneDuplicate      = errors.New("telefone já cadastrado para outro cliente")
	ErrCustomerEmailInvalid        = errors.New("email inválido")
	ErrCustomerEmailDuplicate      = errors.New("email já cadastrado para outro cliente")
	ErrCustomerCPFInvalid          = errors.New("CPF inválido (deve ter 11 dígitos)")
	ErrCustomerCPFDuplicate        = errors.New("CPF já cadastrado para outro cliente")
	ErrCustomerCEPInvalid          = errors.New("CEP inválido (deve ter 8 dígitos)")
	ErrCustomerStateInvalid        = errors.New("UF inválido (deve ter 2 caracteres)")
	ErrCustomerBirthDateFuture     = errors.New("data de nascimento não pode ser futura")
	ErrCustomerDateInvalid         = errors.New("formato de data inválido (esperado: YYYY-MM-DD)")
	ErrCustomerGenderInvalid       = errors.New("gênero inválido (M, F, NB ou PNI)")
	ErrCustomerObservationsTooLong = errors.New("observações devem ter no máximo 500 caracteres")
	ErrCustomerMaxTagsExceeded     = errors.New("máximo de 10 tags por cliente")
	ErrCustomerTagTooLong          = errors.New("tag deve ter no máximo 50 caracteres")
	ErrCustomerNotFound            = errors.New("cliente não encontrado")
	ErrCustomerInactive            = errors.New("cliente está inativo")

	// Erros de Lista da Vez (Barber Turn)
	ErrBarberTurnProfessionalRequired  = errors.New("professional_id é obrigatório")
	ErrBarberTurnProfessionalNotBarber = errors.New("apenas profissionais do tipo BARBEIRO podem ser adicionados à lista da vez")
	ErrBarberTurnProfessionalNotActive = errors.New("profissional deve estar ativo para participar da lista da vez")
	ErrBarberTurnProfessionalNotFound  = errors.New("profissional não encontrado na lista da vez")
	ErrBarberTurnAlreadyInList         = errors.New("profissional já está na lista da vez")
	ErrBarberTurnNotFound              = errors.New("registro não encontrado na lista da vez")
	ErrBarberTurnInvalidPoints         = errors.New("pontuação inválida (não pode ser negativa)")
	ErrBarberTurnCannotRecord          = errors.New("não é possível registrar atendimento para barbeiro pausado")
	ErrBarberTurnMonthYearInvalid      = errors.New("formato de mês/ano inválido (esperado: YYYY-MM)")
	ErrBarberTurnHistoryNotFound       = errors.New("histórico não encontrado para o período especificado")

	// Erros de Categoria de Serviço
	ErrCategoriaNomeRequired     = errors.New("nome da categoria é obrigatório")
	ErrCategoriaNomeTooShort     = errors.New("nome da categoria deve ter pelo menos 2 caracteres")
	ErrCategoriaNomeTooLong      = errors.New("nome da categoria deve ter no máximo 100 caracteres")
	ErrCategoriaNomeDuplicate    = errors.New("já existe uma categoria com este nome")
	ErrCategoriaCorInvalida      = errors.New("cor inválida (formato esperado: #RRGGBB)")
	ErrCategoriaNotFound         = errors.New("categoria não encontrada")
	ErrCategoriaHasServices      = errors.New("não é possível excluir categoria com serviços vinculados")
	ErrCategoriaDescricaoTooLong = errors.New("descrição deve ter no máximo 500 caracteres")
)
