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

package controller

import (
	"context"
	"errors"

	"github.com/radius-project/radius/pkg/corerp/backend/deployment"
	"github.com/radius-project/radius/pkg/ucp/store"

	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// Options represents controller options.
type Options struct {
	// StorageClient is the data storage client.
	StorageClient store.StorageClient

	// KubeClient is the Kubernetes controller runtime client.
	KubeClient runtimeclient.Client

	// ResourceType is the string that represents the resource type.
	ResourceType string

	// GetDeploymentProcessor is the factory function to create core rp DeploymentProcessor instance.
	GetDeploymentProcessor func() deployment.DeploymentProcessor
}

// Validate validates that required fields are set on the options.
func (o Options) Validate() error {
	var err error
	if o.StorageClient == nil {
		err = errors.Join(err, errors.New("StorageClient is required"))
	}
	if o.ResourceType == "" {
		err = errors.Join(err, errors.New("ResourceType is required"))
	}

	// KubeClient and GetDeploymentProcessor are not used by the majority of the code, so they
	// are not validated here.

	return err
}

// Controller is an interface to implement async operation controller.
type Controller interface {
	// Run runs async request operation.
	Run(ctx context.Context, request *Request) (Result, error)

	// StorageClient gets the storage client for resource type.
	StorageClient() store.StorageClient
}

// BaseController is the base struct of async operation controller.
type BaseController struct {
	options Options
}

// NewBaseAsyncController creates a new BaseController instance with the given Options for Async Operation.
func NewBaseAsyncController(options Options) BaseController {
	return BaseController{options}
}

// StorageClient gets storage client for this controller.
func (b *BaseController) StorageClient() store.StorageClient {
	return b.options.StorageClient
}

// KubeClient gets Kubernetes client for this controller.
func (b *BaseController) KubeClient() runtimeclient.Client {
	return b.options.KubeClient
}

// ResourceType gets the resource type for this controller.
func (b *BaseController) ResourceType() string {
	return b.options.ResourceType
}

// DeploymentProcessor gets the core rp deployment processor for this controller.
func (b *BaseController) DeploymentProcessor() deployment.DeploymentProcessor {
	return b.options.GetDeploymentProcessor()
}
