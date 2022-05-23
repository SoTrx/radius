//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package v20220315privatepreview

import (
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"reflect"
	"time"
)

// BasicResourceProperties - Basic properties of a Radius resource.
type BasicResourceProperties struct {
	// Status of the resource
	Status *ResourceStatus `json:"status,omitempty"`
}

// DaprSecretStoreList - Object that includes an array of DaprSecretStore and a possible link for next set
type DaprSecretStoreList struct {
	// The link used to fetch the next page of DaprSecretStore list.
	NextLink *string `json:"nextLink,omitempty"`

	// List of DaprSecretStore resources
	Value []*DaprSecretStoreResource `json:"value,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type DaprSecretStoreList.
func (d DaprSecretStoreList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", d.NextLink)
	populate(objectMap, "value", d.Value)
	return json.Marshal(objectMap)
}

// DaprSecretStoreProperties - DaprSecretStore connector properties
type DaprSecretStoreProperties struct {
	BasicResourceProperties
	// REQUIRED; Fully qualified resource ID for the environment that the connector is linked to
	Environment *string `json:"environment,omitempty"`

	// REQUIRED; Radius kind for Dapr Secret Store
	Kind *DaprSecretStorePropertiesKind `json:"kind,omitempty"`

	// REQUIRED; Metadata for the Secret Store resource. This should match the values specified in Dapr component spec
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// REQUIRED; Dapr Secret Store type. These strings match the types defined in Dapr Component format: https://docs.dapr.io/reference/components-reference/supported-secret-stores/
	Type *string `json:"type,omitempty"`

	// REQUIRED; Dapr component version
	Version *string `json:"version,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the application that the connector is consumed by
	Application *string `json:"application,omitempty" azure:"ro"`

	// READ-ONLY; Provisioning state of the dapr secret store connector at the time the operation was called
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`
}

// DaprSecretStoreResource - DaprSecretStore connector
type DaprSecretStoreResource struct {
	TrackedResource
	// REQUIRED; DaprSecretStore connector properties
	Properties *DaprSecretStoreProperties `json:"properties,omitempty"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type DaprSecretStoreResource.
func (d DaprSecretStoreResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	d.TrackedResource.marshalInternal(objectMap)
	populate(objectMap, "properties", d.Properties)
	populate(objectMap, "systemData", d.SystemData)
	return json.Marshal(objectMap)
}

// DaprSecretStoresCreateOrUpdateOptions contains the optional parameters for the DaprSecretStores.CreateOrUpdate method.
type DaprSecretStoresCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// DaprSecretStoresDeleteOptions contains the optional parameters for the DaprSecretStores.Delete method.
type DaprSecretStoresDeleteOptions struct {
	// placeholder for future optional parameters
}

// DaprSecretStoresGetOptions contains the optional parameters for the DaprSecretStores.Get method.
type DaprSecretStoresGetOptions struct {
	// placeholder for future optional parameters
}

// DaprSecretStoresListBySubscriptionOptions contains the optional parameters for the DaprSecretStores.ListBySubscription method.
type DaprSecretStoresListBySubscriptionOptions struct {
	// placeholder for future optional parameters
}

// DaprSecretStoresListOptions contains the optional parameters for the DaprSecretStores.List method.
type DaprSecretStoresListOptions struct {
	// placeholder for future optional parameters
}

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
	// READ-ONLY; The additional info.
	Info map[string]interface{} `json:"info,omitempty" azure:"ro"`

	// READ-ONLY; The additional info type.
	Type *string `json:"type,omitempty" azure:"ro"`
}

// ErrorDetail - The error detail.
type ErrorDetail struct {
	// READ-ONLY; The error additional info.
	AdditionalInfo []*ErrorAdditionalInfo `json:"additionalInfo,omitempty" azure:"ro"`

	// READ-ONLY; The error code.
	Code *string `json:"code,omitempty" azure:"ro"`

	// READ-ONLY; The error details.
	Details []*ErrorDetail `json:"details,omitempty" azure:"ro"`

	// READ-ONLY; The error message.
	Message *string `json:"message,omitempty" azure:"ro"`

	// READ-ONLY; The error target.
	Target *string `json:"target,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type ErrorDetail.
func (e ErrorDetail) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "additionalInfo", e.AdditionalInfo)
	populate(objectMap, "code", e.Code)
	populate(objectMap, "details", e.Details)
	populate(objectMap, "message", e.Message)
	populate(objectMap, "target", e.Target)
	return json.Marshal(objectMap)
}

// ErrorResponse - Common error response for all Azure Resource Manager APIs to return error details for failed operations. (This also follows the OData
// error response format.).
// Implements the error and azcore.HTTPResponse interfaces.
type ErrorResponse struct {
	raw string
	// The error object.
	InnerError *ErrorDetail `json:"error,omitempty"`
}

// Error implements the error interface for type ErrorResponse.
// The contents of the error text are not contractual and subject to change.
func (e ErrorResponse) Error() string {
	return e.raw
}

// MongoDatabaseList - Object that includes an array of MongoDatabase and a possible link for next set
type MongoDatabaseList struct {
	// The link used to fetch the next page of MongoDatabase list.
	NextLink *string `json:"nextLink,omitempty"`

	// List of MongoDatabase resources
	Value []*MongoDatabaseResource `json:"value,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type MongoDatabaseList.
func (m MongoDatabaseList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", m.NextLink)
	populate(objectMap, "value", m.Value)
	return json.Marshal(objectMap)
}

// MongoDatabaseProperties - MongoDatabse connector properties
type MongoDatabaseProperties struct {
	BasicResourceProperties
	// REQUIRED; Fully qualified resource ID for the environment that the connector is linked to
	Environment *string `json:"environment,omitempty"`

	// Host name of the target Mongo database
	Host *string `json:"host,omitempty"`

	// Port value of the target Mongo database
	Port *int32 `json:"port,omitempty"`

	// Fully qualified resource ID of a supported resource with Mongo API to use for this connector
	Resource *string `json:"resource,omitempty"`

	// Secrets values provided for the resource
	Secrets *MongoDatabasePropertiesSecrets `json:"secrets,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the application that the connector is consumed by
	Application *string `json:"application,omitempty" azure:"ro"`

	// READ-ONLY; Provisioning state of the mongo database connector at the time the operation was called
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`
}

// MongoDatabasePropertiesSecrets - Secrets values provided for the resource
type MongoDatabasePropertiesSecrets struct {
	// Connection string used to connect to the target Mongo database
	ConnectionString *string `json:"connectionString,omitempty"`

	// Password to use when connecting to the target Mongo database
	Password *string `json:"password,omitempty"`

	// Username to use when connecting to the target Mongo database
	Username *string `json:"username,omitempty"`
}

// MongoDatabaseResource - MongoDatabse connector
type MongoDatabaseResource struct {
	TrackedResource
	// REQUIRED; MongoDatabse connector properties
	Properties *MongoDatabaseProperties `json:"properties,omitempty"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type MongoDatabaseResource.
func (m MongoDatabaseResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	m.TrackedResource.marshalInternal(objectMap)
	populate(objectMap, "properties", m.Properties)
	populate(objectMap, "systemData", m.SystemData)
	return json.Marshal(objectMap)
}

// MongoDatabasesCreateOrUpdateOptions contains the optional parameters for the MongoDatabases.CreateOrUpdate method.
type MongoDatabasesCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// MongoDatabasesDeleteOptions contains the optional parameters for the MongoDatabases.Delete method.
type MongoDatabasesDeleteOptions struct {
	// placeholder for future optional parameters
}

// MongoDatabasesGetOptions contains the optional parameters for the MongoDatabases.Get method.
type MongoDatabasesGetOptions struct {
	// placeholder for future optional parameters
}

// MongoDatabasesListBySubscriptionOptions contains the optional parameters for the MongoDatabases.ListBySubscription method.
type MongoDatabasesListBySubscriptionOptions struct {
	// placeholder for future optional parameters
}

// MongoDatabasesListOptions contains the optional parameters for the MongoDatabases.List method.
type MongoDatabasesListOptions struct {
	// placeholder for future optional parameters
}

// RedisCacheList - Object that includes an array of RedisCache and a possible link for next set
type RedisCacheList struct {
	// The link used to fetch the next page of RedisCache list.
	NextLink *string `json:"nextLink,omitempty"`

	// List of RedisCache resources
	Value []*RedisCacheResource `json:"value,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type RedisCacheList.
func (r RedisCacheList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", r.NextLink)
	populate(objectMap, "value", r.Value)
	return json.Marshal(objectMap)
}

// RedisCacheProperties - RedisCache connector properties
type RedisCacheProperties struct {
	BasicResourceProperties
	// REQUIRED; Fully qualified resource ID for the environment that the connector is linked to
	Environment *string `json:"environment,omitempty"`

	// The host name of the target redis cache
	Host *string `json:"host,omitempty"`

	// The port value of the target redis cache
	Port *int32 `json:"port,omitempty"`

	// Fully qualified resource ID of a supported resource with Redis API to use for this connector
	Resource *string `json:"resource,omitempty"`
	Secrets *RedisCachePropertiesSecrets `json:"secrets,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the application that the connector is consumed by
	Application *string `json:"application,omitempty" azure:"ro"`

	// READ-ONLY; Provisioning state of the redis cache connector at the time the operation was called
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`
}

type RedisCachePropertiesSecrets struct {
	// The Redis connection string used to connect to the redis cache
	ConnectionString *string `json:"connectionString,omitempty"`

	// The password for this Redis instance
	Password *string `json:"password,omitempty"`
}

// RedisCacheResource - RedisCache connector
type RedisCacheResource struct {
	TrackedResource
	// REQUIRED; RedisCache connector properties
	Properties *RedisCacheProperties `json:"properties,omitempty"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type RedisCacheResource.
func (r RedisCacheResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	r.TrackedResource.marshalInternal(objectMap)
	populate(objectMap, "properties", r.Properties)
	populate(objectMap, "systemData", r.SystemData)
	return json.Marshal(objectMap)
}

// RedisCachesCreateOrUpdateOptions contains the optional parameters for the RedisCaches.CreateOrUpdate method.
type RedisCachesCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// RedisCachesDeleteOptions contains the optional parameters for the RedisCaches.Delete method.
type RedisCachesDeleteOptions struct {
	// placeholder for future optional parameters
}

// RedisCachesGetOptions contains the optional parameters for the RedisCaches.Get method.
type RedisCachesGetOptions struct {
	// placeholder for future optional parameters
}

// RedisCachesListBySubscriptionOptions contains the optional parameters for the RedisCaches.ListBySubscription method.
type RedisCachesListBySubscriptionOptions struct {
	// placeholder for future optional parameters
}

// RedisCachesListOptions contains the optional parameters for the RedisCaches.List method.
type RedisCachesListOptions struct {
	// placeholder for future optional parameters
}

// Resource - Common fields that are returned in the response for all Azure Resource Manager resources
type Resource struct {
	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type Resource.
func (r Resource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	r.marshalInternal(objectMap)
	return json.Marshal(objectMap)
}

func (r Resource) marshalInternal(objectMap map[string]interface{}) {
	populate(objectMap, "id", r.ID)
	populate(objectMap, "name", r.Name)
	populate(objectMap, "type", r.Type)
}

// ResourceStatus - Status of a resource.
type ResourceStatus struct {
	OutputResources []map[string]interface{} `json:"outputResources,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type ResourceStatus.
func (r ResourceStatus) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "outputResources", r.OutputResources)
	return json.Marshal(objectMap)
}

// SQLDatabaseList - Object that includes an array of SQLDatabase and a possible link for next set
type SQLDatabaseList struct {
	// The link used to fetch the next page of SQLDatabase list.
	NextLink *string `json:"nextLink,omitempty"`

	// List of SQLDatabase resources
	Value []*SQLDatabaseResource `json:"value,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type SQLDatabaseList.
func (s SQLDatabaseList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", s.NextLink)
	populate(objectMap, "value", s.Value)
	return json.Marshal(objectMap)
}

// SQLDatabaseProperties - SQLDatabse connector properties
type SQLDatabaseProperties struct {
	BasicResourceProperties
	// REQUIRED; The resource id of the environment linked to the sqlDatabase connector
	Environment *string `json:"environment,omitempty"`

	// The name of the SQL database.
	Database *string `json:"database,omitempty"`

	// Fully qualified resource ID of a supported resource with SQL API to use for this connector
	Resource *string `json:"resource,omitempty"`

	// The fully qualified domain name of the SQL database.
	Server *string `json:"server,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the environment that the connector is linked to
	Application *string `json:"application,omitempty" azure:"ro"`

	// READ-ONLY; Provisioning state of the SQL database connector at the time the operation was called
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`
}

// SQLDatabaseResource - SQLDatabse connector
type SQLDatabaseResource struct {
	TrackedResource
	// REQUIRED; SQLDatabse connector properties
	Properties *SQLDatabaseProperties `json:"properties,omitempty"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type SQLDatabaseResource.
func (s SQLDatabaseResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	s.TrackedResource.marshalInternal(objectMap)
	populate(objectMap, "properties", s.Properties)
	populate(objectMap, "systemData", s.SystemData)
	return json.Marshal(objectMap)
}

// SQLDatabasesCreateOrUpdateOptions contains the optional parameters for the SQLDatabases.CreateOrUpdate method.
type SQLDatabasesCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// SQLDatabasesDeleteOptions contains the optional parameters for the SQLDatabases.Delete method.
type SQLDatabasesDeleteOptions struct {
	// placeholder for future optional parameters
}

// SQLDatabasesGetOptions contains the optional parameters for the SQLDatabases.Get method.
type SQLDatabasesGetOptions struct {
	// placeholder for future optional parameters
}

// SQLDatabasesListBySubscriptionOptions contains the optional parameters for the SQLDatabases.ListBySubscription method.
type SQLDatabasesListBySubscriptionOptions struct {
	// placeholder for future optional parameters
}

// SQLDatabasesListOptions contains the optional parameters for the SQLDatabases.List method.
type SQLDatabasesListOptions struct {
	// placeholder for future optional parameters
}

// SystemData - Metadata pertaining to creation and last modification of the resource.
type SystemData struct {
	// The timestamp of resource creation (UTC).
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	// The identity that created the resource.
	CreatedBy *string `json:"createdBy,omitempty"`

	// The type of identity that created the resource.
	CreatedByType *CreatedByType `json:"createdByType,omitempty"`

	// The timestamp of resource last modification (UTC)
	LastModifiedAt *time.Time `json:"lastModifiedAt,omitempty"`

	// The identity that last modified the resource.
	LastModifiedBy *string `json:"lastModifiedBy,omitempty"`

	// The type of identity that last modified the resource.
	LastModifiedByType *CreatedByType `json:"lastModifiedByType,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type SystemData.
func (s SystemData) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "createdAt", (*timeRFC3339)(s.CreatedAt))
	populate(objectMap, "createdBy", s.CreatedBy)
	populate(objectMap, "createdByType", s.CreatedByType)
	populate(objectMap, "lastModifiedAt", (*timeRFC3339)(s.LastModifiedAt))
	populate(objectMap, "lastModifiedBy", s.LastModifiedBy)
	populate(objectMap, "lastModifiedByType", s.LastModifiedByType)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type SystemData.
func (s *SystemData) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return err
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "createdAt":
				var aux timeRFC3339
				err = unpopulate(val, &aux)
				s.CreatedAt = (*time.Time)(&aux)
				delete(rawMsg, key)
		case "createdBy":
				err = unpopulate(val, &s.CreatedBy)
				delete(rawMsg, key)
		case "createdByType":
				err = unpopulate(val, &s.CreatedByType)
				delete(rawMsg, key)
		case "lastModifiedAt":
				var aux timeRFC3339
				err = unpopulate(val, &aux)
				s.LastModifiedAt = (*time.Time)(&aux)
				delete(rawMsg, key)
		case "lastModifiedBy":
				err = unpopulate(val, &s.LastModifiedBy)
				delete(rawMsg, key)
		case "lastModifiedByType":
				err = unpopulate(val, &s.LastModifiedByType)
				delete(rawMsg, key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// TrackedResource - The resource model definition for an Azure Resource Manager tracked top level resource which has 'tags' and a 'location'
type TrackedResource struct {
	Resource
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type TrackedResource.
func (t TrackedResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	t.marshalInternal(objectMap)
	return json.Marshal(objectMap)
}

func (t TrackedResource) marshalInternal(objectMap map[string]interface{}) {
	t.Resource.marshalInternal(objectMap)
	populate(objectMap, "location", t.Location)
	populate(objectMap, "tags", t.Tags)
}

func populate(m map[string]interface{}, k string, v interface{}) {
	if v == nil {
		return
	} else if azcore.IsNullValue(v) {
		m[k] = nil
	} else if !reflect.ValueOf(v).IsNil() {
		m[k] = v
	}
}

func unpopulate(data json.RawMessage, v interface{}) error {
	if data == nil {
		return nil
	}
	return json.Unmarshal(data, v)
}

