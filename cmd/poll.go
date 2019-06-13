package cmd

import (
	"github.com/spf13/cobra"
	"time"
)

func init() {
	rootCmd.AddCommand(poll)
}

var poll = &cobra.Command{
	Use:   "poll",
	Short: "Start a polling worker.",
	Long:  `start a polling worker to run mixer.`,
	Run: func(cmd *cobra.Command, args []string) {

		Dependencies.MixerService.CyclePool()

		time.Sleep(time.Duration(10) * time.Second)

	},
}
