# Plan: Address Truncated API Key Issue in corax_model_provider

## 1. Objective

To prevent Terraform from reporting an "inconsistent values for sensitive attribute" error for the `configuration.api_key` in the `corax_model_provider` resource. This error occurs when the Corax API returns a truncated version of the `api_key` for security reasons, which then mismatches the full `api_key` provided in the Terraform configuration.

## 2. Background

The `corax_model_provider` resource manages configurations for LLM providers, including sensitive data like API keys within a `configuration` map.

- The user provides a full `api_key` in their Terraform configuration.
- Upon resource creation or read, the Corax API is called.
- The API, for security, returns the `api_key` in a truncated format (e.g., `sk-...xxxx`).
- The Terraform provider, by default, updates its state based on this API response, leading to the state containing the truncated key.
- Terraform then compares the user's configuration (with the full key) against the provider's state (with the truncated key) and detects a difference, resulting in the "inconsistent values" error.

## 3. Core Strategy

The primary strategy is to ensure that the full `api_key` value, as provided in the Terraform configuration (plan) or as it exists in the prior Terraform state, is always preserved in the Terraform state. This means the provider will intentionally override the truncated `api_key` received from the API with the full version it already knows from the user's input or previous state.

## 4. Detailed Changes to `internal/provider/resource_model_provider.go`

The following modifications will be made to the [`internal/provider/resource_model_provider.go`](internal/provider/resource_model_provider.go) file:

### a. `Create` Function

1.  **Cache Planned Configuration:** Before making the API call to create the resource, store a copy of the `Configuration` map from the incoming plan (`req.Plan.Get(ctx, &plan)` then `plan.Configuration`).
2.  **API Call & Initial Mapping:** Proceed with the API call to create the model provider. The `mapAPIModelProviderToResourceModel` function will populate the `plan` variable, including `plan.Configuration` which will now likely hold the truncated `api_key` from the API response.
3.  **Restore Full API Key:** After the API response is mapped, check the cached planned configuration. If the `api_key` was present and non-empty in the cached planned configuration, explicitly set `plan.Configuration["api_key"]` to this full, cached value.
4.  **Set State:** Set the Terraform state using this corrected `plan` variable, which now contains the full `api_key`.

### b. `Read` Function

1.  **Cache Prior State Configuration:** Before making the API call to read the resource, store a copy of the `Configuration` map from the current Terraform state (`req.State.Get(ctx, &state)` then `state.Configuration`).
2.  **API Call & Initial Mapping:** Proceed with the API call to get the current state of the model provider. The `mapAPIModelProviderToResourceModel` function will populate the `state` variable, including `state.Configuration` which will now likely hold the truncated `api_key` from the API response.
3.  **Restore Full API Key:** After the API response is mapped, check the cached prior state configuration. If the `api_key` was present and non-empty in the cached prior state configuration, explicitly set `state.Configuration["api_key"]` to this full, cached value.
4.  **Set State:** Set the Terraform state using this corrected `state` variable.

### c. `Update` Function

The existing logic within the `Update` function is already robust for this scenario. It preserves the entire `plan.Configuration` (which includes the user's intended full `api_key`) and applies it to the `finalState.Configuration` after the API update call. This effectively overwrites any truncated key that might come from the API response during an update.
**No changes are anticipated for the `Update` function regarding this specific issue.**

## 5. Illustrative Diagram (Create Operation)

```mermaid
graph TD
    subgraph Revised Solution for Create (Truncated Key Scenario)
        direction LR
        P_Create[Initial Plan (req.Plan): config has FULL_api_key] --> Cache_PConfig_Create[Cache Planned Configuration (plannedConfig)];
        P_Create --> M_Create[Model for API Call: config has FULL_api_key];
        M_Create --> AC_Create[API Create Call];
        AC_Create --> AR_Create[API Response: config has TRUNCATED_api_key];
        AR_Create --> Map_Create[Provider maps API response to 'plan' variable];
        Map_Create --> S_Create['plan.Configuration' now has TRUNCATED_api_key, other fields (ID, CreatedAt) populated];
        Cache_PConfig_Create ----> Restore_FullKey_Create[Iterate plannedConfig: If 'api_key' exists and is non-empty in plannedConfig, overwrite 'plan.Configuration[\"api_key\"]' with FULL_api_key from plannedConfig];
        S_Create ----> Restore_FullKey_Create;
        Restore_FullKey_Create --> FS_Create[Final 'plan' var for State: 'plan.Configuration' has FULL_api_key, other fields from API response];
        FS_Create --> SetState_Create[Set Terraform State];
        P_Create ----> Compare_Create[Terraform Compares Plan & State: OK (both have FULL_api_key)];
        SetState_Create ----> Compare_Create;
    end
```

This plan ensures that the Terraform state accurately reflects the user's intended configuration for the `api_key`, preventing spurious diffs caused by the API's security measure of returning truncated keys.
