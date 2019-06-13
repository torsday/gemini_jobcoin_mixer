package model

type TransferRequest struct {
	FromAddress               string
	ToAddress                 string
	Amount                    float64
	OriginatingDepositAddress string
}
