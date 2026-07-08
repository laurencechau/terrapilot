package descriptor

import (
	"path/filepath"
	"testing"
)

func TestResolvePath(t *testing.T) {
	cases := []struct {
		path     string
		stackDir string
		root     string
		want     string
	}{
		{"../networking", "/root/stacks/dev/eks", "/root", "/root/stacks/dev/networking"},
		{"//shared/backend.tf.tpl", "/root/stacks/dev/eks", "/root", "/root/shared/backend.tf.tpl"},
		{"//stacks/dev/env.tfvars", "/root/stacks/dev/eks", "/root", "/root/stacks/dev/env.tfvars"},
	}
	for _, c := range cases {
		got := ResolvePath(c.path, c.stackDir, c.root)
		if got != c.want {
			t.Errorf("ResolvePath(%q): want %q, got %q", c.path, c.want, got)
		}
	}
}

func TestParse_Valid(t *testing.T) {
	stack, err := Parse("testdata/stack.tp.hcl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stack.Name != "eks" {
		t.Errorf("name: want %q, got %q", "eks", stack.Name)
	}
	if stack.Description != "EKS cluster for production" {
		t.Errorf("description: want %q, got %q", "EKS cluster for production", stack.Description)
	}
	if stack.Runner != "tofu" {
		t.Errorf("runner: want %q, got %q", "tofu", stack.Runner)
	}
	if !stack.Enabled {
		t.Errorf("enabled: want true, got false")
	}
	if len(stack.VarFiles) != 2 {
		t.Errorf("var_files: want 2, got %d", len(stack.VarFiles))
	}
	if len(stack.Tags) != 3 {
		t.Errorf("tags: want 3, got %d", len(stack.Tags))
	}
	if len(stack.DependsOn) != 2 {
		t.Errorf("depends_on: want 2, got %d", len(stack.DependsOn))
	}
	wantVPCPath, _ := filepath.Abs("vpc")
	if stack.DependsOn[0].Path != wantVPCPath {
		t.Errorf("depends_on[0].path: want %q, got %q", wantVPCPath, stack.DependsOn[0].Path)
	}
	if stack.DependsOn[0].MockOutputs["vpc_id"] != "vpc-mock-12345" {
		t.Errorf("mock_outputs vpc_id: want %q, got %q", "vpc-mock-12345", stack.DependsOn[0].MockOutputs["vpc_id"])
	}
	if stack.Meta["aws_region"] != "us-east-1" {
		t.Errorf("meta aws_region: want %q, got %q", "us-east-1", stack.Meta["aws_region"])
	}
	if len(stack.Imports) != 2 {
		t.Errorf("imports: want 2, got %d", len(stack.Imports))
	}
}

func TestParse_Minimal(t *testing.T) {
	stack, err := Parse("testdata/minimal_stack.tp.hcl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stack.Name != "networking" {
		t.Errorf("name: want %q, got %q", "networking", stack.Name)
	}
	if !stack.Enabled {
		t.Error("enabled should default to true")
	}
	if len(stack.DependsOn) != 0 {
		t.Errorf("depends_on: want 0, got %d", len(stack.DependsOn))
	}
}

func TestParse_InvalidHCL(t *testing.T) {
	_, err := Parse("testdata/invalid_hcl_stack.tp.hcl")
	if err == nil {
		t.Fatal("expected error for invalid HCL, got nil")
	}
}

func TestParse_InvalidRunner(t *testing.T) {
	_, err := Parse("testdata/invalid_runner_stack.tp.hcl")
	if err == nil {
		t.Fatal("expected error for invalid runner, got nil")
	}
}

func TestResolveTemplateContext(t *testing.T) {
	stack, err := Parse("testdata/varfile_stack/stack.tp.hcl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, err := ResolveTemplateContext(stack)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// region comes from var_files
	if ctx["region"] != "ap-southeast-1" {
		t.Errorf("region: want %q, got %q", "ap-southeast-1", ctx["region"])
	}
	// env is overridden by meta
	if ctx["env"] != "override" {
		t.Errorf("env: want %q, got %q", "override", ctx["env"])
	}
	// backend_key comes from meta
	if ctx["backend_key"] != "dev/eks/terraform.tfstate" {
		t.Errorf("backend_key: want %q, got %q", "dev/eks/terraform.tfstate", ctx["backend_key"])
	}
	// count (number) should be skipped
	if _, ok := ctx["count"]; ok {
		t.Error("count should be skipped — not a string value")
	}
}

func TestDiscover_LocalsInheritance(t *testing.T) {
	stacks, err := Discover("testdata/locals_inherit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stacks) != 1 {
		t.Fatalf("want 1 stack, got %d", len(stacks))
	}

	s := stacks[0]
	// project comes from root locals.tp.hcl
	if s.Meta["project"] != "myproject" {
		t.Errorf("meta.project: want %q, got %q", "myproject", s.Meta["project"])
	}
	// env overridden by dev/locals.tp.hcl
	if s.Meta["env"] != "dev" {
		t.Errorf("meta.env: want %q, got %q", "dev", s.Meta["env"])
	}
}
