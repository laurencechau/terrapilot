stack "eks" {
  description = "EKS cluster for production"
  runner      = "tofu"
  enabled     = true
  var_files   = ["../common.tfvars", "eks.tfvars"]
  tags        = ["production", "eks", "ap-southeast-1"]
}

depends_on {
  path = "../vpc"
  mock_outputs = {
    vpc_id     = "vpc-mock-12345"
    subnet_ids = "subnet-mock-1"
  }
}

depends_on {
  path = "../networking"
}

meta {
  key        = "eks/terraform.tfstate"
  aws_region = "us-east-1"
}

import {
  files = [
    "../../shared/backend.tf.tpl",
    "../../shared/providers.tf.tpl"
  ]
}
