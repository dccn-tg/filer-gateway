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

// GetTasksTypeIDReader is a Reader for the GetTasksTypeID structure.
type GetTasksTypeIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetTasksTypeIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetTasksTypeIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGetTasksTypeIDBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewGetTasksTypeIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetTasksTypeIDInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /tasks/{type}/{id}] GetTasksTypeID", response, response.Code())
	}
}

// NewGetTasksTypeIDOK creates a GetTasksTypeIDOK with default headers values
func NewGetTasksTypeIDOK() *GetTasksTypeIDOK {
	return &GetTasksTypeIDOK{}
}

/*
GetTasksTypeIDOK describes a response with status code 200, with default header values.

success
*/
type GetTasksTypeIDOK struct {
	Payload *models.ResponseBodyTaskResource
}

// IsSuccess returns true when this get tasks type Id o k response has a 2xx status code
func (o *GetTasksTypeIDOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get tasks type Id o k response has a 3xx status code
func (o *GetTasksTypeIDOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get tasks type Id o k response has a 4xx status code
func (o *GetTasksTypeIDOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get tasks type Id o k response has a 5xx status code
func (o *GetTasksTypeIDOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get tasks type Id o k response a status code equal to that given
func (o *GetTasksTypeIDOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get tasks type Id o k response
func (o *GetTasksTypeIDOK) Code() int {
	return 200
}

func (o *GetTasksTypeIDOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /tasks/{type}/{id}][%d] getTasksTypeIdOK %s", 200, payload)
}

func (o *GetTasksTypeIDOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /tasks/{type}/{id}][%d] getTasksTypeIdOK %s", 200, payload)
}

func (o *GetTasksTypeIDOK) GetPayload() *models.ResponseBodyTaskResource {
	return o.Payload
}

func (o *GetTasksTypeIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBodyTaskResource)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetTasksTypeIDBadRequest creates a GetTasksTypeIDBadRequest with default headers values
func NewGetTasksTypeIDBadRequest() *GetTasksTypeIDBadRequest {
	return &GetTasksTypeIDBadRequest{}
}

/*
GetTasksTypeIDBadRequest describes a response with status code 400, with default header values.

bad request
*/
type GetTasksTypeIDBadRequest struct {
	Payload *models.ResponseBody400
}

// IsSuccess returns true when this get tasks type Id bad request response has a 2xx status code
func (o *GetTasksTypeIDBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get tasks type Id bad request response has a 3xx status code
func (o *GetTasksTypeIDBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get tasks type Id bad request response has a 4xx status code
func (o *GetTasksTypeIDBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this get tasks type Id bad request response has a 5xx status code
func (o *GetTasksTypeIDBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this get tasks type Id bad request response a status code equal to that given
func (o *GetTasksTypeIDBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the get tasks type Id bad request response
func (o *GetTasksTypeIDBadRequest) Code() int {
	return 400
}

func (o *GetTasksTypeIDBadRequest) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /tasks/{type}/{id}][%d] getTasksTypeIdBadRequest %s", 400, payload)
}

func (o *GetTasksTypeIDBadRequest) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /tasks/{type}/{id}][%d] getTasksTypeIdBadRequest %s", 400, payload)
}

func (o *GetTasksTypeIDBadRequest) GetPayload() *models.ResponseBody400 {
	return o.Payload
}

func (o *GetTasksTypeIDBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody400)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetTasksTypeIDNotFound creates a GetTasksTypeIDNotFound with default headers values
func NewGetTasksTypeIDNotFound() *GetTasksTypeIDNotFound {
	return &GetTasksTypeIDNotFound{}
}

/*
GetTasksTypeIDNotFound describes a response with status code 404, with default header values.

task not found
*/
type GetTasksTypeIDNotFound struct {
}

// IsSuccess returns true when this get tasks type Id not found response has a 2xx status code
func (o *GetTasksTypeIDNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get tasks type Id not found response has a 3xx status code
func (o *GetTasksTypeIDNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get tasks type Id not found response has a 4xx status code
func (o *GetTasksTypeIDNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this get tasks type Id not found response has a 5xx status code
func (o *GetTasksTypeIDNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this get tasks type Id not found response a status code equal to that given
func (o *GetTasksTypeIDNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the get tasks type Id not found response
func (o *GetTasksTypeIDNotFound) Code() int {
	return 404
}

func (o *GetTasksTypeIDNotFound) Error() string {
	return fmt.Sprintf("[GET /tasks/{type}/{id}][%d] getTasksTypeIdNotFound", 404)
}

func (o *GetTasksTypeIDNotFound) String() string {
	return fmt.Sprintf("[GET /tasks/{type}/{id}][%d] getTasksTypeIdNotFound", 404)
}

func (o *GetTasksTypeIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetTasksTypeIDInternalServerError creates a GetTasksTypeIDInternalServerError with default headers values
func NewGetTasksTypeIDInternalServerError() *GetTasksTypeIDInternalServerError {
	return &GetTasksTypeIDInternalServerError{}
}

/*
GetTasksTypeIDInternalServerError describes a response with status code 500, with default header values.

failure
*/
type GetTasksTypeIDInternalServerError struct {
	Payload *models.ResponseBody500
}

// IsSuccess returns true when this get tasks type Id internal server error response has a 2xx status code
func (o *GetTasksTypeIDInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get tasks type Id internal server error response has a 3xx status code
func (o *GetTasksTypeIDInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get tasks type Id internal server error response has a 4xx status code
func (o *GetTasksTypeIDInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get tasks type Id internal server error response has a 5xx status code
func (o *GetTasksTypeIDInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get tasks type Id internal server error response a status code equal to that given
func (o *GetTasksTypeIDInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the get tasks type Id internal server error response
func (o *GetTasksTypeIDInternalServerError) Code() int {
	return 500
}

func (o *GetTasksTypeIDInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /tasks/{type}/{id}][%d] getTasksTypeIdInternalServerError %s", 500, payload)
}

func (o *GetTasksTypeIDInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /tasks/{type}/{id}][%d] getTasksTypeIdInternalServerError %s", 500, payload)
}

func (o *GetTasksTypeIDInternalServerError) GetPayload() *models.ResponseBody500 {
	return o.Payload
}

func (o *GetTasksTypeIDInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody500)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
