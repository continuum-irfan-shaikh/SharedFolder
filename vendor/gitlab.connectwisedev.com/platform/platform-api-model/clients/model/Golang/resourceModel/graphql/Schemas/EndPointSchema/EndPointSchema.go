package EndPointSchema

import (
	"errors"
	"reflect"
	"strings"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/graphql-go/graphql"
)

//EndPointData : EndPointData Structure
type EndPointData struct {
	RegID                 int64  `json:"regId"`
	MachineName           string `json:"machineName"`
	FriendlyName          string `json:"friendlyName"`
	SiteName              string `json:"siteName"`
	OperatingSystem       string `json:"operatingSystem"`
	Availability          int64  `json:"availability"`
	IPAddress             string `json:"ipAddress"`
	RegType               string `json:"regType"`
	LmiStatus             int    `json:"lmiStatus"`
	ResType               string `json:"resType"`
	SiteID                int64  `json:"siteId"`
	SmartDisk             int64  `json:"smartDisk"`
	AMT                   int64  `json:"amt"`
	MBSyncstatus          int64  `json:"mbSyncstatus"`
	EncryptedResourceName string `json:"encryptedResourceName"`
	EncryptedSiteName     string `json:"encryptedSiteName"`
	LmiHostID             int64  `json:"lmiHostId"`
	RequestRegID          int64  `json:"requestRegId"`
}

//EndPointType : EndPoint GraphQL Schema
var EndPointType = graphql.NewObject(graphql.ObjectConfig{
	Name: "endpoint",
	Fields: graphql.Fields{
		"regId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointData); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},

		"machineName": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the Machine",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointData); ok {
					return CurData.MachineName, nil
				}
				return nil, nil
			},
		},

		"friendlyName": &graphql.Field{
			Type:        graphql.String,
			Description: "This is the machine friendly name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointData); ok {
					return CurData.FriendlyName, nil
				}
				return nil, nil
			},
		},

		"siteName": &graphql.Field{
			Type:        graphql.String,
			Description: "Name for the site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointData); ok {
					return CurData.SiteName, nil
				}
				return nil, nil
			},
		},

		"operatingSystem": &graphql.Field{
			Type:        graphql.String,
			Description: "Name for the Operating System",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointData); ok {
					return CurData.OperatingSystem, nil
				}
				return nil, nil
			},
		},

		"availability": &graphql.Field{
			Type:        graphql.String,
			Description: "Availability status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointData); ok {
					return CurData.Availability, nil
				}
				return nil, nil
			},
		},

		"ipAddress": &graphql.Field{
			Type:        graphql.String,
			Description: "This is the IP Address",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointData); ok {
					return CurData.IPAddress, nil
				}
				return nil, nil
			},
		},

		"regType": &graphql.Field{
			Type:        graphql.String,
			Description: "Type of the Reg",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointData); ok {
					return CurData.RegType, nil
				}
				return nil, nil
			},
		},

		"lmiStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "This is the LMI status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointData); ok {
					return CurData.LmiStatus, nil
				}
				return nil, nil
			},
		},

		"resType": &graphql.Field{
			Type:        graphql.String,
			Description: "This is the Res Type",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointData); ok {
					return CurData.ResType, nil
				}
				return nil, nil
			},
		},

		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"smartDisk": &graphql.Field{
			Type:        graphql.String,
			Description: "This is the smart disk",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointData); ok {
					return CurData.SmartDisk, nil
				}
				return nil, nil
			},
		},

		"amt": &graphql.Field{
			Type:        graphql.String,
			Description: "This is the amt",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointData); ok {
					return CurData.AMT, nil
				}
				return nil, nil
			},
		},

		"mbSyncstatus": &graphql.Field{
			Type:        graphql.String,
			Description: "This is the MB Sync status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointData); ok {
					return CurData.MBSyncstatus, nil
				}
				return nil, nil
			},
		},

		"encryptedResourceName": &graphql.Field{
			Type:        graphql.String,
			Description: "This is the encrypted resource name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointData); ok {
					return CurData.EncryptedResourceName, nil
				}
				return nil, nil
			},
		},

		"encryptedSiteName": &graphql.Field{
			Type:        graphql.String,
			Description: "This is the encrypted site name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointData); ok {
					return CurData.EncryptedSiteName, nil
				}
				return nil, nil
			},
		},

		"lmiHostId": &graphql.Field{
			Type:        graphql.String,
			Description: "This is the LMI host id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointData); ok {
					return CurData.LmiHostID, nil
				}
				return nil, nil
			},
		},

		"requestRegId": &graphql.Field{
			Type:        graphql.String,
			Description: "This is the request reg id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointData); ok {
					return CurData.RequestRegID, nil
				}
				return nil, nil
			},
		},
	},
})

//EndPointConnectionDefinition : EndPointConnectionDefinition structure
var EndPointConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "EndPoint",
	NodeType: EndPointType,
})

//EndPointListData : EndPointListData Structure
type EndPointListData struct {
	EndPoints  []EndPointData `json:"endPointList"`
	TotalCount int64          `json:"totalCount"`
}

//EndPointListType : EndPointList GraphQL Schema
var EndPointListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "EndPointList",
	Fields: graphql.Fields{
		"endPointList": &graphql.Field{
			Type:        EndPointConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "endpoint list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(EndPointListData); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.EndPoints {
						arraySliceRet = append(arraySliceRet, CurData.EndPoints[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY {
						var err error
						val := reflect.Indirect(reflect.ValueOf(&EndPointData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						LmiStatusASC := func(p1, p2 interface{}) bool {
							return p1.(EndPointData).LmiStatus < p2.(EndPointData).LmiStatus
						}
						LmiStatusDESC := func(p1, p2 interface{}) bool {
							return p1.(EndPointData).LmiStatus > p2.(EndPointData).LmiStatus
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "LMISTATUS" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(LmiStatusASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(LmiStatusDESC).Sort(arraySliceRet)
								}
							} else {
								return nil, errors.New("EndPointData Sort [" + Column + "] No such column exist!!!")
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
			Description: "totalCount of endpoint list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(EndPointListData); ok {
					return CurData.TotalCount, nil
				}
				return nil, nil
			},
		},
	},
})
