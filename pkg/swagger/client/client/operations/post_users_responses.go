// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/Donders-Institute/filer-gateway/pkg/swagger/client/models"
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
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewPostUsersOK creates a PostUsersOK with default headers values
func NewPostUsersOK() *PostUsersOK {
	return &PostUsersOK{}
}

/*PostUsersOK handles this case with default header values.

success
*/
type PostUsersOK struct {
	Payload *models.ResponseBodyTaskResource
}

func (o *PostUsersOK) Error() string {
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

/*PostUsersBadRequest handles this case with default header values.

bad request
*/
type PostUsersBadRequest struct {
	Payload *models.ResponseBody400
}

func (o *PostUsersBadRequest) Error() string {
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

/*PostUsersInternalServerError handles this case with default header values.

failure
*/
type PostUsersInternalServerError struct {
	Payload *models.ResponseBody500
}

func (o *PostUsersInternalServerError) Error() string {
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