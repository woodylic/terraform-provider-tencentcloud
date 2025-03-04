## 1.24.2 (Unreleased)
## 1.24.1 (November 26, 2019)

ENHANCEMENTS:

* Resource: `tencentcloud_kubernetes_cluster` add support for `PREPAID` instance type. Thanks to @woodylic ([#204](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/204)).
* Resource: `tencentcloud_cos_bucket` add optional argument tags
* Data Source: `tencentcloud_cos_buckets` add optional argument tags

BUG FIXES:
* Fixed docs issues of `tencentcloud_nat_gateway`

## 1.24.0 (November 20, 2019)

FEATURES:

* **New Data Source**: `tencentcloud_ha_vips`
* **New Data Source**: `tencentcloud_ha_vip_eip_attachments`
* **New Resource**: `tencentcloud_ha_vip`
* **New Resource**: `tencentcloud_ha_vip_eip_attachment`

ENHANCEMENTS:

* Resource: `tencentcloud_kubernetes_cluster` cluster_os add new support: `centos7.6x86_64` and `ubuntu18.04.1 LTSx86_64` 
* Resource: `tencentcloud_nat_gateway` add computed argument `created_time`.

BUG FIXES:

* Fixed docs issues of CAM, DNAT and NAT_GATEWAY
* Fixed query issue that paged-query was not supported in data source `tencentcloud_dnats`
* Fixed query issue that filter `address_ip` was set incorrectly in data source `tencentcloud_eips`

## 1.23.0 (November 14, 2019)

FEATURES:

* **New Data Source**: `tencentcloud_images`
* **New Data Source**: `tencentcloud_vpn_gateways`
* **New Data Source**: `tencentcloud_customer_gateways`
* **New Data Source**: `tencentcloud_vpn_connections`
* **New Resource**: `tencentcloud_vpn_gateway`
* **New Resource**: `tencentcloud_customer_gateway`
* **New Resource**: `tencentcloud_vpn_connection`
* **Provider TencentCloud**: add `security_token` argument

ENHANCEMENTS:

* All api calls now using api3.0
* Resource: `tencentcloud_eip` add optional argument `tags`.
* Data Source: `tencentcloud_eips` add optional argument `tags`.

BUG FIXES:

* Fixed docs of CAM

## 1.22.0 (November 05, 2019)

FEATURES:

* **New Data Source**: `tencentcloud_cfs_file_systems`
* **New Data Source**: `tencentcloud_cfs_access_groups`
* **New Data Source**: `tencentcloud_cfs_access_rules`
* **New Data Source**: `tencentcloud_scf_functions`
* **New Data Source**: `tencentcloud_scf_namespaces`
* **New Data Source**: `tencentcloud_scf_logs`
* **New Resource**: `tencentcloud_cfs_file_system`
* **New Resource**: `tencentcloud_cfs_access_group`
* **New Resource**: `tencentcloud_cfs_access_rule`
* **New Resource**: `tencentcloud_scf_function`
* **New Resource**: `tencentcloud_scf_namespace`

## 1.21.2 (October 29, 2019)

BUG FIXES:

* Resource: `tencentcloud_gaap_realserver` add ip/domain exists check
* Resource: `tencentcloud_kubernetes_cluster` add error handling logic and optional argument `tags`.
* Resource: `tencentcloud_kubernetes_scale_worker` add error handling logic.
* Data Source: `tencentcloud_kubernetes_clusters` add optional argument `tags`.

## 1.21.1 (October 23, 2019)

ENHANCEMENTS:

* Updated golang to version 1.13.x

BUG FIXES:

* Fixed docs of CAM

## 1.21.0 (October 15, 2019)

FEATURES:

* **New Data Source**: `tencentcloud_cam_users`
* **New Data Source**: `tencentcloud_cam_groups`
* **New Data Source**: `tencentcloud_cam_policies`
* **New Data Source**: `tencentcloud_cam_roles`
* **New Data Source**: `tencentcloud_cam_user_policy_attachments`
* **New Data Source**: `tencentcloud_cam_group_policy_attachments`
* **New Data Source**: `tencentcloud_cam_role_policy_attachments`
* **New Data Source**: `tencentcloud_cam_group_memberships`
* **New Data Source**: `tencentcloud_cam_saml_providers`
* **New Data Source**: `tencentcloud_reserved_instance_configs`
* **New Data Source**: `tencentcloud_reserved_instances`
* **New Resource**: `tencentcloud_cam_user`
* **New Resource**: `tencentcloud_cam_group`
* **New Resource**: `tencentcloud_cam_role`
* **New Resource**: `tencentcloud_cam_policy`
* **New Resource**: `tencentcloud_cam_user_policy_attachment`
* **New Resource**: `tencentcloud_cam_group_policy_attachment`
* **New Resource**: `tencentcloud_cam_role_policy_attachment`
* **New Resource**: `tencentcloud_cam_group_membership`
* **New Resource**: `tencentcloud_cam_saml_provider`
* **New Resource**: `tencentcloud_reserved_instance`

ENHANCEMENTS:

* Resource: `tencentcloud_gaap_http_domain` support import
* Resource: `tencentcloud_gaap_layer7_listener` support import

BUG FIXES:

* Resource: `tencentcloud_gaap_http_domain` fix sometimes can't enable realserver auth

## 1.20.1 (October 08, 2019)

ENHANCEMENTS:

* Data Source: `tencentcloud_availability_zones` refactor logic with api3.0 .
* Data Source: `tencentcloud_as_scaling_groups` add optional argument `tags` and attribute `tags` for `scaling_group_list`.
* Resource: `tencentcloud_eip` add optional argument `type`, `anycast_zone`, `internet_service_provider`, etc.
* Resource: `tencentcloud_as_scaling_group` add optional argument `tags`.

BUG FIXES:

* Data Source: `tencentcloud_gaap_http_domains` set response `certificate_id`, `client_certificate_id`, `realserver_auth`, `basic_auth` and `gaap_auth` default value when they are nil.
* Resource: `tencentcloud_gaap_http_domain` set response `certificate_id`, `client_certificate_id`, `realserver_auth`, `basic_auth` and `gaap_auth` default value when they are nil.

## 1.20.0 (September 24, 2019)

FEATURES:

* **New Data Source**: `tencentcloud_eips`
* **New Data Source**: `tencentcloud_instances`
* **New Data Source**: `tencentcloud_key_pairs`
* **New Data Source**: `tencentcloud_placement_groups`
* **New Resource**: `tencentcloud_placement_group`

ENHANCEMENTS:

* Data Source: `tencentcloud_redis_instances` add optional argument `tags`.
* Data Source: `tencentcloud_mongodb_instances` add optional argument `tags`.
* Data Source: `tencentcloud_instance_types` add optional argument `availability_zone` and `gpu_core_count`.
* Data Source: `tencentcloud_gaap_http_rules` add optional argument `forward_host` and attributes `forward_host` in `rules`.
* Resource: `tencentcloud_redis_instance` add optional argument `tags`.
* Resource: `tencentcloud_mongodb_instance` add optional argument `tags`.
* Resource: `tencentcloud_mongodb_sharding_instance` add optional argument `tags`.
* Resource: `tencentcloud_instance` add optional argument `placement_group_id`.
* Resource: `tencentcloud_eip` refactor logic with api3.0 .
* Resource: `tencentcloud_eip_association` refactor logic with api3.0 .
* Resource: `tencentcloud_key_pair` refactor logic with api3.0 .
* Resource: `tencentcloud_gaap_http_rule` add optional argument `forward_host`.

BUG FIXES:
* Resource: `tencentcloud_mysql_instance`: miss argument `availability_zone` causes the instance to be recreated.

DEPRECATED:

* Data Source: `tencentcloud_eip` replaced by `tencentcloud_eips`.

## 1.19.0 (September 19, 2019)

FEATURES:

* **New Resource**: `tencentcloud_security_group_lite_rule`.

ENHANCEMENTS:

* Data Source: `tencentcloud_security_groups`: add optional argument `tags`.
* Data Source: `tencentcloud_security_groups`: add optional argument `result_output_file` and new attributes `ingress`, `egress` for list `security_groups`.
* Resource: `tencentcloud_security_group`: add optional argument `tags`.
* Resource: `tencentcloud_as_scaling_config`: internet charge type support `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.

BUG FIXES:
* Resource: `tencentcloud_clb_listener_rule`: fix unclear description and errors in example.
* Resource: `tencentcloud_instance`: fix hostname is not work.

## 1.18.1 (September 17, 2019)

FEATURES:

* **Update Data Source**: `tencentcloud_vpc_instances` add optional argument `tags`
* **Update Data Source**: `tencentcloud_vpc_subnets` add optional argument `tags`
* **Update Data Source**: `tencentcloud_route_tables` add optional argument `tags`
* **Update Resource**: `tencentcloud_vpc` add optional argument `tags`
* **Update Resource**: `tencentcloud_subnet` add optional argument `tags`
* **Update Resource**: `tencentcloud_route_table` add optional argument `tags`

ENHANCEMENTS:

* Data Source:`tencentcloud_kubernetes_clusters`  support pull out authentication information for cluster access too.
* Resource:`tencentcloud_kubernetes_cluster`  support pull out authentication information for cluster access.

BUG FIXES:

* Resource: `tencentcloud_mysql_instance`: when the mysql is abnormal state, read the basic information report error

DEPRECATED:

* Data Source: `tencentcloud_kubernetes_clusters`:`container_runtime` is no longer supported. 

## 1.18.0 (September 10, 2019)

FEATURES:

* **New Data Source**: `tencentcloud_ssl_certificates`
* **New Data Source**: `tencentcloud_dnats`
* **New Data Source**: `tencentcloud_nat_gateways`
* **New Resource**: `tencentcloud_ssl_certificate`
* **Update Resource**: `tencentcloud_clb_redirection` add optional argument `is_auto_rewrite`
* **Update Resource**: `tencentcloud_nat_gateway` , add more configurable items.
* **Update Resource**: `tencentcloud_nat` , add more configurable items.

DEPRECATED:
* Data Source: `tencentcloud_nats` replaced by `tencentcloud_nat_gateways`.

## 1.17.0 (September 04, 2019)

FEATURES:
* **New Data Source**: `tencentcloud_gaap_proxies`
* **New Data Source**: `tencentcloud_gaap_realservers`
* **New Data Source**: `tencentcloud_gaap_layer4_listeners`
* **New Data Source**: `tencentcloud_gaap_layer7_listeners`
* **New Data Source**: `tencentcloud_gaap_http_domains`
* **New Data Source**: `tencentcloud_gaap_http_rules`
* **New Data Source**: `tencentcloud_gaap_security_policies`
* **New Data Source**: `tencentcloud_gaap_security_rules`
* **New Data Source**: `tencentcloud_gaap_certificates`
* **New Resource**: `tencentcloud_gaap_proxy`
* **New Resource**: `tencentcloud_gaap_realserver`
* **New Resource**: `tencentcloud_gaap_layer4_listener`
* **New Resource**: `tencentcloud_gaap_layer7_listener`
* **New Resource**: `tencentcloud_gaap_http_domain`
* **New Resource**: `tencentcloud_gaap_http_rule`
* **New Resource**: `tencentcloud_gaap_certificate`
* **New Resource**: `tencentcloud_gaap_security_policy`
* **New Resource**: `tencentcloud_gaap_security_rule`

## 1.16.3 (August 30, 2019)

BUG FIXIES:

* Resource: `tencentcloud_kubernetes_cluster`: cgi error retry.
* Resource: `tencentcloud_kubernetes_scale_worker`: cgi error retry.

## 1.16.2 (August 28, 2019)

BUG FIXIES:

* Resource: `tencentcloud_instance`: fixed cvm data disks missing computed.
* Resource: `tencentcloud_mysql_backup_policy`: `backup_model` remove logical backup support. 
* Resource: `tencentcloud_mysql_instance`: `tags` adapt to the new official api.

## 1.16.1 (August 27, 2019)

ENHANCEMENTS:
* `tencentcloud_instance`: refactor logic with api3.0 .

## 1.16.0 (August 20, 2019)

FEATURES:
* **New Data Source**: `tencentcloud_kubernetes_clusters`
* **New Resource**: `tencentcloud_kubernetes_scale_worker`
* **New Resource**: `tencentcloud_kubernetes_cluster`

DEPRECATED:
* Data Source: `tencentcloud_container_clusters` replaced by `tencentcloud_kubernetes_clusters`.
* Data Source: `tencentcloud_container_cluster_instances` replaced by `tencentcloud_kubernetes_clusters`.
* Resource: `tencentcloud_container_cluster` replaced by `tencentcloud_kubernetes_cluster`.
* Resource: `tencentcloud_container_cluster_instance` replaced by `tencentcloud_kubernetes_scale_worker`.

## 1.15.2 (August 14, 2019)

ENHANCEMENTS:

* `tencentcloud_as_scaling_group`: fixed issue that binding scaling group to load balancer does not work.
* `tencentcloud_clb_attachements`: rename `rewrite_source_rule_id` with `source_rule_id` and rename `rewrite_target_rule_id` with `target_rule_id`.

## 1.15.1 (August 13, 2019)

ENHANCEMENTS:

* `tencentcloud_instance`: changed `image_id` property to ForceNew ([#78](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/78))
* `tencentcloud_instance`: improved with retry ([#82](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/82))
* `tencentcloud_cbs_storages`: improved with retry ([#82](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/82))
* `tencentcloud_clb_instance`: bug fixed and improved with retry ([#37](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/37))

## 1.15.0 (August 07, 2019)

FEATURES:
* **New Data Source**: `tencentcloud_clb_instances`
* **New Data Source**: `tencentcloud_clb_listeners`
* **New Data Source**: `tencentcloud_clb_listener_rules`
* **New Data Source**: `tencentcloud_clb_attachments`
* **New Data Source**: `tencentcloud_clb_redirections`
* **New Resource**: `tencentcloud_clb_instance`
* **New Resource**: `tencentcloud_clb_listener`
* **New Resource**: `tencentcloud_clb_listener_rule`
* **New Resource**: `tencentcloud_clb_attachment`
* **New Resource**: `tencentcloud_clb_redirection`

DEPRECATED:
* Resource: `tencentcloud_lb` replaced by `tencentcloud_clb_instance`.
* Resource: `tencentcloud_alb_server_attachment` replaced by `tencentcloud_clb_attachment`.

## 1.14.1 (August 05, 2019)

BUG FIXIES:

* resource/tencentcloud_security_group_rule: fixed security group rule id is not compatible with previous version.

## 1.14.0 (July 30, 2019)

FEATURES:
* **New Data Source**: `tencentcloud_security_groups`
* **New Data Source**: `tencentcloud_mongodb_instances`
* **New Data Source**: `tencentcloud_mongodb_zone_config`
* **New Resource**: `tencentcloud_mongodb_instance`
* **New Resource**: `tencentcloud_mongodb_sharding_instance`
* **Update Resource**: `tencentcloud_security_group_rule` add optional argument `description`

DEPRECATED:
* Data Source: `tencnetcloud_security_group` replaced by `tencentcloud_security_groups`

ENHANCEMENTS:
* Refactoring security_group logic with api3.0

## 1.13.0 (July 23, 2019)

FEATURES:
* **New Data Source**: `tencentcloud_dc_gateway_instances`
* **New Data Source**: `tencentcloud_dc_gateway_ccn_routes`
* **New Resource**: `tencentcloud_dc_gateway`
* **New Resource**: `tencentcloud_dc_gateway_ccn_route`

## 1.12.0 (July 16, 2019)

FEATURES:
* **New Data Source**: `tencentcloud_dc_instances`
* **New Data Source**: `tencentcloud_dcx_instances`
* **New Resource**: `tencentcloud_dcx`
* **UPDATE Resource**:`tencentcloud_mysql_instance` and `tencentcloud_mysql_readonly_instance` completely delete instance. 

BUG FIXIES:

* resource/tencentcloud_instance: fixed issue when data disks set as delete_with_instance not works.
* resource/tencentcloud_instance: if managed public_ip manually, please don't define `allocate_public_ip` ([#62](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/62)).
* resource/tencentcloud_eip_association: fixed issue when instances were manually deleted ([#60](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/60)).
* resource/tencentcloud_mysql_readonly_instance:remove an unsupported property `gtid`

## 1.11.0 (July 02, 2019)

FEATURES:
* **New Data Source**: `tencentcloud_ccn_instances`
* **New Data Source**: `tencentcloud_ccn_bandwidth_limits`
* **New Resource**: `tencentcloud_ccn`
* **New Resource**: `tencentcloud_ccn_attachment`
* **New Resource**: `tencentcloud_ccn_bandwidth_limit`

## 1.10.0 (June 27, 2019)

ENHANCEMENTS:

* Refactoring vpc logic with api3.0
* Refactoring cbs logic with api3.0

FEATURES:
* **New Data Source**: `tencentcloud_vpc_instances`
* **New Data Source**: `tencentcloud_vpc_subnets`
* **New Data Source**: `tencentcloud_vpc_route_tables`
* **New Data Source**: `tencentcloud_cbs_storages`
* **New Data Source**: `tencentcloud_cbs_snapshots`
* **New Resource**: `tencentcloud_route_table_entry`
* **New Resource**: `tencentcloud_cbs_snapshot_policy`
* **Update Resource**: `tencentcloud_vpc` , add more configurable items.
* **Update Resource**: `tencentcloud_subnet` , add more configurable items.
* **Update Resource**: `tencentcloud_route_table`, add more configurable items.
* **Update Resource**: `tencentcloud_cbs_storage`, add more configurable items.
* **Update Resource**: `tencentcloud_instance`: add optional argument `tags`.
* **Update Resource**: `tencentcloud_security_group_rule`: add optional argument `source_sgid`.
 
DEPRECATED:
* Data Source: `tencentcloud_vpc` replaced by `tencentcloud_vpc_instances`.
* Data Source: `tencentcloud_subnet` replaced by  `tencentcloud_vpc_subnets`.
* Data Source: `tencentcloud_route_table` replaced by `tencentcloud_vpc_route_tables`.
* Resource: `tencentcloud_route_entry` replaced by `tencentcloud_route_table_entry`.

## 1.9.1 (June 24, 2019)

BUG FIXIES:

* data/tencentcloud_instance: fixed vpc ip is in use error when re-creating with private ip ([#46](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/46)).

## 1.9.0 (June 18, 2019)

ENHANCEMENTS:

* update to `v0.12.1` Terraform SDK version

BUG FIXIES:

* data/tencentcloud_security_group: `project_id` remote API return sometime is string type.
* resource/tencentcloud_security_group: just like `data/tencentcloud_security_group`

## 1.8.0 (June 11, 2019)

FEATURES:
* **New Data Source**: `tencentcloud_as_scaling_configs`
* **New Data Source**: `tencentcloud_as_scaling_groups`
* **New Data Source**: `tencentcloud_as_scaling_policies`
* **New Resource**: `tencentcloud_as_scaling_config`
* **New Resource**: `tencentcloud_as_scaling_group`
* **New Resource**: `tencentcloud_as_attachment`
* **New Resource**: `tencentcloud_as_scaling_policy`
* **New Resource**: `tencentcloud_as_schedule`
* **New Resource**: `tencentcloud_as_lifecycle_hook`
* **New Resource**: `tencentcloud_as_notification`

## 1.7.0 (May 23, 2019)

FEATURES:
* **New Data Source**: `tencentcloud_redis_zone_config`
* **New Data Source**: `tencentcloud_redis_instances`
* **New Resource**: `tencentcloud_redis_instance`
* **New Resource**: `tencentcloud_redis_backup_config`

ENHANCEMENTS:

* resource/tencentcloud_instance: Add `hostname`, `project_id`, `delete_with_instance` argument.
* Update tencentcloud-sdk-go to better support redis api.

## 1.6.0 (May 15, 2019)

FEATURES:

* **New Data Source**: `tencentcloud_cos_buckets`
* **New Data Source**: `tencentcloud_cos_bucket_object`
* **New Resource**: `tencentcloud_cos_bucket`
* **New Resource**: `tencentcloud_cos_bucket_object`

ENHANCEMENTS:

* Add the framework of auto generating terraform docs

## 1.5.0 (April 26, 2019)

FEATURES:

* **New Data Source**: `tencentcloud_mysql_backup_list`
* **New Data Source**: `tencentcloud_mysql_zone_config`
* **New Data Source**: `tencentcloud_mysql_parameter_list`
* **New Data Source**: `tencentcloud_mysql_instance`
* **New Resource**: `tencentcloud_mysql_backup_policy`
* **New Resource**: `tencentcloud_mysql_account`
* **New Resource**: `tencentcloud_mysql_account_privilege`
* **New Resource**: `tencentcloud_mysql_instance`
* **New Resource**: `tencentcloud_mysql_readonly_instance`

ENHANCEMENTS:

* resource/tencentcloud_subnet: `route_table_id` now is an optional argument

## 1.4.0 (April 12, 2019)

ENHANCEMENTS:

* data/tencentcloud_image: add `image_name` attribute to this data source.
* resource/tencentcloud_instance: data disk count limit now is upgrade from 1 to 10, as API has supported more disks.
* resource/tencentcloud_instance: PREPAID instance now can be deleted, but still have some limit in API.

BUG FIXIES:

* resource/tencentcloud_instance: `allocate_public_ip` doesn't work properly when it is set to false.

## 1.3.0 (March 12, 2019)

FEATURES:

* **New Resource**: `tencentcloud_lb` ([#3](https://github.com/terraform-providers/terraform-provider-scaffolding/issues/3))

ENHANCEMENTS:

* resource/tencentcloud_instance: Add `user_data_raw` argument ([#4](https://github.com/terraform-providers/terraform-provider-scaffolding/issues/4))

## 1.2.2 (September 28, 2018)

BUG FIXES:

* resource/tencentcloud_cbs_storage: make name to be required ([#25](https://github.com/tencentyun/terraform-provider-tencentcloud/issues/25))
* resource/tencentcloud_instance: support user data and private ip

## 1.2.0 (April 3, 2018)

FEATURES:

* **New Resource**: `tencentcloud_container_cluster`
* **New Resource**: `tencentcloud_container_cluster_instance`
* **New Data Source**: `tencentcloud_container_clusters`
* **New Data Source**: `tencentcloud_container_cluster_instances`

## 1.1.0 (March 9, 2018)

FEATURES:

* **New Resource**: `tencentcloud_eip`
* **New Resource**: `tencentcloud_eip_association`
* **New Data Source**: `tencentcloud_eip`
* **New Resource**: `tencentcloud_nat_gateway`
* **New Resource**: `tencentcloud_dnat`
* **New Data Source**: `tencentcloud_nats`
* **New Resource**: `tencentcloud_cbs_snapshot`
* **New Resource**: `tencentcloud_alb_server_attachment`

## 1.0.0 (January 19, 2018)

FEATURES:

### CVM

RESOURCES:

* instance create
* instance read
* instance update
    * reset instance
    * reset password
    * update instance name
    * update security groups
* instance delete
* key pair create
* key pair read
* key pair delete

DATA SOURCES:

* image read
* instance\_type read
* zone read

### VPC

RESOURCES:

* vpc create
* vpc read
* vpc update (update name)
* vpc delete
* subnet create
* subnet read
* subnet update (update name)
* subnet delete
* security group create
* security group read
* security group update (update name, description)
* security group delete
* security group rule create
* security group rule read
* security group rule delete
* route table create
* route table read
* route table update (update name)
* route table delete
* route entry create
* route entry read
* route entry delete

DATA SOURCES:

* vpc read
* subnet read
* security group read
* route table read

### CBS

RESOURCES:

* storage create
* storage read
* storage update (update name)
* storage attach
* storage detach
