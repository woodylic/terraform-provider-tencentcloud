package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudClbServerAttachment_tcp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbServerAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbServerAttachment_tcp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbServerAttachmentExists("tencentcloud_clb_attachment.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_attachment.foo", "clb_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_attachment.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.foo", "protocol_type", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.foo", "targets.#", "1"),
				),
			}, {
				Config: testAccClbServerAttachment_tcp_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbServerAttachmentExists("tencentcloud_clb_attachment.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_attachment.foo", "clb_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_attachment.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.foo", "protocol_type", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.foo", "targets.#", "1"),
				),
			},
		},
	})
}

func TestAccTencentCloudClbServerAttachment_http(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbServerAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbServerAttachment_http,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbServerAttachmentExists("tencentcloud_clb_attachment.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_attachment.foo", "clb_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_attachment.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.foo", "protocol_type", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.foo", "targets.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_attachment.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClbServerAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	clbService := ClbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_clb_attachment" {
			continue
		}
		time.Sleep(5 * time.Second)
		items := strings.Split(rs.Primary.ID, "#")
		if len(items) != 3 {
			return fmt.Errorf("id of resource.tencentcloud_clb_attachment is wrong")
		}
		locationId := items[0]
		listenerId := items[1]
		clbId := items[2]
		_, err := clbService.DescribeAttachmentByPara(ctx, clbId, listenerId, locationId)
		if err == nil {
			return fmt.Errorf("clb ServerAttachment still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClbServerAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("clb ServerAttachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("clb ServerAttachment id is not set")
		}
		clbService := ClbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		items := strings.Split(rs.Primary.ID, "#")
		if len(items) != 3 {
			return fmt.Errorf("id of resource.tencentcloud_clb_attachment is wrong")
		}
		locationId := items[0]
		listenerId := items[1]
		clbId := items[2]
		_, err := clbService.DescribeAttachmentByPara(ctx, clbId, listenerId, locationId)
		if err != nil {
			return err
		}
		return nil
	}
}

const testAccClbServerAttachment_tcp = instanceCommonTestCase + `
resource "tencentcloud_clb_instance" "foo" {
  network_type = "OPEN"
  clb_name     = "${var.instance_name}"
  vpc_id       = "${var.vpc_id}"
}

resource "tencentcloud_clb_listener" "foo" {
  clb_id                     = "${tencentcloud_clb_instance.foo.id}"
  listener_name              = "${var.instance_name}"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
}

resource "tencentcloud_clb_attachment" "foo" {
  clb_id      = "${tencentcloud_clb_instance.foo.id}"
  listener_id = "${tencentcloud_clb_listener.foo.id}"

  targets {
    instance_id = "${tencentcloud_instance.default.id}"
    port        = 23
    weight      = 10
  }
}
`

const testAccClbServerAttachment_tcp_update = instanceCommonTestCase + `
resource "tencentcloud_clb_instance" "foo" {
  network_type = "OPEN"
  clb_name     = "${var.instance_name}"
  vpc_id       = "${var.vpc_id}"
}

resource "tencentcloud_clb_listener" "foo" {
  clb_id                     = "${tencentcloud_clb_instance.foo.id}"
  listener_name              = "${var.instance_name}"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
}

resource "tencentcloud_clb_attachment" "foo" {
  clb_id      = "${tencentcloud_clb_instance.foo.id}"
  listener_id = "${tencentcloud_clb_listener.foo.id}"

  targets {
    instance_id = "${tencentcloud_instance.default.id}"
    port        = 23
    weight      = 50
  }
}
`

const testAccClbServerAttachment_http = instanceCommonTestCase + `
resource "tencentcloud_clb_instance" "foo" {
  network_type = "OPEN"
  clb_name     = "${var.instance_name}"
  vpc_id       = "${var.vpc_id}"
}

resource "tencentcloud_clb_listener" "foo" {
  clb_id               = "${tencentcloud_clb_instance.foo.id}"
  listener_name        = "${var.instance_name}"
  port                 = 77
  protocol             = "HTTPS"
  certificate_ssl_mode = "UNIDIRECTIONAL"
  certificate_id       = "VjANRdz8"
}

resource "tencentcloud_clb_listener_rule" "foo" {
  clb_id              = "${tencentcloud_clb_instance.foo.id}"
  listener_id         = "${tencentcloud_clb_listener.foo.id}"
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}

resource "tencentcloud_clb_attachment" "foo" {
  clb_id      = "${tencentcloud_clb_instance.foo.id}"
  listener_id = "${tencentcloud_clb_listener.foo.id}"
  rule_id     = "${tencentcloud_clb_listener_rule.foo.id}"

  targets {
    instance_id = "${tencentcloud_instance.default.id}"
    port        = 23
    weight      = 10
  }
}
`
