// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/dccn-tg/filer-gateway/pkg/swagger/client/models"
)

// PatchProjectsIDReader is a Reader for the PatchProjectsID structure.
type PatchProjectsIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PatchProjectsIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPatchProjectsIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 204:
		result := NewPatchProjectsIDNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewPatchProjectsIDBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewPatchProjectsIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewPatchProjectsIDInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[PATCH /projects/{id}] PatchProjectsID", response, response.Code())
	}
}

// NewPatchProjectsIDOK creates a PatchProjectsIDOK with default headers values
func NewPatchProjectsIDOK() *PatchProjectsIDOK {
	return &PatchProjectsIDOK{}
}

/*
PatchProjectsIDOK describes a response with status code 200, with default header values.

success
*/
type PatchProjectsIDOK struct {
	Payload *models.ResponseBodyTaskResource
}

// IsSuccess returns true when this patch projects Id o k response has a 2xx status code
func (o *PatchProjectsIDOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this patch projects Id o k response has a 3xx status code
func (o *PatchProjectsIDOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this patch projects Id o k response has a 4xx status code
func (o *PatchProjectsIDOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this patch projects Id o k response has a 5xx status code
func (o *PatchProjectsIDOK) IsServerError() bool {
	return false
}

// IsCode returns true when this patch projects Id o k response a status code equal to that given
func (o *PatchProjectsIDOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the patch projects Id o k response
func (o *PatchProjectsIDOK) Code() int {
	return 200
}

func (o *PatchProjectsIDOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PATCH /projects/{id}][%d] patchProjectsIdOK %s", 200, payload)
}

func (o *PatchProjectsIDOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PATCH /projects/{id}][%d] patchProjectsIdOK %s", 200, payload)
}

func (o *PatchProjectsIDOK) GetPayload() *models.ResponseBodyTaskResource {
	return o.Payload
}

func (o *PatchProjectsIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBodyTaskResource)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPatchProjectsIDNoContent creates a PatchProjectsIDNoContent with default headers values
func NewPatchProjectsIDNoContent() *PatchProjectsIDNoContent {
	return &PatchProjectsIDNoContent{}
}

/*
PatchProjectsIDNoContent describes a response with status code 204, with default header values.

no content
*/
type PatchProjectsIDNoContent struct {
}

// IsSuccess returns true when this patch projects Id no content response has a 2xx status code
func (o *PatchProjectsIDNoContent) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this patch projects Id no content response has a 3xx status code
func (o *PatchProjectsIDNoContent) IsRedirect() bool {
	return false
}

// IsClientError returns true when this patch projects Id no content response has a 4xx status code
func (o *PatchProjectsIDNoContent) IsClientError() bool {
	return false
}

// IsServerError returns true when this patch projects Id no content response has a 5xx status code
func (o *PatchProjectsIDNoContent) IsServerError() bool {
	return false
}

// IsCode returns true when this patch projects Id no content response a status code equal to that given
func (o *PatchProjectsIDNoContent) IsCode(code int) bool {
	return code == 204
}

// Code gets the status code for the patch projects Id no content response
func (o *PatchProjectsIDNoContent) Code() int {
	return 204
}

func (o *PatchProjectsIDNoContent) Error() string {
	return fmt.Sprintf("[PATCH /projects/{id}][%d] patchProjectsIdNoContent", 204)
}

func (o *PatchProjectsIDNoContent) String() string {
	return fmt.Sprintf("[PATCH /projects/{id}][%d] patchProjectsIdNoContent", 204)
}

func (o *PatchProjectsIDNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewPatchProjectsIDBadRequest creates a PatchProjectsIDBadRequest with default headers values
func NewPatchProjectsIDBadRequest() *PatchProjectsIDBadRequest {
	return &PatchProjectsIDBadRequest{}
}

/*
PatchProjectsIDBadRequest describes a response with status code 400, with default header values.

bad request
*/
type PatchProjectsIDBadRequest struct {
	Payload *models.ResponseBody400
}

// IsSuccess returns true when this patch projects Id bad request response has a 2xx status code
func (o *PatchProjectsIDBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this patch projects Id bad request response has a 3xx status code
func (o *PatchProjectsIDBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this patch projects Id bad request response has a 4xx status code
func (o *PatchProjectsIDBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this patch projects Id bad request response has a 5xx status code
func (o *PatchProjectsIDBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this patch projects Id bad request response a status code equal to that given
func (o *PatchProjectsIDBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the patch projects Id bad request response
func (o *PatchProjectsIDBadRequest) Code() int {
	return 400
}

func (o *PatchProjectsIDBadRequest) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PATCH /projects/{id}][%d] patchProjectsIdBadRequest %s", 400, payload)
}

func (o *PatchProjectsIDBadRequest) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PATCH /projects/{id}][%d] patchProjectsIdBadRequest %s", 400, payload)
}

func (o *PatchProjectsIDBadRequest) GetPayload() *models.ResponseBody400 {
	return o.Payload
}

func (o *PatchProjectsIDBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody400)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPatchProjectsIDNotFound creates a PatchProjectsIDNotFound with default headers values
func NewPatchProjectsIDNotFound() *PatchProjectsIDNotFound {
	return &PatchProjectsIDNotFound{}
}

/*
PatchProjectsIDNotFound describes a response with status code 404, with default header values.

project not found
*/
type PatchProjectsIDNotFound struct {
	Payload string
}

// IsSuccess returns true when this patch projects Id not found response has a 2xx status code
func (o *PatchProjectsIDNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this patch projects Id not found response has a 3xx status code
func (o *PatchProjectsIDNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this patch projects Id not found response has a 4xx status code
func (o *PatchProjectsIDNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this patch projects Id not found response has a 5xx status code
func (o *PatchProjectsIDNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this patch projects Id not found response a status code equal to that given
func (o *PatchProjectsIDNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the patch projects Id not found response
func (o *PatchProjectsIDNotFound) Code() int {
	return 404
}

func (o *PatchProjectsIDNotFound) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PATCH /projects/{id}][%d] patchProjectsIdNotFound %s", 404, payload)
}

func (o *PatchProjectsIDNotFound) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PATCH /projects/{id}][%d] patchProjectsIdNotFound %s", 404, payload)
}

func (o *PatchProjectsIDNotFound) GetPayload() string {
	return o.Payload
}

func (o *PatchProjectsIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPatchProjectsIDInternalServerError creates a PatchProjectsIDInternalServerError with default headers values
func NewPatchProjectsIDInternalServerError() *PatchProjectsIDInternalServerError {
	return &PatchProjectsIDInternalServerError{}
}

/*
PatchProjectsIDInternalServerError describes a response with status code 500, with default header values.

failure
*/
type PatchProjectsIDInternalServerError struct {
	Payload *models.ResponseBody500
}

// IsSuccess returns true when this patch projects Id internal server error response has a 2xx status code
func (o *PatchProjectsIDInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this patch projects Id internal server error response has a 3xx status code
func (o *PatchProjectsIDInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this patch projects Id internal server error response has a 4xx status code
func (o *PatchProjectsIDInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this patch projects Id internal server error response has a 5xx status code
func (o *PatchProjectsIDInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this patch projects Id internal server error response a status code equal to that given
func (o *PatchProjectsIDInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the patch projects Id internal server error response
func (o *PatchProjectsIDInternalServerError) Code() int {
	return 500
}

func (o *PatchProjectsIDInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PATCH /projects/{id}][%d] patchProjectsIdInternalServerError %s", 500, payload)
}

func (o *PatchProjectsIDInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PATCH /projects/{id}][%d] patchProjectsIdInternalServerError %s", 500, payload)
}

func (o *PatchProjectsIDInternalServerError) GetPayload() *models.ResponseBody500 {
	return o.Payload
}

func (o *PatchProjectsIDInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody500)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
