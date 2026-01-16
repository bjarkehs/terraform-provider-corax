# \MCPServersAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CheckMcpServerHealthV1McpServersServerIdHealthGet**](MCPServersAPI.md#CheckMcpServerHealthV1McpServersServerIdHealthGet) | **Get** /v1/mcp-servers/{server_id}/health | Check MCP server health status
[**CreateMcpServerV1McpServersPost**](MCPServersAPI.md#CreateMcpServerV1McpServersPost) | **Post** /v1/mcp-servers | Create a new MCP server
[**DeleteMcpServerV1McpServersServerIdDelete**](MCPServersAPI.md#DeleteMcpServerV1McpServersServerIdDelete) | **Delete** /v1/mcp-servers/{server_id} | Delete an MCP server
[**GetMcpServerEntitiesV1McpServersServerIdEntitiesGet**](MCPServersAPI.md#GetMcpServerEntitiesV1McpServersServerIdEntitiesGet) | **Get** /v1/mcp-servers/{server_id}/entities | Get all entities from an MCP server
[**GetMcpServerPromptsV1McpServersServerIdPromptsGet**](MCPServersAPI.md#GetMcpServerPromptsV1McpServersServerIdPromptsGet) | **Get** /v1/mcp-servers/{server_id}/prompts | Get all prompts from an MCP server
[**GetMcpServerResourcesV1McpServersServerIdResourcesGet**](MCPServersAPI.md#GetMcpServerResourcesV1McpServersServerIdResourcesGet) | **Get** /v1/mcp-servers/{server_id}/resources | Get all resources from an MCP server
[**GetMcpServerToolsV1McpServersServerIdToolsGet**](MCPServersAPI.md#GetMcpServerToolsV1McpServersServerIdToolsGet) | **Get** /v1/mcp-servers/{server_id}/tools | Get all tools from an MCP server
[**GetMcpServerV1McpServersServerIdGet**](MCPServersAPI.md#GetMcpServerV1McpServersServerIdGet) | **Get** /v1/mcp-servers/{server_id} | Get MCP server by ID
[**ListMcpServersV1McpServersGet**](MCPServersAPI.md#ListMcpServersV1McpServersGet) | **Get** /v1/mcp-servers | List MCP servers
[**UpdateMcpServerV1McpServersServerIdPut**](MCPServersAPI.md#UpdateMcpServerV1McpServersServerIdPut) | **Put** /v1/mcp-servers/{server_id} | Update an existing MCP server



## CheckMcpServerHealthV1McpServersServerIdHealthGet

> map[string]interface{} CheckMcpServerHealthV1McpServersServerIdHealthGet(ctx, serverId).Execute()

Check MCP server health status



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
	serverId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MCPServersAPI.CheckMcpServerHealthV1McpServersServerIdHealthGet(context.Background(), serverId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MCPServersAPI.CheckMcpServerHealthV1McpServersServerIdHealthGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CheckMcpServerHealthV1McpServersServerIdHealthGet`: map[string]interface{}
	fmt.Fprintf(os.Stdout, "Response from `MCPServersAPI.CheckMcpServerHealthV1McpServersServerIdHealthGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**serverId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiCheckMcpServerHealthV1McpServersServerIdHealthGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

**map[string]interface{}**

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateMcpServerV1McpServersPost

> MCPServerResponse CreateMcpServerV1McpServersPost(ctx).MCPServerBase(mCPServerBase).Execute()

Create a new MCP server



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
	mCPServerBase := *openapiclient.NewMCPServerBase("Name_example", "Url_example") // MCPServerBase | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MCPServersAPI.CreateMcpServerV1McpServersPost(context.Background()).MCPServerBase(mCPServerBase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MCPServersAPI.CreateMcpServerV1McpServersPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateMcpServerV1McpServersPost`: MCPServerResponse
	fmt.Fprintf(os.Stdout, "Response from `MCPServersAPI.CreateMcpServerV1McpServersPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateMcpServerV1McpServersPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mCPServerBase** | [**MCPServerBase**](MCPServerBase.md) |  | 

### Return type

[**MCPServerResponse**](MCPServerResponse.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteMcpServerV1McpServersServerIdDelete

> DeleteMcpServerV1McpServersServerIdDelete(ctx, serverId).Execute()

Delete an MCP server



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
	serverId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.MCPServersAPI.DeleteMcpServerV1McpServersServerIdDelete(context.Background(), serverId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MCPServersAPI.DeleteMcpServerV1McpServersServerIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**serverId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteMcpServerV1McpServersServerIdDeleteRequest struct via the builder pattern


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


## GetMcpServerEntitiesV1McpServersServerIdEntitiesGet

> []map[string]interface{} GetMcpServerEntitiesV1McpServersServerIdEntitiesGet(ctx, serverId).Execute()

Get all entities from an MCP server



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
	serverId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MCPServersAPI.GetMcpServerEntitiesV1McpServersServerIdEntitiesGet(context.Background(), serverId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MCPServersAPI.GetMcpServerEntitiesV1McpServersServerIdEntitiesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetMcpServerEntitiesV1McpServersServerIdEntitiesGet`: []map[string]interface{}
	fmt.Fprintf(os.Stdout, "Response from `MCPServersAPI.GetMcpServerEntitiesV1McpServersServerIdEntitiesGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**serverId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetMcpServerEntitiesV1McpServersServerIdEntitiesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]map[string]interface{}**](map.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetMcpServerPromptsV1McpServersServerIdPromptsGet

> []*map[string]interface{} GetMcpServerPromptsV1McpServersServerIdPromptsGet(ctx, serverId).Execute()

Get all prompts from an MCP server



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
	serverId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MCPServersAPI.GetMcpServerPromptsV1McpServersServerIdPromptsGet(context.Background(), serverId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MCPServersAPI.GetMcpServerPromptsV1McpServersServerIdPromptsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetMcpServerPromptsV1McpServersServerIdPromptsGet`: []*map[string]interface{}
	fmt.Fprintf(os.Stdout, "Response from `MCPServersAPI.GetMcpServerPromptsV1McpServersServerIdPromptsGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**serverId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetMcpServerPromptsV1McpServersServerIdPromptsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]*map[string]interface{}**](map.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetMcpServerResourcesV1McpServersServerIdResourcesGet

> []*map[string]interface{} GetMcpServerResourcesV1McpServersServerIdResourcesGet(ctx, serverId).Execute()

Get all resources from an MCP server



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
	serverId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MCPServersAPI.GetMcpServerResourcesV1McpServersServerIdResourcesGet(context.Background(), serverId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MCPServersAPI.GetMcpServerResourcesV1McpServersServerIdResourcesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetMcpServerResourcesV1McpServersServerIdResourcesGet`: []*map[string]interface{}
	fmt.Fprintf(os.Stdout, "Response from `MCPServersAPI.GetMcpServerResourcesV1McpServersServerIdResourcesGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**serverId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetMcpServerResourcesV1McpServersServerIdResourcesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]*map[string]interface{}**](map.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetMcpServerToolsV1McpServersServerIdToolsGet

> []map[string]interface{} GetMcpServerToolsV1McpServersServerIdToolsGet(ctx, serverId).Execute()

Get all tools from an MCP server



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
	serverId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MCPServersAPI.GetMcpServerToolsV1McpServersServerIdToolsGet(context.Background(), serverId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MCPServersAPI.GetMcpServerToolsV1McpServersServerIdToolsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetMcpServerToolsV1McpServersServerIdToolsGet`: []map[string]interface{}
	fmt.Fprintf(os.Stdout, "Response from `MCPServersAPI.GetMcpServerToolsV1McpServersServerIdToolsGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**serverId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetMcpServerToolsV1McpServersServerIdToolsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]map[string]interface{}**](map.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetMcpServerV1McpServersServerIdGet

> MCPServerResponse GetMcpServerV1McpServersServerIdGet(ctx, serverId).Execute()

Get MCP server by ID



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
	serverId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MCPServersAPI.GetMcpServerV1McpServersServerIdGet(context.Background(), serverId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MCPServersAPI.GetMcpServerV1McpServersServerIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetMcpServerV1McpServersServerIdGet`: MCPServerResponse
	fmt.Fprintf(os.Stdout, "Response from `MCPServersAPI.GetMcpServerV1McpServersServerIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**serverId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetMcpServerV1McpServersServerIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**MCPServerResponse**](MCPServerResponse.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListMcpServersV1McpServersGet

> PagedResponseModelMCPServerResponse ListMcpServersV1McpServersGet(ctx).Page(page).Size(size).Sort(sort).Filter(filter).Execute()

List MCP servers



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
	resp, r, err := apiClient.MCPServersAPI.ListMcpServersV1McpServersGet(context.Background()).Page(page).Size(size).Sort(sort).Filter(filter).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MCPServersAPI.ListMcpServersV1McpServersGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListMcpServersV1McpServersGet`: PagedResponseModelMCPServerResponse
	fmt.Fprintf(os.Stdout, "Response from `MCPServersAPI.ListMcpServersV1McpServersGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListMcpServersV1McpServersGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** |  | [default to 1]
 **size** | **int32** |  | [default to 10]
 **sort** | **string** |  | [default to &quot;id&quot;]
 **filter** | **string** |  | 

### Return type

[**PagedResponseModelMCPServerResponse**](PagedResponseModelMCPServerResponse.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateMcpServerV1McpServersServerIdPut

> MCPServerResponse UpdateMcpServerV1McpServersServerIdPut(ctx, serverId).MCPServerBase(mCPServerBase).Execute()

Update an existing MCP server



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
	serverId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	mCPServerBase := *openapiclient.NewMCPServerBase("Name_example", "Url_example") // MCPServerBase | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MCPServersAPI.UpdateMcpServerV1McpServersServerIdPut(context.Background(), serverId).MCPServerBase(mCPServerBase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MCPServersAPI.UpdateMcpServerV1McpServersServerIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateMcpServerV1McpServersServerIdPut`: MCPServerResponse
	fmt.Fprintf(os.Stdout, "Response from `MCPServersAPI.UpdateMcpServerV1McpServersServerIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**serverId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateMcpServerV1McpServersServerIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **mCPServerBase** | [**MCPServerBase**](MCPServerBase.md) |  | 

### Return type

[**MCPServerResponse**](MCPServerResponse.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

