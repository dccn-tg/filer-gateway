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

// GetPingReader is a Reader for the GetPing structure.
type GetPingReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPingReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetPingOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewGetPingInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /ping] GetPing", response, response.Code())
	}
}

// NewGetPingOK creates a GetPingOK with default headers values
func NewGetPingOK() *GetPingOK {
	return &GetPingOK{}
}

/*
GetPingOK describes a response with status code 200, with default header values.

success
*/
type GetPingOK struct {
	Payload string
}

// IsSuccess returns true when this get ping o k response has a 2xx status code
func (o *GetPingOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get ping o k response has a 3xx status code
func (o *GetPingOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get ping o k response has a 4xx status code
func (o *GetPingOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get ping o k response has a 5xx status code
func (o *GetPingOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get ping o k response a status code equal to that given
func (o *GetPingOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get ping o k response
func (o *GetPingOK) Code() int {
	return 200
}

func (o *GetPingOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /ping][%d] getPingOK %s", 200, payload)
}

func (o *GetPingOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /ping][%d] getPingOK %s", 200, payload)
}

func (o *GetPingOK) GetPayload() string {
	return o.Payload
}

func (o *GetPingOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetPingInternalServerError creates a GetPingInternalServerError with default headers values
func NewGetPingInternalServerError() *GetPingInternalServerError {
	return &GetPingInternalServerError{}
}

/*
GetPingInternalServerError describes a response with status code 500, with default header values.

failure
*/
type GetPingInternalServerError struct {
	Payload *models.ResponseBody500
}

// IsSuccess returns true when this get ping internal server error response has a 2xx status code
func (o *GetPingInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get ping internal server error response has a 3xx status code
func (o *GetPingInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get ping internal server error response has a 4xx status code
func (o *GetPingInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get ping internal server error response has a 5xx status code
func (o *GetPingInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get ping internal server error response a status code equal to that given
func (o *GetPingInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the get ping internal server error response
func (o *GetPingInternalServerError) Code() int {
	return 500
}

func (o *GetPingInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /ping][%d] getPingInternalServerError %s", 500, payload)
}

func (o *GetPingInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /ping][%d] getPingInternalServerError %s", 500, payload)
}

func (o *GetPingInternalServerError) GetPayload() *models.ResponseBody500 {
	return o.Payload
}

func (o *GetPingInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody500)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
