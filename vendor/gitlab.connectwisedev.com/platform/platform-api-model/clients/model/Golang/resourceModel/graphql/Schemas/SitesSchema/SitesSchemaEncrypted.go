package SitesSchema

import (
	"errors"
	"reflect"
	"strings"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/graphql-go/graphql"
)

//SitesDataEnc : SitesData Structure
type SitesDataEnc struct {
	SiteName string `json:"siteName"`
	SiteID   string `json:"encSiteId"`
}

//SitesTypeEnc : Sites GraphQL Schema
var SitesTypeEnc = graphql.NewObject(graphql.ObjectConfig{
	Name: "Sites",
	Fields: graphql.Fields{
		"siteName": &graphql.Field{
			Type:        graphql.String,
			Description: "Site Name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesDataEnc); ok {
					return CurData.SiteName, nil
				}
				return nil, nil
			},
		},

		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "Site ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesDataEnc); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},
	},
})

//SitesEncConnectionDefinition : SitesEncConnectionDefinition structure
var SitesEncConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "Sites",
	NodeType: SitesTypeEnc,
})

//SitesEncListData : SitesEncListData Structure
type SitesEncListData struct {
	Sites []SitesDataEnc `json:"siteDetailList"`
}

//SitesEncListType : SitesEncListType GraphQL Schema
var SitesEncListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SitesEnc",
	Fields: graphql.Fields{
		"siteDetailList": &graphql.Field{
			Type:        SitesEncConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "Sites  list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(SitesEncListData); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.Sites {
						arraySliceRet = append(arraySliceRet, CurData.Sites[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY {
						var err error
						val := reflect.Indirect(reflect.ValueOf(&SitesDataEnc{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						SiteIDASC := func(p1, p2 interface{}) bool {
							return p1.(SitesDataEnc).SiteID < p2.(SitesDataEnc).SiteID
						}
						SiteIDDESC := func(p1, p2 interface{}) bool {
							return p1.(SitesDataEnc).SiteID > p2.(SitesDataEnc).SiteID
						}

						SiteNameASC := func(p1, p2 interface{}) bool {
							return p1.(SitesDataEnc).SiteName < p2.(SitesDataEnc).SiteName
						}
						SiteNameDESC := func(p1, p2 interface{}) bool {
							return p1.(SitesDataEnc).SiteName > p2.(SitesDataEnc).SiteName
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
