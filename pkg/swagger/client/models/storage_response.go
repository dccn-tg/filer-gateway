// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// StorageResponse JSON object for storage resource data.
//
// swagger:model storageResponse
type StorageResponse struct {

	// assigned storage quota in GiB.
	// Required: true
	QuotaGb *int64 `json:"quotaGb"`

	// the targeting filer on which the storage resource is allocated.
	// Required: true
	// Enum: ["netapp","cephfs"]
	System *string `json:"system"`

	// used storage quota in MiB.
	// Required: true
	UsageMb *int64 `json:"usageMb"`
}

// Validate validates this storage response
func (m *StorageResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateQuotaGb(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSystem(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUsageMb(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *StorageResponse) validateQuotaGb(formats strfmt.Registry) error {

	if err := validate.Required("quotaGb", "body", m.QuotaGb); err != nil {
		return err
	}

	return nil
}

var storageResponseTypeSystemPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["netapp","cephfs"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		storageResponseTypeSystemPropEnum = append(storageResponseTypeSystemPropEnum, v)
	}
}

const (

	// StorageResponseSystemNetapp captures enum value "netapp"
	StorageResponseSystemNetapp string = "netapp"

	// StorageResponseSystemCephfs captures enum value "cephfs"
	StorageResponseSystemCephfs string = "cephfs"
)

// prop value enum
func (m *StorageResponse) validateSystemEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, storageResponseTypeSystemPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *StorageResponse) validateSystem(formats strfmt.Registry) error {

	if err := validate.Required("system", "body", m.System); err != nil {
		return err
	}

	// value enum
	if err := m.validateSystemEnum("system", "body", *m.System); err != nil {
		return err
	}

	return nil
}

func (m *StorageResponse) validateUsageMb(formats strfmt.Registry) error {

	if err := validate.Required("usageMb", "body", m.UsageMb); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this storage response based on context it is used
func (m *StorageResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *StorageResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *StorageResponse) UnmarshalBinary(b []byte) error {
	var res StorageResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
