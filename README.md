---
Project Name: terrapilot
Website: https://terrapilot.sh
License: MIT
---

> **Work in progress** — core stack orchestration (`terrapilot run`, `terrapilot list`) is functional. Template syncing (`terrapilot sync`) is not yet implemented.

# terrapilot

**Run Terraform/OpenTofu stacks in dependency order, patch repeated HCL blocks, touch nothing else.**

terrapilot is an open source CLI tool that orchestrates Terraform and OpenTofu stacks. It reads a `.tp.hcl` descriptor in each stack folder, builds a dependency graph, and runs your commands in the correct order.

Each stack stays a self-contained, valid Terraform/OpenTofu stack — terrapilot is never a runtime dependency. Remove it at any time.

---

## Why terrapilot?

Terragrunt and Terramate solve multi-stack orchestration but come with a cost: new syntax, abstraction layers, and code generation that produces files you didn't write. Terragrunt's inheritance model in particular makes it hard to understand what's actually being applied.

terrapilot takes the opposite approach — no code generation, no abstraction layer, no lock-in. If you know Terraform, you already know how to use it.

> Inspired by Terramate and Terragrunt, but intentionally simpler.

1. **Pure HCL** — no new syntax, `.tp.hcl` uses standard HCL
2. **Terraform/OpenTofu native** — feels like a missing feature, not a competing tool
3. **Auto-detects runner** — uses `tofu` or `terraform` from your `$PATH`
4. **Zero runtime dependency** — each stack runs with plain `terraform apply`
5. **WET code is fine** — you own your `.tf` files, no generation
6. **Touch nothing else** — orchestration and template sync only

---

## Quick example

Add a `.terrapilot.hcl` at your project root — terrapilot walks up from the current directory to find it, so you can run commands from any subdirectory. Each stack folder contains a `stack.tp.hcl` file (any filename ending in `.tp.hcl` is valid).

```
cloud-resources/
  .terrapilot.hcl          ← marks the project root, run terrapilot from anywhere
  modules/
  stacks/
    dev/
      env.tfvars           ← environment-level vars (env, account ID, etc.)
      ap-southeast-1/
        region.tfvars      ← region-level vars (region, availability zones, etc.)
        eks/
          main.tf
          backend.tf
          variables.tf
          outputs.tf
          stack.tp.hcl     ← stack "eks"
        networking/
          main.tf
          backend.tf
          variables.tf
          outputs.tf
          stack.tp.hcl     ← stack "networking"
```

```hcl
# stacks/dev/ap-southeast-1/eks/stack.tp.hcl
stack "eks" {
  description = "EKS cluster for production"
  runner      = "tofu"                              # "terraform" or "tofu" (default: auto-detect)
  enabled     = true                                # set false to skip without deleting this file
  var_files   = ["../../env.tfvars", "eks.tfvars"]  # passed to terraform/tofu at runtime
  tags        = ["production", "ap-southeast-1"]
}

depends_on {
  stack "networking" {
    mock_outputs = {                                # mock values for planning without upstream state
      vpc_id = "vpc-mock-12345"
    }
  }
}

locals {
  key = "dev/ap-southeast-1/eks/terraform.tfstate" # compile-time values for template rendering
}

import {
  files = ["../../../../modules/backend.tf.tpl"]   # shared HCL templates to sync into this stack
}
```

```bash
terrapilot list                      # list all stacks in dependency order
terrapilot run plan                  # plan all stacks in dependency order
terrapilot run apply                 # apply all stacks in dependency order
terrapilot run plan --tag production # target stacks by tag
```

---

## Installation

```bash
go install github.com/terrapilot/terrapilot@latest
```

Or build from source:

```bash
git clone https://github.com/terrapilot/terrapilot
cd terrapilot
make install
```

