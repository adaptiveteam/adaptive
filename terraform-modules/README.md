# adaptive-terraform-modules
Terraform modules that can be used for Adaptive infrastructure

## Modules

### adaptive-dynamo-autoscaler
This module enables to apply auto-scaling to a dynamo table.

#### Parameters
- `table_name`
    - Required: Yes
    - Description: Name of the table to apply auto-scaling to
- `table_arn`
    - Required: Yes
    - Description: ARN of the table toa apply auto-scaling to
- `min_read_capacity`
    - Required: No
    - Default value: 5
    - Description: Minimum read capacity of the table
- `max_read_capacity`
    - Required: No
    - Default value: 5
    - Description: Maximum read capacity of the table
- `min_write_capacity`
    - Required: No
    - Default value: 5
    - DescriptionL Minimum write capacity of the table
- `max_write_capacity`
    - Required: No
    - Default value: 5
    - Description: Maximum write capacity of the table
- `target_tracking_read_threshold`
    - Required: No
    - Default value: 70
    - Description: Target tracking threshold for read capacity. When read capacity reaches 70% of the specified value, it scales up
- `target_tracking_write_threshold`
    - Required: No
    - Default value: 70
    - Description: Target tracking threshold for write capacity. When write capacity reaches 70% of the specified value, it scales up
- `client_id`
    - Required: Yes
    - Description: Unique id of the client

##### Usage

```$xslt
module "table_scaling" {
  source = "../terraform/modules/adaptive-dynamo-autoscaler"

  client_id = "${var.client_id}"
  table_name = "${aws_dynamodb_table.<table>.name}"
  table_arn = "${aws_dynamodb_table.<table>.arn}"
}
```

### adaptive-lambda
This module enables to create a lambda without specifying much of boiler-plate code required to create roles, logs and policies

#### Parameters
- `table_name`
    - Required: Yes
    - Description: Name of the lambda function
- `description`
    - Required: No
    - Default value: Managed by Terraform
    - Description: Description of the lambda function
- `filename`
    - Required: Yes
    - Description: Filename of the lambda function that should be uploaded as a zip
- `source_hash`
    - Required: Yes
    - Description: Hash of the file that is being uploaded
- `function_name`
    - Required: Yes
    - Description: Name of the lambda function
- `handler`
    - Required: Yes
    - Description: Lambda handler for the function
- `runtime`
    - Required: No
    - Default value: go1.x
    - Description: Runtime of the lambda
- `timeout`
    - Required: No
    - Default value: 30
    - Description: Timeout for the lambda function
- `environment`
    - Required: No
    - Default value: {}
    - Description: Environment variables for the lambda function
- `memory_size`
    - Required: No
    - Default: 128
    - Description: Memory allocated for the lambda
- `attach_vpc_config`
    - Required: No
    - Default value: false
    - Description: Flag to indicate if vpc should be associated with the lambda
- `attach_dl_config`
    - Required: No
    - Default value: false
    - Description: Flag to indicate if dead letter should be associated with the lambda
- `reserved_concurrent_executions`
    - Required: No
    - Default: 0
    - Description: Number of reserved concurrent executions for the lambda
- `tags`
    - Required: No
    - Default value: <global tags>
    - Description: Tags to attach to the lambda
- `dl_config`
    - Required: No
    - Default value: {}
    - Description: Dead letter config for the lambda. This is required when `attach_dl_config` is set to true
- `vpc_config`
    - Required: No
    - Default value: {}
    - Description: VPC config for the lambda. This is required when `attach_vpc_config` is set to true
- `enable_cloudwatch_logs`
    - Required: No
    - Default: true
    - Description: Flag to enable lambda logs
- `attach_policy`
    - Required: No
    - Default value: false
    - Description: Flag to indicate if extra IAM policy should be attached with the lambda role
- `policy`
    - Required: No
    - Default value: ""
    - Description: Extra policy to attach to lambda. This should be specified when `attach_policy` is set to true.
- `schedule`    
    - Required: No
    - Default value: false
    - Description: Flag to indicate if lambda should be scheduled
- `schedule_name`
    - Required: No
    - Default value: ""
    - Description: Name for the lambda schedule. This is required when `schedule` is set to true.
- `schedule_description`
    - Required: No
    - Default value: ""
    - Description: Description for the lambda schedule. This is required when `schedule` is set to true.
- `schedule_expression`
    - Required: No
    - Default value: ""
    - Description: Cron expression for the lambda schedule. This is required when `schedule` is set to true.
- `client_id`
    - Required: Yes
    - Description: Unique id for the client
   
##### Usage

```$xslt
module "adaptive_user_settings_lambda" {
  source = "../terraform/modules/adaptive-lambda"

  client_id = "${var.client_id}"
  filename = "<filename>"
  source_hash = "${base64sha256(file(<filepath>))}"
  function_name = "<function-name>"
  handler = "<function-handler>"
  runtime = "${var.lambda_runtime}"
  timeout = "${var.lambda_timeout}"

  environment {
    variables = {
        "foo" = "var"
    }
  }

  // Attach extra policy
  attach_policy = true
  policy = "${data.aws_iam_policy_document.adaptive_user_settings_policy.json}"

  tags = "${var.global_tags}"
}
```

- Default lambda module creates lambda with a role that is configured to write to appropriate log streams
- Additional policies specific to the lambda can be added by setting `attach_policy` to `true` and passing in a policy expression which can be constructed like below

```$xslt
data "aws_iam_policy_document" "policy" {
  statement {
    effect = "Allow"
    actions = [
      "dynamodb:GetRecords",
      "dynamodb:GetShardIterator",
      "dynamodb:DescribeStream",
      "dynamodb:ListStreams",]
    resources = [
      "<arn>",]
  }
}
```

### adaptive-go-lambda

This module enables to create a lambda from golang code and subscribe it
to SNS topic.

#### Parameters

See [variables](adaptive-go-lambda/variables.tf)

#### Usage


```$xslt
module "adaptive_holidays_lambda" {
    name = "adaptive-holidays-lambda"
    modpath = "../adaptive-user-community-lambdas/adaptive-holidays-lambda/src/main/golang"
    id = "client_id"
    input_sns_topic_arn = "${aws_sns_topic.namespace_payload.arn}"
    tags = "${var.global_tags}"
}
```

It'll create a lambda based on the go module in the given folder. Also this lambda will be subscribed to the SNS topic.
