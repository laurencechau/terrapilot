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
	Locals      map[string]string
	Imports     []string
	// Dir is the directory containing the .tp.hcl file.
	Dir string
}

// Dependency represents a single upstream stack declared in the depends_on block.
type Dependency struct {
	Name        string
	MockOutputs map[string]string
}
