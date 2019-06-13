package service

import (
	"fmt"
	apiAmService "github.com/torsday/gemini_jobcoin_mixer/api_ambassador/service"
	"github.com/torsday/gemini_jobcoin_mixer/mixer/domain/model"
	pathwayModel "github.com/torsday/gemini_jobcoin_mixer/pathway/domain/model"
	"github.com/torsday/gemini_jobcoin_mixer/pathway/domain/service"
	"log"
	"math/rand"
	"time"
)

type MixerService struct {
	pathwayService       *service.PathwayService
	apiAmbassadorService *apiAmService.ApiAmbassadorService
	activePoolAddress    string
}

func NewMixerService(
	pathwayService *service.PathwayService,
	apiAmbassadorService *apiAmService.ApiAmbassadorService,
	activePoolAddress string,
) *MixerService {
	return &MixerService{
		pathwayService:       pathwayService,
		apiAmbassadorService: apiAmbassadorService,
		activePoolAddress:    activePoolAddress,
	}
}

// withdraws coin from pool, and into member deposit addresses (not necessarily all of it in one go).
func (mixerService *MixerService) CyclePool() {
	fmt.Println("\nCycling Pool")
	// check pathway deposit addresses for new deposits
	// move empty deposit addresses to pool && update pathway debt amount (with 1% cut)
	fmt.Println("\n\nChecking deposit addresses for new deposits")
	mixerService.checkAndEmptyDepositAddresses()

	fmt.Println("\n\nPruning Pool (moving funds to members)")
	// go through all pathways with outstanding debt, and "prune" each pathway
	mixerService.prunePathwayDebt()
}

func (mixerService *MixerService) checkAndEmptyDepositAddresses() {
	// get known deposit addresses
	var knownPathways []pathwayModel.Pathway
	knownPathways = mixerService.pathwayService.GetAllPathways()
	// for each one, check to see if it is not empty
	// if it's not empty, withdraw those funds into

	for i := 0; i < len(knownPathways); i++ {
		// get most up-to-date information on an address
		currentPathway := knownPathways[i]
		remoteAddressInfo, err := mixerService.apiAmbassadorService.GetAddressInfo(currentPathway.DepositAddress)
		if err != nil {
			log.Fatalln(err)
		}

		if remoteAddressInfo.Balance > 0 {
			// 1. move those funds (with a cut going to us) to the pool
			ourCut := remoteAddressInfo.Balance * 0.01 // as designed, our cut is stored in the pool itself
			amountForDebt := remoteAddressInfo.Balance - ourCut

			_, err := mixerService.apiAmbassadorService.TransferCoin(remoteAddressInfo.Name, mixerService.activePoolAddress, remoteAddressInfo.Balance)
			if err != nil {
				log.Fatalln(err)
			}

			// 2. Update the persisted pathway to reflect the increase in debt (note: this requires locking db)
			mixerService.pathwayService.UpdatePathwayAmount(currentPathway, currentPathway.AmountOfDebt+amountForDebt)
		}

		// if the amount with that address is empty, continue to the next
	}
}

// for all pathways with outstanding debt stored in the pool,
// go through and prune some of that money to the requisite output addresses
func (mixerService *MixerService) prunePathwayDebt() {
	var pendingTransfers []model.TransferRequest
	// get all pathways with outstanding debt
	var ripePathways []pathwayModel.Pathway
	ripePathways = mixerService.pathwayService.GetPathwaysWithDebt()

	// generate transfer requests for each pathway that has a debt
	for _, pathway := range ripePathways {
		pendingTransfers = append(
			pendingTransfers,
			createTransferRequestsFromPathway(pathway, mixerService.activePoolAddress)...)
	}
	mixerService.processTransfers(pendingTransfers)
}

// Process transfer requests.
func (mixerService *MixerService) processTransfers(transfers []model.TransferRequest) {

	// scramble order of transfer requests so outside observers have a greater challenge to reverse engineer pathways
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(
		len(transfers),
		func(i, j int) {
			transfers[i], transfers[j] = transfers[j], transfers[i]
		})

	// for each transfer request
	for _, tReq := range transfers {
		// 1. attempt transfer
		res, _ := mixerService.apiAmbassadorService.TransferCoin(
			tReq.FromAddress,
			tReq.ToAddress,
			tReq.Amount,
		)
		if res != 200 {
			fmt.Println("Transfer didn't work for: " + fmt.Sprintf("%.2f", tReq.Amount)+ " coins")
			fmt.Println(res)
		}
		// 2. if transfer succeeds, update pathway DB to reflect decrease in debt
		currentPathway, err := mixerService.pathwayService.GetPathwayByDepositAddress(tReq.OriginatingDepositAddress)
		if err != nil {
			fmt.Println(err)
			fmt.Println("failed search for this deposit address: " + tReq.FromAddress)
		}
		mixerService.pathwayService.UpdatePathwayAmount(currentPathway, currentPathway.AmountOfDebt-tReq.Amount)
	}
}

func createTransferRequestsFromPathway(pathway pathwayModel.Pathway, fundingAddress string) []model.TransferRequest {
	var transferRequests []model.TransferRequest

	// TODO Set 1.00 to a constant, or env
	if pathway.AmountOfDebt <= float64(len(pathway.OutputAddresses)) {
		// if pathway below a minimum, dump remaining funds into the first of output addresses
		transferRequests = append(transferRequests, model.TransferRequest{
			FromAddress:               fundingAddress,
			ToAddress:                 pathway.OutputAddresses[0],
			Amount:                    pathway.AmountOfDebt,
			OriginatingDepositAddress: pathway.DepositAddress,
		})
	} else {
		// Here's the logic to split the funds into different cycles and increases randomization of output amounts
		maxTithePerReceivingAddress := pathway.AmountOfDebt / float64(len(pathway.OutputAddresses))

		// randomize amount exported from pool to aid in randomizing data
		for _, oAdd := range pathway.OutputAddresses {
			deductionAmt := int(maxTithePerReceivingAddress / float64(1+rand.Intn(4))) // randomize and truncate to int
			if deductionAmt == 0 {
				deductionAmt = 1
			}
			transferRequests = append(transferRequests, model.TransferRequest{
				FromAddress:               fundingAddress,
				ToAddress:                 oAdd,
				Amount:                    float64(deductionAmt),
				OriginatingDepositAddress: pathway.DepositAddress,
			})
		}
	}

	return transferRequests
}
