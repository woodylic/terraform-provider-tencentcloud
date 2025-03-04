/*
Use this data source to query detailed information of CAM policies

Example Usage

```hcl
data "tencentcloud_cam_policies" "foo" {
  policy_id   = "26655801"
  name        = "cam-policy-test"
  type        = 1
  create_mode = 1
  description = "test"
}
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

func dataSourceTencentCloudCamPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamPoliciesRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the CAM policy to be queried.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of CAM policy to be queried to be queried.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the CAM policy.",
			},
			"type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{1, 2}),
				Description:  "Type of the policy strategy. 1 means customer strategy and 2 means preset strategy.",
			},
			"create_mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{1, 2}),
				Description:  "Mode of creation of policy strategy. 1 means policy was created with console, and 2 means it was created by strategies.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"policy_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of CAM policies. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of CAM policy.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of CAM policy.",
						},
						"attachments": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of attached users.",
						},
						"service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of attached products.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the CAM policy.",
						},
						"type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Type of the policy strategy. 1 means customer strategy and 2 means preset strategy.",
						},
						"create_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Mode of creation of policy strategy. 1 means policy was created with console, and 2 means it was created by strategies.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCamPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cam_policies.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	params := make(map[string]interface{})
	if v, ok := d.GetOk("policy_id"); ok {
		params["policy_id"], _ = strconv.Atoi(v.(string))
	}
	if v, ok := d.GetOk("name"); ok {
		params["name"] = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		params["description"] = v.(string)
	}
	if v, ok := d.GetOk("create_mode"); ok {
		params["create_mode"] = v.(int)
	}
	if v, ok := d.GetOk("type"); ok {
		params["type"] = v.(int)
	}

	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var policies []*cam.StrategyInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := camService.DescribePoliciesByFilter(ctx, params)
		if e != nil {
			return retryError(e)
		}
		policies = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM policies failed, reason:%s\n", logId, err.Error())
		return err
	}
	policyList := make([]map[string]interface{}, 0, len(policies))
	ids := make([]string, 0, len(policies))
	for _, policy := range policies {
		mapping := map[string]interface{}{
			"name":         *policy.PolicyName,
			"attachments":  int(*policy.Attachments),
			"description":  *policy.Description,
			"create_time":  *policy.AddTime,
			"service_type": *policy.ServiceType,
			"create_mode":  int(*policy.CreateMode),
			"type":         int(*policy.Type),
		}
		policyList = append(policyList, mapping)
		ids = append(ids, strconv.Itoa(int(*policy.PolicyId)))
	}

	d.SetId(dataResourceIdsHash(ids))
	if e := d.Set("policy_list", policyList); e != nil {
		log.Printf("[CRITAL]%s provider set policy list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), policyList); e != nil {
			return e
		}
	}

	return nil
}
