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
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewPatchProjectsIDOK creates a PatchProjectsIDOK with default headers values
func NewPatchProjectsIDOK() *PatchProjectsIDOK {
	return &PatchProjectsIDOK{}
}

/*PatchProjectsIDOK handles this case with default header values.

success
*/
type PatchProjectsIDOK struct {
	Payload *models.ResponseBodyTaskResource
}

func (o *PatchProjectsIDOK) Error() string {
	return fmt.Sprintf("[PATCH /projects/{id}][%d] patchProjectsIdOK  %+v", 200, o.Payload)
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

// NewPatchProjectsIDBadRequest creates a PatchProjectsIDBadRequest with default headers values
func NewPatchProjectsIDBadRequest() *PatchProjectsIDBadRequest {
	return &PatchProjectsIDBadRequest{}
}

/*PatchProjectsIDBadRequest handles this case with default header values.

bad request
*/
type PatchProjectsIDBadRequest struct {
	Payload *models.ResponseBody400
}

func (o *PatchProjectsIDBadRequest) Error() string {
	return fmt.Sprintf("[PATCH /projects/{id}][%d] patchProjectsIdBadRequest  %+v", 400, o.Payload)
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

/*PatchProjectsIDNotFound handles this case with default header values.

project not found
*/
type PatchProjectsIDNotFound struct {
	Payload string
}

func (o *PatchProjectsIDNotFound) Error() string {
	return fmt.Sprintf("[PATCH /projects/{id}][%d] patchProjectsIdNotFound  %+v", 404, o.Payload)
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

/*PatchProjectsIDInternalServerError handles this case with default header values.

failure
*/
type PatchProjectsIDInternalServerError struct {
	Payload *models.ResponseBody500
}

func (o *PatchProjectsIDInternalServerError) Error() string {
	return fmt.Sprintf("[PATCH /projects/{id}][%d] patchProjectsIdInternalServerError  %+v", 500, o.Payload)
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
