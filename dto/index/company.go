package index

//Company is a company whose stock is listed in Nifty 50
type Company struct {
	Symbol    string
	Name      string
	Industry  string
	Series    Series
	Weightage float32
}
