// Package customer contém os use cases de clientes do NEXO.
// Implementa operações CRUD, busca, exportação LGPD e validações de negócio.
package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// =============================================================================
// CreateCustomerUseCase
// =============================================================================

// CreateCustomerUseCase cria um novo cliente
type CreateCustomerUseCase struct {
	repo   port.CustomerRepository
	logger *zap.Logger
}

// NewCreateCustomerUseCase cria uma nova instância do use case
func NewCreateCustomerUseCase(repo port.CustomerRepository, logger *zap.Logger) *CreateCustomerUseCase {
	return &CreateCustomerUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute cria um novo cliente
func (uc *CreateCustomerUseCase) Execute(ctx context.Context, tenantID string, req dto.CreateCustomerRequest) (*entity.Customer, error) {
	// Verificar duplicidade de telefone
	exists, err := uc.repo.CheckPhoneExists(ctx, tenantID, req.Telefone, nil)
	if err != nil {
		uc.logger.Error("erro ao verificar telefone duplicado", zap.Error(err))
		return nil, err
	}
	if exists {
		return nil, domain.ErrCustomerPhoneDuplicate
	}

	// Verificar duplicidade de CPF se informado
	if req.CPF != nil && *req.CPF != "" {
		cpfExists, err := uc.repo.CheckCPFExists(ctx, tenantID, *req.CPF, nil)
		if err != nil {
			uc.logger.Error("erro ao verificar CPF duplicado", zap.Error(err))
			return nil, err
		}
		if cpfExists {
			return nil, domain.ErrCustomerCPFDuplicate
		}
	}

	// Converter tenant_id de string para uuid.UUID
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant_id inválido: %w", err)
	}

	// Criar entidade
	customer, err := entity.NewCustomer(tenantUUID, req.Nome, req.Telefone)
	if err != nil {
		return nil, err
	}

	// Campos opcionais
	if req.Email != nil {
		if err := customer.SetEmail(*req.Email); err != nil {
			return nil, err
		}
	}

	if req.CPF != nil {
		if err := customer.SetCPF(*req.CPF); err != nil {
			return nil, err
		}
	}

	if req.DataNascimento != nil {
		dataNasc, err := time.Parse("2006-01-02", *req.DataNascimento)
		if err != nil {
			return nil, domain.ErrCustomerDateInvalid
		}
		customer.DataNascimento = &dataNasc
	}

	if req.Genero != nil {
		customer.Genero = req.Genero
	}

	// Endereço
	if req.EnderecoLogradouro != nil {
		if err := customer.SetEndereco(
			req.EnderecoLogradouro,
			req.EnderecoNumero,
			req.EnderecoComplemento,
			req.EnderecoBairro,
			req.EnderecoCidade,
			req.EnderecoEstado,
			req.EnderecoCEP,
		); err != nil {
			return nil, err
		}
	}

	if req.Observacoes != nil {
		if err := customer.SetObservacoes(*req.Observacoes); err != nil {
			return nil, err
		}
	}

	// Tags
	if len(req.Tags) > 0 {
		for _, tag := range req.Tags {
			if err := customer.AddTag(tag); err != nil {
				return nil, err
			}
		}
	}

	// Persistir
	if err := uc.repo.Create(ctx, customer); err != nil {
		uc.logger.Error("erro ao criar cliente", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("cliente criado com sucesso",
		zap.String("customer_id", customer.ID),
		zap.String("tenant_id", tenantID),
	)

	return customer, nil
}

// =============================================================================
// UpdateCustomerUseCase
// =============================================================================

// UpdateCustomerUseCase atualiza um cliente existente
type UpdateCustomerUseCase struct {
	repo   port.CustomerRepository
	logger *zap.Logger
}

// NewUpdateCustomerUseCase cria uma nova instância do use case
func NewUpdateCustomerUseCase(repo port.CustomerRepository, logger *zap.Logger) *UpdateCustomerUseCase {
	return &UpdateCustomerUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute atualiza um cliente existente
func (uc *UpdateCustomerUseCase) Execute(ctx context.Context, tenantID, customerID string, req dto.UpdateCustomerRequest) (*entity.Customer, error) {
	// Buscar cliente existente
	customer, err := uc.repo.FindByID(ctx, tenantID, customerID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, domain.ErrCustomerNotFound
	}

	// Verificar duplicidade de telefone
	if req.Telefone != nil {
		exists, err := uc.repo.CheckPhoneExists(ctx, tenantID, *req.Telefone, &customerID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domain.ErrCustomerPhoneDuplicate
		}
		if err := customer.UpdateTelefone(*req.Telefone); err != nil {
			return nil, err
		}
	}

	// Verificar duplicidade de CPF
	if req.CPF != nil && *req.CPF != "" {
		exists, err := uc.repo.CheckCPFExists(ctx, tenantID, *req.CPF, &customerID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domain.ErrCustomerCPFDuplicate
		}
		if err := customer.SetCPF(*req.CPF); err != nil {
			return nil, err
		}
	}

	// Atualizar campos
	if req.Nome != nil {
		if err := customer.UpdateNome(*req.Nome); err != nil {
			return nil, err
		}
	}

	if req.Email != nil {
		if err := customer.SetEmail(*req.Email); err != nil {
			return nil, err
		}
	}

	if req.DataNascimento != nil {
		dataNasc, err := time.Parse("2006-01-02", *req.DataNascimento)
		if err != nil {
			return nil, domain.ErrCustomerDateInvalid
		}
		if err := customer.SetDataNascimento(&dataNasc); err != nil {
			return nil, err
		}
	}

	if req.Genero != nil {
		if err := customer.SetGenero(*req.Genero); err != nil {
			return nil, err
		}
	}

	// Endereço
	if req.EnderecoLogradouro != nil || req.EnderecoNumero != nil ||
		req.EnderecoComplemento != nil || req.EnderecoBairro != nil ||
		req.EnderecoCidade != nil || req.EnderecoEstado != nil ||
		req.EnderecoCEP != nil {
		if err := customer.SetEndereco(
			req.EnderecoLogradouro,
			req.EnderecoNumero,
			req.EnderecoComplemento,
			req.EnderecoBairro,
			req.EnderecoCidade,
			req.EnderecoEstado,
			req.EnderecoCEP,
		); err != nil {
			return nil, err
		}
	}

	if req.Observacoes != nil {
		if err := customer.SetObservacoes(*req.Observacoes); err != nil {
			return nil, err
		}
	}

	if req.Tags != nil {
		if err := customer.SetTags(req.Tags); err != nil {
			return nil, err
		}
	}

	customer.UpdatedAt = time.Now()

	// Persistir
	if err := uc.repo.Update(ctx, customer); err != nil {
		uc.logger.Error("erro ao atualizar cliente", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("cliente atualizado com sucesso",
		zap.String("customer_id", customer.ID),
		zap.String("tenant_id", tenantID),
	)

	return customer, nil
}

// =============================================================================
// ListCustomersUseCase
// =============================================================================

// ListCustomersUseCase lista clientes com paginação
type ListCustomersUseCase struct {
	repo   port.CustomerRepository
	logger *zap.Logger
}

// NewListCustomersUseCase cria uma nova instância do use case
func NewListCustomersUseCase(repo port.CustomerRepository, logger *zap.Logger) *ListCustomersUseCase {
	return &ListCustomersUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute lista clientes com filtros
func (uc *ListCustomersUseCase) Execute(ctx context.Context, tenantID string, filter port.CustomerFilter) ([]*entity.Customer, int64, error) {
	customers, total, err := uc.repo.List(ctx, tenantID, filter)
	if err != nil {
		uc.logger.Error("erro ao listar clientes", zap.Error(err))
		return nil, 0, err
	}

	return customers, total, nil
}

// =============================================================================
// GetCustomerUseCase
// =============================================================================

// GetCustomerUseCase busca um cliente por ID
type GetCustomerUseCase struct {
	repo   port.CustomerRepository
	logger *zap.Logger
}

// NewGetCustomerUseCase cria uma nova instância do use case
func NewGetCustomerUseCase(repo port.CustomerRepository, logger *zap.Logger) *GetCustomerUseCase {
	return &GetCustomerUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute busca um cliente por ID
func (uc *GetCustomerUseCase) Execute(ctx context.Context, tenantID, customerID string) (*entity.Customer, error) {
	customer, err := uc.repo.FindByID(ctx, tenantID, customerID)
	if err != nil {
		uc.logger.Error("erro ao buscar cliente", zap.Error(err))
		return nil, err
	}
	if customer == nil {
		return nil, domain.ErrCustomerNotFound
	}

	return customer, nil
}

// =============================================================================
// GetCustomerWithHistoryUseCase
// =============================================================================

// GetCustomerWithHistoryUseCase busca um cliente com histórico de atendimentos
type GetCustomerWithHistoryUseCase struct {
	repo   port.CustomerRepository
	logger *zap.Logger
}

// NewGetCustomerWithHistoryUseCase cria uma nova instância do use case
func NewGetCustomerWithHistoryUseCase(repo port.CustomerRepository, logger *zap.Logger) *GetCustomerWithHistoryUseCase {
	return &GetCustomerWithHistoryUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute busca um cliente com histórico completo
func (uc *GetCustomerWithHistoryUseCase) Execute(ctx context.Context, tenantID, customerID string) (*port.CustomerWithHistory, error) {
	cwh, err := uc.repo.GetWithHistory(ctx, tenantID, customerID)
	if err != nil {
		uc.logger.Error("erro ao buscar cliente com histórico", zap.Error(err))
		return nil, err
	}
	if cwh == nil {
		return nil, domain.ErrCustomerNotFound
	}

	return cwh, nil
}

// =============================================================================
// InactivateCustomerUseCase
// =============================================================================

// InactivateCustomerUseCase inativa um cliente
type InactivateCustomerUseCase struct {
	repo   port.CustomerRepository
	logger *zap.Logger
}

// NewInactivateCustomerUseCase cria uma nova instância do use case
func NewInactivateCustomerUseCase(repo port.CustomerRepository, logger *zap.Logger) *InactivateCustomerUseCase {
	return &InactivateCustomerUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute inativa um cliente (soft delete)
func (uc *InactivateCustomerUseCase) Execute(ctx context.Context, tenantID, customerID string) error {
	// Verificar se existe
	customer, err := uc.repo.FindByID(ctx, tenantID, customerID)
	if err != nil {
		return err
	}
	if customer == nil {
		return domain.ErrCustomerNotFound
	}

	if err := uc.repo.Inactivate(ctx, tenantID, customerID); err != nil {
		uc.logger.Error("erro ao inativar cliente", zap.Error(err))
		return err
	}

	uc.logger.Info("cliente inativado com sucesso",
		zap.String("customer_id", customerID),
		zap.String("tenant_id", tenantID),
	)

	return nil
}

// =============================================================================
// SearchCustomersUseCase
// =============================================================================

// SearchCustomersUseCase busca clientes por termo
type SearchCustomersUseCase struct {
	repo   port.CustomerRepository
	logger *zap.Logger
}

// NewSearchCustomersUseCase cria uma nova instância do use case
func NewSearchCustomersUseCase(repo port.CustomerRepository, logger *zap.Logger) *SearchCustomersUseCase {
	return &SearchCustomersUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute busca clientes por nome, telefone ou email
func (uc *SearchCustomersUseCase) Execute(ctx context.Context, tenantID, term string) ([]*port.CustomerSummary, error) {
	customers, err := uc.repo.Search(ctx, tenantID, term)
	if err != nil {
		uc.logger.Error("erro ao buscar clientes", zap.Error(err))
		return nil, err
	}

	return customers, nil
}

// =============================================================================
// ExportCustomerDataUseCase
// =============================================================================

// ExportCustomerDataUseCase exporta todos os dados de um cliente (LGPD)
type ExportCustomerDataUseCase struct {
	repo   port.CustomerRepository
	logger *zap.Logger
}

// NewExportCustomerDataUseCase cria uma nova instância do use case
func NewExportCustomerDataUseCase(repo port.CustomerRepository, logger *zap.Logger) *ExportCustomerDataUseCase {
	return &ExportCustomerDataUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute exporta todos os dados do cliente para LGPD
func (uc *ExportCustomerDataUseCase) Execute(ctx context.Context, tenantID, customerID string) (*port.CustomerExport, error) {
	export, err := uc.repo.GetDataForExport(ctx, tenantID, customerID)
	if err != nil {
		uc.logger.Error("erro ao exportar dados do cliente", zap.Error(err))
		return nil, err
	}
	if export == nil {
		return nil, domain.ErrCustomerNotFound
	}

	uc.logger.Info("dados do cliente exportados (LGPD)",
		zap.String("customer_id", customerID),
		zap.String("tenant_id", tenantID),
	)

	return export, nil
}

// =============================================================================
// GetCustomerStatsUseCase
// =============================================================================

// GetCustomerStatsUseCase retorna estatísticas de clientes
type GetCustomerStatsUseCase struct {
	repo   port.CustomerRepository
	logger *zap.Logger
}

// NewGetCustomerStatsUseCase cria uma nova instância do use case
func NewGetCustomerStatsUseCase(repo port.CustomerRepository, logger *zap.Logger) *GetCustomerStatsUseCase {
	return &GetCustomerStatsUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute retorna estatísticas dos clientes
func (uc *GetCustomerStatsUseCase) Execute(ctx context.Context, tenantID string) (*port.CustomerStats, error) {
	stats, err := uc.repo.GetStats(ctx, tenantID)
	if err != nil {
		uc.logger.Error("erro ao obter estatísticas de clientes", zap.Error(err))
		return nil, err
	}

	return stats, nil
}

// =============================================================================
// CheckPhoneDuplicateUseCase
// =============================================================================

// CheckPhoneDuplicateUseCase verifica se telefone já está cadastrado
type CheckPhoneDuplicateUseCase struct {
	repo   port.CustomerRepository
	logger *zap.Logger
}

// NewCheckPhoneDuplicateUseCase cria uma nova instância do use case
func NewCheckPhoneDuplicateUseCase(repo port.CustomerRepository, logger *zap.Logger) *CheckPhoneDuplicateUseCase {
	return &CheckPhoneDuplicateUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute verifica duplicidade de telefone
func (uc *CheckPhoneDuplicateUseCase) Execute(ctx context.Context, tenantID, phone string, excludeID *string) (bool, error) {
	return uc.repo.CheckPhoneExists(ctx, tenantID, phone, excludeID)
}

// =============================================================================
// CheckCPFDuplicateUseCase
// =============================================================================

// CheckCPFDuplicateUseCase verifica se CPF já está cadastrado
type CheckCPFDuplicateUseCase struct {
	repo   port.CustomerRepository
	logger *zap.Logger
}

// NewCheckCPFDuplicateUseCase cria uma nova instância do use case
func NewCheckCPFDuplicateUseCase(repo port.CustomerRepository, logger *zap.Logger) *CheckCPFDuplicateUseCase {
	return &CheckCPFDuplicateUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute verifica duplicidade de CPF
func (uc *CheckCPFDuplicateUseCase) Execute(ctx context.Context, tenantID, cpf string, excludeID *string) (bool, error) {
	return uc.repo.CheckCPFExists(ctx, tenantID, cpf, excludeID)
}

// =============================================================================
// ListActiveCustomersUseCase
// =============================================================================

// ListActiveCustomersUseCase lista clientes ativos para autocomplete
type ListActiveCustomersUseCase struct {
	repo   port.CustomerRepository
	logger *zap.Logger
}

// NewListActiveCustomersUseCase cria uma nova instância do use case
func NewListActiveCustomersUseCase(repo port.CustomerRepository, logger *zap.Logger) *ListActiveCustomersUseCase {
	return &ListActiveCustomersUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute lista clientes ativos para selects/autocomplete
func (uc *ListActiveCustomersUseCase) Execute(ctx context.Context, tenantID string) ([]*port.CustomerSummary, error) {
	customers, err := uc.repo.ListActive(ctx, tenantID)
	if err != nil {
		uc.logger.Error("erro ao listar clientes ativos", zap.Error(err))
		return nil, err
	}

	return customers, nil
}
