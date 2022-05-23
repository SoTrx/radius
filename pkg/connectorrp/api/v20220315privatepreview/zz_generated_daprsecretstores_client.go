//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package v20220315privatepreview

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// DaprSecretStoresClient contains the methods for the DaprSecretStores group.
// Don't use this type directly, use NewDaprSecretStoresClient() instead.
type DaprSecretStoresClient struct {
	con *connection
	subscriptionID string
}

// NewDaprSecretStoresClient creates a new instance of DaprSecretStoresClient with the specified values.
func NewDaprSecretStoresClient(con *connection, subscriptionID string) *DaprSecretStoresClient {
	return &DaprSecretStoresClient{con: con, subscriptionID: subscriptionID}
}

// CreateOrUpdate - Creates or updates a DaprSecretStore resource
// If the operation fails it returns the *ErrorResponse error type.
func (client *DaprSecretStoresClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, daprSecretStoreName string, daprSecretStoreParameters DaprSecretStoreResource, options *DaprSecretStoresCreateOrUpdateOptions) (DaprSecretStoresCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, daprSecretStoreName, daprSecretStoreParameters, options)
	if err != nil {
		return DaprSecretStoresCreateOrUpdateResponse{}, err
	}
	resp, err := 	client.con.Pipeline().Do(req)
	if err != nil {
		return DaprSecretStoresCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return DaprSecretStoresCreateOrUpdateResponse{}, client.createOrUpdateHandleError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *DaprSecretStoresClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, daprSecretStoreName string, daprSecretStoreParameters DaprSecretStoreResource, options *DaprSecretStoresCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Applications.Connector/daprSecretStores/{daprSecretStoreName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if daprSecretStoreName == "" {
		return nil, errors.New("parameter daprSecretStoreName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{daprSecretStoreName}", url.PathEscape(daprSecretStoreName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, daprSecretStoreParameters)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *DaprSecretStoresClient) createOrUpdateHandleResponse(resp *http.Response) (DaprSecretStoresCreateOrUpdateResponse, error) {
	result := DaprSecretStoresCreateOrUpdateResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.DaprSecretStoreResource); err != nil {
		return DaprSecretStoresCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client *DaprSecretStoresClient) createOrUpdateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
		errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Delete - Deletes an existing daprSecretStore resource
// If the operation fails it returns the *ErrorResponse error type.
func (client *DaprSecretStoresClient) Delete(ctx context.Context, resourceGroupName string, daprSecretStoreName string, options *DaprSecretStoresDeleteOptions) (DaprSecretStoresDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, daprSecretStoreName, options)
	if err != nil {
		return DaprSecretStoresDeleteResponse{}, err
	}
	resp, err := 	client.con.Pipeline().Do(req)
	if err != nil {
		return DaprSecretStoresDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return DaprSecretStoresDeleteResponse{}, client.deleteHandleError(resp)
	}
	return DaprSecretStoresDeleteResponse{RawResponse: resp}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *DaprSecretStoresClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, daprSecretStoreName string, options *DaprSecretStoresDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Applications.Connector/daprSecretStores/{daprSecretStoreName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if daprSecretStoreName == "" {
		return nil, errors.New("parameter daprSecretStoreName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{daprSecretStoreName}", url.PathEscape(daprSecretStoreName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client *DaprSecretStoresClient) deleteHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
		errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Get - Retrieves information about a daprSecretStore resource
// If the operation fails it returns the *ErrorResponse error type.
func (client *DaprSecretStoresClient) Get(ctx context.Context, resourceGroupName string, daprSecretStoreName string, options *DaprSecretStoresGetOptions) (DaprSecretStoresGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, daprSecretStoreName, options)
	if err != nil {
		return DaprSecretStoresGetResponse{}, err
	}
	resp, err := 	client.con.Pipeline().Do(req)
	if err != nil {
		return DaprSecretStoresGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DaprSecretStoresGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *DaprSecretStoresClient) getCreateRequest(ctx context.Context, resourceGroupName string, daprSecretStoreName string, options *DaprSecretStoresGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Applications.Connector/daprSecretStores/{daprSecretStoreName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if daprSecretStoreName == "" {
		return nil, errors.New("parameter daprSecretStoreName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{daprSecretStoreName}", url.PathEscape(daprSecretStoreName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *DaprSecretStoresClient) getHandleResponse(resp *http.Response) (DaprSecretStoresGetResponse, error) {
	result := DaprSecretStoresGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.DaprSecretStoreResource); err != nil {
		return DaprSecretStoresGetResponse{}, err
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *DaprSecretStoresClient) getHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
		errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// List - Lists information about all daprSecretStore resources in the given subscription and resource group
// If the operation fails it returns the *ErrorResponse error type.
func (client *DaprSecretStoresClient) List(resourceGroupName string, options *DaprSecretStoresListOptions) (*DaprSecretStoresListPager) {
	return &DaprSecretStoresListPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listCreateRequest(ctx, resourceGroupName, options)
		},
		advancer: func(ctx context.Context, resp DaprSecretStoresListResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.DaprSecretStoreList.NextLink)
		},
	}
}

// listCreateRequest creates the List request.
func (client *DaprSecretStoresClient) listCreateRequest(ctx context.Context, resourceGroupName string, options *DaprSecretStoresListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Applications.Connector/daprSecretStores"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *DaprSecretStoresClient) listHandleResponse(resp *http.Response) (DaprSecretStoresListResponse, error) {
	result := DaprSecretStoresListResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.DaprSecretStoreList); err != nil {
		return DaprSecretStoresListResponse{}, err
	}
	return result, nil
}

// listHandleError handles the List error response.
func (client *DaprSecretStoresClient) listHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
		errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// ListBySubscription - Lists information about all daprSecretStore resources in the given subscription
// If the operation fails it returns the *ErrorResponse error type.
func (client *DaprSecretStoresClient) ListBySubscription(options *DaprSecretStoresListBySubscriptionOptions) (*DaprSecretStoresListBySubscriptionPager) {
	return &DaprSecretStoresListBySubscriptionPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listBySubscriptionCreateRequest(ctx, options)
		},
		advancer: func(ctx context.Context, resp DaprSecretStoresListBySubscriptionResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.DaprSecretStoreList.NextLink)
		},
	}
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *DaprSecretStoresClient) listBySubscriptionCreateRequest(ctx context.Context, options *DaprSecretStoresListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Applications.Connector/daprSecretStores"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *DaprSecretStoresClient) listBySubscriptionHandleResponse(resp *http.Response) (DaprSecretStoresListBySubscriptionResponse, error) {
	result := DaprSecretStoresListBySubscriptionResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.DaprSecretStoreList); err != nil {
		return DaprSecretStoresListBySubscriptionResponse{}, err
	}
	return result, nil
}

// listBySubscriptionHandleError handles the ListBySubscription error response.
func (client *DaprSecretStoresClient) listBySubscriptionHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
		errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

