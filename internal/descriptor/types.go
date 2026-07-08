package descriptor

// Stack represents a parsed .tp.hcl stack descriptor.
type Stack struct {
	Name        string
	Description string
	Runner      string
	Enabled     bool
	VarFiles    []string
	Tags        []string
	DependsOn   []Dependency
	Meta        map[string]string // compile-time only, for template rendering
	Imports     []string
	// Dir is the absolute directory containing the .tp.hcl file.
	Dir string
	// Root is the absolute project root (directory containing .terrapilot.hcl or .git).
	Root string
}

// Dependency represents a single upstream stack declared in the depends_on block.
// Path is relative to the stack's own directory.
type Dependency struct {
	Path        string
	MockOutputs map[string]string
}
