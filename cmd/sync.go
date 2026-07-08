package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/terrapilot/terrapilot/internal/dag"
	"github.com/terrapilot/terrapilot/internal/descriptor"
	"github.com/terrapilot/terrapilot/internal/rootdir"
	"github.com/terrapilot/terrapilot/internal/syncer"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Render shared HCL templates into each stack directory",
	Long: `Reads each stack's import block, renders the referenced .tpl files using
the stack's compile-time context (var_files + meta), and writes the output
into the stack directory. Output files are overwritten on every sync.`,
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

		return syncer.Sync(sorted)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
