package index

import "github.com/tekramkcots/sdk/utilities"

//Index the info about an index
type Index struct {
	Symbol      string
	Name        string
	Description string
	Composition map[string]Company
}

func (i Index) GetSymbol() string {
	return i.Symbol
}

func (i Index) GetIndexStockSymbols() []string {
	return utilities.GetKeys(i.Composition)
}

func (i Index) GetStockWeightage() map[string]float32 {
	result := make(map[string]float32)
	for k, v := range i.Composition {
		result[k] = v.Weightage
	}
	return result
}
