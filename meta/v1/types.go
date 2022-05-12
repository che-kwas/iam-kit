package v1

import (
	"time"

	"encoding/json"

	"gorm.io/gorm"
)

// ObjectMeta is metadata that all persisted resources must have, which includes all objects
// ObjectMeta is also used by gorm.
type ObjectMeta struct {
	// ID is the unique in time and space value for this object. It is typically generated by
	// the storage on successful creation of a resource and is not allowed to change on PUT
	// operations.
	//
	// Populated by the system.
	// Read-only.
	ID uint64 `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT;column:id"`

	// InstanceID defines a string type resource identifier,
	// use prefixed to distinguish resource types, easy to remember, Url-friendly.
	InstanceID string `json:"instanceID,omitempty" gorm:"unique;column:instanceID;type:varchar(32);not null"`

	// Extend store the fields that need to be added, but do not want to add a new table column, will not be stored in db.
	Extend Extend `json:"extend,omitempty" gorm:"-" validate:"omitempty"`

	// ExtendShadow is the shadow of Extend. DO NOT modify directly.
	ExtendShadow string `json:"-" gorm:"column:extendShadow" validate:"omitempty"`

	// CreatedAt is a timestamp representing the server time when this object was
	// created. It is not guaranteed to be set in happens-before order across separate operations.
	// Clients may not set this value. It is represented in RFC3339 form and is in UTC.
	//
	// Populated by the system.
	// Read-only.
	// Null for lists.
	CreatedAt time.Time `json:"createdAt,omitempty" gorm:"column:createdAt"`

	// UpdatedAt is a timestamp representing the server time when this object was updated.
	// Clients may not set this value. It is represented in RFC3339 form and is in UTC.
	//
	// Populated by the system.
	// Read-only.
	// Null for lists.
	UpdatedAt time.Time `json:"updatedAt,omitempty" gorm:"column:updatedAt"`
}

// BeforeCreate run before create database record.
func (obj *ObjectMeta) BeforeCreate(tx *gorm.DB) error {
	obj.ExtendShadow = obj.Extend.String()

	return nil
}

// BeforeUpdate run before update database record.
func (obj *ObjectMeta) BeforeUpdate(tx *gorm.DB) error {
	obj.ExtendShadow = obj.Extend.String()

	return nil
}

// AfterFind run after find to unmarshal a extend shadown string into metav1.Extend struct.
func (obj *ObjectMeta) AfterFind(tx *gorm.DB) error {
	if err := json.Unmarshal([]byte(obj.ExtendShadow), &obj.Extend); err != nil {
		return err
	}

	return nil
}

// ListMeta describes metadata that synthetic resources must have, including lists and
// various status objects. A resource may have only one of {ObjectMeta, ListMeta}.
type ListMeta struct {
	TotalCount int64 `json:"totalCount,omitempty"`
}

// Extend defines a new type used to store extended fields.
type Extend map[string]interface{}

// String returns the string format of Extend.
func (ext Extend) String() string {
	data, _ := json.Marshal(ext)
	return string(data)
}

// TypeMeta describes an individual object in an API response or request
// with strings representing the type of the object and its API schema version.
// Structures that are versioned or persisted should inline TypeMeta.
type TypeMeta struct {
	// Kind is a string value representing the REST resource this object represents.
	// Servers may infer this from the endpoint the client submits requests to.
	// Cannot be updated.
	// In CamelCase.
	// required: false
	Kind string `json:"kind,omitempty"`

	// APIVersion defines the versioned schema of this representation of an object.
	// Servers should convert recognized schemas to the latest internal value, and
	// may reject unrecognized values.
	APIVersion string `json:"apiVersion,omitempty"`
}

// CreateOptions may be provided when creating an API object.
type CreateOptions struct {
	TypeMeta `json:",inline"`
}

// GetOptions is the standard query options to the standard REST get call.
type GetOptions struct {
	TypeMeta `json:",inline"`
}

// UpdateOptions may be provided when updating an API object.
// All fields in UpdateOptions should also be present in PatchOptions.
type UpdateOptions struct {
	TypeMeta `json:",inline"`
}

// DeleteOptions may be provided when deleting an API object.
type DeleteOptions struct {
	TypeMeta `json:",inline"`
}

// ListOptions is the query options to a standard REST list call.
type ListOptions struct {
	TypeMeta `json:",inline"`

	// Offset specify the number of records to skip before starting to return the records.
	Offset *int64 `json:"offset,omitempty"`

	// Limit specify the number of records to be retrieved.
	Limit *int64 `json:"limit,omitempty"`
}
