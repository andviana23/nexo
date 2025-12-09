package entity

import (
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

// DespesaFixa representa uma despesa recorrente que gera contas a pagar automaticamente
type DespesaFixa struct {
	ID        string
	TenantID  uuid.UUID
	UnidadeID string // Opcional: associar a uma unidade específica

	Descricao   string
	CategoriaID string // FK para categorias
	Fornecedor  string
	Valor       valueobject.Money

	DiaVencimento int  // Dia do mês (1-31)
	Ativo         bool // Se false, não gera conta no próximo mês

	Observacoes string

	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// NewDespesaFixa cria uma nova despesa fixa com validações
func NewDespesaFixa(
	tenantID uuid.UUID, descricao string,
	valor valueobject.Money,
	diaVencimento int,
) (*DespesaFixa, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}
	if descricao == "" {
		return nil, domain.ErrDescricaoRequired
	}
	if valor.IsNegative() || valor.IsZero() {
		return nil, domain.ErrValorNegativo
	}
	if diaVencimento < 1 || diaVencimento > 31 {
		return nil, domain.ErrDiaVencimentoInvalido
	}

	now := time.Now()
	return &DespesaFixa{
		ID:            uuid.NewString(),
		TenantID:      tenantID,
		Descricao:     descricao,
		Valor:         valor,
		DiaVencimento: diaVencimento,
		Ativo:         true, // Por padrão, despesas são criadas ativas
		CriadoEm:      now,
		AtualizadoEm:  now,
	}, nil
}

// Validate valida a entidade
func (d *DespesaFixa) Validate() error {
	if d.TenantID == uuid.Nil {
		return domain.ErrTenantIDRequired
	}
	if d.Descricao == "" {
		return domain.ErrDescricaoRequired
	}
	if d.Valor.IsNegative() || d.Valor.IsZero() {
		return domain.ErrValorNegativo
	}
	if d.DiaVencimento < 1 || d.DiaVencimento > 31 {
		return domain.ErrDiaVencimentoInvalido
	}
	return nil
}

// Desativar desativa a despesa fixa (não gerará mais contas)
func (d *DespesaFixa) Desativar() {
	d.Ativo = false
	d.AtualizadoEm = time.Now()
}

// Ativar ativa a despesa fixa
func (d *DespesaFixa) Ativar() {
	d.Ativo = true
	d.AtualizadoEm = time.Now()
}

// Toggle alterna o estado ativo/inativo
func (d *DespesaFixa) Toggle() {
	d.Ativo = !d.Ativo
	d.AtualizadoEm = time.Now()
}

// IsAtiva verifica se a despesa está ativa
func (d *DespesaFixa) IsAtiva() bool {
	return d.Ativo
}

// CalcularDataVencimento retorna a data de vencimento para um mês específico
// Ajusta automaticamente para meses com menos de 31 dias
func (d *DespesaFixa) CalcularDataVencimento(ano, mes int) time.Time {
	// Encontrar o último dia do mês
	ultimoDia := ultimoDiaDoMes(ano, mes)

	// Se o dia de vencimento é maior que o último dia do mês, usar o último dia
	dia := d.DiaVencimento
	if dia > ultimoDia {
		dia = ultimoDia
	}

	return time.Date(ano, time.Month(mes), dia, 0, 0, 0, 0, time.Local)
}

// ultimoDiaDoMes retorna o último dia de um mês específico
func ultimoDiaDoMes(ano, mes int) int {
	// Truque: dia 0 do próximo mês é o último dia do mês atual
	return time.Date(ano, time.Month(mes+1), 0, 0, 0, 0, 0, time.UTC).Day()
}

// ToContaPagar converte a despesa fixa em uma conta a pagar para um mês específico
func (d *DespesaFixa) ToContaPagar(ano, mes int) (*ContaPagar, error) {
	if !d.Ativo {
		return nil, domain.ErrDespesaInativa
	}

	dataVencimento := d.CalcularDataVencimento(ano, mes)

	now := time.Now()
	return &ContaPagar{
		ID:             uuid.NewString(),
		TenantID:       d.TenantID,
		Descricao:      d.Descricao,
		CategoriaID:    d.CategoriaID,
		Fornecedor:     d.Fornecedor,
		Valor:          d.Valor,
		Tipo:           valueobject.TipoCustoFixo,
		Recorrente:     true,
		Periodicidade:  "MENSAL",
		DataVencimento: dataVencimento,
		Status:         valueobject.StatusContaPendente,
		Observacoes:    d.gerarObservacaoContaGerada(ano, mes),
		CriadoEm:       now,
		AtualizadoEm:   now,
	}, nil
}

// gerarObservacaoContaGerada gera a observação padrão para rastreabilidade
func (d *DespesaFixa) gerarObservacaoContaGerada(ano, mes int) string {
	return fmt.Sprintf("Gerado automaticamente | Despesa Fixa: %s | Ref: %02d/%d", d.ID, mes, ano)
}

// SetCategoria define a categoria da despesa
func (d *DespesaFixa) SetCategoria(categoriaID string) {
	d.CategoriaID = categoriaID
	d.AtualizadoEm = time.Now()
}

// SetUnidade define a unidade da despesa
func (d *DespesaFixa) SetUnidade(unidadeID string) {
	d.UnidadeID = unidadeID
	d.AtualizadoEm = time.Now()
}

// SetFornecedor define o fornecedor da despesa
func (d *DespesaFixa) SetFornecedor(fornecedor string) {
	d.Fornecedor = fornecedor
	d.AtualizadoEm = time.Now()
}

// SetObservacoes define as observações
func (d *DespesaFixa) SetObservacoes(obs string) {
	d.Observacoes = obs
	d.AtualizadoEm = time.Now()
}

// Update atualiza os campos principais da despesa
func (d *DespesaFixa) Update(
	descricao string,
	valor valueobject.Money,
	diaVencimento int,
	categoriaID, fornecedor, unidadeID, observacoes string,
) error {
	if descricao == "" {
		return domain.ErrDescricaoRequired
	}
	if valor.IsNegative() || valor.IsZero() {
		return domain.ErrValorNegativo
	}
	if diaVencimento < 1 || diaVencimento > 31 {
		return domain.ErrDiaVencimentoInvalido
	}

	d.Descricao = descricao
	d.Valor = valor
	d.DiaVencimento = diaVencimento
	d.CategoriaID = categoriaID
	d.Fornecedor = fornecedor
	d.UnidadeID = unidadeID
	d.Observacoes = observacoes
	d.AtualizadoEm = time.Now()

	return nil
}
