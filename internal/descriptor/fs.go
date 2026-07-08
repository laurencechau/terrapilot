package descriptor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func readFile(path string) ([]byte, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", path, err)
	}
	return src, nil
}

// ResolvePath resolves a path from a stack descriptor.
// Paths starting with "//" are relative to the project root.
// All other paths are relative to the stack directory.
func ResolvePath(path, stackDir, root string) string {
	if strings.HasPrefix(path, "//") {
		return filepath.Join(root, path[2:])
	}
	return filepath.Join(stackDir, path)
}
