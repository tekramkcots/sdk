package index

//Series is the listing series of a company stock
type Series string

const (
	EQ Series = "EQ"
)

func (s Series) String() string {
	return string(s)
}
