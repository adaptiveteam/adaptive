// See https://github.com/cloudposse/terraform-aws-efs-backup/blob/master/iam.tf
data "aws_iam_policy_document" "assume_role" {
  statement {
    sid     = "EC2AssumeRole"
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = [
        "ec2.amazonaws.com",
        "elasticmapreduce.amazonaws.com",
        "datapipeline.amazonaws.com"]
    }
  }
}


# module "resource_role_label" {
#   source     = "git::https://github.com/cloudposse/terraform-null-label.git?ref=tags/0.3.1"
#   namespace  = "${var.namespace}"
#   stage      = "${var.stage}"
#   name       = "${var.name}"
#   delimiter  = "${var.delimiter}"
#   attributes = "${var.attributes}"
#   tags       = "${var.tags}"
# }

resource "aws_iam_role" "resource_role" {
  name               = "${var.client_id}_backup_resource_role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_role_policy_attachment" "resource_role_DataPipeline_policy" {
  role       = aws_iam_role.resource_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforDataPipelineRole"
}

resource "aws_iam_instance_profile" "resource_role" {
  name = "${var.client_id}_backup_resource_role"
  role = aws_iam_role.resource_role.name
}
resource "aws_iam_instance_profile" "role" {
  name = "${var.client_id}_backup_role"
  role = aws_iam_role.role.name
}

# data "aws_iam_policy_document" "role" {
#   statement {
#     sid     = "AssumeRole"
#     effect  = "Allow"
#     actions = ["sts:AssumeRole"]

#     principals {
#       type = "Service"

#       identifiers = [
#         "ec2.amazonaws.com",
#         "elasticmapreduce.amazonaws.com",
#         "datapipeline.amazonaws.com",
#       ]
#     }
#   }
# }

# module "role_label" {
#   source     = "git::https://github.com/cloudposse/terraform-null-label.git?ref=tags/0.3.1"
#   namespace  = "${var.namespace}"
#   stage      = "${var.stage}"
#   name       = "${var.name}"
#   delimiter  = "${var.delimiter}"
#   attributes = ["${compact(concat(var.attributes, list("role")))}"]
#   tags       = "${var.tags}"
# }

resource "aws_iam_role" "role" {
  name               = "${var.client_id}_backup_role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_role_policy_attachment" "role" {
  role       = aws_iam_role.role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSDataPipelineRole"
}
//ser: arn:aws:sts::221851954636:assumed-role/backup_my_resource_role/EDPSession is not authorized to perform:
// elasticmapreduce:ListClusters on resource: * (Service: AmazonElasticMapReduce; Status Code: 400; Error Code: AccessDeniedException; Request ID: 9a8969f7-f62a-4e23-9fb6-62db979cabf3)

data "aws_iam_policy_document" "allow_role_to_start_emr_cluster" {
  statement {
    sid     = "StartEmrCluster"
    effect  = "Allow"
    actions = ["elasticmapreduce:*",
      "iam:PassRole", "iam:CreateRole", "iam:*"]

    # iam:CreateRole
    # iam:PutRolePolicy
    # iam:CreateInstanceProfile
    # iam:AddRoleToInstanceProfile
    # iam:ListRoles
    # iam:GetPolicy
    # iam:GetInstanceProfile
    # iam:GetPolicyVersion
    # iam:AttachRolePolicy
    # iam:PassRole

    resources = ["*"]

  }
  statement {
    sid     = "AccessToBackupBucket"
    effect  = "Allow"
    actions = ["s3:*"]
    resources = [
      "${aws_s3_bucket.backup_bucket.arn}",
      "${aws_s3_bucket.backup_bucket.arn}/*",
      ]
  }
  statement {
    sid = "AllowEC2Operations"
    effect = "Allow"
    actions = ["ec2:*"]
    resources = ["*"]
  }
}
// ws:sts::221851954636:assumed-role/backup_my_resource_role/EDPSession is not authorized to perform:
// iam:PassRole on resource: arn:aws:iam::221851954636:role/backup_my_resource_role (Service: AmazonElasticMapReduce; Status Code: 400; Error Code: AccessDeniedException; Request ID: bf7820f4-2fff-4aef-a81e-21d9d92611d0)

resource "aws_iam_policy" "role_DataPipeline_policy" {
  name   = "${var.client_id}_role_DataPipeline_policy"
  policy = data.aws_iam_policy_document.allow_role_to_start_emr_cluster.json
}

resource "aws_iam_policy_attachment" "resource_role_DataPipeline_policy" {
  name       = "${var.client_id}_role_DataPipeline_policy_attachment"
  roles      = [aws_iam_role.role.name]
  policy_arn = aws_iam_policy.role_DataPipeline_policy.arn
}