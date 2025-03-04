package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/pkg/errors"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type scfFunctionInfo struct {
	name        string
	handler     *string
	desc        *string
	memSize     *int
	timeout     *int
	environment map[string]string
	runtime     *string
	vpcId       *string
	subnetId    *string
	role        *string
	clsLogsetId *string
	clsTopicId  *string
	namespace   *string
	l5Enable    *bool

	cosBucketName   *string
	cosObjectName   *string
	cosBucketRegion *string

	zipFile *string
}

type scfTrigger struct {
	name        string
	triggerType string
	triggerDesc string
}

type ScfService struct {
	client *connectivity.TencentCloudClient
}

func (me *ScfService) CreateFunction(ctx context.Context, info scfFunctionInfo) error {
	client := me.client.UseScfClient()

	request := scf.NewCreateFunctionRequest()
	request.FunctionName = &info.name
	request.Handler = info.handler
	request.Description = info.desc
	request.MemorySize = int64ToPointer(*info.memSize)
	request.Timeout = int64ToPointer(*info.timeout)
	for k, v := range info.environment {
		if request.Environment == nil {
			request.Environment = new(scf.Environment)
		}
		request.Environment.Variables = append(request.Environment.Variables, &scf.Variable{
			Key:   stringToPointer(k),
			Value: stringToPointer(v),
		})
	}
	request.Runtime = info.runtime

	if info.vpcId != nil {
		request.VpcConfig = &scf.VpcConfig{
			VpcId:    info.vpcId,
			SubnetId: info.subnetId,
		}
	}

	request.Namespace = info.namespace
	request.Role = info.role
	request.ClsLogsetId = info.clsLogsetId
	request.ClsTopicId = info.clsTopicId
	request.Type = stringToPointer(SCF_FUNCTION_TYPE_EVENT)

	request.Code = &scf.Code{
		CosBucketName:   info.cosBucketName,
		CosObjectName:   info.cosObjectName,
		CosBucketRegion: info.cosBucketRegion,
		ZipFile:         info.zipFile,
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.CreateFunction(request); err != nil {
			return retryError(errors.WithStack(err), "InternalError")
		}
		return nil
	}); err != nil {
		return err
	}

	return waitScfFunctionReady(ctx, info.name, *info.namespace, client)
}

func (me *ScfService) DescribeFunction(ctx context.Context, name, namespace string) (resp *scf.GetFunctionResponse, err error) {
	request := scf.NewGetFunctionRequest()
	request.FunctionName = &name
	request.Namespace = &namespace

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseScfClient().GetFunction(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				for _, code := range SCF_FUNCTIONS_NOT_FOUND_SET {
					if sdkError.Code == code {
						return nil
					}
				}
			}

			return retryError(errors.WithStack(err), "InternalError")
		}

		resp = response
		return nil
	}); err != nil {
		return nil, err
	}

	return
}

func (me *ScfService) DescribeFunctions(ctx context.Context, name, namespace, desc *string, tags map[string]string) (functions []*scf.Function, err error) {
	request := scf.NewListFunctionsRequest()
	request.SearchKey = name
	request.Namespace = namespace
	request.Description = desc
	for k, v := range tags {
		request.Filters = append(request.Filters, &scf.Filter{
			Name:   stringToPointer("tag-" + k),
			Values: []*string{stringToPointer(v)},
		})
	}
	request.Limit = int64ToPointer(SCF_FUNCTION_DESCRIBE_LIMIT)

	var offset int64
	count := SCF_FUNCTION_DESCRIBE_LIMIT

	// at least run loop once
	for count == SCF_FUNCTION_DESCRIBE_LIMIT {
		request.Offset = &offset

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := me.client.UseScfClient().ListFunctions(request)
			if err != nil {
				return retryError(errors.WithStack(err))
			}

			functions = append(functions, response.Response.Functions...)
			count = len(response.Response.Functions)

			return nil
		}); err != nil {
			return nil, err
		}

		offset += int64(count)
	}

	return
}

func (me *ScfService) ModifyFunctionCode(ctx context.Context, info scfFunctionInfo) error {
	client := me.client.UseScfClient()

	request := scf.NewUpdateFunctionCodeRequest()
	request.FunctionName = &info.name
	request.Handler = info.handler
	request.Namespace = info.namespace
	request.Code = &scf.Code{
		CosBucketName:   info.cosBucketName,
		CosObjectName:   info.cosObjectName,
		CosBucketRegion: info.cosBucketRegion,
		ZipFile:         info.zipFile,
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.UpdateFunctionCode(request); err != nil {
			return retryError(errors.WithStack(err), "InternalError")
		}
		return nil
	}); err != nil {
		return err
	}

	return waitScfFunctionReady(ctx, info.name, *info.namespace, client)
}

func (me *ScfService) ModifyFunctionConfig(ctx context.Context, info scfFunctionInfo) error {
	client := me.client.UseScfClient()

	request := scf.NewUpdateFunctionConfigurationRequest()
	request.FunctionName = &info.name
	request.Description = info.desc
	if info.memSize != nil {
		request.MemorySize = int64ToPointer(*info.memSize)
	}
	if info.timeout != nil {
		request.Timeout = int64ToPointer(*info.timeout)
	}
	request.Runtime = info.runtime

	request.Environment = new(scf.Environment)
	for k, v := range info.environment {
		request.Environment.Variables = append(request.Environment.Variables, &scf.Variable{
			Key:   stringToPointer(k),
			Value: stringToPointer(v),
		})
	}
	// clean all environments
	if len(request.Environment.Variables) == 0 {
		request.Environment.Variables = []*scf.Variable{
			{
				Key:   stringToPointer(""),
				Value: stringToPointer(""),
			},
		}
	}

	request.Namespace = info.namespace
	if info.vpcId != nil {
		request.VpcConfig = &scf.VpcConfig{VpcId: info.vpcId}
	}
	if info.subnetId != nil {
		if request.VpcConfig == nil {
			request.VpcConfig = new(scf.VpcConfig)
		}
		request.VpcConfig.SubnetId = info.subnetId
	}
	request.Role = info.role
	request.ClsLogsetId = info.clsLogsetId
	request.ClsTopicId = info.clsTopicId
	if info.l5Enable != nil {
		request.L5Enable = stringToPointer("FALSE")
		if *info.l5Enable {
			request.L5Enable = stringToPointer("TRUE")
		}
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.UpdateFunctionConfiguration(request); err != nil {
			return retryError(errors.WithStack(err), "InternalError")
		}
		return nil
	}); err != nil {
		return err
	}

	return waitScfFunctionReady(ctx, info.name, *info.namespace, client)
}

func (me *ScfService) DeleteFunction(ctx context.Context, name, namespace string) error {
	client := me.client.UseScfClient()

	deleteRequest := scf.NewDeleteFunctionRequest()
	deleteRequest.FunctionName = &name
	deleteRequest.Namespace = &namespace

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(deleteRequest.GetAction())

		if _, err := client.DeleteFunction(deleteRequest); err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				for _, code := range SCF_FUNCTIONS_NOT_FOUND_SET {
					if sdkError.Code == code {
						return nil
					}
				}
			}
			return retryError(errors.WithStack(err), "InternalError")
		}

		return nil
	}); err != nil {
		return err
	}

	descRequest := scf.NewGetFunctionRequest()
	descRequest.FunctionName = &name
	descRequest.Namespace = &namespace

	return resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(descRequest.GetAction())

		if _, err := client.GetFunction(descRequest); err == nil {
			return resource.RetryableError(errors.New("function still exists"))
		} else {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				for _, code := range SCF_FUNCTIONS_NOT_FOUND_SET {
					if sdkError.Code == code {
						return nil
					}
				}
			}

			return retryError(errors.WithStack(err), "InternalError")
		}
	})
}

func (me *ScfService) CreateNamespace(ctx context.Context, namespace, desc string) error {
	request := scf.NewCreateNamespaceRequest()
	request.Namespace = &namespace
	request.Description = &desc

	return resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := me.client.UseScfClient().CreateNamespace(request); err != nil {
			return retryError(errors.WithStack(err))
		}

		return nil
	})
}

func (me *ScfService) DescribeNamespace(ctx context.Context, namespace string) (ns *scf.Namespace, err error) {
	request := scf.NewListNamespacesRequest()
	request.Limit = int64ToPointer(SCF_NAMESPACE_DESCRIBE_LIMIT)

	var offset int64
	count := SCF_NAMESPACE_DESCRIBE_LIMIT

	// at least run loop once
	for count == SCF_NAMESPACE_DESCRIBE_LIMIT {
		request.Offset = &offset

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := me.client.UseScfClient().ListNamespaces(request)
			if err != nil {
				return retryError(errors.WithStack(err))
			}

			for _, respNs := range response.Response.Namespaces {
				if *respNs.Name == namespace {
					ns = respNs
					return nil
				}
			}

			count = len(response.Response.Namespaces)
			return nil
		}); err != nil {
			return nil, err
		}

		if ns != nil {
			return
		}

		offset += int64(count)
	}

	return
}

func (me *ScfService) DescribeNamespaces(ctx context.Context) (nss []*scf.Namespace, err error) {
	request := scf.NewListNamespacesRequest()
	request.Limit = int64ToPointer(SCF_NAMESPACE_DESCRIBE_LIMIT)

	var offset int64
	count := SCF_NAMESPACE_DESCRIBE_LIMIT

	// at least run loop once
	for count == SCF_NAMESPACE_DESCRIBE_LIMIT {
		request.Offset = &offset

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := me.client.UseScfClient().ListNamespaces(request)
			if err != nil {
				return retryError(errors.WithStack(err))
			}

			count = len(response.Response.Namespaces)
			nss = append(nss, response.Response.Namespaces...)

			return nil
		}); err != nil {
			return nil, err
		}

		offset += int64(count)
	}

	return
}

func (me *ScfService) ModifyNamespace(ctx context.Context, namespace, desc string) error {
	request := scf.NewUpdateNamespaceRequest()
	request.Namespace = &namespace
	request.Description = &desc

	return resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := me.client.UseScfClient().UpdateNamespace(request); err != nil {
			return retryError(errors.WithStack(err))
		}
		return nil
	})
}

func (me *ScfService) DeleteNamespace(ctx context.Context, namespace string) error {
	request := scf.NewDeleteNamespaceRequest()
	request.Namespace = &namespace

	return resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := me.client.UseScfClient().DeleteNamespace(request); err != nil {
			return retryError(errors.WithStack(err))
		}

		return nil
	})
}

func (me *ScfService) CreateTriggers(ctx context.Context, functionName, namespace string, triggers []scfTrigger) error {
	for _, trigger := range triggers {
		request := scf.NewCreateTriggerRequest()
		request.FunctionName = &functionName
		request.TriggerName = &trigger.name
		request.Type = &trigger.triggerType
		request.TriggerDesc = &trigger.triggerDesc
		request.Namespace = &namespace
		request.Enable = stringToPointer("OPEN")

		if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			if _, err := me.client.UseScfClient().CreateTrigger(request); err != nil {
				return retryError(errors.WithStack(err))
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func (me *ScfService) DeleteTriggers(ctx context.Context, functionName, namespace string, triggers []scfTrigger) error {
	for _, trigger := range triggers {
		request := scf.NewDeleteTriggerRequest()
		request.FunctionName = &functionName
		request.Namespace = &namespace
		request.TriggerName = &trigger.name
		request.Type = &trigger.triggerType
		request.TriggerDesc = &trigger.triggerDesc

		if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			if _, err := me.client.UseScfClient().DeleteTrigger(request); err != nil {
				return retryError(errors.WithStack(err))
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func (me *ScfService) DescribeLogs(
	ctx context.Context,
	fnName, namespace, order, orderBy string,
	offset, limit int,
	retCode, invokeRequestId, startTime, endTime *string,
) (logs []*scf.FunctionLog, err error) {
	request := scf.NewGetFunctionLogsRequest()
	request.FunctionName = &fnName
	request.Offset = int64ToPointer(offset)
	request.Limit = int64ToPointer(limit)
	request.Order = &order
	request.OrderBy = &orderBy
	if retCode != nil {
		request.Filter = &scf.LogFilter{RetCode: retCode}
	}
	request.Namespace = &namespace
	request.FunctionRequestId = invokeRequestId
	request.StartTime = startTime
	request.EndTime = endTime

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseScfClient().GetFunctionLogs(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				for _, code := range SCF_FUNCTIONS_NOT_FOUND_SET {
					if sdkError.Code == code {
						return nil
					}
				}
			}
			return retryError(errors.WithStack(err))
		}

		logs = response.Response.Data
		return nil
	}); err != nil {
		return nil, err
	}

	return
}

func waitScfFunctionReady(ctx context.Context, name, namespace string, client *scf.Client) error {
	request := scf.NewGetFunctionRequest()
	request.FunctionName = &name
	request.Namespace = &namespace

	return resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := client.GetFunction(request)
		if err != nil {
			return retryError(errors.WithStack(err), "InternalError")
		}

		switch *response.Response.Status {
		case SCF_FUNCTION_STATUS_CREATING, SCF_FUNCTION_STATUS_UPDATING:
			return resource.RetryableError(errors.New("function is not ready"))

		case SCF_FUNCTION_STATUS_ACTIVE:
			return nil

		default:
			return resource.NonRetryableError(errors.Errorf("function status is %s", *response.Response.Status))
		}
	})
}
