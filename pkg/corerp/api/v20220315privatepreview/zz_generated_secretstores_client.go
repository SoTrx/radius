//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package v20220315privatepreview

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// SecretStoresClient contains the methods for the SecretStores group.
// Don't use this type directly, use NewSecretStoresClient() instead.
type SecretStoresClient struct {
	host string
	rootScope string
	pl runtime.Pipeline
}

// NewSecretStoresClient creates a new instance of SecretStoresClient with the specified values.
// rootScope - The scope in which the resource is present. For Azure resource this would be /subscriptions/{subscriptionID}/resourceGroups/{resourcegroupID}
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewSecretStoresClient(rootScope string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SecretStoresClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &SecretStoresClient{
		rootScope: rootScope,
		host: ep,
pl: pl,
	}
	return client, nil
}

// CreateOrUpdate - Create or update a secret store.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-15-privatepreview
// secretStoreName - The name of the secret store.
// secretStoreResource - SecretStore details
// options - SecretStoresClientCreateOrUpdateOptions contains the optional parameters for the SecretStoresClient.CreateOrUpdate
// method.
func (client *SecretStoresClient) CreateOrUpdate(ctx context.Context, secretStoreName string, secretStoreResource SecretStoreResource, options *SecretStoresClientCreateOrUpdateOptions) (SecretStoresClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, secretStoreName, secretStoreResource, options)
	if err != nil {
		return SecretStoresClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return SecretStoresClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return SecretStoresClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *SecretStoresClient) createOrUpdateCreateRequest(ctx context.Context, secretStoreName string, secretStoreResource SecretStoreResource, options *SecretStoresClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/secretStores/{secretStoreName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if secretStoreName == "" {
		return nil, errors.New("parameter secretStoreName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secretStoreName}", url.PathEscape(secretStoreName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, secretStoreResource)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *SecretStoresClient) createOrUpdateHandleResponse(resp *http.Response) (SecretStoresClientCreateOrUpdateResponse, error) {
	result := SecretStoresClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SecretStoreResource); err != nil {
		return SecretStoresClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Delete a secret store.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-15-privatepreview
// secretStoreName - The name of the secret store.
// options - SecretStoresClientDeleteOptions contains the optional parameters for the SecretStoresClient.Delete method.
func (client *SecretStoresClient) Delete(ctx context.Context, secretStoreName string, options *SecretStoresClientDeleteOptions) (SecretStoresClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, secretStoreName, options)
	if err != nil {
		return SecretStoresClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return SecretStoresClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return SecretStoresClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return SecretStoresClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *SecretStoresClient) deleteCreateRequest(ctx context.Context, secretStoreName string, options *SecretStoresClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/secretStores/{secretStoreName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if secretStoreName == "" {
		return nil, errors.New("parameter secretStoreName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secretStoreName}", url.PathEscape(secretStoreName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the properties of a secret store.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-15-privatepreview
// secretStoreName - The name of the secret store.
// options - SecretStoresClientGetOptions contains the optional parameters for the SecretStoresClient.Get method.
func (client *SecretStoresClient) Get(ctx context.Context, secretStoreName string, options *SecretStoresClientGetOptions) (SecretStoresClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, secretStoreName, options)
	if err != nil {
		return SecretStoresClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return SecretStoresClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return SecretStoresClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *SecretStoresClient) getCreateRequest(ctx context.Context, secretStoreName string, options *SecretStoresClientGetOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/secretStores/{secretStoreName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if secretStoreName == "" {
		return nil, errors.New("parameter secretStoreName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secretStoreName}", url.PathEscape(secretStoreName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *SecretStoresClient) getHandleResponse(resp *http.Response) (SecretStoresClientGetResponse, error) {
	result := SecretStoresClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SecretStoreResource); err != nil {
		return SecretStoresClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List all secret stores in the given scope.
// Generated from API version 2022-03-15-privatepreview
// options - SecretStoresClientListOptions contains the optional parameters for the SecretStoresClient.List method.
func (client *SecretStoresClient) NewListPager(options *SecretStoresClientListOptions) (*runtime.Pager[SecretStoresClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[SecretStoresClientListResponse]{
		More: func(page SecretStoresClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *SecretStoresClientListResponse) (SecretStoresClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return SecretStoresClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return SecretStoresClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return SecretStoresClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *SecretStoresClient) listCreateRequest(ctx context.Context, options *SecretStoresClientListOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/secretStores"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *SecretStoresClient) listHandleResponse(resp *http.Response) (SecretStoresClientListResponse, error) {
	result := SecretStoresClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SecretStoreResourceList); err != nil {
		return SecretStoresClientListResponse{}, err
	}
	return result, nil
}

// ListSecrets - List the secrets of a secret stores.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-15-privatepreview
// secretStoreName - The name of the secret store.
// options - SecretStoresClientListSecretsOptions contains the optional parameters for the SecretStoresClient.ListSecrets
// method.
func (client *SecretStoresClient) ListSecrets(ctx context.Context, secretStoreName string, options *SecretStoresClientListSecretsOptions) (SecretStoresClientListSecretsResponse, error) {
	req, err := client.listSecretsCreateRequest(ctx, secretStoreName, options)
	if err != nil {
		return SecretStoresClientListSecretsResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return SecretStoresClientListSecretsResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return SecretStoresClientListSecretsResponse{}, runtime.NewResponseError(resp)
	}
	return client.listSecretsHandleResponse(resp)
}

// listSecretsCreateRequest creates the ListSecrets request.
func (client *SecretStoresClient) listSecretsCreateRequest(ctx context.Context, secretStoreName string, options *SecretStoresClientListSecretsOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/secretStores/{secretStoreName}/listSecrets"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if secretStoreName == "" {
		return nil, errors.New("parameter secretStoreName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secretStoreName}", url.PathEscape(secretStoreName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listSecretsHandleResponse handles the ListSecrets response.
func (client *SecretStoresClient) listSecretsHandleResponse(resp *http.Response) (SecretStoresClientListSecretsResponse, error) {
	result := SecretStoresClientListSecretsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SecretStoreListSecretsResult); err != nil {
		return SecretStoresClientListSecretsResponse{}, err
	}
	return result, nil
}

// Update - Update the properties of an existing secret store.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-15-privatepreview
// secretStoreName - The name of the secret store.
// secretStoreResource - SecretStore details
// options - SecretStoresClientUpdateOptions contains the optional parameters for the SecretStoresClient.Update method.
func (client *SecretStoresClient) Update(ctx context.Context, secretStoreName string, secretStoreResource SecretStoreResource, options *SecretStoresClientUpdateOptions) (SecretStoresClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, secretStoreName, secretStoreResource, options)
	if err != nil {
		return SecretStoresClientUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return SecretStoresClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return SecretStoresClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *SecretStoresClient) updateCreateRequest(ctx context.Context, secretStoreName string, secretStoreResource SecretStoreResource, options *SecretStoresClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/secretStores/{secretStoreName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if secretStoreName == "" {
		return nil, errors.New("parameter secretStoreName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secretStoreName}", url.PathEscape(secretStoreName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, secretStoreResource)
}

// updateHandleResponse handles the Update response.
func (client *SecretStoresClient) updateHandleResponse(resp *http.Response) (SecretStoresClientUpdateResponse, error) {
	result := SecretStoresClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SecretStoreResource); err != nil {
		return SecretStoresClientUpdateResponse{}, err
	}
	return result, nil
}

