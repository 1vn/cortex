/*
Copyright 2019 Cortex Labs, Inc.

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

package context

import (
	"fmt"

	"github.com/cortexlabs/cortex/pkg/lib/errors"
	"github.com/cortexlabs/cortex/pkg/lib/sets/strset"
	"github.com/cortexlabs/cortex/pkg/operator/api/resource"
	userconfig "github.com/cortexlabs/cortex/pkg/operator/api/userconfig"
	"github.com/cortexlabs/cortex/pkg/operator/config"
)

type Context struct {
	ID                string               `json:"id"`
	Key               string               `json:"key"`
	CreatedEpoch      int64                `json:"created_epoch"`
	CortexConfig      *config.CortexConfig `json:"cortex_config"`
	DeploymentVersion string               `json:"deployment_version"`
	Root              string               `json:"root"`
	MetadataRoot      string               `json:"metadata_root"`
	StatusPrefix      string               `json:"status_prefix"`
	App               *App                 `json:"app"`
	APIs              APIs                 `json:"apis"`
	ProjectID         string               `json:"project_id"`
	ProjectKey        string               `json:"project_key"`
}

type Resource interface {
	userconfig.Resource
	GetID() string
}

type ComputedResource interface {
	Resource
	GetWorkloadID() string
	SetWorkloadID(string)
}

type ResourceFields struct {
	ID           string        `json:"id"`
	ResourceType resource.Type `json:"resource_type"`
}

type ComputedResourceFields struct {
	*ResourceFields
	WorkloadID string `json:"workload_id"`
}

func (r *ResourceFields) GetID() string {
	return r.ID
}

func (r *ComputedResourceFields) GetWorkloadID() string {
	return r.WorkloadID
}

func (r *ComputedResourceFields) SetWorkloadID(workloadID string) {
	r.WorkloadID = workloadID
}

func ExtractResourceWorkloadIDs(resources []ComputedResource) map[string]string {
	resourceWorkloadIDs := make(map[string]string, len(resources))
	for _, res := range resources {
		resourceWorkloadIDs[res.GetID()] = res.GetWorkloadID()
	}
	return resourceWorkloadIDs
}

func (ctx *Context) DataComputedResources() []ComputedResource {
	var resources []ComputedResource
	return resources
}

func (ctx *Context) APIResources() []ComputedResource {
	resources := make([]ComputedResource, len(ctx.APIs))
	i := 0
	for _, api := range ctx.APIs {
		resources[i] = api
		i++
	}
	return resources
}

func (ctx *Context) ComputedResources() []ComputedResource {
	return append(ctx.DataComputedResources(), ctx.APIResources()...)
}

func (ctx *Context) AllResources() []Resource {
	var resources []Resource
	for _, res := range ctx.ComputedResources() {
		resources = append(resources, res)
	}
	return resources
}

func (ctx *Context) ComputedResourceIDs() strset.Set {
	resourceIDs := strset.New()
	for _, res := range ctx.ComputedResources() {
		resourceIDs.Add(res.GetID())
	}
	return resourceIDs
}

func (ctx *Context) DataResourceWorkloadIDs() map[string]string {
	return ExtractResourceWorkloadIDs(ctx.DataComputedResources())
}

func (ctx *Context) APIResourceWorkloadIDs() map[string]string {
	return ExtractResourceWorkloadIDs(ctx.APIResources())
}

func (ctx *Context) ComputedResourceResourceWorkloadIDs() map[string]string {
	return ExtractResourceWorkloadIDs(ctx.ComputedResources())
}

func (ctx *Context) ComputedResourceWorkloadIDs() strset.Set {
	workloadIDs := strset.New()
	for _, workloadID := range ExtractResourceWorkloadIDs(ctx.ComputedResources()) {
		workloadIDs.Add(workloadID)
	}
	return workloadIDs
}

// Note: there may be >1 resources with the ID, this returns one of them
func (ctx *Context) OneResourceByID(resourceID string) Resource {
	for _, res := range ctx.AllResources() {
		if res.GetID() == resourceID {
			return res
		}
	}
	return nil
}

func (ctx *Context) AllResourcesByName(name string) []Resource {
	var resources []Resource
	for _, res := range ctx.AllResources() {
		if res.GetName() == name {
			resources = append(resources, res)
		}
	}
	return resources
}

func (ctx *Context) CheckAllWorkloadIDsPopulated() error {
	for _, res := range ctx.ComputedResources() {
		if res.GetWorkloadID() == "" {
			return errors.New(ctx.App.Name, "workload ID missing", fmt.Sprintf("%s (ID: %s)", res.GetName(), res.GetID())) // unexpected
		}
	}
	return nil
}

func (ctx *Context) VisibleResourcesMap() map[string][]ComputedResource {
	resources := make(map[string][]ComputedResource)
	for name, api := range ctx.APIs {
		resources[name] = append(resources[name], api)
	}
	return resources
}

func (ctx *Context) VisibleResourcesByName(name string) []ComputedResource {
	return ctx.VisibleResourcesMap()[name]
}

func (ctx *Context) VisibleResourceByName(name string) (ComputedResource, error) {
	resources := ctx.VisibleResourcesByName(name)
	if len(resources) == 0 {
		return nil, resource.ErrorNameNotFound(name)
	}
	if len(resources) > 1 {
		validStrs := make([]string, len(resources))
		for i, resource := range resources {
			resourceTypeStr := resource.GetResourceType().String()
			validStrs[i] = resourceTypeStr + " " + name
		}
		return nil, resource.ErrorBeMoreSpecific(validStrs...)
	}
	return resources[0], nil
}

func (ctx *Context) VisibleResourceByNameAndType(name string, resourceTypeStr string) (ComputedResource, error) {
	resourceType := resource.TypeFromString(resourceTypeStr)

	switch resourceType {
	case resource.APIType:
		res := ctx.APIs[name]
		if res == nil {
			return nil, resource.ErrorNotFound(name, resourceType)
		}
		return res, nil
	}

	return nil, resource.ErrorInvalidType(resourceTypeStr)
}

func (ctx *Context) Validate() error {
	return nil
}
