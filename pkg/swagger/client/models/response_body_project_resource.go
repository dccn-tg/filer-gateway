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

// ResponseBodyProjectResource JSON object containing project resources.
//
// swagger:model responseBodyProjectResource
type ResponseBodyProjectResource struct {

	// members
	// Required: true
	Members Members `json:"members"`

	// project ID
	// Required: true
	ProjectID *ProjectID `json:"projectID"`

	// storage
	// Required: true
	Storage *StorageResponse `json:"storage"`
}

// Validate validates this response body project resource
func (m *ResponseBodyProjectResource) Validate(formats strfmt.Registry) error {
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

func (m *ResponseBodyProjectResource) validateMembers(formats strfmt.Registry) error {

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

func (m *ResponseBodyProjectResource) validateProjectID(formats strfmt.Registry) error {

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

func (m *ResponseBodyProjectResource) validateStorage(formats strfmt.Registry) error {

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

// ContextValidate validate this response body project resource based on the context it is used
func (m *ResponseBodyProjectResource) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
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

func (m *ResponseBodyProjectResource) contextValidateMembers(ctx context.Context, formats strfmt.Registry) error {

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

func (m *ResponseBodyProjectResource) contextValidateProjectID(ctx context.Context, formats strfmt.Registry) error {

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

func (m *ResponseBodyProjectResource) contextValidateStorage(ctx context.Context, formats strfmt.Registry) error {

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
func (m *ResponseBodyProjectResource) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ResponseBodyProjectResource) UnmarshalBinary(b []byte) error {
	var res ResponseBodyProjectResource
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
