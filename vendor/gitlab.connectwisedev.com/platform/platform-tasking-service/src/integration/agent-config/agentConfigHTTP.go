package agentConfig

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gocql/gocql"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	transactionID "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
)

const (
	wrongStatusCode        = "agent-config MS returned wrong status code: %v"
	taskingIdentifier      = "platform-tasking-service"
	description            = "Configuration change by tasking MS according to trigger definition with id: %v."
	eventlogConfigFileName = "plugin/eventlog/eventlog_agent_plugin_cfg.json"
	eventlogPluginName     = "platform-eventlog-plugin"
	pathToPatch            = "/rules/-"
	addOperation           = "add"
	agentConfURL           = "%s/partners/%s/profiles"
	agentConfDeleteURL     = "%s/partners/%s/profiles/%v"
	contentType            = "Content-Type"
	contentTypeValue       = "application/json"
	uidKey                 = "uid"
)

// NewAgentConfClient returns AgentConfClient
func NewAgentConfClient(cli integration.HTTPClient, domain string, log logger.Logger) *agentConfClient {
	return &agentConfClient{
		cli:               cli,
		log:               log,
		agentConfigDomain: domain,
	}
}

// agentConfClient is a service to communicate with sysEvent MS
type agentConfClient struct {
	cli               integration.HTTPClient
	log               logger.Logger
	agentConfigDomain string
}

// Activate sends HTTP request to AgentConfClient MS to activate triggers using payload
func (s *agentConfClient) Activate(ctx context.Context, content entities.Rule, managedEndpointsIDs map[string]entities.Endpoints, partnerID string) (profileID gocql.UUID, err error) {
	agentConfigPayload := s.buildAgentConfigActivatePayload(content, managedEndpointsIDs, partnerID)

	b, err := json.Marshal(agentConfigPayload)
	if err != nil {
		return profileID, fmt.Errorf("AgentConfClient.Activate err during marshaling: %v", err)
	}

	url := fmt.Sprintf(agentConfURL, s.agentConfigDomain, partnerID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return profileID, fmt.Errorf("AgentConfClient.Activate err during creating request: %v", err)
	}

	req.Header.Add(contentType, contentTypeValue)
	req.Header.Add(uidKey, taskingIdentifier)
	req.Header.Add(transactionID.Key, transactionID.FromContext(ctx))

	s.log.DebugfCtx(ctx, "Performing request by url %v with payload '%s'", url, string(b))

	resp, err := s.cli.Do(req)
	if err != nil {
		return profileID, fmt.Errorf("AgentConfClient.Activate err: %v", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			s.log.WarnfCtx(ctx,"AgentConfClient.Activate err closing body: %v", err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return profileID, err
	}

	if resp.StatusCode != http.StatusCreated {
		return profileID, fmt.Errorf("agent-config MS returned wrong status code: %v", resp.StatusCode)
	}

	var respStructure entities.AgentActivateResp
	err = json.Unmarshal(body, &respStructure)
	if err != nil {
		return profileID, err
	}
	s.log.InfofCtx(ctx, "AgentConfClient.Activate: agent-config returned %v and ID %v", resp.StatusCode, respStructure.ProfileID)
	return respStructure.ProfileID, nil
}

func (s *agentConfClient) buildAgentConfigActivatePayload(content entities.Rule, managedEndpointsIDs map[string]entities.Endpoints, partnerID string) entities.AgentConfigPayload {
	patch := entities.Patch{
		Op:    addOperation,
		Path:  pathToPatch,
		Value: content, // this is what going to be deployed in plugin's config
	}

	configuration := entities.Configuration{
		Patch:       []entities.Patch{patch},
		FileName:    eventlogConfigFileName,
		PackageName: eventlogPluginName, // plugin we deploy trigger config to
	}

	configurations := []entities.Configuration{configuration}

	targets := make([]entities.AgentTargets, 0)
	for id, e := range managedEndpointsIDs {
		targets = append(targets, entities.AgentTargets{
			PartnerID:  partnerID,
			ClientID:   e.ClientID,
			SiteID:     e.SiteID,
			EndpointID: id,
		})
	}

	agentConfigPayload := entities.AgentConfigPayload{
		Tag:           taskingIdentifier,
		Description:   fmt.Sprintf(description, content.TriggerID),
		Configuration: configurations,
		Targets:       targets,
	}
	return agentConfigPayload
}

// Deactivate sends HTTP request to AgentConfClient MS to deactivate triggers using payload
func (s *agentConfClient) Deactivate(ctx context.Context, profileID gocql.UUID, partnerID string) (err error) {
	url := fmt.Sprintf(agentConfDeleteURL, s.agentConfigDomain, partnerID, profileID.String())
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("AgentConfClient.Deactivate err: %v", err)
	}

	req.Header.Add(transactionID.Key, transactionID.FromContext(ctx))

	resp, err := s.cli.Do(req)
	if err != nil {
		return fmt.Errorf("AgentConfClient.Deactivate err: %v", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			s.log.WarnfCtx(ctx,"AgentConfClient.Deactivate err closing body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf(wrongStatusCode, resp.StatusCode)
	}
	return nil
}

// Update sends HTTP request to AgentConfClient MS to update triggers using payload
func (s *agentConfClient) Update(ctx context.Context, content entities.Rule, managedEndpointsIDs map[string]entities.Endpoints, partnerID string, profileID gocql.UUID) (err error) {
	agentConfigPayload := s.buildAgentConfigActivatePayload(content, managedEndpointsIDs, partnerID)

	b, err := json.Marshal(agentConfigPayload)
	if err != nil {
		return fmt.Errorf("AgentConfClient.Update err during marshaling: %v", err)
	}

	url := fmt.Sprintf(agentConfDeleteURL, s.agentConfigDomain, partnerID, profileID.String())
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("AgentConfClient.Update err during creating request: %v", err)
	}

	req.Header.Add(contentType, contentTypeValue)
	req.Header.Add(uidKey, taskingIdentifier)
	req.Header.Add(transactionID.Key, transactionID.FromContext(ctx))

	resp, err := s.cli.Do(req)
	if err != nil {
		return fmt.Errorf("AgentConfClient.Update  err: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.log.WarnfCtx(ctx,"AgentConfClient.Update  err closing body: %v", err)
		}
	}()

	s.log.DebugfCtx(ctx, "AgentConfClient.Update: sent body %s to agent-config service", b)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(wrongStatusCode, resp.StatusCode)
	}

	var respStructure entities.AgentActivateResp
	err = json.Unmarshal(body, &respStructure)
	if err != nil {
		return err
	}
	s.log.DebugfCtx(ctx, "AgentConfClient.Update: agent-config returned %v and ID %v", resp.StatusCode, respStructure.ProfileID)
	return nil
}
