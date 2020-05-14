package site

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type sites struct {
	SiteList []struct {
		ID int64 `json:"siteId"`
	} `json:"siteDetailList"`
}

// NewSite returns new Sites repo
func NewSite(client *http.Client, siteDomain string) *Site {
	return &Site{
		client:     client,
		siteDomain: siteDomain,
	}
}

// Site represents site client
type Site struct {
	client     *http.Client
	siteDomain string
}

// Sites returns sites
func (s *Site) Sites(partnerID, token string) ([]string, error) {
	url := fmt.Sprintf("%s/partner/%s/sites", s.siteDomain, partnerID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(`iPlanetDirectoryPro`, token)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	sitesList := sites{}
	if err := json.Unmarshal(body, &sitesList); err != nil {
		return nil, err
	}

	sites := make([]string, 0)
	for _, v := range sitesList.SiteList {
		sites = append(sites, strconv.FormatInt(v.ID, 10))
	}

	return sites, nil
}
