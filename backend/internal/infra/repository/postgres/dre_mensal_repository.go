// Package postgres implementa os repositórios usando PostgreSQL e sqlc.
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
)

// DREMensalRepository implementa port.DREMensalRepository usando sqlc.
type DREMensalRepository struct {
	queries *db.Queries
}

// NewDREMensalRepository cria uma nova instância do repositório de DRE Mensal.
func NewDREMensalRepository(queries *db.Queries) *DREMensalRepository {
	return &DREMensalRepository{
		queries: queries,
	}
}

// Create persiste uma nova DRE Mensal no banco de dados.
func (r *DREMensalRepository) Create(ctx context.Context, dre *entity.DREMensal) error {
	tenantUUID := uuidStringToPgtype(dre.TenantID)

	params := db.CreateDREMensalParams{
		TenantID:             tenantUUID,
		MesAno:               dre.MesAno.String(),
		ReceitaServicos:      moneyToNumeric(dre.ReceitaServicos),
		ReceitaProdutos:      moneyToNumeric(dre.ReceitaProdutos),
		ReceitaPlanos:        moneyToNumeric(dre.ReceitaPlanos),
		ReceitaTotal:         moneyToNumeric(dre.ReceitaTotal),
		CustoComissoes:       moneyToNumeric(dre.CustoComissoes),
		CustoInsumos:         moneyToNumeric(dre.CustoInsumos),
		CustoVariavelTotal:   moneyToNumeric(dre.CustoVariavelTotal),
		DespesaFixa:          moneyToNumeric(dre.DespesaFixa),
		DespesaVariavel:      moneyToNumeric(dre.DespesaVariavel),
		DespesaTotal:         moneyToNumeric(dre.DespesaTotal),
		ResultadoBruto:       moneyToNumeric(dre.ResultadoBruto),
		ResultadoOperacional: moneyToNumeric(dre.ResultadoOperacional),
		MargemBruta:          percentageToNumeric(dre.MargemBruta),
		MargemOperacional:    percentageToNumeric(dre.MargemOperacional),
		LucroLiquido:         moneyToNumeric(dre.LucroLiquido),
		ProcessadoEm:         timestampToTimestamptz(dre.ProcessadoEm),
	}

	result, err := r.queries.CreateDREMensal(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar DRE Mensal: %w", err)
	}

	dre.ID = pgUUIDToString(result.ID)
	dre.TenantID = pgUUIDToString(result.TenantID)
	dre.CriadoEm = timestamptzToTime(result.CriadoEm)
	dre.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)

	return nil
}

// FindByID busca uma DRE Mensal pelo ID.
func (r *DREMensalRepository) FindByID(ctx context.Context, tenantID, id string) (*entity.DREMensal, error) {
	idPg := uuidStringToPgtype(id)
	tenantPg := uuidStringToPgtype(tenantID)

	result, err := r.queries.GetDREMensalByID(ctx, db.GetDREMensalByIDParams{
		ID:       idPg,
		TenantID: tenantPg,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar DRE Mensal por ID: %w", err)
	}

	return r.modelToDREMensal(&result)
}

// FindByMesAno busca uma DRE Mensal por mês/ano.
func (r *DREMensalRepository) FindByMesAno(ctx context.Context, tenantID string, mesAno valueobject.MesAno) (*entity.DREMensal, error) {
	tenantPg := uuidStringToPgtype(tenantID)

	result, err := r.queries.GetDREMensalByMesAno(ctx, db.GetDREMensalByMesAnoParams{
		TenantID: tenantPg,
		MesAno:   mesAno.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar DRE Mensal por mês/ano: %w", err)
	}

	return r.modelToDREMensal(&result)
}

// List lista DREs Mensais paginadas.
func (r *DREMensalRepository) List(ctx context.Context, tenantID string, filters port.DREListFilters) ([]*entity.DREMensal, error) {
	tenantPg := uuidStringToPgtype(tenantID)

	page := filters.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filters.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	results, err := r.queries.ListDREMensalByTenant(ctx, db.ListDREMensalByTenantParams{
		TenantID: tenantPg,
		Limit:    int32(pageSize),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar DREs Mensais: %w", err)
	}

	dres := make([]*entity.DREMensal, len(results))
	for i, res := range results {
		dre, err := r.modelToDREMensal(&res)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter DRE Mensal %d: %w", i, err)
		}
		dres[i] = dre
	}

	return dres, nil
}

// ListByPeriod lista DREs Mensais por período.
func (r *DREMensalRepository) ListByPeriod(ctx context.Context, tenantID string, inicio, fim valueobject.MesAno) ([]*entity.DREMensal, error) {
	tenantPg := uuidStringToPgtype(tenantID)

	results, err := r.queries.ListDREMensalByPeriod(ctx, db.ListDREMensalByPeriodParams{
		TenantID: tenantPg,
		MesAno:   inicio.String(),
		MesAno_2: fim.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar DREs Mensais por período: %w", err)
	}

	dres := make([]*entity.DREMensal, len(results))
	for i, res := range results {
		dre, err := r.modelToDREMensal(&res)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter DRE Mensal %d: %w", i, err)
		}
		dres[i] = dre
	}

	return dres, nil
}

// Update atualiza uma DRE Mensal existente.
func (r *DREMensalRepository) Update(ctx context.Context, dre *entity.DREMensal) error {
	idPg := uuidStringToPgtype(dre.ID)
	tenantPg := uuidStringToPgtype(dre.TenantID)

	params := db.UpdateDREMensalParams{
		ID:                   idPg,
		TenantID:             tenantPg,
		ReceitaServicos:      moneyToNumeric(dre.ReceitaServicos),
		ReceitaProdutos:      moneyToNumeric(dre.ReceitaProdutos),
		ReceitaPlanos:        moneyToNumeric(dre.ReceitaPlanos),
		ReceitaTotal:         moneyToNumeric(dre.ReceitaTotal),
		CustoComissoes:       moneyToNumeric(dre.CustoComissoes),
		CustoInsumos:         moneyToNumeric(dre.CustoInsumos),
		CustoVariavelTotal:   moneyToNumeric(dre.CustoVariavelTotal),
		DespesaFixa:          moneyToNumeric(dre.DespesaFixa),
		DespesaVariavel:      moneyToNumeric(dre.DespesaVariavel),
		DespesaTotal:         moneyToNumeric(dre.DespesaTotal),
		ResultadoBruto:       moneyToNumeric(dre.ResultadoBruto),
		ResultadoOperacional: moneyToNumeric(dre.ResultadoOperacional),
		MargemBruta:          percentageToNumeric(dre.MargemBruta),
		MargemOperacional:    percentageToNumeric(dre.MargemOperacional),
		LucroLiquido:         moneyToNumeric(dre.LucroLiquido),
		ProcessadoEm:         timestampToTimestamptz(dre.ProcessadoEm),
	}

	result, err := r.queries.UpdateDREMensal(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao atualizar DRE Mensal: %w", err)
	}

	dre.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)

	return nil
}

// Delete remove uma DRE Mensal.
func (r *DREMensalRepository) Delete(ctx context.Context, tenantID, id string) error {
	idPg := uuidStringToPgtype(id)
	tenantPg := uuidStringToPgtype(tenantID)

	err := r.queries.DeleteDREMensal(ctx, db.DeleteDREMensalParams{
		ID:       idPg,
		TenantID: tenantPg,
	})
	if err != nil {
		return fmt.Errorf("erro ao deletar DRE Mensal: %w", err)
	}

	return nil
}

// SumReceitasByPeriod soma as receitas por período.
func (r *DREMensalRepository) SumReceitasByPeriod(ctx context.Context, tenantID string, inicio, fim valueobject.MesAno) (valueobject.Money, error) {
	tenantPg := uuidStringToPgtype(tenantID)

	result, err := r.queries.SumReceitasByPeriod(ctx, db.SumReceitasByPeriodParams{
		TenantID: tenantPg,
		MesAno:   inicio.String(),
		MesAno_2: fim.String(),
	})
	if err != nil {
		return valueobject.Money{}, fmt.Errorf("erro ao somar receitas: %w", err)
	}

	dec, err := aggregateNumericToDecimal(result)
	if err != nil {
		return valueobject.Money{}, fmt.Errorf("erro ao converter soma de receitas: %w", err)
	}

	return valueobject.NewMoneyFromDecimal(dec), nil
}

// SumDespesasByPeriod soma as despesas por período.
func (r *DREMensalRepository) SumDespesasByPeriod(ctx context.Context, tenantID string, inicio, fim valueobject.MesAno) (valueobject.Money, error) {
	tenantPg := uuidStringToPgtype(tenantID)

	result, err := r.queries.SumDespesasByPeriod(ctx, db.SumDespesasByPeriodParams{
		TenantID: tenantPg,
		MesAno:   inicio.String(),
		MesAno_2: fim.String(),
	})
	if err != nil {
		return valueobject.Money{}, fmt.Errorf("erro ao somar despesas: %w", err)
	}

	dec, err := aggregateNumericToDecimal(result)
	if err != nil {
		return valueobject.Money{}, fmt.Errorf("erro ao converter soma de despesas: %w", err)
	}

	return valueobject.NewMoneyFromDecimal(dec), nil
}

// AvgMargemBrutaByPeriod calcula a margem bruta média por período.
func (r *DREMensalRepository) AvgMargemBrutaByPeriod(ctx context.Context, tenantID string, inicio, fim valueobject.MesAno) (valueobject.Percentage, error) {
	tenantPg := uuidStringToPgtype(tenantID)

	result, err := r.queries.AvgMargemBrutaByPeriod(ctx, db.AvgMargemBrutaByPeriodParams{
		TenantID: tenantPg,
		MesAno:   inicio.String(),
		MesAno_2: fim.String(),
	})
	if err != nil {
		return valueobject.Percentage{}, fmt.Errorf("erro ao calcular média de margem bruta: %w", err)
	}

	dec, err := aggregateNumericToDecimal(result)
	if err != nil {
		return valueobject.Percentage{}, fmt.Errorf("erro ao converter média de margem bruta: %w", err)
	}

	return valueobject.NewPercentage(dec)
}

// AvgMargemOperacionalByPeriod calcula a margem operacional média por período.
func (r *DREMensalRepository) AvgMargemOperacionalByPeriod(ctx context.Context, tenantID string, inicio, fim valueobject.MesAno) (valueobject.Percentage, error) {
	tenantPg := uuidStringToPgtype(tenantID)

	result, err := r.queries.AvgMargemOperacionalByPeriod(ctx, db.AvgMargemOperacionalByPeriodParams{
		TenantID: tenantPg,
		MesAno:   inicio.String(),
		MesAno_2: fim.String(),
	})
	if err != nil {
		return valueobject.Percentage{}, fmt.Errorf("erro ao calcular média de margem operacional: %w", err)
	}

	dec, err := aggregateNumericToDecimal(result)
	if err != nil {
		return valueobject.Percentage{}, fmt.Errorf("erro ao converter média de margem operacional: %w", err)
	}

	return valueobject.NewPercentage(dec)
}

// modelToDREMensal converte o modelo do sqlc para a entidade de domínio.
func (r *DREMensalRepository) modelToDREMensal(m *db.DreMensal) (*entity.DREMensal, error) {
	mesAno, err := valueobject.NewMesAno(m.MesAno)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter mes_ano: %w", err)
	}

	margemBruta, err := numericToPercentage(m.MargemBruta)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter margem_bruta: %w", err)
	}

	margemOp, err := numericToPercentage(m.MargemOperacional)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter margem_operacional: %w", err)
	}

	dre := &entity.DREMensal{
		ID:                   pgUUIDToString(m.ID),
		TenantID:             pgUUIDToString(m.TenantID),
		MesAno:               mesAno,
		ReceitaServicos:      numericToMoney(m.ReceitaServicos),
		ReceitaProdutos:      numericToMoney(m.ReceitaProdutos),
		ReceitaPlanos:        numericToMoney(m.ReceitaPlanos),
		ReceitaTotal:         numericToMoney(m.ReceitaTotal),
		CustoComissoes:       numericToMoney(m.CustoComissoes),
		CustoInsumos:         numericToMoney(m.CustoInsumos),
		CustoVariavelTotal:   numericToMoney(m.CustoVariavelTotal),
		DespesaFixa:          numericToMoney(m.DespesaFixa),
		DespesaVariavel:      numericToMoney(m.DespesaVariavel),
		DespesaTotal:         numericToMoney(m.DespesaTotal),
		ResultadoBruto:       numericToMoney(m.ResultadoBruto),
		ResultadoOperacional: numericToMoney(m.ResultadoOperacional),
		MargemBruta:          margemBruta,
		MargemOperacional:    margemOp,
		LucroLiquido:         numericToMoney(m.LucroLiquido),
		ProcessadoEm:         timestamptzToTime(m.ProcessadoEm),
		CriadoEm:             timestamptzToTime(m.CriadoEm),
		AtualizadoEm:         timestamptzToTime(m.AtualizadoEm),
	}

	return dre, nil
}

// aggregateNumericToDecimal converte resultados agregados (sqlc retorna interface{}).
func aggregateNumericToDecimal(result interface{}) (decimal.Decimal, error) {
	switch v := result.(type) {
	case decimal.Decimal:
		return v, nil
	case pgtype.Numeric:
		return numericToDecimal(v), nil
	case float64:
		return decimal.NewFromFloat(v), nil
	case int64:
		return decimal.NewFromInt(v), nil
	default:
		return decimal.Zero, fmt.Errorf("tipo de resultado não suportado: %T", result)
	}
}
