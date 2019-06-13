package base

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// GET http://jobcoin.gemini.com/carnation/api/addresses/{address}
const ApiUrlAddressInfo = "http://jobcoin.gemini.com/carnation/api/addresses/"

// GET:  Get the list of all Jobcoin transactions; we currently don't care about this.
// POST: Send Jobcoins from one address to another.
// See docs: https://jobcoin.gemini.com/carnation/api
const ApiUrlTransactions = "http://jobcoin.gemini.com/carnation/api/transactions"

type Ambassador struct {
	HttpClient *http.Client
}

// Get an instance of the Ambassador
func NewAmbassador(httpClient *http.Client) *Ambassador {
	return &Ambassador{HttpClient: httpClient}
}

func (ambassador *Ambassador) GetAddressInfo(addressName string) (string, error) {
	//return `{"balance":"97.4","transactions":[{"timestamp":"2019-05-23T20:45:12.214Z","fromAddress":"Alice","toAddress":"Chris","amount":"2"},{"timestamp":"2019-05-23T20:45:18.938Z","toAddress":"Chris","amount":"50"},{"timestamp":"2019-05-23T21:03:10.404Z","fromAddress":"Chris","toAddress":"Alice","amount":"1"},{"timestamp":"2019-05-23T21:18:43.115Z","fromAddress":"Chris","toAddress":"Alice","amount":"0.9"},{"timestamp":"2019-05-23T21:18:44.696Z","fromAddress":"Chris","toAddress":"Alice","amount":"0.9"},{"timestamp":"2019-05-23T21:18:57.902Z","fromAddress":"Chris","toAddress":"Alice","amount":"0.9"},{"timestamp":"2019-05-23T21:19:21.236Z","fromAddress":"Chris","toAddress":"Alice","amount":"0.9"},{"timestamp":"2019-05-23T21:31:47.065Z","toAddress":"Chris","amount":"50"}]}`

	return ambassador.getRequest(ApiUrlAddressInfo + addressName)
}

func (ambassador *Ambassador) TransferCoin(withdrawalAddress string, depositAddress string, amount float64) (int, error) {
	if amount <= 0.0 {

		return 0, errors.New("can't transfer zero")
	}
	fmt.Println("Moving " + fmt.Sprintf("%.17f", amount) + " coin from " + withdrawalAddress + " to " + depositAddress)

	response, err := http.PostForm(ApiUrlTransactions, url.Values{
		"fromAddress": {withdrawalAddress},
		"toAddress":   {depositAddress},
		"amount":      {fmt.Sprintf("%.17f", amount)}})

	if err != nil {
		log.Println(err) // TODO better handle this
		return 0, errors.Wrap(err, "POST failed")
	}

	defer response.Body.Close()

	if response.StatusCode == 422 {
		return 422, errors.Wrap(err, "Insufficient Funds")
	}

	if response.StatusCode != 200 {
		return response.StatusCode, errors.Wrap(err, "POST failed, no 200")
	}

	return response.StatusCode, nil
}

// get the response of a get request given a URL.
func (ambassador *Ambassador) getRequest(url string) (string, error) {

	resp, err := ambassador.HttpClient.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "GET request failed for: "+url)
	}

	body, err := ioutil.ReadAll(resp.Body) // TODO close the body: http://polyglot.ninja/golang-making-http-requests/
	if err != nil {
		return "", errors.Wrap(err, "reading GET request failed for: "+url)
	}
	defer resp.Body.Close()

	return string(body), nil
}
