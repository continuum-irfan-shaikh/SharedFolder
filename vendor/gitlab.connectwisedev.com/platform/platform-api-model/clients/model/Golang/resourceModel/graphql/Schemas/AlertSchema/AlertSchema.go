package AlertSchema

import (
	"github.com/graphql-go/graphql"
)

//PartnerAlerts : PartnerAlerts Structure
type PartnerAlerts struct {
	Client          	string	`json:"client"`
	Location       		string 	`json:"location"`
	Groupname           	string 	`json:"groupname"`
	TaskID       		int64 	`json:"taskid"`
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
}

//PartnerAlertsType : PartnerAlerts GraphQL Schema
var PartnerAlertsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "partnerAlerts",
	Fields: graphql.Fields{
		"client": &graphql.Field{
			Type:        graphql.String,
			Description: "client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlerts); ok {
					return CurData.Client, nil
				}
				return nil, nil
			},
		},

		"location": &graphql.Field{
			Type:        graphql.String,
			Description: "location",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlerts); ok {
					return CurData.Location, nil
				}
				return nil, nil
			},
		},

		"groupname": &graphql.Field{
			Type:        graphql.String,
			Description: "groupname",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlerts); ok {
					return CurData.Groupname, nil
				}
				return nil, nil
			},
		},

		"taskid": &graphql.Field{
			Type:        graphql.String,
			Description: "taskid",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlerts); ok {
					return CurData.TaskID, nil
				}
				return nil, nil
			},
		},

		"jobDescription": &graphql.Field{
			Type:        graphql.String,
			Description: "jobDescription",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlerts); ok {
					return CurData.JobDescription, nil
				}
				return nil, nil
			},
		},

		"jobName": &graphql.Field{
			Type:        graphql.String,
			Description: "jobName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlerts); ok {
					return CurData.JobName, nil
				}
				return nil, nil
			},
		},

		"resource": &graphql.Field{
			Type:        graphql.String,
			Description: "resource",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlerts); ok {
					return CurData.Resource, nil
				}
				return nil, nil
			},
		},

		"jobDateTime": &graphql.Field{
			Type:        graphql.String,
			Description: "jobDateTime",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlerts); ok {
					return CurData.JobDateTime, nil
				}
				return nil, nil
			},
		},

		"categoryname": &graphql.Field{
			Type:        graphql.String,
			Description: "categoryname",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlerts); ok {
					return CurData.Categoryname, nil
				}
				return nil, nil
			},
		},

		"statusName": &graphql.Field{
			Type:        graphql.String,
			Description: "statusName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlerts); ok {
					return CurData.StatusName, nil
				}
				return nil, nil
			},
		},

		"mnid": &graphql.Field{
			Type:        graphql.String,
			Description: "mnid",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlerts); ok {
					return CurData.MnID, nil
				}
				return nil, nil
			},
		},

		"totalMins": &graphql.Field{
			Type:        graphql.String,
			Description: "totalMins",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlerts); ok {
					return CurData.TotalMins, nil
				}
				return nil, nil
			},
		},

		"jobIdDisplay": &graphql.Field{
			Type:        graphql.String,
			Description: "jobIdDisplay",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlerts); ok {
					return CurData.JobIDDisplay, nil
				}
				return nil, nil
			},
		},

		"statusUpdateOn": &graphql.Field{
			Type:        graphql.String,
			Description: "statusUpdateOn",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlerts); ok {
					return CurData.StatusUpdateOn, nil
				}
				return nil, nil
			},
		},

		"regid": &graphql.Field{
			Type:        graphql.String,
			Description: "regid",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlerts); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},

		"thresholdCounter": &graphql.Field{
			Type:        graphql.String,
			Description: "thresholdCounter",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlerts); ok {
					return CurData.ThresholdCounter, nil
				}
				return nil, nil
			},
		},
	},
})

//PartnerAlertsList : PartnerAlertsList Structure
type PartnerAlertsList struct {
	RegID          					string			`json:"regId"`
	SiteID       					int64 			`json:"siteId"`
	PAlerts          				[]PartnerAlerts 	`json:"partnerDefinedAlertData"`
	TotalPartnerDefinedAlertCount			int64 			`json:"totalPartnerDefinedAlertCount"`
}

//PartnerAlertsListType : PartnerAlertsList GraphQL Schema
var PartnerAlertsListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "partnerDefinedAlertList",
	Fields: graphql.Fields{
		"regId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlertsList); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},

		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlertsList); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"partnerDefinedAlertData": &graphql.Field{
			Type:        graphql.NewList(PartnerAlertsType),
			Description: "partner defined list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlertsList); ok {
					return CurData.PAlerts, nil
				}
				return nil, nil
			},
		},

		"totalPartnerDefinedAlertCount": &graphql.Field{
			Type:        graphql.String,
			Description: "total partner defined alert count",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(PartnerAlertsList); ok {
					return CurData.TotalPartnerDefinedAlertCount, nil
				}
				return nil, nil
			},
		},
	},
})

//NOCAlerts : NOCAlerts Structure
type NOCAlerts struct {
	Client          	string	`json:"client"`
	Location       		string 	`json:"location"`
	Groupname           	string 	`json:"groupname"`
	TaskID       		int64 	`json:"taskId"`
	JobID       		int64 	`json:"jobid"`
	JobName 		string  `json:"jobName"`
	JobDescription      	string 	`json:"jobDescription"`
	Resource       		string 	`json:"resource"`
	JobDateTime         	string 	`json:"jobDateTime"`
	Categoryname 		string  `json:"categoryname"`
	StatusName       	string 	`json:"statusName"`
	MnID         		int64 	`json:"mnid"`
	TotalMins 		int64  	`json:"totalMins"`
	JobIDDisplay       	string 	`json:"jobIdDisplay"`
	StatusUpdateOn     	string 	`json:"statusUpdateOn"`
	DefinedBy 		string 	`json:"definedBy"`
	AlertGroup 		string 	`json:"alertGroup"`
	ThresholdCounter	string 	`json:"thresholdCounter"`
}

//NOCAlertsType : NOCAlerts GraphQL Schema
var NOCAlertsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "nocAlerts",
	Fields: graphql.Fields{
		"client": &graphql.Field{
			Type:        graphql.String,
			Description: "client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlerts); ok {
					return CurData.Client, nil
				}
				return nil, nil
			},
		},

		"location": &graphql.Field{
			Type:        graphql.String,
			Description: "location",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlerts); ok {
					return CurData.Location, nil
				}
				return nil, nil
			},
		},

		"groupname": &graphql.Field{
			Type:        graphql.String,
			Description: "groupname",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlerts); ok {
					return CurData.Groupname, nil
				}
				return nil, nil
			},
		},

		"taskId": &graphql.Field{
			Type:        graphql.String,
			Description: "taskId",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlerts); ok {
					return CurData.TaskID, nil
				}
				return nil, nil
			},
		},

		"jobid": &graphql.Field{
			Type:        graphql.String,
			Description: "jobid",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlerts); ok {
					return CurData.JobID, nil
				}
				return nil, nil
			},
		},

		"jobName": &graphql.Field{
			Type:        graphql.String,
			Description: "jobName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlerts); ok {
					return CurData.JobName, nil
				}
				return nil, nil
			},
		},

		"jobDescription": &graphql.Field{
			Type:        graphql.String,
			Description: "jobDescription",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlerts); ok {
					return CurData.JobDescription, nil
				}
				return nil, nil
			},
		},

		"resource": &graphql.Field{
			Type:        graphql.String,
			Description: "resource",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlerts); ok {
					return CurData.Resource, nil
				}
				return nil, nil
			},
		},

		"jobDateTime": &graphql.Field{
			Type:        graphql.String,
			Description: "jobDateTime",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlerts); ok {
					return CurData.JobDateTime, nil
				}
				return nil, nil
			},
		},

		"categoryname": &graphql.Field{
			Type:        graphql.String,
			Description: "categoryname",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlerts); ok {
					return CurData.Categoryname, nil
				}
				return nil, nil
			},
		},

		"statusName": &graphql.Field{
			Type:        graphql.String,
			Description: "statusName",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlerts); ok {
					return CurData.StatusName, nil
				}
				return nil, nil
			},
		},

		"mnid": &graphql.Field{
			Type:        graphql.String,
			Description: "mnid",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlerts); ok {
					return CurData.MnID, nil
				}
				return nil, nil
			},
		},

		"totalMins": &graphql.Field{
			Type:        graphql.String,
			Description: "totalMins",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlerts); ok {
					return CurData.TotalMins, nil
				}
				return nil, nil
			},
		},

		"jobIdDisplay": &graphql.Field{
			Type:        graphql.String,
			Description: "jobIdDisplay",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlerts); ok {
					return CurData.JobIDDisplay, nil
				}
				return nil, nil
			},
		},

		"statusUpdateOn": &graphql.Field{
			Type:        graphql.String,
			Description: "statusUpdateOn",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlerts); ok {
					return CurData.StatusUpdateOn, nil
				}
				return nil, nil
			},
		},

		"definedBy": &graphql.Field{
			Type:        graphql.String,
			Description: "definedBy",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlerts); ok {
					return CurData.DefinedBy, nil
				}
				return nil, nil
			},
		},

		"alertGroup": &graphql.Field{
			Type:        graphql.String,
			Description: "alertGroup",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlerts); ok {
					return CurData.AlertGroup, nil
				}
				return nil, nil
			},
		},

		"thresholdCounter": &graphql.Field{
			Type:        graphql.String,
			Description: "thresholdCounter",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlerts); ok {
					return CurData.ThresholdCounter, nil
				}
				return nil, nil
			},
		},
	},
})

//NOCAlertsList : NOCAlertsList Structure
type NOCAlertsList struct {
	RegID          			string			`json:"regId"`
	SiteID       			int64 			`json:"siteId"`
	NAlerts          		[]NOCAlerts 		`json:"nocDefinedAlertData"`
	TotalNocDefinedAlertCount	int64 			`json:"totalNocDefinedAlertCount"`
}

//NOCAlertsListType : NOCAlertsList GraphQL Schema
var NOCAlertsListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "nocDefinedAlertList",
	Fields: graphql.Fields{
		"regId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for a specific endpoint",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlertsList); ok {
					return CurData.RegID, nil
				}
				return nil, nil
			},
		},

		"siteId": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier representing a specific client",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlertsList); ok {
					return CurData.SiteID, nil
				}
				return nil, nil
			},
		},

		"nocDefinedAlertData": &graphql.Field{
			Type:        graphql.NewList(NOCAlertsType),
			Description: "noc defined list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlertsList); ok {
					return CurData.NAlerts, nil
				}
				return nil, nil
			},
		},

		"totalNocDefinedAlertCount": &graphql.Field{
			Type:        graphql.String,
			Description: "total Noc defined alert count",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(NOCAlertsList); ok {
					return CurData.TotalNocDefinedAlertCount, nil
				}
				return nil, nil
			},
		},
	},
})

//AlertsList : AlertsList Structure
type AlertsList struct {
	PAlertList	[]PartnerAlertsList	`json:"partnerDefinedAlertList"`
	NAlertList  	[]NOCAlertsList 	`json:"nocDefinedAlertList"`
	TotalCount  	int64 			`json:"totalCount"`
}


//AlertsListType : AlertsList GraphQL Schema
var AlertsListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "alertsList",
	Fields: graphql.Fields{
		"partnerDefinedAlertList": &graphql.Field{
			Type:        graphql.NewList(PartnerAlertsListType),
			Description: "partner defined alert list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertsList); ok {
					return CurData.PAlertList, nil
				}
				return nil, nil
			},
		},

		"nocDefinedAlertList": &graphql.Field{
			Type:        graphql.NewList(NOCAlertsListType),
			Description: "noc defined alert list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertsList); ok {
					return CurData.NAlertList, nil
				}
				return nil, nil
			},
		},

		"totalCount": &graphql.Field{
			Type:        graphql.String,
			Description: "totalCount",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AlertsList); ok {
					return CurData.TotalCount, nil
				}
				return nil, nil
			},
		},
	},
})
