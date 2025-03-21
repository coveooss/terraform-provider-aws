---
subcategory: "WAF Classic"
layout: "aws"
page_title: "AWS: aws_waf_web_acl"
description: |-
  Retrieves a WAF Web ACL id.
---


<!-- Please do not edit this file, it is generated. -->
# Data Source: aws_waf_web_acl

`aws_waf_web_acl` Retrieves a WAF Web ACL Resource Id.

## Example Usage

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { DataAwsWafWebAcl } from "./.gen/providers/aws/data-aws-waf-web-acl";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    new DataAwsWafWebAcl(this, "example", {
      name: "tfWAFWebACL",
    });
  }
}

```

## Argument Reference

This data source supports the following arguments:

* `name` - (Required) Name of the WAF Web ACL.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `id` - ID of the WAF Web ACL.

<!-- cache-key: cdktf-0.20.8 input-d097cc64235699a72862520fb07194eebd805ec1ae4dcf4ceff43f0fb78e0e97 -->