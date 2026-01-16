# \CapabilityTypesAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetCapabilityTypeV1CapabilityTypesCapabilityTypeGet**](CapabilityTypesAPI.md#GetCapabilityTypeV1CapabilityTypesCapabilityTypeGet) | **Get** /v1/capability-types/{capability_type} | Get Capability Type
[**ListCapabilityTypesV1CapabilityTypesGet**](CapabilityTypesAPI.md#ListCapabilityTypesV1CapabilityTypesGet) | **Get** /v1/capability-types | List Capability Types
[**UpdateCapabilityTypeV1CapabilityTypesCapabilityTypePut**](CapabilityTypesAPI.md#UpdateCapabilityTypeV1CapabilityTypesCapabilityTypePut) | **Put** /v1/capability-types/{capability_type} | Update Capability Type



## GetCapabilityTypeV1CapabilityTypesCapabilityTypeGet

> CapabilityTypeRepresentation GetCapabilityTypeV1CapabilityTypesCapabilityTypeGet(ctx, capabilityType).Execute()

Get Capability Type



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
	capabilityType := "capabilityType_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilityTypesAPI.GetCapabilityTypeV1CapabilityTypesCapabilityTypeGet(context.Background(), capabilityType).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilityTypesAPI.GetCapabilityTypeV1CapabilityTypesCapabilityTypeGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetCapabilityTypeV1CapabilityTypesCapabilityTypeGet`: CapabilityTypeRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CapabilityTypesAPI.GetCapabilityTypeV1CapabilityTypesCapabilityTypeGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**capabilityType** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetCapabilityTypeV1CapabilityTypesCapabilityTypeGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CapabilityTypeRepresentation**](CapabilityTypeRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListCapabilityTypesV1CapabilityTypesGet

> CapabilityTypesRepresentation ListCapabilityTypesV1CapabilityTypesGet(ctx).Execute()

List Capability Types



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilityTypesAPI.ListCapabilityTypesV1CapabilityTypesGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilityTypesAPI.ListCapabilityTypesV1CapabilityTypesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListCapabilityTypesV1CapabilityTypesGet`: CapabilityTypesRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CapabilityTypesAPI.ListCapabilityTypesV1CapabilityTypesGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListCapabilityTypesV1CapabilityTypesGetRequest struct via the builder pattern


### Return type

[**CapabilityTypesRepresentation**](CapabilityTypesRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateCapabilityTypeV1CapabilityTypesCapabilityTypePut

> CapabilityTypeRepresentation UpdateCapabilityTypeV1CapabilityTypesCapabilityTypePut(ctx, capabilityType).DefaultModelDeploymentUpdate(defaultModelDeploymentUpdate).Execute()

Update Capability Type



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
	capabilityType := "capabilityType_example" // string | 
	defaultModelDeploymentUpdate := *openapiclient.NewDefaultModelDeploymentUpdate("DefaultModelDeploymentId_example") // DefaultModelDeploymentUpdate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilityTypesAPI.UpdateCapabilityTypeV1CapabilityTypesCapabilityTypePut(context.Background(), capabilityType).DefaultModelDeploymentUpdate(defaultModelDeploymentUpdate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilityTypesAPI.UpdateCapabilityTypeV1CapabilityTypesCapabilityTypePut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateCapabilityTypeV1CapabilityTypesCapabilityTypePut`: CapabilityTypeRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CapabilityTypesAPI.UpdateCapabilityTypeV1CapabilityTypesCapabilityTypePut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**capabilityType** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateCapabilityTypeV1CapabilityTypesCapabilityTypePutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **defaultModelDeploymentUpdate** | [**DefaultModelDeploymentUpdate**](DefaultModelDeploymentUpdate.md) |  | 

### Return type

[**CapabilityTypeRepresentation**](CapabilityTypeRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

