import $file.Meta
import Meta._

import $file.Templates
import Templates._

import $file.DynamoTemplates
import DynamoTemplates._

def tfName(name: SimpleName): String = snakeCaseName(name)

def tfField(field: Field): List[String] = 
lines(s"""attribute {
    name = "${dynamoName(field.dbName)}"
    type = "${dynamoType(field.tpe)}"
}
""")

def tfSecondaryIndex(defaultIndex: Index)(index: Index): List[String] = 
  if(defaultIndex.hashKey == index.hashKey && 
     index.rangeKey.isDefined) // range key is mandatory for local secondary index
    tfLocalSecondaryIndex(index)
  else
    tfGlobalSecondaryIndex(index)

def tfGlobalSecondaryIndex(index: Index): List[String] =
  // val indexName = Name("var.dynamo" :: index.name.parts)
  blockNamed("global_secondary_index", lines(
    s"""name            = "${indexName(index)}"
       |hash_key        = "${tfName(index.hashKey.dbName)}"
       |${index.rangeKey.map(rk => "range_key       = \"" + tfName(rk.dbName) + "\"").getOrElse("") }
       |projection_type = "${projectionTypeTemplate(index.projectionType)}"
       |write_capacity  = var.dynamo_ondemand_write_capacity
       |read_capacity   = var.dynamo_ondemand_read_capacity
       |""".stripMargin
  ))

def projectionTypeTemplate(projectionType: ProjectionType): String = projectionType match {
  case ProjectionType.ALL => "ALL"
  case ProjectionType.INCLUDE => "INCLUDE"
}

def nonKeyAttributesTemplate(nonKeyAttributes: List[Field]): List[String] =
  if(nonKeyAttributes.isEmpty)
    Nil
  else
    bracketBlockNamed("non_key_attributes = ", nonKeyAttributes.map(_.dbName).map(tfName).map("\"" + _ + "\","))

def tfLocalSecondaryIndex(index: Index): List[String] =
  blockNamed("local_secondary_index",
    lines(
      s"""name            = "${goPublicName(index.name)}"
         |projection_type = "${projectionTypeTemplate(index.projectionType)}"
         |${index.rangeKey.map(rk => "range_key       = \"" + tfName(rk.dbName) + "\"").getOrElse("") }
         |""".stripMargin
    ) :::
    nonKeyAttributesTemplate(index.nonKeyAttributes)
  )

def aws_dynamodb_table(table: Table): List[String] = {
  blockNamed(s"""resource "aws_dynamodb_table" "${tfName(table.entity.name)}_dynamodb_table" """,
    lines(
      s"""name           = "$${var.client_id}_${tfName(table.entity.name)}"
         |billing_mode   = "PAY_PER_REQUEST"
         |
         |tags           = local.default_tags
         |hash_key       = "${tfName(table.hashKey.dbName)}"
         |""".stripMargin
    ) ::: 
    table.defaultIndex.rangeKey.map(rk => 
      s"range_key      = " + quote(tfName(rk.dbName))
    ).toList :::
    blockNamed("point_in_time_recovery", lines("enabled = true")) :::
    (if(table.encrypted) 
      blockNamed("server_side_encryption", lines("enabled = true"))
    else Nil)
    :::
    table.fieldsInIndices.flatMap(tfField)
    :::
    table.indices.flatMap(tfSecondaryIndex(table.defaultIndex))
  ) 
}

def adaptive_table_scaling(table: Table): List[String] = lines(
  s"""
  |module "${tfName(table.entity.name)}_table_scaling" {
  |    source = "github.com/adaptiveteam/adaptive-terraform-modules//adaptive-dynamo-autoscaler?ref=v1.0.0"
  |  
  |    client_id  = var.client_id
  |    table_name = aws_dynamodb_table.${tfName(table.entity.name)}_dynamodb_table.name
  |    table_arn  = aws_dynamodb_table.${tfName(table.entity.name)}_dynamodb_table.arn
  |}  
  |""".stripMargin
)

def adaptive_backup(table: Table): List[String] = lines(
  s"""module "backup_${tfName(table.entity.name)}" {
     |  source = "github.com/adaptiveteam/adaptive-terraform-modules//adaptive-backup?ref=backup-module"
     |  bucket = aws_s3_bucket.backup_bucket.bucket
     |  name_prefix = "$${local.client_id}${goPublicName(table.entity.name)}"
     |  role_name = aws_iam_role.role.name
     |  resource_role_name = aws_iam_role.resource_role.name
     |  table_name = aws_dynamodb_table.${tfName(table.entity.name)}_dynamodb_table.name
     |}
     |""".stripMargin)

def output_name(table: Table): List[String] = 
  blockNamed(
    s"""output "${tfName(table.entity.name)}_table_name"""",
    List(
      s"""description = "Name of the ${tfName(table.entity.name)} table"""",
      s"""value = aws_dynamodb_table.${tfName(table.entity.name)}_dynamodb_table.name"""
    )
  )

def output_arn(table: Table): List[String] = 
  blockNamed(
    s"""output "${tfName(table.entity.name)}_table_arn"""",
    List(
      s"""description = "ARN of the ${tfName(table.entity.name)} table"""",
      s"""value = aws_dynamodb_table.${tfName(table.entity.name)}_dynamodb_table.arn"""
    )
  )

def terraformTemplate(table: Table): List[String] = {
    val tbl = aws_dynamodb_table(table) 
    // val scaling = adaptive_table_scaling(table)
    // val backup = adaptive_backup(table)
    val arn = output_arn(table)
    val outName = output_name(table)
    tbl ::: arn ::: outName // ::: scaling ::: backup
}
