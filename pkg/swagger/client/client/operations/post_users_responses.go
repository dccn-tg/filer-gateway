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

// PostUsersReader is a Reader for the PostUsers structure.
type PostUsersReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostUsersReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPostUsersOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewPostUsersBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewPostUsersInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[POST /users] PostUsers", response, response.Code())
	}
}

// NewPostUsersOK creates a PostUsersOK with default headers values
func NewPostUsersOK() *PostUsersOK {
	return &PostUsersOK{}
}

/*
PostUsersOK describes a response with status code 200, with default header values.

success
*/
type PostUsersOK struct {
	Payload *models.ResponseBodyTaskResource
}

// IsSuccess returns true when this post users o k response has a 2xx status code
func (o *PostUsersOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this post users o k response has a 3xx status code
func (o *PostUsersOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this post users o k response has a 4xx status code
func (o *PostUsersOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this post users o k response has a 5xx status code
func (o *PostUsersOK) IsServerError() bool {
	return false
}

// IsCode returns true when this post users o k response a status code equal to that given
func (o *PostUsersOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the post users o k response
func (o *PostUsersOK) Code() int {
	return 200
}

func (o *PostUsersOK) Error() string {
	return fmt.Sprintf("[POST /users][%d] postUsersOK  %+v", 200, o.Payload)
}

func (o *PostUsersOK) String() string {
	return fmt.Sprintf("[POST /users][%d] postUsersOK  %+v", 200, o.Payload)
}

func (o *PostUsersOK) GetPayload() *models.ResponseBodyTaskResource {
	return o.Payload
}

func (o *PostUsersOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBodyTaskResource)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostUsersBadRequest creates a PostUsersBadRequest with default headers values
func NewPostUsersBadRequest() *PostUsersBadRequest {
	return &PostUsersBadRequest{}
}

/*
PostUsersBadRequest describes a response with status code 400, with default header values.

bad request
*/
type PostUsersBadRequest struct {
	Payload *models.ResponseBody400
}

// IsSuccess returns true when this post users bad request response has a 2xx status code
func (o *PostUsersBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this post users bad request response has a 3xx status code
func (o *PostUsersBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this post users bad request response has a 4xx status code
func (o *PostUsersBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this post users bad request response has a 5xx status code
func (o *PostUsersBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this post users bad request response a status code equal to that given
func (o *PostUsersBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the post users bad request response
func (o *PostUsersBadRequest) Code() int {
	return 400
}

func (o *PostUsersBadRequest) Error() string {
	return fmt.Sprintf("[POST /users][%d] postUsersBadRequest  %+v", 400, o.Payload)
}

func (o *PostUsersBadRequest) String() string {
	return fmt.Sprintf("[POST /users][%d] postUsersBadRequest  %+v", 400, o.Payload)
}

func (o *PostUsersBadRequest) GetPayload() *models.ResponseBody400 {
	return o.Payload
}

func (o *PostUsersBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody400)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostUsersInternalServerError creates a PostUsersInternalServerError with default headers values
func NewPostUsersInternalServerError() *PostUsersInternalServerError {
	return &PostUsersInternalServerError{}
}

/*
PostUsersInternalServerError describes a response with status code 500, with default header values.

failure
*/
type PostUsersInternalServerError struct {
	Payload *models.ResponseBody500
}

// IsSuccess returns true when this post users internal server error response has a 2xx status code
func (o *PostUsersInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this post users internal server error response has a 3xx status code
func (o *PostUsersInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this post users internal server error response has a 4xx status code
func (o *PostUsersInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this post users internal server error response has a 5xx status code
func (o *PostUsersInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this post users internal server error response a status code equal to that given
func (o *PostUsersInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the post users internal server error response
func (o *PostUsersInternalServerError) Code() int {
	return 500
}

func (o *PostUsersInternalServerError) Error() string {
	return fmt.Sprintf("[POST /users][%d] postUsersInternalServerError  %+v", 500, o.Payload)
}

func (o *PostUsersInternalServerError) String() string {
	return fmt.Sprintf("[POST /users][%d] postUsersInternalServerError  %+v", 500, o.Payload)
}

func (o *PostUsersInternalServerError) GetPayload() *models.ResponseBody500 {
	return o.Payload
}

func (o *PostUsersInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody500)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
