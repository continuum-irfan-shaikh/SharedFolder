package models

import "github.com/gocql/gocql"

type ScriptDefinition struct {
	ID                       gocql.UUID  `json:"id"`
	PartnerID                string      `json:"partnerId"`
	Category                 []string    `json:"category"`
	Description              string      `json:"description"`
	Engine                   string      `json:"engine"`
	EngineMaxVersion         int         `json:"engineMaxVersion"`
	ExpectedExecutionTimeSec int         `json:"expectedExecutionTimeSec"`
	FailureMessage           string      `json:"failureMessage"`
	SuccessMessage           string      `json:"successMessage"`
	Internal                 bool        `json:"internal"`
	Name                     string      `json:"name"`
	Tags                     []string    `json:"tags"`
	Sequence                 bool        `json:"sequence"`
	Content                  string      `json:"content"`
	JSONSchema               interface{} `json:"jsonSchema"`
	UISchema                 interface{} `json:"uiSchema"`
	IsHDScript               bool        `json:"isHDScript"`
	Disabled                 bool        `json:"disabled,omitempty"`
}

func NewScriptDefinition() ScriptDefinition {
	return ScriptDefinition{
		ID:                       gocql.TimeUUID(),
		PartnerID:                "00000000-0000-0000-0000-000000000000",
		Engine:                   "powershell",
		EngineMaxVersion:         5,
		ExpectedExecutionTimeSec: 300,
		FailureMessage:           "Executed with errors",
		Internal:                 false,
		Sequence:                 true,
		Tags:                     []string{"Windows 7", "Windows 10"},
		IsHDScript:               true, // visible only for NOC users by default
	}
}
