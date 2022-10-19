package cmd

import (
	"fmt"

	"github.com/kakao/detek/cases"
	"github.com/kakao/detek/pkg/detek"
	"github.com/kakao/detek/pkg/renderer"
	"github.com/spf13/cobra"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Before running the test, verify current test can be executed",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetSet := cases.DefaultSet
		if len(args) != 0 {
			targetSet = args[0]
		}
		m := detek.NewManager(
			cases.CollectorSet[targetSet](map[string]string{}),
			cases.DetectorSet[targetSet](map[string]string{}),
		)
		fmt.Print(
			renderer.RenderTablePlan(m.ShowPlan()),
		)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
