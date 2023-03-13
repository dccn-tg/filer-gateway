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

	"github.com/Donders-Institute/filer-gateway/pkg/swagger/client/models"
)

// NewPostUsersParams creates a new PostUsersParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPostUsersParams() *PostUsersParams {
	return &PostUsersParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPostUsersParamsWithTimeout creates a new PostUsersParams object
// with the ability to set a timeout on a request.
func NewPostUsersParamsWithTimeout(timeout time.Duration) *PostUsersParams {
	return &PostUsersParams{
		timeout: timeout,
	}
}

// NewPostUsersParamsWithContext creates a new PostUsersParams object
// with the ability to set a context for a request.
func NewPostUsersParamsWithContext(ctx context.Context) *PostUsersParams {
	return &PostUsersParams{
		Context: ctx,
	}
}

// NewPostUsersParamsWithHTTPClient creates a new PostUsersParams object
// with the ability to set a custom HTTPClient for a request.
func NewPostUsersParamsWithHTTPClient(client *http.Client) *PostUsersParams {
	return &PostUsersParams{
		HTTPClient: client,
	}
}

/*
PostUsersParams contains all the parameters to send to the API endpoint

	for the post users operation.

	Typically these are written to a http.Request.
*/
type PostUsersParams struct {

	/* UserProvisionData.

	   data for user provisioning
	*/
	UserProvisionData *models.RequestBodyUserProvision

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the post users params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostUsersParams) WithDefaults() *PostUsersParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the post users params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostUsersParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the post users params
func (o *PostUsersParams) WithTimeout(timeout time.Duration) *PostUsersParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post users params
func (o *PostUsersParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post users params
func (o *PostUsersParams) WithContext(ctx context.Context) *PostUsersParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post users params
func (o *PostUsersParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post users params
func (o *PostUsersParams) WithHTTPClient(client *http.Client) *PostUsersParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post users params
func (o *PostUsersParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithUserProvisionData adds the userProvisionData to the post users params
func (o *PostUsersParams) WithUserProvisionData(userProvisionData *models.RequestBodyUserProvision) *PostUsersParams {
	o.SetUserProvisionData(userProvisionData)
	return o
}

// SetUserProvisionData adds the userProvisionData to the post users params
func (o *PostUsersParams) SetUserProvisionData(userProvisionData *models.RequestBodyUserProvision) {
	o.UserProvisionData = userProvisionData
}

// WriteToRequest writes these params to a swagger request
func (o *PostUsersParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.UserProvisionData != nil {
		if err := r.SetBodyParam(o.UserProvisionData); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
