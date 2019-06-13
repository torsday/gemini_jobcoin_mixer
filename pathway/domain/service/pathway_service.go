package service

import (
	"errors"
	"github.com/torsday/gemini_jobcoin_mixer/pathway/domain/model"
	"github.com/torsday/gemini_jobcoin_mixer/pathway/domain/repository"
)

type PathwayService struct {
	pathwayRepository *repository.PathwayRepository
}

// get a new pathway service.
func NewPathwayService(pathwayRepository *repository.PathwayRepository) *PathwayService {
	return &PathwayService{pathwayRepository: pathwayRepository}
}

// get a pathway by searching for its deposit address.
func (pWayService *PathwayService) GetPathwayByDepositAddress(depositAddress string) (model.Pathway, error) {
	res, err := pWayService.pathwayRepository.GetPathwayByDepositAddress(depositAddress)
	if err != nil {
		return model.Pathway{}, errors.New("GetPathwayByDepositAddress didn't find a pathway")
	}
	return res, nil
}

// get (and create if necessary) a pathway.
func (pWayService *PathwayService) GetOrCreatePathway(outputAddresses []string) model.Pathway {
	newPathway := pWayService.pathwayRepository.FindOrCreatePathway(outputAddresses)

	return newPathway
}

// update the amount of debt a pathway has.
func (pWayService *PathwayService) UpdatePathwayAmount(pWay model.Pathway, newAmount float64) {
	pWayService.pathwayRepository.UpdatePathwayAmount(pWay, newAmount)
}

func (pWayService *PathwayService) UpdateWhenPathwayLastChecked(pWay model.Pathway, whenLastChecked int) {
	pWayService.pathwayRepository.UpdateWhenPathwayLastChecked(pWay, whenLastChecked)
}

// get pathways with debt.
func (pWayService *PathwayService) GetPathwaysWithDebt() []model.Pathway {
	return pWayService.pathwayRepository.GetPathwaysWithDebt()
}

// get a slice of ALL pathways.
func (pWayService *PathwayService) GetAllPathways() []model.Pathway {
	return pWayService.pathwayRepository.GetAllPathways()
}
