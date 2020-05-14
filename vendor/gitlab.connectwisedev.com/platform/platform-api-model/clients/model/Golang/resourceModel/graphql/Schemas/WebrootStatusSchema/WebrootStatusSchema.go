package WebrootStatusSchema

import (
	"reflect"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/graphql-go/graphql"
)

// WebrootStatusData : WebrootStatusData structure
type WebrootStatusData struct {
	PartnerID  string `json:"partnerId"`
	ClientID   string `json:"clientId"`
	SiteID     string `json:"siteId"`
	EndpointID string `json:"endpointId"`
	SyncStatus string `json:"syncStatus"`
}

// WebrootStatusType : WebrootStatusType GraphQL Schema
var WebrootStatusType = graphql.NewObject(graphql.ObjectConfig{
	Name: "webrootstatus",
	Fields: graphql.Fields{
		"partnerId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootStatusData); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},
		"clientId": &graphql.Field{
			Type:        graphql.String,
			Description: "Client ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootStatusData); ok {
					return CurData.ClientID, nil
				}
				return nil, nil
			},
		},
		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "Site ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootStatusData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},
		"endpointId": &graphql.Field{
			Type:        graphql.String,
			Description: "Juno endpoint ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootStatusData); ok {
					return CurData.EndpointID, nil
				}
				return nil, nil
			},
		},
		"syncStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "Sync status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(WebrootStatusData); ok {
					return CurData.SyncStatus, nil
				}
				return nil, nil
			},
		},
	},
})

// WebrootStatusConnectionDefinition : WebrootStatusConnectionDefinition structure
var WebrootStatusConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "WebrootStatus",
	NodeType: WebrootStatusType,
})

// WebrootStatusListData : WebrootStatusListData Structure
type WebrootStatusListData struct {
	WebrootsStatus []WebrootStatusData `json:"webrootsStatus"`
}

// WebrootStatusListType : WebrootStatusListType GraphQL Schema
var WebrootStatusListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "WebrootStatusList",
	Fields: graphql.Fields{
		"webrootsStatus": &graphql.Field{
			Type:        WebrootStatusConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "webroot status list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(WebrootStatusListData); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.WebrootsStatus {
						arraySliceRet = append(arraySliceRet, CurData.WebrootsStatus[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY {
						var err error
						val := reflect.Indirect(reflect.ValueOf(&WebrootStatusData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					return Relay.ConnectionFromArray(arraySliceRet, args, ""), nil
				}
				return nil, nil
			},
		},
	},
})
