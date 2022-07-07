package index

func NormalisedValue(i Index, symbol string, value float64) (float64, bool) {
	weightage, ok := i.Composition[symbol]
	if !ok {
		return 0, false
	}
	return value * weightage / 100, true
}
