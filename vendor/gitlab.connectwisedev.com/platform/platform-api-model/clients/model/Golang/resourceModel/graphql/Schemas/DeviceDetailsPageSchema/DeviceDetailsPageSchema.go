package DeviceDetailsPageSchema

import (
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/AssetSchema"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/PerformanceSchema"
	"github.com/graphql-go/graphql"
)

//DeviceMachineSpecsDetails : DeviceMachineSpecsDetails Machine Specs details
type DeviceMachineSpecsDetails struct {
	Location     string                        `json:"location"`
	SiteName     string                        `json:"siteName"`
	MachineName  string                        `json:"machineName"`
	FriendlyName string                        `json:"friendlyName"`
	Os           string                        `json:"os"`
	DeviceModel  string                        `json:"deviceModel"`
	RAM          []AssetSchema.AssetMemoryData `json:"PhysicalMemory"`
	Processor    string                        `json:"processor"`
	SerialNumber string                        `json:"serialNumber"`
	HardDrive    string                        `json:"hardDrive"`
}

//DeviceDetailsPageData : DeviceDetailsPageData
type DeviceDetailsPageData struct {
	MachineSpecs DeviceMachineSpecsDetails                 `json:"machineSpecs"`
	Memory       []PerformanceSchema.PerformanceMemoryData `json:"memory"`
	ErrorMessage []string                                  `json:"errorMessage"`
}

//DeviceDetailsPageType : DeviceDetailsPageData GraphQL Schema
var DeviceDetailsPageType = graphql.NewObject(graphql.ObjectConfig{
	Name: "deviceDetailsPage",
	Fields: graphql.Fields{
		"machineSpecs": &graphql.Field{
			Type:        DeviceDetailsMachineSpecsType,
			Description: "machineSpecs",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsPageData); ok {
					return CurData.MachineSpecs, nil
				}
				return nil, nil
			},
		},
		"memory": &graphql.Field{
			Type:        graphql.NewList(PerformanceSchema.PerformanceMemoryType),
			Description: "memory",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsPageData); ok {
					return CurData.Memory, nil
				}
				return nil, nil
			},
		},
		"errorMessage": &graphql.Field{
			Type:        graphql.String,
			Description: "errorMessage",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsPageData); ok {
					return CurData.ErrorMessage, nil
				}
				return nil, nil
			},
		},
	},
})

//DeviceDetailsMachineSpecsType : combined details for device
var DeviceDetailsMachineSpecsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "deviceDetailsMachineSpecs",
	Fields: graphql.Fields{
		"location": &graphql.Field{
			Type:        graphql.String,
			Description: "location",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceMachineSpecsDetails); ok {
					return CurData.Location, nil
				}
				return nil, nil
			},
		},
		"siteName": &graphql.Field{
			Type:        graphql.String,
			Description: "siteName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceMachineSpecsDetails); ok {
					return CurData.SiteName, nil
				}
				return nil, nil
			},
		},
		"machineName": &graphql.Field{
			Type:        graphql.String,
			Description: "machineName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceMachineSpecsDetails); ok {
					return CurData.MachineName, nil
				}
				return nil, nil
			},
		},
		"friendlyName": &graphql.Field{
			Type:        graphql.String,
			Description: "friendlyName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceMachineSpecsDetails); ok {
					return CurData.FriendlyName, nil
				}
				return nil, nil
			},
		},
		"os": &graphql.Field{
			Type:        graphql.String,
			Description: "os",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceMachineSpecsDetails); ok {
					return CurData.Os, nil
				}
				return nil, nil
			},
		},
		"deviceModel": &graphql.Field{
			Type:        graphql.String,
			Description: "deviceModel",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceMachineSpecsDetails); ok {
					return CurData.DeviceModel, nil
				}
				return nil, nil
			},
		},
		"ram": &graphql.Field{
			Type:        graphql.String,
			Description: "ram",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceMachineSpecsDetails); ok {
					return CurData.RAM, nil
				}
				return nil, nil
			},
		},
		"processor": &graphql.Field{
			Type:        graphql.String,
			Description: "processor",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceMachineSpecsDetails); ok {
					return CurData.Processor, nil
				}
				return nil, nil
			},
		},
		"serialNumber": &graphql.Field{
			Type:        graphql.String,
			Description: "serialNumber",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceMachineSpecsDetails); ok {
					return CurData.SerialNumber, nil
				}
				return nil, nil
			},
		},
		"hardDrive": &graphql.Field{
			Type:        graphql.String,
			Description: "hardDrive",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceMachineSpecsDetails); ok {
					return CurData.HardDrive, nil
				}
				return nil, nil
			},
		},
	},
})
