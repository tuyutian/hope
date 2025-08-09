// Package utils pkg/utils/money.go
package utils

import (
	"strconv"

	"github.com/shopspring/decimal"
)

// ParseMoneyFloat 将字符串金额转换为 float64
func ParseMoneyFloat(amount string) float64 {
	val, _ := strconv.ParseFloat(amount, 64)
	return val
}

// ParseMoneyDecimal 将字符串金额转换为 decimal.Decimal
func ParseMoneyDecimal(amount string) decimal.Decimal {
	dec, _ := decimal.NewFromString(amount)
	return dec
}

// DecimalToFloat 将 decimal 转换为 float64
func DecimalToFloat(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}
