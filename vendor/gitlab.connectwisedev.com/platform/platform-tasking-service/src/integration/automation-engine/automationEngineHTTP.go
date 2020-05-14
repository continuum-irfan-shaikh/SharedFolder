package automationEngine

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	t "html/template"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	e "gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	transactionID "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
)

const (
	policyParam           = "policy"
	policyVersion         = "1"
	policyFileName        = "policy.zip"
	policyMessageType     = "message"
	policyMessageDataType = "string"
	apiPath               = `/automationengine/policies`
	closingFileError      = "error while closing the file %v"
	partnerIDField        = "partnerID"
	clientIDField         = "clientID"
	siteIDField           = "siteID"
	endpointIDField       = "endpointID"
	httpProtocol          = "http"
	actionName            = "action1"
	mode                  = "PARALLEL"
	endResource           = "execute-trigger"
	uid                   = "AUTOMATION_ENGINE"
	realm                 = "//activedirectory"
	//triggerURL path variables
	partners      = "partners"
	clients       = "clients"
	sites         = "sites"
	endpoints     = "endpoints"
	taskingSource = "tasking"
)

type automationEngineResponse struct {
	ID string `json:"systemidhash"`
}

// Client represents automation engine client to work via http
type Client struct {
	log           logger.Logger
	cli           integration.HTTPClient
	engineDomain  string
	taskingDomain string
}

// New returns new AE http client
func New(l logger.Logger, cli *http.Client, ae, tasking string) *Client {
	return &Client{
		log:           l,
		cli:           cli,
		engineDomain:  ae,
		taskingDomain: tasking,
	}
}

// GeneratePolicyFiles generates and creates policy files from trigger config file. returns slice of file names
func (c *Client) generatePolicyFiles(data []e.TriggerDefinition) (tempFolderName string, fileNames []string, err error) {
	policies, err := c.buildPolicies(data)
	if err != nil {
		return tempFolderName, fileNames, err
	}

	var template *t.Template

	tempFolderName = policyParam + time.Now().String()
	err = os.Mkdir(tempFolderName, os.ModePerm)
	if err != nil {
		return tempFolderName, fileNames, err
	}

	for _, policy := range policies {
		fileName := fmt.Sprintf("%v.policy", policy.ID)
		file, err := os.Create(tempFolderName + "/" + fileName)
		if err != nil {
			return tempFolderName, fileNames, err
		}

		templateFormat := ""
		if policy.IsAlertingPolicy() {
			templateFormat = alertPolicyTemplate
		} else {
			templateFormat = genericPolicyTemplate
		}

		systemIdentifierCounter := 0
		actionCounter := 0
		template = t.Must(t.New("policy").Funcs(t.FuncMap{
			"incrementCounter": func() string {
				systemIdentifierCounter++
				return ""
			},
			"incrementActionCounter": func() string {
				actionCounter++
				return ""
			},
			"IsNotLastSystemIdentifier": func() bool {
				return !(systemIdentifierCounter == len(policy.Variables))
			},
			"IsNotLastAction": func() bool {
				return !(actionCounter == len(policy.Actions))
			},
		}).Parse(templateFormat))

		// template needed vars
		if err = template.Execute(file, policy); err != nil {
			return tempFolderName, fileNames, err
		}
		fileNames = append(fileNames, fileName)
	}
	return
}

// buildPolicies build policies from trigger definitions
func (c *Client) buildPolicies(defs []e.TriggerDefinition) (policies []e.Policy, err error) {
	for _, def := range defs {
		variablesMap := make(map[string]e.PolicyVars)
		payloadMap := make(map[string]string)
		variablesMap[partnerIDField] = e.PolicyVars{
			Type:     policyMessageType,
			DataType: policyMessageDataType,
			Key:      def.EventDetails.EndpointIdentifiers.PartnerID.Value,
		}
		variablesMap[siteIDField] = e.PolicyVars{
			Type:     policyMessageType,
			DataType: policyMessageDataType,
			Key:      def.EventDetails.EndpointIdentifiers.SiteID.Value,
		}
		variablesMap[clientIDField] = e.PolicyVars{
			Type:     policyMessageType,
			DataType: policyMessageDataType,
			Key:      def.EventDetails.EndpointIdentifiers.ClientID.Value,
		}
		variablesMap[endpointIDField] = e.PolicyVars{
			Type:     policyMessageType,
			DataType: policyMessageDataType,
			Key:      def.EventDetails.EndpointIdentifiers.EndpointID.Value,
		}

		// adding other keys, that used in payload
		for k, payloadVar := range def.EventDetails.PayloadIdentifiers {
			variablesMap[k] = e.PolicyVars{
				Type:     policyMessageType,
				DataType: policyMessageDataType,
				Key:      payloadVar.Value,
			}
			payloadMap[k] = k
		}

		taskingURL, err := url.Parse(c.taskingDomain)
		if err != nil {
			return policies, err
		}

		actions := []e.Action{
			{
				TriggerID:   def.ID,
				Name:        actionName,
				Mode:        mode,
				Protocol:    httpProtocol,
				Endpoint:    taskingURL.Host,
				Method:      http.MethodPost,
				Context:     taskingURL.Path[1:],
				EndResource: endResource,
				PathVariables: []e.KeyValue{
					{
						Key:   partners,
						Value: partnerIDField,
					},
					{
						Key:   clients,
						Value: clientIDField,
					},
					{
						Key:   sites,
						Value: siteIDField,
					},
					{
						Key:   endpoints,
						Value: endpointIDField,
					},
				},
				Payload: payloadMap,
				Headers: e.Headers{
					Uid:   uid,
					Realm: realm,
				},
			},
		}
		policy := e.Policy{
			ID:               def.ID,
			Description:      def.Description,
			Type:             policyParam,
			SystemIdentifier: def.EventDetails.MessageIdentifier,
			Name:             def.DisplayName,
			Version:          policyVersion,
			Topic:            def.EventDetails.Topic,
			Variables:        variablesMap,
			Actions:          actions,
		}
		policies = append(policies, policy)
	}
	return
}

const updatePolicyErr = "UpdateRemotePolicies: %v"

// UpdateRemotePolicies creates zip file from given policy files and sends it to AE MS via http
func (c *Client) UpdateRemotePolicies(ctx context.Context, defs []e.TriggerDefinition) (policyID string, err error) {
	folderName, files, err := c.generatePolicyFiles(defs)
	if err != nil {
		return
	}

	// compresses policy files into zip file
	if err = c.zipFiles(ctx, folderName, files); err != nil {
		return policyID, fmt.Errorf(updatePolicyErr, err)
	}
	defer func() {
		if err := c.removeCreatedFiles(folderName); err != nil {
			c.log.WarnfCtx(ctx,"UpdateRemotePolicies: error during clearing files: %v", err)
		}
	}()

	req, err := c.fileUploadRequest(ctx, c.engineDomain+apiPath, folderName+"/"+policyFileName)
	if err != nil {
		return policyID, fmt.Errorf(updatePolicyErr, err)
	}
	req.Header.Set(transactionID.Key, transactionID.FromContext(ctx))

	c.log.DebugfCtx(ctx, "Request POST url %v", req.URL)
	res, err := c.cli.Do(req)
	if err != nil {
		return policyID, fmt.Errorf(updatePolicyErr, err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			logger.Log.WarnfCtx(ctx,"CreateExecutionPayload: error while closing body: %v", err)
		}
	}()
	c.log.DebugfCtx(ctx, "Response from POST url %v , status code '%v'", req.URL, res.StatusCode)

	if res.StatusCode != http.StatusCreated {
		return policyID, fmt.Errorf("UpdateRemotePolicies: bad status, expected 201 but got %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	c.log.DebugfCtx(ctx, "Response from POST url %v , payload '%s'", req.URL, string(body))

	var resp []automationEngineResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return
	}
	if len(resp) == 0 {
		return "", fmt.Errorf("got empty response from AE")
	}

	if resp[0].ID == "" {
		return "", fmt.Errorf("got empty policy ID, something went wrong")
	}
	return resp[0].ID, nil
}

type removePayload struct {
	SystemIdentifier map[string]interface{} `json:"systemidentifier"`
	SourceSystem     string                 `json:"sourcesystem"`
}

// RemovePolicy creates zip file from given policy files and sends it to AE MS via http
func (c *Client) RemovePolicy(ctx context.Context, sysIdentifier map[string]interface{}) (err error) {
	payload := removePayload{
		SystemIdentifier: sysIdentifier,
		SourceSystem:     taskingSource,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("RemovePolicy: marshal %v", err)
	}

	url := c.engineDomain + apiPath
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("RemovePolicy: create request %v", err)
	}

	c.log.DebugfCtx(ctx, "Performing request by url %v with payload '%s'", url, string(b))

	req.Header.Set(transactionID.Key, transactionID.FromContext(ctx))
	req.Header.Add("Content-type", "application/json")
	res, err := c.cli.Do(req)
	if err != nil {
		return fmt.Errorf("RemovePolicy: send error %v", err)
	}
	c.log.DebugfCtx(ctx, "Response from POST url %v , status code '%v'", req.URL, res.StatusCode)

	// means that policy has been already removed, nothing to remove
	if res.StatusCode == http.StatusInternalServerError {
		return nil
	}

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("RemovePolicy: bad status, expected 204 but got %s", res.Status)
	}
	return
}

// Creates a new file upload http request
func (c *Client) fileUploadRequest(ctx context.Context, url, filename string) (req *http.Request, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer c.fileCloser(ctx, file)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(policyParam, filepath.Base(file.Name()))
	if err != nil {
		return
	}
	if _, err = io.Copy(part, file); err != nil {
		return
	}
	if err = writer.Close(); err != nil {
		return
	}

	if req, err = http.NewRequest(http.MethodPost, url, body); err != nil {
		return
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return
}

// creates zip file from given policies
func (c *Client) zipFiles(ctx context.Context, folderName string, files []string) error {
	policiesArchive, err := os.Create(folderName + "/" + policyFileName)
	if err != nil {
		return err
	}
	defer c.fileCloser(ctx, policiesArchive)

	zipWriter := zip.NewWriter(policiesArchive)
	defer func() {
		if err := zipWriter.Close(); err != nil {
			c.log.WarnfCtx(ctx, closingFileError, err)
		}
	}()

	for _, file := range files {
		if err = c.addFileToZip(ctx, zipWriter, folderName, file); err != nil {
			return err
		}
	}
	return nil
}

// adds file to existing zip file
func (c *Client) addFileToZip(ctx context.Context, zipWriter *zip.Writer, folderName, filename string) error {
	fileToZip, err := os.Open(folderName + "/" + filename)
	if err != nil {
		return err
	}
	defer c.fileCloser(ctx, fileToZip)

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filename
	header.Method = zip.Deflate
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

// ClearGeneratedPolicies remove previously created policy files
func (c *Client) removeCreatedFiles(folderName string) error {
	return os.RemoveAll(folderName)
}

func (c *Client) fileCloser(ctx context.Context, file *os.File) {
	if err := file.Close(); err != nil {
		c.log.WarnfCtx(ctx, closingFileError, err)
	}
}

const alertPolicyTemplate = `{
	"policyid": "{{.ID}}",
	"description": "{{.Description}}",
	"type": "{{.Type}}",
	"systemidentifier": {
	    {{ range $key, $value := .SystemIdentifier}}
		    "{{ $key }}": {{ $value }}
	    {{end}}
	},
	"name": "{{.Name}}",
	"sourcesystem":"tasking",
	"version": "{{.Version}}",
	"eventTopic": "{{.Topic}}",
	"variables": {
         {{ range $key, $value := .Variables}}
        		"{{ $key }}": {
        		    "type":"{{ $value.Type }}",
        		    "key":"{{ $value.Key }}",
                    "dataType":"{{ $value.DataType}}"
        		}{{ incrementCounter}}{{if IsNotLastSystemIdentifier }},{{end}}
        {{end}}
	},
	"actions": [
        {{range $i, $a := .Actions}}{
            "Name":"{{ $a.Name }}",
            "Mode":"{{ $a.Mode }}",
            "Protocol":"{{ $a.Protocol }}",
            "Endpoint":"{{ $a.Endpoint }}",
            "Method":"{{ $a.Method }}",
            "Context":"{{ $a.Context }}",
            "EndResource":"{{ $a.EndResource }}",
            "PathVariables": [
             {{range $j, $b := $a.PathVariables}}
                {
                    "{{$b.Key}}": ${{$b.Value}}
                },
             {{end}}
                {
                    "triggers": "{{.TriggerID}}"
                }
            ],
            "payload": {
                {{range $key, $value := $a.Payload}}
                    "{{ $key }}": ${{ $value }},
                {{end}}
				"dummy": 1
            },
		    "headers": {
			    "uid": "{{ .Headers.Uid }}",
			    "realm":"{{ .Headers.Realm }}"
		    }
        }
 		{{incrementActionCounter}}
        {{if IsNotLastAction }}
        ,
        {{end}}
    {{end}}
	]
}`

const genericPolicyTemplate = `{
	"policyid": "{{.ID}}",
	"description": "{{.Description}}",
	"type": "{{.Type}}",
	"systemidentifier": {
	    {{ range $key, $value := .SystemIdentifier}}
		    "{{ $key }}": "{{ $value }}"
	    {{end}}
	},
	"name": "{{.Name}}",
	"sourcesystem":"tasking",
	"version": "{{.Version}}",
	"eventTopic": "{{.Topic}}",
	"variables": {
         {{ range $key, $value := .Variables}}
        		"{{ $key }}": {
        		    "type":"{{ $value.Type }}",
        		    "key":"{{ $value.Key }}",
                    "dataType":"{{ $value.DataType}}"
        		}{{ incrementCounter}}{{if IsNotLastSystemIdentifier }},{{end}}
        {{end}}
	},
	"actions": [
        {{range $i, $a := .Actions}}{
            "Name":"{{ $a.Name }}",
            "Mode":"{{ $a.Mode }}",
            "Protocol":"{{ $a.Protocol }}",
            "Endpoint":"{{ $a.Endpoint }}",
            "Method":"{{ $a.Method }}",
            "Context":"{{ $a.Context }}",
            "EndResource":"{{ $a.EndResource }}",
            "PathVariables": [
             {{range $j, $b := $a.PathVariables}}
                {
                    "{{$b.Key}}": ${{$b.Value}}
                },
             {{end}}
                {
                    "triggers": "{{.TriggerID}}"
                }
            ],
            "payload": {
                {{range $key, $value := $a.Payload}}
                    "{{ $key }}": ${{ $value }},
                {{end}}
				"dummy": 1
            },
		    "headers": {
			    "uid": "{{ .Headers.Uid }}",
			    "realm":"{{ .Headers.Realm }}"
		    }
        }
 		{{incrementActionCounter}}
        {{if IsNotLastAction }}
        ,
        {{end}}
    {{end}}
	]
}`
