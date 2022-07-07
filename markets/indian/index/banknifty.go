package index

func BankNifty() Index {
	return Index{
		Symbol:      BankNiftySymbol,
		Name:        BankNiftyName,
		Description: BankNiftyDescription,
		Composition: map[string]Company{
			HDFCBankSymbol:      {Symbol: HDFCBankSymbol, Name: HDFCBankName, Industry: Banking, Series: EQ, Weightage: 32.78},
			ICICIBankSymbol:     {Symbol: ICICIBankSymbol, Name: ICICIBankName, Industry: Banking, Series: EQ, Weightage: 26.68},
			KotakBankSymbol:     {Symbol: KotakBankSymbol, Name: KotakBankName, Industry: Banking, Series: EQ, Weightage: 13.25},
			AxisBankSymbol:      {Symbol: AxisBankSymbol, Name: AxisBankName, Industry: Banking, Series: EQ, Weightage: 10.41},
			SBISymbol:           {Symbol: SBISymbol, Name: SBIName, Industry: Banking, Series: EQ, Weightage: 9.53},
			IndusIndBankSymbol:  {Symbol: IndusIndBankSymbol, Name: IndusIndBankName, Industry: Banking, Series: EQ, Weightage: 2.49},
			BandhanBankSymbol:   {Symbol: BandhanBankSymbol, Name: BandhanBankName, Industry: Banking, Series: EQ, Weightage: 1.3},
			BankOfBarodaSymbol:  {Symbol: BankOfBarodaSymbol, Name: BankOfBarodaName, Industry: Banking, Series: EQ, Weightage: 1.01},
			FederalBankSymbol:   {Symbol: FederalBankSymbol, Name: FederalBankName, Industry: Banking, Series: EQ, Weightage: 0.95},
			AuBankSymbol:        {Symbol: AuBankSymbol, Name: AuBankName, Industry: Banking, Series: EQ, Weightage: 0.64},
			PNBBankSymbol:       {Symbol: PNBBankSymbol, Name: PNBBankName, Industry: Banking, Series: EQ, Weightage: 0.50},
			IDFCFirstBankSymbol: {Symbol: IDFCFirstBankSymbol, Name: IDFCFirstBankName, Industry: Banking, Series: EQ, Weightage: 0.47},
		},
	}
}
