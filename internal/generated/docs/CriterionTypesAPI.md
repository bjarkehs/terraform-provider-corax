# \CriterionTypesAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetEvaluationCriterionTypeV1CriterionTypesCriterionTypeGet**](CriterionTypesAPI.md#GetEvaluationCriterionTypeV1CriterionTypesCriterionTypeGet) | **Get** /v1/criterion-types/{criterion_type} | Get Evaluation Criterion Type
[**ListEvaluationCriteriaTypesV1CriterionTypesGet**](CriterionTypesAPI.md#ListEvaluationCriteriaTypesV1CriterionTypesGet) | **Get** /v1/criterion-types | List Evaluation Criteria Types



## GetEvaluationCriterionTypeV1CriterionTypesCriterionTypeGet

> CapabilityEvaluationCriterionTypeRepresentation GetEvaluationCriterionTypeV1CriterionTypesCriterionTypeGet(ctx, criterionType).Execute()

Get Evaluation Criterion Type



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
	criterionType := "criterionType_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CriterionTypesAPI.GetEvaluationCriterionTypeV1CriterionTypesCriterionTypeGet(context.Background(), criterionType).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CriterionTypesAPI.GetEvaluationCriterionTypeV1CriterionTypesCriterionTypeGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetEvaluationCriterionTypeV1CriterionTypesCriterionTypeGet`: CapabilityEvaluationCriterionTypeRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CriterionTypesAPI.GetEvaluationCriterionTypeV1CriterionTypesCriterionTypeGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**criterionType** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetEvaluationCriterionTypeV1CriterionTypesCriterionTypeGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CapabilityEvaluationCriterionTypeRepresentation**](CapabilityEvaluationCriterionTypeRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListEvaluationCriteriaTypesV1CriterionTypesGet

> PagedResponseModelCapabilityEvaluationCriterionTypeRepresentation ListEvaluationCriteriaTypesV1CriterionTypesGet(ctx).Page(page).Size(size).Sort(sort).Filter(filter).Execute()

List Evaluation Criteria Types



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
	resp, r, err := apiClient.CriterionTypesAPI.ListEvaluationCriteriaTypesV1CriterionTypesGet(context.Background()).Page(page).Size(size).Sort(sort).Filter(filter).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CriterionTypesAPI.ListEvaluationCriteriaTypesV1CriterionTypesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListEvaluationCriteriaTypesV1CriterionTypesGet`: PagedResponseModelCapabilityEvaluationCriterionTypeRepresentation
	fmt.Fprintf(os.Stdout, "Response from `CriterionTypesAPI.ListEvaluationCriteriaTypesV1CriterionTypesGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListEvaluationCriteriaTypesV1CriterionTypesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** |  | [default to 1]
 **size** | **int32** |  | [default to 10]
 **sort** | **string** |  | [default to &quot;id&quot;]
 **filter** | **string** |  | 

### Return type

[**PagedResponseModelCapabilityEvaluationCriterionTypeRepresentation**](PagedResponseModelCapabilityEvaluationCriterionTypeRepresentation.md)

### Authorization

[APIKeyHeader](../README.md#APIKeyHeader), [HTTPBearer](../README.md#HTTPBearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

