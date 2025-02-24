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
        "/analyze": {
            "post": {
                "description": "Analyzes a webpage for HTML version, headings, links, and login form",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "analyzer"
                ],
                "summary": "Analyze a webpage",
                "parameters": [
                    {
                        "description": "URL to analyze",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.AnalysisRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.PageAnalysis"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/domain.APIError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/domain.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/domain.APIError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.APIError": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "statusCode": {
                    "type": "integer"
                }
            }
        },
        "domain.AnalysisRequest": {
            "type": "object",
            "required": [
                "url"
            ],
            "properties": {
                "url": {
                    "type": "string"
                }
            }
        },
        "domain.HeadingCount": {
            "type": "object",
            "properties": {
                "h1": {
                    "type": "integer"
                },
                "h2": {
                    "type": "integer"
                },
                "h3": {
                    "type": "integer"
                },
                "h4": {
                    "type": "integer"
                },
                "h5": {
                    "type": "integer"
                },
                "h6": {
                    "type": "integer"
                }
            }
        },
        "domain.LinkAnalysis": {
            "type": "object",
            "properties": {
                "external": {
                    "type": "integer"
                },
                "inaccessible": {
                    "type": "integer"
                },
                "internal": {
                    "type": "integer"
                }
            }
        },
        "domain.PageAnalysis": {
            "type": "object",
            "properties": {
                "hasLoginForm": {
                    "type": "boolean"
                },
                "headings": {
                    "$ref": "#/definitions/domain.HeadingCount"
                },
                "htmlVersion": {
                    "type": "string"
                },
                "links": {
                    "$ref": "#/definitions/domain.LinkAnalysis"
                },
                "pageTitle": {
                    "type": "string"
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
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
