package index

func NormalisedValue(i Index, symbol string, value float64) (float64, bool) {
	company, ok := i.Composition[symbol]
	if !ok {
		return 0, false
	}
	return value * float64(company.Weightage) / 100, true
}
