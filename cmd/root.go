package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/torsday/gemini_jobcoin_mixer/mixer/domain/service"
	"github.com/torsday/gemini_jobcoin_mixer/pathway/domain/repository"
	service2 "github.com/torsday/gemini_jobcoin_mixer/pathway/domain/service"
	"os"
)

var Dependencies struct {
	MixerService   *service.MixerService
	PathwayService *service2.PathwayService
	PathwayTbl     *repository.PathwayTbl
}

var rootCmd = &cobra.Command{
	Use:   "jmixer",
	Short: "Jobcoin Mixer is a cryptocoin anonymizer.",
	Long:  `It pools together many different deposits before forwarding them onto user defined addresses.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
