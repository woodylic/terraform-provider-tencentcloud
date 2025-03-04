package tencentcloud

/*
---------------------------------------------------
If you want to run through the test cases,
the following must be changed to your resource id.
---------------------------------------------------
*/

const appid string = "1259649581"

const defaultRegion = "ap-guangzhou"
const defaultVpcId = "vpc-h70b6b49"
const defaultVpcCidr = "172.16.0.0/16"
const defaultVpcCidrLess = "172.16.0.0/18"

const defaultAZone = "ap-guangzhou-3"
const defaultSubnetId = "subnet-1uwh63so"
const defaultSubnetCidr = "172.16.0.0/20"
const defaultSubnetCidrLess = "172.16.0.0/22"

const defaultInsName = "tf-ci-test"
const defaultInsNameUpdate = "tf-ci-test-update"

/*
---------------------------------------------------
The folling are common test case used as templates.
---------------------------------------------------
*/

const defaultVpcVariable = `
variable "instance_name" {
  default = "` + defaultInsName + `"
}

variable "instance_name_update" {
  default = "` + defaultInsNameUpdate + `"
}

variable "availability_zone" {
  default = "` + defaultAZone + `"
}

variable "vpc_id" {
  default = "` + defaultVpcId + `"
}

variable "vpc_cidr" {
  default = "` + defaultVpcCidr + `"
}

variable "vpc_cidr_less" {
  default = "` + defaultVpcCidrLess + `"
}

variable "subnet_id" {
  default = "` + defaultSubnetId + `"
}

variable "subnet_cidr" {
  default = "` + defaultSubnetCidr + `"
}

variable "subnet_cidr_less" {
  default = "` + defaultSubnetCidrLess + `"
}
`

const defaultInstanceVariable = defaultVpcVariable + `
data "tencentcloud_availability_zones" "default" {
}

data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

data "tencentcloud_instance_types" "default" {
  filter {
    name   = "instance-family"
    values = ["S1"]
  }

  cpu_core_count = 1
  memory_size    = 1
}
`

const instanceCommonTestCase = defaultInstanceVariable + `
resource "tencentcloud_instance" "default" {
  instance_name              = "${var.instance_name}"
  availability_zone          = "${data.tencentcloud_availability_zones.default.zones.0.name}"
  image_id                   = "${data.tencentcloud_images.default.images.0.image_id}"
  instance_type              = "${data.tencentcloud_instance_types.default.instance_types.0.instance_type}"
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  allocate_public_ip         = true
  internet_max_bandwidth_out = 10
  vpc_id                     = "${var.vpc_id}"
  subnet_id                  = "${var.subnet_id}"
}
`

const mysqlInstanceCommonTestCase = defaultVpcVariable + `
resource "tencentcloud_mysql_instance" "default" {
  mem_size = 1000
  volume_size = 25
  instance_name = "${var.instance_name}"
  engine_version = "5.7"
  root_password = "0153Y474"
  availability_zone = "${var.availability_zone}"
}
`
