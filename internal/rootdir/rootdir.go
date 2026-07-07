package rootdir

import (
	"os"
	"path/filepath"
)

const ProjectFile = ".terrapilot.hcl"

// Find walks up from start looking for .terrapilot.hcl, then .git.
// Falls back to start if neither is found.
func Find(start string) (string, error) {
	abs, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	if root, ok := walkUp(abs, ProjectFile); ok {
		return root, nil
	}
	if root, ok := walkUp(abs, ".git"); ok {
		return root, nil
	}

	return abs, nil
}

func walkUp(dir, marker string) (string, bool) {
	for {
		if _, err := os.Stat(filepath.Join(dir, marker)); err == nil {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}
