stack "eks" {
  var_files = ["env.tfvars"]
}

meta {
  backend_key = "dev/eks/terraform.tfstate"
  env         = "override"
}
