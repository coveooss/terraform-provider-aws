---
subcategory: "CE (Cost Explorer)"
layout: "aws"
page_title: "AWS: aws_ce_cost_category"
description: |-
  Provides details about a specific CostExplorer Cost Category Definition
---


<!-- Please do not edit this file, it is generated. -->
# Resource: aws_ce_cost_category

Provides details about a specific CostExplorer Cost Category.

## Example Usage

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { DataAwsCeCostCategory } from "./.gen/providers/aws/data-aws-ce-cost-category";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    new DataAwsCeCostCategory(this, "example", {
      costCategoryArn: "costCategoryARN",
    });
  }
}

```

## Argument Reference

The following arguments are required:

* `costCategoryArn` - (Required) Unique name for the Cost Category.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `arn` - ARN of the cost category.
* `defaultValue` - Default value for the cost category.
* `effectiveEnd` - Effective end data of your Cost Category.
* `effectiveStart` - Effective state data of your Cost Category.
* `id` - Unique ID of the cost category.
* `rule` - Configuration block for the Cost Category rules used to categorize costs. See below.
* `ruleVersion` - Rule schema version in this particular Cost Category.
* `splitChargeRule` - Configuration block for the split charge rules used to allocate your charges between your Cost Category values. See below.
* `tags` - Resource tags.

### `rule`

* `inheritedValue` - Configuration block for the value the line item is categorized as if the line item contains the matched dimension. See below.
* `rule` - Configuration block for the `Expression` object used to categorize costs. See below.
* `type` - You can define the CostCategoryRule rule type as either `REGULAR` or `INHERITED_VALUE`.
* `value` - Default value for the cost category.

### `inheritedValue`

* `dimensionKey` - Key to extract cost category values.
* `dimensionName` - Name of the dimension that's used to group costs. If you specify `LINKED_ACCOUNT_NAME`, the cost category value is based on account name. If you specify `TAG`, the cost category value will be based on the value of the specified tag key. Valid values are `LINKED_ACCOUNT_NAME`, `TAG`

### `rule`

* `and` - Return results that match both `Dimension` objects.
* `costCategory` - Configuration block for the filter that's based on `CostCategory` values. See below.
* `dimension` - Configuration block for the specific `Dimension` to use for `Expression`. See below.
* `not` - Return results that do not match the `Dimension` object.
* `or` - Return results that match either `Dimension` object.
* `tags` - Configuration block for the specific `Tag` to use for `Expression`. See below.

### `costCategory`

* `key` - Unique name of the Cost Category.
* `matchOptions` - Match options that you can use to filter your results. MatchOptions is only applicable for actions related to cost category. The default values for MatchOptions is `EQUALS` and `CASE_SENSITIVE`. Valid values are: `EQUALS`,  `ABSENT`, `STARTS_WITH`, `ENDS_WITH`, `CONTAINS`, `CASE_SENSITIVE`, `CASE_INSENSITIVE`.
* `values` - Specific value of the Cost Category.

### `dimension`

* `key` - Unique name of the Cost Category.
* `matchOptions` - Match options that you can use to filter your results. MatchOptions is only applicable for actions related to cost category. The default values for MatchOptions is `EQUALS` and `CASE_SENSITIVE`. Valid values are: `EQUALS`,  `ABSENT`, `STARTS_WITH`, `ENDS_WITH`, `CONTAINS`, `CASE_SENSITIVE`, `CASE_INSENSITIVE`.
* `values` - Specific value of the Cost Category.

### `tags`

* `key` - Key for the tag.
* `matchOptions` - Match options that you can use to filter your results. MatchOptions is only applicable for actions related to cost category. The default values for MatchOptions is `EQUALS` and `CASE_SENSITIVE`. Valid values are: `EQUALS`,  `ABSENT`, `STARTS_WITH`, `ENDS_WITH`, `CONTAINS`, `CASE_SENSITIVE`, `CASE_INSENSITIVE`.
* `values` - Specific value of the Cost Category.

### `splitChargeRule`

* `method` - Method that's used to define how to split your source costs across your targets. Valid values are `FIXED`, `PROPORTIONAL`, `EVEN`
* `parameter` - Configuration block for the parameters for a split charge method. This is only required for the `FIXED` method. See below.
* `source` - Cost Category value that you want to split.
* `targets` - Cost Category values that you want to split costs across. These values can't be used as a source in other split charge rules.

### `parameter`

* `type` - Parameter type.
* `values` - Parameter values.

<!-- cache-key: cdktf-0.20.8 input-4cb2fbc95adf79676205ba698e97a322301998e777b3e00de5567e2938e38ac6 -->