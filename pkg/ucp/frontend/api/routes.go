/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"context"
	"fmt"

	"github.com/gorilla/mux"
	v1 "github.com/project-radius/radius/pkg/armrpc/api/v1"
	frontend_ctrl "github.com/project-radius/radius/pkg/armrpc/frontend/controller"
	"github.com/project-radius/radius/pkg/armrpc/frontend/defaultoperation"
	"github.com/project-radius/radius/pkg/armrpc/frontend/server"
	"github.com/project-radius/radius/pkg/ucp/api/v20220901privatepreview"
	ucp_aws "github.com/project-radius/radius/pkg/ucp/aws"
	"github.com/project-radius/radius/pkg/ucp/datamodel"
	"github.com/project-radius/radius/pkg/ucp/datamodel/converter"
	awsproxy_ctrl "github.com/project-radius/radius/pkg/ucp/frontend/controller/awsproxy"
	aws_credential_ctrl "github.com/project-radius/radius/pkg/ucp/frontend/controller/credentials/aws"
	azure_credential_ctrl "github.com/project-radius/radius/pkg/ucp/frontend/controller/credentials/azure"
	kubernetes_ctrl "github.com/project-radius/radius/pkg/ucp/frontend/controller/kubernetes"
	planes_ctrl "github.com/project-radius/radius/pkg/ucp/frontend/controller/planes"
	resourcegroups_ctrl "github.com/project-radius/radius/pkg/ucp/frontend/controller/resourcegroups"
	"github.com/project-radius/radius/pkg/ucp/secret"
	"github.com/project-radius/radius/pkg/ucp/ucplog"
	"github.com/project-radius/radius/pkg/validator"
	"github.com/project-radius/radius/swagger"
)

// TODO: Use variables and construct the path as we add more APIs.
const (
	planeCollectionPath           = "/planes"
	awsPlaneType                  = "/planes/aws"
	planeItemPath                 = "/planes/{planeType}/{planeName}"
	planeCollectionByType         = "/planes/{planeType}"
	awsOperationResultsPath       = "/{AWSPlaneName}/accounts/{AccountID}/regions/{Region}/providers/{Provider}/locations/{Location}/operationResults/{operationID}"
	awsOperationStatusesPath      = "/{AWSPlaneName}/accounts/{AccountID}/regions/{Region}/providers/{Provider}/locations/{Location}/operationStatuses/{operationID}"
	awsResourceCollectionPath     = "/{AWSPlaneName}/accounts/{AccountID}/regions/{Region}/providers/{Provider}/{ResourceType}"
	awsResourcePath               = "/{AWSPlaneName}/accounts/{AccountID}/regions/{Region}/providers/{Provider}/{ResourceType}/{ResourceName}"
	putPath                       = "put"
	getPath                       = "get"
	deletePath                    = "delete"
	azureCredentialCollectionPath = "/planes/azure/{planeName}/providers/{Provider}/{ResourceType}"
	azureCredentialResourcePath   = "/planes/azure/{planeName}/providers/{Provider}/{ResourceType}/{ResourceName}"
	awsCredentialCollectionPath   = "/planes/aws/{planeName}/providers/{Provider}/{ResourceType}"
	awsCredentialResourcePath     = "/planes/aws/{planeName}/providers/{Provider}/{ResourceType}/{ResourceName}"
)

// Register registers the routes for UCP
func Register(ctx context.Context, router *mux.Router, ctrlOpts frontend_ctrl.Options, secretClient secret.Client, awsClients ucp_aws.Clients) error {
	logger := ucplog.FromContextOrDiscard(ctx)
	logger.Info(fmt.Sprintf("Registering routes with path base: %s", ctrlOpts.PathBase))

	router.NotFoundHandler = validator.APINotFoundHandler()
	router.MethodNotAllowedHandler = validator.APIMethodNotAllowedHandler()

	handlerOptions := []server.HandlerOptions{}

	// If we're in Kubernetes we have some required routes to implement.
	if ctrlOpts.PathBase != "" {
		// NOTE: the Kubernetes API Server does not include the gvr (base path) in
		// the URL for swagger routes.
		handlerOptions = append(handlerOptions, []server.HandlerOptions{
			{
				ParentRouter:      router.Path("/openapi/v2").Subrouter(),
				OperationType:     &v1.OperationType{Type: "KUBERNETESOPENAPIV2DOC", Method: v1.OperationGet},
				Method:            v1.OperationGet,
				ControllerFactory: kubernetes_ctrl.NewOpenAPIv2Doc,
			},
			{
				ParentRouter:      router.Path(ctrlOpts.PathBase).Subrouter(),
				OperationType:     &v1.OperationType{Type: "KUBERNETESDISCOVERYDOC", Method: v1.OperationGet},
				Method:            v1.OperationGet,
				ControllerFactory: kubernetes_ctrl.NewDiscoveryDoc,
			},
		}...)
	}

	specLoader, err := validator.LoadSpec(ctx, "ucp", swagger.SpecFilesUCP, []string{ctrlOpts.PathBase}, "")
	if err != nil {
		return err
	}

	rootScopeRouter := router.PathPrefix(ctrlOpts.PathBase).Subrouter()
	rootScopeRouter.Use(validator.APIValidatorUCP(specLoader))

	planeCollectionSubRouter := rootScopeRouter.Path(planeCollectionPath).Subrouter()
	planeCollectionByTypeSubRouter := rootScopeRouter.Path(planeCollectionByType).Subrouter()
	planeSubRouter := rootScopeRouter.Path(planeItemPath).Subrouter()

	var resourceGroupCollectionPath = fmt.Sprintf("%s/%s", planeItemPath, "resourcegroups")
	var resourceGroupItemPath = fmt.Sprintf("%s/%s", resourceGroupCollectionPath, "{resourceGroupName}")
	resourceGroupCollectionSubRouter := rootScopeRouter.Path(resourceGroupCollectionPath).Subrouter()
	resourceGroupSubRouter := rootScopeRouter.Path(resourceGroupItemPath).Subrouter()

	awsResourcesSubRouter := router.PathPrefix(fmt.Sprintf("%s%s", ctrlOpts.PathBase, awsPlaneType)).Subrouter()
	awsResourceCollectionSubRouter := awsResourcesSubRouter.Path(awsResourceCollectionPath).Subrouter()
	awsSingleResourceSubRouter := awsResourcesSubRouter.Path(awsResourcePath).Subrouter()
	awsOperationStatusesSubRouter := awsResourcesSubRouter.PathPrefix(awsOperationStatusesPath).Subrouter()
	awsOperationResultsSubRouter := awsResourcesSubRouter.PathPrefix(awsOperationResultsPath).Subrouter()
	awsPutResourceSubRouter := awsResourcesSubRouter.Path(fmt.Sprintf("%s/:%s", awsResourceCollectionPath, putPath)).Subrouter()
	awsGetResourceSubRouter := awsResourcesSubRouter.Path(fmt.Sprintf("%s/:%s", awsResourceCollectionPath, getPath)).Subrouter()
	awsDeleteResourceSubRouter := awsResourcesSubRouter.Path(fmt.Sprintf("%s/:%s", awsResourceCollectionPath, deletePath)).Subrouter()

	azureCredentialCollectionSubRouter := router.Path(fmt.Sprintf("%s%s", ctrlOpts.PathBase, azureCredentialCollectionPath)).Subrouter()
	azureCredentialResourceSubRouter := router.Path(fmt.Sprintf("%s%s", ctrlOpts.PathBase, azureCredentialResourcePath)).Subrouter()
	awsCredentialCollectionSubRouter := router.Path(fmt.Sprintf("%s%s", ctrlOpts.PathBase, awsCredentialCollectionPath)).Subrouter()
	awsCredentialResourceSubRouter := router.Path(fmt.Sprintf("%s%s", ctrlOpts.PathBase, awsCredentialResourcePath)).Subrouter()

	handlerOptions = append(handlerOptions, []server.HandlerOptions{
		// Planes resource handler registration.
		{
			// This is scope query unlike the default list handler.
			ParentRouter:      planeCollectionSubRouter,
			Method:            v1.OperationList,
			OperationType:     &v1.OperationType{Type: "PLANES", Method: v1.OperationList},
			ControllerFactory: planes_ctrl.NewListPlanes,
		},
		{
			// This is scope query unlike the default list handler.
			ParentRouter:      planeCollectionByTypeSubRouter,
			Method:            v1.OperationList,
			OperationType:     &v1.OperationType{Type: "PLANESBYTYPE", Method: v1.OperationList},
			ControllerFactory: planes_ctrl.NewListPlanesByType,
		},
		{
			ParentRouter:  planeSubRouter,
			Method:        v1.OperationGet,
			OperationType: &v1.OperationType{Type: "PLANESBYTYPE", Method: v1.OperationGet},
			ControllerFactory: func(opt frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return defaultoperation.NewGetResource(opt,
					frontend_ctrl.ResourceOptions[datamodel.Plane]{
						RequestConverter:  converter.PlaneDataModelFromVersioned,
						ResponseConverter: converter.PlaneDataModelToVersioned,
					},
				)
			},
		},
		{
			ParentRouter:  planeSubRouter,
			Method:        v1.OperationPut,
			OperationType: &v1.OperationType{Type: "PLANESBYTYPE", Method: v1.OperationPut},
			ControllerFactory: func(opt frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return defaultoperation.NewDefaultSyncPut(opt,
					frontend_ctrl.ResourceOptions[datamodel.Plane]{
						RequestConverter:  converter.PlaneDataModelFromVersioned,
						ResponseConverter: converter.PlaneDataModelToVersioned,
					},
				)
			},
		},
		{
			ParentRouter:  planeSubRouter,
			Method:        v1.OperationDelete,
			OperationType: &v1.OperationType{Type: "PLANESBYTYPE", Method: v1.OperationDelete},
			ControllerFactory: func(opt frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return defaultoperation.NewDefaultSyncDelete(opt,
					frontend_ctrl.ResourceOptions[datamodel.Plane]{
						RequestConverter:  converter.PlaneDataModelFromVersioned,
						ResponseConverter: converter.PlaneDataModelToVersioned,
					},
				)
			},
		},

		// Resource group handler registration
		{
			// This is scope query unlike the default list handler.
			ParentRouter:      resourceGroupCollectionSubRouter,
			ResourceType:      v20220901privatepreview.ResourceGroupType,
			Method:            v1.OperationList,
			ControllerFactory: resourcegroups_ctrl.NewListResourceGroups,
		},
		{
			ParentRouter: resourceGroupSubRouter,
			ResourceType: v20220901privatepreview.ResourceGroupType,
			Method:       v1.OperationGet,
			ControllerFactory: func(opt frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return defaultoperation.NewGetResource(opt,
					frontend_ctrl.ResourceOptions[datamodel.ResourceGroup]{
						RequestConverter:  converter.ResourceGroupDataModelFromVersioned,
						ResponseConverter: converter.ResourceGroupDataModelToVersioned,
					},
				)
			},
		},
		{
			ParentRouter: resourceGroupSubRouter,
			ResourceType: v20220901privatepreview.ResourceGroupType,
			Method:       v1.OperationPut,
			ControllerFactory: func(opt frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return defaultoperation.NewDefaultSyncPut(opt,
					frontend_ctrl.ResourceOptions[datamodel.ResourceGroup]{
						RequestConverter:  converter.ResourceGroupDataModelFromVersioned,
						ResponseConverter: converter.ResourceGroupDataModelToVersioned,
					},
				)
			},
		},
		{
			ParentRouter: resourceGroupSubRouter,
			ResourceType: v20220901privatepreview.ResourceGroupType,
			Method:       v1.OperationDelete,
			ControllerFactory: func(opt frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return defaultoperation.NewDefaultSyncDelete(opt,
					frontend_ctrl.ResourceOptions[datamodel.ResourceGroup]{
						RequestConverter:  converter.ResourceGroupDataModelFromVersioned,
						ResponseConverter: converter.ResourceGroupDataModelToVersioned,
					},
				)
			},
		},

		// AWS Plane handlers
		{
			ParentRouter:  awsOperationResultsSubRouter,
			Method:        v1.OperationGetOperationResult,
			OperationType: &v1.OperationType{Type: "AWSRESOURCE", Method: v1.OperationGetOperationResult},
			ControllerFactory: func(opt frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return awsproxy_ctrl.NewGetAWSOperationResults(opt, awsClients)
			},
		},
		{
			ParentRouter:  awsOperationStatusesSubRouter,
			Method:        v1.OperationGetOperationStatuses,
			OperationType: &v1.OperationType{Type: "AWSRESOURCE", Method: v1.OperationGetOperationStatuses},
			ControllerFactory: func(opts frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return awsproxy_ctrl.NewGetAWSOperationStatuses(opts, awsClients)
			},
		},
		{
			ParentRouter:  awsResourceCollectionSubRouter,
			Method:        v1.OperationList,
			OperationType: &v1.OperationType{Type: "AWSRESOURCE", Method: v1.OperationList},
			ControllerFactory: func(opts frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return awsproxy_ctrl.NewListAWSResources(opts, awsClients)
			},
		},
		{
			ParentRouter:  awsSingleResourceSubRouter,
			Method:        v1.OperationPut,
			OperationType: &v1.OperationType{Type: "AWSRESOURCE", Method: v1.OperationPut},
			ControllerFactory: func(opts frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return awsproxy_ctrl.NewCreateOrUpdateAWSResource(opts, awsClients)
			},
		},
		{
			ParentRouter:  awsSingleResourceSubRouter,
			Method:        v1.OperationDelete,
			OperationType: &v1.OperationType{Type: "AWSRESOURCE", Method: v1.OperationDelete},
			ControllerFactory: func(opts frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return awsproxy_ctrl.NewDeleteAWSResource(opts, awsClients)
			},
		},
		{
			ParentRouter:  awsSingleResourceSubRouter,
			Method:        v1.OperationGet,
			OperationType: &v1.OperationType{Type: "AWSRESOURCE", Method: v1.OperationGet},
			ControllerFactory: func(opt frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return awsproxy_ctrl.NewGetAWSResource(opt, awsClients)
			},
		},
		{
			ParentRouter:  awsPutResourceSubRouter,
			Method:        v1.OperationPutImperative,
			OperationType: &v1.OperationType{Type: "AWSRESOURCE", Method: v1.OperationPutImperative},
			ControllerFactory: func(opt frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return awsproxy_ctrl.NewCreateOrUpdateAWSResourceWithPost(opt, awsClients)
			},
		},
		{
			ParentRouter:  awsGetResourceSubRouter,
			Method:        v1.OperationGetImperative,
			OperationType: &v1.OperationType{Type: "AWSRESOURCE", Method: v1.OperationGetImperative},
			ControllerFactory: func(opt frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return awsproxy_ctrl.NewGetAWSResourceWithPost(opt, awsClients)
			},
		},
		{
			ParentRouter:  awsDeleteResourceSubRouter,
			Method:        v1.OperationDeleteImperative,
			OperationType: &v1.OperationType{Type: "AWSRESOURCE", Method: v1.OperationDeleteImperative},
			ControllerFactory: func(opts frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return awsproxy_ctrl.NewDeleteAWSResourceWithPost(opts, awsClients)
			},
		},

		// Azure Credential Handlers
		{
			ParentRouter: azureCredentialCollectionSubRouter,
			ResourceType: v20220901privatepreview.AzureCredentialType,
			Method:       v1.OperationList,
			ControllerFactory: func(opt frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return defaultoperation.NewListResources(opt,
					frontend_ctrl.ResourceOptions[datamodel.AzureCredential]{
						RequestConverter:  converter.AzureCredentialDataModelFromVersioned,
						ResponseConverter: converter.AzureCredentialDataModelToVersioned,
					},
				)
			},
		},
		{
			ParentRouter: azureCredentialResourceSubRouter,
			ResourceType: v20220901privatepreview.AzureCredentialType,
			Method:       v1.OperationGet,
			ControllerFactory: func(opt frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return defaultoperation.NewGetResource(opt,
					frontend_ctrl.ResourceOptions[datamodel.AzureCredential]{
						RequestConverter:  converter.AzureCredentialDataModelFromVersioned,
						ResponseConverter: converter.AzureCredentialDataModelToVersioned,
					},
				)
			},
		},
		{
			ParentRouter: azureCredentialResourceSubRouter,
			ResourceType: v20220901privatepreview.AzureCredentialType,
			Method:       v1.OperationPut,
			ControllerFactory: func(opt frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return azure_credential_ctrl.NewCreateOrUpdateAzureCredential(opt, secretClient)
			},
		},
		{
			ParentRouter: azureCredentialResourceSubRouter,
			ResourceType: v20220901privatepreview.AzureCredentialType,
			Method:       v1.OperationDelete,
			ControllerFactory: func(opt frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return azure_credential_ctrl.NewDeleteAzureCredential(opt, secretClient)
			},
		},

		// AWS Credential Handlers
		{
			ParentRouter: awsCredentialCollectionSubRouter,
			ResourceType: v20220901privatepreview.AWSCredentialType,
			Method:       v1.OperationList,
			ControllerFactory: func(opt frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return defaultoperation.NewListResources(opt,
					frontend_ctrl.ResourceOptions[datamodel.AWSCredential]{
						RequestConverter:  converter.AWSCredentialDataModelFromVersioned,
						ResponseConverter: converter.AWSCredentialDataModelToVersioned,
					},
				)
			},
		},
		{
			ParentRouter: awsCredentialResourceSubRouter,
			ResourceType: v20220901privatepreview.AWSCredentialType,
			Method:       v1.OperationGet,
			ControllerFactory: func(opt frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return defaultoperation.NewGetResource(opt,
					frontend_ctrl.ResourceOptions[datamodel.AWSCredential]{
						RequestConverter:  converter.AWSCredentialDataModelFromVersioned,
						ResponseConverter: converter.AWSCredentialDataModelToVersioned,
					},
				)
			},
		},
		{
			ParentRouter: awsCredentialResourceSubRouter,
			ResourceType: v20220901privatepreview.AWSCredentialType,
			Method:       v1.OperationPut,
			ControllerFactory: func(opt frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return aws_credential_ctrl.NewCreateOrUpdateAWSCredential(opt, secretClient)
			},
		},
		{
			ParentRouter: awsCredentialResourceSubRouter,
			ResourceType: v20220901privatepreview.AWSCredentialType,
			Method:       v1.OperationDelete,
			ControllerFactory: func(opt frontend_ctrl.Options) (frontend_ctrl.Controller, error) {
				return aws_credential_ctrl.NewDeleteAWSCredential(opt, secretClient)
			},
		},

		// Proxy request should take the least priority in routing and should therefore be last
		//
		// Note that the API validation is not applied to the router used for proxying
		{
			// Method deliberately omitted. This is a catch-all route for proxying.
			ParentRouter:      router.PathPrefix(fmt.Sprintf("%s%s", ctrlOpts.PathBase, planeItemPath)).Subrouter(),
			OperationType:     &v1.OperationType{Type: "UCPPROXY", Method: v1.OperationPutImperative},
			ControllerFactory: planes_ctrl.NewProxyPlane,
		},
	}...)

	for _, h := range handlerOptions {
		if err := server.RegisterHandler(ctx, h, ctrlOpts); err != nil {
			return err
		}
	}

	return nil
}
