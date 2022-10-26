package index

type Exchange string

const (
	NSE Exchange = "NSE"
)

func (e Exchange) String() string {
	return string(e)
}
