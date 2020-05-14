package MobileAndOtherDevicesAvailabilitySchema

import (
	"github.com/graphql-go/graphql"
)

//MobileAndOtherDeviceDetailsData : MobileAndOtherDeviceDetailsData Structure
// Site id, onlineDevices, totalDevices
type MobileAndOtherDeviceDetailsData struct {
	OfflineDevices int `json:"offlineDevices"`
	OnlineDevices  int `json:"onlineDevices"`
	TotalDevices   int `json:"totalDevices"`
}

//MobileAndOtherDeviceDetailsType : MobileAndOtherDeviceDetails GraphQL Schema
var MobileAndOtherDeviceDetailsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "mobileAndOtherDeviceDetails",
	Fields: graphql.Fields{
		"offlineDevices": &graphql.Field{
			Type:        graphql.String,
			Description: "offline devices",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(MobileAndOtherDeviceDetailsData); ok {
					return CurData.OfflineDevices, nil
				}
				return nil, nil
			},
		},

		"totalDevices": &graphql.Field{
			Type:        graphql.String,
			Description: "total devices",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(MobileAndOtherDeviceDetailsData); ok {
					return CurData.TotalDevices, nil
				}
				return nil, nil
			},
		},

		"onlineDevices": &graphql.Field{
			Type:        graphql.String,
			Description: "online devices",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(MobileAndOtherDeviceDetailsData); ok {
					return CurData.OnlineDevices, nil
				}
				return nil, nil
			},
		},
	},
})

//MobileAndOtherDevicesAvailabilityData : MobileAndOtherDevicesAvailabilityData Structure
type MobileAndOtherDevicesAvailabilityData struct {
	Status       int64                             `json:"status"`
	DeviceStatus []MobileAndOtherDeviceDetailsData `json:"outdata"`
}

//MobileAndOtherDevicesAvailabilityType : MobileAndOtherDevicesAvailabilityType GraphQL Schema
var MobileAndOtherDevicesAvailabilityType = graphql.NewObject(graphql.ObjectConfig{
	Name: "MobileAndOtherDevicesAvailability",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "status code",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(MobileAndOtherDevicesAvailabilityData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},

		"outdata": &graphql.Field{
			Type:        graphql.NewList(MobileAndOtherDeviceDetailsType),
			Description: "devices status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(MobileAndOtherDevicesAvailabilityData); ok {
					return CurData.DeviceStatus, nil
				}
				return nil, nil
			},
		},
	},
})
