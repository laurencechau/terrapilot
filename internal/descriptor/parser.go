package descriptor

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

const DescriptorSuffix = ".tp.hcl"

// Parse reads and validates a .tp.hcl file, returning a Stack.
func Parse(path string) (*Stack, error) {
	src, err := readFile(path)
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolving path %s: %w", path, err)
	}

	file, diags := hclsyntax.ParseConfig(src, absPath, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("HCL parse error in %s: %w", absPath, diags)
	}

	return decode(file.Body, absPath)
}

func decode(body hcl.Body, path string) (*Stack, error) {
	content, diags := body.Content(stackSchema)
	if diags.HasErrors() {
		return nil, fmt.Errorf("schema error in %s: %w", path, diags)
	}

	stack := &Stack{
		Enabled: true,
		Dir:     filepath.Dir(path),
	}

	if err := decodeStackBlock(content, stack, path); err != nil {
		return nil, err
	}
	if err := decodeDependsOn(content, stack, path); err != nil {
		return nil, err
	}
	if err := decodeMeta(content, stack, path); err != nil {
		return nil, err
	}
	if err := decodeImport(content, stack, path); err != nil {
		return nil, err
	}

	return stack, nil
}

func decodeStackBlock(content *hcl.BodyContent, stack *Stack, path string) error {
	blocks := content.Blocks.OfType("stack")
	if len(blocks) == 0 {
		return fmt.Errorf("%s: missing required stack block", path)
	}
	if len(blocks) > 1 {
		return fmt.Errorf("%s: only one stack block is allowed", path)
	}

	block := blocks[0]
	stack.Name = block.Labels[0]

	inner, diags := block.Body.Content(stackBodySchema)
	if diags.HasErrors() {
		return fmt.Errorf("stack block error in %s: %w", path, diags)
	}

	if attr, ok := inner.Attributes["description"]; ok {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return fmt.Errorf("description error in %s: %w", path, diags)
		}
		stack.Description = val.AsString()
	}

	if attr, ok := inner.Attributes["runner"]; ok {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return fmt.Errorf("runner error in %s: %w", path, diags)
		}
		r := val.AsString()
		if !validRunners[r] {
			return fmt.Errorf("%s: runner must be \"terraform\" or \"tofu\", got %q", path, r)
		}
		stack.Runner = r
	}

	if attr, ok := inner.Attributes["enabled"]; ok {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return fmt.Errorf("enabled error in %s: %w", path, diags)
		}
		stack.Enabled = val.True()
	}

	if attr, ok := inner.Attributes["var_files"]; ok {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return fmt.Errorf("var_files error in %s: %w", path, diags)
		}
		for it := val.ElementIterator(); it.Next(); {
			_, v := it.Element()
			stack.VarFiles = append(stack.VarFiles, v.AsString())
		}
	}

	if attr, ok := inner.Attributes["tags"]; ok {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return fmt.Errorf("tags error in %s: %w", path, diags)
		}
		for it := val.ElementIterator(); it.Next(); {
			_, v := it.Element()
			stack.Tags = append(stack.Tags, v.AsString())
		}
	}

	return nil
}

func decodeDependsOn(content *hcl.BodyContent, stack *Stack, path string) error {
	blocks := content.Blocks.OfType("depends_on")
	if len(blocks) == 0 {
		return nil
	}

	for _, block := range blocks {
		inner, diags := block.Body.Content(dependsOnSchema)
		if diags.HasErrors() {
			return fmt.Errorf("depends_on error in %s: %w", path, diags)
		}

		pathAttr := inner.Attributes["path"]
		val, diags := pathAttr.Expr.Value(nil)
		if diags.HasErrors() {
			return fmt.Errorf("depends_on path error in %s: %w", path, diags)
		}

		dep := Dependency{
			Path: filepath.Clean(filepath.Join(filepath.Dir(path), val.AsString())),
		}

		if attr, ok := inner.Attributes["mock_outputs"]; ok {
			val, diags := attr.Expr.Value(nil)
			if diags.HasErrors() {
				return fmt.Errorf("depends_on mock_outputs error in %s: %w", path, diags)
			}
			dep.MockOutputs = make(map[string]string)
			for it := val.ElementIterator(); it.Next(); {
				k, v := it.Element()
				dep.MockOutputs[k.AsString()] = v.AsString()
			}
		}

		stack.DependsOn = append(stack.DependsOn, dep)
	}

	return nil
}

func decodeMeta(content *hcl.BodyContent, stack *Stack, path string) error {
	blocks := content.Blocks.OfType("meta")
	if len(blocks) == 0 {
		return nil
	}
	if len(blocks) > 1 {
		return fmt.Errorf("%s: only one meta block is allowed", path)
	}

	attrs, diags := blocks[0].Body.JustAttributes()
	if diags.HasErrors() {
		return fmt.Errorf("meta error in %s: %w", path, diags)
	}

	stack.Meta = make(map[string]string)
	for name, attr := range attrs {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return fmt.Errorf("meta.%s error in %s: %w", name, path, diags)
		}
		stack.Meta[name] = val.AsString()
	}

	return nil
}

func decodeImport(content *hcl.BodyContent, stack *Stack, path string) error {
	blocks := content.Blocks.OfType("import")
	if len(blocks) == 0 {
		return nil
	}
	if len(blocks) > 1 {
		return fmt.Errorf("%s: only one import block is allowed", path)
	}

	inner, diags := blocks[0].Body.Content(importSchema)
	if diags.HasErrors() {
		return fmt.Errorf("import error in %s: %w", path, diags)
	}

	if attr, ok := inner.Attributes["files"]; ok {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return fmt.Errorf("import.files error in %s: %w", path, diags)
		}
		for it := val.ElementIterator(); it.Next(); {
			_, v := it.Element()
			stack.Imports = append(stack.Imports, v.AsString())
		}
	}

	return nil
}
