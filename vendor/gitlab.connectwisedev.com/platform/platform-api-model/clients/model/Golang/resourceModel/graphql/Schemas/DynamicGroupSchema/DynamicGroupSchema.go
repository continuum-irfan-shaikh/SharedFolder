package DynamicGroupSchema

import (
	"errors"
	"strings"
	"reflect"
	"github.com/graphql-go/graphql"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
)

//ManagedEndpointsData : ManagedEndpointsData Structure
type ManagedEndpointsData struct {
	ID	string `json:"id"`
	Client	string `json:"client"`
	Partner	string `json:"partner"`
	Site	string `json:"site"`
}

//ManagedEndpointsType : ManagedEndpoints GraphQL Schema
var ManagedEndpointsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "managedEndpoints",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.String,
			Description: "Id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ManagedEndpointsData); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},

		"client": &graphql.Field{
			Type:        graphql.String,
			Description: "Client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ManagedEndpointsData); ok {
					return CurData.Client, nil
				}
				return nil, nil
			},
		},

		"partner": &graphql.Field{
			Type:        graphql.String,
			Description: "Partner",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ManagedEndpointsData); ok {
					return CurData.Partner, nil
				}
				return nil, nil
			},
		},

		"site": &graphql.Field{
			Type:        graphql.String,
			Description: "Site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ManagedEndpointsData); ok {
					return CurData.Site, nil
				}
				return nil, nil
			},
		},
	},
})

//ManagedEndpointsConnectionDefinition : ManagedEndpointsConnectionDefinition structure
var ManagedEndpointsConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "ManagedEndpoints",
	NodeType: ManagedEndpointsType,
})

//ManagedEndpointsListData : ManagedEndpointsListData Structure
type ManagedEndpointsListData struct {
	ManagedEndpoints []ManagedEndpointsData `json:"managedEndpointsList"`
}

//ManagedEndpointsListType : ManagedEndpointsList GraphQL Schema
var ManagedEndpointsListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ManagedEndpointsList",
	Fields: graphql.Fields{
		"managedEndpointsList": &graphql.Field{
			Type:        ManagedEndpointsConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "managed endpoints list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(ManagedEndpointsListData); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.ManagedEndpoints {
						arraySliceRet = append(arraySliceRet, CurData.ManagedEndpoints[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY{
						var err error
						val := reflect.Indirect(reflect.ValueOf(&ManagedEndpointsData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						ClientASC := func(p1, p2 interface{}) bool {
							return p1.(ManagedEndpointsData).Client < p2.(ManagedEndpointsData).Client
						}
						ClientDESC := func(p1, p2 interface{}) bool {
							return p1.(ManagedEndpointsData).Client > p2.(ManagedEndpointsData).Client
						}

						SiteASC := func(p1, p2 interface{}) bool {
							return p1.(ManagedEndpointsData).Site < p2.(ManagedEndpointsData).Site
						}
						SiteDESC := func(p1, p2 interface{}) bool {
							return p1.(ManagedEndpointsData).Site > p2.(ManagedEndpointsData).Site
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "CLIENT" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(ClientASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(ClientDESC).Sort(arraySliceRet)
								}
							}else if strings.ToUpper(Column) == "SITE" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(SiteASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(SiteDESC).Sort(arraySliceRet)
								}
							}else{
								return nil, errors.New("ManagedEndpointsData Sort [" + Column + "] No such column exist!!!")
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
