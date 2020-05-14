package main

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"chilkat"

	"github.com/gocql/gocql"
)

const (
	startOfRedundantSymbols = 15
	endOfRedundantSymbols   = 13
)

func getPowerShellTaskDef(sc []ScriptMsSQL) (td []TaskDefinition, err error) {
	for _, s := range sc {
		var (
			xmlData    PowershellXMLData
			t          TaskDefinition
			userParams CustomScriptBody
			script     string
		)

		err = xml.Unmarshal([]byte("<b>"+s.ScriptXML+"</b>"), &xmlData)
		if err != nil {
			return
		}

		//--- assigning to TD
		t.Categories = append(t.Categories, xmlData.Category)
		t.CreatedBy = s.CreatedBy
		t.CreatedAt = s.CreatedAt
		t.PartnerID = strconv.Itoa(s.MemberID)
		t.Description = s.ScriptDesc
		t.Name = s.ScriptName
		t.OriginID = powershellOriginID
		t.ID = gocql.TimeUUID()
		t.Type = scriptType
		t.Deleted = false

		c := chilkat.NewCrypt2()

		//--- user parameters XML part
		hasParameters := false
		for _, step := range xmlData.Steps {
			if step.Parameters != "" {
				hasParameters = true
				break
			}

			// this step removes `<unnamed>=??B?` from the start and  `?=</unnamed>` from the end
			scriptBody := step.Body[startOfRedundantSymbols : len(step.Body)-endOfRedundantSymbols]

			b := c.DecodeString(scriptBody, "ASCII", "base64")
			if b == nil {
				return nil, fmt.Errorf("can not decode stringBody: %s", scriptBody)
			}

			userParams.ExpExecTime += step.Timeout * 60 // stored in mins
			script += *b + "\n"
		}

		// skip scripts with parameters for now...
		if hasParameters {
			continue
		}

		userParams.ScriptBody = script

		b, err := json.Marshal(userParams)
		if err != nil {
			return nil, err
		}

		t.UserParameters = string(b)
		td = append(td, t)
	}
	return
}
