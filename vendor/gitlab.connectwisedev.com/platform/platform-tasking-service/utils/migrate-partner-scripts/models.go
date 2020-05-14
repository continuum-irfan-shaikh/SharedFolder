package main

import (
	"encoding/xml"
	"time"

	"github.com/gocql/gocql"
)

const (
	bashTemplateID          = "970"
	cmdTemplateID           = "908"
	powershellTemplateID    = "848"
	bashTemplateIDint       = 970
	cmdTemplateIDint        = 908
	powershellTemplateIDint = 848

	scriptType          = "script"
	CassandraTimeFormat = `2006-01-02 15:04:05-0700`
	csvFileName         = "result.csv"
)

var (
	cmdOriginID, _        = gocql.ParseUUID("e3d2c26b-c5ba-49cf-a089-7637f6de949e")
	powershellOriginID, _ = gocql.ParseUUID("51a74346-e19b-11e7-9809-0800279505d9")
	bashOriginID, _       = gocql.ParseUUID("37f7f19f-40e8-11e9-a643-e0d55e1ce78a") //https://gitlab.connectwisedev.com/platform/rmm-scripts/pull/196
)

type ScriptMsSQL struct {
	Category   string
	ScriptID   int
	MemberID   int
	ScriptName string
	ScriptDesc string
	ScriptXML  string // ScriptData
	TemplateID int
	CreatedBy  string
	CreatedAt  time.Time //DCDtime
	UpdatedAt  time.Time //ModifiedOn
	UpdatedBy  string    //ModifiedBy
}

type PowershellXMLData struct {
	XmlName  xml.Name            `xml:"b"` //tag that is needed to parse XML, there is no such tag in 1.0 xml, so add it manualy
	Category string              `xml:"agentcategory"`
	Steps    []StepPowershellTag `xml:"step"`
}

type CmdXMLData struct {
	XmlName  xml.Name     `xml:"b"` //tag that is needed to parse XML, there is no such tag in 1.0 xml, so add it manualy
	Category string       `xml:"agentcategory"`
	Steps    []StepCmdTag `xml:"step"`
}

type StepPowershellTag struct {
	Body       string `xml:"data"`
	Parameters string `xml:"cmdpara"`
	Timeout    int    `xml:"timeoutval"`
}

type StepCmdTag struct {
	Body       string `xml:"doscommand"`
	Parameters string `xml:"parameters"`
	PauseSec   int    `xml:"pausesec"`
	TimeoutMin int    `xml:"execmin"`
}

type StepBashTag struct {
	Body    string `xml:"data"`
	Timeout int    `xml:"timeoutval"`
}

type BashXMLData struct {
	XmlName  xml.Name    `xml:"b"` //tag that is needed to parse XML, there is no such tag in 1.0 xml, so add it manualy
	Category string      `xml:"agentcategory"`
	Steps    StepBashTag `xml:"step"`
}

// TaskDefinitionDetails structure contains detailed data about particular user defined task definition
type TaskDefinition struct {
	ID             gocql.UUID `json:"id"         valid:"unsettableByUsers"`
	PartnerID      string     `json:"partnerId"  valid:"unsettableByUsers"`
	OriginID       gocql.UUID `json:"originId"   valid:"requiredForUsers"`
	Name           string     `json:"name"       valid:"requiredForUsers"`
	Type           string     `json:"type"       valid:"validType"`
	Categories     []string   `json:"categories" valid:"validCategories"`
	Deleted        bool       `json:"-"`
	Description    string     `json:"description"    valid:"optional"`
	CreatedAt      time.Time  `json:"createdAt"      valid:"unsettableByUsers"`
	CreatedBy      string     `json:"createdBy"      valid:"unsettableByUsers"`
	UpdatedAt      time.Time  `json:"updatedAt"      valid:"unsettableByUsers"`
	UpdatedBy      string     `json:"updatedBy"      valid:"unsettableByUsers"`
	UserParameters string     `json:"userParameters" valid:"json"`
	UISchema       string     `json:"UISchema"       valid:"unsettableByUsers"`
	JSONSchema     string     `json:"JSONSchema"     valid:"unsettableByUsers"`
}

type CustomScriptBody struct {
	ScriptBody  string `json:"body"`
	ExpExecTime int    `json:"expectedExecutionTimeSec"`
}
