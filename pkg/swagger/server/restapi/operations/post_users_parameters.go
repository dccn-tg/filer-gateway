// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/models"
)

// NewPostUsersParams creates a new PostUsersParams object
// no default values defined in spec.
func NewPostUsersParams() PostUsersParams {

	return PostUsersParams{}
}

// PostUsersParams contains all the bound params for the post users operation
// typically these are obtained from a http.Request
//
// swagger:parameters PostUsers
type PostUsersParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*data for user provisioning
	  Required: true
	  In: body
	*/
	UserProvisionData *models.RequestBodyUserProvision
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostUsersParams() beforehand.
func (o *PostUsersParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.RequestBodyUserProvision
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("userProvisionData", "body", ""))
			} else {
				res = append(res, errors.NewParseError("userProvisionData", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.UserProvisionData = &body
			}
		}
	} else {
		res = append(res, errors.Required("userProvisionData", "body", ""))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
