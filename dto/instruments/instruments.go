package instruments

type InstrumentFetcher interface {
	Fetch() ([]Instrument, error)
}

type Instrument struct{}
