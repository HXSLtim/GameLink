package model

// Currency 定义系统支持的三字母货币编码。
type Currency string

// Currency values define supported ISO-4217 3-letter codes.
const (
	CurrencyCNY Currency = "CNY"
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
)

// SupportedCurrencies 返回受支持的货币列表。
func SupportedCurrencies() []Currency {
	return []Currency{
		CurrencyCNY,
		CurrencyUSD,
		CurrencyEUR,
	}
}

// IsValidCurrency 判断输入货币是否受支持。
func IsValidCurrency(value Currency) bool {
	for _, c := range SupportedCurrencies() {
		if c == value {
			return true
		}
	}
	return false
}

// GormDataType 指定默认存储类型。
func (Currency) GormDataType() string {
	return "char(3)"
}
