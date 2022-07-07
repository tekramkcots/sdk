package index

//Index the info about an index
type Index struct {
	Symbol      string
	Name        string
	Description string
	Composition map[string]Company
}
