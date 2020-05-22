module "reporting_db" {
  source  = "terraform-aws-modules/rds/aws"
  version = "~> 2.0"

  identifier = "${local.client_id}-reporting"
  name = var.RDS_DB_NAME

  engine            = "mysql"
  engine_version    = "8.0.17"
  major_engine_version = "8.0"
  family = "mysql8.0"

  # engine            = "postgres"
  # engine_version    = "9.6.9"
  # major_engine_version = "9.6"
  # family = "postgres9.6"
  
  instance_class    = "db.t2.micro"
  allocated_storage = 5
  # storage_encrypted = false

  port     = var.RDS_PORT
  username = var.RDS_USER
  password = var.RDS_PASSWORD


  iam_database_authentication_enabled = true

  vpc_security_group_ids              = [
    aws_security_group.reporting.id,
  ]

  # DB subnet group
  subnet_ids = module.reporting_vpc.database_subnets

  multi_az = false

  # disable backups to create DB faster
  backup_retention_period = 0

  publicly_accessible = true

  # Snapshot name upon DB deletion
  final_snapshot_identifier = "${local.client_id}-reporting-snapshot-final"

  parameters = [
    {
      name  = "character_set_client"
      value = "utf8"
    },
    {
      name  = "character_set_server"
      value = "utf8"
    },
  ]

  maintenance_window = "Mon:00:00-Mon:03:00"
  backup_window      = "03:00-06:00"
}

# resource "null_resource" "mysql_bootstrap" {
#   depends_on = [
#     module.reporting_db.this_db_instance_id,
#   ]

# //  provisioner "local-exec" {
#   //    command = "mysql -u ${module.reporting_db.this_db_instance_username} -p ${module.reporting_db.this_db_instance_password} -h ${module.reporting_db.this_db_instance_endpoint} < bootstrap.sql"
#   //  }
# }