// Package postgres implementa os repositórios usando PostgreSQL e sqlc.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
)

// CompensacaoBancariaRepository implementa port.CompensacaoBancariaRepository usando sqlc.
type CompensacaoBancariaRepository struct {
	queries *db.Queries
}

// NewCompensacaoBancariaRepository cria uma nova instância do repositório.
func NewCompensacaoBancariaRepository(queries *db.Queries) *CompensacaoBancariaRepository {
	return &CompensacaoBancariaRepository{
		queries: queries,
	}
}

// Create persiste uma nova compensação bancária.
func (r *CompensacaoBancariaRepository) Create(ctx context.Context, comp *entity.CompensacaoBancaria) error {
	tenantUUID := entityUUIDToPgtype(comp.TenantID)
	receitaUUID := uuidStringToPgtype(comp.ReceitaID)
	meioPagamentoUUID := uuidStringToPgtype(comp.MeioPagamentoID)

	statusStr := comp.Status.String()

	params := db.CreateCompensacaoBancariaParams{
		TenantID:        tenantUUID,
		ReceitaID:       receitaUUID,
		DataTransacao:   dateToDate(comp.DataTransacao),
		DataCompensacao: dateToDate(comp.DataCompensacao),
		DataCompensado:  timePtrToDate(comp.DataCompensado),
		ValorBruto:      moneyToRawDecimal(comp.ValorBruto),
		TaxaPercentual:  percentageToNumeric(comp.TaxaPercentual),
		TaxaFixa:        moneyToNumeric(comp.TaxaFixa),
		ValorLiquido:    moneyToRawDecimal(comp.ValorLiquido),
		MeioPagamentoID: meioPagamentoUUID,
		DMais:           dmaisToInt32(comp.DMais),
		Status:          &statusStr,
	}

	result, err := r.queries.CreateCompensacaoBancaria(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar compensação bancária: %w", err)
	}

	// Atualizar ID da entidade
	comp.ID = pgUUIDToString(result.ID)
	comp.CriadoEm = timestamptzToTime(result.CriadoEm)
	comp.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)

	return nil
}

// FindByID busca uma compensação por ID.
func (r *CompensacaoBancariaRepository) FindByID(ctx context.Context, tenantID, id string) (*entity.CompensacaoBancaria, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	idUUID := uuidStringToPgtype(id)

	result, err := r.queries.GetCompensacaoBancariaByID(ctx, db.GetCompensacaoBancariaByIDParams{
		ID:       idUUID,
		TenantID: tenantUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar compensação: %w", err)
	}

	return r.toDomain(&result)
}

// FindByReceitaID busca compensação de uma receita.
func (r *CompensacaoBancariaRepository) FindByReceitaID(ctx context.Context, tenantID, receitaID string) (*entity.CompensacaoBancaria, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	receitaUUID := uuidStringToPgtype(receitaID)

	results, err := r.queries.ListCompensacoesByReceita(ctx, db.ListCompensacoesByReceitaParams{
		TenantID:  tenantUUID,
		ReceitaID: receitaUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar compensações por receita: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("compensação não encontrada para receita %s", receitaID)
	}

	return r.toDomain(&results[0])
}

// Update atualiza uma compensação existente.
func (r *CompensacaoBancariaRepository) Update(ctx context.Context, comp *entity.CompensacaoBancaria) error {
	tenantUUID := entityUUIDToPgtype(comp.TenantID)
	idUUID := uuidStringToPgtype(comp.ID)

	statusStr := comp.Status.String()

	params := db.UpdateCompensacaoBancariaParams{
		ID:              idUUID,
		TenantID:        tenantUUID,
		DataCompensacao: dateToDate(comp.DataCompensacao),
		DataCompensado:  timePtrToDate(comp.DataCompensado),
		ValorBruto:      moneyToRawDecimal(comp.ValorBruto),
		TaxaPercentual:  percentageToNumeric(comp.TaxaPercentual),
		TaxaFixa:        moneyToNumeric(comp.TaxaFixa),
		ValorLiquido:    moneyToRawDecimal(comp.ValorLiquido),
		Status:          &statusStr,
	}

	result, err := r.queries.UpdateCompensacaoBancaria(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao atualizar compensação: %w", err)
	}

	comp.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)
	return nil
}

// Delete remove uma compensação.
func (r *CompensacaoBancariaRepository) Delete(ctx context.Context, tenantID, id string) error {
	tenantUUID := uuidStringToPgtype(tenantID)
	idUUID := uuidStringToPgtype(id)

	err := r.queries.DeleteCompensacaoBancaria(ctx, db.DeleteCompensacaoBancariaParams{
		ID:       idUUID,
		TenantID: tenantUUID,
	})
	if err != nil {
		return fmt.Errorf("erro ao deletar compensação: %w", err)
	}

	return nil
}

// List lista compensações com filtros.
func (r *CompensacaoBancariaRepository) List(ctx context.Context, tenantID string, filters port.CompensacaoListFilters) ([]*entity.CompensacaoBancaria, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	pageSize := int32(filters.PageSize)
	if pageSize <= 0 {
		pageSize = 50
	}
	offset := int32(filters.Page * filters.PageSize)

	var results []db.CompensacoesBancaria
	var err error

	if filters.Status != nil {
		statusStr := filters.Status.String()
		results, err = r.queries.ListCompensacoesByStatus(ctx, db.ListCompensacoesByStatusParams{
			TenantID: tenantUUID,
			Status:   &statusStr,
			Limit:    pageSize,
			Offset:   offset,
		})
	} else {
		results, err = r.queries.ListCompensacoesBancariasByTenant(ctx, db.ListCompensacoesBancariasByTenantParams{
			TenantID: tenantUUID,
			Limit:    pageSize,
			Offset:   offset,
		})
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao listar compensações: %w", err)
	}

	return r.toDomainList(results)
}

// ListByStatus lista compensações por status.
func (r *CompensacaoBancariaRepository) ListByStatus(ctx context.Context, tenantID string, status valueobject.StatusCompensacao) ([]*entity.CompensacaoBancaria, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	statusStr := status.String()
	results, err := r.queries.ListCompensacoesByStatus(ctx, db.ListCompensacoesByStatusParams{
		TenantID: tenantUUID,
		Status:   &statusStr,
		Limit:    1000,
		Offset:   0,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar compensações por status: %w", err)
	}

	return r.toDomainList(results)
}

// ListPendentesCompensacao lista compensações pendentes de compensação (data <= hoje).
func (r *CompensacaoBancariaRepository) ListPendentesCompensacao(ctx context.Context, tenantID string) ([]*entity.CompensacaoBancaria, error) {
	statuses := []valueobject.StatusCompensacao{
		valueobject.StatusCompensacaoPrevisto,
		valueobject.StatusCompensacaoConfirmado,
	}

	// Buscar compensações previstas ou confirmadas
	var compensacoes []*entity.CompensacaoBancaria
	for _, st := range statuses {
		items, err := r.ListByStatus(ctx, tenantID, st)
		if err != nil {
			return nil, err
		}
		compensacoes = append(compensacoes, items...)
	}

	// Filtrar apenas as que têm data_compensacao <= hoje
	hoje := time.Now()
	var pendentes []*entity.CompensacaoBancaria
	for _, comp := range compensacoes {
		if !comp.DataCompensacao.After(hoje) {
			pendentes = append(pendentes, comp)
		}
	}

	return pendentes, nil
}

// ListByDateRange lista compensações em um período.
func (r *CompensacaoBancariaRepository) ListByDateRange(ctx context.Context, tenantID string, inicio, fim time.Time) ([]*entity.CompensacaoBancaria, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	results, err := r.queries.ListCompensacoesByDataCompensacao(ctx, db.ListCompensacoesByDataCompensacaoParams{
		TenantID:          tenantUUID,
		DataCompensacao:   dateToDate(inicio),
		DataCompensacao_2: dateToDate(fim),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar compensações por período: %w", err)
	}

	return r.toDomainList(results)
}

// toDomain converte modelo sqlc para entidade de domínio.
func (r *CompensacaoBancariaRepository) toDomain(model *db.CompensacoesBancaria) (*entity.CompensacaoBancaria, error) {
	var status valueobject.StatusCompensacao
	if model.Status != nil {
		status = valueobject.StatusCompensacao(*model.Status)
	} else {
		status = valueobject.StatusCompensacaoPrevisto
	}

	if !status.IsValid() {
		return nil, fmt.Errorf("status inválido: %s", status)
	}

	taxaPercentual, err := numericToPercentage(model.TaxaPercentual)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter taxa percentual: %w", err)
	}

	comp := &entity.CompensacaoBancaria{
		ID:              pgUUIDToString(model.ID),
		TenantID:        pgtypeToEntityUUID(model.TenantID),
		ReceitaID:       pgUUIDToString(model.ReceitaID),
		DataTransacao:   dateToTime(model.DataTransacao),
		DataCompensacao: dateToTime(model.DataCompensacao),
		DataCompensado:  dateToTimePtr(model.DataCompensado),
		ValorBruto:      rawDecimalToMoney(model.ValorBruto),
		TaxaPercentual:  taxaPercentual,
		TaxaFixa:        numericToMoney(model.TaxaFixa),
		ValorLiquido:    rawDecimalToMoney(model.ValorLiquido),
		MeioPagamentoID: pgUUIDToString(model.MeioPagamentoID),
		DMais:           int32ToDMais(model.DMais),
		Status:          status,
		CriadoEm:        timestamptzToTime(model.CriadoEm),
		AtualizadoEm:    timestamptzToTime(model.AtualizadoEm),
	}

	return comp, nil
}

// toDomainList converte lista de modelos para lista de entidades.
func (r *CompensacaoBancariaRepository) toDomainList(models []db.CompensacoesBancaria) ([]*entity.CompensacaoBancaria, error) {
	result := make([]*entity.CompensacaoBancaria, 0, len(models))
	for i := range models {
		comp, err := r.toDomain(&models[i])
		if err != nil {
			return nil, err
		}
		result = append(result, comp)
	}
	return result, nil
}
