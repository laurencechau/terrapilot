package descriptor

import (
	"fmt"
	"os"
)

func readFile(path string) ([]byte, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", path, err)
	}
	return src, nil
}
