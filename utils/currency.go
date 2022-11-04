package utils

var validCurrency = map[string]bool{
	"USD": true,
	"CAD": true,
	"EUR": true,
}

func IsSupportedCurrency(currency string) bool {
	ok, _ := validCurrency[currency]
	return ok
}
