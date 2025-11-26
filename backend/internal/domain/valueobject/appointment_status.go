package valueobject

// AppointmentStatus representa o status de um agendamento
type AppointmentStatus string

const (
	AppointmentStatusCreated   AppointmentStatus = "CREATED"
	AppointmentStatusConfirmed AppointmentStatus = "CONFIRMED"
	AppointmentStatusInService AppointmentStatus = "IN_SERVICE"
	AppointmentStatusDone      AppointmentStatus = "DONE"
	AppointmentStatusNoShow    AppointmentStatus = "NO_SHOW"
	AppointmentStatusCanceled  AppointmentStatus = "CANCELED"
)

// IsValid verifica se o status é válido
func (s AppointmentStatus) IsValid() bool {
	switch s {
	case AppointmentStatusCreated,
		AppointmentStatusConfirmed,
		AppointmentStatusInService,
		AppointmentStatusDone,
		AppointmentStatusNoShow,
		AppointmentStatusCanceled:
		return true
	}
	return false
}

// String retorna a string do status
func (s AppointmentStatus) String() string {
	return string(s)
}

// CanTransitionTo verifica se a transição de status é válida
// Fluxo normal: CREATED -> CONFIRMED -> IN_SERVICE -> DONE
// A qualquer momento (antes de DONE): pode cancelar ou marcar no-show
func (s AppointmentStatus) CanTransitionTo(newStatus AppointmentStatus) bool {
	switch s {
	case AppointmentStatusCreated:
		return newStatus == AppointmentStatusConfirmed ||
			newStatus == AppointmentStatusCanceled ||
			newStatus == AppointmentStatusNoShow

	case AppointmentStatusConfirmed:
		return newStatus == AppointmentStatusInService ||
			newStatus == AppointmentStatusCanceled ||
			newStatus == AppointmentStatusNoShow

	case AppointmentStatusInService:
		return newStatus == AppointmentStatusDone ||
			newStatus == AppointmentStatusCanceled

	case AppointmentStatusDone, AppointmentStatusNoShow, AppointmentStatusCanceled:
		// Estados finais - não permitem transição
		return false
	}
	return false
}

// IsFinal verifica se é um status final (não pode mudar)
func (s AppointmentStatus) IsFinal() bool {
	return s == AppointmentStatusDone ||
		s == AppointmentStatusNoShow ||
		s == AppointmentStatusCanceled
}

// IsActive verifica se o agendamento está ativo
func (s AppointmentStatus) IsActive() bool {
	return s == AppointmentStatusCreated ||
		s == AppointmentStatusConfirmed ||
		s == AppointmentStatusInService
}

// DisplayName retorna o nome amigável em português
func (s AppointmentStatus) DisplayName() string {
	switch s {
	case AppointmentStatusCreated:
		return "Criado"
	case AppointmentStatusConfirmed:
		return "Confirmado"
	case AppointmentStatusInService:
		return "Em Atendimento"
	case AppointmentStatusDone:
		return "Concluído"
	case AppointmentStatusNoShow:
		return "Não Compareceu"
	case AppointmentStatusCanceled:
		return "Cancelado"
	default:
		return "Desconhecido"
	}
}

// Color retorna a cor associada ao status (para UI)
func (s AppointmentStatus) Color() string {
	switch s {
	case AppointmentStatusCreated:
		return "#3B82F6" // blue-500
	case AppointmentStatusConfirmed:
		return "#10B981" // emerald-500
	case AppointmentStatusInService:
		return "#F59E0B" // amber-500
	case AppointmentStatusDone:
		return "#22C55E" // green-500
	case AppointmentStatusNoShow:
		return "#EF4444" // red-500
	case AppointmentStatusCanceled:
		return "#6B7280" // gray-500
	default:
		return "#9CA3AF" // gray-400
	}
}

// ParseAppointmentStatus converte string para AppointmentStatus
func ParseAppointmentStatus(s string) (AppointmentStatus, bool) {
	status := AppointmentStatus(s)
	return status, status.IsValid()
}
