# \CapabilitiesAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateCapabilityV1CapabilitiesPost**](CapabilitiesAPI.md#CreateCapabilityV1CapabilitiesPost) | **Post** /v1/capabilities | Create Capability
[**DeleteCapabilityV1CapabilitiesCapabilityIdDelete**](CapabilitiesAPI.md#DeleteCapabilityV1CapabilitiesCapabilityIdDelete) | **Delete** /v1/capabilities/{capability_id} | Delete Capability
[**ExecuteCapabilityV1CapabilitiesCapabilityIdExecutionsPost**](CapabilitiesAPI.md#ExecuteCapabilityV1CapabilitiesCapabilityIdExecutionsPost) | **Post** /v1/capabilities/{capability_id}/executions | Execute Capability
[**GetCapabilityVersionV1CapabilitiesCapabilityIdVersionsVersionGet**](CapabilitiesAPI.md#GetCapabilityVersionV1CapabilitiesCapabilityIdVersionsVersionGet) | **Get** /v1/capabilities/{capability_id}/versions/{version} | Get Capability Version
[**GetExecutionResultV1CapabilitiesCapabilityIdExecutionsExecutionIdResultGet**](CapabilitiesAPI.md#GetExecutionResultV1CapabilitiesCapabilityIdExecutionsExecutionIdResultGet) | **Get** /v1/capabilities/{capability_id}/executions/{execution_id}/result | Get Execution Result
[**GetExecutionUsageV1CapabilitiesCapabilityIdExecutionsExecutionIdUsageGet**](CapabilitiesAPI.md#GetExecutionUsageV1CapabilitiesCapabilityIdExecutionsExecutionIdUsageGet) | **Get** /v1/capabilities/{capability_id}/executions/{execution_id}/usage | Get Execution Usage
[**GetExecutionV1CapabilitiesCapabilityIdExecutionsExecutionIdGet**](CapabilitiesAPI.md#GetExecutionV1CapabilitiesCapabilityIdExecutionsExecutionIdGet) | **Get** /v1/capabilities/{capability_id}/executions/{execution_id} | Get Execution
[**ListCapabilitiesV1CapabilitiesGet**](CapabilitiesAPI.md#ListCapabilitiesV1CapabilitiesGet) | **Get** /v1/capabilities | List Capabilities
[**ListCapabilityVersionsV1CapabilitiesCapabilityIdVersionsGet**](CapabilitiesAPI.md#ListCapabilityVersionsV1CapabilitiesCapabilityIdVersionsGet) | **Get** /v1/capabilities/{capability_id}/versions | List Capability Versions
[**ListExecutionsV1CapabilitiesCapabilityIdExecutionsGet**](CapabilitiesAPI.md#ListExecutionsV1CapabilitiesCapabilityIdExecutionsGet) | **Get** /v1/capabilities/{capability_id}/executions | List Executions
[**ProvideFeedbackV1CapabilitiesCapabilityIdExecutionsExecutionIdFeedbackPost**](CapabilitiesAPI.md#ProvideFeedbackV1CapabilitiesCapabilityIdExecutionsExecutionIdFeedbackPost) | **Post** /v1/capabilities/{capability_id}/executions/{execution_id}/feedback | Provide Feedback
[**ReadCapabilityV1CapabilitiesCapabilityIdGet**](CapabilitiesAPI.md#ReadCapabilityV1CapabilitiesCapabilityIdGet) | **Get** /v1/capabilities/{capability_id} | Read Capability
[**SetDefaultCapabilityVersionV1CapabilitiesCapabilityIdDefaultVersionPut**](CapabilitiesAPI.md#SetDefaultCapabilityVersionV1CapabilitiesCapabilityIdDefaultVersionPut) | **Put** /v1/capabilities/{capability_id}/default-version | Set Default Capability Version
[**StreamExecutionEventsV1CapabilitiesCapabilityIdExecutionsExecutionIdStreamGet**](CapabilitiesAPI.md#StreamExecutionEventsV1CapabilitiesCapabilityIdExecutionsExecutionIdStreamGet) | **Get** /v1/capabilities/{capability_id}/executions/{execution_id}/stream | Stream Execution Events
[**UpdateCapabilityV1CapabilitiesCapabilityIdPut**](CapabilitiesAPI.md#UpdateCapabilityV1CapabilitiesCapabilityIdPut) | **Put** /v1/capabilities/{capability_id} | Update Capability



## CreateCapabilityV1CapabilitiesPost

> ResponseCreateCapabilityV1CapabilitiesPost CreateCapabilityV1CapabilitiesPost(ctx).Capability1(capability1).Execute()

Create Capability



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
	capability1 := *openapiclient.NewCapability1("Name_example", "Type_example", "SystemPrompt_example", "CompletionPrompt_example", "OutputType_example") // Capability1 | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilitiesAPI.CreateCapabilityV1CapabilitiesPost(context.Background()).Capability1(capability1).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilitiesAPI.CreateCapabilityV1CapabilitiesPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateCapabilityV1CapabilitiesPost`: ResponseCreateCapabilityV1CapabilitiesPost
	fmt.Fprintf(os.Stdout, "Response from `CapabilitiesAPI.CreateCapabilityV1CapabilitiesPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateCapabilityV1CapabilitiesPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **capability1** | [**Capability1**](Capability1.md) |  | 

### Return type

[**ResponseCreateCapabilityV1CapabilitiesPost**](ResponseCreateCapabilityV1CapabilitiesPost.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteCapabilityV1CapabilitiesCapabilityIdDelete

> DeleteCapabilityV1CapabilitiesCapabilityIdDelete(ctx, capabilityId).Execute()

Delete Capability



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
	capabilityId := *openapiclient.NewCapabilityId1() // CapabilityId1 | The ID or semantic ID of the capability

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.CapabilitiesAPI.DeleteCapabilityV1CapabilitiesCapabilityIdDelete(context.Background(), capabilityId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilitiesAPI.DeleteCapabilityV1CapabilitiesCapabilityIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**capabilityId** | [**CapabilityId1**](.md) | The ID or semantic ID of the capability | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteCapabilityV1CapabilitiesCapabilityIdDeleteRequest struct via the builder pattern


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


## ExecuteCapabilityV1CapabilitiesCapabilityIdExecutionsPost

> Execution ExecuteCapabilityV1CapabilitiesCapabilityIdExecutionsPost(ctx, capabilityId).ExecutionCreate(executionCreate).Stream(stream).StreamTransport(streamTransport).Execute()

Execute Capability



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
	capabilityId := *openapiclient.NewCapabilityId1() // CapabilityId1 | The ID or semantic ID of the capability
	executionCreate := *openapiclient.NewExecutionCreate() // ExecutionCreate | 
	stream := true // bool | Stream incremental chat capability responses. (optional) (default to false)
	streamTransport := openapiclient.StreamingTransport("streamablehttp") // StreamingTransport | Output transport when streaming. Use 'streamablehttp' for newline-delimited JSON chunks or 'sse' for legacy Server-Sent Events. (optional) (default to "streamablehttp")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilitiesAPI.ExecuteCapabilityV1CapabilitiesCapabilityIdExecutionsPost(context.Background(), capabilityId).ExecutionCreate(executionCreate).Stream(stream).StreamTransport(streamTransport).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilitiesAPI.ExecuteCapabilityV1CapabilitiesCapabilityIdExecutionsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExecuteCapabilityV1CapabilitiesCapabilityIdExecutionsPost`: Execution
	fmt.Fprintf(os.Stdout, "Response from `CapabilitiesAPI.ExecuteCapabilityV1CapabilitiesCapabilityIdExecutionsPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**capabilityId** | [**CapabilityId1**](.md) | The ID or semantic ID of the capability | 

### Other Parameters

Other parameters are passed through a pointer to a apiExecuteCapabilityV1CapabilitiesCapabilityIdExecutionsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **executionCreate** | [**ExecutionCreate**](ExecutionCreate.md) |  | 
 **stream** | **bool** | Stream incremental chat capability responses. | [default to false]
 **streamTransport** | [**StreamingTransport**](StreamingTransport.md) | Output transport when streaming. Use &#39;streamablehttp&#39; for newline-delimited JSON chunks or &#39;sse&#39; for legacy Server-Sent Events. | [default to &quot;streamablehttp&quot;]

### Return type

[**Execution**](Execution.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetCapabilityVersionV1CapabilitiesCapabilityIdVersionsVersionGet

> CapabilityRepresentation GetCapabilityVersionV1CapabilitiesCapabilityIdVersionsVersionGet(ctx, version, capabilityId).Execute()

Get Capability Version



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
	version := int32(56) // int32 | 
	capabilityId := *openapiclient.NewCapabilityId1() // CapabilityId1 | The ID or semantic ID of the capability

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilitiesAPI.GetCapabilityVersionV1CapabilitiesCapabilityIdVersionsVersionGet(context.Background(), version, capabilityId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilitiesAPI.GetCapabilityVersionV1CapabilitiesCapabilityIdVersionsVersionGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetCapabilityVersionV1CapabilitiesCapabilityIdVersionsVersionGet`: CapabilityRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CapabilitiesAPI.GetCapabilityVersionV1CapabilitiesCapabilityIdVersionsVersionGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**version** | **int32** |  | 
**capabilityId** | [**CapabilityId1**](.md) | The ID or semantic ID of the capability | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetCapabilityVersionV1CapabilitiesCapabilityIdVersionsVersionGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**CapabilityRepresentation**](CapabilityRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetExecutionResultV1CapabilitiesCapabilityIdExecutionsExecutionIdResultGet

> ExecutionResult GetExecutionResultV1CapabilitiesCapabilityIdExecutionsExecutionIdResultGet(ctx, executionId, capabilityId).Execute()

Get Execution Result



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
	executionId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	capabilityId := *openapiclient.NewCapabilityId1() // CapabilityId1 | The ID or semantic ID of the capability

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilitiesAPI.GetExecutionResultV1CapabilitiesCapabilityIdExecutionsExecutionIdResultGet(context.Background(), executionId, capabilityId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilitiesAPI.GetExecutionResultV1CapabilitiesCapabilityIdExecutionsExecutionIdResultGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetExecutionResultV1CapabilitiesCapabilityIdExecutionsExecutionIdResultGet`: ExecutionResult
	fmt.Fprintf(os.Stdout, "Response from `CapabilitiesAPI.GetExecutionResultV1CapabilitiesCapabilityIdExecutionsExecutionIdResultGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**executionId** | **string** |  | 
**capabilityId** | [**CapabilityId1**](.md) | The ID or semantic ID of the capability | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetExecutionResultV1CapabilitiesCapabilityIdExecutionsExecutionIdResultGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**ExecutionResult**](ExecutionResult.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetExecutionUsageV1CapabilitiesCapabilityIdExecutionsExecutionIdUsageGet

> ExecutionUsage GetExecutionUsageV1CapabilitiesCapabilityIdExecutionsExecutionIdUsageGet(ctx, executionId, capabilityId).Execute()

Get Execution Usage



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
	executionId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	capabilityId := *openapiclient.NewCapabilityId1() // CapabilityId1 | The ID or semantic ID of the capability

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilitiesAPI.GetExecutionUsageV1CapabilitiesCapabilityIdExecutionsExecutionIdUsageGet(context.Background(), executionId, capabilityId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilitiesAPI.GetExecutionUsageV1CapabilitiesCapabilityIdExecutionsExecutionIdUsageGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetExecutionUsageV1CapabilitiesCapabilityIdExecutionsExecutionIdUsageGet`: ExecutionUsage
	fmt.Fprintf(os.Stdout, "Response from `CapabilitiesAPI.GetExecutionUsageV1CapabilitiesCapabilityIdExecutionsExecutionIdUsageGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**executionId** | **string** |  | 
**capabilityId** | [**CapabilityId1**](.md) | The ID or semantic ID of the capability | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetExecutionUsageV1CapabilitiesCapabilityIdExecutionsExecutionIdUsageGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**ExecutionUsage**](ExecutionUsage.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetExecutionV1CapabilitiesCapabilityIdExecutionsExecutionIdGet

> Execution GetExecutionV1CapabilitiesCapabilityIdExecutionsExecutionIdGet(ctx, executionId, capabilityId).Execute()

Get Execution



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
	executionId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	capabilityId := *openapiclient.NewCapabilityId1() // CapabilityId1 | The ID or semantic ID of the capability

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilitiesAPI.GetExecutionV1CapabilitiesCapabilityIdExecutionsExecutionIdGet(context.Background(), executionId, capabilityId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilitiesAPI.GetExecutionV1CapabilitiesCapabilityIdExecutionsExecutionIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetExecutionV1CapabilitiesCapabilityIdExecutionsExecutionIdGet`: Execution
	fmt.Fprintf(os.Stdout, "Response from `CapabilitiesAPI.GetExecutionV1CapabilitiesCapabilityIdExecutionsExecutionIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**executionId** | **string** |  | 
**capabilityId** | [**CapabilityId1**](.md) | The ID or semantic ID of the capability | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetExecutionV1CapabilitiesCapabilityIdExecutionsExecutionIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**Execution**](Execution.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListCapabilitiesV1CapabilitiesGet

> PagedResponseModelCapability ListCapabilitiesV1CapabilitiesGet(ctx).Page(page).Size(size).Sort(sort).Filter(filter).Execute()

List Capabilities



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
	resp, r, err := apiClient.CapabilitiesAPI.ListCapabilitiesV1CapabilitiesGet(context.Background()).Page(page).Size(size).Sort(sort).Filter(filter).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilitiesAPI.ListCapabilitiesV1CapabilitiesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListCapabilitiesV1CapabilitiesGet`: PagedResponseModelCapability
	fmt.Fprintf(os.Stdout, "Response from `CapabilitiesAPI.ListCapabilitiesV1CapabilitiesGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListCapabilitiesV1CapabilitiesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** |  | [default to 1]
 **size** | **int32** |  | [default to 10]
 **sort** | **string** |  | [default to &quot;id&quot;]
 **filter** | **string** |  | 

### Return type

[**PagedResponseModelCapability**](PagedResponseModelCapability.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListCapabilityVersionsV1CapabilitiesCapabilityIdVersionsGet

> PagedResponseModelCapabilityVersion ListCapabilityVersionsV1CapabilitiesCapabilityIdVersionsGet(ctx, capabilityId).Page(page).Size(size).Sort(sort).Filter(filter).Execute()

List Capability Versions



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
	capabilityId := *openapiclient.NewCapabilityId1() // CapabilityId1 | The ID or semantic ID of the capability
	page := int32(56) // int32 |  (optional) (default to 1)
	size := int32(56) // int32 |  (optional) (default to 10)
	sort := "sort_example" // string |  (optional) (default to "id")
	filter := "filter_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilitiesAPI.ListCapabilityVersionsV1CapabilitiesCapabilityIdVersionsGet(context.Background(), capabilityId).Page(page).Size(size).Sort(sort).Filter(filter).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilitiesAPI.ListCapabilityVersionsV1CapabilitiesCapabilityIdVersionsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListCapabilityVersionsV1CapabilitiesCapabilityIdVersionsGet`: PagedResponseModelCapabilityVersion
	fmt.Fprintf(os.Stdout, "Response from `CapabilitiesAPI.ListCapabilityVersionsV1CapabilitiesCapabilityIdVersionsGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**capabilityId** | [**CapabilityId1**](.md) | The ID or semantic ID of the capability | 

### Other Parameters

Other parameters are passed through a pointer to a apiListCapabilityVersionsV1CapabilitiesCapabilityIdVersionsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **page** | **int32** |  | [default to 1]
 **size** | **int32** |  | [default to 10]
 **sort** | **string** |  | [default to &quot;id&quot;]
 **filter** | **string** |  | 

### Return type

[**PagedResponseModelCapabilityVersion**](PagedResponseModelCapabilityVersion.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListExecutionsV1CapabilitiesCapabilityIdExecutionsGet

> PagedResponseModelExecution ListExecutionsV1CapabilitiesCapabilityIdExecutionsGet(ctx, capabilityId).Page(page).Size(size).Sort(sort).Filter(filter).Execute()

List Executions



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
	capabilityId := *openapiclient.NewCapabilityId1() // CapabilityId1 | The ID or semantic ID of the capability
	page := int32(56) // int32 |  (optional) (default to 1)
	size := int32(56) // int32 |  (optional) (default to 10)
	sort := "sort_example" // string |  (optional) (default to "id")
	filter := "filter_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilitiesAPI.ListExecutionsV1CapabilitiesCapabilityIdExecutionsGet(context.Background(), capabilityId).Page(page).Size(size).Sort(sort).Filter(filter).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilitiesAPI.ListExecutionsV1CapabilitiesCapabilityIdExecutionsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListExecutionsV1CapabilitiesCapabilityIdExecutionsGet`: PagedResponseModelExecution
	fmt.Fprintf(os.Stdout, "Response from `CapabilitiesAPI.ListExecutionsV1CapabilitiesCapabilityIdExecutionsGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**capabilityId** | [**CapabilityId1**](.md) | The ID or semantic ID of the capability | 

### Other Parameters

Other parameters are passed through a pointer to a apiListExecutionsV1CapabilitiesCapabilityIdExecutionsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **page** | **int32** |  | [default to 1]
 **size** | **int32** |  | [default to 10]
 **sort** | **string** |  | [default to &quot;id&quot;]
 **filter** | **string** |  | 

### Return type

[**PagedResponseModelExecution**](PagedResponseModelExecution.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProvideFeedbackV1CapabilitiesCapabilityIdExecutionsExecutionIdFeedbackPost

> ExecutionResult ProvideFeedbackV1CapabilitiesCapabilityIdExecutionsExecutionIdFeedbackPost(ctx, executionId, capabilityId).FeedbackCreate(feedbackCreate).Execute()

Provide Feedback



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
	executionId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	capabilityId := *openapiclient.NewCapabilityId1() // CapabilityId1 | The ID or semantic ID of the capability
	feedbackCreate := *openapiclient.NewFeedbackCreate("Feedback_example") // FeedbackCreate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilitiesAPI.ProvideFeedbackV1CapabilitiesCapabilityIdExecutionsExecutionIdFeedbackPost(context.Background(), executionId, capabilityId).FeedbackCreate(feedbackCreate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilitiesAPI.ProvideFeedbackV1CapabilitiesCapabilityIdExecutionsExecutionIdFeedbackPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ProvideFeedbackV1CapabilitiesCapabilityIdExecutionsExecutionIdFeedbackPost`: ExecutionResult
	fmt.Fprintf(os.Stdout, "Response from `CapabilitiesAPI.ProvideFeedbackV1CapabilitiesCapabilityIdExecutionsExecutionIdFeedbackPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**executionId** | **string** |  | 
**capabilityId** | [**CapabilityId1**](.md) | The ID or semantic ID of the capability | 

### Other Parameters

Other parameters are passed through a pointer to a apiProvideFeedbackV1CapabilitiesCapabilityIdExecutionsExecutionIdFeedbackPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **feedbackCreate** | [**FeedbackCreate**](FeedbackCreate.md) |  | 

### Return type

[**ExecutionResult**](ExecutionResult.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReadCapabilityV1CapabilitiesCapabilityIdGet

> CapabilityRepresentation ReadCapabilityV1CapabilitiesCapabilityIdGet(ctx, capabilityId).Execute()

Read Capability



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
	capabilityId := *openapiclient.NewCapabilityId() // CapabilityId | The ID or semantic ID of the capability

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilitiesAPI.ReadCapabilityV1CapabilitiesCapabilityIdGet(context.Background(), capabilityId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilitiesAPI.ReadCapabilityV1CapabilitiesCapabilityIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReadCapabilityV1CapabilitiesCapabilityIdGet`: CapabilityRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CapabilitiesAPI.ReadCapabilityV1CapabilitiesCapabilityIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**capabilityId** | [**CapabilityId**](.md) | The ID or semantic ID of the capability | 

### Other Parameters

Other parameters are passed through a pointer to a apiReadCapabilityV1CapabilitiesCapabilityIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CapabilityRepresentation**](CapabilityRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SetDefaultCapabilityVersionV1CapabilitiesCapabilityIdDefaultVersionPut

> CapabilityRepresentation SetDefaultCapabilityVersionV1CapabilitiesCapabilityIdDefaultVersionPut(ctx, capabilityId).CapabilitySetDefaultVersionRequest(capabilitySetDefaultVersionRequest).Execute()

Set Default Capability Version



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
	capabilityId := *openapiclient.NewCapabilityId1() // CapabilityId1 | The ID or semantic ID of the capability
	capabilitySetDefaultVersionRequest := *openapiclient.NewCapabilitySetDefaultVersionRequest(int32(123)) // CapabilitySetDefaultVersionRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilitiesAPI.SetDefaultCapabilityVersionV1CapabilitiesCapabilityIdDefaultVersionPut(context.Background(), capabilityId).CapabilitySetDefaultVersionRequest(capabilitySetDefaultVersionRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilitiesAPI.SetDefaultCapabilityVersionV1CapabilitiesCapabilityIdDefaultVersionPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SetDefaultCapabilityVersionV1CapabilitiesCapabilityIdDefaultVersionPut`: CapabilityRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CapabilitiesAPI.SetDefaultCapabilityVersionV1CapabilitiesCapabilityIdDefaultVersionPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**capabilityId** | [**CapabilityId1**](.md) | The ID or semantic ID of the capability | 

### Other Parameters

Other parameters are passed through a pointer to a apiSetDefaultCapabilityVersionV1CapabilitiesCapabilityIdDefaultVersionPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **capabilitySetDefaultVersionRequest** | [**CapabilitySetDefaultVersionRequest**](CapabilitySetDefaultVersionRequest.md) |  | 

### Return type

[**CapabilityRepresentation**](CapabilityRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## StreamExecutionEventsV1CapabilitiesCapabilityIdExecutionsExecutionIdStreamGet

> interface{} StreamExecutionEventsV1CapabilitiesCapabilityIdExecutionsExecutionIdStreamGet(ctx, executionId, capabilityId).Transport(transport).Execute()

Stream Execution Events



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
	executionId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	capabilityId := *openapiclient.NewCapabilityId1() // CapabilityId1 | The ID or semantic ID of the capability
	transport := openapiclient.StreamingTransport("streamablehttp") // StreamingTransport | Preferred streaming transport. Use 'streamablehttp' for newline-delimited JSON chunks or 'sse' for Server-Sent Events. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilitiesAPI.StreamExecutionEventsV1CapabilitiesCapabilityIdExecutionsExecutionIdStreamGet(context.Background(), executionId, capabilityId).Transport(transport).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilitiesAPI.StreamExecutionEventsV1CapabilitiesCapabilityIdExecutionsExecutionIdStreamGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `StreamExecutionEventsV1CapabilitiesCapabilityIdExecutionsExecutionIdStreamGet`: interface{}
	fmt.Fprintf(os.Stdout, "Response from `CapabilitiesAPI.StreamExecutionEventsV1CapabilitiesCapabilityIdExecutionsExecutionIdStreamGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**executionId** | **string** |  | 
**capabilityId** | [**CapabilityId1**](.md) | The ID or semantic ID of the capability | 

### Other Parameters

Other parameters are passed through a pointer to a apiStreamExecutionEventsV1CapabilitiesCapabilityIdExecutionsExecutionIdStreamGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **transport** | [**StreamingTransport**](StreamingTransport.md) | Preferred streaming transport. Use &#39;streamablehttp&#39; for newline-delimited JSON chunks or &#39;sse&#39; for Server-Sent Events. | 

### Return type

**interface{}**

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateCapabilityV1CapabilitiesCapabilityIdPut

> CapabilityRepresentation UpdateCapabilityV1CapabilitiesCapabilityIdPut(ctx, capabilityId).Capability2(capability2).Execute()

Update Capability



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
	capabilityId := *openapiclient.NewCapabilityId1() // CapabilityId1 | The ID or semantic ID of the capability
	capability2 := *openapiclient.NewCapability2("Name_example", "Type_example") // Capability2 | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilitiesAPI.UpdateCapabilityV1CapabilitiesCapabilityIdPut(context.Background(), capabilityId).Capability2(capability2).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilitiesAPI.UpdateCapabilityV1CapabilitiesCapabilityIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateCapabilityV1CapabilitiesCapabilityIdPut`: CapabilityRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CapabilitiesAPI.UpdateCapabilityV1CapabilitiesCapabilityIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**capabilityId** | [**CapabilityId1**](.md) | The ID or semantic ID of the capability | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateCapabilityV1CapabilitiesCapabilityIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **capability2** | [**Capability2**](Capability2.md) |  | 

### Return type

[**CapabilityRepresentation**](CapabilityRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

