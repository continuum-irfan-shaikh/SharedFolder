package PatchOSUpdateSchema

import (
	"errors"
	"strings"
	"reflect"
	"github.com/graphql-go/graphql"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
)

//PatchOSUpdateData : PatchOSUpdateData Structure
type PatchOSUpdateData struct {
	PatchType           string `json:"patchType"`
	MacIconName         string `json:"macIconName"`
	MacRestartReq       string `json:"macRestartReq"`
	MacOrder            int64  `json:"macOrder"`
	KbArticle           string `json:"kbArticle"`
	MacWhitelist        string `json:"macWhitelist"`
	MacProduct          string `json:"macProduct"`
	MacInstallerVersion string `json:"macInstallerVersion"`
	PatchName           string `json:"patchName"`
	NocTestStatus       string `json:"nocTestStatus"`
	RestartNeeded       string `json:"restartNeeded"`
	MacStatus           string `json:"macStatus"`
	MacCurrentVersion   string `json:"macCurrentVersion"`
	PatchSeverity       string `json:"patchSeverity"`
	MacPatchName        string `json:"macPatchName"`
	PatchStatus         string `json:"patchStatus"`
}

//PatchOSUpdateDataType : PatchOSUpdateData GraphQL Schema
var PatchOSUpdateDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PatchOSUpdateData",
	Fields: graphql.Fields{
		"patchType": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a patch type.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchOSUpdateData); ok {
					return CurData.PatchType, nil
				}
				return nil, nil
			},
		},
		"macIconName": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a KBArticle.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchOSUpdateData); ok {
					return CurData.MacIconName, nil
				}
				return nil, nil
			},
		},
		"macRestartReq": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a patch type.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchOSUpdateData); ok {
					return CurData.MacRestartReq, nil
				}
				return nil, nil
			},
		},
		"macOrder": &graphql.Field{
			Type:        graphql.String,
			Description: "Patch Severity.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchOSUpdateData); ok {
					return CurData.MacOrder, nil
				}
				return nil, nil
			},
		},
		"kbArticle": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a KBArticle.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchOSUpdateData); ok {
					return CurData.KbArticle, nil
				}
				return nil, nil
			},
		},
		"macWhitelist": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing patch name.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchOSUpdateData); ok {
					return CurData.MacWhitelist, nil
				}
				return nil, nil
			},
		},
		"macProduct": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a patch type.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchOSUpdateData); ok {
					return CurData.MacProduct, nil
				}
				return nil, nil
			},
		},
		"macInstallerVersion": &graphql.Field{
			Type:        graphql.String,
			Description: "Patch Severity.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchOSUpdateData); ok {
					return CurData.MacInstallerVersion, nil
				}
				return nil, nil
			},
		},
		"patchName": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing patch name.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchOSUpdateData); ok {
					return CurData.PatchName, nil
				}
				return nil, nil
			},
		},
		"nocTestStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing patch name.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchOSUpdateData); ok {
					return CurData.NocTestStatus, nil
				}
				return nil, nil
			},
		},
		"restartNeeded": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a KBArticle.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchOSUpdateData); ok {
					return CurData.RestartNeeded, nil
				}
				return nil, nil
			},
		},
		"macStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "Patch Severity.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchOSUpdateData); ok {
					return CurData.MacStatus, nil
				}
				return nil, nil
			},
		},
		"macCurrentVersion": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a KBArticle.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchOSUpdateData); ok {
					return CurData.MacCurrentVersion, nil
				}
				return nil, nil
			},
		},
		"patchSeverity": &graphql.Field{
			Type:        graphql.String,
			Description: "Patch Severity.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchOSUpdateData); ok {
					return CurData.PatchSeverity, nil
				}
				return nil, nil
			},
		},
		"macPatchName": &graphql.Field{
			Type:        graphql.String,
			Description: "Patch Severity.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchOSUpdateData); ok {
					return CurData.MacPatchName, nil
				}
				return nil, nil
			},
		},
		"patchStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a patch type.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchOSUpdateData); ok {
					return CurData.PatchStatus, nil
				}
				return nil, nil
			},
		},
	},
})

//PatchOSUpdateDataConnectionDefinition : PatchOSUpdateDataConnectionDefinition structure
var PatchOSUpdateDataConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "PatchOSUpdateData",
	NodeType: PatchOSUpdateDataType,
})

//PatchOSUpdate : PatchOSUpdate Structure
type PatchOSUpdate struct {
	PatchOSUpdatesList []PatchOSUpdateData `json:"patchOSUpdatesList"`
	TotalCount         int64               `json:"totalCount"`
}

//PatchOSUpdateType : PatchOSUpdate GraphQL Schema
var PatchOSUpdateType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PatchOSUpdate",
	Fields: graphql.Fields{
		"patchOSUpdatesList": &graphql.Field{
			Type:        PatchOSUpdateDataConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "array of patchOSUpdateData.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(PatchOSUpdate); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.PatchOSUpdatesList {
						arraySliceRet = append(arraySliceRet, CurData.PatchOSUpdatesList[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY{
						var err error
						val := reflect.Indirect(reflect.ValueOf(&PatchOSUpdateData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						PatchTypeASC := func(p1, p2 interface{}) bool {
							return p1.(PatchOSUpdateData).PatchType < p2.(PatchOSUpdateData).PatchType
						}
						PatchTypeDESC := func(p1, p2 interface{}) bool {
							return p1.(PatchOSUpdateData).PatchType > p2.(PatchOSUpdateData).PatchType
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "PATCHTYPE" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(PatchTypeASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(PatchTypeDESC).Sort(arraySliceRet)
								}
							}else{
								return nil, errors.New("PatchOSUpdateData Sort [" + Column + "] No such column exist!!!")
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
				if CurData, ok := p.Source.(PatchOSUpdate); ok {
					return CurData.TotalCount, nil
				}
				return nil, nil
			},
		},
	},
})
