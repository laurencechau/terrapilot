package syncer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/terrapilot/terrapilot/internal/descriptor"
)

func TestSync(t *testing.T) {
	stack, err := descriptor.Parse("testdata/networking/stack.tp.hcl")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if err := Sync([]*descriptor.Stack{stack}); err != nil {
		t.Fatalf("sync: %v", err)
	}

	t.Cleanup(func() {
		os.Remove(filepath.Join(stack.Dir, "backend.tf"))
		os.Remove(filepath.Join(stack.Dir, "providers.tf"))
	})

	// backend.tf
	backend, err := os.ReadFile(filepath.Join(stack.Dir, "backend.tf"))
	if err != nil {
		t.Fatalf("reading backend.tf: %v", err)
	}
	backendStr := string(backend)
	if !strings.Contains(backendStr, header) {
		t.Error("backend.tf: missing header")
	}
	if !strings.Contains(backendStr, `bucket = "my-tfstate"`) {
		t.Error("backend.tf: missing bucket")
	}
	if !strings.Contains(backendStr, `key    = "networking/terraform.tfstate"`) {
		t.Error("backend.tf: missing key")
	}
	if !strings.Contains(backendStr, `region = "ap-southeast-1"`) {
		t.Error("backend.tf: missing region")
	}

	// providers.tf
	providers, err := os.ReadFile(filepath.Join(stack.Dir, "providers.tf"))
	if err != nil {
		t.Fatalf("reading providers.tf: %v", err)
	}
	providersStr := string(providers)
	if !strings.Contains(providersStr, `version = "5.50.0"`) {
		t.Error("providers.tf: missing provider_version")
	}
	if !strings.Contains(providersStr, `region = "ap-southeast-1"`) {
		t.Error("providers.tf: missing region")
	}
}

func TestSync_MissingKey(t *testing.T) {
	stack := &descriptor.Stack{
		Name:    "test",
		Dir:     t.TempDir(),
		Imports: []string{"testdata/shared/backend.tf.tpl"},
		Meta:    map[string]string{}, // no bucket, key, region
	}

	err := Sync([]*descriptor.Stack{stack})
	if err == nil {
		t.Fatal("expected error for missing template key, got nil")
	}
}
