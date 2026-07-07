package dag

import (
	"testing"

	"github.com/terrapilot/terrapilot/internal/descriptor"
)

func TestBuild_Order(t *testing.T) {
	stacks := []*descriptor.Stack{
		{Name: "app", DependsOn: []descriptor.Dependency{{Name: "vpc"}}},
		{Name: "vpc", DependsOn: []descriptor.Dependency{{Name: "networking"}}},
		{Name: "networking"},
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
		{Name: "a", DependsOn: []descriptor.Dependency{{Name: "b"}}},
		{Name: "b", DependsOn: []descriptor.Dependency{{Name: "a"}}},
	}

	_, err := Build(stacks)
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
}

func TestBuild_MissingDep(t *testing.T) {
	stacks := []*descriptor.Stack{
		{Name: "app", DependsOn: []descriptor.Dependency{{Name: "ghost"}}},
	}

	_, err := Build(stacks)
	if err == nil {
		t.Fatal("expected missing dependency error, got nil")
	}
}

func TestBuild_NoDeps(t *testing.T) {
	stacks := []*descriptor.Stack{
		{Name: "a"},
		{Name: "b"},
	}

	sorted, err := Build(stacks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sorted) != 2 {
		t.Errorf("want 2 stacks, got %d", len(sorted))
	}
}
