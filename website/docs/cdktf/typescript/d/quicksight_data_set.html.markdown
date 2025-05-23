---
subcategory: "QuickSight"
layout: "aws"
page_title: "AWS: aws_quicksight_data_set"
description: |-
  Use this data source to fetch information about a QuickSight Data Set.
---


<!-- Please do not edit this file, it is generated. -->
# Data Source: aws_quicksight_data_set

Data source for managing a QuickSight Data Set.

## Example Usage

### Basic Usage

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { DataAwsQuicksightDataSet } from "./.gen/providers/aws/data-aws-quicksight-data-set";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    new DataAwsQuicksightDataSet(this, "example", {
      dataSetId: "example-id",
    });
  }
}

```

## Argument Reference

This data source supports the following arguments:

* `dataSetId` - (Required) Identifier for the data set.
* `awsAccountId` - (Optional) AWS account ID.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

See the [Data Set Resource](/docs/providers/aws/r/quicksight_data_set.html) for details on the
returned attributes - they are identical.

<!-- cache-key: cdktf-0.20.8 input-33b66a824e4abf0d8fd597d57ed0f5abb63f7e0300c1c2a1ccfe1ba1c233a1ee -->