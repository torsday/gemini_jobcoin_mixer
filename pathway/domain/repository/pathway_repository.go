package repository

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/torsday/gemini_jobcoin_mixer/pathway/domain/model"
	"log"
)

type PathwayRepository struct {
	pathwayTbl *PathwayTbl
}

func NewPathwayRepository(pathwayTbl *PathwayTbl) *PathwayRepository {
	return &PathwayRepository{pathwayTbl: pathwayTbl}
}

// Get a pathway object by searching the deposit address
func (pWayRepo *PathwayRepository) GetPathwayByDepositAddress(depositAddress string) (model.Pathway, error) {
	rawDbPathway := pWayRepo.pathwayTbl.GetPathwayByDepositAddress(depositAddress)

	res, err := CreatePathwayFromRawSqlInput(rawDbPathway)
	if err != nil {
		return model.Pathway{}, errors.New("GetPathwayByDepositAddress found no pathway")
	}
	return res, nil
}

// FindOrCreatePathway find or create a pathway.
func (pWayRepo *PathwayRepository) FindOrCreatePathway(outputAddresses []string) model.Pathway {
	var retPathway model.Pathway
	// create or find pathway using output_addresses as lookup

	// TODO DRY + code smell; method chain?
	unifiedOutAdd := PackSliceOfAddresses(NormalizeAddressSlice(outputAddresses))
	res := pWayRepo.pathwayTbl.FindPathwayByUnifiedOutputAddress(unifiedOutAdd)

	if len(res) == 4 {
		// pathway exists, create instance from old data
		var err error
		retPathway, err = CreatePathwayFromRawSqlInput(res)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Previous Pathway found")
	} else {
		fmt.Println("This is a new Pathway")
		// pathway unknown, create instance and persist
		// the reason we create a domain object we discard soon after is so that if we add any validation checking,
		// it can be localized to the domain object itself.
		depositAddress := GenerateDepositAddress(outputAddresses)
		outputAddresses := NormalizeAddressSlice(outputAddresses)

		retPathway = model.Pathway{
			DepositAddress:  depositAddress,
			OutputAddresses: outputAddresses,
			AmountOfDebt:    0,
		}

		// persist new pathway to db
		pWayRepo.pathwayTbl.CreatePathway(depositAddress, unifiedOutAdd)

		// we'll leave it to the use to create this address, as the API is fine with that.
		// in practice we would want to create it ourselves, assuming we'd then have the keys to the address.
	}
	return retPathway
}

// Update the amount of debt for a pathway
func (pWayRepo *PathwayRepository) UpdatePathwayAmount(pWay model.Pathway, newAmount float64) {
	stringOfFloatAmt := fmt.Sprintf("%.17f", newAmount)
	pWayRepo.pathwayTbl.UpdatePathwayAmount(pWay.DepositAddress, stringOfFloatAmt)
}

// Update when a pathway was last checked
func (pWayRepo *PathwayRepository) UpdateWhenPathwayLastChecked(pWay model.Pathway, whenLastChecked int) {
	pWayRepo.pathwayTbl.UpdateWhenPathwayLastChecked(pWay.DepositAddress, whenLastChecked)
}

// TODO this doesn't scale well, necessarily; keep track of when last checked to optimize.
// There are business decisions to be made (eternal monitoring of input addresses?)
func (pWayRepo *PathwayRepository) GetAllPathways() []model.Pathway {
	var retPathways []model.Pathway
	rawPways := pWayRepo.pathwayTbl.GetAllPathways()

	for _, rawRow := range rawPways {
		pway, err := CreatePathwayFromRawSqlInput(rawRow)
		if err != nil {
			log.Fatalln(err)
		}
		retPathways = append(retPathways, pway)
	}

	return retPathways
}

// get pathways with debt.
func (pWayRepo *PathwayRepository) GetPathwaysWithDebt() []model.Pathway {
	var retPaths []model.Pathway
	rawPathwaysWithDebt := pWayRepo.pathwayTbl.GetPathwaysWithDebt()

	for _, rawPath := range rawPathwaysWithDebt {
		pway, err := CreatePathwayFromRawSqlInput(rawPath)
		if err != nil {
			log.Fatalln(err)
		}
		retPaths = append(retPaths, pway)
	}

	return retPaths
}
