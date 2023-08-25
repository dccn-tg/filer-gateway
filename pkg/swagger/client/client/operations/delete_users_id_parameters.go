// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewDeleteUsersIDParams creates a new DeleteUsersIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewDeleteUsersIDParams() *DeleteUsersIDParams {
	return &DeleteUsersIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteUsersIDParamsWithTimeout creates a new DeleteUsersIDParams object
// with the ability to set a timeout on a request.
func NewDeleteUsersIDParamsWithTimeout(timeout time.Duration) *DeleteUsersIDParams {
	return &DeleteUsersIDParams{
		timeout: timeout,
	}
}

// NewDeleteUsersIDParamsWithContext creates a new DeleteUsersIDParams object
// with the ability to set a context for a request.
func NewDeleteUsersIDParamsWithContext(ctx context.Context) *DeleteUsersIDParams {
	return &DeleteUsersIDParams{
		Context: ctx,
	}
}

// NewDeleteUsersIDParamsWithHTTPClient creates a new DeleteUsersIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewDeleteUsersIDParamsWithHTTPClient(client *http.Client) *DeleteUsersIDParams {
	return &DeleteUsersIDParams{
		HTTPClient: client,
	}
}

/*
DeleteUsersIDParams contains all the parameters to send to the API endpoint

	for the delete users ID operation.

	Typically these are written to a http.Request.
*/
type DeleteUsersIDParams struct {

	/* ID.

	   user identifier
	*/
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the delete users ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteUsersIDParams) WithDefaults() *DeleteUsersIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the delete users ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteUsersIDParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the delete users ID params
func (o *DeleteUsersIDParams) WithTimeout(timeout time.Duration) *DeleteUsersIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete users ID params
func (o *DeleteUsersIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete users ID params
func (o *DeleteUsersIDParams) WithContext(ctx context.Context) *DeleteUsersIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete users ID params
func (o *DeleteUsersIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete users ID params
func (o *DeleteUsersIDParams) WithHTTPClient(client *http.Client) *DeleteUsersIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete users ID params
func (o *DeleteUsersIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the delete users ID params
func (o *DeleteUsersIDParams) WithID(id string) *DeleteUsersIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the delete users ID params
func (o *DeleteUsersIDParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteUsersIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
