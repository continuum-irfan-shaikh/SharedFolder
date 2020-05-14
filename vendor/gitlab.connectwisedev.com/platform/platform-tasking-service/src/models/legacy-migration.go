package models

type (
	// LegacyScriptInfo represents legacy script report and mapping info struct
	LegacyScriptInfo struct {
		PartnerID         string                `json:"partnerID"`
		LegacyScriptID    string                `json:"legacyScriptID"`
		LegacyTemplateID  string                `json:"legacyTemplateID"`
		DefinitionID      string                `json:"definitionID,omitempty"`
		OriginID          string                `json:"originID"`
		IsSequence        bool                  `json:"isSequence,omitempty"` // this means that script was converted to sequence
		IsParametrized    bool                  `json:"isParametrized"`       //
		ErrorReason       string                `json:"reason,omitempty"`
		DefinitionDetails TaskDefinitionDetails `json:"definitionDetails,omitempty"`
	}

	// LegacyJobInfo represents legacy job report and mapping info struct
	LegacyJobInfo struct {
		PartnerID        string `json:"partnerID"`
		LegacyJobID      string `json:"jobID"`
		LegacyScriptID   string `json:"legacyScriptID"`
		LegacyTemplateID string `json:"legacyTemplateID"`
		Type             string `json:"type"` // sequence or script
		DefinitionID     string `json:"definitionID,omitempty"`
		OriginID         string `json:"scriptID"`
		TaskID           string `json:"taskID,omitempty"`
		ErrorReason      string `json:"reason"`
	}
)
