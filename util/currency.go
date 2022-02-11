package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	PND = "PND"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, PND:
		return true
	}
	return false
}
