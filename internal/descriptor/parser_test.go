package descriptor

import (
	"testing"
)

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
	if stack.DependsOn[0].Name != "vpc" {
		t.Errorf("depends_on[0].name: want %q, got %q", "vpc", stack.DependsOn[0].Name)
	}
	if stack.DependsOn[0].MockOutputs["vpc_id"] != "vpc-mock-12345" {
		t.Errorf("mock_outputs vpc_id: want %q, got %q", "vpc-mock-12345", stack.DependsOn[0].MockOutputs["vpc_id"])
	}
	if stack.Locals["aws_region"] != "us-east-1" {
		t.Errorf("locals aws_region: want %q, got %q", "us-east-1", stack.Locals["aws_region"])
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
