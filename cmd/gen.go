package cmd

import (
	"github.com/learnergo/cuttle/invoke"
	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate certs",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if args != nil && len(args) != 0 {
			if args[0] == "all" {
				invoke.RunConfig()
			} else if args[0] == "some" {
				invoke.RunSpeConfig()
			} else {
				invoke.RunConfig()
			}
		} else {
			invoke.RunConfig()
		}
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}
