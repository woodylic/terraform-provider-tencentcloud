package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudCcnV3AttachmentBasic(t *testing.T) {
	keyName := "tencentcloud_ccn_attachment.attachment"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCcnAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCcnAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCcnAttachmentExists(keyName),
					resource.TestCheckResourceAttrSet(keyName, "ccn_id"),
					resource.TestCheckResourceAttrSet(keyName, "instance_type"),
					resource.TestCheckResourceAttrSet(keyName, "instance_region"),
					resource.TestCheckResourceAttrSet(keyName, "instance_id"),
					resource.TestCheckResourceAttrSet(keyName, "state"),
					resource.TestCheckResourceAttrSet(keyName, "attached_time"),
					resource.TestCheckResourceAttrSet(keyName, "cidr_block.#"),
				),
			},
		},
	})
}

func testAccCheckCcnAttachmentExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeCcnAttachedInstance(ctx,
			rs.Primary.Attributes["ccn_id"],
			rs.Primary.Attributes["instance_region"],
			rs.Primary.Attributes["instance_type"],
			rs.Primary.Attributes["instance_id"])

		if err != nil {
			return err
		}
		if has > 0 {
			return nil
		}
		return fmt.Errorf("ccn attachment not exists.")
	}
}

func testAccCheckCcnAttachmentDestroy(s *terraform.State) error {

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ccn_attachment" {
			continue
		}
		time.Sleep(5 * time.Second)
		_, has, err := service.DescribeCcnAttachedInstance(ctx,
			rs.Primary.Attributes["ccn_id"], rs.Primary.Attributes["instance_region"],
			rs.Primary.Attributes["instance_type"],
			rs.Primary.Attributes["instance_id"])
		if err != nil {
			return err
		}
		if has == 0 {
			return nil
		}
		return fmt.Errorf("ccn  attachment not delete ok")
	}
	return nil
}

const testAccCcnAttachmentConfig = `
variable "region" {
  default = "ap-guangzhou"
}

resource tencentcloud_vpc vpc {
  name         = "ci-temp-test-vpc"
  cidr_block   = "10.0.0.0/16"
  dns_servers  = ["119.29.29.29", "8.8.8.8"]
  is_multicast = false
}

resource tencentcloud_ccn main {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource tencentcloud_ccn_attachment attachment {
  ccn_id          = "${tencentcloud_ccn.main.id}"
  instance_type   = "VPC"
  instance_id     = "${tencentcloud_vpc.vpc.id}"
  instance_region = "${var.region}"
}
`
