# \KnowledgeCollectionsAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateCollectionAsyncV1CollectionsPost**](KnowledgeCollectionsAPI.md#CreateCollectionAsyncV1CollectionsPost) | **Post** /v1/collections | Create Collection Async
[**DeleteCollectionV1CollectionsCollectionIdDelete**](KnowledgeCollectionsAPI.md#DeleteCollectionV1CollectionsCollectionIdDelete) | **Delete** /v1/collections/{collection_id} | Delete Collection
[**DeleteDocumentV1CollectionsCollectionIdDocumentsDocumentIdDelete**](KnowledgeCollectionsAPI.md#DeleteDocumentV1CollectionsCollectionIdDocumentsDocumentIdDelete) | **Delete** /v1/collections/{collection_id}/documents/{document_id} | Delete Document
[**EmbedDocumentsToCollectionAsyncV1CollectionsCollectionIdEmbedPost**](KnowledgeCollectionsAPI.md#EmbedDocumentsToCollectionAsyncV1CollectionsCollectionIdEmbedPost) | **Post** /v1/collections/{collection_id}/embed | Embed Documents To Collection Async
[**ReadCollectionAsyncV1CollectionsCollectionIdGet**](KnowledgeCollectionsAPI.md#ReadCollectionAsyncV1CollectionsCollectionIdGet) | **Get** /v1/collections/{collection_id} | Read Collection Async
[**ReadCollectionDocumentsAsyncV1CollectionsCollectionIdDocumentsGet**](KnowledgeCollectionsAPI.md#ReadCollectionDocumentsAsyncV1CollectionsCollectionIdDocumentsGet) | **Get** /v1/collections/{collection_id}/documents | Read Collection Documents Async
[**ReadCollectionsAsyncV1CollectionsGet**](KnowledgeCollectionsAPI.md#ReadCollectionsAsyncV1CollectionsGet) | **Get** /v1/collections | Read Collections Async
[**UpdateCollectionAsyncV1CollectionsCollectionIdPut**](KnowledgeCollectionsAPI.md#UpdateCollectionAsyncV1CollectionsCollectionIdPut) | **Put** /v1/collections/{collection_id} | Update Collection Async



## CreateCollectionAsyncV1CollectionsPost

> Collection CreateCollectionAsyncV1CollectionsPost(ctx).CollectionCreate(collectionCreate).Execute()

Create Collection Async

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
	collectionCreate := *openapiclient.NewCollectionCreate("Name_example") // CollectionCreate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.KnowledgeCollectionsAPI.CreateCollectionAsyncV1CollectionsPost(context.Background()).CollectionCreate(collectionCreate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `KnowledgeCollectionsAPI.CreateCollectionAsyncV1CollectionsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateCollectionAsyncV1CollectionsPost`: Collection
	fmt.Fprintf(os.Stdout, "Response from `KnowledgeCollectionsAPI.CreateCollectionAsyncV1CollectionsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateCollectionAsyncV1CollectionsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **collectionCreate** | [**CollectionCreate**](CollectionCreate.md) |  | 

### Return type

[**Collection**](Collection.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteCollectionV1CollectionsCollectionIdDelete

> DeleteCollectionV1CollectionsCollectionIdDelete(ctx, collectionId).Execute()

Delete Collection

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
	collectionId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.KnowledgeCollectionsAPI.DeleteCollectionV1CollectionsCollectionIdDelete(context.Background(), collectionId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `KnowledgeCollectionsAPI.DeleteCollectionV1CollectionsCollectionIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**collectionId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteCollectionV1CollectionsCollectionIdDeleteRequest struct via the builder pattern


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


## DeleteDocumentV1CollectionsCollectionIdDocumentsDocumentIdDelete

> DeleteDocumentV1CollectionsCollectionIdDocumentsDocumentIdDelete(ctx, collectionId, documentId).Execute()

Delete Document

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
	collectionId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	documentId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.KnowledgeCollectionsAPI.DeleteDocumentV1CollectionsCollectionIdDocumentsDocumentIdDelete(context.Background(), collectionId, documentId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `KnowledgeCollectionsAPI.DeleteDocumentV1CollectionsCollectionIdDocumentsDocumentIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**collectionId** | **string** |  | 
**documentId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteDocumentV1CollectionsCollectionIdDocumentsDocumentIdDeleteRequest struct via the builder pattern


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


## EmbedDocumentsToCollectionAsyncV1CollectionsCollectionIdEmbedPost

> DocumentEmbeddingResponse EmbedDocumentsToCollectionAsyncV1CollectionsCollectionIdEmbedPost(ctx, collectionId).Files(files).Execute()

Embed Documents To Collection Async



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
	collectionId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	files := []*os.File{"TODO"} // []*os.File | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.KnowledgeCollectionsAPI.EmbedDocumentsToCollectionAsyncV1CollectionsCollectionIdEmbedPost(context.Background(), collectionId).Files(files).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `KnowledgeCollectionsAPI.EmbedDocumentsToCollectionAsyncV1CollectionsCollectionIdEmbedPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `EmbedDocumentsToCollectionAsyncV1CollectionsCollectionIdEmbedPost`: DocumentEmbeddingResponse
	fmt.Fprintf(os.Stdout, "Response from `KnowledgeCollectionsAPI.EmbedDocumentsToCollectionAsyncV1CollectionsCollectionIdEmbedPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**collectionId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiEmbedDocumentsToCollectionAsyncV1CollectionsCollectionIdEmbedPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **files** | **[]*os.File** |  | 

### Return type

[**DocumentEmbeddingResponse**](DocumentEmbeddingResponse.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReadCollectionAsyncV1CollectionsCollectionIdGet

> Collection ReadCollectionAsyncV1CollectionsCollectionIdGet(ctx, collectionId).Execute()

Read Collection Async

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
	collectionId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.KnowledgeCollectionsAPI.ReadCollectionAsyncV1CollectionsCollectionIdGet(context.Background(), collectionId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `KnowledgeCollectionsAPI.ReadCollectionAsyncV1CollectionsCollectionIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReadCollectionAsyncV1CollectionsCollectionIdGet`: Collection
	fmt.Fprintf(os.Stdout, "Response from `KnowledgeCollectionsAPI.ReadCollectionAsyncV1CollectionsCollectionIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**collectionId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReadCollectionAsyncV1CollectionsCollectionIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Collection**](Collection.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReadCollectionDocumentsAsyncV1CollectionsCollectionIdDocumentsGet

> PagedResponseModelDocument ReadCollectionDocumentsAsyncV1CollectionsCollectionIdDocumentsGet(ctx, collectionId).Page(page).Size(size).Sort(sort).Filter(filter).Execute()

Read Collection Documents Async

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
	collectionId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	page := int32(56) // int32 |  (optional) (default to 1)
	size := int32(56) // int32 |  (optional) (default to 10)
	sort := "sort_example" // string |  (optional) (default to "id")
	filter := "filter_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.KnowledgeCollectionsAPI.ReadCollectionDocumentsAsyncV1CollectionsCollectionIdDocumentsGet(context.Background(), collectionId).Page(page).Size(size).Sort(sort).Filter(filter).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `KnowledgeCollectionsAPI.ReadCollectionDocumentsAsyncV1CollectionsCollectionIdDocumentsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReadCollectionDocumentsAsyncV1CollectionsCollectionIdDocumentsGet`: PagedResponseModelDocument
	fmt.Fprintf(os.Stdout, "Response from `KnowledgeCollectionsAPI.ReadCollectionDocumentsAsyncV1CollectionsCollectionIdDocumentsGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**collectionId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReadCollectionDocumentsAsyncV1CollectionsCollectionIdDocumentsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **page** | **int32** |  | [default to 1]
 **size** | **int32** |  | [default to 10]
 **sort** | **string** |  | [default to &quot;id&quot;]
 **filter** | **string** |  | 

### Return type

[**PagedResponseModelDocument**](PagedResponseModelDocument.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReadCollectionsAsyncV1CollectionsGet

> PagedResponseModelCollection ReadCollectionsAsyncV1CollectionsGet(ctx).Page(page).Size(size).Sort(sort).Filter(filter).Execute()

Read Collections Async

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
	resp, r, err := apiClient.KnowledgeCollectionsAPI.ReadCollectionsAsyncV1CollectionsGet(context.Background()).Page(page).Size(size).Sort(sort).Filter(filter).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `KnowledgeCollectionsAPI.ReadCollectionsAsyncV1CollectionsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReadCollectionsAsyncV1CollectionsGet`: PagedResponseModelCollection
	fmt.Fprintf(os.Stdout, "Response from `KnowledgeCollectionsAPI.ReadCollectionsAsyncV1CollectionsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiReadCollectionsAsyncV1CollectionsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** |  | [default to 1]
 **size** | **int32** |  | [default to 10]
 **sort** | **string** |  | [default to &quot;id&quot;]
 **filter** | **string** |  | 

### Return type

[**PagedResponseModelCollection**](PagedResponseModelCollection.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateCollectionAsyncV1CollectionsCollectionIdPut

> Collection UpdateCollectionAsyncV1CollectionsCollectionIdPut(ctx, collectionId).CollectionUpdate(collectionUpdate).Execute()

Update Collection Async

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
	collectionId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	collectionUpdate := *openapiclient.NewCollectionUpdate("Name_example") // CollectionUpdate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.KnowledgeCollectionsAPI.UpdateCollectionAsyncV1CollectionsCollectionIdPut(context.Background(), collectionId).CollectionUpdate(collectionUpdate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `KnowledgeCollectionsAPI.UpdateCollectionAsyncV1CollectionsCollectionIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateCollectionAsyncV1CollectionsCollectionIdPut`: Collection
	fmt.Fprintf(os.Stdout, "Response from `KnowledgeCollectionsAPI.UpdateCollectionAsyncV1CollectionsCollectionIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**collectionId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateCollectionAsyncV1CollectionsCollectionIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **collectionUpdate** | [**CollectionUpdate**](CollectionUpdate.md) |  | 

### Return type

[**Collection**](Collection.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

