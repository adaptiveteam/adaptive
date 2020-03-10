module "reporting_vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "${local.client_id}_reporting"
  cidr = var.vpc_cidr

  azs = local.availability_zones

  private_subnets = var.vpc_private_subnets

  public_subnets = var.vpc_public_subnets

  database_subnets = var.vpc_database_subnets

  enable_nat_gateway = true
  enable_vpn_gateway = false

  # To address: InvalidVPCNetworkStateFault: Cannot create a publicly accessible DBInstance.  The specified VPC
  # does not support DNS resolution, DNS hostnames, or both. Update the VPC and then try again
  enable_dns_hostnames = true
  enable_dns_support = true

  create_database_subnet_group = true
  create_database_subnet_route_table = true
  create_database_internet_gateway_route = true

  private_subnet_tags = {
    Layer = "private"
  }

  public_subnet_tags = {
    Layer = "public"
  }
}

resource "aws_security_group" "reporting" {
  name = "${local.client_id}_reporting"
  description = "Security group for AWS lambda and AWS RDS connection"
  vpc_id = module.reporting_vpc.vpc_id

  ingress {
    cidr_blocks = [
      "0.0.0.0/0",
    ]
    from_port = 0
    ipv6_cidr_blocks = [
      "::/0",
    ]
    prefix_list_ids = []
    protocol = "-1"
    security_groups = []
    self = false
    to_port = 0
  }
  ingress {
    cidr_blocks = [
      "0.0.0.0/0",
    ]
    description = ""
    from_port = 3306
    ipv6_cidr_blocks = [
      "::/0",
    ]
    prefix_list_ids = []
    protocol = "tcp"
    security_groups = []
    self = false
    to_port = 3306
  }
  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = [
      "0.0.0.0/0"]
  }
}