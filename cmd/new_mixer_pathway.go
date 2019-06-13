package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	rootCmd.AddCommand(newMixerPathway)
}

var newMixerPathway = &cobra.Command{
	Use:   "new",
	Short: "Create a new mixer pathway.",
	Long:  `given a list of addresses to flush the money to, this returns a deposit address, completing the pathway.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 1 {
			fmt.Println("Creating new pathway")

			pway := Dependencies.PathwayService.GetOrCreatePathway(args)

			fmt.Println("Deposit to this address: " + pway.DepositAddress)
			fmt.Println("To receive funds at these addresses: " + strings.Join(args, ", "))
			fmt.Println("NOTE: in order to keep the money trail as random as possible, the money won't be evenly distributed, nor deposited all at once")
		} else {
			fmt.Println("Create a new pathway by entering:")
			fmt.Println("gemini_jobcoin_mixer new <address_1> <address_2> ...")
			fmt.Println("For 1 or more addresses; only alphanumeric please") // TODO validate this
		}

	},
}
