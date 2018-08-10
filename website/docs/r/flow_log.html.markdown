---
layout: "aws"
page_title: "AWS: aws_flow_log"
sidebar_current: "docs-aws-resource-flow-log"
description: |-
  Provides a VPC/Subnet/ENI Flow Log
---

# aws_flow_log

Provides a VPC, Subnet, or ENI Flow Log to capture IP traffic for a list of VPC,
subnet, or network interface IDs. Logs are sent to either a CloudWatch Log Group
or an S3 bucket.

## Example Usage

```hcl
resource "aws_flow_log" "test_flow_log" {
  log_group_name = "${aws_cloudwatch_log_group.test_log_group.name}"
  iam_role_arn   = "${aws_iam_role.test_role.arn}"
  vpc_id         = "${aws_vpc.default.id}"
  traffic_type   = "ALL"
}

resource "aws_cloudwatch_log_group" "test_log_group" {
  name = "test_log_group"
}

resource "aws_iam_role" "test_role" {
  name = "test_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "vpc-flow-logs.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "test_policy" {
  name = "test_policy"
  role = "${aws_iam_role.test_role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents",
        "logs:DescribeLogGroups",
        "logs:DescribeLogStreams"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}
```

## Argument Reference

The following arguments are supported:

* `resource_type` - (Required) The type of resource on which to create the flow log.
  Valid Values: `VPC`, `Subnet`, `NetworkInterface`
* `resource_id` - (Required) One or more subnet, network interface, or VPC IDs.
  Constraints: Maximum of 1000 resources
* `traffic_type` - (Required) The type of traffic to capture. Valid values:
  `ACCEPT`, `REJECT`, `ALL`
* `log_destination_type` - (Optional) Type of destination to which the flow log
  data is to be published. Flow log data can be published to CloudWatch Logs or
  Amazon S3. To publish flow log data to CloudWatch Logs, specify `cloud-watch-logs`.
  To publish flow log data to Amazon S3, specify `s3`. Default: `cloud-watch-logs`
* `log_destination` - (Optional) The destination to which the flow log data is to
  be published. Flow log data can be published to an CloudWatch Logs log group or
  an Amazon S3 bucket. The value specified for this parameter depends on the value
  specified for `log_destination_type`.
* `log_group_name` - (Optional) The name of the CloudWatch Log group to which the
  flow log data is to be published.
* `iam_role_arn` - (Optional) The ARN for the IAM role that's used to post flow
  logs to a CloudWatch Logs log group. **Required** if the sending logs to
  CloudWatch Logs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Flow Log ID

## Import

Flow Logs can be imported using the `id`, e.g.

```
$ terraform import aws_flow_log.test_flow_log fl-1a2b3c4d
```
