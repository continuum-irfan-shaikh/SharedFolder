package PatchTPOSUpdateSchema

import (
	"errors"
	"strings"
	"reflect"
	"github.com/graphql-go/graphql"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
)

//PatchTPOSUpdateData : PatchTPOSUpdateData Structure
type PatchTPOSUpdateData struct {
	ProductAffected      string `json:"productAffected"`
	MacStatus            string `json:"macStatus"`
	MacOrder             int64  `json:"macOrder"`
	PatchName            string `json:"patchName"`
	InstalledVersion     string `json:"installedVersion"`
	VersionToBeInstalled string `json:"versionToBeInstalled"`
	PatchStatus          string `json:"patchStatus"`
	MacProduct           string `json:"macProduct"`
	MacPatchName         string `json:"macPatchName"`
	MacCurrentVersion    string `json:"macCurrentVersion"`
	MacRestartReq        string `json:"macRestartReq"`
	VendorName           string `json:"vendorName"`
	RestartNeeded        string `json:"restartNeeded"`
	MacWhitelist         string `json:"macWhitelist"`
	MacIconName          string `json:"macIconName"`
	MacVendorName        string `json:"macVendorName"`
	MacInstallerVersion  string `json:"macInstallerVersion"`
}

//PatchTPOSUpdateDataType : PatchTPOSUpdateData GraphQL Schema
var PatchTPOSUpdateDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PatchTPOSUpdateData",
	Fields: graphql.Fields{
		"productAffected": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a KBArticle.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPOSUpdateData); ok {
					return CurData.ProductAffected, nil
				}
				return nil, nil
			},
		},
		"macStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "Patch Severity.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPOSUpdateData); ok {
					return CurData.MacStatus, nil
				}
				return nil, nil
			},
		},
		"macOrder": &graphql.Field{
			Type:        graphql.String,
			Description: "Patch Severity.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPOSUpdateData); ok {
					return CurData.MacOrder, nil
				}
				return nil, nil
			},
		},
		"patchName": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing patch name.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPOSUpdateData); ok {
					return CurData.PatchName, nil
				}
				return nil, nil
			},
		},
		"installedVersion": &graphql.Field{
			Type:        graphql.String,
			Description: "Patch Severity.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPOSUpdateData); ok {
					return CurData.InstalledVersion, nil
				}
				return nil, nil
			},
		},
		"versionToBeInstalled": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing patch name.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPOSUpdateData); ok {
					return CurData.VersionToBeInstalled, nil
				}
				return nil, nil
			},
		},
		"patchStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a patch type.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPOSUpdateData); ok {
					return CurData.PatchStatus, nil
				}
				return nil, nil
			},
		},
		"macProduct": &graphql.Field{
			Type:        graphql.String,
			Description: "Patch Severity.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPOSUpdateData); ok {
					return CurData.MacProduct, nil
				}
				return nil, nil
			},
		},
		"macPatchName": &graphql.Field{
			Type:        graphql.String,
			Description: "Patch Severity.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPOSUpdateData); ok {
					return CurData.MacPatchName, nil
				}
				return nil, nil
			},
		},
		"macCurrentVersion": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a patch type.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPOSUpdateData); ok {
					return CurData.MacCurrentVersion, nil
				}
				return nil, nil
			},
		},
		"macRestartReq": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a patch type.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPOSUpdateData); ok {
					return CurData.MacRestartReq, nil
				}
				return nil, nil
			},
		},
		"vendorName": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a patch type.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPOSUpdateData); ok {
					return CurData.VendorName, nil
				}
				return nil, nil
			},
		},
		"restartNeeded": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a KBArticle.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPOSUpdateData); ok {
					return CurData.RestartNeeded, nil
				}
				return nil, nil
			},
		},
		"macWhitelist": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing patch name.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPOSUpdateData); ok {
					return CurData.MacWhitelist, nil
				}
				return nil, nil
			},
		},
		"macIconName": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a KBArticle.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPOSUpdateData); ok {
					return CurData.MacIconName, nil
				}
				return nil, nil
			},
		},
		"macVendorName": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a patch type.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPOSUpdateData); ok {
					return CurData.MacVendorName, nil
				}
				return nil, nil
			},
		},
		"macInstallerVersion": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a KBArticle.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PatchTPOSUpdateData); ok {
					return CurData.MacInstallerVersion, nil
				}
				return nil, nil
			},
		},
	},
})

//PatchTPOSUpdateDataConnectionDefinition : PatchTPOSUpdateDataConnectionDefinition structure
var PatchTPOSUpdateDataConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "PatchTPOSUpdateData",
	NodeType: PatchTPOSUpdateDataType,
})

//PatchTPOSUpdate : PatchTPOSUpdate Structure
type PatchTPOSUpdate struct {
	PatchTPOSUpdatesList []PatchTPOSUpdateData `json:"patchTPOSUpdatesList"`
	TotalCount           int64                 `json:"totalCount"`
}

//PatchTPOSUpdateType : PatchTPOSUpdate GraphQL Schema
var PatchTPOSUpdateType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PatchTPOSUpdate",
	Fields: graphql.Fields{
		"patchTPOSUpdatesList": &graphql.Field{
			Type:        PatchTPOSUpdateDataConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "array of patchTPOSUpdateData.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(PatchTPOSUpdate); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.PatchTPOSUpdatesList {
						arraySliceRet = append(arraySliceRet, CurData.PatchTPOSUpdatesList[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY{
						var err error
						val := reflect.Indirect(reflect.ValueOf(&PatchTPOSUpdateData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						ProductAffectedASC := func(p1, p2 interface{}) bool {
							return p1.(PatchTPOSUpdateData).ProductAffected < p2.(PatchTPOSUpdateData).ProductAffected
						}
						ProductAffectedDESC := func(p1, p2 interface{}) bool {
							return p1.(PatchTPOSUpdateData).ProductAffected > p2.(PatchTPOSUpdateData).ProductAffected
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "PRODUCTAFFECTED" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(ProductAffectedASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(ProductAffectedDESC).Sort(arraySliceRet)
								}
							}else{
								return nil, errors.New("PatchTPOSUpdateData Sort [" + Column + "] No such column exist!!!")
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
				if CurData, ok := p.Source.(PatchTPOSUpdate); ok {
					return CurData.TotalCount, nil
				}
				return nil, nil
			},
		},
	},
})
