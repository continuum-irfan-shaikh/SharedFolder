package schemas

import (
	"fmt"
)

// Validation schemas that uses for A-CMS
var (
	ExtendedEndpointID = `
		{
			"type": "object",
			"properties": {
				"partnerID": {
					"type": "string",
					"minLength": 1
				},
				"clientID": {
					"type": "string",
					"minLength": 1
				},
				"siteID": {
					"type": "string",
					"minLength": 1
				},
				"endpointID": {
					"type": "string",
					"pattern": "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
				}
			},
			"required": ["partnerID", "clientID", "siteID", "endpointID"]
		}`

	Audit = `
		{
			"type": "object",
			"properties": {
				"partnerID": {
					"type": "string",
					"minLength": 1
				},
				"clientID": {
					"type": "string",
					"minLength": 1
				},
				"siteID": {
					"type": "string",
					"minLength": 1
				},
				"endpointID": {
					"type": "string",
					"pattern": "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
				},
				"forceUpdate": {
					"type": "boolean"
				}
			},
			"required": ["partnerID", "clientID", "siteID", "endpointID"]
		}`

	InstallationStatusMessage = `
		{
			"type": "object",
			"properties": {
				"partnerID": {
					"type": "string",
					"minLength": 1
				},
				"clientID": {
					"type": "string",
					"minLength": 1
				},
				"siteID": {
					"type": "string",
					"minLength": 1
				},
				"endpointID": {
					"type": "string",
					"pattern": "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
				},
				"status": {
					"type": "string",
					"minLength": 1
				}
			},
			"required": ["partnerID", "clientID", "siteID", "endpointID", "status"]
		}`

	DeviceStatusMessage = `
		{
			"type": "object",
			"properties": {
				"partnerID": {
					"type": "string",
					"minLength": 1
				},
				"clientID": {
					"type": "string",
					"minLength": 1
				},
				"siteID": {
					"type": "string",
					"minLength": 1
				},
				"endpointID": {
					"type": "string",
					"pattern": "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
				},
				"Message": {
					"type": "object",
					"properties": {
						"notificationType": {
							"type": "string",
							"pattern": "^(DeviceUp)|(DeviceDown)$"
						}
					},
					"required": ["notificationType"]
				}
			},
			"required": ["partnerID", "clientID", "siteID", "endpointID", "Message"]
		}`

	Configurations = `
		{
			"type": "array",
			"minItems": 1,
			"items": {
				"type": "object",
				"properties": {
					"packageName": {
						"type": "string",
						"minLength": 1
					},
					"fileName": {
						"type": "string",
						"minLength": 1
					},
					"minSupportedVersion": {
						"type": "string"
					},
					"maxSupportedVersion": {
						"type": "string"
					},
					"patch": {
						"type": "array",
						"minItems": 1,
						"items": {
							"oneOf": [
								{
									"type": "object",
									"properties": {
										"op": {
											"type": "string",
											"enum": ["add", "replace", "test"]
										},
										"path": {
											"type": "string"
										},
										"value": {
											"oneOf": [
												{ "type": "boolean" },
												{ "type": "number" },
												{ "type": "string" },
												{ "type": "null" },
												{ "type": "array" },
												{ "type": "object" }
											]
										}
									},
									"required": ["op", "path", "value"],
									"additionalProperties": false
								},
								{
									"type": "object",
									"properties": {
										"op": {
											"type": "string",
											"enum": ["remove"]
										},
										"path": {
											"type": "string"
										}
									},
									"required": ["op", "path"],
									"additionalProperties": false
								},
								{
									"type": "object",
									"properties": {
										"op": {
											"type": "string",
											"enum": ["move", "copy"]
										},
										"from": {
											"type": "string"
										},
										"path": {
											"type": "string"
										}
									},
									"required": ["op", "from", "path"],
									"additionalProperties": false
								}
							]
						}
					}
				},
				"required": ["packageName"],
				"additionalProperties": false,
				"dependencies": {
					"fileName": {
					    "required":[
						  "patch"
					    ]
					},
					"patch":{
					    "required":[
						  "fileName"
					    ]
					}
				}
			}
		}`

	Profile = fmt.Sprintf(`
		{
			"type": "object",
			"properties": {
				"id": {
					"type": "string",
					"pattern": "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
				},
				"description": {
					"type": "string"
				},
				"tag": {
					"type": "string"
				},
				"targets": {
					"type": "array",
					"items": {
						"type": "object",
						"properties": {
							"partnerID": {
								"type": "string"
							},
							"clientID": {
								"type": "string"
							},
							"siteID": {
								"type": "string"
							},
							"endpointID": {
								"type": "string"
							}
						},
						"additionalProperties": false
					}
				},
				"condition": {
					"type": "string"
				},
				"configurations": %s,
				"configurationTemplate": {
					"type": "string"
				},
				"createdTime": {
					"type": "string",
					"format": "date-time"
				},
				"createdBy": {
					"type": "string"
				},
				"modifiedTime": {
					"type": "string",
					"format": "date-time"
				},
				"modifiedBy": {
					"type": "string"
				}
			},
			"required": ["tag"],
			"additionalProperties": false,
			"oneOf" : [
				{ "required": [ "configurations" ] },
   				{ "required": [ "configurationTemplate" ] }
			]
		}`, Configurations)

	Update = `{
		"type": "object",
		"properties": {
			"partnerID": {
				"type": "string",
				"minLength": 1
			},
			"clientID": {
				"type": "string",
				"minLength": 1
			},
			"siteID": {
				"type": "string",
				"minLength": 1
			},
			"endpointID": {
				"type": "string",
				"pattern": "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
			},
			"forceUpdate": {
				"type": "boolean"
			},
			"transactionID": {
				"type": "string"
			},
			"originator": {
				"type": "string"
			},
			"mailboxMsgID": {
				"type": "string"
			},
			"version": {
				"type": "string"
			},
			"status": {
				"type": "string"
			},
			"startedAt": {
				"type": "string",
				"format": "date-time"
			},
			"lastUpdated": {
				"type": "string",
				"format": "date-time"
			},
			"finishedAt": {
				"type": "string",
				"format": "date-time"
			},
			"errorDetails": {
				"type": "string"
			},
			"retryCount": {
				"type": "integer"
			},
			"profiles": {
				"type": "array",
				"items": {
					"type": "string",
					"pattern": "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
				}
			}
		},
		"required": ["partnerID", "clientID", "siteID", "endpointID", "transactionID"]
	}`
)
