// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ResponseBody403 JSON object containing error message concerning forbidden request.
//
// swagger:model responseBody403
type ResponseBody403 struct {

	// error message specifying the forbidden request.
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// Validate validates this response body403
func (m *ResponseBody403) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this response body403 based on context it is used
func (m *ResponseBody403) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ResponseBody403) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ResponseBody403) UnmarshalBinary(b []byte) error {
	var res ResponseBody403
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
