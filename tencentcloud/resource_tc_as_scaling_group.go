/*
Provides a resource to create a group of AS (Auto scaling) instances.

Example Usage

```hcl
resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name   = "tf-as-scaling-group"
  configuration_id     = "asc-oqio4yyj"
  max_size             = 1
  min_size             = 0
  vpc_id               = "vpc-3efmz0z"
  subnet_ids           = ["subnet-mc3egos"]
  project_id           = 0
  default_cooldown     = 400
  desired_capacity     = 1
  termination_policies = ["NEWEST_INSTANCE"]
  retry_policy         = "INCREMENTAL_INTERVALS"

  forward_balancer_ids {
    load_balancer_id = "lb-hk693b1l"
    listener_id      = "lbl-81wr497k"
    rule_id          = "loc-kiodx943"

    target_attribute {
      port   = 80
      weight = 90
    }
  }
}
```

Import

AutoScaling Groups can be imported using the id, e.g.

```
$ terraform import tencentcloud_as_scaling_group.scaling_group asg-n32ymck2
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudAsScalingGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsScalingGroupCreate,
		Read:   resourceTencentCloudAsScalingGroupRead,
		Update: resourceTencentCloudAsScalingGroupUpdate,
		Delete: resourceTencentCloudAsScalingGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"scaling_group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 55),
				Description:  "Name of a scaling group.",
			},
			"configuration_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "An available ID for a launch configuration.",
			},
			"max_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(0, 2000),
				Description:  "Maximum number of CVM instances (0~2000).",
			},
			"min_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(0, 2000),
				Description:  "Minimum number of CVM instances (0~2000).",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of VPC network.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Specifys to which project the scaling group belongs.",
			},
			"subnet_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "ID list of subnet, and for VPC it is required.",
			},
			"zones": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of available zones, for Basic network it is required.",
			},
			"default_cooldown": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Default cooldown time in second, and default value is 300.",
			},
			"desired_capacity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Desired volume of CVM instances, which is between max_size and min_size.",
			},
			"load_balancer_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"forward_balancer_ids"},
				Description:   "ID list of traditional load balancers.",
			},
			"forward_balancer_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"load_balancer_ids"},
				Description:   "List of application load balancers, which can't be specified with load_balancer_ids together.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of available load balancers.",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Listener ID for application load balancers.",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of forwarding rules.",
						},
						"target_attribute": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Attribute list of target rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Port number.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Weight.",
									},
								},
							},
						},
					},
				},
			},
			"termination_policies": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Available values for termination policies include OLDEST_INSTANCE and NEWEST_INSTANCE.",
				Elem: &schema.Schema{
					Type:    schema.TypeString,
					Default: SCALING_GROUP_TERMINATION_POLICY_OLDEST_INSTANCE,
					ValidateFunc: validateAllowedStringValue([]string{SCALING_GROUP_TERMINATION_POLICY_OLDEST_INSTANCE,
						SCALING_GROUP_TERMINATION_POLICY_NEWEST_INSTANCE}),
				},
			},
			"retry_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Available values for retry policies include IMMEDIATE_RETRY and INCREMENTAL_INTERVALS.",
				Default:     SCALING_GROUP_RETRY_POLICY_IMMEDIATE_RETRY,
				ValidateFunc: validateAllowedStringValue([]string{SCALING_GROUP_RETRY_POLICY_IMMEDIATE_RETRY,
					SCALING_GROUP_RETRY_POLICY_INCREMENTAL_INTERVALS}),
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of a scaling group.",
			},

			// computed value
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current status of a scaling group.",
			},
			"instance_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Instance number of a scaling group.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the AS group was created.",
			},
		},
	}
}

func resourceTencentCloudAsScalingGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_scaling_group.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	request := as.NewCreateAutoScalingGroupRequest()

	request.AutoScalingGroupName = stringToPointer(d.Get("scaling_group_name").(string))
	request.LaunchConfigurationId = stringToPointer(d.Get("configuration_id").(string))
	request.MaxSize = intToPointer(d.Get("max_size").(int))
	request.MinSize = intToPointer(d.Get("min_size").(int))
	request.VpcId = stringToPointer(d.Get("vpc_id").(string))
	if v, ok := d.GetOk("default_cooldown"); ok {
		request.DefaultCooldown = intToPointer(v.(int))
	}
	if v, ok := d.GetOk("desired_capacity"); ok {
		request.DesiredCapacity = intToPointer(v.(int))
	}
	if v, ok := d.GetOk("retry_policy"); ok {
		request.RetryPolicy = stringToPointer(v.(string))
	}

	if v, ok := d.GetOk("subnet_ids"); ok {
		subnetIds := v.([]interface{})
		request.SubnetIds = make([]*string, 0, len(subnetIds))
		for i := range subnetIds {
			subnetId := subnetIds[i].(string)
			request.SubnetIds = append(request.SubnetIds, &subnetId)
		}
	}

	if v, ok := d.GetOk("zones"); ok {
		zones := v.([]interface{})
		request.Zones = make([]*string, 0, len(zones))
		for i := range zones {
			zone := zones[i].(string)
			request.Zones = append(request.Zones, &zone)
		}
	}

	if v, ok := d.GetOk("load_balancer_ids"); ok {
		loadBalancerIds := v.([]interface{})
		request.LoadBalancerIds = make([]*string, 0, len(loadBalancerIds))
		for i := range loadBalancerIds {
			loadBalancerId := loadBalancerIds[i].(string)
			request.LoadBalancerIds = append(request.LoadBalancerIds, &loadBalancerId)
		}
	}

	if v, ok := d.GetOk("forward_balancer_ids"); ok {
		forwardBalancers := v.([]interface{})
		request.ForwardLoadBalancers = make([]*as.ForwardLoadBalancer, 0, len(forwardBalancers))
		for _, v := range forwardBalancers {
			vv := v.(map[string]interface{})
			targets := vv["target_attribute"].([]interface{})
			forwardBalancer := as.ForwardLoadBalancer{
				LoadBalancerId: stringToPointer(vv["load_balancer_id"].(string)),
				ListenerId:     stringToPointer(vv["listener_id"].(string)),
				LocationId:     stringToPointer(vv["rule_id"].(string)),
			}
			forwardBalancer.TargetAttributes = make([]*as.TargetAttribute, 0, len(targets))
			for _, target := range targets {
				t := target.(map[string]interface{})
				targetAttribute := as.TargetAttribute{
					Port:   intToPointer(t["port"].(int)),
					Weight: intToPointer(t["weight"].(int)),
				}
				forwardBalancer.TargetAttributes = append(forwardBalancer.TargetAttributes, &targetAttribute)
			}

			request.ForwardLoadBalancers = append(request.ForwardLoadBalancers, &forwardBalancer)
		}
	}

	if v, ok := d.GetOk("termination_policies"); ok {
		terminationPolicies := v.([]interface{})
		request.TerminationPolicies = make([]*string, 0, len(terminationPolicies))
		for i := range terminationPolicies {
			terminationPolicy := terminationPolicies[i].(string)
			request.TerminationPolicies = append(request.TerminationPolicies, &terminationPolicy)
		}
	}

	if tags := getTags(d, "tags"); len(tags) > 0 {
		for k, v := range tags {
			request.Tags = append(request.Tags, &as.Tag{
				Key:   stringToPointer(k),
				Value: stringToPointer(v),
			})
		}
	}

	var id string
	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := meta.(*TencentCloudClient).apiV3Conn.UseAsClient().CreateAutoScalingGroup(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response.Response.AutoScalingGroupId == nil {
			err = fmt.Errorf("Auto scaling group id is nil")
			return resource.NonRetryableError(err)
		}

		id = *response.Response.AutoScalingGroupId

		return nil
	}); err != nil {
		return err
	}
	d.SetId(id)

	// wait for status
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		scalingGroup, _, errRet := asService.DescribeAutoScalingGroupById(ctx, id)
		if errRet != nil {
			return retryError(errRet, "InternalError")
		}
		if scalingGroup != nil && *scalingGroup.InActivityStatus == SCALING_GROUP_NOT_IN_ACTIVITY_STATUS {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("scaling group status is %s, retry...", *scalingGroup.InActivityStatus))
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudAsScalingGroupRead(d, meta)
}

func resourceTencentCloudAsScalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_scaling_group.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	scalingGroupId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var (
		scalingGroup *as.AutoScalingGroup
		e            error
		has          int
	)
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		scalingGroup, has, e = asService.DescribeAutoScalingGroupById(ctx, scalingGroupId)
		if e != nil {
			return retryError(e)
		}
		if has == 0 {
			d.SetId("")
		}
		return nil
	})
	if err != nil {
		return err
	}
	if has == 0 {
		return nil
	}

	d.Set("scaling_group_name", scalingGroup.AutoScalingGroupName)
	d.Set("configuration_id", scalingGroup.LaunchConfigurationId)
	d.Set("status", scalingGroup.AutoScalingGroupStatus)
	d.Set("instance_count", scalingGroup.InstanceCount)
	d.Set("max_size", scalingGroup.MaxSize)
	d.Set("min_size", scalingGroup.MinSize)
	d.Set("vpc_id", scalingGroup.VpcId)
	d.Set("project_id", scalingGroup.ProjectId)
	d.Set("subnet_ids", flattenStringList(scalingGroup.SubnetIdSet))
	d.Set("zones", flattenStringList(scalingGroup.ZoneSet))
	d.Set("default_cooldown", scalingGroup.DefaultCooldown)
	d.Set("desired_capacity", scalingGroup.DesiredCapacity)
	d.Set("load_balancer_ids", flattenStringList(scalingGroup.LoadBalancerIdSet))
	d.Set("termination_policies", flattenStringList(scalingGroup.TerminationPolicySet))
	d.Set("retry_policy", scalingGroup.RetryPolicy)
	d.Set("create_time", scalingGroup.CreatedTime)

	if scalingGroup.ForwardLoadBalancerSet != nil && len(scalingGroup.ForwardLoadBalancerSet) > 0 {
		forwardLoadBalancers := make([]map[string]interface{}, 0, len(scalingGroup.ForwardLoadBalancerSet))
		for _, v := range scalingGroup.ForwardLoadBalancerSet {
			targetAttributes := make([]map[string]interface{}, 0, len(v.TargetAttributes))
			for _, vv := range v.TargetAttributes {
				targetAttribute := map[string]interface{}{
					"port":   *vv.Port,
					"weight": *vv.Weight,
				}
				targetAttributes = append(targetAttributes, targetAttribute)
			}
			forwardLoadBalancer := map[string]interface{}{
				"load_balancer_id":  *v.LoadBalancerId,
				"listener_id":       *v.ListenerId,
				"target_attributes": targetAttributes,
				"rule_id":           *v.LocationId,
			}
			forwardLoadBalancers = append(forwardLoadBalancers, forwardLoadBalancer)
		}
		d.Set("forward_balancer_ids", forwardLoadBalancers)
	}

	tags := make(map[string]string, len(scalingGroup.Tags))
	for _, tag := range scalingGroup.Tags {
		tags[*tag.Key] = *tag.Value
	}
	d.Set("tags", tags)

	return nil
}

func resourceTencentCloudAsScalingGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_scaling_group.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	client := meta.(*TencentCloudClient).apiV3Conn
	tagService := TagService{client: client}
	region := client.Region

	request := as.NewModifyAutoScalingGroupRequest()
	scalingGroupId := d.Id()

	d.Partial(true)

	var updateAttrs []string

	request.AutoScalingGroupId = &scalingGroupId
	if d.HasChange("scaling_group_name") {
		updateAttrs = append(updateAttrs, "scaling_group_name")
		request.AutoScalingGroupName = stringToPointer(d.Get("scaling_group_name").(string))
	}
	if d.HasChange("max_size") {
		updateAttrs = append(updateAttrs, "max_size")
		request.MaxSize = intToPointer(d.Get("max_size").(int))
	}
	if d.HasChange("min_size") {
		updateAttrs = append(updateAttrs, "max_size")
		request.MinSize = intToPointer(d.Get("min_size").(int))
	}
	if d.HasChange("vpc_id") {
		updateAttrs = append(updateAttrs, "vpc_id")
		request.VpcId = stringToPointer(d.Get("vpc_id").(string))
	}
	if d.HasChange("project_id") {
		updateAttrs = append(updateAttrs, "project_id")
		request.ProjectId = intToPointer(d.Get("project_id").(int))
	}
	if d.HasChange("default_cooldown") {
		updateAttrs = append(updateAttrs, "default_cooldown")
		request.DefaultCooldown = intToPointer(d.Get("default_cooldown").(int))
	}
	if d.HasChange("desired_capacity") {
		updateAttrs = append(updateAttrs, "desired_capacity")
		request.DesiredCapacity = intToPointer(d.Get("desired_capacity").(int))
	}
	if d.HasChange("retry_policy") {
		updateAttrs = append(updateAttrs, "retry_policy")
		request.RetryPolicy = stringToPointer(d.Get("retry_policy").(string))
	}
	if d.HasChange("subnet_ids") {
		updateAttrs = append(updateAttrs, "subnet_ids")
		subnetIds := d.Get("subnet_ids").([]interface{})
		request.SubnetIds = make([]*string, 0, len(subnetIds))
		for i := range subnetIds {
			subnetId := subnetIds[i].(string)
			request.SubnetIds = append(request.SubnetIds, &subnetId)
		}
	}
	if d.HasChange("zones") {
		updateAttrs = append(updateAttrs, "zones")
		zones := d.Get("zones").([]interface{})
		request.Zones = make([]*string, 0, len(zones))
		for i := range zones {
			zone := zones[i].(string)
			request.Zones = append(request.Zones, &zone)
		}
	}
	if d.HasChange("termination_policies") {
		updateAttrs = append(updateAttrs, "termination_policies")
		terminationPolicies := d.Get("termination_policies").([]interface{})
		request.TerminationPolicies = make([]*string, 0, len(terminationPolicies))
		for i := range terminationPolicies {
			terminationPolicy := terminationPolicies[i].(string)
			request.TerminationPolicies = append(request.TerminationPolicies, &terminationPolicy)
		}
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := client.UseAsClient().ModifyAutoScalingGroup(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		return nil
	}); err != nil {
		return err
	}

	for _, attr := range updateAttrs {
		d.SetPartial(attr)
	}
	updateAttrs = updateAttrs[:0]

	balancerRequest := as.NewModifyLoadBalancersRequest()
	balancerRequest.AutoScalingGroupId = &scalingGroupId
	if d.HasChange("load_balancer_ids") {
		updateAttrs = append(updateAttrs, "load_balancer_ids")

		loadBalancerIds := d.Get("load_balancer_ids").([]interface{})
		balancerRequest.LoadBalancerIds = make([]*string, 0, len(loadBalancerIds))
		for i := range loadBalancerIds {
			loadBalancerId := loadBalancerIds[i].(string)
			balancerRequest.LoadBalancerIds = append(balancerRequest.LoadBalancerIds, &loadBalancerId)
		}
	}

	if d.HasChange("forward_balancer_ids") {
		updateAttrs = append(updateAttrs, "forward_balancer_ids")

		forwardBalancers := d.Get("forward_balancer_ids").([]interface{})
		balancerRequest.ForwardLoadBalancers = make([]*as.ForwardLoadBalancer, 0, len(forwardBalancers))
		for _, v := range forwardBalancers {
			vv := v.(map[string]interface{})
			targets := vv["target_attribute"].([]interface{})
			forwardBalancer := as.ForwardLoadBalancer{
				LoadBalancerId: stringToPointer(vv["load_balancer_id"].(string)),
				ListenerId:     stringToPointer(vv["listener_id"].(string)),
				LocationId:     stringToPointer(vv["rule_id"].(string)),
			}
			forwardBalancer.TargetAttributes = make([]*as.TargetAttribute, 0, len(targets))
			for _, target := range targets {
				t := target.(map[string]interface{})
				targetAttribute := as.TargetAttribute{
					Port:   intToPointer(t["port"].(int)),
					Weight: intToPointer(t["weight"].(int)),
				}
				forwardBalancer.TargetAttributes = append(forwardBalancer.TargetAttributes, &targetAttribute)
			}

			balancerRequest.ForwardLoadBalancers = append(balancerRequest.ForwardLoadBalancers, &forwardBalancer)
		}
	}

	if len(updateAttrs) > 0 {
		if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(balancerRequest.GetAction())

			balancerResponse, err := client.UseAsClient().ModifyLoadBalancers(balancerRequest)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, balancerRequest.GetAction(), balancerRequest.ToJsonString(), err.Error())
				return retryError(err)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, balancerRequest.GetAction(), balancerRequest.ToJsonString(), balancerResponse.ToJsonString())

			return nil
		}); err != nil {
			return err
		}

		for _, attr := range updateAttrs {
			d.SetPartial(attr)
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := BuildTagResourceName("as", "auto-scaling-group", region, scalingGroupId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceTencentCloudAsScalingGroupRead(d, meta)
}

func resourceTencentCloudAsScalingGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_scaling_group.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	scalingGroupId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	// We need read the scaling group in order to check if there are instances.
	// If so, we need to remove those first.
	scalingGroup, has, err := asService.DescribeAutoScalingGroupById(ctx, scalingGroupId)
	if err != nil {
		return err
	}
	if has == 0 {
		return nil
	}
	if *scalingGroup.InstanceCount > 0 || *scalingGroup.DesiredCapacity > 0 {
		if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			return retryError(asService.ClearScalingGroupInstance(ctx, scalingGroupId))
		}); err != nil {
			return err
		}
	}

	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		if errRet := asService.DeleteScalingGroup(ctx, scalingGroupId); errRet != nil {
			if sdkErr, ok := errRet.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == AsScalingGroupNotFound {
					return nil
				} else if sdkErr.Code == AsScalingGroupInProgress || sdkErr.Code == AsScalingGroupInstanceInGroup {
					return resource.RetryableError(sdkErr)
				}
			}
			return resource.NonRetryableError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
