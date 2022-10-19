package cmd

import (
	"context"
	"fmt"

	"github.com/kakao/detek/cases"
	"github.com/kakao/detek/pkg/detek"
	"github.com/kakao/detek/pkg/renderer"
	"github.com/kakao/detek/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	kubeconfigPath string
	outputFormstS  string
	outputFormat   renderer.Format
	renderOpts     renderer.RenderOpts
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "try detekting known issues from the Kubernetes cluster",
	Long: fmt.Sprintf(`try detekting known issues from the Kubernetes cluster
// this will run "default" test set
detek run

// if you want to run the other test set (like "kakao"), than
detek run kakao
// this will run kakao test set

// currently available test sets are %v

// detek will use client configuration defined in
//   1. kubeconfig file located by "--kubeconfig" flag
//   2. kubeconfig file located by "KUBECONFIG" env
//   3. in-cluster client configuration (useful when using detek in a kubernetes cluster)
//   4. kubeconfig file located in default directory ($HOME/.kube/config)`,
		utils.Keys(cases.DetectorSet)),
	RunE: func(cmd *cobra.Command, args []string) error {
		{
			// pre-validation
			if err := outputFormat.IsValid(); err != nil {
				return err
			}
		}
		targetSet := cases.DefaultSet
		if len(args) != 0 {
			targetSet = args[0]
		}
		m := detek.NewManager(
			cases.CollectorSet[targetSet](map[string]string{
				cases.CONFIG_KUBECONFIG: kubeconfigPath,
			}),
			cases.DetectorSet[targetSet](map[string]string{}),
		)
		list, err := m.Run(context.Background(), &detek.MangerRunOptions{})
		if err != nil {
			return err
		}
		fmt.Println(
			renderer.RenderReports(list, outputFormat, renderOpts),
		)
		return nil
	},
	SilenceUsage: true,
}

func init() {
	flags := runCmd.Flags()
	flags.StringVar(&kubeconfigPath, "kubeconfig", "", "set kubeconfig path")
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringVarP(&outputFormstS, "format", "f", "html", "set output format. [json|table|html] ")
	runCmd.PersistentFlags().IntVar(&renderOpts.Table.MaxWidth, "table-max-width", 0, "truncate overflowed contents in table")
	runCmd.PersistentFlags().BoolVar(&renderOpts.JSON.Pretty, "json-pretty", true, "prettify json output")
}
