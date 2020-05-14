package AssetSchema

import (
	"github.com/graphql-go/graphql"
)

//AssetServiceData : AssetServiceData Structure
type AssetServiceData struct {
	ServiceName		string
	DisplayName		string
	ExecutablePath		string
	StartupType		string
	ServiceStatus		string
	LogOnAs			string
	StopEnableAction	bool
	DelayedAutoStart	bool
	Win32ExitCode		int64
	ServiceSpecificExitCode	int64
}

//AssetServiceType : AssetService GraphQL Schema
var AssetServiceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "assetService",
	Fields: graphql.Fields{
		"serviceName": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the service",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetServiceData); ok {
					return CurData.ServiceName, nil
				}
				return nil, nil
			},
		},

		"displayName": &graphql.Field{
			Type:        graphql.String,
			Description: "Display name of the service",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetServiceData); ok {
					return CurData.DisplayName, nil
				}
				return nil, nil
			},
		},

		"executablePath": &graphql.Field{
			Type:        graphql.String,
			Description: "Path to executable",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetServiceData); ok {
					return CurData.ExecutablePath, nil
				}
				return nil, nil
			},
		},

		"startupType": &graphql.Field{
			Type:        graphql.String,
			Description: "Startup type of the service",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetServiceData); ok {
					return CurData.StartupType, nil
				}
				return nil, nil
			},
		},

		"serviceStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "Status of the service",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetServiceData); ok {
					return CurData.ServiceStatus, nil
				}
				return nil, nil
			},
		},

		"logOnAs": &graphql.Field{
			Type:        graphql.String,
			Description: "Log on as. Example- Local system",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetServiceData); ok {
					return CurData.LogOnAs, nil
				}
				return nil, nil
			},
		},

		"stopEnableAction": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Enable actions for stops",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetServiceData); ok {
					return CurData.StopEnableAction, nil
				}
				return nil, nil
			},
		},

		"delayedAutoStart": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "If True, the service is started after other auto-start services are started plus a short delay",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetServiceData); ok {
					return CurData.DelayedAutoStart, nil
				}
				return nil, nil
			},
		},

		"win32ExitCode": &graphql.Field{
			Type:        graphql.String,
			Description: "Error code that defines errors encountered in starting or stopping the service",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetServiceData); ok {
					return CurData.Win32ExitCode, nil
				}
				return nil, nil
			},
		},

		"serviceSpecificExitCode": &graphql.Field{
			Type:        graphql.String,
			Description: "Service-specific error code for errors that occur while the service is either starting or stopping",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetServiceData); ok {
					return CurData.ServiceSpecificExitCode, nil
				}
				return nil, nil
			},
		},
	},
})
