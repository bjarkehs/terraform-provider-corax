# \APIKeysAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateApiKeyV1ApiKeysPost**](APIKeysAPI.md#CreateApiKeyV1ApiKeysPost) | **Post** /v1/api-keys | Create Api Key
[**DeleteApiKeyV1ApiKeysKeyIdDelete**](APIKeysAPI.md#DeleteApiKeyV1ApiKeysKeyIdDelete) | **Delete** /v1/api-keys/{key_id} | Delete Api Key
[**GetApiKeysV1ApiKeysGet**](APIKeysAPI.md#GetApiKeysV1ApiKeysGet) | **Get** /v1/api-keys | Get Api Keys



## CreateApiKeyV1ApiKeysPost

> ApiKey CreateApiKeyV1ApiKeysPost(ctx).ApiKeyCreate(apiKeyCreate).Execute()

Create Api Key



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
    "time"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/api"
)

func main() {
	apiKeyCreate := *openapiclient.NewApiKeyCreate("Name_example", time.Now()) // ApiKeyCreate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.APIKeysAPI.CreateApiKeyV1ApiKeysPost(context.Background()).ApiKeyCreate(apiKeyCreate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `APIKeysAPI.CreateApiKeyV1ApiKeysPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateApiKeyV1ApiKeysPost`: ApiKey
	fmt.Fprintf(os.Stdout, "Response from `APIKeysAPI.CreateApiKeyV1ApiKeysPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateApiKeyV1ApiKeysPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **apiKeyCreate** | [**ApiKeyCreate**](ApiKeyCreate.md) |  | 

### Return type

[**ApiKey**](ApiKey.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteApiKeyV1ApiKeysKeyIdDelete

> DeleteApiKeyV1ApiKeysKeyIdDelete(ctx, keyId).Execute()

Delete Api Key



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
	keyId := "keyId_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.APIKeysAPI.DeleteApiKeyV1ApiKeysKeyIdDelete(context.Background(), keyId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `APIKeysAPI.DeleteApiKeyV1ApiKeysKeyIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**keyId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteApiKeyV1ApiKeysKeyIdDeleteRequest struct via the builder pattern


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


## GetApiKeysV1ApiKeysGet

> PagedResponseModelApiKey GetApiKeysV1ApiKeysGet(ctx).Page(page).Size(size).Sort(sort).Filter(filter).Execute()

Get Api Keys



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
	sort := "sort_example" // string | Field to sort by, prefix with '-' for descending order (optional) (default to "id")
	filter := "filter_example" // string | Fields to filter by. Format 'key::value', use * for wildcards, use | to apply multiple filters (AND) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.APIKeysAPI.GetApiKeysV1ApiKeysGet(context.Background()).Page(page).Size(size).Sort(sort).Filter(filter).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `APIKeysAPI.GetApiKeysV1ApiKeysGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetApiKeysV1ApiKeysGet`: PagedResponseModelApiKey
	fmt.Fprintf(os.Stdout, "Response from `APIKeysAPI.GetApiKeysV1ApiKeysGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetApiKeysV1ApiKeysGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** |  | [default to 1]
 **size** | **int32** |  | [default to 10]
 **sort** | **string** | Field to sort by, prefix with &#39;-&#39; for descending order | [default to &quot;id&quot;]
 **filter** | **string** | Fields to filter by. Format &#39;key::value&#39;, use * for wildcards, use | to apply multiple filters (AND) | 

### Return type

[**PagedResponseModelApiKey**](PagedResponseModelApiKey.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

