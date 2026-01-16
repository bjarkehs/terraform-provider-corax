# \ModelProvidersAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateModelProviderV1ModelProvidersPost**](ModelProvidersAPI.md#CreateModelProviderV1ModelProvidersPost) | **Post** /v1/model-providers | Create Model Provider
[**DeleteModelProviderV1ModelProvidersProviderIdDelete**](ModelProvidersAPI.md#DeleteModelProviderV1ModelProvidersProviderIdDelete) | **Delete** /v1/model-providers/{provider_id} | Delete Model Provider
[**GetModelProviderV1ModelProvidersProviderIdGet**](ModelProvidersAPI.md#GetModelProviderV1ModelProvidersProviderIdGet) | **Get** /v1/model-providers/{provider_id} | Get Model Provider
[**ListModelProvidersV1ModelProvidersGet**](ModelProvidersAPI.md#ListModelProvidersV1ModelProvidersGet) | **Get** /v1/model-providers | List Model Providers
[**UpdateModelProviderV1ModelProvidersProviderIdPut**](ModelProvidersAPI.md#UpdateModelProviderV1ModelProvidersProviderIdPut) | **Put** /v1/model-providers/{provider_id} | Update Model Provider



## CreateModelProviderV1ModelProvidersPost

> ModelProvider CreateModelProviderV1ModelProvidersPost(ctx).ModelProviderCreate(modelProviderCreate).Execute()

Create Model Provider



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
	modelProviderCreate := *openapiclient.NewModelProviderCreate("Name_example", "ProviderType_example", map[string]interface{}{"key": interface{}(123)}) // ModelProviderCreate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ModelProvidersAPI.CreateModelProviderV1ModelProvidersPost(context.Background()).ModelProviderCreate(modelProviderCreate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ModelProvidersAPI.CreateModelProviderV1ModelProvidersPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateModelProviderV1ModelProvidersPost`: ModelProvider
	fmt.Fprintf(os.Stdout, "Response from `ModelProvidersAPI.CreateModelProviderV1ModelProvidersPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateModelProviderV1ModelProvidersPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **modelProviderCreate** | [**ModelProviderCreate**](ModelProviderCreate.md) |  | 

### Return type

[**ModelProvider**](ModelProvider.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteModelProviderV1ModelProvidersProviderIdDelete

> DeleteModelProviderV1ModelProvidersProviderIdDelete(ctx, providerId).Execute()

Delete Model Provider



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
	providerId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ModelProvidersAPI.DeleteModelProviderV1ModelProvidersProviderIdDelete(context.Background(), providerId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ModelProvidersAPI.DeleteModelProviderV1ModelProvidersProviderIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**providerId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteModelProviderV1ModelProvidersProviderIdDeleteRequest struct via the builder pattern


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


## GetModelProviderV1ModelProvidersProviderIdGet

> ModelProvider GetModelProviderV1ModelProvidersProviderIdGet(ctx, providerId).Execute()

Get Model Provider



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
	providerId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ModelProvidersAPI.GetModelProviderV1ModelProvidersProviderIdGet(context.Background(), providerId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ModelProvidersAPI.GetModelProviderV1ModelProvidersProviderIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetModelProviderV1ModelProvidersProviderIdGet`: ModelProvider
	fmt.Fprintf(os.Stdout, "Response from `ModelProvidersAPI.GetModelProviderV1ModelProvidersProviderIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**providerId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetModelProviderV1ModelProvidersProviderIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ModelProvider**](ModelProvider.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListModelProvidersV1ModelProvidersGet

> PagedResponseModelModelProvider ListModelProvidersV1ModelProvidersGet(ctx).Page(page).Size(size).Sort(sort).Filter(filter).Execute()

List Model Providers



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
	resp, r, err := apiClient.ModelProvidersAPI.ListModelProvidersV1ModelProvidersGet(context.Background()).Page(page).Size(size).Sort(sort).Filter(filter).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ModelProvidersAPI.ListModelProvidersV1ModelProvidersGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListModelProvidersV1ModelProvidersGet`: PagedResponseModelModelProvider
	fmt.Fprintf(os.Stdout, "Response from `ModelProvidersAPI.ListModelProvidersV1ModelProvidersGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListModelProvidersV1ModelProvidersGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** |  | [default to 1]
 **size** | **int32** |  | [default to 10]
 **sort** | **string** |  | [default to &quot;id&quot;]
 **filter** | **string** |  | 

### Return type

[**PagedResponseModelModelProvider**](PagedResponseModelModelProvider.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateModelProviderV1ModelProvidersProviderIdPut

> ModelProvider UpdateModelProviderV1ModelProvidersProviderIdPut(ctx, providerId).ModelProviderUpdate(modelProviderUpdate).Execute()

Update Model Provider



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
	providerId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	modelProviderUpdate := *openapiclient.NewModelProviderUpdate("Name_example", "ProviderType_example", map[string]interface{}{"key": interface{}(123)}, "Id_example") // ModelProviderUpdate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ModelProvidersAPI.UpdateModelProviderV1ModelProvidersProviderIdPut(context.Background(), providerId).ModelProviderUpdate(modelProviderUpdate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ModelProvidersAPI.UpdateModelProviderV1ModelProvidersProviderIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateModelProviderV1ModelProvidersProviderIdPut`: ModelProvider
	fmt.Fprintf(os.Stdout, "Response from `ModelProvidersAPI.UpdateModelProviderV1ModelProvidersProviderIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**providerId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateModelProviderV1ModelProvidersProviderIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **modelProviderUpdate** | [**ModelProviderUpdate**](ModelProviderUpdate.md) |  | 

### Return type

[**ModelProvider**](ModelProvider.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

