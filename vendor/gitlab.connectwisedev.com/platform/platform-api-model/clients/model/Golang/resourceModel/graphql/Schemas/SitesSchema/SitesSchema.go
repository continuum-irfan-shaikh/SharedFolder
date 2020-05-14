package SitesSchema

import (
	"errors"
	"strings"
	"reflect"
	"github.com/graphql-go/graphql"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
)

//SitesData : SitesData Structure
type SitesData struct {
	SiteName         string `json:"siteName"`
	SiteCode         string `json:"siteCode"`
	SiteID           int64  `json:"siteId"`
	HelpDesk         string `json:"helpDesk"`
	DesktopOption    string `json:"desktopOption"`
	ServerOption     string `json:"serverOption"`
	ServerServiceID  int64  `json:"serverServiceId"`
	DesktopServiceID int64  `json:"desktopServiceId"`
}

//SitesType : Sites GraphQL Schema
var SitesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Sites",
	Fields: graphql.Fields{
		"siteName": &graphql.Field{
			Type:        graphql.String,
			Description: "Site Name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesData); ok {
					return CurData.SiteName, nil
				}
				return nil, nil
			},
		},

		"siteCode": &graphql.Field{
			Type:        graphql.String,
			Description: "Site Code",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesData); ok {
					return CurData.SiteCode, nil
				}
				return nil, nil
			},
		},

		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "Site ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"helpDesk": &graphql.Field{
			Type:        graphql.String,
			Description: "HelpDesk",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesData); ok {
					return CurData.HelpDesk, nil
				}
				return nil, nil
			},
		},

		"desktopOption": &graphql.Field{
			Type:        graphql.String,
			Description: "DesktopOption",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesData); ok {
					return CurData.DesktopOption, nil
				}
				return nil, nil
			},
		},

		"serverOption": &graphql.Field{
			Type:        graphql.String,
			Description: "ServerOption",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesData); ok {
					return CurData.ServerOption, nil
				}
				return nil, nil
			},
		},

		"serverServiceId": &graphql.Field{
			Type:        graphql.String,
			Description: "ServerServiceID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesData); ok {
					return CurData.ServerServiceID, nil
				}
				return nil, nil
			},
		},

		"desktopServiceId": &graphql.Field{
			Type:        graphql.String,
			Description: "DesktopServiceID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesData); ok {
					return CurData.DesktopServiceID, nil
				}
				return nil, nil
			},
		},
	},
})

//SitesConnectionDefinition : SitesConnectionDefinition structure
var SitesConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "Sites",
	NodeType: SitesType,
})

//SitesListData : SitesListData Structure
type SitesListData struct {
	Sites []SitesData `json:"siteDetailList"`
}

//SitesListType : SitesListType GraphQL Schema
var SitesListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Sites",
	Fields: graphql.Fields{
		"siteDetailList": &graphql.Field{
			Type:        SitesConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "Sites list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(SitesListData); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.Sites {
						arraySliceRet = append(arraySliceRet, CurData.Sites[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY {
						var err error
						val := reflect.Indirect(reflect.ValueOf(&SitesData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						SiteIDASC := func(p1, p2 interface{}) bool {
							return p1.(SitesData).SiteID < p2.(SitesData).SiteID
						}
						SiteIDDESC := func(p1, p2 interface{}) bool {
							return p1.(SitesData).SiteID > p2.(SitesData).SiteID
						}

						SiteNameASC := func(p1, p2 interface{}) bool {
							return p1.(SitesData).SiteName < p2.(SitesData).SiteName
						}
						SiteNameDESC := func(p1, p2 interface{}) bool {
							return p1.(SitesData).SiteName > p2.(SitesData).SiteName
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "SITEID" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(SiteIDASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(SiteIDDESC).Sort(arraySliceRet)
								}
							} else if strings.ToUpper(Column) == "SITENAME" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(SiteNameASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(SiteNameDESC).Sort(arraySliceRet)
								}
							} else {
								return nil, errors.New("SitesData Sort [" + Column + "] No such column exist!!!")
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
