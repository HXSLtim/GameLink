package model

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCurrencyHelpers(t *testing.T) {
	expected := []Currency{CurrencyCNY, CurrencyUSD, CurrencyEUR}
	require.ElementsMatch(t, expected, SupportedCurrencies())
	require.True(t, IsValidCurrency(CurrencyEUR))
	require.False(t, IsValidCurrency(Currency("JPY")))
	require.Equal(t, "char(3)", Currency("").GormDataType())
}

func TestOrderPointerHelpers(t *testing.T) {
	order := &Order{TotalPriceCents: 12345}
	require.Equal(t, int64(12345), order.GetPriceCents())

	require.Zero(t, order.GetPlayerID())
	order.SetPlayerID(99)
	require.Equal(t, uint64(99), order.GetPlayerID())

	require.Zero(t, order.GetGameID())
	order.SetGameID(77)
	require.Equal(t, uint64(77), order.GetGameID())

	require.False(t, order.IsGiftOrder())
	recipient := uint64(10)
	order.RecipientPlayerID = &recipient
	require.True(t, order.IsGiftOrder())
}

func TestOrderNumberGenerators(t *testing.T) {
	assertFormat := func(prefix string, generator func() string) {
		re := regexp.MustCompile("^" + prefix + `\d{20}$`)
		value := generator()
		require.True(t, re.MatchString(value), "value %s did not match expected format", value)
		require.NotEqual(t, value, generator(), "two successive values should not be identical")
	}

	assertFormat("ESC", GenerateEscortOrderNo)
	assertFormat("GIFT", GenerateGiftOrderNo)
}

func TestServiceItemHelpers(t *testing.T) {
	item := &ServiceItem{
		SubCategory:    SubCategoryGift,
		BasePriceCents: 1000,
		CommissionRate: 0.25,
	}
	require.Equal(t, "service_items", ServiceItem{}.TableName())
	require.True(t, item.IsGift())

	commission, income := item.CalculateCommission(2)
	require.Equal(t, int64(500), commission)
	require.Equal(t, int64(1500), income)
}

func TestFinancialTableNames(t *testing.T) {
	require.Equal(t, "financial_accounts", FinancialAccount{}.TableName())
	require.Equal(t, "financial_vouchers", FinancialVoucher{}.TableName())
	require.Equal(t, "financial_voucher_entries", FinancialVoucherEntry{}.TableName())
	require.Equal(t, "financial_transactions", FinancialTransaction{}.TableName())
	require.Equal(t, "reconciliations", Reconciliation{}.TableName())
	require.Equal(t, "reconciliation_details", ReconciliationDetail{}.TableName())
	require.Equal(t, "financial_reports", FinancialReport{}.TableName())
	require.Equal(t, "financial_account_settings", FinancialAccountSetting{}.TableName())
}

func TestCommissionTableNames(t *testing.T) {
	require.Equal(t, "commission_rules", CommissionRule{}.TableName())
	require.Equal(t, "commission_records", CommissionRecord{}.TableName())
	require.Equal(t, "monthly_settlements", MonthlySettlement{}.TableName())
}

func TestRoleHelpers(t *testing.T) {
	role := &RoleModel{Slug: string(RoleSlugSuperAdmin), IsSystem: true}
	require.Equal(t, "roles", RoleModel{}.TableName())
	require.True(t, role.IsSystemRole())
	require.True(t, role.IsSuperAdmin())

	role.Slug = string(RoleSlugAdmin)
	require.False(t, role.IsSuperAdmin())
	role.IsSystem = false
	require.False(t, role.IsSystemRole())
}
