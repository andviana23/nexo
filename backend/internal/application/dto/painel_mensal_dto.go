package dto

// PainelMensalResponse representa a resposta do painel financeiro mensal
type PainelMensalResponse struct {
	Ano     int    `json:"ano"`
	Mes     int    `json:"mes"`
	NomeMes string `json:"nome_mes"`

	// Receitas
	ReceitaRealizada string `json:"receita_realizada"`
	ReceitaPendente  string `json:"receita_pendente"`
	ReceitaTotal     string `json:"receita_total"`

	// Despesas
	DespesasFixas     string `json:"despesas_fixas"`
	DespesasVariaveis string `json:"despesas_variaveis"`
	DespesasPagas     string `json:"despesas_pagas"`
	DespesasPendentes string `json:"despesas_pendentes"`
	DespesasTotal     string `json:"despesas_total"`

	// Resultados
	LucroBruto    string `json:"lucro_bruto"`
	LucroLiquido  string `json:"lucro_liquido"`
	MargemLiquida string `json:"margem_liquida"` // Em percentual

	// Metas
	MetaMensal     string `json:"meta_mensal"`
	PercentualMeta string `json:"percentual_meta"` // Em percentual
	DiferencaMeta  string `json:"diferenca_meta"`
	StatusMeta     string `json:"status_meta"` // "Atingida", "Em andamento", "Abaixo", "Sem meta"

	// Caixa
	SaldoCaixaAtual string `json:"saldo_caixa_atual"`

	// Comparativo
	VariacaoMesAnterior string `json:"variacao_mes_anterior"` // Em percentual
	TendenciaVariacao   string `json:"tendencia_variacao"`    // "up", "down", "stable"
}

// ProjecaoMensalResponse representa a projeção de um mês específico
type ProjecaoMensalResponse struct {
	Ano                int    `json:"ano"`
	Mes                int    `json:"mes"`
	NomeMes            string `json:"nome_mes"`
	ReceitaProjetada   string `json:"receita_projetada"`
	DespesasProjetadas string `json:"despesas_projetadas"`
	DespesasFixas      string `json:"despesas_fixas"`
	LucroProjetado     string `json:"lucro_projetado"`
	DiasUteis          int    `json:"dias_uteis"`
	MetaDiaria         string `json:"meta_diaria"`
	Confianca          string `json:"confianca"` // "Alta", "Média", "Baixa"
}

// ProjecoesResponse representa o resultado das projeções financeiras
type ProjecoesResponse struct {
	Projecoes           []ProjecaoMensalResponse `json:"projecoes"`
	MediaReceita3Meses  string                   `json:"media_receita_3_meses"`
	MediaDespesas3Meses string                   `json:"media_despesas_3_meses"`
	TendenciaReceita    string                   `json:"tendencia_receita"` // "Crescente", "Estável", "Decrescente"
	DataGeracao         string                   `json:"data_geracao"`
}
