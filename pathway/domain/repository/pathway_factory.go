package repository

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/torsday/gemini_jobcoin_mixer/pathway/domain/model"
	"log"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

func GenerateDepositAddress(unifiedOutputAddresses []string) string {
	salt := "stopgap salt" // Unify this, and comment about rainbow tables
	normalizedSlice := NormalizeAddressSlice(unifiedOutputAddresses)
	unifiedStr := PackSliceOfAddresses(normalizedSlice)

	// Use a salted hash to construct our deposit address
	theHash := sha256.Sum256([]byte(unifiedStr + salt)) // TODO confirm this concatenation
	base64Version := base64.StdEncoding.EncodeToString(theHash[:])
	urlEncodedVer := url.QueryEscape(base64Version) // url encode the result so it plays nice with HTTP
	return strings.Replace(urlEncodedVer, "%", "", -1)
}

// Create pathway from raw sql data.
func CreatePathwayFromRawSqlInput(rawSingularTblData []string) (model.Pathway, error) {

	if len(rawSingularTblData) == 0 {
		return model.Pathway{}, errors.New("CreatePathwayFromRawSqlInput was passed an empty slice")
	}

	// convert stored string of debt into unsigned 64 bit int
	amountOfDebt, err := strconv.ParseFloat(rawSingularTblData[2], 64)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	return model.Pathway{DepositAddress: rawSingularTblData[0],
		OutputAddresses: UnPackSliceOfAddresses(rawSingularTblData[1]),
		AmountOfDebt:    amountOfDebt,
	}, nil
}

func UnPackSliceOfAddresses(unifiedStr string) []string {
	return strings.Split(unifiedStr, ",")
}

func NormalizeAddressSlice(addresses []string) []string {
	// we know that addresses are case sensitive, and spaces matter.
	// so we currently only sort the addresses.
	// other normalization processes would go here if determined apt.
	sort.Strings(addresses)
	return addresses
}

func PackSliceOfAddresses(addresses []string) string {
	NormalizeAddressSlice(addresses)
	return strings.Join(addresses, ",")
}
