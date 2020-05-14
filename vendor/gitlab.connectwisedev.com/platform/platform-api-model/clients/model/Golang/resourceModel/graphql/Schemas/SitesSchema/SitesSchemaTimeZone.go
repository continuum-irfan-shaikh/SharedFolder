package SitesSchema

import (
	"errors"
	"reflect"
	"strings"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/graphql-go/graphql"
)

//SitesDataTimeZone : SitesData Structure
type SitesDataTimeZone struct {
	SiteName         string  `json:"siteName"`
	SiteCode         string  `json:"siteCode"`
	SiteID           int64   `json:"siteId"`
	HelpDesk         string  `json:"helpDesk"`
	DesktopOption    string  `json:"desktopOption"`
	ServerOption     string  `json:"serverOption"`
	ServerServiceID  int64   `json:"serverServiceId"`
	DesktopServiceID int64   `json:"desktopServiceId"`
	TimeZone         string  `json:"timezone"`
	TimeDiff         float64 `json:"timediff"`
	ZoneDisplayName  string  `json:"zonedisplayname"`
}

//SitesTypeTimeZone : Sites GraphQL Schema
var SitesTypeTimeZone = graphql.NewObject(graphql.ObjectConfig{
	Name: "SitesTimeZone",
	Fields: graphql.Fields{
		"siteName": &graphql.Field{
			Type:        graphql.String,
			Description: "Site Name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesDataTimeZone); ok {
					return CurData.SiteName, nil
				}
				return nil, nil
			},
		},

		"siteCode": &graphql.Field{
			Type:        graphql.String,
			Description: "Site Code",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesDataTimeZone); ok {
					return CurData.SiteCode, nil
				}
				return nil, nil
			},
		},

		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "Site ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesDataTimeZone); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"helpDesk": &graphql.Field{
			Type:        graphql.String,
			Description: "HelpDesk",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesDataTimeZone); ok {
					return CurData.HelpDesk, nil
				}
				return nil, nil
			},
		},

		"desktopOption": &graphql.Field{
			Type:        graphql.String,
			Description: "DesktopOption",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesDataTimeZone); ok {
					return CurData.DesktopOption, nil
				}
				return nil, nil
			},
		},

		"serverOption": &graphql.Field{
			Type:        graphql.String,
			Description: "ServerOption",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesDataTimeZone); ok {
					return CurData.ServerOption, nil
				}
				return nil, nil
			},
		},

		"serverServiceId": &graphql.Field{
			Type:        graphql.String,
			Description: "ServerServiceID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesDataTimeZone); ok {
					return CurData.ServerServiceID, nil
				}
				return nil, nil
			},
		},

		"desktopServiceId": &graphql.Field{
			Type:        graphql.String,
			Description: "DesktopServiceID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesDataTimeZone); ok {
					return CurData.DesktopServiceID, nil
				}
				return nil, nil
			},
		},

		"timeZone": &graphql.Field{
			Type:        graphql.String,
			Description: "TimeZone",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesDataTimeZone); ok {
					return CurData.TimeZone, nil
				}
				return nil, nil
			},
		},

		"timeDiff": &graphql.Field{
			Type:        graphql.String,
			Description: "TimeDiff",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesDataTimeZone); ok {
					return CurData.TimeDiff, nil
				}
				return nil, nil
			},
		},

		"zoneDisplayName": &graphql.Field{
			Type:        graphql.String,
			Description: "ZoneDisplayName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesDataTimeZone); ok {
					return CurData.ZoneDisplayName, nil
				}
				return nil, nil
			},
		},
	},
})

//SitesTimeZoneConnectionDefinition : SitesTimeZoneConnectionDefinition structure
var SitesTimeZoneConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "SitesTimeZone",
	NodeType: SitesTypeTimeZone,
})

//SitesTimeZoneListData : SitesTimeZoneListData Structure
type SitesTimeZoneListData struct {
	Sites []SitesDataTimeZone `json:"siteDetailList"`
}

//SitesTimeZoneListType : SitesTimeZoneListType GraphQL Schema
var SitesTimeZoneListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SitesTimeZone",
	Fields: graphql.Fields{
		"siteTimeZoneDetailList": &graphql.Field{
			Type:        SitesTimeZoneConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "Sites TimeZone list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(SitesTimeZoneListData); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.Sites {
						arraySliceRet = append(arraySliceRet, CurData.Sites[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY {
						var err error
						val := reflect.Indirect(reflect.ValueOf(&SitesDataTimeZone{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						SiteIDASC := func(p1, p2 interface{}) bool {
							return p1.(SitesDataTimeZone).SiteID < p2.(SitesDataTimeZone).SiteID
						}
						SiteIDDESC := func(p1, p2 interface{}) bool {
							return p1.(SitesDataTimeZone).SiteID > p2.(SitesDataTimeZone).SiteID
						}

						SiteNameASC := func(p1, p2 interface{}) bool {
							return p1.(SitesDataTimeZone).SiteName < p2.(SitesDataTimeZone).SiteName
						}
						SiteNameDESC := func(p1, p2 interface{}) bool {
							return p1.(SitesDataTimeZone).SiteName > p2.(SitesDataTimeZone).SiteName
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
