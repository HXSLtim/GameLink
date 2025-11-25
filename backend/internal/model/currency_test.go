package model_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gamelink/internal/model"
)

func TestCurrencyConstants(t *testing.T) {
	// 测试货币常量
	assert.Equal(t, model.Currency("CNY"), model.CurrencyCNY)
	assert.Equal(t, model.Currency("USD"), model.CurrencyUSD)
	assert.Equal(t, model.Currency("EUR"), model.CurrencyEUR)
}

func TestSupportedCurrencies(t *testing.T) {
	currencies := model.SupportedCurrencies()
	
	assert.Len(t, currencies, 3)
	assert.Contains(t, currencies, model.CurrencyCNY)
	assert.Contains(t, currencies, model.CurrencyUSD)
	assert.Contains(t, currencies, model.CurrencyEUR)
}

func TestIsValidCurrency(t *testing.T) {
	// 测试有效的货币
	validCurrencies := []model.Currency{
		model.CurrencyCNY,
		model.CurrencyUSD,
		model.CurrencyEUR,
	}

	for _, currency := range validCurrencies {
		assert.True(t, model.IsValidCurrency(currency), "Currency %s should be valid", currency)
	}

	// 测试无效的货币
	invalidCurrencies := []model.Currency{
		"JPY", // 日元（不支持）
		"GBP", // 英镑（不支持）
		"AUD", // 澳元（不支持）
		"CAD", // 加元（不支持）
		"",    // 空字符串
		"INVALID", // 无效货币
	}

	for _, currency := range invalidCurrencies {
		assert.False(t, model.IsValidCurrency(currency), "Currency %s should be invalid", currency)
	}
}

func TestCurrencyJSONSerialization(t *testing.T) {
	// 测试JSON序列化
	currency := model.CurrencyCNY
	data, err := json.Marshal(currency)
	assert.NoError(t, err)
	assert.Equal(t, `"CNY"`, string(data))

	// 测试JSON反序列化
	var decoded model.Currency
	err = json.Unmarshal([]byte(`"USD"`), &decoded)
	assert.NoError(t, err)
	assert.Equal(t, model.CurrencyUSD, decoded)

	// 测试EUR
	err = json.Unmarshal([]byte(`"EUR"`), &decoded)
	assert.NoError(t, err)
	assert.Equal(t, model.CurrencyEUR, decoded)
}

func TestCurrencyGormDataType(t *testing.T) {
	currency := model.Currency("")
	dataType := currency.GormDataType()
	assert.Equal(t, "char(3)", dataType)
}

func TestCurrencyInStructs(t *testing.T) {
	// 测试在结构体中的使用
	type TestStruct struct {
		ID       uint64       `json:"id"`
		Amount int64        `json:"amount"`
		Currency model.Currency `json:"currency"`
	}

	testObj := TestStruct{
		ID:       1,
		Amount:   10000,
		Currency: model.CurrencyCNY,
	}

	assert.Equal(t, uint64(1), testObj.ID)
	assert.Equal(t, int64(10000), testObj.Amount)
	assert.Equal(t, model.CurrencyCNY, testObj.Currency)
	assert.True(t, model.IsValidCurrency(testObj.Currency))

	// 测试JSON序列化
	data, err := json.Marshal(testObj)
	assert.NoError(t, err)

	// 验证JSON结构
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	assert.Contains(t, result, "id")
	assert.Contains(t, result, "amount")
	assert.Contains(t, result, "currency")
	assert.Equal(t, float64(1), result["id"])
	assert.Equal(t, float64(10000), result["amount"])
	assert.Equal(t, "CNY", result["currency"])
}

func TestCurrencyComparison(t *testing.T) {
	cny := model.CurrencyCNY
	usd := model.CurrencyUSD
	eur := model.CurrencyEUR

	// 测试相等性
	assert.True(t, cny == model.CurrencyCNY)
	assert.True(t, usd == model.CurrencyUSD)
	assert.True(t, eur == model.CurrencyEUR)

	// 测试不等性
	assert.True(t, cny != usd)
	assert.True(t, usd != eur)
	assert.True(t, cny != eur)
}

func TestCurrencyStringConversion(t *testing.T) {
	// 测试字符串转换
	cny := model.CurrencyCNY
	assert.Equal(t, "CNY", string(cny))

	usd := model.CurrencyUSD
	assert.Equal(t, "USD", string(usd))

	eur := model.CurrencyEUR
	assert.Equal(t, "EUR", string(eur))
}

func TestCurrencyEdgeCases(t *testing.T) {
	// 测试空字符串
	emptyCurrency := model.Currency("")
	assert.False(t, model.IsValidCurrency(emptyCurrency))
	assert.Equal(t, "", string(emptyCurrency))

	// 测试小写货币代码（应该无效）
	lowercaseCNY := model.Currency("cny")
	assert.False(t, model.IsValidCurrency(lowercaseCNY))

	// 测试大小写混合
	mixedCase := model.Currency("Cny")
	assert.False(t, model.IsValidCurrency(mixedCase))

	// 测试过长字符串
	longCurrency := model.Currency("CNYEXTRA")
	assert.False(t, model.IsValidCurrency(longCurrency))

	// 测试过短字符串
	shortCurrency := model.Currency("CN")
	assert.False(t, model.IsValidCurrency(shortCurrency))
}

func TestCurrencyInFinancialContext(t *testing.T) {
	// 测试在金融场景中的使用
	type Transaction struct {
		ID          uint64       `json:"id"`
		Amount      int64        `json:"amount"`      // 金额（分）
		Currency    model.Currency `json:"currency"`    // 货币
		Description string       `json:"description"`
	}

	transactions := []Transaction{
		{ID: 1, Amount: 10000, Currency: model.CurrencyCNY, Description: "人民币交易"},
		{ID: 2, Amount: 15000, Currency: model.CurrencyUSD, Description: "美元交易"},
		{ID: 3, Amount: 12000, Currency: model.CurrencyEUR, Description: "欧元交易"},
	}

	for _, tx := range transactions {
		assert.True(t, model.IsValidCurrency(tx.Currency))
		assert.NotEmpty(t, string(tx.Currency))
		assert.Len(t, string(tx.Currency), 3)
	}
}

func TestCurrencyJSONWithDifferentFormats(t *testing.T) {
	// 测试在复杂结构中的JSON序列化
	type PriceInfo struct {
		Amount   int64        `json:"amount"`
		Currency model.Currency `json:"currency"`
		Symbol   string       `json:"symbol"`
	}

	type Product struct {
		ID          uint64     `json:"id"`
		Name        string     `json:"name"`
		Price       PriceInfo  `json:"price"`
		Currencies  []model.Currency `json:"currencies"`
	}

	product := Product{
		ID:   1,
		Name: "测试产品",
		Price: PriceInfo{
			Amount:   9999,
			Currency: model.CurrencyCNY,
			Symbol:   "¥",
		},
		Currencies: []model.Currency{
			model.CurrencyCNY,
			model.CurrencyUSD,
			model.CurrencyEUR,
		},
	}

	data, err := json.Marshal(product)
	assert.NoError(t, err)

	// 验证JSON结构
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	assert.Contains(t, result, "price")
	price := result["price"].(map[string]interface{})
	assert.Equal(t, "CNY", price["currency"])

	assert.Contains(t, result, "currencies")
	currencies := result["currencies"].([]interface{})
	assert.Len(t, currencies, 3)
	assert.Equal(t, "CNY", currencies[0])
	assert.Equal(t, "USD", currencies[1])
	assert.Equal(t, "EUR", currencies[2])
}

func TestCurrencyZeroValue(t *testing.T) {
	var currency model.Currency
	assert.Equal(t, model.Currency(""), currency)
	assert.False(t, model.IsValidCurrency(currency))
	assert.Equal(t, "", string(currency))
}

func TestCurrencyCaseSensitivity(t *testing.T) {
	// 测试大小写敏感性
	upperCNY := model.Currency("CNY")
	lowerCNY := model.Currency("cny")
	mixedCNY := model.Currency("Cny")

	// 只有大写是有效的
	assert.True(t, model.IsValidCurrency(upperCNY))
	assert.False(t, model.IsValidCurrency(lowerCNY))
	assert.False(t, model.IsValidCurrency(mixedCNY))

	// 字符串表示也不同
	assert.Equal(t, "CNY", string(upperCNY))
	assert.Equal(t, "cny", string(lowerCNY))
	assert.Equal(t, "Cny", string(mixedCNY))
}