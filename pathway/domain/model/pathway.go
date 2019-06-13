package model

type Pathway struct {
	DepositAddress  string
	OutputAddresses []string
	AmountOfDebt    float64
	WhenLastChecked int // timestamp
}

// Max amount of positions to the right of the decimal: 17 -> 0.12345678901234568
