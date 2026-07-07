package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/terrapilot/terrapilot/internal/dag"
	"github.com/terrapilot/terrapilot/internal/descriptor"
	"github.com/terrapilot/terrapilot/internal/rootdir"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all stacks in dependency order",
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
			fmt.Println("no stacks found")
			return nil
		}

		sorted, err := dag.Build(stacks)
		if err != nil {
			return fmt.Errorf("building dependency graph: %w", err)
		}

		for i, s := range sorted {
			enabled := ""
			if !s.Enabled {
				enabled = " (disabled)"
			}
			tags := ""
			if len(s.Tags) > 0 {
				tags = " [" + strings.Join(s.Tags, ", ") + "]"
			}
			desc := ""
			if s.Description != "" {
				desc = " — " + s.Description
			}
			fmt.Printf("%d. %s%s%s%s\n", i+1, s.Name, desc, tags, enabled)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
