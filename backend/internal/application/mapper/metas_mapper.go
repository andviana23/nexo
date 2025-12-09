package mapper

import (
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/shopspring/decimal"
)

// ToMetaMensalResponse converte entidade MetaMensal para DTO Response
func ToMetaMensalResponse(meta *entity.MetaMensal) dto.MetaMensalResponse {
	return dto.MetaMensalResponse{
		ID:              meta.ID,
		MesAno:          meta.MesAno.String(),
		MetaFaturamento: meta.MetaFaturamento.Raw(),
		Origem:          string(meta.Origem),
		Status:          meta.Status,
		Realizado:       meta.Realizado.Raw(),
		Percentual:      meta.Percentual.String(),
		CriadoEm:        meta.CriadoEm.Format(time.RFC3339),
		AtualizadoEm:    meta.AtualizadoEm.Format(time.RFC3339),
	}
}

// ToMetaBarbeiroResponse converte entidade MetaBarbeiro para DTO Response
func ToMetaBarbeiroResponse(meta *entity.MetaBarbeiro) dto.MetaBarbeiroResponse {
	return dto.MetaBarbeiroResponse{
		ID:                       meta.ID,
		BarbeiroID:               meta.BarbeiroID,
		MesAno:                   meta.MesAno.String(),
		MetaServicosGerais:       meta.MetaServicosGerais.Raw(),
		MetaServicosExtras:       meta.MetaServicosExtras.Raw(),
		MetaProdutos:             meta.MetaProdutos.Raw(),
		RealizadoServicosGerais:  meta.RealizadoServicosGerais.Raw(),
		RealizadoServicosExtras:  meta.RealizadoServicosExtras.Raw(),
		RealizadoProdutos:        meta.RealizadoProdutos.Raw(),
		PercentualServicosGerais: meta.PercentualServicosGerais.String(),
		PercentualServicosExtras: meta.PercentualServicosExtras.String(),
		PercentualProdutos:       meta.PercentualProdutos.String(),
		CriadoEm:                 meta.CriadoEm.Format(time.RFC3339),
		AtualizadoEm:             meta.AtualizadoEm.Format(time.RFC3339),
	}
}

// ToMetaTicketResponse converte entidade MetaTicketMedio para DTO Response
func ToMetaTicketResponse(meta *entity.MetaTicketMedio) dto.MetaTicketResponse {
	return dto.MetaTicketResponse{
		ID:                   meta.ID,
		MesAno:               meta.MesAno.String(),
		Tipo:                 string(meta.Tipo),
		BarbeiroID:           meta.BarbeiroID,
		MetaValor:            meta.MetaValor.Raw(),
		TicketMedioRealizado: meta.TicketMedioRealizado.Raw(),
		Percentual:           meta.Percentual.String(),
		CriadoEm:             meta.CriadoEm.Format(time.RFC3339),
		AtualizadoEm:         meta.AtualizadoEm.Format(time.RFC3339),
	}
}

// FromSetMetaMensalRequest converte DTO Request para parâmetros do use case
func FromSetMetaMensalRequest(req dto.SetMetaMensalRequest) (
	mesAno valueobject.MesAno,
	metaFaturamento valueobject.Money,
	origem valueobject.OrigemMeta,
	err error,
) {
	// Parse mesAno
	mesAno, err = valueobject.NewMesAno(req.MesAno)
	if err != nil {
		return valueobject.MesAno{}, valueobject.Money{}, valueobject.OrigemMeta(""), fmt.Errorf("mes_ano inválido: %w", err)
	}

	// Parse meta
	metaDecimal, err := decimal.NewFromString(req.MetaFaturamento)
	if err != nil {
		return valueobject.MesAno{}, valueobject.Money{}, valueobject.OrigemMeta(""), fmt.Errorf("meta_faturamento inválido: %w", err)
	}
	metaFaturamento = valueobject.NewMoneyFromDecimal(metaDecimal)

	// Parse origem
	origem = valueobject.OrigemMeta(req.Origem)
	if !origem.IsValid() {
		return valueobject.MesAno{}, valueobject.Money{}, valueobject.OrigemMeta(""), fmt.Errorf("origem inválida")
	}

	return mesAno, metaFaturamento, origem, nil
}

// FromSetMetaBarbeiroRequest converte DTO Request para parâmetros do use case
func FromSetMetaBarbeiroRequest(req dto.SetMetaBarbeiroRequest) (
	mesAno valueobject.MesAno,
	metaServicosGerais, metaServicosExtras, metaProdutos valueobject.Money,
	err error,
) {
	// Parse mesAno
	mesAno, err = valueobject.NewMesAno(req.MesAno)
	if err != nil {
		return valueobject.MesAno{}, valueobject.Money{}, valueobject.Money{}, valueobject.Money{}, fmt.Errorf("mes_ano inválido: %w", err)
	}

	// Parse metas
	metaGeraisDecimal, err := decimal.NewFromString(req.MetaServicosGerais)
	if err != nil {
		return valueobject.MesAno{}, valueobject.Money{}, valueobject.Money{}, valueobject.Money{}, fmt.Errorf("meta_servicos_gerais inválido: %w", err)
	}
	metaServicosGerais = valueobject.NewMoneyFromDecimal(metaGeraisDecimal)

	metaExtrasDecimal, err := decimal.NewFromString(req.MetaServicosExtras)
	if err != nil {
		return valueobject.MesAno{}, valueobject.Money{}, valueobject.Money{}, valueobject.Money{}, fmt.Errorf("meta_servicos_extras inválido: %w", err)
	}
	metaServicosExtras = valueobject.NewMoneyFromDecimal(metaExtrasDecimal)

	metaProdutosDecimal, err := decimal.NewFromString(req.MetaProdutos)
	if err != nil {
		return valueobject.MesAno{}, valueobject.Money{}, valueobject.Money{}, valueobject.Money{}, fmt.Errorf("meta_produtos inválido: %w", err)
	}
	metaProdutos = valueobject.NewMoneyFromDecimal(metaProdutosDecimal)

	return mesAno, metaServicosGerais, metaServicosExtras, metaProdutos, nil
}

// FromSetMetaTicketRequest converte DTO Request para parâmetros do use case
func FromSetMetaTicketRequest(req dto.SetMetaTicketRequest) (
	mesAno valueobject.MesAno,
	tipo valueobject.TipoMetaTicket,
	barbeiroID *string,
	metaValor valueobject.Money,
	err error,
) {
	// Parse mesAno
	mesAno, err = valueobject.NewMesAno(req.MesAno)
	if err != nil {
		return valueobject.MesAno{}, valueobject.TipoMetaTicket(""), nil, valueobject.Money{}, fmt.Errorf("mes_ano inválido: %w", err)
	}

	// Parse tipo
	tipo = valueobject.TipoMetaTicket(req.Tipo)
	if !tipo.IsValid() {
		return valueobject.MesAno{}, valueobject.TipoMetaTicket(""), nil, valueobject.Money{}, fmt.Errorf("tipo inválido")
	}

	// BarbeiroID opcional
	if req.BarbeiroID != nil {
		barbeiroID = req.BarbeiroID
	}

	// Parse meta valor
	metaDecimal, err := decimal.NewFromString(req.MetaValor)
	if err != nil {
		return valueobject.MesAno{}, valueobject.TipoMetaTicket(""), nil, valueobject.Money{}, fmt.Errorf("meta_valor inválido: %w", err)
	}
	metaValor = valueobject.NewMoneyFromDecimal(metaDecimal)

	return mesAno, tipo, barbeiroID, metaValor, nil
}
