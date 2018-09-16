package aws

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
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
		CustomizeDiff: resourceAwsLogFlowCustomizeDiff,

		MigrateState:  resourceAwsFlowLogMigrateState,
		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"iam_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"log_group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"log_destination", "log_destination_type"},
				Deprecated:    "Attribute log_group_name is deprecated on aws_flow_log resources. Use log_destination_type in combination with log_destination instead.",
			},

			"log_destination": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"log_group_name"},
				ValidateFunc:  validateArn,
			},

			"log_destination_type": {
				Type: schema.TypeString,
				// Default must be nil if computed.
				//Default:       ec2.LogDestinationTypeCloudWatchLogs,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"log_group_name"},
				ValidateFunc: validation.StringInSlice([]string{
					ec2.LogDestinationTypeCloudWatchLogs,
					ec2.LogDestinationTypeS3,
				}, false),
			},

			"vpc_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"eni_id", "subnet_id", "resource_id", "resource_type"},
				Deprecated:    "Attribute vpc_id is deprecated on aws_flow_log resources. Use resource_type in combination with resource_id instead.",
			},

			"eni_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"subnet_id", "vpc_id", "resource_id", "resource_type"},
				Deprecated:    "Attribute eni_id is deprecated on aws_flow_log resources. Use resource_type in combination with resource_id instead.",
			},

			"subnet_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"eni_id", "vpc_id", "resource_id", "resource_type"},
				Deprecated:    "Attribute subnet_id is deprecated on aws_flow_log resources. Use resource_type in combination with resource_id instead.",
			},

			"resource_id": {
				Type:          schema.TypeString,
				Optional:      true, // should be switched to Required when old format is deprecated
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"eni_id", "subnet_id", "vpc_id"},
			},

			"resource_type": {
				Type:          schema.TypeString,
				Optional:      true, // should be switched to Required when old format is deprecated
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"eni_id", "subnet_id", "vpc_id"},
				ValidateFunc: validation.StringInSlice([]string{
					ec2.FlowLogsResourceTypeVpc,
					ec2.FlowLogsResourceTypeSubnet,
					ec2.FlowLogsResourceTypeNetworkInterface,
				}, false),
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

			"flow_log_status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"deliver_logs_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAwsLogFlowCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn

	req := &ec2.CreateFlowLogsInput{
		TrafficType: aws.String(d.Get("traffic_type").(string)),
	}

	if v := d.Get("log_destination").(string); v != "" {
		req.LogDestination = aws.String(v)
		logDestinationType := d.Get("log_destination_type").(string)
		if logDestinationType == "" {
			logDestinationType = ec2.LogDestinationTypeCloudWatchLogs
		}
		req.LogDestinationType = aws.String(logDestinationType)
	} else if v = d.Get("log_group_name").(string); v != "" {
		arn := arn.ARN{
			Partition: meta.(*AWSClient).partition,
			Region:    meta.(*AWSClient).region,
			Service:   "logs",
			AccountID: meta.(*AWSClient).accountid,
			Resource:  fmt.Sprintf("log-group:%s", v),
		}.String()
		req.LogDestination = aws.String(arn)
		req.LogDestinationType = aws.String(ec2.LogDestinationTypeCloudWatchLogs)
	} else {
		return errors.New("Either 'log_destination' or 'log_group_name' must be set")
	}

	if v := d.Get("iam_role_arn").(string); v != "" {
		req.DeliverLogsPermissionArn = aws.String(v)
	} else if aws.StringValue(req.LogDestinationType) == ec2.LogDestinationTypeCloudWatchLogs {
		return errors.New("'iam_role_arn' must be set for CloudWatch Logs destination")
	}

	if v := d.Get("resource_id").(string); v != "" {
		req.ResourceIds = aws.StringSlice([]string{v})
		if resourceType := d.Get("resource_type").(string); resourceType != "" {
			req.ResourceType = aws.String(resourceType)
		} else {
			return errors.New("'resource_type' must be set if 'resource_id' is set")
		}
	} else if v := d.Get("vpc_id").(string); v != "" {
		req.ResourceIds = aws.StringSlice([]string{v})
		req.ResourceType = aws.String(ec2.FlowLogsResourceTypeVpc)
	} else if v := d.Get("subnet_id").(string); v != "" {
		req.ResourceIds = aws.StringSlice([]string{v})
		req.ResourceType = aws.String(ec2.FlowLogsResourceTypeSubnet)
	} else if v := d.Get("eni_id").(string); v != "" {
		req.ResourceIds = aws.StringSlice([]string{v})
		req.ResourceType = aws.String(ec2.FlowLogsResourceTypeNetworkInterface)
	} else {
		return errors.New("One of 'resource_id', 'vpc_id', 'subnet_id' or 'eni_id' must be set")
	}

	log.Printf("[DEBUG] Creating Flow Log: %#v", req)
	resp, err := conn.CreateFlowLogs(req)
	if err != nil {
		return err
	}

	d.SetId(aws.StringValue(resp.FlowLogIds[0]))

	return resourceAwsLogFlowRead(d, meta)
}

func resourceAwsLogFlowRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn

	resp, err := conn.DescribeFlowLogs(&ec2.DescribeFlowLogsInput{
		FlowLogIds: aws.StringSlice([]string{d.Id()}),
	})
	if err != nil {
		return err
	}
	if len(resp.FlowLogs) == 0 {
		log.Printf("[WARN] Flow Log (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	fl := resp.FlowLogs[0]
	d.Set("iam_role_arn", fl.DeliverLogsPermissionArn)
	d.Set("traffic_type", fl.TrafficType)
	d.Set("flow_log_status", fl.FlowLogStatus)
	d.Set("deliver_logs_status", fl.DeliverLogsStatus)

	if logDestination := aws.StringValue(fl.LogDestination); logDestination != "" {
		d.Set("log_destination", logDestination)
		d.Set("log_destination_type", fl.LogDestinationType)
	} else if logGroupName := aws.StringValue(fl.LogGroupName); logGroupName != "" {
		arn := arn.ARN{
			Partition: meta.(*AWSClient).partition,
			Region:    meta.(*AWSClient).region,
			Service:   "logs",
			AccountID: meta.(*AWSClient).accountid,
			Resource:  fmt.Sprintf("log-group:%s", logGroupName),
		}.String()
		d.Set("log_destination", arn)
		d.Set("log_destination_type", ec2.LogDestinationTypeCloudWatchLogs)
	}

	d.Set("resource_id", fl.ResourceId)
	prefix, _ := parseEc2ResourceId(aws.StringValue(fl.ResourceId))
	switch prefix {
	case "vpc":
		d.Set("resource_type", ec2.FlowLogsResourceTypeVpc)
	case "subnet":
		d.Set("resource_type", ec2.FlowLogsResourceTypeSubnet)
	case "eni":
		d.Set("resource_type", ec2.FlowLogsResourceTypeNetworkInterface)
	}

	// Clear legacy attributes.
	d.Set("log_group_name", nil)
	d.Set("vpc_id", nil)
	d.Set("subnet_id", nil)
	d.Set("eni_id", nil)

	return nil
}

func resourceAwsLogFlowDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn

	log.Printf("[DEBUG] Deleting Flow Log: %s", d.Id())
	_, err := conn.DeleteFlowLogs(&ec2.DeleteFlowLogsInput{
		FlowLogIds: aws.StringSlice([]string{d.Id()}),
	})
	if err != nil {
		if isAWSErr(err, "InvalidFlowLogId.NotFound", "") {
			return nil
		}
		return err
	}

	return nil
}

func resourceAwsLogFlowCustomizeDiff(diff *schema.ResourceDiff, meta interface{}) error {
	if diff.Id() != "" {
		// Existing resource.
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
