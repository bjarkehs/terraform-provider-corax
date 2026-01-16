# \CapabilityEvaluationsAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaPost**](CapabilityEvaluationsAPI.md#CreateEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaPost) | **Post** /v1/capability-evaluations/{evaluation_id}/criteria | Create Evaluation Criterion
[**CreateEvaluationV1CapabilityEvaluationsPost**](CapabilityEvaluationsAPI.md#CreateEvaluationV1CapabilityEvaluationsPost) | **Post** /v1/capability-evaluations | Create Evaluation
[**DeleteEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdDelete**](CapabilityEvaluationsAPI.md#DeleteEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdDelete) | **Delete** /v1/capability-evaluations/{evaluation_id}/criteria/{criterion_id} | Delete Evaluation Criterion
[**DeleteEvaluationV1CapabilityEvaluationsEvaluationIdDelete**](CapabilityEvaluationsAPI.md#DeleteEvaluationV1CapabilityEvaluationsEvaluationIdDelete) | **Delete** /v1/capability-evaluations/{evaluation_id} | Delete Evaluation
[**ExecuteEvaluationV1CapabilityEvaluationsEvaluationIdExecutionsPost**](CapabilityEvaluationsAPI.md#ExecuteEvaluationV1CapabilityEvaluationsEvaluationIdExecutionsPost) | **Post** /v1/capability-evaluations/{evaluation_id}/executions | Execute Evaluation
[**GetEvaluationCriteriaV1CapabilityEvaluationsEvaluationIdCriteriaGet**](CapabilityEvaluationsAPI.md#GetEvaluationCriteriaV1CapabilityEvaluationsEvaluationIdCriteriaGet) | **Get** /v1/capability-evaluations/{evaluation_id}/criteria | Get Evaluation Criteria
[**GetEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdGet**](CapabilityEvaluationsAPI.md#GetEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdGet) | **Get** /v1/capability-evaluations/{evaluation_id}/criteria/{criterion_id} | Get Evaluation Criterion
[**GetEvaluationExecutionCriteriaExecutionsV1CapabilityEvaluationsEvaluationIdExecutionsEvaluationExecutionIdCriteriaExecutionsGet**](CapabilityEvaluationsAPI.md#GetEvaluationExecutionCriteriaExecutionsV1CapabilityEvaluationsEvaluationIdExecutionsEvaluationExecutionIdCriteriaExecutionsGet) | **Get** /v1/capability-evaluations/{evaluation_id}/executions/{evaluation_execution_id}/criteria-executions | Get Evaluation Execution Criteria Executions
[**GetEvaluationExecutionV1CapabilityEvaluationsEvaluationIdExecutionsEvaluationExecutionIdGet**](CapabilityEvaluationsAPI.md#GetEvaluationExecutionV1CapabilityEvaluationsEvaluationIdExecutionsEvaluationExecutionIdGet) | **Get** /v1/capability-evaluations/{evaluation_id}/executions/{evaluation_execution_id} | Get Evaluation Execution
[**GetEvaluationExecutionsV1CapabilityEvaluationsEvaluationIdExecutionsGet**](CapabilityEvaluationsAPI.md#GetEvaluationExecutionsV1CapabilityEvaluationsEvaluationIdExecutionsGet) | **Get** /v1/capability-evaluations/{evaluation_id}/executions | Get Evaluation Executions
[**GetEvaluationV1CapabilityEvaluationsEvaluationIdGet**](CapabilityEvaluationsAPI.md#GetEvaluationV1CapabilityEvaluationsEvaluationIdGet) | **Get** /v1/capability-evaluations/{evaluation_id} | Get Evaluation
[**ListEvaluationsV1CapabilityEvaluationsGet**](CapabilityEvaluationsAPI.md#ListEvaluationsV1CapabilityEvaluationsGet) | **Get** /v1/capability-evaluations | List Evaluations
[**UpdateEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdPut**](CapabilityEvaluationsAPI.md#UpdateEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdPut) | **Put** /v1/capability-evaluations/{evaluation_id}/criteria/{criterion_id} | Update Evaluation Criterion
[**UpdateEvaluationV1CapabilityEvaluationsEvaluationIdPut**](CapabilityEvaluationsAPI.md#UpdateEvaluationV1CapabilityEvaluationsEvaluationIdPut) | **Put** /v1/capability-evaluations/{evaluation_id} | Update Evaluation



## CreateEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaPost

> CapabilityEvaluationCriterionRepresentation CreateEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaPost(ctx, evaluationId).CapabilityEvaluationCriterionCreate(capabilityEvaluationCriterionCreate).Execute()

Create Evaluation Criterion



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
	evaluationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	capabilityEvaluationCriterionCreate := *openapiclient.NewCapabilityEvaluationCriterionCreate(openapiclient.CapabilityCriterionType("correctness"), map[string]interface{}{"key": interface{}(123)}) // CapabilityEvaluationCriterionCreate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilityEvaluationsAPI.CreateEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaPost(context.Background(), evaluationId).CapabilityEvaluationCriterionCreate(capabilityEvaluationCriterionCreate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilityEvaluationsAPI.CreateEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaPost`: CapabilityEvaluationCriterionRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CapabilityEvaluationsAPI.CreateEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**evaluationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **capabilityEvaluationCriterionCreate** | [**CapabilityEvaluationCriterionCreate**](CapabilityEvaluationCriterionCreate.md) |  | 

### Return type

[**CapabilityEvaluationCriterionRepresentation**](CapabilityEvaluationCriterionRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateEvaluationV1CapabilityEvaluationsPost

> CapabilityEvaluationRepresentation CreateEvaluationV1CapabilityEvaluationsPost(ctx).CapabilityEvaluationCreate(capabilityEvaluationCreate).Execute()

Create Evaluation



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
	capabilityEvaluationCreate := *openapiclient.NewCapabilityEvaluationCreate("Name_example", "CapabilityId_example", "EvaluationDatasetId_example", *openapiclient.NewCapabilityEvaluationConfiguration()) // CapabilityEvaluationCreate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilityEvaluationsAPI.CreateEvaluationV1CapabilityEvaluationsPost(context.Background()).CapabilityEvaluationCreate(capabilityEvaluationCreate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilityEvaluationsAPI.CreateEvaluationV1CapabilityEvaluationsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateEvaluationV1CapabilityEvaluationsPost`: CapabilityEvaluationRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CapabilityEvaluationsAPI.CreateEvaluationV1CapabilityEvaluationsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateEvaluationV1CapabilityEvaluationsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **capabilityEvaluationCreate** | [**CapabilityEvaluationCreate**](CapabilityEvaluationCreate.md) |  | 

### Return type

[**CapabilityEvaluationRepresentation**](CapabilityEvaluationRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdDelete

> DeleteEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdDelete(ctx, evaluationId, criterionId).Execute()

Delete Evaluation Criterion



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
	evaluationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	criterionId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.CapabilityEvaluationsAPI.DeleteEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdDelete(context.Background(), evaluationId, criterionId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilityEvaluationsAPI.DeleteEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**evaluationId** | **string** |  | 
**criterionId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdDeleteRequest struct via the builder pattern


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


## DeleteEvaluationV1CapabilityEvaluationsEvaluationIdDelete

> DeleteEvaluationV1CapabilityEvaluationsEvaluationIdDelete(ctx, evaluationId).Execute()

Delete Evaluation



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
	evaluationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.CapabilityEvaluationsAPI.DeleteEvaluationV1CapabilityEvaluationsEvaluationIdDelete(context.Background(), evaluationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilityEvaluationsAPI.DeleteEvaluationV1CapabilityEvaluationsEvaluationIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**evaluationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteEvaluationV1CapabilityEvaluationsEvaluationIdDeleteRequest struct via the builder pattern


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


## ExecuteEvaluationV1CapabilityEvaluationsEvaluationIdExecutionsPost

> EvaluationExecutionRepresentation ExecuteEvaluationV1CapabilityEvaluationsEvaluationIdExecutionsPost(ctx, evaluationId).EvaluationExecutionCreate(evaluationExecutionCreate).Execute()

Execute Evaluation



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
	evaluationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	evaluationExecutionCreate := *openapiclient.NewEvaluationExecutionCreate("EvaluationId_example") // EvaluationExecutionCreate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilityEvaluationsAPI.ExecuteEvaluationV1CapabilityEvaluationsEvaluationIdExecutionsPost(context.Background(), evaluationId).EvaluationExecutionCreate(evaluationExecutionCreate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilityEvaluationsAPI.ExecuteEvaluationV1CapabilityEvaluationsEvaluationIdExecutionsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExecuteEvaluationV1CapabilityEvaluationsEvaluationIdExecutionsPost`: EvaluationExecutionRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CapabilityEvaluationsAPI.ExecuteEvaluationV1CapabilityEvaluationsEvaluationIdExecutionsPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**evaluationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiExecuteEvaluationV1CapabilityEvaluationsEvaluationIdExecutionsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **evaluationExecutionCreate** | [**EvaluationExecutionCreate**](EvaluationExecutionCreate.md) |  | 

### Return type

[**EvaluationExecutionRepresentation**](EvaluationExecutionRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetEvaluationCriteriaV1CapabilityEvaluationsEvaluationIdCriteriaGet

> PagedResponseModelCapabilityEvaluationCriterionRepresentation GetEvaluationCriteriaV1CapabilityEvaluationsEvaluationIdCriteriaGet(ctx, evaluationId).Page(page).Size(size).Sort(sort).Filter(filter).Execute()

Get Evaluation Criteria



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
	evaluationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	page := int32(56) // int32 |  (optional) (default to 1)
	size := int32(56) // int32 |  (optional) (default to 10)
	sort := "sort_example" // string |  (optional) (default to "id")
	filter := "filter_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilityEvaluationsAPI.GetEvaluationCriteriaV1CapabilityEvaluationsEvaluationIdCriteriaGet(context.Background(), evaluationId).Page(page).Size(size).Sort(sort).Filter(filter).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilityEvaluationsAPI.GetEvaluationCriteriaV1CapabilityEvaluationsEvaluationIdCriteriaGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetEvaluationCriteriaV1CapabilityEvaluationsEvaluationIdCriteriaGet`: PagedResponseModelCapabilityEvaluationCriterionRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CapabilityEvaluationsAPI.GetEvaluationCriteriaV1CapabilityEvaluationsEvaluationIdCriteriaGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**evaluationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetEvaluationCriteriaV1CapabilityEvaluationsEvaluationIdCriteriaGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **page** | **int32** |  | [default to 1]
 **size** | **int32** |  | [default to 10]
 **sort** | **string** |  | [default to &quot;id&quot;]
 **filter** | **string** |  | 

### Return type

[**PagedResponseModelCapabilityEvaluationCriterionRepresentation**](PagedResponseModelCapabilityEvaluationCriterionRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdGet

> CapabilityEvaluationCriterionRepresentation GetEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdGet(ctx, evaluationId, criterionId).Execute()

Get Evaluation Criterion



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
	evaluationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	criterionId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilityEvaluationsAPI.GetEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdGet(context.Background(), evaluationId, criterionId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilityEvaluationsAPI.GetEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdGet`: CapabilityEvaluationCriterionRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CapabilityEvaluationsAPI.GetEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**evaluationId** | **string** |  | 
**criterionId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**CapabilityEvaluationCriterionRepresentation**](CapabilityEvaluationCriterionRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetEvaluationExecutionCriteriaExecutionsV1CapabilityEvaluationsEvaluationIdExecutionsEvaluationExecutionIdCriteriaExecutionsGet

> PagedResponseModelEvaluationCriterionExecutionRepresentation GetEvaluationExecutionCriteriaExecutionsV1CapabilityEvaluationsEvaluationIdExecutionsEvaluationExecutionIdCriteriaExecutionsGet(ctx, evaluationId, evaluationExecutionId).Page(page).Size(size).Sort(sort).Filter(filter).Execute()

Get Evaluation Execution Criteria Executions



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
	evaluationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	evaluationExecutionId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	page := int32(56) // int32 |  (optional) (default to 1)
	size := int32(56) // int32 |  (optional) (default to 10)
	sort := "sort_example" // string |  (optional) (default to "id")
	filter := "filter_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilityEvaluationsAPI.GetEvaluationExecutionCriteriaExecutionsV1CapabilityEvaluationsEvaluationIdExecutionsEvaluationExecutionIdCriteriaExecutionsGet(context.Background(), evaluationId, evaluationExecutionId).Page(page).Size(size).Sort(sort).Filter(filter).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilityEvaluationsAPI.GetEvaluationExecutionCriteriaExecutionsV1CapabilityEvaluationsEvaluationIdExecutionsEvaluationExecutionIdCriteriaExecutionsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetEvaluationExecutionCriteriaExecutionsV1CapabilityEvaluationsEvaluationIdExecutionsEvaluationExecutionIdCriteriaExecutionsGet`: PagedResponseModelEvaluationCriterionExecutionRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CapabilityEvaluationsAPI.GetEvaluationExecutionCriteriaExecutionsV1CapabilityEvaluationsEvaluationIdExecutionsEvaluationExecutionIdCriteriaExecutionsGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**evaluationId** | **string** |  | 
**evaluationExecutionId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetEvaluationExecutionCriteriaExecutionsV1CapabilityEvaluationsEvaluationIdExecutionsEvaluationExecutionIdCriteriaExecutionsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **page** | **int32** |  | [default to 1]
 **size** | **int32** |  | [default to 10]
 **sort** | **string** |  | [default to &quot;id&quot;]
 **filter** | **string** |  | 

### Return type

[**PagedResponseModelEvaluationCriterionExecutionRepresentation**](PagedResponseModelEvaluationCriterionExecutionRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetEvaluationExecutionV1CapabilityEvaluationsEvaluationIdExecutionsEvaluationExecutionIdGet

> EvaluationExecutionRepresentation GetEvaluationExecutionV1CapabilityEvaluationsEvaluationIdExecutionsEvaluationExecutionIdGet(ctx, evaluationId, evaluationExecutionId).Execute()

Get Evaluation Execution



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
	evaluationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	evaluationExecutionId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilityEvaluationsAPI.GetEvaluationExecutionV1CapabilityEvaluationsEvaluationIdExecutionsEvaluationExecutionIdGet(context.Background(), evaluationId, evaluationExecutionId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilityEvaluationsAPI.GetEvaluationExecutionV1CapabilityEvaluationsEvaluationIdExecutionsEvaluationExecutionIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetEvaluationExecutionV1CapabilityEvaluationsEvaluationIdExecutionsEvaluationExecutionIdGet`: EvaluationExecutionRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CapabilityEvaluationsAPI.GetEvaluationExecutionV1CapabilityEvaluationsEvaluationIdExecutionsEvaluationExecutionIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**evaluationId** | **string** |  | 
**evaluationExecutionId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetEvaluationExecutionV1CapabilityEvaluationsEvaluationIdExecutionsEvaluationExecutionIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**EvaluationExecutionRepresentation**](EvaluationExecutionRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetEvaluationExecutionsV1CapabilityEvaluationsEvaluationIdExecutionsGet

> PagedResponseModelEvaluationExecutionRepresentation GetEvaluationExecutionsV1CapabilityEvaluationsEvaluationIdExecutionsGet(ctx, evaluationId).CapabilityId(capabilityId).Page(page).Size(size).Sort(sort).Filter(filter).Execute()

Get Evaluation Executions



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
	evaluationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	capabilityId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string |  (optional)
	page := int32(56) // int32 |  (optional) (default to 1)
	size := int32(56) // int32 |  (optional) (default to 10)
	sort := "sort_example" // string |  (optional) (default to "id")
	filter := "filter_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilityEvaluationsAPI.GetEvaluationExecutionsV1CapabilityEvaluationsEvaluationIdExecutionsGet(context.Background(), evaluationId).CapabilityId(capabilityId).Page(page).Size(size).Sort(sort).Filter(filter).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilityEvaluationsAPI.GetEvaluationExecutionsV1CapabilityEvaluationsEvaluationIdExecutionsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetEvaluationExecutionsV1CapabilityEvaluationsEvaluationIdExecutionsGet`: PagedResponseModelEvaluationExecutionRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CapabilityEvaluationsAPI.GetEvaluationExecutionsV1CapabilityEvaluationsEvaluationIdExecutionsGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**evaluationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetEvaluationExecutionsV1CapabilityEvaluationsEvaluationIdExecutionsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **capabilityId** | **string** |  | 
 **page** | **int32** |  | [default to 1]
 **size** | **int32** |  | [default to 10]
 **sort** | **string** |  | [default to &quot;id&quot;]
 **filter** | **string** |  | 

### Return type

[**PagedResponseModelEvaluationExecutionRepresentation**](PagedResponseModelEvaluationExecutionRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetEvaluationV1CapabilityEvaluationsEvaluationIdGet

> CapabilityEvaluationRepresentation GetEvaluationV1CapabilityEvaluationsEvaluationIdGet(ctx, evaluationId).Execute()

Get Evaluation



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
	evaluationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilityEvaluationsAPI.GetEvaluationV1CapabilityEvaluationsEvaluationIdGet(context.Background(), evaluationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilityEvaluationsAPI.GetEvaluationV1CapabilityEvaluationsEvaluationIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetEvaluationV1CapabilityEvaluationsEvaluationIdGet`: CapabilityEvaluationRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CapabilityEvaluationsAPI.GetEvaluationV1CapabilityEvaluationsEvaluationIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**evaluationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetEvaluationV1CapabilityEvaluationsEvaluationIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CapabilityEvaluationRepresentation**](CapabilityEvaluationRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListEvaluationsV1CapabilityEvaluationsGet

> PagedResponseModelCapabilityEvaluationRepresentation ListEvaluationsV1CapabilityEvaluationsGet(ctx).Page(page).Size(size).Sort(sort).Filter(filter).Execute()

List Evaluations



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
	resp, r, err := apiClient.CapabilityEvaluationsAPI.ListEvaluationsV1CapabilityEvaluationsGet(context.Background()).Page(page).Size(size).Sort(sort).Filter(filter).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilityEvaluationsAPI.ListEvaluationsV1CapabilityEvaluationsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListEvaluationsV1CapabilityEvaluationsGet`: PagedResponseModelCapabilityEvaluationRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CapabilityEvaluationsAPI.ListEvaluationsV1CapabilityEvaluationsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListEvaluationsV1CapabilityEvaluationsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** |  | [default to 1]
 **size** | **int32** |  | [default to 10]
 **sort** | **string** |  | [default to &quot;id&quot;]
 **filter** | **string** |  | 

### Return type

[**PagedResponseModelCapabilityEvaluationRepresentation**](PagedResponseModelCapabilityEvaluationRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdPut

> CapabilityEvaluationCriterionRepresentation UpdateEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdPut(ctx, evaluationId, criterionId).CapabilityEvaluationCriterionCreate(capabilityEvaluationCriterionCreate).Execute()

Update Evaluation Criterion



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
	evaluationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	criterionId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	capabilityEvaluationCriterionCreate := *openapiclient.NewCapabilityEvaluationCriterionCreate(openapiclient.CapabilityCriterionType("correctness"), map[string]interface{}{"key": interface{}(123)}) // CapabilityEvaluationCriterionCreate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilityEvaluationsAPI.UpdateEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdPut(context.Background(), evaluationId, criterionId).CapabilityEvaluationCriterionCreate(capabilityEvaluationCriterionCreate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilityEvaluationsAPI.UpdateEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdPut`: CapabilityEvaluationCriterionRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CapabilityEvaluationsAPI.UpdateEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**evaluationId** | **string** |  | 
**criterionId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateEvaluationCriterionV1CapabilityEvaluationsEvaluationIdCriteriaCriterionIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **capabilityEvaluationCriterionCreate** | [**CapabilityEvaluationCriterionCreate**](CapabilityEvaluationCriterionCreate.md) |  | 

### Return type

[**CapabilityEvaluationCriterionRepresentation**](CapabilityEvaluationCriterionRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateEvaluationV1CapabilityEvaluationsEvaluationIdPut

> CapabilityEvaluationRepresentation UpdateEvaluationV1CapabilityEvaluationsEvaluationIdPut(ctx, evaluationId).CapabilityEvaluationUpdate(capabilityEvaluationUpdate).Execute()

Update Evaluation



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
	evaluationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	capabilityEvaluationUpdate := *openapiclient.NewCapabilityEvaluationUpdate("Name_example", "CapabilityId_example", "EvaluationDatasetId_example", *openapiclient.NewCapabilityEvaluationConfiguration()) // CapabilityEvaluationUpdate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CapabilityEvaluationsAPI.UpdateEvaluationV1CapabilityEvaluationsEvaluationIdPut(context.Background(), evaluationId).CapabilityEvaluationUpdate(capabilityEvaluationUpdate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CapabilityEvaluationsAPI.UpdateEvaluationV1CapabilityEvaluationsEvaluationIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateEvaluationV1CapabilityEvaluationsEvaluationIdPut`: CapabilityEvaluationRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CapabilityEvaluationsAPI.UpdateEvaluationV1CapabilityEvaluationsEvaluationIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**evaluationId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateEvaluationV1CapabilityEvaluationsEvaluationIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **capabilityEvaluationUpdate** | [**CapabilityEvaluationUpdate**](CapabilityEvaluationUpdate.md) |  | 

### Return type

[**CapabilityEvaluationRepresentation**](CapabilityEvaluationRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

