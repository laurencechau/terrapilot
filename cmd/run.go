package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/terrapilot/terrapilot/internal/dag"
	"github.com/terrapilot/terrapilot/internal/descriptor"
	"github.com/terrapilot/terrapilot/internal/rootdir"
	"github.com/terrapilot/terrapilot/internal/runner"
)

var runTags []string

var runCmd = &cobra.Command{
	Use:   "run <command> [args...]",
	Short: "Run a terraform/tofu command across all stacks in dependency order",
	Example: `  terrapilot run plan
  terrapilot run apply
  terrapilot run plan --tag production`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		root, err := rootdir.Find(cwd)
		if err != nil {
			return err
		}

		stacks, err := descriptor.Discover(root)
		if err != nil {
			return fmt.Errorf("discovering stacks: %w", err)
		}
		if len(stacks) == 0 {
			fmt.Println("[terrapilot] no stacks found")
			return nil
		}

		sorted, err := dag.Build(stacks)
		if err != nil {
			return fmt.Errorf("building dependency graph: %w", err)
		}

		return runner.Run(sorted, args, runTags)
	},
}

func init() {
	runCmd.Flags().StringSliceVar(&runTags, "tag", nil, "only run stacks with this tag (repeatable)")
	rootCmd.AddCommand(runCmd)
}
