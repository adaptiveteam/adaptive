# core-infra-terraform [![Build Status](https://travis-ci.com/adaptiveteam/core-infra-terraform.svg?token=BSM7265i3ndP9kG2qsqY&branch=develop)](https://travis-ci.com/adaptiveteam/core-infra-terraform)
Adaptive core infrastructure using terraform

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
  source_hash = "${base64sha256(filebase64(<filepath>))}"
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

## Execution

### Infrastructure
Terraform scripts are present inside `terraform` folder. We can deploy the infrastructure by running terraform in that folder. This folder is also packaged as part of the build generated zip file.

### Remote state
We use S3 to store state of terraform execution, by using locking on the state. With locking, only one person will be able to make modifications to the existing infrastructure at a time. This ensure accidental override of the state when 2 or more services are simultaneously updating the state.

We need to specify the backend S3 configuration during initialization of the terraform. 

```bash
./backend/init.sh
```
(You should configure `ADAPTIVE_CLIENT_ID` in your `.bashrc`.)

With this terraform will be initialized with provided backend config. This decoupling enables us to run against different buckets when working with prod/dev/testing environments.


### Testing

We currently have an end-to-end integration test that deploys the infrastructure into a region specified with provided S3 backend. We accomplish this using [terratest](https://github.com/gruntwork-io/terratest) which provides Go API to initiate terraform commands. Once infrastructure is deployed, we have a suite of tests that are run against this deployment. Once tests are done, infrastructure is destroyed.


## Lambda flow

![Flow diagram](/files/Adaptive_Core.png?raw=true "Optional Title")

## 