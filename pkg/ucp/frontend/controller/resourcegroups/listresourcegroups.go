// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------
package resourcegroups

import (
	"context"
	"fmt"
	http "net/http"

	armrpc_rest "github.com/project-radius/radius/pkg/armrpc/rest"
	"github.com/project-radius/radius/pkg/middleware"
	ctrl "github.com/project-radius/radius/pkg/ucp/frontend/controller"
	"github.com/project-radius/radius/pkg/ucp/resources"
	"github.com/project-radius/radius/pkg/ucp/rest"
	"github.com/project-radius/radius/pkg/ucp/store"
	"github.com/project-radius/radius/pkg/ucp/ucplog"
)

var _ ctrl.Controller = (*ListResourceGroups)(nil)

// ListResourceGroups is the controller implementation to get the list of UCP resource groups.
type ListResourceGroups struct {
	ctrl.BaseController
}

// NewListResourceGroups creates a new ListResourceGroups.
func NewListResourceGroups(opts ctrl.Options) (ctrl.Controller, error) {
	return &ListResourceGroups{ctrl.NewBaseController(opts)}, nil
}

func (r *ListResourceGroups) Run(ctx context.Context, w http.ResponseWriter, req *http.Request) (armrpc_rest.Response, error) {
	path := middleware.GetRelativePath(r.Options.BasePath, req.URL.Path)
	logger := ucplog.GetLogger(ctx)
	var query store.Query
	planeType, planeName, _, err := resources.ExtractPlanesPrefixFromURLPath(path)
	if err != nil {
		return nil, err
	}
	query.RootScope = resources.SegmentSeparator + resources.PlanesSegment + resources.SegmentSeparator + planeType + resources.SegmentSeparator + planeName
	query.IsScopeQuery = true
	query.ResourceType = "resourcegroups"
	logger.Info(fmt.Sprintf("Listing resource groups in scope %s", query.RootScope))

	result, err := r.StorageClient().Query(ctx, query)
	if err != nil {
		return nil, err
	}
	listOfResourceGroups, err := r.createResponse(ctx, req, result)
	if err != nil {
		return nil, err
	}

	var ok = armrpc_rest.NewOKResponse(listOfResourceGroups)
	return ok, nil
}

func (e *ListResourceGroups) createResponse(ctx context.Context, req *http.Request, result *store.ObjectQueryResult) (rest.ResourceGroupList, error) {
	listOfResourceGroups := rest.ResourceGroupList{}
	if result != nil && len(result.Items) > 0 {
		for _, item := range result.Items {
			var rg rest.ResourceGroup
			err := item.As(&rg)
			if err != nil {
				return listOfResourceGroups, err
			}
			listOfResourceGroups.Value = append(listOfResourceGroups.Value, rg)
		}
	}
	return listOfResourceGroups, nil
}
