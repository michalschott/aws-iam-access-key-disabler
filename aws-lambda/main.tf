variable "low_threshold_days" {
  default = 90
}

variable "high_threshold_days" {
  default = 180
}

variable "dry_run" {
  default = 1
}

variable "whitelist" {
  default = ""
}

variable "schedule" {
  default = "rate(1 hour)"
}

variable "lambda_timeout" {
  default = 15
}

variable "lambda_memory_size" {
  default = 128
}

data "aws_caller_identity" "this" {}

data "aws_region" "this" {}

resource "aws_lambda_function" "this" {
  filename         = "${path.module}/aws-iam-access-key-disabler.zip"
  function_name    = "aws-iam-access-key-disabler"
  role             = "${aws_iam_role.this.arn}"
  handler          = "aws-iam-access-key-disabler"
  source_code_hash = "${base64sha256(file("${path.module}/aws-iam-access-key-disabler.zip"))}"
  runtime          = "go1.x"
  timeout          = "${var.lambda_timeout}"
  memory_size      = "${var.lambda_memory_size}"

  environment {
    variables = {
      LOWTHRESHOLDDAYS  = "${var.low_threshold_days}"
      HIGHTHRESHOLDDAYS = "${var.high_threshold_days}"
      DRYRUN            = "${var.dry_run}"
      WHITELIST         = "${var.whitelist}"
    }
  }
}

resource "aws_lambda_permission" "this" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.this.function_name}"
  principal     = "events.amazonaws.com"
  source_arn    = "${aws_cloudwatch_event_rule.this.arn}"
}

resource "aws_iam_role" "this" {
  name = "aws-iam-access-key-disabler"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_policy" "this" {
  name        = "aws-iam-access-key-disabler"
  path        = "/"
  description = "Policy allowing lambda function aws-iam-access-key-disabler to disable IAM keys."
  policy      = "${data.aws_iam_policy_document.this.json}"
}

resource "aws_iam_role_policy_attachment" "this-attachment" {
  role       = "${aws_iam_role.this.name}"
  policy_arn = "${aws_iam_policy.this.arn}"
}

data "aws_iam_policy_document" "this" {
  "statement" {
    effect = "Allow"

    actions = [
      "logs:CreateLogGroup",
    ]

    resources = ["arn:aws:logs:${data.aws_region.this.name}:${data.aws_caller_identity.this.account_id}:*"]
  }

  "statement" {
    effect = "Allow"

    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = ["arn:aws:logs:${data.aws_region.this.name}:${data.aws_caller_identity.this.account_id}:log-group:/aws/lambda/aws-iam-access-key-disabler:*"]
  }

  "statement" {
    effect = "Allow"

    actions = [
      "iam:UpdateAccessKey",
      "iam:ListUsers",
      "iam:ListAccessKeys",
    ]

    resources = ["*"]
  }
}

resource "aws_cloudwatch_event_rule" "this" {
  name                = "aws-iam-access-key-disabler"
  description         = "Event rule for scheduling aws-iam-access-key-disabler lambda function."
  schedule_expression = "${var.schedule}"
}

resource "aws_cloudwatch_event_target" "this" {
  rule      = "${aws_cloudwatch_event_rule.this.name}"
  target_id = "aws-iam-access-key-disabler"
  arn       = "${aws_lambda_function.this.arn}"
}
