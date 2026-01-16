# \ModelProviderTypesAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetModelProviderConfigurationV1ModelProviderTypesProviderTypeModelProviderConfigurationGet**](ModelProviderTypesAPI.md#GetModelProviderConfigurationV1ModelProviderTypesProviderTypeModelProviderConfigurationGet) | **Get** /v1/model-provider-types/{provider_type}/model-provider-configuration | Get Model Provider Configuration
[**GetModelProviderDeploymentConfigurationV1ModelProviderTypesProviderTypeModelDeploymentConfigurationGet**](ModelProviderTypesAPI.md#GetModelProviderDeploymentConfigurationV1ModelProviderTypesProviderTypeModelDeploymentConfigurationGet) | **Get** /v1/model-provider-types/{provider_type}/model-deployment-configuration | Get Model Provider Deployment Configuration
[**GetModelProviderTypeV1ModelProviderTypesProviderTypeGet**](ModelProviderTypesAPI.md#GetModelProviderTypeV1ModelProviderTypesProviderTypeGet) | **Get** /v1/model-provider-types/{provider_type} | Get Model Provider Type
[**ListModelProviderTypesV1ModelProviderTypesGet**](ModelProviderTypesAPI.md#ListModelProviderTypesV1ModelProviderTypesGet) | **Get** /v1/model-provider-types | List Model Provider Types



## GetModelProviderConfigurationV1ModelProviderTypesProviderTypeModelProviderConfigurationGet

> Configuration GetModelProviderConfigurationV1ModelProviderTypesProviderTypeModelProviderConfigurationGet(ctx, providerType).Execute()

Get Model Provider Configuration



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
	providerType := "providerType_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ModelProviderTypesAPI.GetModelProviderConfigurationV1ModelProviderTypesProviderTypeModelProviderConfigurationGet(context.Background(), providerType).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ModelProviderTypesAPI.GetModelProviderConfigurationV1ModelProviderTypesProviderTypeModelProviderConfigurationGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetModelProviderConfigurationV1ModelProviderTypesProviderTypeModelProviderConfigurationGet`: Configuration
	fmt.Fprintf(os.Stdout, "Response from `ModelProviderTypesAPI.GetModelProviderConfigurationV1ModelProviderTypesProviderTypeModelProviderConfigurationGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**providerType** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetModelProviderConfigurationV1ModelProviderTypesProviderTypeModelProviderConfigurationGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Configuration**](Configuration.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetModelProviderDeploymentConfigurationV1ModelProviderTypesProviderTypeModelDeploymentConfigurationGet

> Configuration GetModelProviderDeploymentConfigurationV1ModelProviderTypesProviderTypeModelDeploymentConfigurationGet(ctx, providerType).Execute()

Get Model Provider Deployment Configuration



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
	providerType := "providerType_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ModelProviderTypesAPI.GetModelProviderDeploymentConfigurationV1ModelProviderTypesProviderTypeModelDeploymentConfigurationGet(context.Background(), providerType).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ModelProviderTypesAPI.GetModelProviderDeploymentConfigurationV1ModelProviderTypesProviderTypeModelDeploymentConfigurationGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetModelProviderDeploymentConfigurationV1ModelProviderTypesProviderTypeModelDeploymentConfigurationGet`: Configuration
	fmt.Fprintf(os.Stdout, "Response from `ModelProviderTypesAPI.GetModelProviderDeploymentConfigurationV1ModelProviderTypesProviderTypeModelDeploymentConfigurationGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**providerType** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetModelProviderDeploymentConfigurationV1ModelProviderTypesProviderTypeModelDeploymentConfigurationGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Configuration**](Configuration.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetModelProviderTypeV1ModelProviderTypesProviderTypeGet

> ModelProviderType GetModelProviderTypeV1ModelProviderTypesProviderTypeGet(ctx, providerType).Execute()

Get Model Provider Type



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
	providerType := "providerType_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ModelProviderTypesAPI.GetModelProviderTypeV1ModelProviderTypesProviderTypeGet(context.Background(), providerType).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ModelProviderTypesAPI.GetModelProviderTypeV1ModelProviderTypesProviderTypeGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetModelProviderTypeV1ModelProviderTypesProviderTypeGet`: ModelProviderType
	fmt.Fprintf(os.Stdout, "Response from `ModelProviderTypesAPI.GetModelProviderTypeV1ModelProviderTypesProviderTypeGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**providerType** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetModelProviderTypeV1ModelProviderTypesProviderTypeGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ModelProviderType**](ModelProviderType.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListModelProviderTypesV1ModelProviderTypesGet

> PagedResponseModelModelProviderType ListModelProviderTypesV1ModelProviderTypesGet(ctx).Page(page).Size(size).Sort(sort).Filter(filter).Execute()

List Model Provider Types



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
	size := int32(56) // int32 |  (optional) (default to 100)
	sort := "sort_example" // string |  (optional) (default to "label")
	filter := "filter_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ModelProviderTypesAPI.ListModelProviderTypesV1ModelProviderTypesGet(context.Background()).Page(page).Size(size).Sort(sort).Filter(filter).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ModelProviderTypesAPI.ListModelProviderTypesV1ModelProviderTypesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListModelProviderTypesV1ModelProviderTypesGet`: PagedResponseModelModelProviderType
	fmt.Fprintf(os.Stdout, "Response from `ModelProviderTypesAPI.ListModelProviderTypesV1ModelProviderTypesGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListModelProviderTypesV1ModelProviderTypesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** |  | [default to 1]
 **size** | **int32** |  | [default to 100]
 **sort** | **string** |  | [default to &quot;label&quot;]
 **filter** | **string** |  | 

### Return type

[**PagedResponseModelModelProviderType**](PagedResponseModelModelProviderType.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

