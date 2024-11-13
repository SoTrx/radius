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

package show

import (
	"context"
	"slices"
	"strings"

	"github.com/radius-project/radius/pkg/cli"
	"github.com/radius-project/radius/pkg/cli/clients"
	"github.com/radius-project/radius/pkg/cli/clierrors"
	"github.com/radius-project/radius/pkg/cli/cmd/commonflags"
	"github.com/radius-project/radius/pkg/cli/cmd/resourcetype/common"
	"github.com/radius-project/radius/pkg/cli/connections"
	"github.com/radius-project/radius/pkg/cli/framework"
	"github.com/radius-project/radius/pkg/cli/output"
	"github.com/radius-project/radius/pkg/cli/workspaces"
	"github.com/spf13/cobra"
)

// NewCommand creates an instance of the `rad resource-type show` command and runner.
func NewCommand(factory framework.Factory) (*cobra.Command, framework.Runner) {
	runner := NewRunner(factory)

	cmd := &cobra.Command{
		Use:   "show [resource type]",
		Short: "Show resource resource type",
		Long: `Show resource resource type
		
Resource types are the entities that can be created and managed by Radius such as 'Applications.Core/containers'. Each resource type can define multiple API versions, and each API version defines a schema that resource instances conform to. Resource types can be configured using resource providers.`,
		Example: `
# Show a resource type
rad resource-type show 'Applications.Core/containers'`,
		Args: cobra.ExactArgs(1),
		RunE: framework.RunCommand(runner),
	}

	commonflags.AddOutputFlag(cmd)
	commonflags.AddWorkspaceFlag(cmd)

	return cmd, runner
}

// Runner is the Runner implementation for the `rad resource-type show` command.
type Runner struct {
	ConnectionFactory connections.Factory
	ConfigHolder      *framework.ConfigHolder
	Output            output.Interface
	Format            string
	Workspace         *workspaces.Workspace

	ResourceTypeName          string
	ResourceProviderNamespace string
	ResourceTypeSuffix        string
}

// NewRunner creates an instance of the runner for the `rad resource-type show` command.
func NewRunner(factory framework.Factory) *Runner {
	return &Runner{
		ConnectionFactory: factory.GetConnectionFactory(),
		ConfigHolder:      factory.GetConfigHolder(),
		Output:            factory.GetOutput(),
	}
}

// Validate runs validation for the `rad resource-type show` command.
func (r *Runner) Validate(cmd *cobra.Command, args []string) error {
	workspace, err := cli.RequireWorkspace(cmd, r.ConfigHolder.Config, r.ConfigHolder.DirectoryConfig)
	if err != nil {
		return err
	}
	r.Workspace = workspace

	format, err := cli.RequireOutput(cmd)
	if err != nil {
		return err
	}
	r.Format = format

	r.ResourceTypeName = args[0]
	parts := strings.Split(r.ResourceTypeName, "/")
	if len(parts) != 2 {
		return clierrors.Message("Invalid resource type %q. Expected format: '<provider>/<type>'", r.ResourceTypeName)
	}

	r.ResourceProviderNamespace = parts[0]
	r.ResourceTypeSuffix = parts[1]

	return nil
}

// Run runs the `rad resource-type show` command.
func (r *Runner) Run(ctx context.Context) error {
	client, err := r.ConnectionFactory.CreateApplicationsManagementClient(ctx, *r.Workspace)
	if err != nil {
		return err
	}

	resourceProvider, err := client.GetResourceProviderSummary(ctx, "local", r.ResourceProviderNamespace)
	if clients.Is404Error(err) {
		return clierrors.Message("The resource provider %q was not found or has been deleted.", r.ResourceProviderNamespace)
	} else if err != nil {
		return err
	}

	resourceTypes := common.ResourceTypesForProvider(&resourceProvider)
	idx := slices.IndexFunc(resourceTypes, func(rt common.ResourceType) bool {
		return rt.Name == r.ResourceTypeName
	})

	if idx < 0 {
		return clierrors.Message("Resource type %q not found in resource provider %q.", r.ResourceTypeSuffix, r.ResourceProviderNamespace)
	}

	err = r.Output.WriteFormatted(r.Format, resourceTypes[idx], common.GetResourceTypeTableFormat())
	if err != nil {
		return err
	}

	return nil
}