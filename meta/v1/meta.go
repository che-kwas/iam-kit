package v1

import (
	"time"
)

// Object lets you work with object metadata from any of the versioned or
// internal API objects. Attempting to set or retrieve a field on an object that does
// not support that field will be a no-op and return a default value.
type Object interface {
	GetID() uint64
	SetID(id uint64)
	GetCreatedAt() time.Time
	SetCreatedAt(createdAt time.Time)
	GetUpdatedAt() time.Time
	SetUpdatedAt(updatedAt time.Time)
}

// List lets you work with list metadata from any of the versioned or
// internal API objects. Attempting to set or retrieve a field on an object that does
// not support that field will be a no-op and return a default value.
type List interface {
	GetTotalCount() int64
	SetTotalCount(count int64)
}

// Type exposes the type and APIVersion of versioned or internal API objects.
type Type interface {
	GetAPIVersion() string
	SetAPIVersion(version string)
	GetKind() string
	SetKind(kind string)
}

var _ Object = &ObjectMeta{}

func (meta *ObjectMeta) GetID() uint64                    { return meta.ID }
func (meta *ObjectMeta) SetID(id uint64)                  { meta.ID = id }
func (meta *ObjectMeta) GetCreatedAt() time.Time          { return meta.CreatedAt }
func (meta *ObjectMeta) SetCreatedAt(createdAt time.Time) { meta.CreatedAt = createdAt }
func (meta *ObjectMeta) GetUpdatedAt() time.Time          { return meta.UpdatedAt }
func (meta *ObjectMeta) SetUpdatedAt(updatedAt time.Time) { meta.UpdatedAt = updatedAt }

var _ List = &ListMeta{}

func (meta *ListMeta) GetTotalCount() int64      { return meta.TotalCount }
func (meta *ListMeta) SetTotalCount(count int64) { meta.TotalCount = count }

var _ Type = &TypeMeta{}

func (meta *TypeMeta) GetAPIVersion() string        { return meta.APIVersion }
func (meta *TypeMeta) SetAPIVersion(version string) { meta.APIVersion = version }
func (meta *TypeMeta) GetKind() string              { return meta.Kind }
func (meta *TypeMeta) SetKind(kind string)          { meta.Kind = kind }
