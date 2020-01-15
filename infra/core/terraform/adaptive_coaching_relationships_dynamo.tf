resource "aws_dynamodb_table" "coaching_relationships" {
  name           = "${var.client_id}_coaching_relationships"
  billing_mode = "PAY_PER_REQUEST"

  hash_key       = "coach_quarter_year"
  range_key      = "coachee"

  attribute {
    name = "coach_quarter_year"
    type = "S"
  }

  attribute {
    name = "coachee"
    type = "S"
  }

  attribute {
    name = "quarter"
    type = "N"
  }

  attribute {
    name = "year"
    type = "N"
  }

  attribute {
    name = "coachee_quarter_year"
    type = "S"
  }

  global_secondary_index {
    name            = var.dynamo_coaching_relationship_coach_index
    hash_key        = "coach_quarter_year"
    projection_type = "ALL"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
  }

  global_secondary_index {
    name            = var.dynamo_coaching_relationship_quarter_year_index
    hash_key        = "quarter"
    range_key       = "year"
    projection_type = "ALL"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
  }

  global_secondary_index {
    name            = var.dynamo_coaching_relationship_coachee_index
    hash_key        = "coachee_quarter_year"
    projection_type = "ALL"
    write_capacity  = var.dynamo_ondemand_write_capacity
    read_capacity  = var.dynamo_ondemand_read_capacity
  }

  tags = local.default_tags
}


