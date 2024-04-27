# tfirmg

This tool generates import blocks, removed blocks, and moved blocks based on differences in the codebase and tfstate.
## Requirements

### Terraform
v1.7 or higher

### OpenTofu
Sorry i don't know.

Need removed block supported.

## Install

### Go install
```bash
go install github.com/gainings/tfirmg
````

### Download


Get binary from releases
https://github.com/gainings/tfirmg/releases

## How to use

```bash
## local file
tfirmg generate --src-dir ./current_state --dst-dir ./new_state --src-tfstate-path file://my-tfstate/tfstate

## s3 
tfirmg geenrate --src-dir ./current_state --dst-dir ./new_state --src-tfstate-path s3://my-example-bucket/tfstate
```

## Introduction

For example, You have Terraform file like this.

```terraform
#current_state/main.tf
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.47.0"
    }
  }
}

provider "aws" {
  region = "ap-northeast-1"
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "my-vpc"
  cidr = "10.0.0.0/16"

  azs             = ["ap-northeast-1a", "ap-northeast-1c", "ap-northeast-1d"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]

  enable_nat_gateway = false
  enable_vpn_gateway = false

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"]
}

resource "aws_instance" "my_perfect_app" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t3.micro"

  tags = {
    Terraform = "true"
    Name      = "HelloWorld"
  }
}
```
When you want to refactor this Terraform to move resources to a different state.

Please move the code yourself using the same resource name.

```terraform
#current_state/main.tf
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.47.0"
    }
  }
}

provider "aws" {
  region = "ap-northeast-1"
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "my-vpc"
  cidr = "10.0.0.0/16"

  azs             = ["ap-northeast-1a", "ap-northeast-1c", "ap-northeast-1d"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]

  enable_nat_gateway = false
  enable_vpn_gateway = false

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}
```

```terraform
#app_state/main.tf
# just move to other directory with code
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.47.0"
    }
  }
}

provider "aws" {
  region = "ap-northeast-1"
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"]
}

resource "aws_instance" "my_perfect_app" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t3.micro"

  tags = {
    Terraform = "true"
    Name      = "HelloWorld"
  }
}
```

Use the following commands to generate the import block and remove block

```bash
tfirmg generate --src-dir ./current_state --dst-dir ./app_state --src-tfstate-path file://$(PWD)/current_state/terraform.tfstate
```

This command create following files.

```terraform
# current_state/removed.tf
removed {
  from = aws_instance.my_perfect_app
  lifecycle {
    destroy = false
  }
}
```

```terraform
# app_state/import.tf
import {
  to = aws_instance.my_perfect_app
  id = "i-123456789abcdef"
}
```

Next, you manually refactor Terraform to make it ready to run plan and apply.

This is because references may be broken due to the resources being moved.

Finally, by running terraform apply in each state, the import and remove are completed.

Also module can move like this.

```terraform
#current_state/main.tf
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.47.0"
    }
  }
}

provider "aws" {
  region = "ap-northeast-1"
}
```

```terraform
#netowrk_state/main.tf
#do refactor yourself
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.47.0"
    }
  }
}

provider "aws" {
  region = "ap-northeast-1"
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "my-vpc"
  cidr = "10.0.0.0/16"

  azs             = ["ap-northeast-1a", "ap-northeast-1c", "ap-northeast-1d"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]

  enable_nat_gateway = false
  enable_vpn_gateway = false

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}
```

```bash
tfirmg generate --src-dir ./current_state --dst-dir ./network_state --src-tfstate-path file://$(PWD)/current_state/terraform.tfstate
```

its generate code like this.

```terraform
# current_state/removed.tf
removed {
  from = module.vpc.aws_default_network_acl.this
  lifecycle {
    destroy = false
  }
}

removed {
  from = module.vpc.aws_default_route_table.default
  lifecycle {
    destroy = false
  }
}

removed {
  from = module.vpc.aws_default_security_group.this
  lifecycle {
    destroy = false
  }
}
...
```

```terraform
# network_state/import.tf
import {
  to = module.vpc.aws_default_network_acl.this[0]
  id = "acl-123456789abcdef"
}

import {
  to = module.vpc.aws_default_security_group.this[0]
  id = "sg-123456789abcdef"
}

import {
  to = module.vpc.aws_internet_gateway.this[0]
  id = "igw-123456789abcdef"
}
...
```
