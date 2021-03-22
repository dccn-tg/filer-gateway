// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/security"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/models"
)

// NewFilerGatewayAPI creates a new FilerGateway instance
func NewFilerGatewayAPI(spec *loads.Document) *FilerGatewayAPI {
	return &FilerGatewayAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		customConsumers:     make(map[string]runtime.Consumer),
		customProducers:     make(map[string]runtime.Producer),
		PreServerShutdown:   func() {},
		ServerShutdown:      func() {},
		spec:                spec,
		useSwaggerUI:        false,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,

		JSONConsumer: runtime.JSONConsumer(),

		JSONProducer: runtime.JSONProducer(),

		GetPingHandler: GetPingHandlerFunc(func(params GetPingParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation GetPing has not yet been implemented")
		}),
		GetProjectsHandler: GetProjectsHandlerFunc(func(params GetProjectsParams) middleware.Responder {
			return middleware.NotImplemented("operation GetProjects has not yet been implemented")
		}),
		GetProjectsIDHandler: GetProjectsIDHandlerFunc(func(params GetProjectsIDParams) middleware.Responder {
			return middleware.NotImplemented("operation GetProjectsID has not yet been implemented")
		}),
		GetTasksTypeIDHandler: GetTasksTypeIDHandlerFunc(func(params GetTasksTypeIDParams) middleware.Responder {
			return middleware.NotImplemented("operation GetTasksTypeID has not yet been implemented")
		}),
		GetUsersHandler: GetUsersHandlerFunc(func(params GetUsersParams) middleware.Responder {
			return middleware.NotImplemented("operation GetUsers has not yet been implemented")
		}),
		GetUsersIDHandler: GetUsersIDHandlerFunc(func(params GetUsersIDParams) middleware.Responder {
			return middleware.NotImplemented("operation GetUsersID has not yet been implemented")
		}),
		PatchProjectsIDHandler: PatchProjectsIDHandlerFunc(func(params PatchProjectsIDParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation PatchProjectsID has not yet been implemented")
		}),
		PatchUsersIDHandler: PatchUsersIDHandlerFunc(func(params PatchUsersIDParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation PatchUsersID has not yet been implemented")
		}),
		PostProjectsHandler: PostProjectsHandlerFunc(func(params PostProjectsParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation PostProjects has not yet been implemented")
		}),
		PostUsersHandler: PostUsersHandlerFunc(func(params PostUsersParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation PostUsers has not yet been implemented")
		}),

		// Applies when the "X-API-Key" header is set
		APIKeyHeaderAuth: func(token string) (*models.Principle, error) {
			return nil, errors.NotImplemented("api key auth (apiKeyHeader) X-API-Key from header param [X-API-Key] has not yet been implemented")
		},
		// Applies when the Authorization header is set with the Basic scheme
		BasicAuthAuth: func(user string, pass string) (*models.Principle, error) {
			return nil, errors.NotImplemented("basic auth  (basicAuth) has not yet been implemented")
		},
		Oauth2Auth: func(token string, scopes []string) (*models.Principle, error) {
			return nil, errors.NotImplemented("oauth2 bearer auth (oauth2) has not yet been implemented")
		},
		// default authorizer is authorized meaning no requests are blocked
		APIAuthorizer: security.Authorized(),
	}
}

/*FilerGatewayAPI filer gateway APIs */
type FilerGatewayAPI struct {
	spec            *loads.Document
	context         *middleware.Context
	handlers        map[string]map[string]http.Handler
	formats         strfmt.Registry
	customConsumers map[string]runtime.Consumer
	customProducers map[string]runtime.Producer
	defaultConsumes string
	defaultProduces string
	Middleware      func(middleware.Builder) http.Handler
	useSwaggerUI    bool

	// BasicAuthenticator generates a runtime.Authenticator from the supplied basic auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BasicAuthenticator func(security.UserPassAuthentication) runtime.Authenticator

	// APIKeyAuthenticator generates a runtime.Authenticator from the supplied token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	APIKeyAuthenticator func(string, string, security.TokenAuthentication) runtime.Authenticator

	// BearerAuthenticator generates a runtime.Authenticator from the supplied bearer token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BearerAuthenticator func(string, security.ScopedTokenAuthentication) runtime.Authenticator

	// JSONConsumer registers a consumer for the following mime types:
	//   - application/json
	JSONConsumer runtime.Consumer

	// JSONProducer registers a producer for the following mime types:
	//   - application/json
	JSONProducer runtime.Producer

	// APIKeyHeaderAuth registers a function that takes a token and returns a principal
	// it performs authentication based on an api key X-API-Key provided in the header
	APIKeyHeaderAuth func(string) (*models.Principle, error)

	// BasicAuthAuth registers a function that takes username and password and returns a principal
	// it performs authentication with basic auth
	BasicAuthAuth func(string, string) (*models.Principle, error)

	// Oauth2Auth registers a function that takes an access token and a collection of required scopes and returns a principal
	// it performs authentication based on an oauth2 bearer token provided in the request
	Oauth2Auth func(string, []string) (*models.Principle, error)

	// APIAuthorizer provides access control (ACL/RBAC/ABAC) by providing access to the request and authenticated principal
	APIAuthorizer runtime.Authorizer

	// GetPingHandler sets the operation handler for the get ping operation
	GetPingHandler GetPingHandler
	// GetProjectsHandler sets the operation handler for the get projects operation
	GetProjectsHandler GetProjectsHandler
	// GetProjectsIDHandler sets the operation handler for the get projects ID operation
	GetProjectsIDHandler GetProjectsIDHandler
	// GetTasksTypeIDHandler sets the operation handler for the get tasks type ID operation
	GetTasksTypeIDHandler GetTasksTypeIDHandler
	// GetUsersHandler sets the operation handler for the get users operation
	GetUsersHandler GetUsersHandler
	// GetUsersIDHandler sets the operation handler for the get users ID operation
	GetUsersIDHandler GetUsersIDHandler
	// PatchProjectsIDHandler sets the operation handler for the patch projects ID operation
	PatchProjectsIDHandler PatchProjectsIDHandler
	// PatchUsersIDHandler sets the operation handler for the patch users ID operation
	PatchUsersIDHandler PatchUsersIDHandler
	// PostProjectsHandler sets the operation handler for the post projects operation
	PostProjectsHandler PostProjectsHandler
	// PostUsersHandler sets the operation handler for the post users operation
	PostUsersHandler PostUsersHandler

	// ServeError is called when an error is received, there is a default handler
	// but you can set your own with this
	ServeError func(http.ResponseWriter, *http.Request, error)

	// PreServerShutdown is called before the HTTP(S) server is shutdown
	// This allows for custom functions to get executed before the HTTP(S) server stops accepting traffic
	PreServerShutdown func()

	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	ServerShutdown func()

	// Custom command line argument groups with their descriptions
	CommandLineOptionsGroups []swag.CommandLineOptionsGroup

	// User defined logger function.
	Logger func(string, ...interface{})
}

// UseRedoc for documentation at /docs
func (o *FilerGatewayAPI) UseRedoc() {
	o.useSwaggerUI = false
}

// UseSwaggerUI for documentation at /docs
func (o *FilerGatewayAPI) UseSwaggerUI() {
	o.useSwaggerUI = true
}

// SetDefaultProduces sets the default produces media type
func (o *FilerGatewayAPI) SetDefaultProduces(mediaType string) {
	o.defaultProduces = mediaType
}

// SetDefaultConsumes returns the default consumes media type
func (o *FilerGatewayAPI) SetDefaultConsumes(mediaType string) {
	o.defaultConsumes = mediaType
}

// SetSpec sets a spec that will be served for the clients.
func (o *FilerGatewayAPI) SetSpec(spec *loads.Document) {
	o.spec = spec
}

// DefaultProduces returns the default produces media type
func (o *FilerGatewayAPI) DefaultProduces() string {
	return o.defaultProduces
}

// DefaultConsumes returns the default consumes media type
func (o *FilerGatewayAPI) DefaultConsumes() string {
	return o.defaultConsumes
}

// Formats returns the registered string formats
func (o *FilerGatewayAPI) Formats() strfmt.Registry {
	return o.formats
}

// RegisterFormat registers a custom format validator
func (o *FilerGatewayAPI) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	o.formats.Add(name, format, validator)
}

// Validate validates the registrations in the FilerGatewayAPI
func (o *FilerGatewayAPI) Validate() error {
	var unregistered []string

	if o.JSONConsumer == nil {
		unregistered = append(unregistered, "JSONConsumer")
	}

	if o.JSONProducer == nil {
		unregistered = append(unregistered, "JSONProducer")
	}

	if o.APIKeyHeaderAuth == nil {
		unregistered = append(unregistered, "XAPIKeyAuth")
	}
	if o.BasicAuthAuth == nil {
		unregistered = append(unregistered, "BasicAuthAuth")
	}
	if o.Oauth2Auth == nil {
		unregistered = append(unregistered, "Oauth2Auth")
	}

	if o.GetPingHandler == nil {
		unregistered = append(unregistered, "GetPingHandler")
	}
	if o.GetProjectsHandler == nil {
		unregistered = append(unregistered, "GetProjectsHandler")
	}
	if o.GetProjectsIDHandler == nil {
		unregistered = append(unregistered, "GetProjectsIDHandler")
	}
	if o.GetTasksTypeIDHandler == nil {
		unregistered = append(unregistered, "GetTasksTypeIDHandler")
	}
	if o.GetUsersHandler == nil {
		unregistered = append(unregistered, "GetUsersHandler")
	}
	if o.GetUsersIDHandler == nil {
		unregistered = append(unregistered, "GetUsersIDHandler")
	}
	if o.PatchProjectsIDHandler == nil {
		unregistered = append(unregistered, "PatchProjectsIDHandler")
	}
	if o.PatchUsersIDHandler == nil {
		unregistered = append(unregistered, "PatchUsersIDHandler")
	}
	if o.PostProjectsHandler == nil {
		unregistered = append(unregistered, "PostProjectsHandler")
	}
	if o.PostUsersHandler == nil {
		unregistered = append(unregistered, "PostUsersHandler")
	}

	if len(unregistered) > 0 {
		return fmt.Errorf("missing registration: %s", strings.Join(unregistered, ", "))
	}

	return nil
}

// ServeErrorFor gets a error handler for a given operation id
func (o *FilerGatewayAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return o.ServeError
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (o *FilerGatewayAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]runtime.Authenticator {
	result := make(map[string]runtime.Authenticator)
	for name := range schemes {
		switch name {
		case "apiKeyHeader":
			scheme := schemes[name]
			result[name] = o.APIKeyAuthenticator(scheme.Name, scheme.In, func(token string) (interface{}, error) {
				return o.APIKeyHeaderAuth(token)
			})

		case "basicAuth":
			result[name] = o.BasicAuthenticator(func(username, password string) (interface{}, error) {
				return o.BasicAuthAuth(username, password)
			})

		case "oauth2":
			result[name] = o.BearerAuthenticator(name, func(token string, scopes []string) (interface{}, error) {
				return o.Oauth2Auth(token, scopes)
			})

		}
	}
	return result
}

// Authorizer returns the registered authorizer
func (o *FilerGatewayAPI) Authorizer() runtime.Authorizer {
	return o.APIAuthorizer
}

// ConsumersFor gets the consumers for the specified media types.
// MIME type parameters are ignored here.
func (o *FilerGatewayAPI) ConsumersFor(mediaTypes []string) map[string]runtime.Consumer {
	result := make(map[string]runtime.Consumer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONConsumer
		}

		if c, ok := o.customConsumers[mt]; ok {
			result[mt] = c
		}
	}
	return result
}

// ProducersFor gets the producers for the specified media types.
// MIME type parameters are ignored here.
func (o *FilerGatewayAPI) ProducersFor(mediaTypes []string) map[string]runtime.Producer {
	result := make(map[string]runtime.Producer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONProducer
		}

		if p, ok := o.customProducers[mt]; ok {
			result[mt] = p
		}
	}
	return result
}

// HandlerFor gets a http.Handler for the provided operation method and path
func (o *FilerGatewayAPI) HandlerFor(method, path string) (http.Handler, bool) {
	if o.handlers == nil {
		return nil, false
	}
	um := strings.ToUpper(method)
	if _, ok := o.handlers[um]; !ok {
		return nil, false
	}
	if path == "/" {
		path = ""
	}
	h, ok := o.handlers[um][path]
	return h, ok
}

// Context returns the middleware context for the filer gateway API
func (o *FilerGatewayAPI) Context() *middleware.Context {
	if o.context == nil {
		o.context = middleware.NewRoutableContext(o.spec, o, nil)
	}

	return o.context
}

func (o *FilerGatewayAPI) initHandlerCache() {
	o.Context() // don't care about the result, just that the initialization happened
	if o.handlers == nil {
		o.handlers = make(map[string]map[string]http.Handler)
	}

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/ping"] = NewGetPing(o.context, o.GetPingHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/projects"] = NewGetProjects(o.context, o.GetProjectsHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/projects/{id}"] = NewGetProjectsID(o.context, o.GetProjectsIDHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/tasks/{type}/{id}"] = NewGetTasksTypeID(o.context, o.GetTasksTypeIDHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/users"] = NewGetUsers(o.context, o.GetUsersHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/users/{id}"] = NewGetUsersID(o.context, o.GetUsersIDHandler)
	if o.handlers["PATCH"] == nil {
		o.handlers["PATCH"] = make(map[string]http.Handler)
	}
	o.handlers["PATCH"]["/projects/{id}"] = NewPatchProjectsID(o.context, o.PatchProjectsIDHandler)
	if o.handlers["PATCH"] == nil {
		o.handlers["PATCH"] = make(map[string]http.Handler)
	}
	o.handlers["PATCH"]["/users/{id}"] = NewPatchUsersID(o.context, o.PatchUsersIDHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/projects"] = NewPostProjects(o.context, o.PostProjectsHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/users"] = NewPostUsers(o.context, o.PostUsersHandler)
}

// Serve creates a http handler to serve the API over HTTP
// can be used directly in http.ListenAndServe(":8000", api.Serve(nil))
func (o *FilerGatewayAPI) Serve(builder middleware.Builder) http.Handler {
	o.Init()

	if o.Middleware != nil {
		return o.Middleware(builder)
	}
	if o.useSwaggerUI {
		return o.context.APIHandlerSwaggerUI(builder)
	}
	return o.context.APIHandler(builder)
}

// Init allows you to just initialize the handler cache, you can then recompose the middleware as you see fit
func (o *FilerGatewayAPI) Init() {
	if len(o.handlers) == 0 {
		o.initHandlerCache()
	}
}

// RegisterConsumer allows you to add (or override) a consumer for a media type.
func (o *FilerGatewayAPI) RegisterConsumer(mediaType string, consumer runtime.Consumer) {
	o.customConsumers[mediaType] = consumer
}

// RegisterProducer allows you to add (or override) a producer for a media type.
func (o *FilerGatewayAPI) RegisterProducer(mediaType string, producer runtime.Producer) {
	o.customProducers[mediaType] = producer
}

// AddMiddlewareFor adds a http middleware to existing handler
func (o *FilerGatewayAPI) AddMiddlewareFor(method, path string, builder middleware.Builder) {
	um := strings.ToUpper(method)
	if path == "/" {
		path = ""
	}
	o.Init()
	if h, ok := o.handlers[um][path]; ok {
		o.handlers[method][path] = builder(h)
	}
}
