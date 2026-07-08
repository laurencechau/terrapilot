stack "networking" {
  var_files = ["env.tfvars"]
}

meta {
  bucket = "my-tfstate"
  key    = "networking/terraform.tfstate"
}

import {
  files = [
    "../shared/backend.tf.tpl",
    "../shared/providers.tf.tpl"
  ]
}
