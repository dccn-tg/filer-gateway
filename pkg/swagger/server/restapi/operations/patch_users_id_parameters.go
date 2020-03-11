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

	strfmt "github.com/go-openapi/strfmt"

	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/models"
)

// NewPatchUsersIDParams creates a new PatchUsersIDParams object
// no default values defined in spec.
func NewPatchUsersIDParams() PatchUsersIDParams {

	return PatchUsersIDParams{}
}

// PatchUsersIDParams contains all the bound params for the patch users ID operation
// typically these are obtained from a http.Request
//
// swagger:parameters PatchUsersID
type PatchUsersIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*user identifier
	  Required: true
	  In: path
	*/
	ID string
	/*data for user update
	  Required: true
	  In: body
	*/
	UserUpdateData *models.RequestBodyUserResource
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPatchUsersIDParams() beforehand.
func (o *PatchUsersIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rID, rhkID, _ := route.Params.GetOK("id")
	if err := o.bindID(rID, rhkID, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.RequestBodyUserResource
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("userUpdateData", "body"))
			} else {
				res = append(res, errors.NewParseError("userUpdateData", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.UserUpdateData = &body
			}
		}
	} else {
		res = append(res, errors.Required("userUpdateData", "body"))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindID binds and validates parameter ID from path.
func (o *PatchUsersIDParams) bindID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ID = raw

	return nil
}
