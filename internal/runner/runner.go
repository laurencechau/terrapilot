package runner

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/terrapilot/terrapilot/internal/descriptor"
)

// Run executes a terraform/tofu command in each stack directory in order.
func Run(stacks []*descriptor.Stack, args []string, tags []string) error {
	for _, s := range stacks {
		if !s.Enabled {
			fmt.Printf("[terrapilot] skipping: %s (disabled)\n", s.Name)
			continue
		}
		if !matchesTags(s, tags) {
			continue
		}

		bin, err := resolveBinary(s.Runner)
		if err != nil {
			return fmt.Errorf("stack %q: %w", s.Name, err)
		}

		cmdArgs := append(args, varFileFlags(s.VarFiles)...)
		fmt.Printf("[terrapilot] running %s: %s %v\n", s.Name, bin, args)

		cmd := exec.Command(bin, cmdArgs...)
		cmd.Dir = s.Dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("stack %q failed: %w", s.Name, err)
		}
	}
	return nil
}

func resolveBinary(runner string) (string, error) {
	if runner != "" {
		return runner, nil
	}
	for _, bin := range []string{"tofu", "terraform"} {
		if path, err := exec.LookPath(bin); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("neither 'tofu' nor 'terraform' found in PATH")
}

func varFileFlags(files []string) []string {
	flags := make([]string, 0, len(files)*2)
	for _, f := range files {
		flags = append(flags, "-var-file="+f)
	}
	return flags
}

func matchesTags(s *descriptor.Stack, tags []string) bool {
	if len(tags) == 0 {
		return true
	}
	for _, want := range tags {
		for _, have := range s.Tags {
			if want == have {
				return true
			}
		}
	}
	return false
}
