package DeviceDetailsExportSchema

import (
	"time"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/AlertTicketSchema"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/AssetSchema"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/NoteDetailsSchema"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/PatchPolicySchema"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/PerformanceSchema"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/PthWinOSUpdatesSchema"
	"github.com/graphql-go/graphql"
)

//DeviceCombinedDetails : DeviceCombinedDetails base details
type DeviceCombinedDetails struct {
	SiteName            string    `json:"siteName"`
	MachineName         string    `json:"machineName"`
	FriendlyName        string    `json:"friendlyName"`
	RemoteAddress       string    `json:"remoteAddress"`
	InternalAddress     string    `json:"internalAddress"`
	Os                  string    `json:"os"`
	TimeZone            string    `json:"timeZone"`
	TimeZoneDescription string    `json:"timeZoneDescription"`
	DeviceModel         string    `json:"deviceModel"`
	Processor           string    `json:"processor"`
	SerialNumber        string    `json:"serialNumber"`
	OsType              string    `json:"osType"`
	LastRestartDate     time.Time `json:"lastRestartDate"`
}

//DeviceHeaderWidget ...
type DeviceHeaderWidget struct {
	AlertCount         int64  `json:"totalAlertsCount"`
	TicketCount        int64  `json:"totalTicketCount"`
	TaskCount          int64  `json:"totalTaskCount"`
	Availability       bool   `json:"availability"`
	OSUpdates          string `json:"osUpdates"`
	TPUpdates          string `json:"tpUpdates"`
	EndPointProtection string `json:"endPointProtection"`
}

//DeviceDetailsExportData : DeviceDetailsExportData
type DeviceDetailsExportData struct {
	DeviceHeader               DeviceHeaderWidget                                  `json:"deviceHeader"`
	Details                    DeviceCombinedDetails                               `json:"details"`
	Services                   []AssetSchema.AssetServiceData                      `json:"services"`
	PartitionData              []AssetSchema.AssetDrivePartition                   `json:"partitionData"`
	Assets                     AssetSchema.AssetCollectionData                     `json:"assets"`
	PatchPolicy                []PatchPolicySchema.PatchPolicyData                 `json:"pathcPolicy"`
	OSUpdates                  []PthWinOSUpdatesSchema.PthWinOSUpdatesData         `json:"osUpdates"`
	AlertTicketsDetailsPartner []AlertTicketSchema.AlertTicketDetails              `json:"alertTicketsdetailsPartner"`
	AlertTicketsDetailsNoc     []AlertTicketSchema.AlertTicketDetails              `json:"alertTicketsdetailsNoc"`
	Memory                     []PerformanceSchema.PerformanceMemoryData           `json:"memory"`
	Processors                 []PerformanceSchema.ProcessorsData                  `json:"processors"`
	Drives                     []PerformanceSchema.PerformanceStoragePartitionData `json:"drives"`
	Notes                      NoteDetails.NoteDetails                             `json:"notes"`
	ErrorMessage               []string                                            `json:"errorMessage"`
}

//DeviceDetailsExportType : DeviceDetailsExportData GraphQL Schema
var DeviceDetailsExportType = graphql.NewObject(graphql.ObjectConfig{
	Name: "deviceDetailsExport",
	Fields: graphql.Fields{
		"details": &graphql.Field{
			Type:        DeviceCombinedDetailsType,
			Description: "details",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsExportData); ok {
					return CurData.Details, nil
				}
				return nil, nil
			},
		},
		"deviceHeader": &graphql.Field{
			Type:        DeviceHeaderType,
			Description: "deviceHeader",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsExportData); ok {
					return CurData.DeviceHeader, nil
				}
				return nil, nil
			},
		},
		"services": &graphql.Field{
			Type:        graphql.NewList(AssetSchema.AssetServiceType),
			Description: "services",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsExportData); ok {
					return CurData.Services, nil
				}
				return nil, nil
			},
		},
		"assets": &graphql.Field{
			Type:        AssetSchema.AssetCollectionType,
			Description: "assets",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsExportData); ok {
					return CurData.Assets, nil
				}
				return nil, nil
			},
		},
		"partitionData": &graphql.Field{
			Type:        graphql.NewList(AssetSchema.AssetDrivePartitionType),
			Description: "partitionData",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsExportData); ok {
					return CurData.PartitionData, nil
				}
				return nil, nil
			},
		},
		"patchPolicy": &graphql.Field{
			Type:        graphql.NewList(PatchPolicySchema.PatchPolicyDataType),
			Description: "patchPolicy",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsExportData); ok {
					return CurData.PatchPolicy, nil
				}
				return nil, nil
			},
		},
		"osUpdates": &graphql.Field{
			Type:        graphql.NewList(PthWinOSUpdatesSchema.PthWinOSUpdatesDataType),
			Description: "osUpdates",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsExportData); ok {
					return CurData.OSUpdates, nil
				}
				return nil, nil
			},
		},
		"alertTicketsdetailsNoc": &graphql.Field{
			Type:        graphql.NewList(AlertTicketSchema.AlertTicketDetailsType),
			Description: "alertTicketsdetailsNoc",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsExportData); ok {
					return CurData.AlertTicketsDetailsNoc, nil
				}
				return nil, nil
			},
		},
		"alertTicketsdetailsPartner": &graphql.Field{
			Type:        graphql.NewList(AlertTicketSchema.AlertTicketDetailsType),
			Description: "alertTicketsdetailsPartner",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsExportData); ok {
					return CurData.AlertTicketsDetailsPartner, nil
				}
				return nil, nil
			},
		},
		"memory": &graphql.Field{
			Type:        graphql.NewList(PerformanceSchema.PerformanceMemoryType),
			Description: "memory",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsExportData); ok {
					return CurData.Memory, nil
				}
				return nil, nil
			},
		},
		"processors": &graphql.Field{
			Type:        graphql.NewList(PerformanceSchema.ProcessorsType),
			Description: "processors",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsExportData); ok {
					return CurData.Processors, nil
				}
				return nil, nil
			},
		},
		"drives": &graphql.Field{
			Type:        graphql.NewList(PerformanceSchema.PerformanceStoragePartitionType),
			Description: "drivers",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsExportData); ok {
					return CurData.Drives, nil
				}
				return nil, nil
			},
		},
		"notes": &graphql.Field{
			Type:        NoteDetails.NoteDetailsType,
			Description: "notes",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsExportData); ok {
					return CurData.Notes, nil
				}
				return nil, nil
			},
		},
		"errorMessage": &graphql.Field{
			Type:        graphql.String,
			Description: "errorMessage",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceDetailsExportData); ok {
					return CurData.ErrorMessage, nil
				}
				return nil, nil
			},
		},
	},
})

//DeviceHeaderType ...
var DeviceHeaderType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DeviceHeaderType",
	Fields: graphql.Fields{
		"totalAlertsCount": &graphql.Field{
			Description: "alertsCount",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(DeviceHeaderWidget); ok {
					return CurrData.AlertCount, nil
				}
				return nil, nil
			},
		},
		"totalTicketsCount": &graphql.Field{
			Description: "ticketCount",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(DeviceHeaderWidget); ok {
					return CurrData.TicketCount, nil
				}
				return nil, nil
			},
		},
		"totalTaskCount": &graphql.Field{
			Description: "taskCount",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(DeviceHeaderWidget); ok {
					return CurrData.TaskCount, nil
				}
				return nil, nil
			},
		},
		"availability": &graphql.Field{
			Description: "Availability",
			Type:        graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(DeviceHeaderWidget); ok {
					return CurrData.Availability, nil
				}
				return nil, nil
			},
		},
		"osUpdates": &graphql.Field{
			Description: "osUpdates",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(DeviceHeaderWidget); ok {
					return CurrData.OSUpdates, nil
				}
				return nil, nil
			},
		},
		"tpUpdates": &graphql.Field{
			Description: "tpUpdates",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(DeviceHeaderWidget); ok {
					return CurrData.TPUpdates, nil
				}
				return nil, nil
			},
		},
		"endPointProtection": &graphql.Field{
			Description: "endPointProtection",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(DeviceHeaderWidget); ok {
					return CurrData.EndPointProtection, nil
				}
				return nil, nil
			},
		},
	},
})

//DeviceCombinedDetailsType : combined details for device
var DeviceCombinedDetailsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "deviceDetailsExport",
	Fields: graphql.Fields{
		"siteName": &graphql.Field{
			Type:        graphql.String,
			Description: "siteName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceCombinedDetails); ok {
					return CurData.SiteName, nil
				}
				return nil, nil
			},
		},
		"machineName": &graphql.Field{
			Type:        graphql.String,
			Description: "machineName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceCombinedDetails); ok {
					return CurData.MachineName, nil
				}
				return nil, nil
			},
		},
		"friendlyName": &graphql.Field{
			Type:        graphql.String,
			Description: "friendlyName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceCombinedDetails); ok {
					return CurData.FriendlyName, nil
				}
				return nil, nil
			},
		},
		"remoteAddress": &graphql.Field{
			Type:        graphql.String,
			Description: "remoteAddress",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceCombinedDetails); ok {
					return CurData.RemoteAddress, nil
				}
				return nil, nil
			},
		},
		"internalAddress": &graphql.Field{
			Type:        graphql.String,
			Description: "internalAddress",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceCombinedDetails); ok {
					return CurData.InternalAddress, nil
				}
				return nil, nil
			},
		},
		"os": &graphql.Field{
			Type:        graphql.String,
			Description: "os",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceCombinedDetails); ok {
					return CurData.Os, nil
				}
				return nil, nil
			},
		},
		"timeZone": &graphql.Field{
			Type:        graphql.String,
			Description: "timeZone",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceCombinedDetails); ok {
					return CurData.TimeZone, nil
				}
				return nil, nil
			},
		},
		"timeZoneDescription": &graphql.Field{
			Type:        graphql.String,
			Description: "timeZoneDescription",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceCombinedDetails); ok {
					return CurData.TimeZoneDescription, nil
				}
				return nil, nil
			},
		},
		"deviceModel": &graphql.Field{
			Type:        graphql.String,
			Description: "deviceModel",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceCombinedDetails); ok {
					return CurData.DeviceModel, nil
				}
				return nil, nil
			},
		},
		"processor": &graphql.Field{
			Type:        graphql.String,
			Description: "processor",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceCombinedDetails); ok {
					return CurData.Processor, nil
				}
				return nil, nil
			},
		},
		"serialNumber": &graphql.Field{
			Type:        graphql.String,
			Description: "serialNumber",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceCombinedDetails); ok {
					return CurData.SerialNumber, nil
				}
				return nil, nil
			},
		},
		"lastRestartDate": &graphql.Field{
			Type:        graphql.String,
			Description: "lastRestartDate",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceCombinedDetails); ok {
					return CurData.LastRestartDate, nil
				}
				return nil, nil
			},
		},
		"osType": &graphql.Field{
			Type:        graphql.String,
			Description: "osType",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceCombinedDetails); ok {
					return CurData.OsType, nil
				}
				return nil, nil
			},
		},
	},
})
