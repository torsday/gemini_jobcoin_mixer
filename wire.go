//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/torsday/gemini_jobcoin_mixer/api_ambassador/base"
	service2 "github.com/torsday/gemini_jobcoin_mixer/api_ambassador/service"
	service3 "github.com/torsday/gemini_jobcoin_mixer/mixer/domain/service"
	"github.com/torsday/gemini_jobcoin_mixer/pathway/domain/repository"
	"github.com/torsday/gemini_jobcoin_mixer/pathway/domain/service"
	"net/http"
)

func InitializePathwayTbl() *repository.PathwayTbl {
	panic(wire.Build(
		repository.NewPathwayTbl,
	))
	return &repository.PathwayTbl{}
}

func InitializePathwayRepository() *repository.PathwayRepository {
	panic(wire.Build(
		repository.NewPathwayTbl,
		repository.NewPathwayRepository,
	))
	return &repository.PathwayRepository{}
}

func InitializePathwayService() *service.PathwayService {
	panic(wire.Build(
		repository.NewPathwayTbl,
		repository.NewPathwayRepository,
		service.NewPathwayService,
	))
	return &service.PathwayService{}
}

// example client code: &http.Client{Timeout: 10 * time.Second},
func InitializeAmbassador(httpClient *http.Client) *base.Ambassador {
	panic(wire.Build(
		base.NewAmbassador,
	))
	return &base.Ambassador{}
}

func InitializeApiAmbassadorService(httpClient *http.Client) *service2.ApiAmbassadorService {
	panic(wire.Build(
		base.NewAmbassador,
		service2.NewApiAmbassadorService,
	))
	return &service2.ApiAmbassadorService{}
}

func InitializeMixerService(activePoolAddress string, httpClient *http.Client) *service3.MixerService {
	panic(wire.Build(
		repository.NewPathwayTbl,
		repository.NewPathwayRepository,
		service.NewPathwayService,
		base.NewAmbassador,
		service2.NewApiAmbassadorService,
		service3.NewMixerService,
	))
	return &service3.MixerService{}
}
