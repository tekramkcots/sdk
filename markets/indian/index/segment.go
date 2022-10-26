package index

type Segment string

const (
	Indices Segment = "INDICES"
)

func (s Segment) String() string {
	return string(s)
}
