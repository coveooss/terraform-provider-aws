package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func resourceAwsFlowLogMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AWS VPC Flow Log State v0; migrating to v1")
		return migrateAwsFlowLogStateV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateAwsFlowLogStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() || is.Attributes == nil {
		log.Println("[DEBUG] Empty AWS VPC Flow Log State; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	types := []struct {
		OldID string
		ID    string
		Type  string
	}{
		{ID: is.Attributes["vpc_id"], Type: ec2.ResourceTypeVpc, OldID: "vpc_id"},
		{ID: is.Attributes["subnet_id"], Type: ec2.ResourceTypeSubnet, OldID: "subnet_id"},
		{ID: is.Attributes["eni_id"], Type: ec2.ResourceTypeNetworkInterface, OldID: "eni_id"},
	}

	for _, t := range types {
		if t.ID == "" {
			continue
		}

		is.Attributes["resource_type"] = t.Type

		entity := resourceAwsFlowLog()
		writer := schema.MapFieldWriter{
			Schema: entity.Schema,
		}

		// Convert the old format that restricted to a single resource to
		// the new format that supports a list of resources
		if err := writer.WriteField([]string{"resource_ids"}, []string{is.Attributes[t.OldID]}); err != nil {
			return is, err
		}

		for k, v := range writer.Map() {
			is.Attributes[k] = v
		}

		delete(is.Attributes, t.OldID)

	}

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
