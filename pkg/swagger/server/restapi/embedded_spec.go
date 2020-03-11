// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "description": "filer gateway APIs",
    "title": "filer-gateway",
    "version": "0.1.0"
  },
  "basePath": "/v1",
  "paths": {
    "/projects": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "provision filer resource for a new project.",
        "parameters": [
          {
            "description": "data for project provisioning",
            "name": "projectProvisionData",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/requestBodyProjectProvision"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyProjectResource"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/projects/{id}": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "get filer resource for an existing project.",
        "parameters": [
          {
            "type": "string",
            "description": "project identifier",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyProjectResource"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "404": {
            "description": "project not found",
            "schema": {
              "type": "string",
              "enum": [
                "project not found"
              ]
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      },
      "patch": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "update filer resource for an existing project.",
        "parameters": [
          {
            "type": "string",
            "description": "project identifier",
            "name": "id",
            "in": "path",
            "required": true
          },
          {
            "description": "data for project update",
            "name": "projectUpdateData",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/requestBodyProjectResource"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyProjectResource"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "404": {
            "description": "project not found",
            "schema": {
              "type": "string",
              "enum": [
                "project not found"
              ]
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/projects/{id}/members": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "retrieves project members and their filer roles.",
        "parameters": [
          {
            "type": "string",
            "description": "project identifier",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyProjectMembers"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "404": {
            "description": "project not found",
            "schema": {
              "type": "string",
              "enum": [
                "project not found"
              ]
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/projects/{id}/storage": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "retrieves storage resource information of a project.",
        "parameters": [
          {
            "type": "string",
            "description": "project identifier",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyProjectStorage"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "404": {
            "description": "project not found",
            "schema": {
              "type": "string",
              "enum": [
                "project not found"
              ]
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/users": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "provision filer resource for a new user.",
        "parameters": [
          {
            "description": "data for user provisioning",
            "name": "userProvisionData",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/requestBodyUserProvision"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyUserResource"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/users/{id}": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "get filer resource for an existing user.",
        "parameters": [
          {
            "type": "string",
            "description": "user identifier",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyUserResource"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "404": {
            "description": "user not found",
            "schema": {
              "type": "string",
              "enum": [
                "user not found"
              ]
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      },
      "patch": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "update filer resource for an existing user.",
        "parameters": [
          {
            "type": "string",
            "description": "user identifier",
            "name": "id",
            "in": "path",
            "required": true
          },
          {
            "description": "data for user update",
            "name": "userUpdateData",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/requestBodyUserResource"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyUserResource"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "404": {
            "description": "user not found",
            "schema": {
              "type": "string",
              "enum": [
                "user not found"
              ]
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "member": {
      "description": "JSON object for a project member.",
      "required": [
        "userID",
        "role"
      ],
      "properties": {
        "role": {
          "description": "role of the member.",
          "type": "string",
          "enum": [
            "manager",
            "contributor",
            "viewer",
            "traverse",
            "none"
          ]
        },
        "userID": {
          "description": "userid of the member.",
          "type": "string"
        }
      }
    },
    "members": {
      "description": "a list of project members.",
      "type": "array",
      "items": {
        "$ref": "#/definitions/member"
      }
    },
    "projectID": {
      "description": "project identifier.",
      "type": "string"
    },
    "requestBodyProjectProvision": {
      "required": [
        "projectID",
        "storage",
        "members"
      ],
      "properties": {
        "members": {
          "$ref": "#/definitions/members"
        },
        "projectID": {
          "$ref": "#/definitions/projectID"
        },
        "storage": {
          "$ref": "#/definitions/storageRequest"
        }
      }
    },
    "requestBodyProjectResource": {
      "description": "JSON object describing resource to be set to the project.",
      "required": [
        "storage",
        "members"
      ],
      "properties": {
        "members": {
          "$ref": "#/definitions/members"
        },
        "storage": {
          "$ref": "#/definitions/storageRequest"
        }
      }
    },
    "requestBodyUserProvision": {
      "required": [
        "userID",
        "storage"
      ],
      "properties": {
        "storage": {
          "$ref": "#/definitions/storageRequest"
        },
        "userID": {
          "$ref": "#/definitions/userID"
        }
      }
    },
    "requestBodyUserResource": {
      "description": "JSON object describing resource to be set to the user.",
      "required": [
        "storage"
      ],
      "properties": {
        "storage": {
          "$ref": "#/definitions/storageRequest"
        }
      }
    },
    "responseBody400": {
      "description": "JSON object containing error message concerning bad client request.",
      "properties": {
        "errorMessage": {
          "description": "error message specifying the bad request.",
          "type": "string"
        }
      }
    },
    "responseBody500": {
      "description": "JSON object containing server side error.",
      "properties": {
        "errorMessage": {
          "description": "server-side error message.",
          "type": "string"
        },
        "exitCode": {
          "description": "server-side exit code.",
          "type": "integer"
        }
      }
    },
    "responseBodyProjectMembers": {
      "description": "JSON object containing current list of members implemented on the filer.",
      "required": [
        "projectID",
        "members"
      ],
      "properties": {
        "members": {
          "$ref": "#/definitions/members"
        },
        "projectID": {
          "$ref": "#/definitions/projectID"
        }
      }
    },
    "responseBodyProjectResource": {
      "description": "JSON object containing project resources.",
      "required": [
        "projectID",
        "storage",
        "members"
      ],
      "properties": {
        "members": {
          "$ref": "#/definitions/members"
        },
        "projectID": {
          "$ref": "#/definitions/projectID"
        },
        "storage": {
          "$ref": "#/definitions/storageResponse"
        }
      }
    },
    "responseBodyProjectStorage": {
      "description": "JSON object containing storage resource information of a project.",
      "required": [
        "projectID",
        "storage"
      ],
      "properties": {
        "projectID": {
          "$ref": "#/definitions/projectID"
        },
        "storage": {
          "$ref": "#/definitions/storageResponse"
        }
      }
    },
    "responseBodyUserResource": {
      "description": "JSON object containing user resources.",
      "required": [
        "userID",
        "storage"
      ],
      "properties": {
        "storage": {
          "$ref": "#/definitions/storageResponse"
        },
        "userID": {
          "$ref": "#/definitions/userID"
        }
      }
    },
    "storageRequest": {
      "description": "JSON object for storage resource data.",
      "required": [
        "system",
        "quotaGb"
      ],
      "properties": {
        "quotaGb": {
          "description": "storage quota in GiB.",
          "type": "integer"
        },
        "system": {
          "description": "the targeting filer on which the storage resource is allocated.",
          "type": "string",
          "enum": [
            "netapp",
            "freenas",
            "ceph"
          ]
        }
      }
    },
    "storageResponse": {
      "description": "JSON object for storage resource data.",
      "required": [
        "system",
        "quotaGb",
        "usageGb"
      ],
      "properties": {
        "quotaGb": {
          "description": "storage quota in GiB.",
          "type": "integer"
        },
        "system": {
          "description": "the targeting filer on which the storage resource is allocated.",
          "type": "string",
          "enum": [
            "netapp",
            "freenas",
            "ceph"
          ]
        },
        "usageGb": {
          "description": "used quota size in GiB (not used for the request data).",
          "type": "integer"
        }
      }
    },
    "userID": {
      "description": "user identifier.",
      "type": "string"
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "description": "filer gateway APIs",
    "title": "filer-gateway",
    "version": "0.1.0"
  },
  "basePath": "/v1",
  "paths": {
    "/projects": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "provision filer resource for a new project.",
        "parameters": [
          {
            "description": "data for project provisioning",
            "name": "projectProvisionData",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/requestBodyProjectProvision"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyProjectResource"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/projects/{id}": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "get filer resource for an existing project.",
        "parameters": [
          {
            "type": "string",
            "description": "project identifier",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyProjectResource"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "404": {
            "description": "project not found",
            "schema": {
              "type": "string",
              "enum": [
                "project not found"
              ]
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      },
      "patch": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "update filer resource for an existing project.",
        "parameters": [
          {
            "type": "string",
            "description": "project identifier",
            "name": "id",
            "in": "path",
            "required": true
          },
          {
            "description": "data for project update",
            "name": "projectUpdateData",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/requestBodyProjectResource"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyProjectResource"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "404": {
            "description": "project not found",
            "schema": {
              "type": "string",
              "enum": [
                "project not found"
              ]
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/projects/{id}/members": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "retrieves project members and their filer roles.",
        "parameters": [
          {
            "type": "string",
            "description": "project identifier",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyProjectMembers"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "404": {
            "description": "project not found",
            "schema": {
              "type": "string",
              "enum": [
                "project not found"
              ]
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/projects/{id}/storage": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "retrieves storage resource information of a project.",
        "parameters": [
          {
            "type": "string",
            "description": "project identifier",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyProjectStorage"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "404": {
            "description": "project not found",
            "schema": {
              "type": "string",
              "enum": [
                "project not found"
              ]
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/users": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "provision filer resource for a new user.",
        "parameters": [
          {
            "description": "data for user provisioning",
            "name": "userProvisionData",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/requestBodyUserProvision"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyUserResource"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/users/{id}": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "get filer resource for an existing user.",
        "parameters": [
          {
            "type": "string",
            "description": "user identifier",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyUserResource"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "404": {
            "description": "user not found",
            "schema": {
              "type": "string",
              "enum": [
                "user not found"
              ]
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      },
      "patch": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "update filer resource for an existing user.",
        "parameters": [
          {
            "type": "string",
            "description": "user identifier",
            "name": "id",
            "in": "path",
            "required": true
          },
          {
            "description": "data for user update",
            "name": "userUpdateData",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/requestBodyUserResource"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyUserResource"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "404": {
            "description": "user not found",
            "schema": {
              "type": "string",
              "enum": [
                "user not found"
              ]
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "member": {
      "description": "JSON object for a project member.",
      "required": [
        "userID",
        "role"
      ],
      "properties": {
        "role": {
          "description": "role of the member.",
          "type": "string",
          "enum": [
            "manager",
            "contributor",
            "viewer",
            "traverse",
            "none"
          ]
        },
        "userID": {
          "description": "userid of the member.",
          "type": "string"
        }
      }
    },
    "members": {
      "description": "a list of project members.",
      "type": "array",
      "items": {
        "$ref": "#/definitions/member"
      }
    },
    "projectID": {
      "description": "project identifier.",
      "type": "string"
    },
    "requestBodyProjectProvision": {
      "required": [
        "projectID",
        "storage",
        "members"
      ],
      "properties": {
        "members": {
          "$ref": "#/definitions/members"
        },
        "projectID": {
          "$ref": "#/definitions/projectID"
        },
        "storage": {
          "$ref": "#/definitions/storageRequest"
        }
      }
    },
    "requestBodyProjectResource": {
      "description": "JSON object describing resource to be set to the project.",
      "required": [
        "storage",
        "members"
      ],
      "properties": {
        "members": {
          "$ref": "#/definitions/members"
        },
        "storage": {
          "$ref": "#/definitions/storageRequest"
        }
      }
    },
    "requestBodyUserProvision": {
      "required": [
        "userID",
        "storage"
      ],
      "properties": {
        "storage": {
          "$ref": "#/definitions/storageRequest"
        },
        "userID": {
          "$ref": "#/definitions/userID"
        }
      }
    },
    "requestBodyUserResource": {
      "description": "JSON object describing resource to be set to the user.",
      "required": [
        "storage"
      ],
      "properties": {
        "storage": {
          "$ref": "#/definitions/storageRequest"
        }
      }
    },
    "responseBody400": {
      "description": "JSON object containing error message concerning bad client request.",
      "properties": {
        "errorMessage": {
          "description": "error message specifying the bad request.",
          "type": "string"
        }
      }
    },
    "responseBody500": {
      "description": "JSON object containing server side error.",
      "properties": {
        "errorMessage": {
          "description": "server-side error message.",
          "type": "string"
        },
        "exitCode": {
          "description": "server-side exit code.",
          "type": "integer"
        }
      }
    },
    "responseBodyProjectMembers": {
      "description": "JSON object containing current list of members implemented on the filer.",
      "required": [
        "projectID",
        "members"
      ],
      "properties": {
        "members": {
          "$ref": "#/definitions/members"
        },
        "projectID": {
          "$ref": "#/definitions/projectID"
        }
      }
    },
    "responseBodyProjectResource": {
      "description": "JSON object containing project resources.",
      "required": [
        "projectID",
        "storage",
        "members"
      ],
      "properties": {
        "members": {
          "$ref": "#/definitions/members"
        },
        "projectID": {
          "$ref": "#/definitions/projectID"
        },
        "storage": {
          "$ref": "#/definitions/storageResponse"
        }
      }
    },
    "responseBodyProjectStorage": {
      "description": "JSON object containing storage resource information of a project.",
      "required": [
        "projectID",
        "storage"
      ],
      "properties": {
        "projectID": {
          "$ref": "#/definitions/projectID"
        },
        "storage": {
          "$ref": "#/definitions/storageResponse"
        }
      }
    },
    "responseBodyUserResource": {
      "description": "JSON object containing user resources.",
      "required": [
        "userID",
        "storage"
      ],
      "properties": {
        "storage": {
          "$ref": "#/definitions/storageResponse"
        },
        "userID": {
          "$ref": "#/definitions/userID"
        }
      }
    },
    "storageRequest": {
      "description": "JSON object for storage resource data.",
      "required": [
        "system",
        "quotaGb"
      ],
      "properties": {
        "quotaGb": {
          "description": "storage quota in GiB.",
          "type": "integer"
        },
        "system": {
          "description": "the targeting filer on which the storage resource is allocated.",
          "type": "string",
          "enum": [
            "netapp",
            "freenas",
            "ceph"
          ]
        }
      }
    },
    "storageResponse": {
      "description": "JSON object for storage resource data.",
      "required": [
        "system",
        "quotaGb",
        "usageGb"
      ],
      "properties": {
        "quotaGb": {
          "description": "storage quota in GiB.",
          "type": "integer"
        },
        "system": {
          "description": "the targeting filer on which the storage resource is allocated.",
          "type": "string",
          "enum": [
            "netapp",
            "freenas",
            "ceph"
          ]
        },
        "usageGb": {
          "description": "used quota size in GiB (not used for the request data).",
          "type": "integer"
        }
      }
    },
    "userID": {
      "description": "user identifier.",
      "type": "string"
    }
  }
}`))
}
