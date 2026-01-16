# \ConversationsAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteConversationV1ConversationsConversationIdDelete**](ConversationsAPI.md#DeleteConversationV1ConversationsConversationIdDelete) | **Delete** /v1/conversations/{conversation_id} | Delete Conversation
[**GetConversationsV1ConversationsGet**](ConversationsAPI.md#GetConversationsV1ConversationsGet) | **Get** /v1/conversations | Get Conversations
[**ReadConversationV1ConversationsConversationIdGet**](ConversationsAPI.md#ReadConversationV1ConversationsConversationIdGet) | **Get** /v1/conversations/{conversation_id} | Read Conversation



## DeleteConversationV1ConversationsConversationIdDelete

> DeleteConversationV1ConversationsConversationIdDelete(ctx, conversationId).Execute()

Delete Conversation



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
	conversationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ConversationsAPI.DeleteConversationV1ConversationsConversationIdDelete(context.Background(), conversationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ConversationsAPI.DeleteConversationV1ConversationsConversationIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**conversationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteConversationV1ConversationsConversationIdDeleteRequest struct via the builder pattern


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


## GetConversationsV1ConversationsGet

> PagedResponseModelChatConversation GetConversationsV1ConversationsGet(ctx).Page(page).Size(size).Sort(sort).Filter(filter).Execute()

Get Conversations



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
	resp, r, err := apiClient.ConversationsAPI.GetConversationsV1ConversationsGet(context.Background()).Page(page).Size(size).Sort(sort).Filter(filter).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ConversationsAPI.GetConversationsV1ConversationsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetConversationsV1ConversationsGet`: PagedResponseModelChatConversation
	fmt.Fprintf(os.Stdout, "Response from `ConversationsAPI.GetConversationsV1ConversationsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetConversationsV1ConversationsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** |  | [default to 1]
 **size** | **int32** |  | [default to 10]
 **sort** | **string** |  | [default to &quot;id&quot;]
 **filter** | **string** |  | 

### Return type

[**PagedResponseModelChatConversation**](PagedResponseModelChatConversation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReadConversationV1ConversationsConversationIdGet

> ChatConversation ReadConversationV1ConversationsConversationIdGet(ctx, conversationId).Execute()

Read Conversation



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
	conversationId := "conversationId_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ConversationsAPI.ReadConversationV1ConversationsConversationIdGet(context.Background(), conversationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ConversationsAPI.ReadConversationV1ConversationsConversationIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReadConversationV1ConversationsConversationIdGet`: ChatConversation
	fmt.Fprintf(os.Stdout, "Response from `ConversationsAPI.ReadConversationV1ConversationsConversationIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**conversationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReadConversationV1ConversationsConversationIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ChatConversation**](ChatConversation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

