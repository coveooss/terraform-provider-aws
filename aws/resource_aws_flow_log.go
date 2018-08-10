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
		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
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

			"log_destination_type":{
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
			},

			"resource_id": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				MinItems: 1,
				MaxItems: 1000,
				Required: true,
				ForceNew: true,
			},

			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					ec2.ResourceTypeVpc,
					ec2.ResourceTypeSubnet,
					ec2.ResourceTypeNetworkInterface,
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
		},
	}
}

func resourceAwsLogFlowCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn

	opts := &ec2.CreateFlowLogsInput{
		DeliverLogsPermissionArn: aws.String(d.Get("iam_role_arn").(string)),
		LogDestination:           aws.String(d.Get("log_destination").(string)),
		LogDestinationType:       aws.String(d.Get("log_destination_type").(string)),
		LogGroupName:             aws.String(d.Get("log_group_name").(string)),
		ResourceIds:              aws.StringSlice(d.Get("resource_id").([]string)),
		ResourceType:             aws.String(d.Get("resource_type").(string)),
		TrafficType:              aws.String(d.Get("traffic_type").(string)),
	}

	log.Printf(
		"[DEBUG] Flow Log Create configuration: %s", opts)
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
	d.Set("resource_id", fl.ResourceId)

	if strings.HasPrefix(*fl.ResourceId, "vpc-") {
		d.Set("resource_type", ec2.ResourceTypeVpc)
	} else if strings.HasPrefix(*fl.ResourceId, "subnet-") {
		d.Set("resource_type", ec2.ResourceTypeSubnet)
	} else if strings.HasPrefix(*fl.ResourceId, "eni-") {
		d.Set("resource_type", ec2.ResourceTypeNetworkInterface)
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
