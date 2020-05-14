package AntiVirusStatusSchema

import (
	"errors"
	"strings"
	"reflect"
	"github.com/graphql-go/graphql"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
)

//AntiVirusStatusData : AntiVirusStatusData Structure
type AntiVirusStatusData struct {
	SiteID          int64  `json:"siteId"`
	PartnerID       int64  `json:"partnerId"`
	RegID           int64  `json:"regId"`
	AntiVirus       string `json:"antiVirus"`
	Version         string `json:"version"`
	AntiVirusStatus string `json:"antiVirusStatus"`
}

//AntiVirusStatusType : AntiVirusStatus GraphQL Schema
var AntiVirusStatusType = graphql.NewObject(graphql.ObjectConfig{
	Name: "antivirusstatus",
	Fields: graphql.Fields{
		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AntiVirusStatusData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"partnerId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific partner",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AntiVirusStatusData); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},

		"regId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AntiVirusStatusData); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},

		"antiVirus": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of antivirus",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AntiVirusStatusData); ok {
					return CurData.AntiVirus, nil
				}
				return nil, nil
			},
		},

		"version": &graphql.Field{
			Type:        graphql.String,
			Description: "Version of antivirus",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AntiVirusStatusData); ok {
					return CurData.Version, nil
				}
				return nil, nil
			},
		},

		"antiVirusStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "Status of antivirus",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AntiVirusStatusData); ok {
					return CurData.AntiVirusStatus, nil
				}
				return nil, nil
			},
		},
	},
})

//AntiVirusStatusConnectionDefinition : AntiVirusStatusConnectionDefinition structure
var AntiVirusStatusConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "AntiVirusStatus",
	NodeType: AntiVirusStatusType,
})

//AntiVirusStatusListData : AntiVirusStatusListData Structure
type AntiVirusStatusListData struct {
	AntiVirusStatus []AntiVirusStatusData `json:"antiVirusList"`
	TotalCount      int64                 `json:"totalCount"`
}

//AntiVirusStatusListType : AntiVirusStatusList GraphQL Schema
var AntiVirusStatusListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AntiVirusStatusList",
	Fields: graphql.Fields{
		"antiVirusStatusList": &graphql.Field{
			Type:        AntiVirusStatusConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "antivirus status list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(AntiVirusStatusListData); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.AntiVirusStatus {
						arraySliceRet = append(arraySliceRet, CurData.AntiVirusStatus[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY{
						var err error
						val := reflect.Indirect(reflect.ValueOf(&AntiVirusStatusData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						AntiVirusStatusASC := func(p1, p2 interface{}) bool {
							return p1.(AntiVirusStatusData).AntiVirusStatus < p2.(AntiVirusStatusData).AntiVirusStatus
						}
						AntiVirusStatusDESC := func(p1, p2 interface{}) bool {
							return p1.(AntiVirusStatusData).AntiVirusStatus > p2.(AntiVirusStatusData).AntiVirusStatus
						}

						RegIDASC := func(p1, p2 interface{}) bool {
							return p1.(AntiVirusStatusData).RegID < p2.(AntiVirusStatusData).RegID
						}
						RegIDDESC := func(p1, p2 interface{}) bool {
							return p1.(AntiVirusStatusData).RegID > p2.(AntiVirusStatusData).RegID
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "ANTIVIRUSSTATUS" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(AntiVirusStatusASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(AntiVirusStatusDESC).Sort(arraySliceRet)
								}
							}else if strings.ToUpper(Column) == "REGID" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(RegIDASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(RegIDDESC).Sort(arraySliceRet)
								}
							}else{
								return nil, errors.New("AntiVirusStatusData Sort [" + Column + "] No such column exist!!!")
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
			Description: "totalCount of patchStatust list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AntiVirusStatusListData); ok {
					return CurData.TotalCount, nil
				}
				return nil, nil
			},
		},
	},
})
