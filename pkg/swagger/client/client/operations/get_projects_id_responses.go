// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/dccn-tg/filer-gateway/pkg/swagger/client/models"
)

// GetProjectsIDReader is a Reader for the GetProjectsID structure.
type GetProjectsIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetProjectsIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetProjectsIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGetProjectsIDBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewGetProjectsIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetProjectsIDInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /projects/{id}] GetProjectsID", response, response.Code())
	}
}

// NewGetProjectsIDOK creates a GetProjectsIDOK with default headers values
func NewGetProjectsIDOK() *GetProjectsIDOK {
	return &GetProjectsIDOK{}
}

/*
GetProjectsIDOK describes a response with status code 200, with default header values.

success
*/
type GetProjectsIDOK struct {
	Payload *models.ResponseBodyProjectResource
}

// IsSuccess returns true when this get projects Id o k response has a 2xx status code
func (o *GetProjectsIDOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get projects Id o k response has a 3xx status code
func (o *GetProjectsIDOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get projects Id o k response has a 4xx status code
func (o *GetProjectsIDOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get projects Id o k response has a 5xx status code
func (o *GetProjectsIDOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get projects Id o k response a status code equal to that given
func (o *GetProjectsIDOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get projects Id o k response
func (o *GetProjectsIDOK) Code() int {
	return 200
}

func (o *GetProjectsIDOK) Error() string {
	return fmt.Sprintf("[GET /projects/{id}][%d] getProjectsIdOK  %+v", 200, o.Payload)
}

func (o *GetProjectsIDOK) String() string {
	return fmt.Sprintf("[GET /projects/{id}][%d] getProjectsIdOK  %+v", 200, o.Payload)
}

func (o *GetProjectsIDOK) GetPayload() *models.ResponseBodyProjectResource {
	return o.Payload
}

func (o *GetProjectsIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBodyProjectResource)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetProjectsIDBadRequest creates a GetProjectsIDBadRequest with default headers values
func NewGetProjectsIDBadRequest() *GetProjectsIDBadRequest {
	return &GetProjectsIDBadRequest{}
}

/*
GetProjectsIDBadRequest describes a response with status code 400, with default header values.

bad request
*/
type GetProjectsIDBadRequest struct {
	Payload *models.ResponseBody400
}

// IsSuccess returns true when this get projects Id bad request response has a 2xx status code
func (o *GetProjectsIDBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get projects Id bad request response has a 3xx status code
func (o *GetProjectsIDBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get projects Id bad request response has a 4xx status code
func (o *GetProjectsIDBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this get projects Id bad request response has a 5xx status code
func (o *GetProjectsIDBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this get projects Id bad request response a status code equal to that given
func (o *GetProjectsIDBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the get projects Id bad request response
func (o *GetProjectsIDBadRequest) Code() int {
	return 400
}

func (o *GetProjectsIDBadRequest) Error() string {
	return fmt.Sprintf("[GET /projects/{id}][%d] getProjectsIdBadRequest  %+v", 400, o.Payload)
}

func (o *GetProjectsIDBadRequest) String() string {
	return fmt.Sprintf("[GET /projects/{id}][%d] getProjectsIdBadRequest  %+v", 400, o.Payload)
}

func (o *GetProjectsIDBadRequest) GetPayload() *models.ResponseBody400 {
	return o.Payload
}

func (o *GetProjectsIDBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody400)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetProjectsIDNotFound creates a GetProjectsIDNotFound with default headers values
func NewGetProjectsIDNotFound() *GetProjectsIDNotFound {
	return &GetProjectsIDNotFound{}
}

/*
GetProjectsIDNotFound describes a response with status code 404, with default header values.

project not found
*/
type GetProjectsIDNotFound struct {
	Payload string
}

// IsSuccess returns true when this get projects Id not found response has a 2xx status code
func (o *GetProjectsIDNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get projects Id not found response has a 3xx status code
func (o *GetProjectsIDNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get projects Id not found response has a 4xx status code
func (o *GetProjectsIDNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this get projects Id not found response has a 5xx status code
func (o *GetProjectsIDNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this get projects Id not found response a status code equal to that given
func (o *GetProjectsIDNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the get projects Id not found response
func (o *GetProjectsIDNotFound) Code() int {
	return 404
}

func (o *GetProjectsIDNotFound) Error() string {
	return fmt.Sprintf("[GET /projects/{id}][%d] getProjectsIdNotFound  %+v", 404, o.Payload)
}

func (o *GetProjectsIDNotFound) String() string {
	return fmt.Sprintf("[GET /projects/{id}][%d] getProjectsIdNotFound  %+v", 404, o.Payload)
}

func (o *GetProjectsIDNotFound) GetPayload() string {
	return o.Payload
}

func (o *GetProjectsIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetProjectsIDInternalServerError creates a GetProjectsIDInternalServerError with default headers values
func NewGetProjectsIDInternalServerError() *GetProjectsIDInternalServerError {
	return &GetProjectsIDInternalServerError{}
}

/*
GetProjectsIDInternalServerError describes a response with status code 500, with default header values.

failure
*/
type GetProjectsIDInternalServerError struct {
	Payload *models.ResponseBody500
}

// IsSuccess returns true when this get projects Id internal server error response has a 2xx status code
func (o *GetProjectsIDInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get projects Id internal server error response has a 3xx status code
func (o *GetProjectsIDInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get projects Id internal server error response has a 4xx status code
func (o *GetProjectsIDInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get projects Id internal server error response has a 5xx status code
func (o *GetProjectsIDInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get projects Id internal server error response a status code equal to that given
func (o *GetProjectsIDInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the get projects Id internal server error response
func (o *GetProjectsIDInternalServerError) Code() int {
	return 500
}

func (o *GetProjectsIDInternalServerError) Error() string {
	return fmt.Sprintf("[GET /projects/{id}][%d] getProjectsIdInternalServerError  %+v", 500, o.Payload)
}

func (o *GetProjectsIDInternalServerError) String() string {
	return fmt.Sprintf("[GET /projects/{id}][%d] getProjectsIdInternalServerError  %+v", 500, o.Payload)
}

func (o *GetProjectsIDInternalServerError) GetPayload() *models.ResponseBody500 {
	return o.Payload
}

func (o *GetProjectsIDInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody500)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
