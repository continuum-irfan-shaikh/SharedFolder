package TargetsExplorerSchema

import (
	"errors"
	"strconv"
	"strings"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/graphql-go/graphql"
)

//NAME column in uppercase format
const NAME = "NAME"

//TargetsExplorer : TargetsExplorer structure
type TargetsExplorer struct {
	Sites        []Site        `json:"sites"`
	DeviceGroups []DeviceGroup `json:"deviceGroups"`
	Devices      []Device      `json:"devices"`
}

//Site : Site structure
type Site struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

//DeviceGroup : DeviceGroup structure
type DeviceGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//Device : Device structure
type Device struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	SiteName     string `json:"siteName"`
	TimeZone     string `json:"timeZone"`
	OSProduct    string `json:"osProduct"`
	FriendlyName string `json:"friendlyName"`
	ResourceType string `json:"resourceType"`
	RegID        string `json:"regId"`
}

//DeviceType : DeviceType GraphQL Schema
var DeviceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DeviceType",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.String,
			Description: "Device ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Device); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Device name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Device); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "Device status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Device); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},
		"siteName": &graphql.Field{
			Type:        graphql.String,
			Description: "Site name of the device",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Device); ok {
					return CurData.SiteName, nil
				}
				return nil, nil
			},
		},
		"regId": &graphql.Field{
			Type:        graphql.String,
			Description: "Reg id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Device); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},
		"timeZone": &graphql.Field{
			Type:        graphql.String,
			Description: "Time zone",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Device); ok {
					return CurData.TimeZone, nil
				}
				return nil, nil
			},
		},
		"osProduct": &graphql.Field{
			Type:        graphql.String,
			Description: "osProduct",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Device); ok {
					return CurData.OSProduct, nil
				}
				return nil, nil
			},
		},
		"friendlyName": &graphql.Field{
			Type:        graphql.String,
			Description: "friendlyName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Device); ok {
					return CurData.FriendlyName, nil
				}
				return nil, nil
			},
		},
		"resourceType": &graphql.Field{
			Type:        graphql.String,
			Description: "resourceType",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Device); ok {
					return CurData.ResourceType, nil
				}
				return nil, nil
			},
		},
	},
})

//DeviceGroupType : DeviceGroupType GraphQL Schema
var DeviceGroupType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DeviceGroupType",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.String,
			Description: "Device group ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceGroup); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Device group name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(DeviceGroup); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},
	},
})

//SiteType : SiteType GraphQL Schema
var SiteType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SiteType",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.String,
			Description: "Site ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Site); ok {
					return strconv.FormatInt(CurData.ID, 10), nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Site name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(Site); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},
	},
})

//SiteListTypeConnDef : SiteTypeConnDef structure
var SiteListTypeConnDef = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "SiteTypeConnDef",
	NodeType: SiteType,
})

//DeviceGroupListTypeConnDef : DeviceGroupTypeConnDef structure
var DeviceGroupListTypeConnDef = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "DeviceGroupTypeConnDef",
	NodeType: DeviceGroupType,
})

//DeviceListTypeConnDef : DeviceTypeConnDef structure
var DeviceListTypeConnDef = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "DeviceTypeConnDef",
	NodeType: DeviceType,
})

//TargetsExplorerType : TargetsExplorerType GraphQL Schema
var TargetsExplorerType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TargetsExplorerType",
	Fields: graphql.Fields{
		"sites": &graphql.Field{
			Type:        SiteListTypeConnDef.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "List of sites",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(TargetsExplorer); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.Sites {
						arraySliceRet = append(arraySliceRet, CurData.Sites[ind])
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							switch strings.ToUpper(Column) {
							case NAME:
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(SiteNameASC).Sort(arraySliceRet)
								} else {
									Relay.SortBy(SiteNameDESC).Sort(arraySliceRet)
								}
							default:
								return nil, errors.New("PatchPolicyData Sort [" + Column + "] No such column exist!!!")
							}

						}
					}
					return Relay.ConnectionFromArray(arraySliceRet, args, ""), nil
				}
				return nil, nil
			},
		},
		"deviceGroups": &graphql.Field{
			Type:        DeviceGroupListTypeConnDef.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "List of device groups",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(TargetsExplorer); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.DeviceGroups {
						arraySliceRet = append(arraySliceRet, CurData.DeviceGroups[ind])
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							switch strings.ToUpper(Column) {
							case NAME:
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(DeviceGroupNameASC).Sort(arraySliceRet)
								} else {
									Relay.SortBy(DeviceGroupNameDESC).Sort(arraySliceRet)
								}
							default:
								return nil, errors.New("PatchPolicyData Sort [" + Column + "] No such column exist!!!")
							}

						}
					}
					return Relay.ConnectionFromArray(arraySliceRet, args, ""), nil
				}
				return nil, nil
			},
		},
		"devices": &graphql.Field{
			Type:        DeviceListTypeConnDef.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "List of devices",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(TargetsExplorer); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.Devices {
						arraySliceRet = append(arraySliceRet, CurData.Devices[ind])
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							switch strings.ToUpper(Column) {
							case NAME:
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(DeviceNameASC).Sort(arraySliceRet)
								} else {
									Relay.SortBy(DeviceNameDESC).Sort(arraySliceRet)
								}
							default:
								return nil, errors.New("PatchPolicyData Sort [" + Column + "] No such column exist!!!")
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

//SiteNameASC : ASC sorting function for Site's Name column
func SiteNameASC(p1, p2 interface{}) bool {
	p1Name := p1.(Site).Name
	p2Name := p2.(Site).Name
	if p1Name == p2Name {
		return p1.(Site).ID < p2.(Site).ID
	}
	return Relay.StringLessOp(p1Name, p2Name)
}

//SiteNameDESC : DESC sorting function for Site's Name column
func SiteNameDESC(p1, p2 interface{}) bool {
	p1Name := p1.(Site).Name
	p2Name := p2.(Site).Name
	if p1Name == p2Name {
		return p1.(Site).ID > p2.(Site).ID
	}
	return !Relay.StringLessOp(p1Name, p2Name)
}

//DeviceGroupNameASC : ASC sorting function for DeviceGroup's Name column
func DeviceGroupNameASC(p1, p2 interface{}) bool {
	p1Name := p1.(DeviceGroup).Name
	p2Name := p2.(DeviceGroup).Name
	if p1Name == p2Name {
		return Relay.StringLessOp(p1.(DeviceGroup).ID, p2.(DeviceGroup).ID)
	}
	return Relay.StringLessOp(p1Name, p2Name)
}

//DeviceGroupNameDESC : DESC sorting function for DeviceGroup's Name column
func DeviceGroupNameDESC(p1, p2 interface{}) bool {
	p1Name := p1.(DeviceGroup).Name
	p2Name := p2.(DeviceGroup).Name
	if p1Name == p2Name {
		return !Relay.StringLessOp(p1.(DeviceGroup).ID, p2.(DeviceGroup).ID)
	}
	return !Relay.StringLessOp(p1Name, p2Name)
}

//DeviceNameASC : ASC sorting function for Device's Name column
func DeviceNameASC(p1, p2 interface{}) bool {
	p1Name := p1.(Device).Name
	p2Name := p2.(Device).Name
	if p1Name == p2Name {
		return Relay.StringLessOp(p1.(Device).ID, p2.(Device).ID)
	}
	return Relay.StringLessOp(p1Name, p2Name)
}

//DeviceNameDESC : DESC sorting function for Device's Name column
func DeviceNameDESC(p1, p2 interface{}) bool {
	p1Name := p1.(Device).Name
	p2Name := p2.(Device).Name
	if p1Name == p2Name {
		return !Relay.StringLessOp(p1.(Device).ID, p2.(Device).ID)
	}
	return !Relay.StringLessOp(p1Name, p2Name)
}
