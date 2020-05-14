package PthWinOSUpdatesSchema

import (
	"errors"
	"reflect"
	"strings"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/graphql-go/graphql"
)

//PthWinOSUpdatesData : PthWinOSUpdatesData Structure
type PthWinOSUpdatesData struct {
	Status         string `json:"Status"`
	PatchType      string `json:"PatchType"`
	IconName       string `json:"IconName"`
	RebootRequired string `json:"RebootRequired"`
	KBArticle      string `json:"KBArticle"`
	Order          int64  `json:"Order"`
	TestStatus     string `json:"TestStatus"`
	Product        string `json:"Product"`
	SupportURL     string `json:"SupportUrl"`
	UpdateTitle    string `json:"UpdateTitle"`
	PatchSeverity  string `json:"PatchSeverity"`
	Age            int64  `json:"Age"`
	ReleaseDate    string `json:"ReleaseDate"`
	InstalledDate  string `json:"InstalledDate"`
}

//PthWinOSUpdatesDataType : PthWinOSUpdatesData GraphQL Schema
var PthWinOSUpdatesDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PthWinOSUpdatesData",
	Fields: graphql.Fields{
		"Status": &graphql.Field{
			Type:        graphql.String,
			Description: "Status.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinOSUpdatesData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},
		"PatchType": &graphql.Field{
			Type:        graphql.String,
			Description: "PatchType.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinOSUpdatesData); ok {
					return CurData.PatchType, nil
				}
				return nil, nil
			},
		},
		"IconName": &graphql.Field{
			Type:        graphql.String,
			Description: "IconName.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinOSUpdatesData); ok {
					return CurData.IconName, nil
				}
				return nil, nil
			},
		},
		"RebootRequired": &graphql.Field{
			Type:        graphql.String,
			Description: "RebootRequired.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinOSUpdatesData); ok {
					return CurData.RebootRequired, nil
				}
				return nil, nil
			},
		},
		"KBArticle": &graphql.Field{
			Type:        graphql.String,
			Description: "KBArticle.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinOSUpdatesData); ok {
					return CurData.KBArticle, nil
				}
				return nil, nil
			},
		},
		"Order": &graphql.Field{
			Type:        graphql.String,
			Description: "Order.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinOSUpdatesData); ok {
					return CurData.Order, nil
				}
				return nil, nil
			},
		},
		"TestStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "TestStatus.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinOSUpdatesData); ok {
					return CurData.TestStatus, nil
				}
				return nil, nil
			},
		},
		"Product": &graphql.Field{
			Type:        graphql.String,
			Description: "Product.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinOSUpdatesData); ok {
					return CurData.Product, nil
				}
				return nil, nil
			},
		},
		"SupportUrl": &graphql.Field{
			Type:        graphql.String,
			Description: "SupportUrl.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinOSUpdatesData); ok {
					return CurData.SupportURL, nil
				}
				return nil, nil
			},
		},
		"UpdateTitle": &graphql.Field{
			Type:        graphql.String,
			Description: "UpdateTitle.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinOSUpdatesData); ok {
					return CurData.UpdateTitle, nil
				}
				return nil, nil
			},
		},
		"PatchSeverity": &graphql.Field{
			Type:        graphql.String,
			Description: "PatchSeverity.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinOSUpdatesData); ok {
					return CurData.PatchSeverity, nil
				}
				return nil, nil
			},
		},
		"Age": &graphql.Field{
			Type:        graphql.String,
			Description: "Age.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PthWinOSUpdatesData); ok {
					return CurData.Age, nil
				}
				return nil, nil
			},
		},
		"ReleaseDate": &graphql.Field{
			Description: "ReleaseDate.",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PthWinOSUpdatesData); ok {
					return CurrData.ReleaseDate, nil
				}
				return nil, nil
			},
		},
		"InstalledDate": &graphql.Field{
			Description: "InstalledDate.",
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurrData, ok := p.Source.(PthWinOSUpdatesData); ok {
					return CurrData.InstalledDate, nil
				}
				return nil, nil
			},
		},
	},
})

//PthWinOSUpdatesDataConnectionDefinition : PthWinOSUpdatesDataConnectionDefinition structure
var PthWinOSUpdatesDataConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "PthWinOSUpdatesData",
	NodeType: PthWinOSUpdatesDataType,
})

//PthWinOSUpdates : PthWinOSUpdates Structure
type PthWinOSUpdates struct {
	PthWinOSUpdatesList []PthWinOSUpdatesData `json:"PthWinOSUpdatesList"`
	TotalCount          int64                 `json:"totalCount"`
}

//PthWinOSUpdatesType : PthWinOSUpdates GraphQL Schema
var PthWinOSUpdatesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PthWinOSUpdates",
	Fields: graphql.Fields{
		"PthWinOSUpdatesList": &graphql.Field{
			Type:        PthWinOSUpdatesDataConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "array of PthWinOSUpdatesData.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(PthWinOSUpdates); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.PthWinOSUpdatesList {
						arraySliceRet = append(arraySliceRet, CurData.PthWinOSUpdatesList[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY {
						var err error
						val := reflect.Indirect(reflect.ValueOf(&PthWinOSUpdatesData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						ProductASC := func(p1, p2 interface{}) bool {
							return p1.(PthWinOSUpdatesData).Product < p2.(PthWinOSUpdatesData).Product
						}
						ProductDESC := func(p1, p2 interface{}) bool {
							return p1.(PthWinOSUpdatesData).Product > p2.(PthWinOSUpdatesData).Product
						}

						AgeASC := func(p1, p2 interface{}) bool {
							return p1.(PthWinOSUpdatesData).Age < p2.(PthWinOSUpdatesData).Age
						}
						AgeDESC := func(p1, p2 interface{}) bool {
							return p1.(PthWinOSUpdatesData).Age > p2.(PthWinOSUpdatesData).Age
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "PRODUCT" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(ProductASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(ProductDESC).Sort(arraySliceRet)
								}
							} else if strings.ToUpper(Column) == "AGE" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(AgeASC).Sort(arraySliceRet)
								} else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(AgeDESC).Sort(arraySliceRet)
								}
							} else {
								return nil, errors.New("PthWinOSUpdatesData Sort [" + Column + "] No such column exist!!!")
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
				if CurData, ok := p.Source.(PthWinOSUpdates); ok {
					return CurData.TotalCount, nil
				}
				return nil, nil
			},
		},
	},
})
