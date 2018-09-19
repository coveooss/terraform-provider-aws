package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSFlowLog_importBasic(t *testing.T) {
	resourceName := "aws_flow_log.test_flow_log_vpc"

	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFlowLogConfig_vpcOldSyntaxCloudWatch(rInt),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSFlowLog_vpcToCloudWatch(t *testing.T) {
	var flowLog ec2.FlowLog

	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "aws_flow_log.test_flow_log_vpc",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFlowLogConfig_vpcOldSyntaxCloudWatch(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists("aws_flow_log.test_flow_log_vpc", &flowLog),
					resource.TestCheckResourceAttr("aws_flow_log.test_flow_log_vpc", "log_destination_type", "cloud-watch-logs"),
					resource.TestCheckResourceAttr("aws_flow_log.test_flow_log_vpc", "resource_type", "VPC"),
					resource.TestCheckResourceAttr("aws_flow_log.test_flow_log_vpc", "traffic_type", "ALL"),
				),
			},
			{
				Config:             testAccFlowLogConfig_vpcNewSyntaxCloudWatch(rInt),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAWSFlowLog_vpcToCloudWatchStartWithNewSyntax(t *testing.T) {
	var flowLog ec2.FlowLog

	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "aws_flow_log.test_flow_log_vpc",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFlowLogConfig_vpcNewSyntaxCloudWatch(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists("aws_flow_log.test_flow_log_vpc", &flowLog),
					resource.TestCheckResourceAttr("aws_flow_log.test_flow_log_vpc", "log_destination_type", "cloud-watch-logs"),
					resource.TestCheckResourceAttr("aws_flow_log.test_flow_log_vpc", "resource_type", "VPC"),
					resource.TestCheckResourceAttr("aws_flow_log.test_flow_log_vpc", "traffic_type", "ALL"),
				),
			},
			{
				Config:             testAccFlowLogConfig_vpcOldSyntaxCloudWatch(rInt),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAWSFlowLog_subnetToCloudWatch(t *testing.T) {
	var flowLog ec2.FlowLog

	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "aws_flow_log.test_flow_log_subnet",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFlowLogConfig_subnetOldSyntaxCloudWatch(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists("aws_flow_log.test_flow_log_subnet", &flowLog),
					resource.TestCheckResourceAttr("aws_flow_log.test_flow_log_subnet", "log_destination_type", "cloud-watch-logs"),
					resource.TestCheckResourceAttr("aws_flow_log.test_flow_log_subnet", "resource_type", "Subnet"),
					resource.TestCheckResourceAttr("aws_flow_log.test_flow_log_subnet", "traffic_type", "ACCEPT"),
				),
			},
			{
				Config:             testAccFlowLogConfig_subnetNewSyntaxCloudWatch(rInt),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAWSFlowLog_eniToCloudWatch(t *testing.T) {
	var flowLog ec2.FlowLog

	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "aws_flow_log.test_flow_log_eni",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFlowLogConfig_eniOldSyntaxCloudWatch(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists("aws_flow_log.test_flow_log_eni", &flowLog),
					resource.TestCheckResourceAttr("aws_flow_log.test_flow_log_eni", "log_destination_type", "cloud-watch-logs"),
					resource.TestCheckResourceAttr("aws_flow_log.test_flow_log_eni", "resource_type", "NetworkInterface"),
					resource.TestCheckResourceAttr("aws_flow_log.test_flow_log_eni", "traffic_type", "REJECT"),
				),
			},
			{
				Config:             testAccFlowLogConfig_eniNewSyntaxCloudWatch(rInt),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAWSFlowLog_vpcToS3(t *testing.T) {
	var flowLog ec2.FlowLog

	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "aws_flow_log.test_flow_log_vpc",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFlowLogConfig_vpcS3(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowLogExists("aws_flow_log.test_flow_log_vpc", &flowLog),
					resource.TestCheckResourceAttr("aws_flow_log.test_flow_log_vpc", "log_destination_type", "s3"),
					resource.TestCheckResourceAttr("aws_flow_log.test_flow_log_vpc", "resource_type", "VPC"),
					resource.TestCheckResourceAttr("aws_flow_log.test_flow_log_vpc", "traffic_type", "ALL"),
				),
			},
		},
	})
}

func TestAccAWSFlowLog_migrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion int
		ID           string
		Attributes   map[string]string
		Expected     map[string]string
	}{
		"v0_1_vpc_id": {
			StateVersion: 0,
			ID:           "some_id",
			Attributes: map[string]string{
				"vpc_id":         "vpc-12345678",
				"log_group_name": "test",
			},
			Expected: map[string]string{
				"resource_id":          "vpc-12345678",
				"resource_type":        "VPC",
				"log_destination_type": "cloud-watch-logs",
			},
		},
		"v0_1_subnet_id": {
			StateVersion: 0,
			ID:           "some_id",
			Attributes: map[string]string{
				"subnet_id":      "sn-12345678",
				"log_group_name": "test",
			},
			Expected: map[string]string{
				"resource_id":          "sn-12345678",
				"resource_type":        "Subnet",
				"log_destination_type": "cloud-watch-logs",
			},
		},
		"v0_1_eni_id": {
			StateVersion: 0,
			ID:           "some_id",
			Attributes: map[string]string{
				"eni_id":         "eni-12345678",
				"log_group_name": "test",
			},
			Expected: map[string]string{
				"resource_id":          "eni-12345678",
				"resource_type":        "NetworkInterface",
				"log_destination_type": "cloud-watch-logs",
			},
		},
	}

	testAccPreCheck(t)

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         tc.ID,
			Attributes: tc.Attributes,
		}
		is, err := resourceAwsFlowLogMigrateState(tc.StateVersion, is, testAccProvider.Meta())
		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		for k, v := range tc.Expected {
			if is.Attributes[k] != v {
				t.Fatalf("Bad migration (%s): %s\n\n expected: %s", k, is.Attributes[k], v)
			}
		}

		_, err = arn.Parse(is.Attributes["log_destination"])
		if err != nil {
			t.Fatalf("invalid ARN: %s, err: %#v", is.Attributes["log_destination"], err)
		}
	}
}

func testAccCheckFlowLogExists(n string, flowLog *ec2.FlowLog) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Flow Log ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).ec2conn
		describeOpts := &ec2.DescribeFlowLogsInput{
			FlowLogIds: []*string{aws.String(rs.Primary.ID)},
		}
		resp, err := conn.DescribeFlowLogs(describeOpts)
		if err != nil {
			return err
		}
		if len(resp.FlowLogs) == 0 {
			return fmt.Errorf("No Flow Logs found for id (%s)", rs.Primary.ID)
		}

		if *resp.FlowLogs[0].FlowLogStatus != "ACTIVE" {
			return fmt.Errorf("Flow Log status is not ACTIVE, got: %s", *resp.FlowLogs[0].FlowLogStatus)
		}

		*flowLog = *resp.FlowLogs[0]

		return nil
	}
}

func testAccCheckFlowLogDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_flow_log" {
			continue
		}

		return nil
	}

	return nil
}

func testAccFlowLogConfig_vpcOldSyntaxCloudWatch(rInt int) string {
	return fmt.Sprintf(`
resource "aws_vpc" "default" {
  cidr_block = "10.0.0.0/16"

  tags {
    Name = "terraform-testacc-flow-log-vpc"
  }
}

resource "aws_iam_role" "test_role" {
  name = "tf_test_flow_log_vpc_%d"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ec2.amazonaws.com"
        ]
      },
      "Action": [
        "sts:AssumeRole"
      ]
    }
  ]
}
EOF
}

resource "aws_cloudwatch_log_group" "foobar" {
  name = "tf-test-fl-%d"
}

resource "aws_flow_log" "test_flow_log_vpc" {
  log_group_name = "${aws_cloudwatch_log_group.foobar.name}"
  iam_role_arn   = "${aws_iam_role.test_role.arn}"
  vpc_id         = "${aws_vpc.default.id}"
  traffic_type   = "ALL"
}
`, rInt, rInt)
}

func testAccFlowLogConfig_vpcNewSyntaxCloudWatch(rInt int) string {
	return fmt.Sprintf(`
resource "aws_vpc" "default" {
  cidr_block = "10.0.0.0/16"

  tags {
    Name = "terraform-testacc-flow-log-vpc"
  }
}

resource "aws_iam_role" "test_role" {
  name = "tf_test_flow_log_vpc_%d"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ec2.amazonaws.com"
        ]
      },
      "Action": [
        "sts:AssumeRole"
      ]
    }
  ]
}
EOF
}

resource "aws_cloudwatch_log_group" "foobar" {
  name = "tf-test-fl-%d"
}

resource "aws_flow_log" "test_flow_log_vpc" {
  log_destination      = "${aws_cloudwatch_log_group.foobar.arn}"
  log_destination_type = "cloud-watch-logs"
  iam_role_arn         = "${aws_iam_role.test_role.arn}"
  resource_id          = "${aws_vpc.default.id}"
  resource_type        = "VPC"
  traffic_type         = "ALL"
}
`, rInt, rInt)
}

func testAccFlowLogConfig_vpcS3(rInt int) string {
	return fmt.Sprintf(`
resource "aws_vpc" "default" {
  cidr_block = "10.0.0.0/16"

  tags {
    Name = "terraform-testacc-flow-log-vpc"
  }
}

resource "aws_s3_bucket" "foobar" {
  bucket = "tf-test-fl-%d"
}

resource "aws_flow_log" "test_flow_log_vpc" {
  log_destination      = "${aws_s3_bucket.foobar.arn}"
  log_destination_type = "s3"
  resource_id          = "${aws_vpc.default.id}"
  resource_type        = "VPC"
  traffic_type         = "ALL"
}
`, rInt)
}

func testAccFlowLogConfig_subnetOldSyntaxCloudWatch(rInt int) string {
	return fmt.Sprintf(`
resource "aws_vpc" "default" {
  cidr_block = "10.0.0.0/16"

  tags {
    Name = "terraform-testacc-flow-log-subnet"
  }
}

resource "aws_subnet" "test_subnet" {
  vpc_id     = "${aws_vpc.default.id}"
  cidr_block = "10.0.1.0/24"

  tags {
    Name = "tf-acc-flow-log-subnet"
  }
}

resource "aws_iam_role" "test_role" {
  name = "tf_test_flow_log_subnet_%d"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ec2.amazonaws.com"
        ]
      },
      "Action": [
        "sts:AssumeRole"
      ]
    }
  ]
}
EOF
}

resource "aws_cloudwatch_log_group" "foobar" {
  name = "tf-test-fl-%d"
}

resource "aws_flow_log" "test_flow_log_subnet" {
  log_group_name = "${aws_cloudwatch_log_group.foobar.name}"
  iam_role_arn   = "${aws_iam_role.test_role.arn}"
  subnet_id      = "${aws_subnet.test_subnet.id}"
  traffic_type   = "ACCEPT"
}
`, rInt, rInt)
}

func testAccFlowLogConfig_subnetNewSyntaxCloudWatch(rInt int) string {
	return fmt.Sprintf(`
resource "aws_vpc" "default" {
  cidr_block = "10.0.0.0/16"

  tags {
    Name = "terraform-testacc-flow-log-subnet"
  }
}

resource "aws_subnet" "test_subnet" {
  vpc_id     = "${aws_vpc.default.id}"
  cidr_block = "10.0.1.0/24"

  tags {
    Name = "tf-acc-flow-log-subnet"
  }
}

resource "aws_iam_role" "test_role" {
  name = "tf_test_flow_log_subnet_%d"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ec2.amazonaws.com"
        ]
      },
      "Action": [
        "sts:AssumeRole"
      ]
    }
  ]
}
EOF
}

resource "aws_cloudwatch_log_group" "foobar" {
  name = "tf-test-fl-%d"
}

resource "aws_flow_log" "test_flow_log_subnet" {
  log_destination = "${aws_cloudwatch_log_group.foobar.arn}"
  iam_role_arn    = "${aws_iam_role.test_role.arn}"
  resource_id     = "${aws_subnet.test_subnet.id}"
  resource_type   = "Subnet"
  traffic_type    = "ACCEPT"
}
`, rInt, rInt)
}

func testAccFlowLogConfig_eniOldSyntaxCloudWatch(rInt int) string {
	return fmt.Sprintf(`
resource "aws_vpc" "default" {
  cidr_block = "10.0.0.0/16"

  tags {
    Name = "terraform-testacc-flow-log-eni"
  }
}

resource "aws_subnet" "test_subnet" {
  vpc_id     = "${aws_vpc.default.id}"
  cidr_block = "10.0.1.0/24"

  tags {
    Name = "tf-acc-flow-log-eni"
  }
}

resource "aws_network_interface" "test_eni" {
  subnet_id = "${aws_subnet.test_subnet.id}"

  tags {
    Name = "tf-acc-flow-log-eni"
  }
}

resource "aws_iam_role" "test_role" {
  name = "tf_test_flow_log_eni_%d"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ec2.amazonaws.com"
        ]
      },
      "Action": [
        "sts:AssumeRole"
      ]
    }
  ]
}
EOF
}

resource "aws_cloudwatch_log_group" "foobar" {
  name = "tf-test-fl-%d"
}

resource "aws_flow_log" "test_flow_log_eni" {
  log_group_name = "${aws_cloudwatch_log_group.foobar.name}"
  iam_role_arn   = "${aws_iam_role.test_role.arn}"
  eni_id         = "${aws_network_interface.test_eni.id}"
  traffic_type   = "REJECT"
}
`, rInt, rInt)
}

func testAccFlowLogConfig_eniNewSyntaxCloudWatch(rInt int) string {
	return fmt.Sprintf(`
resource "aws_vpc" "default" {
  cidr_block = "10.0.0.0/16"

  tags {
    Name = "terraform-testacc-flow-log-eni"
  }
}

resource "aws_subnet" "test_subnet" {
  vpc_id     = "${aws_vpc.default.id}"
  cidr_block = "10.0.1.0/24"

  tags {
    Name = "tf-acc-flow-log-eni"
  }
}

resource "aws_network_interface" "test_eni" {
  subnet_id = "${aws_subnet.test_subnet.id}"

  tags {
    Name = "tf-acc-flow-log-eni"
  }
}

resource "aws_iam_role" "test_role" {
  name = "tf_test_flow_log_eni_%d"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ec2.amazonaws.com"
        ]
      },
      "Action": [
        "sts:AssumeRole"
      ]
    }
  ]
}
EOF
}

resource "aws_cloudwatch_log_group" "foobar" {
  name = "tf-test-fl-%d"
}

resource "aws_flow_log" "test_flow_log_eni" {
  log_destination      = "${aws_cloudwatch_log_group.foobar.arn}"
  log_destination_type = "cloud-watch-logs"
  iam_role_arn         = "${aws_iam_role.test_role.arn}"
  resource_id          = "${aws_network_interface.test_eni.id}"
  resource_type        = "NetworkInterface"
  traffic_type         = "REJECT"
}
`, rInt, rInt)
}
