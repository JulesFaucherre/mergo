package mergo

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/jfaucherre/mergo/cmd/create"
	"gitlab.com/jfaucherre/mergo/tools"
)

var mergoCmd = cobra.Command{
	Use:   "mergo",
	Short: "A tool to create pull requests from the command line",
	Long:  `A tool to create pull requests from the command line`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		return
	},
}

func init() {
	mergoCmd.AddCommand(create.CreateCmd)
	mergoCmd.Flags().BoolVarP(&tools.Verbose, "verbose", "v", false, "Use verbose output")
}

func Execute() {
	if err := mergoCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
