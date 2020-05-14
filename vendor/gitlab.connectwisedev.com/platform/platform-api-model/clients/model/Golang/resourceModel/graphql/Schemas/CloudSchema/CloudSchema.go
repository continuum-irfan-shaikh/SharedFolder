package CloudSchema

import (
	"reflect"
	"strings"

	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
	"github.com/graphql-go/graphql"
)

//CloudResourceHierarchy hierarchy of the resources
type CloudResourceHierarchy struct {
	Title string `json:"title"`
	Name  string `json:"name"`
	ID    string `json:"id"`
	Level int    `json:"level"`
}

//CloudResourceSummary structure for a cloud resource and it health information
type CloudResourceSummary struct {
	ID                string                   `json:"id"`
	Name              string                   `json:"name"`
	ServiceType       string                   `json:"service"`
	Location          string                   `json:"location"`
	ClientID          string                   `json:"clientid"`
	AccountID         string                   `json:"accountid"`
	PartnerID         string                   `json:"partnerid"`
	AccountName       string                   `json:"accountname"`
	AvailablityStatus string                   `json:"availabilitystatus"`
	Hierarchies       []CloudResourceHierarchy `json:"hierarchies"`
	ServiceName       string                   `json:"servicename"`
	SiteName          string
	SiteID            int64
}

//CloudResourceHealthDetails represents a cloud resource
type CloudResourceHealthDetails struct {
	ID                 string                   `json:"id"`
	Name               string                   `json:"name"`
	AvailabilityStatus string                   `json:"availabilitystatus"`
	ImpactingEvents    []CloudImpactingEvent    `json:"impactingevents"`
	RecommendedActions []CloudRecommendedAction `json:"recommendedactions"`
	TicketID           int64
	TicketStatus       string
}

//CloudImpactingEvent details of services
type CloudImpactingEvent struct {
	Summary           string `json:"summary"`
	Description       string `json:"description"`
	Event             string `json:"event"`
	EventType         string `json:"eventtype"`
	Location          string `json:"location"`
	OccurrenceTimeUTC string `json:"occurrencetimeutc"`
	LastUpdateTimeUTC string `json:"lastupdatetimeutc"`
	ResolutionETA     string `json:"resolutioneta"`
	Status            string `json:"status"`
}

//CloudRecommendedAction recommended actions
type CloudRecommendedAction struct {
	Description string `json:"description"`
	URL         string `json:"url"`
	URLTitle    string `json:"urltitle"`
}

//CloudResourceHealthDetailsRecommendedActionsDataType ...
var CloudResourceHealthDetailsRecommendedActionsDataType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "CloudResourceHealthDetailsRecommendedActionsDataType",
		Fields: graphql.Fields{
			"Description": &graphql.Field{
				Type:        graphql.String,
				Description: "description",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudRecommendedAction); ok {
						return CurData.Description, nil
					}
					return nil, nil
				},
			},
			"URL": &graphql.Field{
				Type:        graphql.String,
				Description: "description",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudRecommendedAction); ok {
						return CurData.URL, nil
					}
					return nil, nil
				},
			},
			"URLTitle": &graphql.Field{
				Type:        graphql.String,
				Description: "URL Title",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudRecommendedAction); ok {
						return CurData.URLTitle, nil
					}
					return nil, nil
				},
			},
		},
	},
)

//CloudResourceHealthDetailsImpactingEventsDataType ...
var CloudResourceHealthDetailsImpactingEventsDataType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "CloudResourceHealthDetailsImpactingEventsDataType",
		Fields: graphql.Fields{
			"Summary": &graphql.Field{
				Type:        graphql.String,
				Description: "summary",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudImpactingEvent); ok {
						return CurData.Summary, nil
					}
					return nil, nil
				},
			},
			"Description": &graphql.Field{
				Type:        graphql.String,
				Description: "description",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudImpactingEvent); ok {
						return CurData.Description, nil
					}
					return nil, nil
				},
			},
			"Event": &graphql.Field{
				Type:        graphql.String,
				Description: "ID",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudImpactingEvent); ok {
						return CurData.Event, nil
					}
					return nil, nil
				},
			},
			"EventType": &graphql.Field{
				Type:        graphql.String,
				Description: "eventtype",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudImpactingEvent); ok {
						return CurData.EventType, nil
					}
					return nil, nil
				},
			},
			"Location": &graphql.Field{
				Type:        graphql.String,
				Description: "location",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudImpactingEvent); ok {
						return CurData.Location, nil
					}
					return nil, nil
				},
			},
			"OccurrenceTimeUTC": &graphql.Field{
				Type:        graphql.String,
				Description: "occurrencetimeutc",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudImpactingEvent); ok {
						return CurData.OccurrenceTimeUTC, nil
					}
					return nil, nil
				},
			},
			"ResolutionETA": &graphql.Field{
				Type:        graphql.String,
				Description: "resolutioneta",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudImpactingEvent); ok {
						return CurData.ResolutionETA, nil
					}
					return nil, nil
				},
			},
			"Status": &graphql.Field{
				Type:        graphql.String,
				Description: "status",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudImpactingEvent); ok {
						return CurData.Status, nil
					}
					return nil, nil
				},
			},
		},
	},
)

//CloudResourceHealthDetailsDataType ...
var CloudResourceHealthDetailsDataType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "CloudResourceHealthDetailsDataType",
		Fields: graphql.Fields{
			"ID": &graphql.Field{
				Type:        graphql.String,
				Description: "id for the resource",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudResourceHealthDetails); ok {
						return CurData.ID, nil
					}
					return nil, nil
				},
			},
			"Name": &graphql.Field{
				Type:        graphql.String,
				Description: "name of the resource",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudResourceHealthDetails); ok {
						return CurData.Name, nil
					}
					return nil, nil
				},
			},
			"AvailabilityStatus": &graphql.Field{
				Type:        graphql.String,
				Description: "Availability Status of the resource",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudResourceHealthDetails); ok {
						return CurData.AvailabilityStatus, nil
					}
					return nil, nil
				},
			},
			"ImpactingEvents": &graphql.Field{
				Type:        graphql.NewList(CloudResourceHealthDetailsImpactingEventsDataType),
				Description: "collection of ImpactingEvents",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudResourceHealthDetails); ok {
						return CurData.ImpactingEvents, nil
					}
					return nil, nil
				},
			},
			"RecommendedActions": &graphql.Field{
				Type:        graphql.NewList(CloudResourceHealthDetailsRecommendedActionsDataType),
				Description: "collection of recommended actions",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudResourceHealthDetails); ok {
						return CurData.RecommendedActions, nil
					}
					return nil, nil
				},
			},
			"TicketID": &graphql.Field{
				Type:        graphql.String,
				Description: "TicketID of the resource",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudResourceHealthDetails); ok {
						return CurData.TicketID, nil
					}
					return nil, nil
				},
			},
			"TicketStatus": &graphql.Field{
				Type:        graphql.String,
				Description: "Status of TicketID",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudResourceHealthDetails); ok {
						return CurData.TicketStatus, nil
					}
					return nil, nil
				},
			},
		},
	},
)

//CloudResourceHierarchyDataType data type for cloud hierarchy
var CloudResourceHierarchyDataType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "CloudResourceSummaryDataType",
		Fields: graphql.Fields{
			"Title": &graphql.Field{
				Type:        graphql.String,
				Description: "title",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudResourceHierarchy); ok {
						return CurData.Title, nil
					}
					return nil, nil
				},
			},
			"Name": &graphql.Field{
				Type:        graphql.String,
				Description: "name",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudResourceHierarchy); ok {
						return CurData.Name, nil
					}
					return nil, nil
				},
			},
			"ID": &graphql.Field{
				Type:        graphql.String,
				Description: "ID",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudResourceHierarchy); ok {
						return CurData.ID, nil
					}
					return nil, nil
				},
			},
			"Level": &graphql.Field{
				Type:        graphql.String,
				Description: "Level",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudResourceHierarchy); ok {
						return CurData.Level, nil
					}
					return nil, nil
				},
			},
		},
	},
)

//CloudResourceSummaryDataConnectionDefinition : CloudResourceSummaryDataConnectionDefinition structure
var CloudResourceSummaryDataConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "Resources",
	NodeType: CloudResourceSummaryDataType,
})

//CloudResourceSummaryDataType defines the graph ql data type for cloudresourcesummary
var CloudResourceSummaryDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CloudResourceSummaryDataType",
	Fields: graphql.Fields{
		"ID": &graphql.Field{
			Type:        graphql.String,
			Description: "id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CloudResourceSummary); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},
		"Name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CloudResourceSummary); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},
		"ServiceType": &graphql.Field{
			Type:        graphql.String,
			Description: "ServiceType",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CloudResourceSummary); ok {
					return CurData.ServiceType, nil
				}
				return nil, nil
			},
		},
		"Location": &graphql.Field{
			Type:        graphql.String,
			Description: "Location",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CloudResourceSummary); ok {
					return CurData.Location, nil
				}
				return nil, nil
			},
		},
		"Hierarchies": &graphql.Field{
			Type:        graphql.NewList(CloudResourceHierarchyDataType),
			Description: "Hierarchies",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CloudResourceSummary); ok {
					return CurData.Hierarchies, nil
				}
				return nil, nil
			},
		},
		"AvailablityStatus": &graphql.Field{
			Type:        graphql.String,
			Description: "AvailablityStatus",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CloudResourceSummary); ok {
					return CurData.AvailablityStatus, nil
				}
				return nil, nil
			},
		},
		"ClientID": &graphql.Field{
			Type:        graphql.String,
			Description: "ClientID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CloudResourceSummary); ok {
					return CurData.ClientID, nil
				}
				return nil, nil
			},
		},
		"PartnerID": &graphql.Field{
			Type:        graphql.String,
			Description: "PartnerID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CloudResourceSummary); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},
		"AccountID": &graphql.Field{
			Type:        graphql.String,
			Description: "AccountID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CloudResourceSummary); ok {
					return CurData.AccountID, nil
				}
				return nil, nil
			},
		},
		"AccountName": &graphql.Field{
			Type:        graphql.String,
			Description: "AccountName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CloudResourceSummary); ok {
					return CurData.AccountName, nil
				}
				return nil, nil
			},
		},
		"SiteName": &graphql.Field{
			Type:        graphql.String,
			Description: "SiteName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CloudResourceSummary); ok {
					return CurData.SiteName, nil
				}
				return nil, nil
			},
		},
		"SiteID": &graphql.Field{
			Type:        graphql.String,
			Description: "SiteID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CloudResourceSummary); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},
	},
})

//CloudResourceSummaryList defines the list of resources
type CloudResourceSummaryList struct {
	Resources    []CloudResourceSummary `json:"cloudResourceSummary"`
	ErrorMessage []string               `json:"errorMessage"`
}

//CloudResourceSummaryListDataType defines the graph ql data type for cloudresourcesummary
var CloudResourceSummaryListDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CloudResourceSummaryDataType",
	Fields: graphql.Fields{
		"Resources": &graphql.Field{
			Type:        graphql.NewList(CloudResourceSummaryDataType),
			Description: "Resources",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CloudResourceSummaryList); ok {
					return CurData.Resources, nil
				}
				return nil, nil
			},
		},
		"ErrorMessage": &graphql.Field{
			Type:        graphql.String,
			Description: "ErrorMessage",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CloudResourceSummaryList); ok {
					return CurData.ErrorMessage, nil
				}
				return nil, nil
			},
		},
	},
})

//CloudResourceClientSummaryListDataType defines the graph ql data type for cloudresourcesummary
var CloudResourceClientSummaryListDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CloudResourceSummaryDataType",
	Fields: graphql.Fields{
		"Resources": &graphql.Field{
			//Type:        graphql.NewList(CloudResourceSummaryDataType),
			Type:        CloudResourceSummaryDataConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "Resources",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(CloudResourceSummaryList); ok {
					var arrResData []interface{}
					for ind := range CurData.Resources {
						arrResData = append(arrResData, CurData.Resources[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY {
						var err error
						val := reflect.Indirect(reflect.ValueOf(&CloudResourceSummary{}))
						spliteResult := strings.Split(string(args.Filter), "&")
						if len(spliteResult) > 1 {
							for j := 0; j < len(spliteResult); j++ {
								if len(spliteResult[j]) > 0 {
									arrResData, err = Relay.Filter(string(spliteResult[j]), val, arrResData)
									if err != nil {
										return nil, err
									}
								}
							}

						} else {
							arrResData, err = Relay.Filter(string(args.Filter), val, arrResData)
							if err != nil {
								return nil, err
							}
						}
					}

					if args.Sort != "" && args.Sort != Relay.NILQUERY {

					}
					return Relay.ConnectionFromArray(arrResData, args, ""), nil
					//return CurData.Resources, nil
				}
				return nil, nil
			},
		},
		"ErrorMessage": &graphql.Field{
			Type:        graphql.String,
			Description: "ErrorMessage",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CloudResourceSummaryList); ok {
					return CurData.ErrorMessage, nil
				}
				return nil, nil
			},
		},
	},
})

//CloudResourceStatusTimelineData ...
type CloudResourceStatusTimelineData struct {
	AvailabilityStatus string `json:"availabilitystatus"`
	OccurrenceTimeUTC  string `json:"occurrencetimeutc"`
	Summary            string `json:"summary"`
}

//CloudResourceStatusTimelineDataList : CloudResourceStatusTimelineData
type CloudResourceStatusTimelineDataList struct {
	CloudResourceStatusTimelineDataList []CloudResourceStatusTimelineData `json:"cloudResourceStatusTimelineDataList"`
}

//CloudResourceStatusTimelineListDataType ...
var CloudResourceStatusTimelineListDataType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "CloudResourceStatusTimelineListDataType",
		Fields: graphql.Fields{
			"AvailabilityStatus": &graphql.Field{
				Type:        graphql.String,
				Description: "Availability Status of the resource",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudResourceStatusTimelineData); ok {
						return CurData.AvailabilityStatus, nil
					}
					return nil, nil
				},
			},
			"Summary": &graphql.Field{
				Type:        graphql.String,
				Description: "Summary for the resource",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudResourceStatusTimelineData); ok {
						return CurData.Summary, nil
					}
					return nil, nil
				},
			},
			"OccurrenceTimeUTC": &graphql.Field{
				Type:        graphql.String,
				Description: "OccurrenceTimeUTC of the resource",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(CloudResourceStatusTimelineData); ok {
						return CurData.OccurrenceTimeUTC, nil
					}
					return nil, nil
				},
			},
		},
	},
)

//CloudResourceStatusTimelineDataListType : CRSTLDataList GraphQL Schema
var CloudResourceStatusTimelineDataListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CloudResourceStatusTimelineDataList",
	Fields: graphql.Fields{
		"cloudResourceStatusTimelineDataList": &graphql.Field{
			Type:        graphql.NewList(CloudResourceStatusTimelineListDataType),
			Description: "CloudResourceStatusTimelineDataList",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CloudResourceStatusTimelineDataList); ok {
					return CurData.CloudResourceStatusTimelineDataList, nil
				}
				return nil, nil
			},
		},
	},
})

//InterfaceToString : Function to Convert Interface To String
func InterfaceToString(Itr interface{}) (sReturn string) {
	if Itr != nil {
		sReturn = Itr.(string)
	}
	return sReturn
}

//AzureAccountData ...
type AzureAccountData struct {
	ClientID           string `json:"clientid"`
	EmailID            string `json:"emailid"`
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Status             int32  `json:"status"`
	TransactionID      string `json:"transactionid"`
	UpdateddatetimeUTC string `json:"updateddatetimeutc"`
}

//MapAccountsToSites ...
type MapAccountsToSites struct {
	AccountID string
	SiteName  string
	SiteID    int64
}

//ResourceAlert ...
type ResourceAlert struct {
	PartnerID    int64  `json:"PartnerID"`
	ResourceID   string `json:"ResourceID"`
	ResourceName string `json:"ResourceName"`
	SiteID       int64  `json:"SiteID"`
	TicketID     int64  `json:"TicketId"`
	TicketStatus string `json:"TktStatus"`
}

//ResourceAlertList : list of partnerwise alert
type ResourceAlertList struct {
	Alerts []ResourceAlert `json:"outdata"`
}

//CloudResourceConfigListDataType ...
// var CloudResourceConfigListDataType = graphql.NewObject(graphql.ObjectConfig{
// 	Name: "CloudResourceConfigDataType",
// 	Fields: graphql.Fields{
// 		"ResourceConfig": &graphql.Field{
// 			Type:        graphql.NewList(CloudResourceConfigDataType),
// 			Description: "ResourceConfig",
// 			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 				if CurData, ok := p.Source.(AzureConfigInfoList); ok {
// 					return CurData.ResourceConfig, nil
// 				}
// 				return nil, nil
// 			},
// 		},
// 		"ErrorMessage": &graphql.Field{
// 			Type:        graphql.String,
// 			Description: "ErrorMessage",
// 			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 				if CurData, ok := p.Source.(AzureConfigInfoList); ok {
// 					return CurData.ErrorMessage, nil
// 				}
// 				return nil, nil
// 			},
// 		},
// 	},
// })

//CloudResourceConfigPropDataType data type for cloud hierarchy
var CloudResourceConfigPropDataType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "CloudResourceConfigPropDataType",
		Fields: graphql.Fields{
			"Operatingsystem": &graphql.Field{
				Type:        graphql.String,
				Description: "Operatingsystem",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(ConfigProperties); ok {
						return CurData.Operatingsystem, nil
					}
					return nil, nil
				},
			},
			"Size": &graphql.Field{
				Type:        graphql.String,
				Description: "Size",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(ConfigProperties); ok {
						return CurData.Size, nil
					}
					return nil, nil
				},
			},
			"Tags": &graphql.Field{
				Type:        graphql.String,
				Description: "Tags",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(ConfigProperties); ok {
						return CurData.Tags, nil
					}
					return nil, nil
				},
			},
		},
	},
)

//CloudResourceConfigHierarchies data type for cloud hierarchy
var CloudResourceConfigHierarchies = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "CloudResourceConfigHierarchies",
		Fields: graphql.Fields{
			"Title": &graphql.Field{
				Type:        graphql.String,
				Description: "Title",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(ConfigHierarchies); ok {
						return CurData.Title, nil
					}
					return nil, nil
				},
			},
			"ID": &graphql.Field{
				Type:        graphql.String,
				Description: "ID",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(ConfigHierarchies); ok {
						return CurData.ID, nil
					}
					return nil, nil
				},
			},
			"Level": &graphql.Field{
				Type:        graphql.String,
				Description: "Level",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(ConfigHierarchies); ok {
						return CurData.Level, nil
					}
					return nil, nil
				},
			},
			"Name": &graphql.Field{
				Type:        graphql.String,
				Description: "Name",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(ConfigHierarchies); ok {
						return CurData.Name, nil
					}
					return nil, nil
				},
			},
		},
	},
)

//CloudResourceConfigDataType defines the graph ql data type for CloudResourceConfig
var CloudResourceConfigDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CloudResourceConfigDataType",
	Fields: graphql.Fields{
		"ID": &graphql.Field{
			Type:        graphql.String,
			Description: "id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AzureConfigInfo); ok {
					return CurData.ID, nil
				}
				return nil, nil
			},
		},
		"Name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AzureConfigInfo); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},
		"Service": &graphql.Field{
			Type:        graphql.String,
			Description: "Service",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AzureConfigInfo); ok {
					return CurData.Service, nil
				}
				return nil, nil
			},
		},
		"Location": &graphql.Field{
			Type:        graphql.String,
			Description: "Location",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AzureConfigInfo); ok {
					return CurData.Location, nil
				}
				return nil, nil
			},
		},
		"Properties": &graphql.Field{
			Type:        CloudResourceConfigPropDataType,
			Description: "Properties",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AzureConfigInfo); ok {
					return CurData.Properties, nil
				}
				return nil, nil
			},
		},
		"Hierarchies": &graphql.Field{
			Type:        graphql.NewList(CloudResourceConfigHierarchies),
			Description: "Hierarchies",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AzureConfigInfo); ok {
					return CurData.Hierarchies, nil
				}
				return nil, nil
			},
		},
	},
})

//AzureConfigInfo : configuration of azure resources
type AzureConfigInfo struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Service     string              `json:"service"`
	Location    string              `json:"location"`
	Properties  ConfigProperties    `json:"properties"`
	Hierarchies []ConfigHierarchies `json:"hierarchies"`
}

//ConfigProperties : ConfigProperties
type ConfigProperties struct {
	Operatingsystem string `json:"operatingsystem"`
	Size            string `json:"size"`
	Tags            string `json:"tags"`
}

//ConfigHierarchies : ConfigHierarchies
type ConfigHierarchies struct {
	Title string `json:"title"`
	Name  string `json:"name"`
	ID    string `json:"id"`
	Level int    `json:"level"`
}

//ClientListDataType defines the graph ql data type for ClientList
var ClientListDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ClientListDataType",
	Fields: graphql.Fields{
		"ClientID": &graphql.Field{
			Type:        graphql.String,
			Description: "ClientID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ClientList); ok {
					return CurData.ClientID, nil
				}
				return nil, nil
			},
		},
		"ClientName": &graphql.Field{
			Type:        graphql.String,
			Description: "ClientName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(ClientList); ok {
					return CurData.ClientName, nil
				}
				return nil, nil
			},
		},
	},
})

//CloudClientListDataType ...
var CloudClientListDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CloudClientListDataType",
	Fields: graphql.Fields{
		"Clients": &graphql.Field{
			Type:        graphql.NewList(ClientListDataType),
			Description: "Clients",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CloudClientList); ok {
					return CurData.Clients, nil
				}
				return nil, nil
			},
		},
		"ErrorMessage": &graphql.Field{
			Type:        graphql.String,
			Description: "ErrorMessage",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(CloudClientList); ok {
					return CurData.ErrorMessage, nil
				}
				return nil, nil
			},
		},
	},
})

//ClientList : ClientList
type ClientList struct {
	ClientID   string `json:"clientid"`
	ClientName string
}

//CloudClientList : CloudClientList
type CloudClientList struct {
	Clients      []ClientList `json:"ClientList"`
	ErrorMessage []string     `json:"errorMessage"`
}

//OnboardingSummaryData ...
type OnboardingSummaryData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	TenantID    string `json:"tenantid"`
	Username    string `json:"username"`
	State       string `json:"state"`
	IsMonitored bool   `json:"ismonitored"`
}

//OnboardingSummaryDataList : OnboardingSummaryDataList
type OnboardingSummaryDataList struct {
	OnboardingSummaryDataList []OnboardingSummaryData `json:"onboardingSummaryDataList"`
}

//OnboardingSummaryListDataType : OnboardingSummaryListDataType
var OnboardingSummaryListDataType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "OnboardingSummaryListDataType",
		Fields: graphql.Fields{
			"ID": &graphql.Field{
				Type:        graphql.String,
				Description: "Subscription ID of Account",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(OnboardingSummaryData); ok {
						return CurData.ID, nil
					}
					return nil, nil
				},
			},
			"Name": &graphql.Field{
				Type:        graphql.String,
				Description: "Name of subscription",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(OnboardingSummaryData); ok {
						return CurData.Name, nil
					}
					return nil, nil
				},
			},
			"TenantID": &graphql.Field{
				Type:        graphql.String,
				Description: "TenantID being monitored",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(OnboardingSummaryData); ok {
						return CurData.TenantID, nil
					}
					return nil, nil
				},
			},
			"Username": &graphql.Field{
				Type:        graphql.String,
				Description: "Username for TenantID",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(OnboardingSummaryData); ok {
						return CurData.Username, nil
					}
					return nil, nil
				},
			},
			"IsMonitored": &graphql.Field{
				Type:        graphql.String,
				Description: "If a subscription is being monitored or not",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(OnboardingSummaryData); ok {
						return CurData.IsMonitored, nil
					}
					return nil, nil
				},
			},
			"State": &graphql.Field{
				Type:        graphql.String,
				Description: "Subscription state",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if CurData, ok := p.Source.(OnboardingSummaryData); ok {
						return CurData.State, nil
					}
					return nil, nil
				},
			},
		},
	},
)

//OnboardingSummaryDataListType : OnboardingSummaryDataListType
var OnboardingSummaryDataListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "OnboardingSummaryDataList",
	Fields: graphql.Fields{
		"onboardingSummaryDataList": &graphql.Field{
			Type:        graphql.NewList(OnboardingSummaryListDataType),
			Description: "OnboardingSummaryDataList",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(OnboardingSummaryDataList); ok {
					return CurData.OnboardingSummaryDataList, nil
				}
				return nil, nil
			},
		},
	},
})
