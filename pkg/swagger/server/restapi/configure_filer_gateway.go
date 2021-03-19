// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/rs/cors"

	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/models"
	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/restapi/operations"
)

//go:generate swagger generate server --target ../../server --name FilerGateway --spec ../../swagger.yaml --principal models.Principle --exclude-main

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

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "X-API-Key" header is set
	if api.APIKeyHeaderAuth == nil {
		api.APIKeyHeaderAuth = func(token string) (*models.Principle, error) {
			return nil, errors.NotImplemented("api key auth (apiKeyHeader) X-API-Key from header param [X-API-Key] has not yet been implemented")
		}
	}
	// Applies when the Authorization header is set with the Basic scheme
	if api.BasicAuthAuth == nil {
		api.BasicAuthAuth = func(user string, pass string) (*models.Principle, error) {
			return nil, errors.NotImplemented("basic auth  (basicAuth) has not yet been implemented")
		}
	}
	if api.Oauth2Auth == nil {
		api.Oauth2Auth = func(token string, scopes []string) (*models.Principle, error) {
			return nil, errors.NotImplemented("oauth2 bearer auth (oauth2) has not yet been implemented")
		}
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()

	if api.GetPingHandler == nil {
		api.GetPingHandler = operations.GetPingHandlerFunc(func(params operations.GetPingParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetPing has not yet been implemented")
		})
	}
	if api.GetProjectsHandler == nil {
		api.GetProjectsHandler = operations.GetProjectsHandlerFunc(func(params operations.GetProjectsParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetProjects has not yet been implemented")
		})
	}
	if api.GetProjectsIDHandler == nil {
		api.GetProjectsIDHandler = operations.GetProjectsIDHandlerFunc(func(params operations.GetProjectsIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetProjectsID has not yet been implemented")
		})
	}
	if api.GetTasksTypeIDHandler == nil {
		api.GetTasksTypeIDHandler = operations.GetTasksTypeIDHandlerFunc(func(params operations.GetTasksTypeIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetTasksTypeID has not yet been implemented")
		})
	}
	if api.GetUsersIDHandler == nil {
		api.GetUsersIDHandler = operations.GetUsersIDHandlerFunc(func(params operations.GetUsersIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetUsersID has not yet been implemented")
		})
	}
	if api.PatchProjectsIDHandler == nil {
		api.PatchProjectsIDHandler = operations.PatchProjectsIDHandlerFunc(func(params operations.PatchProjectsIDParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation operations.PatchProjectsID has not yet been implemented")
		})
	}
	if api.PatchUsersIDHandler == nil {
		api.PatchUsersIDHandler = operations.PatchUsersIDHandlerFunc(func(params operations.PatchUsersIDParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation operations.PatchUsersID has not yet been implemented")
		})
	}
	if api.PostProjectsHandler == nil {
		api.PostProjectsHandler = operations.PostProjectsHandlerFunc(func(params operations.PostProjectsParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostProjects has not yet been implemented")
		})
	}
	if api.PostUsersHandler == nil {
		api.PostUsersHandler = operations.PostUsersHandlerFunc(func(params operations.PostUsersParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostUsers has not yet been implemented")
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
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handleCORS := cors.Default().Handler
	return handleCORS(handler)
}
