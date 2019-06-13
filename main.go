package main

import (
	"github.com/torsday/gemini_jobcoin_mixer/cmd"
	"net/http"
	"time"
)

func main() {

	// Setup persistence layer if not already done.
	pwayTbl := InitializePathwayTbl()
	pwayTbl.BuildDbIfNotExists()

	activePoolAddress := "test_pool_address"
	httpClient := &http.Client{Timeout: 10 * time.Second}
	cmd.Dependencies.MixerService = InitializeMixerService(activePoolAddress, httpClient)
	cmd.Dependencies.PathwayService = InitializePathwayService()
	cmd.Dependencies.PathwayTbl = InitializePathwayTbl()

	cmd.Execute()
}
