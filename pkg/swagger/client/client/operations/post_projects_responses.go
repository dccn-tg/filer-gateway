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

// PostProjectsReader is a Reader for the PostProjects structure.
type PostProjectsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostProjectsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPostProjectsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewPostProjectsBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewPostProjectsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewPostProjectsOK creates a PostProjectsOK with default headers values
func NewPostProjectsOK() *PostProjectsOK {
	return &PostProjectsOK{}
}

/*
PostProjectsOK describes a response with status code 200, with default header values.

success
*/
type PostProjectsOK struct {
	Payload *models.ResponseBodyTaskResource
}

// IsSuccess returns true when this post projects o k response has a 2xx status code
func (o *PostProjectsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this post projects o k response has a 3xx status code
func (o *PostProjectsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this post projects o k response has a 4xx status code
func (o *PostProjectsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this post projects o k response has a 5xx status code
func (o *PostProjectsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this post projects o k response a status code equal to that given
func (o *PostProjectsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the post projects o k response
func (o *PostProjectsOK) Code() int {
	return 200
}

func (o *PostProjectsOK) Error() string {
	return fmt.Sprintf("[POST /projects][%d] postProjectsOK  %+v", 200, o.Payload)
}

func (o *PostProjectsOK) String() string {
	return fmt.Sprintf("[POST /projects][%d] postProjectsOK  %+v", 200, o.Payload)
}

func (o *PostProjectsOK) GetPayload() *models.ResponseBodyTaskResource {
	return o.Payload
}

func (o *PostProjectsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBodyTaskResource)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostProjectsBadRequest creates a PostProjectsBadRequest with default headers values
func NewPostProjectsBadRequest() *PostProjectsBadRequest {
	return &PostProjectsBadRequest{}
}

/*
PostProjectsBadRequest describes a response with status code 400, with default header values.

bad request
*/
type PostProjectsBadRequest struct {
	Payload *models.ResponseBody400
}

// IsSuccess returns true when this post projects bad request response has a 2xx status code
func (o *PostProjectsBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this post projects bad request response has a 3xx status code
func (o *PostProjectsBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this post projects bad request response has a 4xx status code
func (o *PostProjectsBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this post projects bad request response has a 5xx status code
func (o *PostProjectsBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this post projects bad request response a status code equal to that given
func (o *PostProjectsBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the post projects bad request response
func (o *PostProjectsBadRequest) Code() int {
	return 400
}

func (o *PostProjectsBadRequest) Error() string {
	return fmt.Sprintf("[POST /projects][%d] postProjectsBadRequest  %+v", 400, o.Payload)
}

func (o *PostProjectsBadRequest) String() string {
	return fmt.Sprintf("[POST /projects][%d] postProjectsBadRequest  %+v", 400, o.Payload)
}

func (o *PostProjectsBadRequest) GetPayload() *models.ResponseBody400 {
	return o.Payload
}

func (o *PostProjectsBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody400)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostProjectsInternalServerError creates a PostProjectsInternalServerError with default headers values
func NewPostProjectsInternalServerError() *PostProjectsInternalServerError {
	return &PostProjectsInternalServerError{}
}

/*
PostProjectsInternalServerError describes a response with status code 500, with default header values.

failure
*/
type PostProjectsInternalServerError struct {
	Payload *models.ResponseBody500
}

// IsSuccess returns true when this post projects internal server error response has a 2xx status code
func (o *PostProjectsInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this post projects internal server error response has a 3xx status code
func (o *PostProjectsInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this post projects internal server error response has a 4xx status code
func (o *PostProjectsInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this post projects internal server error response has a 5xx status code
func (o *PostProjectsInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this post projects internal server error response a status code equal to that given
func (o *PostProjectsInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the post projects internal server error response
func (o *PostProjectsInternalServerError) Code() int {
	return 500
}

func (o *PostProjectsInternalServerError) Error() string {
	return fmt.Sprintf("[POST /projects][%d] postProjectsInternalServerError  %+v", 500, o.Payload)
}

func (o *PostProjectsInternalServerError) String() string {
	return fmt.Sprintf("[POST /projects][%d] postProjectsInternalServerError  %+v", 500, o.Payload)
}

func (o *PostProjectsInternalServerError) GetPayload() *models.ResponseBody500 {
	return o.Payload
}

func (o *PostProjectsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody500)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
