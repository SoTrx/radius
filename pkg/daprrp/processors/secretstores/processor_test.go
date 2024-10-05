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

package secretstores

import (
	"context"
	"testing"

	v1 "github.com/radius-project/radius/pkg/armrpc/api/v1"
	"github.com/radius-project/radius/pkg/daprrp/datamodel"
	dapr_ctrl "github.com/radius-project/radius/pkg/daprrp/frontend/controller"
	"github.com/radius-project/radius/pkg/kubernetes"
	"github.com/radius-project/radius/pkg/portableresources"
	"github.com/radius-project/radius/pkg/portableresources/processors"
	"github.com/radius-project/radius/pkg/portableresources/renderers/dapr"
	"github.com/radius-project/radius/pkg/recipes"
	rpv1 "github.com/radius-project/radius/pkg/rp/v1"
	"github.com/radius-project/radius/pkg/to"
	"github.com/radius-project/radius/test/k8sutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/kubectl/pkg/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Test_Process(t *testing.T) {
	const kubernetesResource = "/planes/kubernetes/local/namespaces/test-namespace/providers/dapr.io/Component/test-component"
	const applicationID = "/planes/radius/local/resourceGroups/test-rg/providers/Applications.Core/applications/test-app"
	const envID = "/planes/radius/local/resourceGroups/test-rg/providers/Applications.Core/environments/test-env"
	const componentName = "test-component"

	t.Run("success - recipe", func(t *testing.T) {
		processor := Processor{
			Client: k8sutil.NewFakeKubeClient(scheme.Scheme),
		}

		resource := &datamodel.DaprSecretStore{
			BaseResource: v1.BaseResource{
				TrackedResource: v1.TrackedResource{
					Name: componentName,
				},
			},
			Properties: datamodel.DaprSecretStoreProperties{
				BasicResourceProperties: rpv1.BasicResourceProperties{
					Application: applicationID,
				},
				BasicDaprResourceProperties: rpv1.BasicDaprResourceProperties{
					ComponentName: componentName,
				},
			},
		}
		options := processors.Options{
			RuntimeConfiguration: recipes.RuntimeConfiguration{
				Kubernetes: &recipes.KubernetesRuntime{
					Namespace: "test-namespace",
				},
			},
			RecipeOutput: &recipes.RecipeOutput{
				Resources: []string{
					kubernetesResource,
				},
				Values:  map[string]any{}, // Component name will be computed for resource name.
				Secrets: map[string]any{},
			},
		}

		err := processor.Process(context.Background(), resource, options)
		require.NoError(t, err)

		require.Equal(t, componentName, resource.Properties.ComponentName)

		expectedValues := map[string]any{
			"componentName": componentName,
		}
		expectedSecrets := map[string]rpv1.SecretValueReference{}

		expectedOutputResources, err := processors.GetOutputResourcesFromRecipe(options.RecipeOutput)
		require.NoError(t, err)

		require.Equal(t, expectedValues, resource.ComputedValues)
		require.Equal(t, expectedSecrets, resource.SecretValues)
		require.Equal(t, expectedOutputResources, resource.Properties.Status.OutputResources)

		components := unstructured.UnstructuredList{}
		components.SetAPIVersion("dapr.io/v1alpha1")
		components.SetKind("Component")

		// No components created for a recipe
		err = processor.Client.List(context.Background(), &components, &client.ListOptions{Namespace: options.RuntimeConfiguration.Kubernetes.Namespace})
		require.NoError(t, err)
		require.Empty(t, components.Items)
	})

	t.Run("success - values", func(t *testing.T) {
		processor := Processor{
			Client: k8sutil.NewFakeKubeClient(scheme.Scheme, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test-namespace"}}),
		}

		resource := &datamodel.DaprSecretStore{
			BaseResource: v1.BaseResource{
				TrackedResource: v1.TrackedResource{
					Name: "some-other-name",
				},
			},
			Properties: datamodel.DaprSecretStoreProperties{
				BasicResourceProperties: rpv1.BasicResourceProperties{
					Application: applicationID,
				},
				BasicDaprResourceProperties: rpv1.BasicDaprResourceProperties{
					ComponentName: componentName,
				},
				ResourceProvisioning: portableresources.ResourceProvisioningManual,
				Metadata: map[string]*rpv1.DaprComponentMetadataValue{
					"config": {
						Value: "extrasecure",
					},
				},
				Type:    "secretstores.kubernetes",
				Version: "v1",
				Scopes:  []string{"test-scope-1"},
			},
		}

		options := processors.Options{
			RuntimeConfiguration: recipes.RuntimeConfiguration{
				Kubernetes: &recipes.KubernetesRuntime{
					Namespace: "test-namespace",
				},
			},
		}

		err := processor.Process(context.Background(), resource, options)
		require.NoError(t, err)

		require.Equal(t, componentName, resource.Properties.ComponentName)

		expectedValues := map[string]any{
			"componentName": componentName,
		}
		expectedSecrets := map[string]rpv1.SecretValueReference{}
		generated := &unstructured.Unstructured{
			Object: map[string]any{
				"apiVersion": dapr.DaprAPIVersion,
				"kind":       dapr.DaprKind,
				"metadata": map[string]any{
					"namespace":       "test-namespace",
					"name":            "test-component",
					"labels":          kubernetes.MakeDescriptiveDaprLabels("test-app", "some-other-name", dapr_ctrl.DaprSecretStoresResourceType),
					"resourceVersion": "1",
				},
				"spec": map[string]any{
					"type":    "secretstores.kubernetes",
					"version": "v1",
					"metadata": []any{
						map[string]any{
							"name":  "config",
							"value": "extrasecure",
						},
					},
				},
				"scopes": []any{"test-scope-1"},
			},
		}

		component := rpv1.NewKubernetesOutputResource("Component", generated, metav1.ObjectMeta{Name: generated.GetName(), Namespace: generated.GetNamespace()})
		component.RadiusManaged = to.Ptr(true)
		expectedOutputResources := []rpv1.OutputResource{component}
		require.NoError(t, err)

		require.Equal(t, expectedValues, resource.ComputedValues)
		require.Equal(t, expectedSecrets, resource.SecretValues)
		require.Equal(t, expectedOutputResources, resource.Properties.Status.OutputResources)

		components := unstructured.UnstructuredList{}
		components.SetAPIVersion("dapr.io/v1alpha1")
		components.SetKind("Component")
		err = processor.Client.List(context.Background(), &components, &client.ListOptions{Namespace: options.RuntimeConfiguration.Kubernetes.Namespace})
		require.NoError(t, err)
		require.NotEmpty(t, components.Items)
		require.Equal(t, []unstructured.Unstructured{*generated}, components.Items)
	})

	t.Run("success - values (no application)", func(t *testing.T) {
		processor := Processor{
			Client: k8sutil.NewFakeKubeClient(scheme.Scheme, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test-namespace"}}),
		}

		resource := &datamodel.DaprSecretStore{
			BaseResource: v1.BaseResource{
				TrackedResource: v1.TrackedResource{
					Name: "some-other-name",
				},
			},
			Properties: datamodel.DaprSecretStoreProperties{
				BasicResourceProperties: rpv1.BasicResourceProperties{
					Environment: envID,
				},
				BasicDaprResourceProperties: rpv1.BasicDaprResourceProperties{
					ComponentName: componentName,
				},
				ResourceProvisioning: portableresources.ResourceProvisioningManual,
				Metadata: map[string]*rpv1.DaprComponentMetadataValue{
					"config": {
						Value: "extrasecure",
					},
				},
				Type:    "secretstores.kubernetes",
				Version: "v1",
			},
		}

		options := processors.Options{
			RuntimeConfiguration: recipes.RuntimeConfiguration{
				Kubernetes: &recipes.KubernetesRuntime{
					Namespace: "test-namespace",
				},
			},
		}

		err := processor.Process(context.Background(), resource, options)
		require.NoError(t, err)

		require.Equal(t, componentName, resource.Properties.ComponentName)

		expectedValues := map[string]any{
			"componentName": componentName,
		}
		expectedSecrets := map[string]rpv1.SecretValueReference{}
		generated := &unstructured.Unstructured{
			Object: map[string]any{
				"apiVersion": dapr.DaprAPIVersion,
				"kind":       dapr.DaprKind,
				"metadata": map[string]any{
					"namespace":       "test-namespace",
					"name":            "test-component",
					"labels":          kubernetes.MakeDescriptiveDaprLabels("", "some-other-name", dapr_ctrl.DaprSecretStoresResourceType),
					"resourceVersion": "1",
				},
				"spec": map[string]any{
					"type":    "secretstores.kubernetes",
					"version": "v1",
					"metadata": []any{
						map[string]any{
							"name":  "config",
							"value": "extrasecure",
						},
					},
				},
			},
		}

		component := rpv1.NewKubernetesOutputResource("Component", generated, metav1.ObjectMeta{Name: generated.GetName(), Namespace: generated.GetNamespace()})
		component.RadiusManaged = to.Ptr(true)
		expectedOutputResources := []rpv1.OutputResource{component}
		require.NoError(t, err)

		require.Equal(t, expectedValues, resource.ComputedValues)
		require.Equal(t, expectedSecrets, resource.SecretValues)
		require.Equal(t, expectedOutputResources, resource.Properties.Status.OutputResources)

		components := unstructured.UnstructuredList{}
		components.SetAPIVersion("dapr.io/v1alpha1")
		components.SetKind("Component")
		err = processor.Client.List(context.Background(), &components, &client.ListOptions{Namespace: options.RuntimeConfiguration.Kubernetes.Namespace})
		require.NoError(t, err)
		require.NotEmpty(t, components.Items)
		require.Equal(t, []unstructured.Unstructured{*generated}, components.Items)
	})

	t.Run("success - recipe with value overrides", func(t *testing.T) {
		processor := Processor{
			Client: k8sutil.NewFakeKubeClient(scheme.Scheme),
		}

		resource := &datamodel.DaprSecretStore{
			BaseResource: v1.BaseResource{
				TrackedResource: v1.TrackedResource{
					Name: "some-other-name",
				},
			},
			Properties: datamodel.DaprSecretStoreProperties{
				BasicDaprResourceProperties: rpv1.BasicDaprResourceProperties{
					ComponentName: componentName,
				},
			},
		}
		options := processors.Options{
			RuntimeConfiguration: recipes.RuntimeConfiguration{
				Kubernetes: &recipes.KubernetesRuntime{
					Namespace: "test-namespace",
				},
			},
			RecipeOutput: &recipes.RecipeOutput{
				Resources: []string{
					kubernetesResource,
				},

				// Values and secrets will be overridden by the resource.
				Values: map[string]any{
					"componentName": "akskdf",
				},
				Secrets: map[string]any{},
			},
		}

		err := processor.Process(context.Background(), resource, options)
		require.NoError(t, err)

		require.Equal(t, componentName, resource.Properties.ComponentName)

		expectedValues := map[string]any{
			"componentName": componentName,
		}
		expectedSecrets := map[string]rpv1.SecretValueReference{}

		expectedOutputResources := []rpv1.OutputResource{}

		recipeOutputResources, err := processors.GetOutputResourcesFromRecipe(options.RecipeOutput)
		require.NoError(t, err)
		expectedOutputResources = append(expectedOutputResources, recipeOutputResources...)
		require.Equal(t, expectedValues, resource.ComputedValues)
		require.Equal(t, expectedSecrets, resource.SecretValues)
		require.Equal(t, expectedOutputResources, resource.Properties.Status.OutputResources)

		components := unstructured.UnstructuredList{}
		components.SetAPIVersion("dapr.io/v1alpha1")
		components.SetKind("Component")
		err = processor.Client.List(context.Background(), &components, &client.ListOptions{Namespace: options.RuntimeConfiguration.Kubernetes.Namespace})
		require.NoError(t, err)
		require.Empty(t, components.Items)
	})

	t.Run("failure - duplicate component", func(t *testing.T) {
		// Create a duplicate with the same component name.
		existing, err := dapr.ConstructDaprGeneric(
			dapr.DaprGeneric{
				Type:     to.Ptr("secretstores.kubernetes"),
				Version:  to.Ptr("v1"),
				Metadata: map[string]*rpv1.DaprComponentMetadataValue{},
			},
			"test-namespace",
			"test-component",
			"test-app",
			"some-other-other-name",
			dapr_ctrl.DaprSecretStoresResourceType)
		require.NoError(t, err)

		processor := Processor{
			Client: k8sutil.NewFakeKubeClient(scheme.Scheme, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test-namespace"}}, &existing),
		}
		resource := &datamodel.DaprSecretStore{
			BaseResource: v1.BaseResource{
				TrackedResource: v1.TrackedResource{
					Name: "some-other-name",
				},
			},
			Properties: datamodel.DaprSecretStoreProperties{
				BasicResourceProperties: rpv1.BasicResourceProperties{
					Application: applicationID,
				},
				BasicDaprResourceProperties: rpv1.BasicDaprResourceProperties{
					ComponentName: componentName,
				},
				ResourceProvisioning: portableresources.ResourceProvisioningManual,
				Metadata: map[string]*rpv1.DaprComponentMetadataValue{
					"config": {
						Value: "extrasecure",
					},
				},
				Type:    "secretstores.kubernetes",
				Version: "v1",
			},
		}

		options := processors.Options{
			RuntimeConfiguration: recipes.RuntimeConfiguration{
				Kubernetes: &recipes.KubernetesRuntime{
					Namespace: "test-namespace",
				},
			},
		}

		err = processor.Process(context.Background(), resource, options)
		require.Error(t, err)
		assert.IsType(t, &processors.ValidationError{}, err)
		assert.Equal(t, "the Dapr component name '\"test-component\"' is already in use by another resource. Dapr component and resource names must be unique across all Dapr types (e.g., StateStores, PubSubBrokers, SecretStores, ConfigurationStores, etc.). Please select a new name and try again.", err.Error())
	})
}
