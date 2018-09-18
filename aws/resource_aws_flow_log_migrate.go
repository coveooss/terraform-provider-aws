package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform/terraform"
)

func resourceAwsFlowLogMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AWS VPC Flow Log State v0; migrating to v1")
		return migrateAwsFlowLogStateV0toV1(is, meta)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateAwsFlowLogStateV0toV1(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() || is.Attributes == nil {
		log.Println("[DEBUG] Empty AWS VPC Flow Log State; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	// Migrate resource ID.
	resourceTypes := []struct {
		OldKey  string
		NewType string
	}{
		{OldKey: "vpc_id", NewType: ec2.FlowLogsResourceTypeVpc},
		{OldKey: "subnet_id", NewType: ec2.FlowLogsResourceTypeSubnet},
		{OldKey: "eni_id", NewType: ec2.FlowLogsResourceTypeNetworkInterface},
	}
	for _, t := range resourceTypes {
		resourceId := is.Attributes[t.OldKey]
		if resourceId != "" {
			is.Attributes["resource_id"] = resourceId
			is.Attributes["resource_type"] = t.NewType
		}
		delete(is.Attributes, t.OldKey)
	}

	// Migrate log destination.
	// Convert CloudWatch log group name to ARN.
	arn := arn.ARN{
		Partition: meta.(*AWSClient).partition,
		Region:    meta.(*AWSClient).region,
		Service:   "logs",
		AccountID: meta.(*AWSClient).accountid,
		Resource:  fmt.Sprintf("log-group:%s:*", is.Attributes["log_group_name"]),
	}.String()
	is.Attributes["log_destination"] = arn
	is.Attributes["log_destination_type"] = ec2.LogDestinationTypeCloudWatchLogs
	delete(is.Attributes, "log_group_name")

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)

	return is, nil
}
