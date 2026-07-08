package descriptor

import (
	"os"
	"path/filepath"
	"strings"
)

// Discover walks root recursively and returns all parsed stacks found.
// Inherited locals from ancestor locals.tp.hcl files are merged into each
// stack's Locals map, with the stack's own locals block taking highest priority.
func Discover(root string) ([]*Stack, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	var stacks []*Stack

	err = filepath.WalkDir(absRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), DescriptorSuffix) {
			return nil
		}
		// skip standalone locals.tp.hcl files — they are not stacks
		if d.Name() == LocalsFile {
			return nil
		}

		s, err := Parse(path)
		if err != nil {
			return err
		}

		inherited, err := InheritMeta(absRoot, s.Dir)
		if err != nil {
			return err
		}

		// merge meta: inherited base, stack's own meta wins
		if len(inherited) > 0 {
			merged := make(map[string]string, len(inherited)+len(s.Meta))
			for k, v := range inherited {
				merged[k] = v
			}
			for k, v := range s.Meta {
				merged[k] = v
			}
			s.Meta = merged
		}

		stacks = append(stacks, s)
		return nil
	})

	return stacks, err
}
