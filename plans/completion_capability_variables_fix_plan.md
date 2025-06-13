# Plan to Make `variables` Order-Insensitive in `corax_completion_capability`

The `variables` attribute in the `corax_completion_capability` resource is currently order-sensitive, leading to "Provider produced inconsistent result after apply" errors when the order of variables changes. This plan outlines the steps to make this attribute order-insensitive by changing its type from a List to a Set.

## Affected File:

- `internal/provider/resource_completion_capability.go`

## Steps:

1.  **Modify Data Model (`CompletionCapabilityResourceModel`)**:

    - In `internal/provider/resource_completion_capability.go`, change the type of the `Variables` field (currently at line 49) from `types.List` to `types.Set`.

    ```mermaid
    graph LR
        A[types.List] --> B(types.Set);
    ```

2.  **Update Resource Schema (`Schema()` method)**:

    - In the `Schema` method (starting at line 63), change the `variables` attribute from `schema.ListAttribute` (line 99) to `schema.SetAttribute`. The `ElementType` will remain `types.StringType`.

3.  **Adjust API-to-Model Mapping (`mapAPICompletionCapabilityToModel()` function)**:

    - In the `mapAPICompletionCapabilityToModel` function (starting at line 313), when populating `model.Variables` from the API response (lines 387-450):
      - Convert the list of strings (either directly from an API list or from sorted keys of an API map) into a `types.Set` using `types.SetValueFrom()`.
      - Ensure `model.Variables` is set to `types.SetNull(types.StringType)` if variables are not present or in case of conversion errors.

4.  **Adjust `Create()` Method**:

    - The line `resp.Diagnostics.Append(plan.Variables.ElementsAs(ctx, &apiPayload.Variables, false)...)` (around line 499) in the `Create` method will not require changes as `ElementsAs` supports `types.Set`. The `apiPayload.Variables` is assumed to be `[]string`.

5.  **Adjust `Update()` Method**:

    - Similarly, the line `resp.Diagnostics.Append(plan.Variables.ElementsAs(ctx, &vars, false)...)` (around line 628) in the `Update` method will not require changes.

6.  **Review and Update Tests**:
    - Associated test files (e.g., `internal/provider/resource_completion_capability_test.go`) will need updates to reflect that `variables` is now a set and assertions should not depend on order.
