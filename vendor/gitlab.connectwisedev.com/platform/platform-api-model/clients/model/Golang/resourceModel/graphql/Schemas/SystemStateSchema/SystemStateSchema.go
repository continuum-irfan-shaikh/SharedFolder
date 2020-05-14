package SystemStateSchema

import (
	"time"
	"errors"
	"strings"
	"reflect"
	"github.com/graphql-go/graphql"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
)

//SystemStateCollection : SystemStateCollection Structure
type SystemStateCollection struct {
	CreateTimeUTC    time.Time        `json:"createTimeUTC"`
	CreatedBy        string           `json:"createdBy"`
	Name             string           `json:"name"`
	Type             string           `json:"type"`
	EndpointID       string           `json:"endpointID"`
	PartnerID        string           `json:"partnerID"`
	ClientID         string           `json:"clientID"`
	SiteID           string           `json:"siteID"`
	StartupStatus    StartupStatus    `json:"startupStatus"`
	LastLoggedOnUser LastLoggedOnUser `json:"lastLoggedOnUser"`
	LoggedOnUsers    []LoggedOnUsers  `json:"loggedOnUsers"`
}

//SystemStateCollectionType : SystemStateCollection GraphQL Schema
var SystemStateCollectionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SystemStateCollection",
	Fields: graphql.Fields{
		"createTimeUTC": &graphql.Field{
			Type:        graphql.String,
			Description: "CreateTimeUTC of agent",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SystemStateCollection); ok {
					return CurData.CreateTimeUTC, nil
				}
				return nil, nil
			},
		},

		"createdBy": &graphql.Field{
			Type:        graphql.String,
			Description: "Created by user",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SystemStateCollection); ok {
					return CurData.CreatedBy, nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SystemStateCollection); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},

		"type": &graphql.Field{
			Type:        graphql.String,
			Description: "Collection of all systemstate information",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SystemStateCollection); ok {
					return CurData.Type, nil
				}
				return nil, nil
			},
		},

		"endpointID": &graphql.Field{
			Type:        graphql.String,
			Description: "Endpoint ID of the managed endpoint resource",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SystemStateCollection); ok {
					return CurData.EndpointID, nil
				}
				return nil, nil
			},
		},

		"partnerID": &graphql.Field{
			Type:        graphql.String,
			Description: "Partner ID of partner",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SystemStateCollection); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},

		"clientID": &graphql.Field{
			Type:        graphql.String,
			Description: "Client ID or company",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SystemStateCollection); ok {
					return CurData.ClientID, nil
				}
				return nil, nil
			},
		},

		"siteID": &graphql.Field{
			Type:        graphql.String,
			Description: "Site ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SystemStateCollection); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"startupStatus": &graphql.Field{
			Type:        StartupStatusType,
			Description: "Startup Status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SystemStateCollection); ok {
					return CurData.StartupStatus, nil
				}
				return nil, nil
			},
		},

		"lastLoggedOnUser": &graphql.Field{
			Type:        LastLoggedOnUserType,
			Description: "Last Logged On User",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SystemStateCollection); ok {
					return CurData.LastLoggedOnUser, nil
				}
				return nil, nil
			},
		},

		"loggedOnUsers": &graphql.Field{
			Type:        graphql.NewList(LoggedOnUsersType),
			Description: "Logged on users information",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SystemStateCollection); ok {
					return CurData.LoggedOnUsers, nil
				}
				return nil, nil
			},
		},
	},
})

//StartupStatus : StartupStatus Structure
type StartupStatus struct {
	LastBootUpTimeUTC time.Time `json:"lastBootUpTimeUTC"`
}

//StartupStatusType : StartupStatus GraphQL Schema
var StartupStatusType = graphql.NewObject(graphql.ObjectConfig{
	Name: "StartupStatus",
	Fields: graphql.Fields{
		"lastBootUpTimeUTC": &graphql.Field{
			Type:        graphql.String,
			Description: "Last bootup time of operating system",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(StartupStatus); ok {
					return CurData.LastBootUpTimeUTC, nil
				}
				return nil, nil
			},
		},
	},
})

//LastLoggedOnUser : LastLoggedOnUser Structure
type LastLoggedOnUser struct {
	Username    string `json:"username"`
	LogonTime   string `json:"logonTime"`
	Status      string `json:"status"`
}

//LastLoggedOnUserType : LastLoggedOnUser GraphQL Schema
var LastLoggedOnUserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "LastLoggedOnUser",
	Fields: graphql.Fields{
		"username": &graphql.Field{
			Type:        graphql.String,
			Description: "User name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(LastLoggedOnUser); ok {
					return CurData.Username, nil
				}
				return nil, nil
			},
		},

		"logonTime": &graphql.Field{
			Type:        graphql.String,
			Description: "User logon Time",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(LastLoggedOnUser); ok {
					return CurData.LogonTime, nil
				}
				return nil, nil
			},
		},

		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "Status of a logon user",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(LastLoggedOnUser); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},
	},
})

//LoggedOnUsers : LoggedOnUsers Structure
type LoggedOnUsers struct {
	Username    string `json:"username"`
	LogonType   string `json:"logonType"`
	SessionID   string `json:"sessionID"`
	SessionName string `json:"sessionName"`
	Status      string `json:"status"`
	Client      string `json:"client"`
}

//LoggedOnUsersType : LoggedOnUsers GraphQL Schema
var LoggedOnUsersType = graphql.NewObject(graphql.ObjectConfig{
	Name: "LoggedOnUsers",
	Fields: graphql.Fields{
		"username": &graphql.Field{
			Type:        graphql.String,
			Description: "User name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(LoggedOnUsers); ok {
					return CurData.Username, nil
				}
				return nil, nil
			},
		},

		"logonType": &graphql.Field{
			Type:        graphql.String,
			Description: "User logon Type",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(LoggedOnUsers); ok {
					return CurData.LogonType, nil
				}
				return nil, nil
			},
		},

		"sessionID": &graphql.Field{
			Type:        graphql.String,
			Description: "Session Id of a logon user",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(LoggedOnUsers); ok {
					return CurData.SessionID, nil
				}
				return nil, nil
			},
		},

		"sessionName": &graphql.Field{
			Type:        graphql.String,
			Description: "Session name of a logon user",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(LoggedOnUsers); ok {
					return CurData.SessionName, nil
				}
				return nil, nil
			},
		},

		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "Status of a logon user",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(LoggedOnUsers); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},

		"client": &graphql.Field{
			Type:        graphql.String,
			Description: "Clien name or ip address for remote logon",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(LoggedOnUsers); ok {
					return CurData.Client, nil
				}
				return nil, nil
			},
		},
	},
})

//SystemStateConnectionDefinition : SystemStateConnectionDefinition structure
var SystemStateConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "SystemState",
	NodeType: SystemStateCollectionType,
})

//SystemStateCollectionListData : SystemStateCollectionListData Structure
type SystemStateCollectionListData struct {
	SystemStateCollection []SystemStateCollection `json:"systemStateCollectionList"`
}

//SystemStateCollectionListType : SystemStateCollectionList GraphQL Schema
var SystemStateCollectionListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SystemStateCollectionList",
	Fields: graphql.Fields{
		"systemStateCollectionList": &graphql.Field{
			Type:        SystemStateConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "SystemStateCollection List",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(SystemStateCollectionListData); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.SystemStateCollection {
						arraySliceRet = append(arraySliceRet, CurData.SystemStateCollection[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY{
						var err error
						val := reflect.Indirect(reflect.ValueOf(&SystemStateCollection{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						SiteIDASC := func(p1, p2 interface{}) bool {
							return p1.(SystemStateCollection).SiteID < p2.(SystemStateCollection).SiteID
						}
						SiteIDDESC := func(p1, p2 interface{}) bool {
							return p1.(SystemStateCollection).SiteID > p2.(SystemStateCollection).SiteID
						}
					
						EndpointIDASC := func(p1, p2 interface{}) bool {
							return p1.(SystemStateCollection).EndpointID < p2.(SystemStateCollection).EndpointID
						}
						EndpointIDDESC := func(p1, p2 interface{}) bool {
							return p1.(SystemStateCollection).EndpointID > p2.(SystemStateCollection).EndpointID
						}

						for i := range subQuery {
							var Key, Column string
							Key, Column, _ = Relay.GetQueryDetails(subQuery[i])
							if strings.ToUpper(Column) == "SITEID" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(SiteIDASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(SiteIDDESC).Sort(arraySliceRet)
								}
							} else if strings.ToUpper(Column) == "ENDPOINTID" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(EndpointIDASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(EndpointIDDESC).Sort(arraySliceRet)
								}
							}else{
								return nil, errors.New("SystemStateCollection Sort [" + Column + "] No such column exist!!!")
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
