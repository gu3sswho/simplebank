package util

// Constants for all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	RUB = "RUB"
)

// IsSupportedCurrency returns true if the currency is supported and false or not
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, RUB:
		return true
	}
	return false
}
