package util

const (
	EUR = "EUR"
	USD = "USD"
	MMK = "MMK"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, MMK:
		return true
	default:
		return false
	}
}
