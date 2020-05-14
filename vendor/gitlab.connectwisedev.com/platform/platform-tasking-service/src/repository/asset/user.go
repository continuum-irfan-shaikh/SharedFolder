package asset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	en "gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
)

type endpoint struct {
	EndpointID string `json:"endpointID"`
	SiteID     string `json:"siteID"`
	ClientID   string `json:"clientID"`
}

func NewUser(client *http.Client, assetDomain string) *user {
	return &user{
		client:      client,
		assetDomain: assetDomain,
	}
}

type user struct {
	client      *http.Client
	assetDomain string
}

// Endpoints - get endpoints by sites and partnerID from Asset service
func (u *user) Endpoints(siteIDs []string, partnerID string) ([]en.Endpoints, error) {
	urlPattern := "%s/partner/%s/sites/%s/summary"
	url := fmt.Sprintf(urlPattern, u.assetDomain, partnerID, strings.Join(siteIDs, "%2C"))

	resp, err := u.client.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	sites := make([]endpoint, 0)
	if err := json.Unmarshal(body, &sites); err != nil {
		return nil, err
	}

	// map endpoints by endpoint
	epBySite := make(map[string]en.Endpoints, 0)
	for _, v := range sites {
		_, ok := epBySite[v.SiteID]
		if !ok {
			epBySite[v.SiteID] = en.Endpoints{
				PartnerID: partnerID,
				SiteID:    v.SiteID,
				ClientID:  v.ClientID,
			}
		}
		ep := epBySite[v.SiteID]
		ep.Endpoints = append(ep.Endpoints, v.EndpointID)
		epBySite[v.SiteID] = ep
	}

	eps := make([]en.Endpoints, 0, len(epBySite))
	for _, v := range epBySite {
		eps = append(eps, v)
	}

	return eps, nil
}
