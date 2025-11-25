package postgres

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

// === UUID CONVERSIONS ===

// uuidToPgUUID converte uuid.UUID para pgtype.UUID
func uuidToPgUUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}
}

// pgUUIDToUUID converte pgtype.UUID para uuid.UUID
func pgUUIDToUUID(u pgtype.UUID) uuid.UUID {
	return u.Bytes
}

// uuidPtrToPgUUID converte *uuid.UUID para pgtype.UUID
func uuidPtrToPgUUID(u *uuid.UUID) pgtype.UUID {
	if u == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{
		Bytes: *u,
		Valid: true,
	}
}

// pgUUIDToUUIDPtr converte pgtype.UUID para *uuid.UUID
func pgUUIDToUUIDPtr(u pgtype.UUID) *uuid.UUID {
	if !u.Valid {
		return nil
	}
	id := uuid.UUID(u.Bytes)
	return &id
}

// === STRING CONVERSIONS ===

// strToText converte string para pgtype.Text
func strToText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  s != "",
	}
}

// textToStr converte pgtype.Text para string
func textToStr(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

// strPtrToText converte *string para pgtype.Text
func strPtrToText(s *string) pgtype.Text {
	if s == nil || *s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

// textToStrPtr converte pgtype.Text para *string
func textToStrPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

// strPtrToPgText converte string para *string (para queries sqlc que esperam *string)
func strPtrToPgText(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// pgTextToStr converte *string do sqlc para string
func pgTextToStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// pgTextToStrPtr converte *string do sqlc para *string do dominio (apenas copia)
func pgTextToStrPtr(s *string) *string {
	return s
}

// boolToPgBool converte bool para *bool (para queries sqlc)
func boolToPgBool(b bool) *bool {
	return &b
}

// === DECIMAL CONVERSIONS ===

// decimalToNumeric converte decimal.Decimal para pgtype.Numeric
func decimalToNumeric(d decimal.Decimal) pgtype.Numeric {
	if d.IsZero() {
		return pgtype.Numeric{Valid: true}
	}

	var num pgtype.Numeric
	if err := num.Scan(d.String()); err != nil {
		return pgtype.Numeric{Valid: false}
	}

	return num
}

// numericToDecimal converte pgtype.Numeric para decimal.Decimal
func numericToDecimal(n pgtype.Numeric) decimal.Decimal {
	if !n.Valid {
		return decimal.Zero
	}

	val, _ := n.Float64Value()
	return decimal.NewFromFloat(val.Float64)
}

// decimalPtrToNumeric converte *decimal.Decimal para pgtype.Numeric
func decimalPtrToNumeric(d *decimal.Decimal) pgtype.Numeric {
	if d == nil {
		return pgtype.Numeric{Valid: false}
	}
	return decimalToNumeric(*d)
}

// numericToDecimalPtr converte pgtype.Numeric para *decimal.Decimal
func numericToDecimalPtr(n pgtype.Numeric) *decimal.Decimal {
	if !n.Valid {
		return nil
	}
	d := numericToDecimal(n)
	return &d
}

// === TIME CONVERSIONS ===

// timeToPgTimestamp converte time.Time para pgtype.Timestamptz
func timeToPgTimestamp(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}

// timePtrToDate converte *time.Time para pgtype.Date
func timePtrToDate(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{
		Time:  *t,
		Valid: true,
	}
}

// dateToTimePtr converte pgtype.Date para *time.Time
func dateToTimePtr(d pgtype.Date) *time.Time {
	if !d.Valid {
		return nil
	}
	return &d.Time
}

// === BOOL CONVERSIONS ===

// boolPtrToBool converte *bool para bool
func boolPtrToBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// === VALUE OBJECT CONVERSIONS (para módulos antigos) ===

// uuidStringToPgtype converte string UUID para pgtype.UUID
func uuidStringToPgtype(s string) pgtype.UUID {
	var pguuid pgtype.UUID
	_ = pguuid.Scan(s)
	return pguuid
}

// dateToDate converte time.Time para pgtype.Date (compatibilidade antiga)
func dateToDate(d time.Time) pgtype.Date {
	return pgtype.Date{
		Time:  d,
		Valid: !d.IsZero(),
	}
}

// moneyToDecimal converte valueobject.Money para pgtype.Numeric
func moneyToDecimal(m valueobject.Money) pgtype.Numeric {
	return decimalToNumeric(m.Value())
}

// moneyToNumeric é alias para moneyToDecimal
func moneyToNumeric(m valueobject.Money) pgtype.Numeric {
	return moneyToDecimal(m)
}

// percentageToNumeric converte valueobject.Percentage para pgtype.Numeric
func percentageToNumeric(p valueobject.Percentage) pgtype.Numeric {
	return decimalToNumeric(p.Value())
}

// dmaisToInt32 converte valueobject.DMais para int32
func dmaisToInt32(d valueobject.DMais) int32 {
	return int32(d.Dias())
}

// === HELPERS ADICIONAIS (compatibilidade) ===

// pgUUIDToString converte pgtype.UUID para string
func pgUUIDToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	return uuid.UUID(u.Bytes).String()
}

// timestamptzToTime converte pgtype.Timestamptz para time.Time
func timestamptzToTime(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// timestamptzToTimePtr converte pgtype.Timestamptz para *time.Time
func timestamptzToTimePtr(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// timestampToTimestamptz converte time.Time para pgtype.Timestamptz
func timestampToTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

// decimalToMoney converte decimal.Decimal para valueobject.Money
func decimalToMoney(d decimal.Decimal) valueobject.Money {
	return valueobject.NewMoneyFromDecimal(d)
}

// numericToMoney converte pgtype.Numeric para valueobject.Money
func numericToMoney(n pgtype.Numeric) valueobject.Money {
	return decimalToMoney(numericToDecimal(n))
}

// int32ToDMais converte int32 para valueobject.DMais
func int32ToDMais(i int32) valueobject.DMais {
	d, _ := valueobject.NewDMais(int(i))
	return d
}

// numericToPercentage converte pgtype.Numeric para valueobject.Percentage
func numericToPercentage(n pgtype.Numeric) (valueobject.Percentage, error) {
	d := numericToDecimal(n)
	return valueobject.NewPercentage(d)
}

// === DECIMAL RAW CONVERSIONS (para campos sqlc que usam decimal.Decimal) ===

// moneyToRawDecimal converte valueobject.Money para decimal.Decimal (para uso direto em sqlc)
func moneyToRawDecimal(m valueobject.Money) decimal.Decimal {
	return m.Value()
}

// rawDecimalToMoney converte decimal.Decimal para valueobject.Money (para uso direto de sqlc)
func rawDecimalToMoney(d decimal.Decimal) valueobject.Money {
	return valueobject.NewMoneyFromDecimal(d)
}

// percentageToRawDecimal converte valueobject.Percentage para decimal.Decimal (para uso direto em sqlc)
func percentageToRawDecimal(p valueobject.Percentage) decimal.Decimal {
	return p.Value()
}

// rawDecimalToPercentage converte decimal.Decimal para valueobject.Percentage (para uso direto de sqlc)
func rawDecimalToPercentage(d decimal.Decimal) (valueobject.Percentage, error) {
	return valueobject.NewPercentage(d)
}

// Aliases para compatibilidade
var (
	percentageToDecimal = percentageToRawDecimal
	decimalToPercentage = rawDecimalToPercentage
)

// dateToTime converte pgtype.Date para time.Time
func dateToTime(d pgtype.Date) time.Time {
	if !d.Valid {
		return time.Time{}
	}
	return d.Time
}

// timeToDate converte time.Time para pgtype.Date
func timeToDate(t time.Time) pgtype.Date {
	return pgtype.Date{
		Time:  t,
		Valid: !t.IsZero(),
	}
}
