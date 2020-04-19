// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/models"
)

// PatchUsersIDHandlerFunc turns a function with the right signature into a patch users ID handler
type PatchUsersIDHandlerFunc func(PatchUsersIDParams, *models.Principle) middleware.Responder

// Handle executing the request and returning a response
func (fn PatchUsersIDHandlerFunc) Handle(params PatchUsersIDParams, principal *models.Principle) middleware.Responder {
	return fn(params, principal)
}

// PatchUsersIDHandler interface for that can handle valid patch users ID params
type PatchUsersIDHandler interface {
	Handle(PatchUsersIDParams, *models.Principle) middleware.Responder
}

// NewPatchUsersID creates a new http.Handler for the patch users ID operation
func NewPatchUsersID(ctx *middleware.Context, handler PatchUsersIDHandler) *PatchUsersID {
	return &PatchUsersID{Context: ctx, Handler: handler}
}

/*PatchUsersID swagger:route PATCH /users/{id} patchUsersId

update filer resource for an existing user.

*/
type PatchUsersID struct {
	Context *middleware.Context
	Handler PatchUsersIDHandler
}

func (o *PatchUsersID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPatchUsersIDParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *models.Principle
	if uprinc != nil {
		principal = uprinc.(*models.Principle) // this is really a models.Principle, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
