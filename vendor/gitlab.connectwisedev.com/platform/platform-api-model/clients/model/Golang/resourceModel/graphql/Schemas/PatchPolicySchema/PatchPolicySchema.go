package PatchPolicySchema

import (
	"errors"
	"strings"
	"reflect"
	"github.com/graphql-go/graphql"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
)

//PatchPolicyData : PatchPolicyData Structure
type PatchPolicyData struct {
	SiteID           		int64  `json:"siteId"`
	RegID           		int64  `json:"regId"`
	MachineName 			string `json:"machineName"`
	Device 				string `json:"device"`
	ResFriendlyName 		string `json:"resFriendlyName"`
	Site 				string `json:"site"`
	PolicyID 			int64  `json:"policyID"`
	PatchPolicyName 		string `json:"patchPolicyName"`
	Status 				string `json:"status"`
	Order 				int64  `json:"order"`
	LastCurrentDate 		string `json:"lastCurrentDate"`
	LastAssessedDate 		string `json:"lastAssessedDate"`
	ResType 			string `json:"resType"`
	IsDeployMissing 		bool   `json:"isDeployMissing"`
	TooltipWin 			string `json:"tooltipWin"`
	TooltipMac 			string `json:"tooltipMac"`
	TooltipTP 			string `json:"tooltipTP"`
	Reboot 				string `json:"reboot"`
	Gateway 			string `json:"gateway"`
	WindowsPatch 			string `json:"windowsPatch"`
	MacPatch 			string `json:"macPatch"`
	ThirdPatch 			string `json:"thirdPatch"`
	PatchStatusWinMac 		string `json:"patchStatusWinMac"`
	PatchStatusWinMacTP 		string `json:"patchStatusWinMacTP"`
	AfterInstallsScheduleType 	int64  `json:"afterInstallsScheduleType"`
	LastSeen 			string `json:"lastSeen"`
	CreatedON 			string `json:"createdON"`
	DownloadScheduleType 		int64  `json:"downloadScheduleType"`
	InstallScheduleType 		int64  `json:"installScheduleType"`
	LastOSCurrentDate 		string `json:"lastOSCurrentDate"`
	LastTPCurrentDate 		string `json:"lastTPCurrentDate"`
	Os 				string `json:"os"`
	ExcludePatch 			bool   `json:"excludePatch"`
	ExcludeTPPatch 			bool   `json:"excludeTPPatch"`
	IsConfigured 			bool   `json:"isConfigured"`
	TpIsConfigured 			bool   `json:"tpIsConfigured"`
	TpVendorID 			string `json:"tpVendorID"`
	ConfigString 			string `json:"configString"`
	Tooltip 			string `json:"tooltip"`
	InterruptDate 			string `json:"interruptDate"`
	LastDeploymentDate 		string `json:"lastDeploymentDate"`
	PolicyChangeDate 		string `json:"policyChangeDate"`
	NoOfMissingWindows 		int64  `json:"noOfMissingWindows"`
	WindowsMissingType 		int64  `json:"windowsMissingType"`
	OsNextDeployStartDate 		string `json:"osNextDeployStartDate"`
	OsNextDeployEndDate 		string `json:"osNextDeployEndDate"`
	TpNextDeployStartDate 		string `json:"tpNextDeployStartDate"`
	TpNextDeployEndDate 		string `json:"tpNextDeployEndDate"`
	Wpla 				string `json:"wpla"`
	Tpla 				string `json:"tpla"`
	Macla 				string `json:"macla"`
	OsDisplayStatus      		string `json:"osDisplayStatus"`
	TpDisplayStatus          	string `json:"tpDisplayStatus"`
	IsNeedsAttention   		bool   `json:"isNeedsAttention"`
	ActualStatus   			string `json:"actualStatus"`
}

//PatchPolicyDataType : PatchPolicyData GraphQL Schema
var PatchPolicyDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PatchPolicyData",
	Fields: graphql.Fields{
		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "SiteID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},
		"regId": &graphql.Field{
			Type:        graphql.String,
			Description: "RegID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},
		"machineName": &graphql.Field{
			Type:        graphql.String,
			Description: "MachineName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.MachineName, nil
				}
				return nil, nil
			},
		},
		"device": &graphql.Field{
			Type:        graphql.String,
			Description: "Device",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.Device, nil
				}
				return nil, nil
			},
		},
		"resFriendlyName": &graphql.Field{
			Type:        graphql.String,
			Description: "ResFriendlyName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.ResFriendlyName, nil
				}
				return nil, nil
			},
		},
		"site": &graphql.Field{
			Type:        graphql.String,
			Description: "Site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.Site, nil
				}
				return nil, nil
			},
		},
		"policyID": &graphql.Field{
			Type:        graphql.String,
			Description: "PolicyID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.PolicyID, nil
				}
				return nil, nil
			},
		},
		"patchPolicyName": &graphql.Field{
			Type:        graphql.String,
			Description: "PatchPolicyName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.PatchPolicyName, nil
				}
				return nil, nil
			},
		},
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "Status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},
		"order": &graphql.Field{
			Type:        graphql.String,
			Description: "Order",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.Order, nil
				}
				return nil, nil
			},
		},
		"lastCurrentDate": &graphql.Field{
			Type:        graphql.String,
			Description: "LastCurrentDate",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.LastCurrentDate, nil
				}
				return nil, nil
			},
		},
		"lastAssessedDate": &graphql.Field{
			Type:        graphql.String,
			Description: "LastAssessedDate",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.LastAssessedDate, nil
				}
				return nil, nil
			},
		},
		"resType": &graphql.Field{
			Type:        graphql.String,
			Description: "ResType",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.ResType, nil
				}
				return nil, nil
			},
		},
		"isDeployMissing": &graphql.Field{
			Type:        graphql.String,
			Description: "IsDeployMissing",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.IsDeployMissing, nil
				}
				return nil, nil
			},
		},
		"tooltipWin": &graphql.Field{
			Type:        graphql.String,
			Description: "TooltipWin",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.TooltipWin, nil
				}
				return nil, nil
			},
		},
		"tooltipMac": &graphql.Field{
			Type:        graphql.String,
			Description: "TooltipMac",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.TooltipMac, nil
				}
				return nil, nil
			},
		},
		"tooltipTP": &graphql.Field{
			Type:        graphql.String,
			Description: "TooltipTP",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.TooltipTP, nil
				}
				return nil, nil
			},
		},
		"reboot": &graphql.Field{
			Type:        graphql.String,
			Description: "Reboot",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.Reboot, nil
				}
				return nil, nil
			},
		},
		"gateway": &graphql.Field{
			Type:        graphql.String,
			Description: "Gateway",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.Gateway, nil
				}
				return nil, nil
			},
		},
		"windowsPatch": &graphql.Field{
			Type:        graphql.String,
			Description: "WindowsPatch",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.WindowsPatch, nil
				}
				return nil, nil
			},
		},
		"macPatch": &graphql.Field{
			Type:        graphql.String,
			Description: "MacPatch",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.MacPatch, nil
				}
				return nil, nil
			},
		},
		"thirdPatch": &graphql.Field{
			Type:        graphql.String,
			Description: "ThirdPatch",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.ThirdPatch, nil
				}
				return nil, nil
			},
		},
		"patchStatusWinMac": &graphql.Field{
			Type:        graphql.String,
			Description: "PatchStatusWinMac",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.PatchStatusWinMac, nil
				}
				return nil, nil
			},
		},
		"patchStatusWinMacTP": &graphql.Field{
			Type:        graphql.String,
			Description: "PatchStatusWinMacTP",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.PatchStatusWinMacTP, nil
				}
				return nil, nil
			},
		},
		"afterInstallsScheduleType": &graphql.Field{
			Type:        graphql.String,
			Description: "AfterInstallsScheduleType",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.AfterInstallsScheduleType, nil
				}
				return nil, nil
			},
		},
		"lastSeen": &graphql.Field{
			Type:        graphql.String,
			Description: "LastSeen",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.LastSeen, nil
				}
				return nil, nil
			},
		},
		"createdON": &graphql.Field{
			Type:        graphql.String,
			Description: "CreatedON",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.CreatedON, nil
				}
				return nil, nil
			},
		},
		"downloadScheduleType": &graphql.Field{
			Type:        graphql.String,
			Description: "DownloadScheduleType",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.DownloadScheduleType, nil
				}
				return nil, nil
			},
		},
		"installScheduleType": &graphql.Field{
			Type:        graphql.String,
			Description: "InstallScheduleType",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.InstallScheduleType, nil
				}
				return nil, nil
			},
		},
		"lastOSCurrentDate": &graphql.Field{
			Type:        graphql.String,
			Description: "LastOSCurrentDate",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.LastOSCurrentDate, nil
				}
				return nil, nil
			},
		},
		"lastTPCurrentDate": &graphql.Field{
			Type:        graphql.String,
			Description: "LastTPCurrentDate",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.LastTPCurrentDate, nil
				}
				return nil, nil
			},
		},
		"os": &graphql.Field{
			Type:        graphql.String,
			Description: "Os",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.Os, nil
				}
				return nil, nil
			},
		},
		"excludePatch": &graphql.Field{
			Type:        graphql.String,
			Description: "ExcludePatch",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.ExcludePatch, nil
				}
				return nil, nil
			},
		},
		"excludeTPPatch": &graphql.Field{
			Type:        graphql.String,
			Description: "ExcludeTPPatch",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.ExcludeTPPatch, nil
				}
				return nil, nil
			},
		},
		"isConfigured": &graphql.Field{
			Type:        graphql.String,
			Description: "IsConfigured",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.IsConfigured, nil
				}
				return nil, nil
			},
		},
		"tpIsConfigured": &graphql.Field{
			Type:        graphql.String,
			Description: "TpIsConfigured",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.TpIsConfigured, nil
				}
				return nil, nil
			},
		},
		"tpVendorID": &graphql.Field{
			Type:        graphql.String,
			Description: "TpVendorID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.TpVendorID, nil
				}
				return nil, nil
			},
		},
		"configString": &graphql.Field{
			Type:        graphql.String,
			Description: "ConfigString",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.ConfigString, nil
				}
				return nil, nil
			},
		},
		"tooltip": &graphql.Field{
			Type:        graphql.String,
			Description: "Tooltip",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.Tooltip, nil
				}
				return nil, nil
			},
		},
		"interruptDate": &graphql.Field{
			Type:        graphql.String,
			Description: "InterruptDate",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.InterruptDate, nil
				}
				return nil, nil
			},
		},
		"lastDeploymentDate": &graphql.Field{
			Type:        graphql.String,
			Description: "LastDeploymentDate",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.LastDeploymentDate, nil
				}
				return nil, nil
			},
		},
		"policyChangeDate": &graphql.Field{
			Type:        graphql.String,
			Description: "PolicyChangeDate",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.PolicyChangeDate, nil
				}
				return nil, nil
			},
		},
		"noOfMissingWindows": &graphql.Field{
			Type:        graphql.String,
			Description: "NoOfMissingWindows",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.NoOfMissingWindows, nil
				}
				return nil, nil
			},
		},
		"windowsMissingType": &graphql.Field{
			Type:        graphql.String,
			Description: "WindowsMissingType",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.WindowsMissingType, nil
				}
				return nil, nil
			},
		},
		"osNextDeployStartDate": &graphql.Field{
			Type:        graphql.String,
			Description: "OsNextDeployStartDate",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.OsNextDeployStartDate, nil
				}
				return nil, nil
			},
		},
		"osNextDeployEndDate": &graphql.Field{
			Type:        graphql.String,
			Description: "OsNextDeployEndDate",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.OsNextDeployEndDate, nil
				}
				return nil, nil
			},
		},
		"tpNextDeployStartDate": &graphql.Field{
			Type:        graphql.String,
			Description: "TpNextDeployStartDate",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.TpNextDeployStartDate, nil
				}
				return nil, nil
			},
		},
		"tpNextDeployEndDate": &graphql.Field{
			Type:        graphql.String,
			Description: "TpNextDeployEndDate",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.TpNextDeployEndDate, nil
				}
				return nil, nil
			},
		},
		"wpla": &graphql.Field{
			Type:        graphql.String,
			Description: "Wpla",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.Wpla, nil
				}
				return nil, nil
			},
		},
		"tpla": &graphql.Field{
			Type:        graphql.String,
			Description: "Tpla",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.Tpla, nil
				}
				return nil, nil
			},
		},
		"macla": &graphql.Field{
			Type:        graphql.String,
			Description: "Macla",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.Macla, nil
				}
				return nil, nil
			},
		},
		"osDisplayStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "OsDisplayStatus",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.OsDisplayStatus, nil
				}
				return nil, nil
			},
		},
		"tpDisplayStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "TpDisplayStatus",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.TpDisplayStatus, nil
				}
				return nil, nil
			},
		},
		"isNeedsAttention": &graphql.Field{
			Type:        graphql.String,
			Description: "IsNeedsAttention",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.IsNeedsAttention, nil
				}
				return nil, nil
			},
		},
		"actualStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "ActualStatus",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicyData); ok {
					return CurData.ActualStatus, nil
				}
				return nil, nil
			},
		},
	},
})

//PatchPolicyDataConnectionDefinition : PatchPolicyDataConnectionDefinition structure
var PatchPolicyDataConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "PatchPolicyData",
	NodeType: PatchPolicyDataType,
})

//PatchPolicy : PatchPolicy Structure
type PatchPolicy struct {
	PatchPolicyDetailList []PatchPolicyData `json:"patchPolicyDetailsList"`
	TotalCount            int64             `json:"totalCount"`
}

//PatchPolicyType : PatchPolicy GraphQL Schema
var PatchPolicyType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PatchPolicy",
	Fields: graphql.Fields{
		"patchPolicyDetailList": &graphql.Field{
			Type:        PatchPolicyDataConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "array of patchPolicyData.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(PatchPolicy); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.PatchPolicyDetailList {
						arraySliceRet = append(arraySliceRet, CurData.PatchPolicyDetailList[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY{
						var err error
						val := reflect.Indirect(reflect.ValueOf(&PatchPolicyData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						PatchPolicyNameASC := func(p1, p2 interface{}) bool {
							return p1.(PatchPolicyData).PatchPolicyName < p2.(PatchPolicyData).PatchPolicyName
						}
						PatchPolicyNameDESC := func(p1, p2 interface{}) bool {
							return p1.(PatchPolicyData).PatchPolicyName > p2.(PatchPolicyData).PatchPolicyName
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "PATCHPOLICYNAME" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(PatchPolicyNameASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(PatchPolicyNameDESC).Sort(arraySliceRet)
								}
							}else{
								return nil, errors.New("PatchPolicyData Sort [" + Column + "] No such column exist!!!")
							}
						}
					}
					return Relay.ConnectionFromArray(arraySliceRet, args, ""), nil
				}
				return nil, nil
			},
		},
		"totalCount": &graphql.Field{
			Type:        graphql.String,
			Description: "total number of count.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchPolicy); ok {
					return CurData.TotalCount, nil
				}
				return nil, nil
			},
		},
	},
})
