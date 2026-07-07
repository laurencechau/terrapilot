package descriptor

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// stackSchema defines the schema for the stack "<name>" block.
var stackSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "stack", LabelNames: []string{"name"}},
		{Type: "depends_on"},
		{Type: "locals"},
		{Type: "import"},
	},
}

// stackBodySchema defines attributes inside the stack block.
var stackBodySchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "description", Required: false},
		{Name: "runner", Required: false},
		{Name: "enabled", Required: false},
		{Name: "var_files", Required: false},
		{Name: "tags", Required: false},
	},
}

// dependsOnSchema defines the schema inside depends_on block.
var dependsOnSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "stack", LabelNames: []string{"name"}},
	},
}

// dependsOnStackSchema defines attributes inside a depends_on > stack block.
var dependsOnStackSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "mock_outputs", Required: false},
	},
}

// importSchema defines attributes inside the import block.
var importSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "files", Required: true},
	},
}

// validRunners is the set of allowed values for the runner attribute.
var validRunners = map[string]bool{
	"terraform": true,
	"tofu":      true,
}

// ctyStringList is a helper cty type for list(string).
var ctyStringList = cty.List(cty.String)
