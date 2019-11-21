/*
Provides a resource to create a CAM user policy attachment.

Example Usage

```hcl
resource "tencentcloud_cam_user_policy_attachment" "foo" {
  user_id   = "cam-test"
  policy_id = "26800353"
}
```

Import

CAM user policy attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_user_policy_attachment.foo cam-test#26800353
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
)

func resourceTencentCloudCamUserPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamUserPolicyAttachmentCreate,
		Read:   resourceTencentCloudCamUserPolicyAttachmentRead,
		Delete: resourceTencentCloudCamUserPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the attached cam user.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the policy.",
			},
			"create_mode": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Mode of Creation of the CAM user policy attachment. 1 means the cam policy attachment is created by production, and the others indicate syntax strategy ways.",
			},
			"policy_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the policy strategy. 'User' means customer strategy and 'QCS' means preset strategy.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the CAM user policy attachment.",
			},
			"policy_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the policy.",
			},
		},
	}
}

func resourceTencentCloudCamUserPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user_policy_attachment.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	userId := d.Get("user_id").(string)
	policyId := d.Get("policy_id").(string)
	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := camService.AddUserPolicyAttachment(ctx, userId, policyId)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete cam user policy attachment failed, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetId(userId + "#" + policyId)

	return resourceTencentCloudCamUserPolicyAttachmentRead(d, meta)
}

func resourceTencentCloudCamUserPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user_policy_attachment.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	userPolicyAttachmentId := d.Id()

	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var instance *cam.AttachPolicyInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := camService.DescribeUserPolicyAttachmentById(ctx, userPolicyAttachmentId)
		if e != nil {
			return retryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read cam user policy attachment failed, reason:%s\n ", logId, err.Error())
		return err
	}
	//split id
	userId, policyId, e := camService.decodeCamPolicyAttachmentId(userPolicyAttachmentId)
	if e != nil {
		return e
	}
	d.Set("user_id", userId)
	d.Set("policy_id", strconv.Itoa(int(policyId)))
	d.Set("policy_name", *instance.PolicyName)
	d.Set("create_time", *instance.AddTime)
	d.Set("create_mode", int(*instance.CreateMode))
	d.Set("policy_type", *instance.PolicyType)
	return nil
}

func resourceTencentCloudCamUserPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user_policy_attachment.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	userPolicyAttachmentId := d.Id()

	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := camService.DeleteUserPolicyAttachmentById(ctx, userPolicyAttachmentId)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete cam user policy attachment failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
