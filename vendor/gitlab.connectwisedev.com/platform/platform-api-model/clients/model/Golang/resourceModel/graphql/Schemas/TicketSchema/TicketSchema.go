package TicketSchema

import (
	"errors"
	"strings"
	"reflect"
	"github.com/graphql-go/graphql"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Relay"
)

//TicketData : TicketData Structure
type TicketData struct {
	SiteID          	string	`json:"siteId"`
	PartnerID       	string	`json:"partnerId"`
	RegID           	string	`json:"regId"`
	ResType       		string	`json:"resType"`
	TaskDesc         	string	`json:"taskDesc"`
	TaskType 		string	`json:"taskType"`
	TaskID       		int64 	`json:"taskID"`
	DisplayTaskID         	string	`json:"displayTaskID"`
	Location 		string 	`json:"location"`
	RegType       		string	`json:"regType"`
	MemberCode         	string	`json:"memberCode"`
	SiteName 		string  `json:"siteName"`
	MachineName       	string	`json:"machineName"`
	ResourceName         	string	`json:"resourceName"`
	Client 			string  `json:"client"`
	TaskSubject       	string	`json:"taskSubject"`
	AssignTo         	string	`json:"assignTo"`
	CategoryName 		string  `json:"categoryName"`
	GroupName       	string	`json:"groupName"`
	Resource         	string	`json:"resource"`
	ResFriendlyName 	string  `json:"resFriendlyName"`
	Status       		int64 	`json:"status"`
	TaskDateTime         	string	`json:"taskDateTime"`
	TaskTypeItsupport247 	string  `json:"taskTypeItsupport247"`
	TimeZone       		string	`json:"timeZone"`
	FName         		string	`json:"fName"`
	RefTaskID 		string  `json:"refTaskId"`
	RefTicketID       	string	`json:"refTicketId"`
	Priority         	int64	`json:"priority"`
	TaskDescription 	string  `json:"taskDescription"`
	AssignTaskTo       	int64	`json:"assignTaskTo"`
	PriorityName         	string	`json:"priorityName"`
	StatusName 		string  `json:"statusName"`
	UserName       		string 	`json:"userName"`
	Generatedby         	string 	`json:"generatedby"`
	AutoTask 		string  `json:"autoTask"`
	PSAId       		string 	`json:"psaId"`
	EscCategory         	string 	`json:"escCategory"`
	ConditionFamily 	string  `json:"conditionFamily"`
	TaskExecutionDate       string 	`json:"taskExecutionDate"`
	ConditionID         	int64 	`json:"conditionId"`
	Duration 		int64  	`json:"duration"`
	Timediff       		float64 `json:"timediff"`
	StatusUpdatedOn         string 	`json:"statusUpdatedOn"`
	MnID 			int64  	`json:"mnId"`
}

//TicketDataType : TicketData GraphQL Schema
var TicketDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ticketdata",
	Fields: graphql.Fields{
		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique site identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"partnerId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique partner identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.PartnerID, nil
				}
				return nil, nil
			},
		},

		"regId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},

		"resType": &graphql.Field{
			Type:        graphql.String,
			Description: "resType",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.ResType, nil
				}
				return nil, nil
			},
		},

		"taskDesc": &graphql.Field{
			Type:        graphql.String,
			Description: "taskDesc",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.TaskDesc, nil
				}
				return nil, nil
			},
		},

		"taskType": &graphql.Field{
			Type:        graphql.String,
			Description: "TaskType",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.TaskType, nil
				}
				return nil, nil
			},
		},

		"taskID": &graphql.Field{
			Type:        graphql.String,
			Description: "TaskID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.TaskID, nil
				}
				return nil, nil
			},
		},

		"displayTaskID": &graphql.Field{
			Type:        graphql.String,
			Description: "The task id to be displayed",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.DisplayTaskID, nil
				}
				return nil, nil
			},
		},

		"location": &graphql.Field{
			Type:        graphql.String,
			Description: "The location",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.Location, nil
				}
				return nil, nil
			},
		},

		"regType": &graphql.Field{
			Type:        graphql.String,
			Description: "Type of the Reg",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.RegType, nil
				}
				return nil, nil
			},
		},

		"memberCode": &graphql.Field{
			Type:        graphql.String,
			Description: "The member code",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.MemberCode, nil
				}
				return nil, nil
			},
		},

		"siteName": &graphql.Field{
			Type:        graphql.String,
			Description: "Name for the site",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.SiteName, nil
				}
				return nil, nil
			},
		},

		"machineName": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the Machine",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.MachineName, nil
				}
				return nil, nil
			},
		},

		"resourceName": &graphql.Field{
			Type:        graphql.String,
			Description: "This is the resource name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.ResourceName, nil
				}
				return nil, nil
			},
		},

		"client": &graphql.Field{
			Type:        graphql.String,
			Description: "This is the client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.Client, nil
				}
				return nil, nil
			},
		},

		"taskSubject": &graphql.Field{
			Type:        graphql.String,
			Description: "The task subject",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.TaskSubject, nil
				}
				return nil, nil
			},
		},

		"assignTo": &graphql.Field{
			Type:        graphql.String,
			Description: "The assign to",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.AssignTo, nil
				}
				return nil, nil
			},
		},

		"categoryName": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the category",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.CategoryName, nil
				}
				return nil, nil
			},
		},

		"groupName": &graphql.Field{
			Type:        graphql.String,
			Description: "This is the group",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.GroupName, nil
				}
				return nil, nil
			},
		},

		"resource": &graphql.Field{
			Type:        graphql.String,
			Description: "The resource",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.Resource, nil
				}
				return nil, nil
			},
		},

		"resFriendlyName": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the resource friendly",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.ResFriendlyName, nil
				}
				return nil, nil
			},
		},

		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "The status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.Status, nil
				}
				return nil, nil
			},
		},

		"taskDateTime": &graphql.Field{
			Type:        graphql.String,
			Description: "The task date time",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.TaskDateTime, nil
				}
				return nil, nil
			},
		},

		"taskTypeItsupport247": &graphql.Field{
			Type:        graphql.String,
			Description: "The task type itsupport247",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.TaskTypeItsupport247, nil
				}
				return nil, nil
			},
		},

		"timeZone": &graphql.Field{
			Type:        graphql.String,
			Description: "The time zone",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.TimeZone, nil
				}
				return nil, nil
			},
		},

		"fName": &graphql.Field{
			Type:        graphql.String,
			Description: "This is the request reg id",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.FName, nil
				}
				return nil, nil
			},
		},

		"refTaskId": &graphql.Field{
			Type:        graphql.String,
			Description: "The reference task identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.RefTaskID, nil
				}
				return nil, nil
			},
		},

		"refTicketId": &graphql.Field{
			Type:        graphql.String,
			Description: "The reference ticket identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.RefTicketID, nil
				}
				return nil, nil
			},
		},

		"priority": &graphql.Field{
			Type:        graphql.String,
			Description: "The priority",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.Priority, nil
				}
				return nil, nil
			},
		},

		"taskDescription": &graphql.Field{
			Type:        graphql.String,
			Description: "The task description",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.TaskDescription, nil
				}
				return nil, nil
			},
		},

		"assignTaskTo": &graphql.Field{
			Type:        graphql.String,
			Description: "The assign task to",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.AssignTaskTo, nil
				}
				return nil, nil
			},
		},

		"priorityName": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the priority",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.PriorityName, nil
				}
				return nil, nil
			},
		},

		"statusName": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the status",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.StatusName, nil
				}
				return nil, nil
			},
		},

		"userName": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the user",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.UserName, nil
				}
				return nil, nil
			},
		},

		"generatedby": &graphql.Field{
			Type:        graphql.String,
			Description: "The generatedby",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.Generatedby, nil
				}
				return nil, nil
			},
		},

		"autoTask": &graphql.Field{
			Type:        graphql.String,
			Description: "The automatic task",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.AutoTask, nil
				}
				return nil, nil
			},
		},

		"psaId": &graphql.Field{
			Type:        graphql.String,
			Description: "The psa identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.PSAId, nil
				}
				return nil, nil
			},
		},

		"escCategory": &graphql.Field{
			Type:        graphql.String,
			Description: "The esc category",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.EscCategory, nil
				}
				return nil, nil
			},
		},

		"conditionFamily": &graphql.Field{
			Type:        graphql.String,
			Description: "The condition family",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.ConditionFamily, nil
				}
				return nil, nil
			},
		},

		"taskExecutionDate": &graphql.Field{
			Type:        graphql.String,
			Description: "TaskExecutionDate",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.TaskExecutionDate, nil
				}
				return nil, nil
			},
		},

		"conditionId": &graphql.Field{
			Type:        graphql.String,
			Description: "The condition identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.ConditionID, nil
				}
				return nil, nil
			},
		},

		"duration": &graphql.Field{
			Type:        graphql.String,
			Description: "The duration",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.Duration, nil
				}
				return nil, nil
			},
		},

		"timediff": &graphql.Field{
			Type:        graphql.String,
			Description: "The timediff",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.Timediff, nil
				}
				return nil, nil
			},
		},

		"statusUpdatedOn": &graphql.Field{
			Type:        graphql.String,
			Description: "The status updated on",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.StatusUpdatedOn, nil
				}
				return nil, nil
			},
		},

		"mnId": &graphql.Field{
			Type:        graphql.String,
			Description: "The MnId",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(TicketData); ok {
					return CurData.MnID, nil
				}
				return nil, nil
			},
		},
	},
})

//TicketDataConnectionDefinition : TicketDataConnectionDefinition structure
var TicketDataConnectionDefinition = Relay.ConnectionDefinitions(Relay.ConnectionConfig{
	Name:     "TicketData",
	NodeType: TicketDataType,
})

//TicketListData : TicketListData Structure
type TicketListData struct {
	Tickets	[]TicketData `json:"ticketList"`
}

//TicketListDataType : TicketListData GraphQL Schema
var TicketListDataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TicketList",
	Fields: graphql.Fields{
		"ticketList": &graphql.Field{
			Type:        TicketDataConnectionDefinition.ConnectionType,
			Args:        Relay.ConnectionArgs,
			Description: "Tickets list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := Relay.NewConnectionArguments(p.Args)
				if CurData, ok := p.Source.(TicketListData); ok {
					var arraySliceRet []interface{}
					for ind := range CurData.Tickets {
						arraySliceRet = append(arraySliceRet, CurData.Tickets[ind])
					}

					if args.Filter != "" && args.Filter != Relay.NILQUERY{
						var err error
						val := reflect.Indirect(reflect.ValueOf(&TicketData{}))
						arraySliceRet, err = Relay.Filter(string(args.Filter), val, arraySliceRet)
						if err != nil {
							return nil, err
						}
					}
					if args.Sort != "" && args.Sort != Relay.NILQUERY {
						subQuery := strings.Split(string(args.Sort), ";")
						SiteIDASC := func(p1, p2 interface{}) bool {
							return p1.(TicketData).SiteID < p2.(TicketData).SiteID
						}
						SiteIDDESC := func(p1, p2 interface{}) bool {
							return p1.(TicketData).SiteID > p2.(TicketData).SiteID
						}

						RegIDASC := func(p1, p2 interface{}) bool {
							return p1.(TicketData).RegID < p2.(TicketData).RegID
						}
						RegIDDESC := func(p1, p2 interface{}) bool {
							return p1.(TicketData).RegID > p2.(TicketData).RegID
						}

						DurationASC := func(p1, p2 interface{}) bool {
							return p1.(TicketData).Duration < p2.(TicketData).Duration
						}
						DurationDESC := func(p1, p2 interface{}) bool {
							return p1.(TicketData).Duration > p2.(TicketData).Duration
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
							}else if strings.ToUpper(Column) == "REGID" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(RegIDASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(RegIDDESC).Sort(arraySliceRet)
								}
							}else if strings.ToUpper(Column) == "DURATION" {
								if strings.ToUpper(Key) == Relay.ORDASC {
									Relay.SortBy(DurationASC).Sort(arraySliceRet)
								}else if strings.ToUpper(Key) == Relay.ORDDESC {
									Relay.SortBy(DurationDESC).Sort(arraySliceRet)
								}
							}else{
								return nil, errors.New("TicketData Sort [" + Column + "] No such column exist!!!")
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
