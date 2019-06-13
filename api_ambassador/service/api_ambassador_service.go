package service

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/torsday/gemini_jobcoin_mixer/api_ambassador/base"
	"github.com/torsday/gemini_jobcoin_mixer/api_ambassador/model"
	"log"
	"strconv"
)

type ApiAmbassadorService struct {
	BaseAmbassador *base.Ambassador
}

type AddressInfoReturnJson struct {
	Balance      string `json:"balance"`
	Transactions []struct {
		Timestamp   string `json:"timestamp"`
		FromAddress string `json:"fromAddress"`
		ToAddress   string `json:"toAddress"`
		Amount      string `json:"amount"`
	} `json:"transactions"`
}

//type TransactionsJson struct {
//	Timestamp   string
//	FromAddress string
//	ToAddress   string
//	Amount      string
//}

func NewApiAmbassadorService(ambassador *base.Ambassador) *ApiAmbassadorService {
	return &ApiAmbassadorService{ambassador}
}

// transfer coin
func (aService *ApiAmbassadorService) TransferCoin(withdrawalAddress string, depositAddress string, amount float64) (int, error) {
	return aService.BaseAmbassador.TransferCoin(withdrawalAddress, depositAddress, amount)
}

// get address info
func (aService *ApiAmbassadorService) GetAddressInfo(address string) (model.AddressInfo, error) {
	rawAddressInfo, err := aService.BaseAmbassador.GetAddressInfo(address)
	if err != nil {
		return model.AddressInfo{}, errors.Wrap(err, "GetAddressInfo call failed for: "+address)
	}

	add, err := buildAddressInfoFromRawResponse(rawAddressInfo)
	if err != nil {
		return model.AddressInfo{Name: address}, errors.Wrap(err, "buildAddressInfoFromRawResponse failed within GetAddressInfo for: "+address)
	} else {
		return add, nil
	}
}

// build an address info domain object from raw data
func buildAddressInfoFromRawResponse(rawResponse string) (add model.AddressInfo, err error) {
	var addressInfoJsonContainer AddressInfoReturnJson

	marshErr := json.Unmarshal([]byte(rawResponse), &addressInfoJsonContainer)
	if marshErr != nil {
		log.Fatalln(err)
	}

	var amount float64
	amount, err = strconv.ParseFloat(addressInfoJsonContainer.Balance, 64)
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}

	if len(addressInfoJsonContainer.Transactions) == 0 {
		return model.AddressInfo{}, fmt.Errorf("That address has no history")
	} else {
		return model.AddressInfo{
			Name:    addressInfoJsonContainer.Transactions[0].ToAddress,
			Balance: amount,
		}, nil
	}

}
