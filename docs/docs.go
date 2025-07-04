// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/login": {
            "post": {
                "description": "Authenticates a user using email and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "description": "Login credentials",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.LoginPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.SignedResponse"
                        }
                    },
                    "400": {
                        "description": "error: invalid request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "error: invalid email or password",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "Registers a user using a one-time key. Key must not have been used before.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User registration data",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.RegisterPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message: user registered successfully",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "error: invalid request or key already used",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.ResponsePayload": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                }
            }
        },
        "handler.SignedResponse": {
            "type": "object",
            "properties": {
                "payload": {
                    "$ref": "#/definitions/handler.ResponsePayload"
                },
                "signature": {
                    "type": "string"
                }
            }
        },
        "models.LoginPayload": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "john@example.com"
                },
                "hwid": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "example": "securepassword123"
                }
            }
        },
        "models.RegisterPayload": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "john@example.com"
                },
                "first_name": {
                    "type": "string",
                    "example": "John"
                },
                "invitation_key": {
                    "type": "string",
                    "example": "abc123"
                },
                "last_name": {
                    "type": "string",
                    "example": "Doe"
                },
                "password": {
                    "type": "string",
                    "example": "securepassword123"
                },
                "phone": {
                    "type": "string",
                    "example": "+998901234567"
                },
                "username": {
                    "type": "string",
                    "example": "johndoe"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
