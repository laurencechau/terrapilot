package descriptor

import (
	"os"
	"path/filepath"
	"strings"
)

// Discover walks root recursively and returns all parsed stacks found.
func Discover(root string) ([]*Stack, error) {
	var stacks []*Stack

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), DescriptorSuffix) {
			s, err := Parse(path)
			if err != nil {
				return err
			}
			stacks = append(stacks, s)
		}
		return nil
	})

	return stacks, err
}
