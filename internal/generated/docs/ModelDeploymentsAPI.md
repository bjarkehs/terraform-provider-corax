# \ModelDeploymentsAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateModelDeploymentV1ModelDeploymentsPost**](ModelDeploymentsAPI.md#CreateModelDeploymentV1ModelDeploymentsPost) | **Post** /v1/model-deployments | Create Model Deployment
[**DeleteModelDeploymentV1ModelDeploymentsDeploymentIdDelete**](ModelDeploymentsAPI.md#DeleteModelDeploymentV1ModelDeploymentsDeploymentIdDelete) | **Delete** /v1/model-deployments/{deployment_id} | Delete Model Deployment
[**GetModelDeploymentV1ModelDeploymentsDeploymentIdGet**](ModelDeploymentsAPI.md#GetModelDeploymentV1ModelDeploymentsDeploymentIdGet) | **Get** /v1/model-deployments/{deployment_id} | Get Model Deployment
[**ListModelDeploymentsV1ModelDeploymentsGet**](ModelDeploymentsAPI.md#ListModelDeploymentsV1ModelDeploymentsGet) | **Get** /v1/model-deployments | List Model Deployments
[**UpdateModelDeploymentV1ModelDeploymentsDeploymentIdPut**](ModelDeploymentsAPI.md#UpdateModelDeploymentV1ModelDeploymentsDeploymentIdPut) | **Put** /v1/model-deployments/{deployment_id} | Update Model Deployment



## CreateModelDeploymentV1ModelDeploymentsPost

> ModelDeployment CreateModelDeploymentV1ModelDeploymentsPost(ctx).ModelDeploymentCreate(modelDeploymentCreate).Execute()

Create Model Deployment



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/api"
)

func main() {
	modelDeploymentCreate := *openapiclient.NewModelDeploymentCreate("Name_example", []openapiclient.CapabilityType{openapiclient.CapabilityType("chat")}, map[string]interface{}{"key": interface{}(123)}, "ProviderId_example") // ModelDeploymentCreate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ModelDeploymentsAPI.CreateModelDeploymentV1ModelDeploymentsPost(context.Background()).ModelDeploymentCreate(modelDeploymentCreate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ModelDeploymentsAPI.CreateModelDeploymentV1ModelDeploymentsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateModelDeploymentV1ModelDeploymentsPost`: ModelDeployment
	fmt.Fprintf(os.Stdout, "Response from `ModelDeploymentsAPI.CreateModelDeploymentV1ModelDeploymentsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateModelDeploymentV1ModelDeploymentsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **modelDeploymentCreate** | [**ModelDeploymentCreate**](ModelDeploymentCreate.md) |  | 

### Return type

[**ModelDeployment**](ModelDeployment.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteModelDeploymentV1ModelDeploymentsDeploymentIdDelete

> DeleteModelDeploymentV1ModelDeploymentsDeploymentIdDelete(ctx, deploymentId).Execute()

Delete Model Deployment



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/api"
)

func main() {
	deploymentId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ModelDeploymentsAPI.DeleteModelDeploymentV1ModelDeploymentsDeploymentIdDelete(context.Background(), deploymentId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ModelDeploymentsAPI.DeleteModelDeploymentV1ModelDeploymentsDeploymentIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**deploymentId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteModelDeploymentV1ModelDeploymentsDeploymentIdDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetModelDeploymentV1ModelDeploymentsDeploymentIdGet

> ModelDeployment GetModelDeploymentV1ModelDeploymentsDeploymentIdGet(ctx, deploymentId).Execute()

Get Model Deployment



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/api"
)

func main() {
	deploymentId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ModelDeploymentsAPI.GetModelDeploymentV1ModelDeploymentsDeploymentIdGet(context.Background(), deploymentId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ModelDeploymentsAPI.GetModelDeploymentV1ModelDeploymentsDeploymentIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetModelDeploymentV1ModelDeploymentsDeploymentIdGet`: ModelDeployment
	fmt.Fprintf(os.Stdout, "Response from `ModelDeploymentsAPI.GetModelDeploymentV1ModelDeploymentsDeploymentIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**deploymentId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetModelDeploymentV1ModelDeploymentsDeploymentIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ModelDeployment**](ModelDeployment.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListModelDeploymentsV1ModelDeploymentsGet

> PagedResponseModelModelDeployment ListModelDeploymentsV1ModelDeploymentsGet(ctx).Page(page).Size(size).Sort(sort).Filter(filter).Execute()

List Model Deployments



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/api"
)

func main() {
	page := int32(56) // int32 |  (optional) (default to 1)
	size := int32(56) // int32 |  (optional) (default to 10)
	sort := "sort_example" // string |  (optional) (default to "id")
	filter := "filter_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ModelDeploymentsAPI.ListModelDeploymentsV1ModelDeploymentsGet(context.Background()).Page(page).Size(size).Sort(sort).Filter(filter).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ModelDeploymentsAPI.ListModelDeploymentsV1ModelDeploymentsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListModelDeploymentsV1ModelDeploymentsGet`: PagedResponseModelModelDeployment
	fmt.Fprintf(os.Stdout, "Response from `ModelDeploymentsAPI.ListModelDeploymentsV1ModelDeploymentsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListModelDeploymentsV1ModelDeploymentsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** |  | [default to 1]
 **size** | **int32** |  | [default to 10]
 **sort** | **string** |  | [default to &quot;id&quot;]
 **filter** | **string** |  | 

### Return type

[**PagedResponseModelModelDeployment**](PagedResponseModelModelDeployment.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateModelDeploymentV1ModelDeploymentsDeploymentIdPut

> ModelDeployment UpdateModelDeploymentV1ModelDeploymentsDeploymentIdPut(ctx, deploymentId).ModelDeploymentUpdate(modelDeploymentUpdate).Execute()

Update Model Deployment



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/api"
)

func main() {
	deploymentId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	modelDeploymentUpdate := *openapiclient.NewModelDeploymentUpdate("Name_example", []openapiclient.CapabilityType{openapiclient.CapabilityType("chat")}, map[string]interface{}{"key": interface{}(123)}, "ProviderId_example") // ModelDeploymentUpdate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ModelDeploymentsAPI.UpdateModelDeploymentV1ModelDeploymentsDeploymentIdPut(context.Background(), deploymentId).ModelDeploymentUpdate(modelDeploymentUpdate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ModelDeploymentsAPI.UpdateModelDeploymentV1ModelDeploymentsDeploymentIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateModelDeploymentV1ModelDeploymentsDeploymentIdPut`: ModelDeployment
	fmt.Fprintf(os.Stdout, "Response from `ModelDeploymentsAPI.UpdateModelDeploymentV1ModelDeploymentsDeploymentIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**deploymentId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateModelDeploymentV1ModelDeploymentsDeploymentIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **modelDeploymentUpdate** | [**ModelDeploymentUpdate**](ModelDeploymentUpdate.md) |  | 

### Return type

[**ModelDeployment**](ModelDeployment.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

