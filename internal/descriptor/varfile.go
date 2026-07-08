package descriptor

import (
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

// parseVarFile reads a .tfvars file and returns its string-valued attributes.
// Non-string values (lists, maps, numbers, bools) are silently skipped — they
// are valid Terraform input variables but not usable in terrapilot templates.
func parseVarFile(path string) (map[string]string, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, nil // missing file is fine, terraform will error at runtime
	}

	file, diags := hclsyntax.ParseConfig(src, path, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, diags
	}

	attrs, diags := file.Body.JustAttributes()
	if diags.HasErrors() {
		return nil, diags
	}

	result := make(map[string]string)
	for name, attr := range attrs {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			continue // skip values that require evaluation context
		}
		if val.Type() == cty.String {
			result[name] = val.AsString()
		}
		// non-string types (list, map, number, bool) are skipped
	}
	return result, nil
}

// ResolveTemplateContext builds the full compile-time key-value map for a stack.
// Precedence (lowest to highest): var_files in order → meta.
func ResolveTemplateContext(s *Stack) (map[string]string, error) {
	ctx := make(map[string]string)

	for _, rel := range s.VarFiles {
		abs := filepath.Join(s.Dir, rel)
		vals, err := parseVarFile(abs)
		if err != nil {
			return nil, err
		}
		for k, v := range vals {
			ctx[k] = v
		}
	}

	for k, v := range s.Meta {
		ctx[k] = v
	}

	return ctx, nil
}
