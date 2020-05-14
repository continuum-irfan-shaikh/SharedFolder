package PthWinTPUpdatesSchema

import (
	"errors"
	"strings"
	"reflect"
	"github.com/graphql-go/graphql"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
)

//PthWinTPUpdatesData : PthWinTPUpdatesData Structure
type PthWinTPUpdatesData struct {
	VendorName           string `json:"VendorName"`
	Order                int64  `json:"Order"`
	Product              string `json:"Product"`
	UpdateTitle          string `json:"UpdateTitle"`
	RebootRequired       string `json:"RebootRequired"`
	SupportInfo          string `json:"SupportInfo"`
	InstalledVersion     string `json:"InstalledVersion"`
	NextAvailableVersion string `json:"NextAvailableVersion"`
	Status               string `json:"Status"`
	IconName             string `json:"IconName"`
	Version              string `json:"Version"`
}

//PthWinTPUpdatesDataType : PthWinTPUpdatesData GraphQL Schema
var PthWinTPUpdatesDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PthWinTPUpdatesData",
	Fields: graphql.Fields{
		"VendorName": &graphql.Field{
			Type:        graphql.String,
			Description: "VendorName.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinTPUpdatesData); ok {
					return CurData.VendorName, nil
				}
				return nil, nil
			},
		},
		"Order": &graphql.Field{
			Type:        graphql.String,
			Description: "Order.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinTPUpdatesData); ok {
					return CurData.Order, nil
				}
				return nil, nil
			},
		},
		"Product": &graphql.Field{
			Type:        graphql.String,
			Description: "Product.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinTPUpdatesData); ok {
					return CurData.Product, nil
				}
				return nil, nil
			},
		},
		"UpdateTitle": &graphql.Field{
			Type:        graphql.String,
			Description: "UpdateTitle.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinTPUpdatesData); ok {
					return CurData.UpdateTitle, nil
				}
				return nil, nil
			},
		},
		"RebootRequired": &graphql.Field{
			Type:        graphql.String,
			Description: "RebootRequired.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinTPUpdatesData); ok {
					return CurData.RebootRequired, nil
				}
				return nil, nil
			},
		},
		"SupportInfo": &graphql.Field{
			Type:        graphql.String,
			Description: "SupportInfo.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinTPUpdatesData); ok {
					return CurData.SupportInfo, nil
				}
				return nil, nil
			},
		},
		"InstalledVersion": &graphql.Field{
			Type:        graphql.String,
			Description: "InstalledVersion.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinTPUpdatesData); ok {
					return CurData.InstalledVersion, nil
				}
				return nil, nil
			},
		},
		"NextAvailableVersion": &graphql.Field{
			Type:        graphql.String,
			Description: "NextAvailableVersion.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinTPUpdatesData); ok {
					return CurData.NextAvailableVersion, nil
				}
				return nil, nil
			},
		},
		"Status": &graphql.Field{
			Type:        graphql.String,
			Description: "Status.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinTPUpdatesData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},
		"IconName": &graphql.Field{
			Type:        graphql.String,
			Description: "IconName.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinTPUpdatesData); ok {
					return CurData.IconName, nil
				}
				return nil, nil
			},
		},
		"Version": &graphql.Field{
			Type:        graphql.String,
			Description: "Version.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinTPUpdatesData); ok {
					return CurData.Version, nil
				}
				return nil, nil
			},
		},
	},
})

//PthWinTPUpdatesDataConnectionDefinition : PthWinTPUpdatesDataConnectionDefinition structure
var PthWinTPUpdatesDataConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "PthWinTPUpdatesData",
	NodeType: PthWinTPUpdatesDataType,
})

//PthWinTPUpdates : PthWinTPUpdates Structure
type PthWinTPUpdates struct {
	TotalCount          int64                 `json:"totalCount"`
	PthWinTPUpdatesList []PthWinTPUpdatesData `json:"PthWinTPUpdatesList"`
}

//PthWinTPUpdatesType : PthWinTPUpdates GraphQL Schema
var PthWinTPUpdatesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PthWinTPUpdates",
	Fields: graphql.Fields{
		"PthWinTPUpdatesList": &graphql.Field{
			Type:        PthWinTPUpdatesDataConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "array of PthWinTPUpdatesData.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(PthWinTPUpdates); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.PthWinTPUpdatesList {
						arraySliceRet = append(arraySliceRet, CurData.PthWinTPUpdatesList[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY{
						var err error
						val := reflect.Indirect(reflect.ValueOf(&PthWinTPUpdatesData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						ProductASC := func(p1, p2 interface{}) bool {
							return p1.(PthWinTPUpdatesData).Product < p2.(PthWinTPUpdatesData).Product
						}
						ProductDESC := func(p1, p2 interface{}) bool {
							return p1.(PthWinTPUpdatesData).Product > p2.(PthWinTPUpdatesData).Product
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "PRODUCT" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(ProductASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(ProductDESC).Sort(arraySliceRet)
								}
							}else{
								return nil, errors.New("PthWinTPUpdatesData Sort [" + Column + "] No such column exist!!!")
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
				if CurData, ok := p.Source.(PthWinTPUpdates); ok {
					return CurData.TotalCount, nil
				}
				return nil, nil
			},
		},
	},
})
