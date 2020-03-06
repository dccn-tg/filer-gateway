// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/restapi/operations"
)

//go:generate swagger generate server --target ../../server --name FilerGateway --spec ../../swagger.yaml --exclude-main

func configureFlags(api *operations.FilerGatewayAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.FilerGatewayAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.GetProjectsIDHandler == nil {
		api.GetProjectsIDHandler = operations.GetProjectsIDHandlerFunc(func(params operations.GetProjectsIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetProjectsID has not yet been implemented")
		})
	}
	if api.GetProjectsIDMembersHandler == nil {
		api.GetProjectsIDMembersHandler = operations.GetProjectsIDMembersHandlerFunc(func(params operations.GetProjectsIDMembersParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetProjectsIDMembers has not yet been implemented")
		})
	}
	if api.GetProjectsIDStorageHandler == nil {
		api.GetProjectsIDStorageHandler = operations.GetProjectsIDStorageHandlerFunc(func(params operations.GetProjectsIDStorageParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetProjectsIDStorage has not yet been implemented")
		})
	}
	if api.PatchProjectsIDHandler == nil {
		api.PatchProjectsIDHandler = operations.PatchProjectsIDHandlerFunc(func(params operations.PatchProjectsIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PatchProjectsID has not yet been implemented")
		})
	}
	if api.PostProjectsHandler == nil {
		api.PostProjectsHandler = operations.PostProjectsHandlerFunc(func(params operations.PostProjectsParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostProjects has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return middleware.Redoc(middleware.RedocOpts{}, handler)
}
