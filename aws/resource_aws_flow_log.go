package aws

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceAwsFlowLog() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsLogFlowCreate,
		Read:   resourceAwsLogFlowRead,
		Delete: resourceAwsLogFlowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		MigrateState:  resourceAwsFlowLogMigrateState,
		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{

			"eni_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ConflictsWith:    []string{"subnet_id", "vpc_id", "resource_type", "resource_ids"},
				Deprecated:       "Attribute eni_id is deprecated on aws_flow_log resources. Use resource_type in combination with resource_ids instead.",
				DiffSuppressFunc: diffIsSameResourcesForType(ec2.FlowLogsResourceTypeNetworkInterface),
			},

			"subnet_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ConflictsWith:    []string{"eni_id", "vpc_id", "resource_type", "resource_ids"},
				Deprecated:       "Attribute subnet_id is deprecated on aws_flow_log resources. Use resource_type in combination with resource_ids instead.",
				DiffSuppressFunc: diffIsSameResourcesForType(ec2.FlowLogsResourceTypeSubnet),
			},

			"vpc_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ConflictsWith:    []string{"eni_id", "subnet_id", "resource_type", "resource_ids"},
				Deprecated:       "Attribute vpc_id is deprecated on aws_flow_log resources. Use resource_type in combination with resource_ids instead.",
				DiffSuppressFunc: diffIsSameResourcesForType(ec2.FlowLogsResourceTypeVpc),
			},

			"iam_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"log_destination": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateArn,
			},

			"log_destination_type": {
				Type:     schema.TypeString,
				Default:  ec2.LogDestinationTypeCloudWatchLogs,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					ec2.LogDestinationTypeCloudWatchLogs,
					ec2.LogDestinationTypeS3,
				}, false),
			},

			"log_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Suppress diff if newer log_destination argument is the same as this legacy value
					logGroupName := d.Get("log_destination")
					logDestType := d.Get("log_destination_type")
					return logGroupName == d.Get("log_group_name") && logDestType == ec2.LogDestinationTypeCloudWatchLogs
				},
			},

			"resource_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				MinItems: 1,
				MaxItems: 1000,
				Optional: true, // should be switched to Required when old format is deprecated
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Suppress diff if newer resource ids are the same as the legacy value
					_, resourceIDs := readLegacyResourceTypeAndIDs(d)
					if resourceIDs == nil {
						return false
					}

					return len(resourceIDs) == 1 && *resourceIDs[0] == old
				},
			},

			"resource_type": {
				Type:     schema.TypeString,
				Optional: true, // should be switched to Required when old format is deprecated
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					ec2.FlowLogsResourceTypeVpc,
					ec2.FlowLogsResourceTypeSubnet,
					ec2.FlowLogsResourceTypeNetworkInterface,
				}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Suppress diff if newer resource type is the same as the legacy value
					resourceType, _ := readLegacyResourceTypeAndIDs(d)
					return resourceType != nil && *resourceType == old
				},
			},

			"traffic_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					ec2.TrafficTypeAccept,
					ec2.TrafficTypeReject,
					ec2.TrafficTypeAll,
				}, false),
			},
		},
	}
}

func resourceAwsLogFlowCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn

	var logDestType, logDest *string
	resourceType, resourceIDs := readLegacyResourceTypeAndIDs(d)

	if len(resourceIDs) == 0 || *resourceType == "" {
		logDest = aws.String(d.Get("log_destination").(string))
		logDestType = aws.String(d.Get("log_destination_type").(string))
		resourceIDs = expandStringList(d.Get("resource_ids").(*schema.Set).List())
		resourceType = aws.String(d.Get("resource_type").(string))
	}

	if len(resourceIDs) == 0 || *resourceType == "" {
		return fmt.Errorf("Error: Flow Logs require a VPC, Subnet, or ENI ID AND a list of one or more IDs")
	}

	opts := &ec2.CreateFlowLogsInput{
		DeliverLogsPermissionArn: aws.String(d.Get("iam_role_arn").(string)),
		LogGroupName:             aws.String(d.Get("log_group_name").(string)),
		ResourceIds:              resourceIDs,
		ResourceType:             resourceType,
		TrafficType:              aws.String(d.Get("traffic_type").(string)),
	}

	if logDest != nil && *logDest != "" && logDestType != nil && *logDestType != "" {
		opts.LogDestination = logDest
		opts.LogDestinationType = logDestType
	}

	log.Printf("[DEBUG] Flow Log Create configuration: %s", opts)

	resp, err := conn.CreateFlowLogs(opts)
	if err != nil {
		return fmt.Errorf("Error creating Flow Log for (%s), error: %s", resourceId, err)
	}

	if len(resp.FlowLogIds) > 1 {
		return fmt.Errorf("Error: multiple Flow Logs created for (%s)", resourceId)
	}

	d.SetId(*resp.FlowLogIds[0])

	return resourceAwsLogFlowRead(d, meta)
}

func resourceAwsLogFlowRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn

	opts := &ec2.DescribeFlowLogsInput{
		FlowLogIds: []*string{aws.String(d.Id())},
	}

	resp, err := conn.DescribeFlowLogs(opts)
	if err != nil {
		log.Printf("[WARN] Error describing Flow Logs for id (%s)", d.Id())
		d.SetId("")
		return nil
	}

	if len(resp.FlowLogs) == 0 {
		log.Printf("[WARN] No Flow Logs found for id (%s)", d.Id())
		d.SetId("")
		return nil
	}

	fl := resp.FlowLogs[0]
	d.Set("traffic_type", fl.TrafficType)
	d.Set("log_group_name", fl.LogGroupName)
	d.Set("iam_role_arn", fl.DeliverLogsPermissionArn)
	d.Set("log_destination", fl.LogDestination)
	d.Set("log_destination_type", fl.LogDestinationType)
	d.Set("resource_ids", schema.NewSet(schema.HashString, []interface{}{*fl.ResourceId}))

	if strings.HasPrefix(*fl.ResourceId, "vpc-") {
		d.Set("resource_type", ec2.FlowLogsResourceTypeVpc)
	} else if strings.HasPrefix(*fl.ResourceId, "subnet-") {
		d.Set("resource_type", ec2.FlowLogsResourceTypeSubnet)
	} else if strings.HasPrefix(*fl.ResourceId, "eni-") {
		d.Set("resource_type", ec2.FlowLogsResourceTypeNetworkInterface)
	}

	return nil
}

func resourceAwsLogFlowDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn

	log.Printf(
		"[DEBUG] Flow Log Destroy: %s", d.Id())
	_, err := conn.DeleteFlowLogs(&ec2.DeleteFlowLogsInput{
		FlowLogIds: []*string{aws.String(d.Id())},
	})

	if err != nil {
		return fmt.Errorf("Error deleting Flow Log with ID (%s), error: %s", d.Id(), err)
	}

	return nil
}

func diffIsSameResourcesForType(resourceType string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		v, ok := d.GetOk("resource_type")
		if ok {
			if v.(string) == resourceType {
				r, ok := d.GetOk("resource_ids")
				if ok {
					resourceIDs := expandStringList(r.(*schema.Set).List())
					if len(resourceIDs) == 0 && *resourceIDs[0] == old {
						return true
					}
				}
			}
		}
		return false
	}
}

func readLegacyResourceTypeAndIDs(d *schema.ResourceData) (resourceType *string, resourceIDs []*string) {
	types := []struct {
		ID   string
		Type string
	}{
		{ID: d.Get("vpc_id").(string), Type: ec2.FlowLogsResourceTypeVpc},
		{ID: d.Get("subnet_id").(string), Type: ec2.FlowLogsResourceTypeSubnet},
		{ID: d.Get("eni_id").(string), Type: ec2.FlowLogsResourceTypeNetworkInterface},
	}
	for _, t := range types {
		if t.ID != "" {
			resourceIDs = []*string{aws.String(t.ID)}
			resourceType = aws.String(t.Type)
			break
		}
	}
	return
}
