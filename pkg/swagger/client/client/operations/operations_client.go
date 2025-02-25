// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// New creates a new operations API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

// New creates a new operations API client with basic auth credentials.
// It takes the following parameters:
// - host: http host (github.com).
// - basePath: any base path for the API client ("/v1", "/v3").
// - scheme: http scheme ("http", "https").
// - user: user for basic authentication header.
// - password: password for basic authentication header.
func NewClientWithBasicAuth(host, basePath, scheme, user, password string) ClientService {
	transport := httptransport.New(host, basePath, []string{scheme})
	transport.DefaultAuthentication = httptransport.BasicAuth(user, password)
	return &Client{transport: transport, formats: strfmt.Default}
}

// New creates a new operations API client with a bearer token for authentication.
// It takes the following parameters:
// - host: http host (github.com).
// - basePath: any base path for the API client ("/v1", "/v3").
// - scheme: http scheme ("http", "https").
// - bearerToken: bearer token for Bearer authentication header.
func NewClientWithBearerToken(host, basePath, scheme, bearerToken string) ClientService {
	transport := httptransport.New(host, basePath, []string{scheme})
	transport.DefaultAuthentication = httptransport.BearerToken(bearerToken)
	return &Client{transport: transport, formats: strfmt.Default}
}

/*
Client for operations API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption may be used to customize the behavior of Client methods.
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	DeleteUsersID(params *DeleteUsersIDParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*DeleteUsersIDOK, error)

	GetMetrics(params *GetMetricsParams, opts ...ClientOption) (*GetMetricsOK, error)

	GetPing(params *GetPingParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetPingOK, error)

	GetProjects(params *GetProjectsParams, opts ...ClientOption) (*GetProjectsOK, error)

	GetProjectsID(params *GetProjectsIDParams, opts ...ClientOption) (*GetProjectsIDOK, error)

	GetTasksTypeID(params *GetTasksTypeIDParams, opts ...ClientOption) (*GetTasksTypeIDOK, error)

	GetUsers(params *GetUsersParams, opts ...ClientOption) (*GetUsersOK, error)

	GetUsersID(params *GetUsersIDParams, opts ...ClientOption) (*GetUsersIDOK, error)

	PatchProjectsID(params *PatchProjectsIDParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PatchProjectsIDOK, *PatchProjectsIDNoContent, error)

	PatchUsersID(params *PatchUsersIDParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PatchUsersIDOK, *PatchUsersIDNoContent, error)

	PostProjects(params *PostProjectsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PostProjectsOK, error)

	PostUsers(params *PostUsersParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PostUsersOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
DeleteUsersID deletes home directory of an existing user
*/
func (a *Client) DeleteUsersID(params *DeleteUsersIDParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*DeleteUsersIDOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteUsersIDParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "DeleteUsersID",
		Method:             "DELETE",
		PathPattern:        "/users/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &DeleteUsersIDReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DeleteUsersIDOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for DeleteUsersID: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetMetrics prometheus metrics
*/
func (a *Client) GetMetrics(params *GetMetricsParams, opts ...ClientOption) (*GetMetricsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetMetricsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetMetrics",
		Method:             "GET",
		PathPattern:        "/metrics",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetMetricsReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetMetricsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetMetrics: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetPing endpoints for API server health check
*/
func (a *Client) GetPing(params *GetPingParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetPingOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetPingParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetPing",
		Method:             "GET",
		PathPattern:        "/ping",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetPingReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetPingOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetPing: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetProjects gets filer resources of all projects
*/
func (a *Client) GetProjects(params *GetProjectsParams, opts ...ClientOption) (*GetProjectsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetProjectsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetProjects",
		Method:             "GET",
		PathPattern:        "/projects",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetProjectsReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetProjectsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetProjects: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetProjectsID gets filer resource for an existing project
*/
func (a *Client) GetProjectsID(params *GetProjectsIDParams, opts ...ClientOption) (*GetProjectsIDOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetProjectsIDParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetProjectsID",
		Method:             "GET",
		PathPattern:        "/projects/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetProjectsIDReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetProjectsIDOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetProjectsID: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetTasksTypeID queries background task status
*/
func (a *Client) GetTasksTypeID(params *GetTasksTypeIDParams, opts ...ClientOption) (*GetTasksTypeIDOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetTasksTypeIDParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetTasksTypeID",
		Method:             "GET",
		PathPattern:        "/tasks/{type}/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetTasksTypeIDReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetTasksTypeIDOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetTasksTypeID: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetUsers gets filer resources of all users
*/
func (a *Client) GetUsers(params *GetUsersParams, opts ...ClientOption) (*GetUsersOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetUsersParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetUsers",
		Method:             "GET",
		PathPattern:        "/users",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetUsersReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetUsersOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetUsers: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetUsersID gets filer resource for an existing user
*/
func (a *Client) GetUsersID(params *GetUsersIDParams, opts ...ClientOption) (*GetUsersIDOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetUsersIDParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetUsersID",
		Method:             "GET",
		PathPattern:        "/users/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetUsersIDReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetUsersIDOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetUsersID: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
PatchProjectsID updates filer resource for an existing project
*/
func (a *Client) PatchProjectsID(params *PatchProjectsIDParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PatchProjectsIDOK, *PatchProjectsIDNoContent, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPatchProjectsIDParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PatchProjectsID",
		Method:             "PATCH",
		PathPattern:        "/projects/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PatchProjectsIDReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, nil, err
	}
	switch value := result.(type) {
	case *PatchProjectsIDOK:
		return value, nil, nil
	case *PatchProjectsIDNoContent:
		return nil, value, nil
	}
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for operations: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
PatchUsersID updates filer resource for an existing user
*/
func (a *Client) PatchUsersID(params *PatchUsersIDParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PatchUsersIDOK, *PatchUsersIDNoContent, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPatchUsersIDParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PatchUsersID",
		Method:             "PATCH",
		PathPattern:        "/users/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PatchUsersIDReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, nil, err
	}
	switch value := result.(type) {
	case *PatchUsersIDOK:
		return value, nil, nil
	case *PatchUsersIDNoContent:
		return nil, value, nil
	}
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for operations: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
PostProjects provisions filer resource for a new project
*/
func (a *Client) PostProjects(params *PostProjectsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PostProjectsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostProjectsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PostProjects",
		Method:             "POST",
		PathPattern:        "/projects",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PostProjectsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PostProjectsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PostProjects: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
PostUsers provisions filer resource for a new user
*/
func (a *Client) PostUsers(params *PostUsersParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PostUsersOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostUsersParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PostUsers",
		Method:             "POST",
		PathPattern:        "/users",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PostUsersReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PostUsersOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PostUsers: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
