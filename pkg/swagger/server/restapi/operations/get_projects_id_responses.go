// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/dccn-tg/filer-gateway/pkg/swagger/server/models"
)

// GetProjectsIDOKCode is the HTTP code returned for type GetProjectsIDOK
const GetProjectsIDOKCode int = 200

/*
GetProjectsIDOK success

swagger:response getProjectsIdOK
*/
type GetProjectsIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.ResponseBodyProjectResource `json:"body,omitempty"`
}

// NewGetProjectsIDOK creates GetProjectsIDOK with default headers values
func NewGetProjectsIDOK() *GetProjectsIDOK {

	return &GetProjectsIDOK{}
}

// WithPayload adds the payload to the get projects Id o k response
func (o *GetProjectsIDOK) WithPayload(payload *models.ResponseBodyProjectResource) *GetProjectsIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get projects Id o k response
func (o *GetProjectsIDOK) SetPayload(payload *models.ResponseBodyProjectResource) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProjectsIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetProjectsIDBadRequestCode is the HTTP code returned for type GetProjectsIDBadRequest
const GetProjectsIDBadRequestCode int = 400

/*
GetProjectsIDBadRequest bad request

swagger:response getProjectsIdBadRequest
*/
type GetProjectsIDBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.ResponseBody400 `json:"body,omitempty"`
}

// NewGetProjectsIDBadRequest creates GetProjectsIDBadRequest with default headers values
func NewGetProjectsIDBadRequest() *GetProjectsIDBadRequest {

	return &GetProjectsIDBadRequest{}
}

// WithPayload adds the payload to the get projects Id bad request response
func (o *GetProjectsIDBadRequest) WithPayload(payload *models.ResponseBody400) *GetProjectsIDBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get projects Id bad request response
func (o *GetProjectsIDBadRequest) SetPayload(payload *models.ResponseBody400) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProjectsIDBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetProjectsIDNotFoundCode is the HTTP code returned for type GetProjectsIDNotFound
const GetProjectsIDNotFoundCode int = 404

/*
GetProjectsIDNotFound project not found

swagger:response getProjectsIdNotFound
*/
type GetProjectsIDNotFound struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetProjectsIDNotFound creates GetProjectsIDNotFound with default headers values
func NewGetProjectsIDNotFound() *GetProjectsIDNotFound {

	return &GetProjectsIDNotFound{}
}

// WithPayload adds the payload to the get projects Id not found response
func (o *GetProjectsIDNotFound) WithPayload(payload string) *GetProjectsIDNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get projects Id not found response
func (o *GetProjectsIDNotFound) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProjectsIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetProjectsIDInternalServerErrorCode is the HTTP code returned for type GetProjectsIDInternalServerError
const GetProjectsIDInternalServerErrorCode int = 500

/*
GetProjectsIDInternalServerError failure

swagger:response getProjectsIdInternalServerError
*/
type GetProjectsIDInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ResponseBody500 `json:"body,omitempty"`
}

// NewGetProjectsIDInternalServerError creates GetProjectsIDInternalServerError with default headers values
func NewGetProjectsIDInternalServerError() *GetProjectsIDInternalServerError {

	return &GetProjectsIDInternalServerError{}
}

// WithPayload adds the payload to the get projects Id internal server error response
func (o *GetProjectsIDInternalServerError) WithPayload(payload *models.ResponseBody500) *GetProjectsIDInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get projects Id internal server error response
func (o *GetProjectsIDInternalServerError) SetPayload(payload *models.ResponseBody500) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProjectsIDInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
