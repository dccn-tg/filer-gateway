// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// RequestBodyProjectProvision request body project provision
//
// swagger:model requestBodyProjectProvision
type RequestBodyProjectProvision struct {

	// members
	// Required: true
	Members Members `json:"members"`

	// project ID
	// Required: true
	ProjectID *ProjectID `json:"projectID"`

	// apply ACL setting for members recursively on existing files/directories.
	Recursion bool `json:"recursion,omitempty"`

	// storage
	// Required: true
	Storage *StorageRequest `json:"storage"`
}

// Validate validates this request body project provision
func (m *RequestBodyProjectProvision) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMembers(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProjectID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStorage(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RequestBodyProjectProvision) validateMembers(formats strfmt.Registry) error {

	if err := validate.Required("members", "body", m.Members); err != nil {
		return err
	}

	if err := m.Members.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("members")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("members")
		}
		return err
	}

	return nil
}

func (m *RequestBodyProjectProvision) validateProjectID(formats strfmt.Registry) error {

	if err := validate.Required("projectID", "body", m.ProjectID); err != nil {
		return err
	}

	if err := validate.Required("projectID", "body", m.ProjectID); err != nil {
		return err
	}

	if m.ProjectID != nil {
		if err := m.ProjectID.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("projectID")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("projectID")
			}
			return err
		}
	}

	return nil
}

func (m *RequestBodyProjectProvision) validateStorage(formats strfmt.Registry) error {

	if err := validate.Required("storage", "body", m.Storage); err != nil {
		return err
	}

	if m.Storage != nil {
		if err := m.Storage.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("storage")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("storage")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this request body project provision based on the context it is used
func (m *RequestBodyProjectProvision) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateMembers(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateProjectID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateStorage(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RequestBodyProjectProvision) contextValidateMembers(ctx context.Context, formats strfmt.Registry) error {

	if err := m.Members.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("members")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("members")
		}
		return err
	}

	return nil
}

func (m *RequestBodyProjectProvision) contextValidateProjectID(ctx context.Context, formats strfmt.Registry) error {

	if m.ProjectID != nil {
		if err := m.ProjectID.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("projectID")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("projectID")
			}
			return err
		}
	}

	return nil
}

func (m *RequestBodyProjectProvision) contextValidateStorage(ctx context.Context, formats strfmt.Registry) error {

	if m.Storage != nil {
		if err := m.Storage.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("storage")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("storage")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *RequestBodyProjectProvision) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RequestBodyProjectProvision) UnmarshalBinary(b []byte) error {
	var res RequestBodyProjectProvision
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
