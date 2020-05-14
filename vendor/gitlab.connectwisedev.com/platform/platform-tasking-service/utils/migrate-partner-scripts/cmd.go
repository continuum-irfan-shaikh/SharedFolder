package main

import (
	"encoding/xml"
	"fmt"
	"github.com/gocql/gocql"
	"strconv"
)

const pingTmp = "ping 127.0.0.1 -n1 -w %v >NUL\n"

func cmdDefs(oldScripts []ScriptMsSQL) ([]TaskDefinition, error) {
	definitions := make([]TaskDefinition, 0, len(oldScripts))

	for _, script := range oldScripts {

		originID, err := getOriginID(script.TemplateID)
		if err != nil {
			return nil, err
		}

		body, categories, err := parseCMDScript(script.ScriptXML)
		if nil != err {
			return nil, err
		}

		definition := TaskDefinition{
			ID:             gocql.TimeUUID(),
			PartnerID:      getPartnerID(script.MemberID),
			OriginID:       originID,
			Name:           script.ScriptName,
			Type:           scriptType,
			Categories:     categories,
			Deleted:        false,
			Description:    script.ScriptDesc,
			CreatedAt:      script.CreatedAt,
			CreatedBy:      script.CreatedBy,
			UpdatedAt:      script.UpdatedAt,
			UpdatedBy:      script.UpdatedBy,
			UserParameters: body,
		}

		definitions = append(definitions, definition)

	}

	return definitions, nil

}

func getPartnerID(memberID int) (partnerID string) {
	//TODO: maybe we need some mapping here as well
	return strconv.Itoa(memberID)
}

func getOriginID(templateID int) (originID gocql.UUID, err error) {
	var (
		originIDStr string
		ok          bool
	)
	if originIDStr, ok = templateToOrigin[templateID]; !ok {
		return gocql.UUID{}, fmt.Errorf("there is no appropriate originID for templateID : %v", templateID)
	}

	originID, err = gocql.ParseUUID(originIDStr)
	if err != nil {
		return gocql.UUID{}, fmt.Errorf("there is no appropriate originID for templateID : %v", templateID)
	}
	return
}

func parseCMDScript(scriptXML string) (scriptBody string, categories []string, err error) {
	var (
		xmlData    CmdXMLData
		userParams CustomScriptBody
	)
	err = xml.Unmarshal([]byte("<b>"+scriptXML+"</b>"), &xmlData)
	if err != nil {
		return "", nil, fmt.Errorf("scriptXML data %v can not be unmarshaled: %s", scriptXML, err.Error())
	}

	categories = append(categories, xmlData.Category)

	for _, step := range xmlData.Steps {
		userParams.ExpExecTime = userParams.ExpExecTime + step.TimeoutMin*60
		userParams.ScriptBody = userParams.ScriptBody + step.Body + " " + step.Parameters + "\n"
		if step.PauseSec > 0 {
			userParams.ScriptBody = userParams.ScriptBody + fmt.Sprintf(pingTmp, step.PauseSec*1000)
		}
	}

	body, err := json.Marshal(userParams)
	if err != nil {
		return "", nil, fmt.Errorf("userParams : %v can not be marshaled to json: %s", userParams, err.Error())
	}
	scriptBody = string(body)

	return scriptBody, categories, nil

}
