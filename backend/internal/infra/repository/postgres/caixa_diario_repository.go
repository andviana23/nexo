// Package postgres implementa os repositórios PostgreSQL do sistema NEXO
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

// CaixaDiarioRepository implementa port.CaixaDiarioRepository usando PostgreSQL/sqlc
type CaixaDiarioRepository struct {
	queries *db.Queries
}

// Compile-time check: garante que CaixaDiarioRepository implementa port.CaixaDiarioRepository
var _ port.CaixaDiarioRepository = (*CaixaDiarioRepository)(nil)

// NewCaixaDiarioRepository cria uma nova instância do repositório
func NewCaixaDiarioRepository(queries *db.Queries) *CaixaDiarioRepository {
	return &CaixaDiarioRepository{queries: queries}
}

// ============================================================
// CREATE
// ============================================================

// Create insere um novo caixa diário
func (r *CaixaDiarioRepository) Create(ctx context.Context, caixa *entity.CaixaDiario) error {
	result, err := r.queries.CreateCaixaDiario(ctx, db.CreateCaixaDiarioParams{
		ID:                uuidToPgUUID(caixa.ID),
		TenantID:          uuidToPgUUID(caixa.TenantID),
		UsuarioAberturaID: uuidToPgUUID(caixa.UsuarioAberturaID),
		DataAbertura:      timeToPgTimestamp(caixa.DataAbertura),
		SaldoInicial:      caixa.SaldoInicial,
		TotalEntradas:     caixa.TotalEntradas,
		TotalSaidas:       caixa.TotalSaidas,
		TotalSangrias:     caixa.TotalSangrias,
		TotalReforcos:     caixa.TotalReforcos,
		SaldoEsperado:     caixa.SaldoEsperado,
		Status:            string(caixa.Status),
	})
	if err != nil {
		return fmt.Errorf("erro ao criar caixa diário: %w", err)
	}

	// Atualizar campos de retorno
	caixa.CreatedAt = timestamptzToTime(result.CreatedAt)
	caixa.UpdatedAt = timestamptzToTime(result.UpdatedAt)

	return nil
}

// ============================================================
// READ
// ============================================================

// FindByID busca um caixa por ID
func (r *CaixaDiarioRepository) FindByID(ctx context.Context, caixaID, tenantID uuid.UUID) (*entity.CaixaDiario, error) {
	result, err := r.queries.GetCaixaDiarioByID(ctx, db.GetCaixaDiarioByIDParams{
		ID:       uuidToPgUUID(caixaID),
		TenantID: uuidToPgUUID(tenantID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCaixaNotFound
		}
		return nil, fmt.Errorf("erro ao buscar caixa por ID: %w", err)
	}

	return r.rowToCaixaDiario(&result), nil
}

// FindAberto busca o caixa aberto do tenant (deve existir apenas 1)
func (r *CaixaDiarioRepository) FindAberto(ctx context.Context, tenantID uuid.UUID) (*entity.CaixaDiario, error) {
	result, err := r.queries.GetCaixaDiarioAberto(ctx, uuidToPgUUID(tenantID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCaixaNaoAberto
		}
		return nil, fmt.Errorf("erro ao buscar caixa aberto: %w", err)
	}

	return r.rowAbertoToCaixaDiario(&result), nil
}

// ============================================================
// UPDATE
// ============================================================

// Update atualiza um caixa existente
func (r *CaixaDiarioRepository) Update(ctx context.Context, caixa *entity.CaixaDiario) error {
	result, err := r.queries.UpdateCaixaDiario(ctx, db.UpdateCaixaDiarioParams{
		ID:            uuidToPgUUID(caixa.ID),
		TenantID:      uuidToPgUUID(caixa.TenantID),
		TotalEntradas: caixa.TotalEntradas,
		TotalSaidas:   caixa.TotalSaidas,
		TotalSangrias: caixa.TotalSangrias,
		TotalReforcos: caixa.TotalReforcos,
		SaldoEsperado: caixa.SaldoEsperado,
	})
	if err != nil {
		return fmt.Errorf("erro ao atualizar caixa: %w", err)
	}

	caixa.UpdatedAt = timestamptzToTime(result.UpdatedAt)
	return nil
}

// UpdateTotais atualiza apenas os totais do caixa
func (r *CaixaDiarioRepository) UpdateTotais(ctx context.Context, caixaID, tenantID uuid.UUID, sangrias, reforcos, entradas decimal.Decimal) error {
	err := r.queries.UpdateCaixaDiarioTotais(ctx, db.UpdateCaixaDiarioTotaisParams{
		ID:            uuidToPgUUID(caixaID),
		TenantID:      uuidToPgUUID(tenantID),
		TotalSangrias: sangrias,
		TotalReforcos: reforcos,
		TotalEntradas: entradas,
	})
	if err != nil {
		return fmt.Errorf("erro ao atualizar totais do caixa: %w", err)
	}
	return nil
}

// Fechar fecha o caixa com os valores de fechamento
func (r *CaixaDiarioRepository) Fechar(ctx context.Context, caixa *entity.CaixaDiario) error {
	if caixa.DataFechamento == nil || caixa.SaldoReal == nil || caixa.Divergencia == nil {
		return fmt.Errorf("caixa não foi fechado corretamente via entidade")
	}

	result, err := r.queries.FecharCaixaDiario(ctx, db.FecharCaixaDiarioParams{
		ID:                       uuidToPgUUID(caixa.ID),
		TenantID:                 uuidToPgUUID(caixa.TenantID),
		UsuarioFechamentoID:      uuidPtrToPgUUID(caixa.UsuarioFechamentoID),
		DataFechamento:           timeToPgTimestamp(*caixa.DataFechamento),
		SaldoReal:                decimalToNumeric(*caixa.SaldoReal),
		Divergencia:              decimalToNumeric(*caixa.Divergencia),
		JustificativaDivergencia: caixa.JustificativaDivergencia,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrCaixaJaFechado
		}
		return fmt.Errorf("erro ao fechar caixa: %w", err)
	}

	caixa.UpdatedAt = timestamptzToTime(result.UpdatedAt)
	return nil
}

// ============================================================
// LIST
// ============================================================

// ListHistorico lista caixas fechados com paginação
func (r *CaixaDiarioRepository) ListHistorico(ctx context.Context, tenantID uuid.UUID, filters port.CaixaFilters) ([]*entity.CaixaDiario, error) {
	results, err := r.queries.ListCaixaDiarioHistorico(ctx, db.ListCaixaDiarioHistoricoParams{
		TenantID: uuidToPgUUID(tenantID),
		Column2:  timePtrToDate(filters.DataInicio),
		Column3:  timePtrToDate(filters.DataFim),
		Column4:  uuidPtrToPgUUID(filters.UsuarioID),
		Limit:    int32(filters.Limit),
		Offset:   int32(filters.Offset),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar histórico de caixas: %w", err)
	}

	caixas := make([]*entity.CaixaDiario, 0, len(results))
	for i := range results {
		caixas = append(caixas, r.rowHistoricoToCaixaDiario(&results[i]))
	}

	return caixas, nil
}

// CountHistorico conta o total de caixas fechados
func (r *CaixaDiarioRepository) CountHistorico(ctx context.Context, tenantID uuid.UUID, filters port.CaixaFilters) (int64, error) {
	count, err := r.queries.CountCaixaDiarioHistorico(ctx, db.CountCaixaDiarioHistoricoParams{
		TenantID: uuidToPgUUID(tenantID),
		Column2:  timePtrToDate(filters.DataInicio),
		Column3:  timePtrToDate(filters.DataFim),
		Column4:  uuidPtrToPgUUID(filters.UsuarioID),
	})
	if err != nil {
		return 0, fmt.Errorf("erro ao contar histórico de caixas: %w", err)
	}
	return count, nil
}

// ============================================================
// OPERAÇÕES
// ============================================================

// CreateOperacao registra uma operação no caixa
func (r *CaixaDiarioRepository) CreateOperacao(ctx context.Context, op *entity.OperacaoCaixa) error {
	result, err := r.queries.CreateOperacaoCaixa(ctx, db.CreateOperacaoCaixaParams{
		ID:        uuidToPgUUID(op.ID),
		CaixaID:   uuidToPgUUID(op.CaixaID),
		TenantID:  uuidToPgUUID(op.TenantID),
		Tipo:      string(op.Tipo),
		Valor:     op.Valor,
		Descricao: op.Descricao,
		Destino:   op.Destino,
		Origem:    op.Origem,
		UsuarioID: uuidToPgUUID(op.UsuarioID),
	})
	if err != nil {
		return fmt.Errorf("erro ao criar operação: %w", err)
	}

	op.CreatedAt = timestamptzToTime(result.CreatedAt)
	return nil
}

// ListOperacoes lista todas as operações de um caixa
func (r *CaixaDiarioRepository) ListOperacoes(ctx context.Context, caixaID, tenantID uuid.UUID) ([]entity.OperacaoCaixa, error) {
	results, err := r.queries.ListOperacoesByCaixa(ctx, db.ListOperacoesByCaixaParams{
		CaixaID:  uuidToPgUUID(caixaID),
		TenantID: uuidToPgUUID(tenantID),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar operações: %w", err)
	}

	operacoes := make([]entity.OperacaoCaixa, 0, len(results))
	for i := range results {
		operacoes = append(operacoes, r.rowToOperacao(&results[i]))
	}

	return operacoes, nil
}

// ListOperacoesByTipo lista operações filtradas por tipo
func (r *CaixaDiarioRepository) ListOperacoesByTipo(ctx context.Context, caixaID, tenantID uuid.UUID, tipo entity.TipoOperacaoCaixa) ([]entity.OperacaoCaixa, error) {
	results, err := r.queries.ListOperacoesByCaixaAndTipo(ctx, db.ListOperacoesByCaixaAndTipoParams{
		CaixaID:  uuidToPgUUID(caixaID),
		TenantID: uuidToPgUUID(tenantID),
		Tipo:     string(tipo),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar operações por tipo: %w", err)
	}

	operacoes := make([]entity.OperacaoCaixa, 0, len(results))
	for i := range results {
		operacoes = append(operacoes, r.rowTipoToOperacao(&results[i]))
	}

	return operacoes, nil
}

// SumOperacoesByTipo soma os valores de operações por tipo
func (r *CaixaDiarioRepository) SumOperacoesByTipo(ctx context.Context, caixaID, tenantID uuid.UUID) (map[entity.TipoOperacaoCaixa]decimal.Decimal, error) {
	results, err := r.queries.SumOperacoesByTipo(ctx, db.SumOperacoesByTipoParams{
		CaixaID:  uuidToPgUUID(caixaID),
		TenantID: uuidToPgUUID(tenantID),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao somar operações por tipo: %w", err)
	}

	sums := make(map[entity.TipoOperacaoCaixa]decimal.Decimal)
	for _, row := range results {
		var val decimal.Decimal
		switch v := row.Total.(type) {
		case float64:
			val = decimal.NewFromFloat(v)
		case int64:
			val = decimal.NewFromInt(v)
		case string:
			val, _ = decimal.NewFromString(v)
		case decimal.Decimal:
			val = v
		case pgtype.Numeric:
			// pgtype.Numeric é o tipo retornado pelo pgx/v5 para colunas numeric
			if v.Valid {
				// Converte para float64 primeiro e depois para decimal
				f, err := v.Float64Value()
				if err == nil && f.Valid {
					val = decimal.NewFromFloat(f.Float64)
				}
			}
		default:
			// Log do tipo desconhecido para debugging
			val = decimal.Zero
		}
		sums[entity.TipoOperacaoCaixa(row.Tipo)] = val
	}

	return sums, nil
}

// ============================================================
// MAPPERS: sqlc row → entity
// ============================================================

// rowToCaixaDiario converte GetCaixaDiarioByIDRow para entity
func (r *CaixaDiarioRepository) rowToCaixaDiario(row *db.GetCaixaDiarioByIDRow) *entity.CaixaDiario {
	caixa := &entity.CaixaDiario{
		ID:                       pgUUIDToUUID(row.ID),
		TenantID:                 pgUUIDToUUID(row.TenantID),
		UsuarioAberturaID:        pgUUIDToUUID(row.UsuarioAberturaID),
		UsuarioFechamentoID:      pgUUIDToUUIDPtr(row.UsuarioFechamentoID),
		DataAbertura:             timestamptzToTime(row.DataAbertura),
		DataFechamento:           timestamptzToTimePtr(row.DataFechamento),
		SaldoInicial:             row.SaldoInicial,
		TotalEntradas:            row.TotalEntradas,
		TotalSaidas:              row.TotalSaidas,
		TotalSangrias:            row.TotalSangrias,
		TotalReforcos:            row.TotalReforcos,
		SaldoEsperado:            row.SaldoEsperado,
		SaldoReal:                numericToDecimalPtr(row.SaldoReal),
		Divergencia:              numericToDecimalPtr(row.Divergencia),
		Status:                   entity.StatusCaixa(row.Status),
		JustificativaDivergencia: row.JustificativaDivergencia,
		CreatedAt:                timestamptzToTime(row.CreatedAt),
		UpdatedAt:                timestamptzToTime(row.UpdatedAt),
		UsuarioAberturaNome:      pgTextToStr(row.UsuarioAberturaNome),
		UsuarioFechamentoNome:    row.UsuarioFechamentoNome,
	}
	return caixa
}

// rowAbertoToCaixaDiario converte GetCaixaDiarioAbertoRow para entity
func (r *CaixaDiarioRepository) rowAbertoToCaixaDiario(row *db.GetCaixaDiarioAbertoRow) *entity.CaixaDiario {
	caixa := &entity.CaixaDiario{
		ID:                       pgUUIDToUUID(row.ID),
		TenantID:                 pgUUIDToUUID(row.TenantID),
		UsuarioAberturaID:        pgUUIDToUUID(row.UsuarioAberturaID),
		UsuarioFechamentoID:      pgUUIDToUUIDPtr(row.UsuarioFechamentoID),
		DataAbertura:             timestamptzToTime(row.DataAbertura),
		DataFechamento:           timestamptzToTimePtr(row.DataFechamento),
		SaldoInicial:             row.SaldoInicial,
		TotalEntradas:            row.TotalEntradas,
		TotalSaidas:              row.TotalSaidas,
		TotalSangrias:            row.TotalSangrias,
		TotalReforcos:            row.TotalReforcos,
		SaldoEsperado:            row.SaldoEsperado,
		SaldoReal:                numericToDecimalPtr(row.SaldoReal),
		Divergencia:              numericToDecimalPtr(row.Divergencia),
		Status:                   entity.StatusCaixa(row.Status),
		JustificativaDivergencia: row.JustificativaDivergencia,
		CreatedAt:                timestamptzToTime(row.CreatedAt),
		UpdatedAt:                timestamptzToTime(row.UpdatedAt),
		UsuarioAberturaNome:      pgTextToStr(row.UsuarioAberturaNome),
		UsuarioFechamentoNome:    row.UsuarioFechamentoNome,
	}
	return caixa
}

// rowHistoricoToCaixaDiario converte ListCaixaDiarioHistoricoRow para entity
func (r *CaixaDiarioRepository) rowHistoricoToCaixaDiario(row *db.ListCaixaDiarioHistoricoRow) *entity.CaixaDiario {
	caixa := &entity.CaixaDiario{
		ID:                       pgUUIDToUUID(row.ID),
		TenantID:                 pgUUIDToUUID(row.TenantID),
		UsuarioAberturaID:        pgUUIDToUUID(row.UsuarioAberturaID),
		UsuarioFechamentoID:      pgUUIDToUUIDPtr(row.UsuarioFechamentoID),
		DataAbertura:             timestamptzToTime(row.DataAbertura),
		DataFechamento:           timestamptzToTimePtr(row.DataFechamento),
		SaldoInicial:             row.SaldoInicial,
		TotalEntradas:            row.TotalEntradas,
		TotalSaidas:              row.TotalSaidas,
		TotalSangrias:            row.TotalSangrias,
		TotalReforcos:            row.TotalReforcos,
		SaldoEsperado:            row.SaldoEsperado,
		SaldoReal:                numericToDecimalPtr(row.SaldoReal),
		Divergencia:              numericToDecimalPtr(row.Divergencia),
		Status:                   entity.StatusCaixa(row.Status),
		JustificativaDivergencia: row.JustificativaDivergencia,
		CreatedAt:                timestamptzToTime(row.CreatedAt),
		UpdatedAt:                timestamptzToTime(row.UpdatedAt),
		UsuarioAberturaNome:      pgTextToStr(row.UsuarioAberturaNome),
		UsuarioFechamentoNome:    row.UsuarioFechamentoNome,
	}
	return caixa
}

// rowToOperacao converte ListOperacoesByCaixaRow para entity
func (r *CaixaDiarioRepository) rowToOperacao(row *db.ListOperacoesByCaixaRow) entity.OperacaoCaixa {
	return entity.OperacaoCaixa{
		ID:          pgUUIDToUUID(row.ID),
		CaixaID:     pgUUIDToUUID(row.CaixaID),
		TenantID:    pgUUIDToUUID(row.TenantID),
		Tipo:        entity.TipoOperacaoCaixa(row.Tipo),
		Valor:       row.Valor,
		Descricao:   row.Descricao,
		Destino:     row.Destino,
		Origem:      row.Origem,
		UsuarioID:   pgUUIDToUUID(row.UsuarioID),
		CreatedAt:   timestamptzToTime(row.CreatedAt),
		UsuarioNome: pgTextToStr(row.UsuarioNome),
	}
}

// rowTipoToOperacao converte ListOperacoesByCaixaAndTipoRow para entity
func (r *CaixaDiarioRepository) rowTipoToOperacao(row *db.ListOperacoesByCaixaAndTipoRow) entity.OperacaoCaixa {
	return entity.OperacaoCaixa{
		ID:          pgUUIDToUUID(row.ID),
		CaixaID:     pgUUIDToUUID(row.CaixaID),
		TenantID:    pgUUIDToUUID(row.TenantID),
		Tipo:        entity.TipoOperacaoCaixa(row.Tipo),
		Valor:       row.Valor,
		Descricao:   row.Descricao,
		Destino:     row.Destino,
		Origem:      row.Origem,
		UsuarioID:   pgUUIDToUUID(row.UsuarioID),
		CreatedAt:   timestamptzToTime(row.CreatedAt),
		UsuarioNome: pgTextToStr(row.UsuarioNome),
	}
}

// ============================================================
// HELPERS locais para conversões específicas
// ============================================================

// numericToDecimalPtr converte pgtype.Numeric para *decimal.Decimal
func numericToDecimalPtrLocal(n pgtype.Numeric) *decimal.Decimal {
	if !n.Valid {
		return nil
	}
	d := numericToDecimal(n)
	return &d
}
