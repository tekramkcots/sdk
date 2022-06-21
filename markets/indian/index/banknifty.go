package index

//Index the info about an index
type Index struct {
	Symbol      string
	Name        string
	Description string
	Composition map[string]float64
}

func BankNifty() Index {
	return Index{
		Symbol:      BankNiftySymbol,
		Name:        BankNiftyName,
		Description: BankNiftyDescription,
		Composition: map[string]float64{
			HDFCBankSymbol:      32.78,
			ICICIBankSymbol:     26.68,
			KotakBankSymbol:     13.25,
			AxisBankSymbol:      10.41,
			SBISymbol:           9.53,
			IndusIndBankSymbol:  2.49,
			BandhanBankSymbol:   1.3,
			BankOfBarodaSymbol:  1.01,
			FederalBankSymbol:   0.95,
			AuBankSymbol:        0.64,
			PNBBankSymbol:       0.50,
			IDFCFirstBankSymbol: 0.47,
		},
	}
}
