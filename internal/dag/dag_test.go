package dag

import (
	"testing"

	"github.com/terrapilot/terrapilot/internal/descriptor"
)

func TestBuild_Order(t *testing.T) {
	stacks := []*descriptor.Stack{
		{Name: "app", Dir: "/root/app", DependsOn: []descriptor.Dependency{{Path: "/root/vpc"}}},
		{Name: "vpc", Dir: "/root/vpc", DependsOn: []descriptor.Dependency{{Path: "/root/networking"}}},
		{Name: "networking", Dir: "/root/networking"},
	}

	sorted, err := Build(stacks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"networking", "vpc", "app"}
	for i, s := range sorted {
		if s.Name != want[i] {
			t.Errorf("position %d: want %q, got %q", i, want[i], s.Name)
		}
	}
}

func TestBuild_Cycle(t *testing.T) {
	stacks := []*descriptor.Stack{
		{Name: "a", Dir: "/root/a", DependsOn: []descriptor.Dependency{{Path: "/root/b"}}},
		{Name: "b", Dir: "/root/b", DependsOn: []descriptor.Dependency{{Path: "/root/a"}}},
	}

	_, err := Build(stacks)
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
}

func TestBuild_MissingDep(t *testing.T) {
	stacks := []*descriptor.Stack{
		{Name: "app", Dir: "/root/app", DependsOn: []descriptor.Dependency{{Path: "/root/ghost"}}},
	}

	_, err := Build(stacks)
	if err == nil {
		t.Fatal("expected missing dependency error, got nil")
	}
}

func TestBuild_NoDeps(t *testing.T) {
	stacks := []*descriptor.Stack{
		{Name: "a", Dir: "/root/a"},
		{Name: "b", Dir: "/root/b"},
	}

	sorted, err := Build(stacks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sorted) != 2 {
		t.Errorf("want 2 stacks, got %d", len(sorted))
	}
}

func TestBuild_DuplicateDir(t *testing.T) {
	stacks := []*descriptor.Stack{
		{Name: "eks", Dir: "/root/eks"},
		{Name: "eks-copy", Dir: "/root/eks"},
	}

	_, err := Build(stacks)
	if err == nil {
		t.Fatal("expected duplicate directory error, got nil")
	}
}
