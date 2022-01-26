package niftybank

import "github.com/tekramkcots/sdk/dto/index"

//niftyBankCompanies is updated based on https://www1.nseindia.com/content/indices/ind_nifty_bank.pdf
var niftyBankCompanies = map[string]index.Company{
	"HDFCBANK": {
		Symbol:    "HDFCBANK",
		Name:      "HDFC Bank",
		Industry:  "Banking",
		Series:    index.EQ,
		Weightage: 27.80,
	},
	"ICICIBANK": {
		Symbol:    "ICICIBANK",
		Name:      "ICICI Bank",
		Industry:  "Banking",
		Series:    index.EQ,
		Weightage: 22.62,
	},
	"KOTAKBANK": {
		Symbol:    "KOTAKBANK",
		Name:      "Kotak Mahindra Bank",
		Industry:  "Banking",
		Series:    index.EQ,
		Weightage: 11.61,
	},
	"AXISBANK": {
		Symbol:    "AXISBANK",
		Name:      "Axis Bank",
		Industry:  "Banking",
		Series:    index.EQ,
		Weightage: 11.52,
	},
	"SBIN": {
		Symbol:    "SBIN",
		Name:      "State Bank of India",
		Industry:  "Banking",
		Series:    index.EQ,
		Weightage: 11.45,
	},
	"INDUSINDBK": {
		Symbol:    "INDUSINDBK",
		Name:      "IndusInd Bank",
		Industry:  "Banking",
		Series:    index.EQ,
		Weightage: 5.91,
	},
	"AUBANK": {
		Symbol:    "AUBANK",
		Name:      "AU Small Finance Bank",
		Industry:  "Banking",
		Series:    index.EQ,
		Weightage: 2.33,
	},
	"BANDHANBNK": {
		Symbol:    "BANDHANBNK",
		Name:      "Bandhan Bank",
		Industry:  "Banking",
		Series:    index.EQ,
		Weightage: 1.75,
	},
	"FEDERALBNK": {
		Symbol:    "FEDERALBNK",
		Name:      "Federal Bank",
		Industry:  "Banking",
		Series:    index.EQ,
		Weightage: 1.7,
	},
	"IDFCFIRSTB": {
		Symbol:    "IDFCFIRSTB",
		Name:      "IDF First Bank",
		Industry:  "Banking",
		Series:    index.EQ,
		Weightage: 1.54,
	},
}

func NormalisedValue(symbol string, value float64) float64 {
	company, ok := niftyBankCompanies[symbol]
	if !ok {
		return 0
	}
	return value / float64(company.Weightage)
}
