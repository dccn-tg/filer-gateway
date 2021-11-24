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

// ResponseBodyTaskResource JSON object containing scheduled task resource.
//
// swagger:model responseBodyTaskResource
type ResponseBodyTaskResource struct {

	// task ID
	// Required: true
	TaskID *TaskID `json:"taskID"`

	// task status
	// Required: true
	TaskStatus *TaskStatus `json:"taskStatus"`
}

// Validate validates this response body task resource
func (m *ResponseBodyTaskResource) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTaskID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTaskStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ResponseBodyTaskResource) validateTaskID(formats strfmt.Registry) error {

	if err := validate.Required("taskID", "body", m.TaskID); err != nil {
		return err
	}

	if err := validate.Required("taskID", "body", m.TaskID); err != nil {
		return err
	}

	if m.TaskID != nil {
		if err := m.TaskID.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("taskID")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("taskID")
			}
			return err
		}
	}

	return nil
}

func (m *ResponseBodyTaskResource) validateTaskStatus(formats strfmt.Registry) error {

	if err := validate.Required("taskStatus", "body", m.TaskStatus); err != nil {
		return err
	}

	if m.TaskStatus != nil {
		if err := m.TaskStatus.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("taskStatus")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("taskStatus")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this response body task resource based on the context it is used
func (m *ResponseBodyTaskResource) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateTaskID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTaskStatus(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ResponseBodyTaskResource) contextValidateTaskID(ctx context.Context, formats strfmt.Registry) error {

	if m.TaskID != nil {
		if err := m.TaskID.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("taskID")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("taskID")
			}
			return err
		}
	}

	return nil
}

func (m *ResponseBodyTaskResource) contextValidateTaskStatus(ctx context.Context, formats strfmt.Registry) error {

	if m.TaskStatus != nil {
		if err := m.TaskStatus.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("taskStatus")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("taskStatus")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ResponseBodyTaskResource) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ResponseBodyTaskResource) UnmarshalBinary(b []byte) error {
	var res ResponseBodyTaskResource
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
