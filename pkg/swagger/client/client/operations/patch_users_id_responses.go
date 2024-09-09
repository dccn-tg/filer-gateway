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

// PatchUsersIDReader is a Reader for the PatchUsersID structure.
type PatchUsersIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PatchUsersIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPatchUsersIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 204:
		result := NewPatchUsersIDNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewPatchUsersIDBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewPatchUsersIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewPatchUsersIDInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[PATCH /users/{id}] PatchUsersID", response, response.Code())
	}
}

// NewPatchUsersIDOK creates a PatchUsersIDOK with default headers values
func NewPatchUsersIDOK() *PatchUsersIDOK {
	return &PatchUsersIDOK{}
}

/*
PatchUsersIDOK describes a response with status code 200, with default header values.

success
*/
type PatchUsersIDOK struct {
	Payload *models.ResponseBodyTaskResource
}

// IsSuccess returns true when this patch users Id o k response has a 2xx status code
func (o *PatchUsersIDOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this patch users Id o k response has a 3xx status code
func (o *PatchUsersIDOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this patch users Id o k response has a 4xx status code
func (o *PatchUsersIDOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this patch users Id o k response has a 5xx status code
func (o *PatchUsersIDOK) IsServerError() bool {
	return false
}

// IsCode returns true when this patch users Id o k response a status code equal to that given
func (o *PatchUsersIDOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the patch users Id o k response
func (o *PatchUsersIDOK) Code() int {
	return 200
}

func (o *PatchUsersIDOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PATCH /users/{id}][%d] patchUsersIdOK %s", 200, payload)
}

func (o *PatchUsersIDOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PATCH /users/{id}][%d] patchUsersIdOK %s", 200, payload)
}

func (o *PatchUsersIDOK) GetPayload() *models.ResponseBodyTaskResource {
	return o.Payload
}

func (o *PatchUsersIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBodyTaskResource)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPatchUsersIDNoContent creates a PatchUsersIDNoContent with default headers values
func NewPatchUsersIDNoContent() *PatchUsersIDNoContent {
	return &PatchUsersIDNoContent{}
}

/*
PatchUsersIDNoContent describes a response with status code 204, with default header values.

no content
*/
type PatchUsersIDNoContent struct {
}

// IsSuccess returns true when this patch users Id no content response has a 2xx status code
func (o *PatchUsersIDNoContent) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this patch users Id no content response has a 3xx status code
func (o *PatchUsersIDNoContent) IsRedirect() bool {
	return false
}

// IsClientError returns true when this patch users Id no content response has a 4xx status code
func (o *PatchUsersIDNoContent) IsClientError() bool {
	return false
}

// IsServerError returns true when this patch users Id no content response has a 5xx status code
func (o *PatchUsersIDNoContent) IsServerError() bool {
	return false
}

// IsCode returns true when this patch users Id no content response a status code equal to that given
func (o *PatchUsersIDNoContent) IsCode(code int) bool {
	return code == 204
}

// Code gets the status code for the patch users Id no content response
func (o *PatchUsersIDNoContent) Code() int {
	return 204
}

func (o *PatchUsersIDNoContent) Error() string {
	return fmt.Sprintf("[PATCH /users/{id}][%d] patchUsersIdNoContent", 204)
}

func (o *PatchUsersIDNoContent) String() string {
	return fmt.Sprintf("[PATCH /users/{id}][%d] patchUsersIdNoContent", 204)
}

func (o *PatchUsersIDNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewPatchUsersIDBadRequest creates a PatchUsersIDBadRequest with default headers values
func NewPatchUsersIDBadRequest() *PatchUsersIDBadRequest {
	return &PatchUsersIDBadRequest{}
}

/*
PatchUsersIDBadRequest describes a response with status code 400, with default header values.

bad request
*/
type PatchUsersIDBadRequest struct {
	Payload *models.ResponseBody400
}

// IsSuccess returns true when this patch users Id bad request response has a 2xx status code
func (o *PatchUsersIDBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this patch users Id bad request response has a 3xx status code
func (o *PatchUsersIDBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this patch users Id bad request response has a 4xx status code
func (o *PatchUsersIDBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this patch users Id bad request response has a 5xx status code
func (o *PatchUsersIDBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this patch users Id bad request response a status code equal to that given
func (o *PatchUsersIDBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the patch users Id bad request response
func (o *PatchUsersIDBadRequest) Code() int {
	return 400
}

func (o *PatchUsersIDBadRequest) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PATCH /users/{id}][%d] patchUsersIdBadRequest %s", 400, payload)
}

func (o *PatchUsersIDBadRequest) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PATCH /users/{id}][%d] patchUsersIdBadRequest %s", 400, payload)
}

func (o *PatchUsersIDBadRequest) GetPayload() *models.ResponseBody400 {
	return o.Payload
}

func (o *PatchUsersIDBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody400)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPatchUsersIDNotFound creates a PatchUsersIDNotFound with default headers values
func NewPatchUsersIDNotFound() *PatchUsersIDNotFound {
	return &PatchUsersIDNotFound{}
}

/*
PatchUsersIDNotFound describes a response with status code 404, with default header values.

user not found
*/
type PatchUsersIDNotFound struct {
	Payload string
}

// IsSuccess returns true when this patch users Id not found response has a 2xx status code
func (o *PatchUsersIDNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this patch users Id not found response has a 3xx status code
func (o *PatchUsersIDNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this patch users Id not found response has a 4xx status code
func (o *PatchUsersIDNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this patch users Id not found response has a 5xx status code
func (o *PatchUsersIDNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this patch users Id not found response a status code equal to that given
func (o *PatchUsersIDNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the patch users Id not found response
func (o *PatchUsersIDNotFound) Code() int {
	return 404
}

func (o *PatchUsersIDNotFound) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PATCH /users/{id}][%d] patchUsersIdNotFound %s", 404, payload)
}

func (o *PatchUsersIDNotFound) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PATCH /users/{id}][%d] patchUsersIdNotFound %s", 404, payload)
}

func (o *PatchUsersIDNotFound) GetPayload() string {
	return o.Payload
}

func (o *PatchUsersIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPatchUsersIDInternalServerError creates a PatchUsersIDInternalServerError with default headers values
func NewPatchUsersIDInternalServerError() *PatchUsersIDInternalServerError {
	return &PatchUsersIDInternalServerError{}
}

/*
PatchUsersIDInternalServerError describes a response with status code 500, with default header values.

failure
*/
type PatchUsersIDInternalServerError struct {
	Payload *models.ResponseBody500
}

// IsSuccess returns true when this patch users Id internal server error response has a 2xx status code
func (o *PatchUsersIDInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this patch users Id internal server error response has a 3xx status code
func (o *PatchUsersIDInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this patch users Id internal server error response has a 4xx status code
func (o *PatchUsersIDInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this patch users Id internal server error response has a 5xx status code
func (o *PatchUsersIDInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this patch users Id internal server error response a status code equal to that given
func (o *PatchUsersIDInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the patch users Id internal server error response
func (o *PatchUsersIDInternalServerError) Code() int {
	return 500
}

func (o *PatchUsersIDInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PATCH /users/{id}][%d] patchUsersIdInternalServerError %s", 500, payload)
}

func (o *PatchUsersIDInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PATCH /users/{id}][%d] patchUsersIdInternalServerError %s", 500, payload)
}

func (o *PatchUsersIDInternalServerError) GetPayload() *models.ResponseBody500 {
	return o.Payload
}

func (o *PatchUsersIDInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody500)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
