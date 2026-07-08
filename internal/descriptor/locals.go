package descriptor

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

const LocalsFile = "locals.tp.hcl"

var localsFileSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "meta"},
	},
}

func parseLocalsFile(path string) (map[string]string, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, nil // missing file is fine
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	file, diags := hclsyntax.ParseConfig(src, absPath, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, diags
	}

	content, diags := file.Body.Content(localsFileSchema)
	if diags.HasErrors() {
		return nil, diags
	}

	blocks := content.Blocks.OfType("meta")
	if len(blocks) == 0 {
		return nil, nil
	}

	attrs, diags := blocks[0].Body.JustAttributes()
	if diags.HasErrors() {
		return nil, diags
	}

	result := make(map[string]string, len(attrs))
	for name, attr := range attrs {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, diags
		}
		result[name] = val.AsString()
	}
	return result, nil
}

// InheritMeta walks from root down to stackDir, merging meta blocks from
// locals.tp.hcl files at each level — child directories win on conflict.
func InheritMeta(root, stackDir string) (map[string]string, error) {
	rel, err := filepath.Rel(root, stackDir)
	if err != nil {
		return nil, err
	}

	merged := make(map[string]string)

	current := root
	segments := []string{"."}
	if rel != "." {
		segments = append(segments, strings.Split(rel, string(filepath.Separator))...)
	}

	for _, seg := range segments {
		if seg != "." {
			current = filepath.Join(current, seg)
		}
		meta, err := parseLocalsFile(filepath.Join(current, LocalsFile))
		if err != nil {
			return nil, err
		}
		for k, v := range meta {
			merged[k] = v
		}
	}

	return merged, nil
}
