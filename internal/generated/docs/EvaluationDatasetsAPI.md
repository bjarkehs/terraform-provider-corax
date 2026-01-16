# \EvaluationDatasetsAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddItemsToEvaluationDatasetV1EvaluationDatasetsDatasetIdItemsPost**](EvaluationDatasetsAPI.md#AddItemsToEvaluationDatasetV1EvaluationDatasetsDatasetIdItemsPost) | **Post** /v1/evaluation-datasets/{dataset_id}/items | Add Items To Evaluation Dataset
[**CreateEvaluationDatasetV1EvaluationDatasetsPost**](EvaluationDatasetsAPI.md#CreateEvaluationDatasetV1EvaluationDatasetsPost) | **Post** /v1/evaluation-datasets | Create Evaluation Dataset
[**DeleteEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdDelete**](EvaluationDatasetsAPI.md#DeleteEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdDelete) | **Delete** /v1/evaluation-datasets/{dataset_id}/items/{item_id} | Delete Evaluation Dataset Item
[**DeleteEvaluationDatasetV1EvaluationDatasetsDatasetIdDelete**](EvaluationDatasetsAPI.md#DeleteEvaluationDatasetV1EvaluationDatasetsDatasetIdDelete) | **Delete** /v1/evaluation-datasets/{dataset_id} | Delete Evaluation Dataset
[**GetEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdGet**](EvaluationDatasetsAPI.md#GetEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdGet) | **Get** /v1/evaluation-datasets/{dataset_id}/items/{item_id} | Get Evaluation Dataset Item
[**GetEvaluationDatasetV1EvaluationDatasetsDatasetIdGet**](EvaluationDatasetsAPI.md#GetEvaluationDatasetV1EvaluationDatasetsDatasetIdGet) | **Get** /v1/evaluation-datasets/{dataset_id} | Get Evaluation Dataset
[**ListEvaluationDatasetItemsV1EvaluationDatasetsDatasetIdItemsGet**](EvaluationDatasetsAPI.md#ListEvaluationDatasetItemsV1EvaluationDatasetsDatasetIdItemsGet) | **Get** /v1/evaluation-datasets/{dataset_id}/items | List Evaluation Dataset Items
[**ListEvaluationDatasetsV1EvaluationDatasetsGet**](EvaluationDatasetsAPI.md#ListEvaluationDatasetsV1EvaluationDatasetsGet) | **Get** /v1/evaluation-datasets | List Evaluation Datasets
[**UpdateEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdPut**](EvaluationDatasetsAPI.md#UpdateEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdPut) | **Put** /v1/evaluation-datasets/{dataset_id}/items/{item_id} | Update Evaluation Dataset Item
[**UpdateEvaluationDatasetV1EvaluationDatasetsDatasetIdPut**](EvaluationDatasetsAPI.md#UpdateEvaluationDatasetV1EvaluationDatasetsDatasetIdPut) | **Put** /v1/evaluation-datasets/{dataset_id} | Update Evaluation Dataset



## AddItemsToEvaluationDatasetV1EvaluationDatasetsDatasetIdItemsPost

> EvaluationDatasetItemRepresentation AddItemsToEvaluationDatasetV1EvaluationDatasetsDatasetIdItemsPost(ctx, datasetId).EvaluationDatasetItemCreate(evaluationDatasetItemCreate).Execute()

Add Items To Evaluation Dataset



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
	datasetId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	evaluationDatasetItemCreate := *openapiclient.NewEvaluationDatasetItemCreate(map[string]string{"key": "Inner_example"}, map[string]interface{}{"key": interface{}(123)}) // EvaluationDatasetItemCreate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.EvaluationDatasetsAPI.AddItemsToEvaluationDatasetV1EvaluationDatasetsDatasetIdItemsPost(context.Background(), datasetId).EvaluationDatasetItemCreate(evaluationDatasetItemCreate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EvaluationDatasetsAPI.AddItemsToEvaluationDatasetV1EvaluationDatasetsDatasetIdItemsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddItemsToEvaluationDatasetV1EvaluationDatasetsDatasetIdItemsPost`: EvaluationDatasetItemRepresentation
	fmt.Fprintf(os.Stdout, "Response from `EvaluationDatasetsAPI.AddItemsToEvaluationDatasetV1EvaluationDatasetsDatasetIdItemsPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**datasetId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddItemsToEvaluationDatasetV1EvaluationDatasetsDatasetIdItemsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **evaluationDatasetItemCreate** | [**EvaluationDatasetItemCreate**](EvaluationDatasetItemCreate.md) |  | 

### Return type

[**EvaluationDatasetItemRepresentation**](EvaluationDatasetItemRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateEvaluationDatasetV1EvaluationDatasetsPost

> EvaluationDatasetRepresentation CreateEvaluationDatasetV1EvaluationDatasetsPost(ctx).EvaluationDatasetCreate(evaluationDatasetCreate).Execute()

Create Evaluation Dataset



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
	evaluationDatasetCreate := *openapiclient.NewEvaluationDatasetCreate("Name_example", *openapiclient.NewEvaluationDatasetConfigurationBase(openapiclient.EvaluationDatasetOutputType("text"))) // EvaluationDatasetCreate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.EvaluationDatasetsAPI.CreateEvaluationDatasetV1EvaluationDatasetsPost(context.Background()).EvaluationDatasetCreate(evaluationDatasetCreate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EvaluationDatasetsAPI.CreateEvaluationDatasetV1EvaluationDatasetsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateEvaluationDatasetV1EvaluationDatasetsPost`: EvaluationDatasetRepresentation
	fmt.Fprintf(os.Stdout, "Response from `EvaluationDatasetsAPI.CreateEvaluationDatasetV1EvaluationDatasetsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateEvaluationDatasetV1EvaluationDatasetsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **evaluationDatasetCreate** | [**EvaluationDatasetCreate**](EvaluationDatasetCreate.md) |  | 

### Return type

[**EvaluationDatasetRepresentation**](EvaluationDatasetRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdDelete

> DeleteEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdDelete(ctx, datasetId, itemId).Execute()

Delete Evaluation Dataset Item



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
	datasetId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	itemId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.EvaluationDatasetsAPI.DeleteEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdDelete(context.Background(), datasetId, itemId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EvaluationDatasetsAPI.DeleteEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**datasetId** | **string** |  | 
**itemId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdDeleteRequest struct via the builder pattern


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


## DeleteEvaluationDatasetV1EvaluationDatasetsDatasetIdDelete

> DeleteEvaluationDatasetV1EvaluationDatasetsDatasetIdDelete(ctx, datasetId).Execute()

Delete Evaluation Dataset



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
	datasetId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.EvaluationDatasetsAPI.DeleteEvaluationDatasetV1EvaluationDatasetsDatasetIdDelete(context.Background(), datasetId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EvaluationDatasetsAPI.DeleteEvaluationDatasetV1EvaluationDatasetsDatasetIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**datasetId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteEvaluationDatasetV1EvaluationDatasetsDatasetIdDeleteRequest struct via the builder pattern


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


## GetEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdGet

> EvaluationDatasetItemRepresentation GetEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdGet(ctx, datasetId, itemId).Execute()

Get Evaluation Dataset Item



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
	datasetId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	itemId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.EvaluationDatasetsAPI.GetEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdGet(context.Background(), datasetId, itemId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EvaluationDatasetsAPI.GetEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdGet`: EvaluationDatasetItemRepresentation
	fmt.Fprintf(os.Stdout, "Response from `EvaluationDatasetsAPI.GetEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**datasetId** | **string** |  | 
**itemId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**EvaluationDatasetItemRepresentation**](EvaluationDatasetItemRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetEvaluationDatasetV1EvaluationDatasetsDatasetIdGet

> EvaluationDatasetRepresentation GetEvaluationDatasetV1EvaluationDatasetsDatasetIdGet(ctx, datasetId).Execute()

Get Evaluation Dataset



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
	datasetId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.EvaluationDatasetsAPI.GetEvaluationDatasetV1EvaluationDatasetsDatasetIdGet(context.Background(), datasetId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EvaluationDatasetsAPI.GetEvaluationDatasetV1EvaluationDatasetsDatasetIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetEvaluationDatasetV1EvaluationDatasetsDatasetIdGet`: EvaluationDatasetRepresentation
	fmt.Fprintf(os.Stdout, "Response from `EvaluationDatasetsAPI.GetEvaluationDatasetV1EvaluationDatasetsDatasetIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**datasetId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetEvaluationDatasetV1EvaluationDatasetsDatasetIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**EvaluationDatasetRepresentation**](EvaluationDatasetRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListEvaluationDatasetItemsV1EvaluationDatasetsDatasetIdItemsGet

> PagedResponseModelEvaluationDatasetItemRepresentation ListEvaluationDatasetItemsV1EvaluationDatasetsDatasetIdItemsGet(ctx, datasetId).Page(page).Size(size).Sort(sort).Filter(filter).Execute()

List Evaluation Dataset Items



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
	datasetId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	page := int32(56) // int32 |  (optional) (default to 1)
	size := int32(56) // int32 |  (optional) (default to 10)
	sort := "sort_example" // string |  (optional) (default to "id")
	filter := "filter_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.EvaluationDatasetsAPI.ListEvaluationDatasetItemsV1EvaluationDatasetsDatasetIdItemsGet(context.Background(), datasetId).Page(page).Size(size).Sort(sort).Filter(filter).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EvaluationDatasetsAPI.ListEvaluationDatasetItemsV1EvaluationDatasetsDatasetIdItemsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListEvaluationDatasetItemsV1EvaluationDatasetsDatasetIdItemsGet`: PagedResponseModelEvaluationDatasetItemRepresentation
	fmt.Fprintf(os.Stdout, "Response from `EvaluationDatasetsAPI.ListEvaluationDatasetItemsV1EvaluationDatasetsDatasetIdItemsGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**datasetId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiListEvaluationDatasetItemsV1EvaluationDatasetsDatasetIdItemsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **page** | **int32** |  | [default to 1]
 **size** | **int32** |  | [default to 10]
 **sort** | **string** |  | [default to &quot;id&quot;]
 **filter** | **string** |  | 

### Return type

[**PagedResponseModelEvaluationDatasetItemRepresentation**](PagedResponseModelEvaluationDatasetItemRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListEvaluationDatasetsV1EvaluationDatasetsGet

> PagedResponseModelEvaluationDatasetRepresentation ListEvaluationDatasetsV1EvaluationDatasetsGet(ctx).Page(page).Size(size).Sort(sort).Filter(filter).Execute()

List Evaluation Datasets



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
	resp, r, err := apiClient.EvaluationDatasetsAPI.ListEvaluationDatasetsV1EvaluationDatasetsGet(context.Background()).Page(page).Size(size).Sort(sort).Filter(filter).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EvaluationDatasetsAPI.ListEvaluationDatasetsV1EvaluationDatasetsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListEvaluationDatasetsV1EvaluationDatasetsGet`: PagedResponseModelEvaluationDatasetRepresentation
	fmt.Fprintf(os.Stdout, "Response from `EvaluationDatasetsAPI.ListEvaluationDatasetsV1EvaluationDatasetsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListEvaluationDatasetsV1EvaluationDatasetsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** |  | [default to 1]
 **size** | **int32** |  | [default to 10]
 **sort** | **string** |  | [default to &quot;id&quot;]
 **filter** | **string** |  | 

### Return type

[**PagedResponseModelEvaluationDatasetRepresentation**](PagedResponseModelEvaluationDatasetRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdPut

> EvaluationDatasetItemRepresentation UpdateEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdPut(ctx, datasetId, itemId).EvaluationDatasetItemUpdate(evaluationDatasetItemUpdate).Execute()

Update Evaluation Dataset Item



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
	datasetId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	itemId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	evaluationDatasetItemUpdate := *openapiclient.NewEvaluationDatasetItemUpdate(map[string]string{"key": "Inner_example"}, map[string]interface{}{"key": interface{}(123)}, "Id_example") // EvaluationDatasetItemUpdate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.EvaluationDatasetsAPI.UpdateEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdPut(context.Background(), datasetId, itemId).EvaluationDatasetItemUpdate(evaluationDatasetItemUpdate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EvaluationDatasetsAPI.UpdateEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdPut`: EvaluationDatasetItemRepresentation
	fmt.Fprintf(os.Stdout, "Response from `EvaluationDatasetsAPI.UpdateEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**datasetId** | **string** |  | 
**itemId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateEvaluationDatasetItemV1EvaluationDatasetsDatasetIdItemsItemIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **evaluationDatasetItemUpdate** | [**EvaluationDatasetItemUpdate**](EvaluationDatasetItemUpdate.md) |  | 

### Return type

[**EvaluationDatasetItemRepresentation**](EvaluationDatasetItemRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateEvaluationDatasetV1EvaluationDatasetsDatasetIdPut

> EvaluationDatasetRepresentation UpdateEvaluationDatasetV1EvaluationDatasetsDatasetIdPut(ctx, datasetId).EvaluationDatasetUpdate(evaluationDatasetUpdate).Execute()

Update Evaluation Dataset



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
	datasetId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	evaluationDatasetUpdate := *openapiclient.NewEvaluationDatasetUpdate("Name_example", *openapiclient.NewEvaluationDatasetConfigurationBase(openapiclient.EvaluationDatasetOutputType("text")), "Id_example") // EvaluationDatasetUpdate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.EvaluationDatasetsAPI.UpdateEvaluationDatasetV1EvaluationDatasetsDatasetIdPut(context.Background(), datasetId).EvaluationDatasetUpdate(evaluationDatasetUpdate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EvaluationDatasetsAPI.UpdateEvaluationDatasetV1EvaluationDatasetsDatasetIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateEvaluationDatasetV1EvaluationDatasetsDatasetIdPut`: EvaluationDatasetRepresentation
	fmt.Fprintf(os.Stdout, "Response from `EvaluationDatasetsAPI.UpdateEvaluationDatasetV1EvaluationDatasetsDatasetIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**datasetId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateEvaluationDatasetV1EvaluationDatasetsDatasetIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **evaluationDatasetUpdate** | [**EvaluationDatasetUpdate**](EvaluationDatasetUpdate.md) |  | 

### Return type

[**EvaluationDatasetRepresentation**](EvaluationDatasetRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

