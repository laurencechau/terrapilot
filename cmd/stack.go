package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/terrapilot/terrapilot/internal/descriptor"
)

const stackTemplate = `stack "{{ .Name }}" {
{{- if .Description }}
  description = {{ printf "%q" .Description }}
{{- end }}
{{- if .Runner }}
  runner      = {{ printf "%q" .Runner }}
{{- end }}
  enabled     = {{ .Enabled }}
{{- if .VarFiles }}
  var_files   = [{{ range $i, $v := .VarFiles }}{{ if $i }}, {{ end }}{{ printf "%q" $v }}{{ end }}]
{{- end }}
{{- if .Tags }}
  tags        = [{{ range $i, $v := .Tags }}{{ if $i }}, {{ end }}{{ printf "%q" $v }}{{ end }}]
{{- end }}
}
{{- range .DependsOn }}

depends_on {
  path = {{ printf "%q" .Path }}
}
{{- end }}
`

var createFlags struct {
	description string
	runner      string
	disabled    bool
	varFiles    []string
	tags        []string
	dependsOn   []string
	dir         string
	filename    string
}

var stackCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a stack descriptor file",
	Example: `  terrapilot stack create eks
  terrapilot stack create eks --description "EKS cluster" --runner tofu --tag production
  terrapilot stack create eks --dir stacks/dev/eks --var-file ../../env.tfvars`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		runner := createFlags.runner
		if runner != "" && runner != "terraform" && runner != "tofu" {
			return fmt.Errorf("runner must be \"terraform\" or \"tofu\", got %q", runner)
		}

		dir := createFlags.dir
		if dir == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			dir = filepath.Join(cwd, name)
		}

		filename := createFlags.filename
		if filename == "" {
			filename = "stack" + descriptor.DescriptorSuffix
		} else if !strings.HasSuffix(filename, descriptor.DescriptorSuffix) {
			filename += descriptor.DescriptorSuffix
		}

		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("creating directory %s: %w", dir, err)
		}

		dest := filepath.Join(dir, filename)
		if _, err := os.Stat(dest); err == nil {
			return fmt.Errorf("%s already exists", dest)
		}

		tmpl, err := template.New("stack").Parse(stackTemplate)
		if err != nil {
			return err
		}

		type depData struct{ Path string }
		var deps []depData
		for _, p := range createFlags.dependsOn {
			deps = append(deps, depData{Path: p})
		}

		data := struct {
			Name        string
			Description string
			Runner      string
			Enabled     bool
			VarFiles    []string
			Tags        []string
			DependsOn   []depData
		}{
			Name:        name,
			Description: createFlags.description,
			Runner:      runner,
			Enabled:     !createFlags.disabled,
			VarFiles:    createFlags.varFiles,
			Tags:        createFlags.tags,
			DependsOn:   deps,
		}

		f, err := os.Create(dest)
		if err != nil {
			return fmt.Errorf("creating %s: %w", dest, err)
		}
		defer f.Close()

		if err := tmpl.Execute(f, data); err != nil {
			return fmt.Errorf("writing %s: %w", dest, err)
		}

		fmt.Printf("created %s\n", dest)
		return nil
	},
}

var stackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Manage stack descriptors",
}

func init() {
	f := stackCreateCmd.Flags()
	f.StringVar(&createFlags.description, "description", "", "stack description")
	f.StringVar(&createFlags.runner, "runner", "", "runner to use: terraform or tofu")
	f.BoolVar(&createFlags.disabled, "disabled", false, "set enabled = false")
	f.StringArrayVar(&createFlags.varFiles, "var-file", nil, "var file to include (repeatable)")
	f.StringArrayVar(&createFlags.tags, "tag", nil, "tag to apply (repeatable)")
	f.StringArrayVar(&createFlags.dependsOn, "depends-on", nil, "relative path to upstream stack (repeatable)")
	f.StringVar(&createFlags.dir, "dir", "", "directory to create the file in (default: current directory)")
	f.StringVar(&createFlags.filename, "filename", "", "filename to use (default: stack.tp.hcl)")

	stackCmd.AddCommand(stackCreateCmd)
	rootCmd.AddCommand(stackCmd)
}
