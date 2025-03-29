package currency

type ConversionRequest struct {
	Amount float64
	From   string
	To     string
}

type ConversionResult struct {
	Amount          float64
	ConvertedAmount float64
	Rate            float64
	From            string
	To              string
}
