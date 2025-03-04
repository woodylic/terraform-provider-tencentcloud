package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudCamGroupMembership_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamGroupMembershipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamGroupMembership_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamGroupMembershipExists("tencentcloud_cam_group_membership.group_membership_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_group_membership.group_membership_basic", "group_id"),
					resource.TestCheckResourceAttr("tencentcloud_cam_group_membership.group_membership_basic", "user_ids.#", "1"),
				),
			}, {
				Config: testAccCamGroupMembership_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamGroupMembershipExists("tencentcloud_cam_group_membership.group_membership_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_group_membership.group_membership_basic", "group_id"),
					resource.TestCheckResourceAttr("tencentcloud_cam_group_membership.group_membership_basic", "user_ids.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_cam_group_membership.group_membership_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCamGroupMembershipDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	camService := CamService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cam_group_membership" {
			continue
		}

		_, err := camService.DescribeGroupMembershipById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("CAM group membership still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCamGroupMembershipExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("CAM group %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("CAM group id is not set")
		}
		camService := CamService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, err := camService.DescribeGroupMembershipById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

const testAccCamGroupMembership_basic = `
resource "tencentcloud_cam_group" "group_basic" {
  name   = "cam-group-membership-test"
  remark = "test"
}

resource "tencentcloud_cam_user" "foo" {
  name                = "cam-user-test2"
  remark              = "test"
  console_login       = true
  use_api             = true
  need_reset_password = true
  password            = "Gail@1234"
  phone_num           = "13631555963"
  country_code        = "86"
  email               = "1234@qq.com"
}

resource "tencentcloud_cam_group_membership" "group_membership_basic" {
  group_id = "${tencentcloud_cam_group.group_basic.id}"
  user_ids = ["${tencentcloud_cam_user.foo.id}",]
}
`

const testAccCamGroupMembership_update = `
resource "tencentcloud_cam_group" "group_basic" {
  name   = "cam-group-membership-test"
  remark = "test"
}

resource "tencentcloud_cam_user" "user_basic" {
  name                = "cam-user-testj"
  remark              = "test"
  console_login       = true
  use_api             = true
  need_reset_password = true
  password            = "Gail@1234"
  phone_num           = "13631555963"
  country_code        = "86"
  email               = "1234@qq.com"
}

resource "tencentcloud_cam_group_membership" "group_membership_basic" {
  group_id = "${tencentcloud_cam_group.group_basic.id}"
  user_ids = ["${tencentcloud_cam_user.user_basic.id}"]
}
`
