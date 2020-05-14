package AssetSchema

import (
	"time"
	"errors"
	"strings"
	"reflect"
	"github.com/graphql-go/graphql"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/CustomDataTypes"
)

//FriendlyNameMutationData : FriendlyNameMutationData Structure
type FriendlyNameMutationData struct {
	FriendlyName	string `json:"friendlyname"`
}

//AssetCollectionData : AssetCollectionData Structure
type AssetCollectionData struct {
	CreateTimeUTC  		time.Time
	CreatedBy      		string
	Name           		string
	Type           		string
	EndpointID     		string
	PartnerID      		string
	ClientID       		string
	SiteID         		string
	RegID          		string
	FriendlyName   		string
	RemoteAddress  		string
	ResourceType		string
	EndpointType		string
	BaseBoard      		AssetBaseBoardData
	Bios           		AssetBiosData
	Drives         		[]AssetDriveData
	PhysicalMemory         	[]AssetMemoryData
	Networks       		[]AssetNetworkData
	Os             		AssetOsData
	Processors     		[]AssetProcessorData
	RaidController 		AssetRaidControllerData
	System         		AssetSystemData
	InstalledSoftwares	[]AssetInstalledSoftwares
	InstalledSoftwaresCnt   int64
	KeyBoards		[]AssetKeyBoardData
	Mouse			[]AssetMouseData
	Monitors		[]AssetMonitorData
	PhysicalDrives		[]AssetPhysicalDriveData
	Users			[]AssetUserData
	Services		[]AssetServiceData
	Shares			[]AssetSharesData
}

//AssetCollectionType : AssetCollection GraphQL Schema
var AssetCollectionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "assetCollection",
	Fields: graphql.Fields{
		"createTimeUTC": &graphql.Field{
			Type:        CustomDataTypes.DateTimeType,
			Description: "CreateTimeUTC of agent",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.CreateTimeUTC, nil
				}
				return nil, nil
			},
		},

		"createdBy": &graphql.Field{
			Type:        graphql.String,
			Description: "Created by user",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.CreatedBy, nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Type:        graphql.String,
			Description: "Collection of all asset information",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.Type, nil
				}
				return nil, nil
			},
		},

		"endpointID": &graphql.Field{
			Type:        graphql.String,
			Description: "Endpoint ID of the managed endpoint resource",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.EndpointID, nil
				}
				return nil, nil
			},
		},

		"partnerID": &graphql.Field{
			Type:        graphql.String,
			Description: "Partner ID of the partner",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},

		"clientID": &graphql.Field{
			Type:        graphql.String,
			Description: "Client ID or company",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.ClientID, nil
				}
				return nil, nil
			},
		},

		"siteID": &graphql.Field{
			Type:        graphql.String,
			Description: "Site ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"regID": &graphql.Field{
			Type:        graphql.String,
			Description: "Registration ID of old agent",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},

		"friendlyName": &graphql.Field{
			Type:        graphql.String,
			Description: "friendly name of managed endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.FriendlyName, nil
				}
				return nil, nil
			},
		},

		"remoteAddress": &graphql.Field{
			Type:        graphql.String,
			Description: "Public IP address of the endpoint from which HTTP request was received",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.RemoteAddress, nil
				}
				return nil, nil
			},
		},

		"resourceType": &graphql.Field{
			Type:        graphql.String,
			Description: "Type of the endpoint. Example- Desktop/Server/Firewall/Mobile Device",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.ResourceType, nil
				}
				return nil, nil
			},
		},

		"endpointType": &graphql.Field{
			Type:        graphql.String,
			Description: "Type of the endpoint. Example- Desktop/Server/Firewall/Mobile Device",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.EndpointType, nil
				}
				return nil, nil
			},
		},

		"baseBoard": &graphql.Field{
			Type:        AssetBaseBoardType,
			Description: "BaseBoard information",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.BaseBoard, nil
				}
				return nil, nil
			},
		},

		"bios": &graphql.Field{
			Type:        AssetBiosType,
			Description: "BIOS information",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.Bios, nil
				}
				return nil, nil
			},
		},

		"drives": &graphql.Field{
			Type:        graphql.NewList(AssetDriveType),
			Description: "Drive interfaces information",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.Drives, nil
				}
				return nil, nil
			},
		},

		"physicalMemory": &graphql.Field{
			Type:        graphql.NewList(AssetMemoryType),
			Description: "Memory information",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.PhysicalMemory, nil
				}
				return nil, nil
			},
		},

		"networks": &graphql.Field{
			Type:        graphql.NewList(AssetNetworkType),
			Description: "Network interfaces information",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.Networks, nil
				}
				return nil, nil
			},
		},

		"os": &graphql.Field{
			Type:        AssetOsType,
			Description: "Operating System information",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.Os, nil
				}
				return nil, nil
			},
		},

		"processors": &graphql.Field{
			Type:        graphql.NewList(AssetProcessorType),
			Description: "Processor information",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.Processors, nil
				}
				return nil, nil
			},
		},

		"raidController": &graphql.Field{
			Type:        AssetRaidControllerType,
			Description: "Raid controller information",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.RaidController, nil
				}
				return nil, nil
			},
		},

		"system": &graphql.Field{
			Type:        AssetSystemType,
			Description: "System information",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.System, nil
				}
				return nil, nil
			},
		},

		"installedSoftwares": &graphql.Field{
			Type:        graphql.NewList(AssetInstalledSoftwaresType),
			Description: "List of softwares installed on the endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.InstalledSoftwares, nil
				}
				return nil, nil
			},
		},

		"installedSoftwaresCnt": &graphql.Field{
			Type:        graphql.String,
			Description: "installed Softwares Count",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.InstalledSoftwaresCnt, nil
				}
				return nil, nil
			},
		},

		"keyBoards": &graphql.Field{
			Type:        graphql.NewList(AssetKeyBoardType),
			Description: "List of keyboards attached to the endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.KeyBoards, nil
				}
				return nil, nil
			},
		},

		"mouse": &graphql.Field{
			Type:        graphql.NewList(AssetMouseType),
			Description: "List of mouse attached to the endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.Mouse, nil
				}
				return nil, nil
			},
		},

		"monitors": &graphql.Field{
			Type:        graphql.NewList(AssetMonitorType),
			Description: "List of monitors attached to the endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.Monitors, nil
				}
				return nil, nil
			},
		},

		"physicalDrives": &graphql.Field{
			Type:        graphql.NewList(AssetPhysicalDriveType),
			Description: "List of physical drive attached to a system",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.PhysicalDrives, nil
				}
				return nil, nil
			},
		},

		"users": &graphql.Field{
			Type:        graphql.NewList(AssetUserType),
			Description: "List of User accounts on the system",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.Users, nil
				}
				return nil, nil
			},
		},

		"services": &graphql.Field{
			Type:        graphql.NewList(AssetServiceType),
			Description: "List of services on the endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.Services, nil
				}
				return nil, nil
			},
		},

		"shares": &graphql.Field{
			Type:        graphql.NewList(AssetSharesType),
			Description: "List of shared drives on the endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetCollectionData); ok {
					return CurData.Shares, nil
				}
				return nil, nil
			},
		},
	},
})

//AssetCollectionConnectionDefinition : AssetCollectionConnectionDefinition structure
var AssetCollectionConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "AssetCollection",
	NodeType: AssetCollectionType,
})

//AssetCollectionListData : AssetCollectionListData Structure
type AssetCollectionListData struct {
	Assets []AssetCollectionData `json:"assetCollectionList"`
}

//AssetCollectionListType : AssetCollectionListType GraphQL Schema
var AssetCollectionListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AssetCollectionList",
	Fields: graphql.Fields{
		"assetCollectionList": &graphql.Field{
			Type:        AssetCollectionConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "Asset Collection List",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(AssetCollectionListData); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.Assets {
						arraySliceRet = append(arraySliceRet, CurData.Assets[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY{
						var err error
						val := reflect.Indirect(reflect.ValueOf(&AssetCollectionData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						SiteIDASC := func(p1, p2 interface{}) bool {
							return p1.(AssetCollectionData).SiteID < p2.(AssetCollectionData).SiteID
						}
						SiteIDDESC := func(p1, p2 interface{}) bool {
							return p1.(AssetCollectionData).SiteID > p2.(AssetCollectionData).SiteID
						}
					
						PartnerIDASC := func(p1, p2 interface{}) bool {
							return p1.(AssetCollectionData).PartnerID < p2.(AssetCollectionData).PartnerID
						}
						PartnerIDDESC := func(p1, p2 interface{}) bool {
							return p1.(AssetCollectionData).PartnerID > p2.(AssetCollectionData).PartnerID
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "SITEID" {
								if strings.ToUpper(Key) == "ASC" {
									Relay.SortBy(SiteIDASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == "DESC" {
									Relay.SortBy(SiteIDDESC).Sort(arraySliceRet)
								}
							} else if strings.ToUpper(Column) == "PARTNERID" {
								if strings.ToUpper(Key) == "ASC" {
									Relay.SortBy(PartnerIDASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == "DESC" {
									Relay.SortBy(PartnerIDDESC).Sort(arraySliceRet)
								}
							}else{
								return nil, errors.New("AssetCollectionData Sort [" + Column + "] No such column exist!!!")
							}
						}
					}
					return Relay.ConnectionFromArray(arraySliceRet, args, ""), nil
				}
				return nil, nil
			},
		},
	},
})
