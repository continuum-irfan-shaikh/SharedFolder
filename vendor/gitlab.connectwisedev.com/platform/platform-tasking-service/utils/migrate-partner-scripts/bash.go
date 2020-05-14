package main

import (
	"encoding/xml"
	"strconv"

	"github.com/gocql/gocql"
)

func getBashTaskDef(sc []ScriptMsSQL) (tds []TaskDefinition, err error) {
	for _, s := range sc {
		var (
			xmlData    BashXMLData
			td         TaskDefinition
			userParams CustomScriptBody
		)

		err = xml.Unmarshal([]byte("<b>"+s.ScriptXML+"</b>"), &xmlData)
		if err != nil {
			return
		}

		//--- user parameters XML part
		if xmlData.Steps.Body != "" {
			userParams.ExpExecTime = xmlData.Steps.Timeout * 60 // stored in seconds
			userParams.ScriptBody = xmlData.Steps.Body          //  must be decoded somehow here

			body, err := json.Marshal(userParams)
			if err != nil {
				return nil, err
			}
			td.UserParameters = string(body)
		}
		//--- assigning to TD
		td.Categories = append(td.Categories, "Custom")
		td.CreatedBy = s.CreatedBy
		td.CreatedAt = s.CreatedAt
		td.PartnerID = strconv.Itoa(s.MemberID)
		td.Description = s.ScriptDesc
		td.Name = s.ScriptName
		td.OriginID = bashOriginID
		td.ID = gocql.TimeUUID()
		td.Type = scriptType
		tds = append(tds, td)
	}
	return
}
