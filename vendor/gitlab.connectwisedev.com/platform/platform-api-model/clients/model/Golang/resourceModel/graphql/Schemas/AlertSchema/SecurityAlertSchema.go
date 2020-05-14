package AlertSchema

import (
	"github.com/graphql-go/graphql"
)

//SecurityAlerts : SecurityAlerts Structure
type SecurityAlerts struct {
	Client          	string	`json:"client"`
	Location       		string 	`json:"location"`
	Groupname           	string 	`json:"groupname"`
	TaskID       		int64 	`json:"taskid"`
	JobID       		int64 	`json:"jobid"`
	JobDescription      	string 	`json:"jobDescription"`
	JobName 		string  `json:"jobName"`
	Resource       		string 	`json:"resource"`
	JobDateTime         	string 	`json:"jobDateTime"`
	Categoryname 		string  `json:"categoryname"`
	StatusName       	string 	`json:"statusName"`
	MnID         		int64 	`json:"mnid"`
	TotalMins 		int64  	`json:"totalMins"`
	JobIDDisplay       	string 	`json:"jobIdDisplay"`
	StatusUpdateOn      	string 	`json:"statusUpdateOn"`
	RegID 			int64  	`json:"regid"`
	ThresholdCounter	string 	`json:"thresholdCounter"`
	ResType	string 	`json:"resType"`	
	RegType	string 	`json:"regType"`
	AlertGroup	string 	`json:"AlertGroup"`
}

//SecurityAlertsType : SecurityAlerts GraphQL Schema
var SecurityAlertsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SecurityAlerts",
	Fields: graphql.Fields{
		"client": &graphql.Field{
			Type:        graphql.String,
			Description: "client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.Client, nil
				}
				return nil, nil
			},
		},

		"location": &graphql.Field{
			Type:        graphql.String,
			Description: "location",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.Location, nil
				}
				return nil, nil
			},
		},

		"groupname": &graphql.Field{
			Type:        graphql.String,
			Description: "groupname",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.Groupname, nil
				}
				return nil, nil
			},
		},

		"taskid": &graphql.Field{
			Type:        graphql.String,
			Description: "taskid",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.TaskID, nil
				}
				return nil, nil
			},
		},

		"jobid": &graphql.Field{
			Type:        graphql.String,
			Description: "jobid",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.JobID, nil
				}
				return nil, nil
			},
		},

		"jobDescription": &graphql.Field{
			Type:        graphql.String,
			Description: "jobDescription",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.JobDescription, nil
				}
				return nil, nil
			},
		},

		"jobName": &graphql.Field{
			Type:        graphql.String,
			Description: "jobName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.JobName, nil
				}
				return nil, nil
			},
		},

		"resource": &graphql.Field{
			Type:        graphql.String,
			Description: "resource",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.Resource, nil
				}
				return nil, nil
			},
		},

		"jobDateTime": &graphql.Field{
			Type:        graphql.String,
			Description: "jobDateTime",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.JobDateTime, nil
				}
				return nil, nil
			},
		},

		"categoryname": &graphql.Field{
			Type:        graphql.String,
			Description: "categoryname",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.Categoryname, nil
				}
				return nil, nil
			},
		},

		"statusName": &graphql.Field{
			Type:        graphql.String,
			Description: "statusName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.StatusName, nil
				}
				return nil, nil
			},
		},

		"mnid": &graphql.Field{
			Type:        graphql.String,
			Description: "mnid",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.MnID, nil
				}
				return nil, nil
			},
		},

		"totalMins": &graphql.Field{
			Type:        graphql.String,
			Description: "totalMins",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.TotalMins, nil
				}
				return nil, nil
			},
		},

		"jobIdDisplay": &graphql.Field{
			Type:        graphql.String,
			Description: "jobIdDisplay",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.JobIDDisplay, nil
				}
				return nil, nil
			},
		},

		"statusUpdateOn": &graphql.Field{
			Type:        graphql.String,
			Description: "statusUpdateOn",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.StatusUpdateOn, nil
				}
				return nil, nil
			},
		},

		"regid": &graphql.Field{
			Type:        graphql.String,
			Description: "regid",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},

		"thresholdCounter": &graphql.Field{
			Type:        graphql.String,
			Description: "thresholdCounter",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.ThresholdCounter, nil
				}
				return nil, nil
			},
		},
		
		"resType": &graphql.Field{
			Type:        graphql.String,
			Description: "resType",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.ResType, nil
				}
				return nil, nil
			},
		},

		"regType": &graphql.Field{
			Type:        graphql.String,
			Description: "regType",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.RegType, nil
				}
				return nil, nil
			},
		},

		"alertGroup": &graphql.Field{
			Type:        graphql.String,
			Description: "alertGroup",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlerts); ok {
					return CurData.AlertGroup, nil
				}
				return nil, nil
			},
		},
	},
})

//SecurityAlertsList : SecurityAlertsList Structure
type SecurityAlertsList struct {
	RegID          					string			`json:"regId"`
	SiteID       					int64 			`json:"siteId"`
	TotalSecurityDefinedAlertCount			int64 			`json:"totalSecurityDefinedAlertCount"`
	SAlerts          				[]SecurityAlerts 	`json:"securityDefinedAlertData"`
	
}

//SecurityAlertsListType : SecurityAlertsList GraphQL Schema
var SecurityAlertsListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "securityDefinedAlertList",
	Fields: graphql.Fields{
		"regId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlertsList); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},

		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlertsList); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"totalSecurityDefinedAlertCount": &graphql.Field{
			Type:        graphql.String,
			Description: "total partner defined alert count",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlertsList); ok {
					return CurData.TotalSecurityDefinedAlertCount, nil
				}
				return nil, nil
			},
		},

		"securityDefinedAlertData": &graphql.Field{
			Type:        graphql.NewList(SecurityAlertsType),
			Description: "partner defined list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlertsList); ok {
					return CurData.SAlerts, nil
				}
				return nil, nil
			},
		},

	},
})

//SecurityAlertsRespList : Security Alerts List Structure
type SecurityAlertsRespList struct {
	SAlertList	[]SecurityAlertsList	`json:"securityDefinedAlertList"`
}


//SecurityAlertsRespListType : security Alerts response List GraphQL Schema
var SecurityAlertsRespListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SecurityAlertsRespList",
	Fields: graphql.Fields{
		"securityDefinedAlertList": &graphql.Field{
			Type:        graphql.NewList(SecurityAlertsListType),
			Description: "security defined alert response list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(SecurityAlertsRespList); ok {
					return CurData.SAlertList, nil
				}
				return nil, nil
			},
		},
	},
})
