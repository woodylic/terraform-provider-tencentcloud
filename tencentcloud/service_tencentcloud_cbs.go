package tencentcloud

import (
	"context"
	"fmt"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	"log"

	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

type CbsService struct {
	client *connectivity.TencentCloudClient
}

func (me *CbsService) DescribeDiskById(ctx context.Context, diskId string) (disk *cbs.Disk, errRet error) {
	logId := getLogId(ctx)
	request := cbs.NewDescribeDisksRequest()
	request.DiskIds = []*string{&diskId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().DescribeDisks(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DiskSet) < 1 {
		errRet = fmt.Errorf("disk id is not found")
		return
	}
	disk = response.Response.DiskSet[0]
	return
}

func (me *CbsService) DescribeDisksByFilter(ctx context.Context, params map[string]string) (disks []*cbs.Disk, errRet error) {
	logId := getLogId(ctx)
	request := cbs.NewDescribeDisksRequest()
	request.Filters = make([]*cbs.Filter, 0, len(params))
	for k, v := range params {
		filter := &cbs.Filter{
			Name:   stringToPointer(k),
			Values: []*string{stringToPointer(v)},
		}
		request.Filters = append(request.Filters, filter)
	}

	offset := 0
	pageSize := 100
	disks = make([]*cbs.Disk, 0)
	for {
		request.Offset = intToPointer(offset)
		request.Limit = intToPointer(pageSize)
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCbsClient().DescribeDisks(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.DiskSet) < 1 {
			break
		}

		disks = append(disks, response.Response.DiskSet...)

		if len(response.Response.DiskSet) < pageSize {
			break
		}
		offset += pageSize
	}
	return
}

func (me *CbsService) ModifyDiskAttributes(ctx context.Context, diskId, diskName string, projectId int) error {
	logId := getLogId(ctx)
	request := cbs.NewModifyDiskAttributesRequest()
	request.DiskIds = []*string{&diskId}
	if diskName != "" {
		request.DiskName = &diskName
	}
	if projectId >= 0 {
		request.ProjectId = intToPointer(projectId)
	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().ModifyDiskAttributes(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CbsService) DeleteDiskById(ctx context.Context, diskId string) error {
	logId := getLogId(ctx)
	request := cbs.NewTerminateDisksRequest()
	request.DiskIds = []*string{&diskId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().TerminateDisks(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CbsService) ResizeDisk(ctx context.Context, diskId string, diskSize int) error {
	logId := getLogId(ctx)
	request := cbs.NewResizeDiskRequest()
	request.DiskId = &diskId
	request.DiskSize = intToPointer(diskSize)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().ResizeDisk(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CbsService) ApplySnapshot(ctx context.Context, diskId, snapshotId string) error {
	logId := getLogId(ctx)
	request := cbs.NewApplySnapshotRequest()
	request.DiskId = &diskId
	request.SnapshotId = &snapshotId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().ApplySnapshot(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CbsService) AttachDisk(ctx context.Context, diskId, instanceId string) error {
	logId := getLogId(ctx)
	request := cbs.NewAttachDisksRequest()
	request.DiskIds = []*string{&diskId}
	request.InstanceId = &instanceId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().AttachDisks(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CbsService) DetachDisk(ctx context.Context, diskId, instanceId string) error {
	logId := getLogId(ctx)
	request := cbs.NewDetachDisksRequest()
	request.DiskIds = []*string{&diskId}
	request.InstanceId = &instanceId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().DetachDisks(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CbsService) CreateSnapshot(ctx context.Context, diskId, snapshotName string) (snapshotId string, errRet error) {
	logId := getLogId(ctx)
	request := cbs.NewCreateSnapshotRequest()
	request.DiskId = &diskId
	request.SnapshotName = &snapshotName
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().CreateSnapshot(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	snapshotId = *response.Response.SnapshotId
	return
}

func (me *CbsService) DescribeSnapshotById(ctx context.Context, snapshotId string) (snapshot *cbs.Snapshot, errRet error) {
	logId := getLogId(ctx)
	request := cbs.NewDescribeSnapshotsRequest()
	request.SnapshotIds = []*string{&snapshotId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().DescribeSnapshots(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.SnapshotSet) < 1 {
		errRet = fmt.Errorf("snapshot id %s is not found", snapshotId)
		return
	}
	snapshot = response.Response.SnapshotSet[0]
	return
}

func (me *CbsService) DescribeSnapshotsByFilter(ctx context.Context, params map[string]string) (snapshots []*cbs.Snapshot, errRet error) {
	logId := getLogId(ctx)
	request := cbs.NewDescribeSnapshotsRequest()
	request.Filters = make([]*cbs.Filter, 0, len(params))
	for k, v := range params {
		filter := &cbs.Filter{
			Name:   stringToPointer(k),
			Values: []*string{stringToPointer(v)},
		}
		request.Filters = append(request.Filters, filter)
	}

	offset := 0
	pageSize := 100
	for {
		request.Offset = intToPointer(offset)
		request.Limit = intToPointer(pageSize)
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCbsClient().DescribeSnapshots(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.SnapshotSet) < 1 {
			break
		}

		snapshots = append(snapshots, response.Response.SnapshotSet...)

		if len(response.Response.SnapshotSet) < pageSize {
			break
		}
		offset += pageSize
	}
	return
}

func (me *CbsService) ModifySnapshotName(ctx context.Context, snapshotId, snapshotName string) error {
	logId := getLogId(ctx)
	request := cbs.NewModifySnapshotAttributeRequest()
	request.SnapshotId = &snapshotId
	request.SnapshotName = &snapshotName
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().ModifySnapshotAttribute(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CbsService) DeleteSnapshot(ctx context.Context, snapshotId string) error {
	logId := getLogId(ctx)
	request := cbs.NewDeleteSnapshotsRequest()
	request.SnapshotIds = []*string{&snapshotId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().DeleteSnapshots(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CbsService) DescribeSnapshotPolicyById(ctx context.Context, policyId string) (policy *cbs.AutoSnapshotPolicy, errRet error) {
	logId := getLogId(contextNil)
	request := cbs.NewDescribeAutoSnapshotPoliciesRequest()
	request.AutoSnapshotPolicyIds = []*string{&policyId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().DescribeAutoSnapshotPolicies(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AutoSnapshotPolicySet) < 1 {
		errRet = fmt.Errorf("snapshot policy id %s is not found", policyId)
		return
	}
	policy = response.Response.AutoSnapshotPolicySet[0]
	return
}

func (me *CbsService) DeleteSnapshotPolicy(ctx context.Context, policyId string) error {
	logId := getLogId(ctx)
	request := cbs.NewDeleteAutoSnapshotPoliciesRequest()
	request.AutoSnapshotPolicyIds = []*string{&policyId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().DeleteAutoSnapshotPolicies(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func flattenCbsTagsMapping(tags []*cbs.Tag) (mapping map[string]string) {
	mapping = make(map[string]string)
	for _, tag := range tags {
		mapping[*tag.Key] = *tag.Value
	}
	return
}
