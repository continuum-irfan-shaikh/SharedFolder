package SitesSchema

import (
	"errors"
	"reflect"
	"strings"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/graphql-go/graphql"
)

//SitesDataQA : SitesData Structure
type SitesDataQA struct {
	SiteName string `json:"siteName"`
	SiteID   string `json:"siteId"`
}

//SitesTypeQA : Sites GraphQL Schema
var SitesTypeQA = graphql.NewObject(graphql.ObjectConfig{
	Name: "SitesQA",
	Fields: graphql.Fields{
		"siteName": &graphql.Field{
			Type:        graphql.String,
			Description: "Site Name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesDataQA); ok {
					return CurData.SiteName, nil
				}
				return nil, nil
			},
		},

		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "Site ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SitesDataQA); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},
	},
})

//SitesQAConnectionDefinition : SitesQAConnectionDefinition structure
var SitesQAConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "SitesQA",
	NodeType: SitesTypeQA,
})

//SitesQAListData : SitesQAListData Structure
type SitesQAListData struct {
	Sites []SitesDataQA `json:"siteDetailsQuickAccess"`
}

//SitesQAListType : SitesQAListType GraphQL Schema
var SitesQAListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SitesQA",
	Fields: graphql.Fields{
		"siteDetailsQuickAccess": &graphql.Field{
			Type:        SitesQAConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "Sites QA list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(SitesQAListData); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.Sites {
						arraySliceRet = append(arraySliceRet, CurData.Sites[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY {
						var err error
						val := reflect.Indirect(reflect.ValueOf(&SitesDataQA{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						SiteIDASC := func(p1, p2 interface{}) bool {
							return p1.(SitesDataQA).SiteID < p2.(SitesDataQA).SiteID
						}
						SiteIDDESC := func(p1, p2 interface{}) bool {
							return p1.(SitesDataQA).SiteID > p2.(SitesDataQA).SiteID
						}

						SiteNameASC := func(p1, p2 interface{}) bool {
							return p1.(SitesDataQA).SiteName < p2.(SitesDataQA).SiteName
						}
						SiteNameDESC := func(p1, p2 interface{}) bool {
							return p1.(SitesDataQA).SiteName > p2.(SitesDataQA).SiteName
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
