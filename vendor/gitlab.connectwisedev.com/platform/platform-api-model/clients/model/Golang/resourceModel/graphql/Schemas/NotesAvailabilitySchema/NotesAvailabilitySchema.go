package NotesAvailabilitySchema

import (
	"errors"
	"strings"
	"reflect"
	"github.com/graphql-go/graphql"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
)

//NotesAvailabilityData : NotesAvailabilityData Structure
type NotesAvailabilityData struct {
	RegID             int64 `json:"regId"`
	SiteID            int64 `json:"siteId"`
	PartnerID         int64 `json:"partnerId"`
	NotesAvailability int64 `json:"notesAvailability"`
}

//NotesAvailabilityType : NotesAvailability GraphQL Schema
var NotesAvailabilityType = graphql.NewObject(graphql.ObjectConfig{
	Name: "NotesAvailability",
	Fields: graphql.Fields{
		"regId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NotesAvailabilityData); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},

		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NotesAvailabilityData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"partnerId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific partner",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NotesAvailabilityData); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},

		"notesAvailability": &graphql.Field{
			Type:        graphql.String,
			Description: "Status of notes",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NotesAvailabilityData); ok {
					return CurData.NotesAvailability, nil
				}
				return nil, nil
			},
		},
	},
})

//NotesAvailabilityConnectionDefinition : NotesAvailabilityConnectionDefinition structure
var NotesAvailabilityConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "NotesAvailability",
	NodeType: NotesAvailabilityType,
})

//NotesAvailabilityListData : NotesAvailabilityListData Structure
type NotesAvailabilityListData struct {
	NotesAvailability []NotesAvailabilityData `json:"notesAvailabilityResponseList"`
	TotalCount        int64                   `json:"totalCount"`
}

//NotesAvailabilityListType : NotesAvailabilityList GraphQL Schema
var NotesAvailabilityListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "NotesAvailability",
	Fields: graphql.Fields{
		"notesAvailabilityList": &graphql.Field{
			Type:        NotesAvailabilityConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "NotesAvailability list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(NotesAvailabilityListData); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.NotesAvailability {
						arraySliceRet = append(arraySliceRet, CurData.NotesAvailability[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY{
						var err error
						val := reflect.Indirect(reflect.ValueOf(&NotesAvailabilityData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						NotesAvailabilityASC := func(p1, p2 interface{}) bool {
							return p1.(NotesAvailabilityData).NotesAvailability < p2.(NotesAvailabilityData).NotesAvailability
						}
						NotesAvailabilityDESC := func(p1, p2 interface{}) bool {
							return p1.(NotesAvailabilityData).NotesAvailability > p2.(NotesAvailabilityData).NotesAvailability
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "NOTESAVAILABILITY" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(NotesAvailabilityASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(NotesAvailabilityDESC).Sort(arraySliceRet)
								}
							}else{
								return nil, errors.New("NotesAvailabilityData Sort [" + Column + "] No such column exist!!!")
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
			Description: "totalCount of NotesAvailability list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NotesAvailabilityListData); ok {
					return CurData.TotalCount, nil
				}
				return nil, nil
			},
		},
	},
})
