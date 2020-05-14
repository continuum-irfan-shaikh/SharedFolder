package sites

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	transactionID "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
)

const partnerSitesURLPattern = "%s/partners/%s/sites?Operation=activesites"

// GetSiteIDs gets site IDs
func GetSiteIDs(ctx context.Context, httpClient *http.Client, partnerID, sitesELB, token string) (siteIDs []int64, err error) {
	const iPlanetDirectoryPro = `iPlanetDirectoryPro`
	if token == "" {
		return getPartnerSites(ctx, httpClient, partnerID)
	}

	var (
		sitesData struct {
			SiteList []struct {
				ID int64 `json:"siteId"`
			} `json:"siteDetailList"`
		}
		request  *http.Request
		response *http.Response
		url      = fmt.Sprintf("%s/partner/%s/sites", sitesELB, partnerID)
	)
	siteIDs = make([]int64, 0)

	request, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	request.Header.Add(iPlanetDirectoryPro, token)
	request.Header.Add(transactionID.Key, transactionID.FromContext(ctx))

	response, err = httpClient.Do(request)
	if err != nil {
		return
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			logger.Log.WarnfCtx(ctx, "GetSiteIDs: error while closing response body: %v", err)
		}
	}()

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("got wrong http status [%d]; expected status [%d]", response.StatusCode, http.StatusOK)
		return
	}

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		err = readErr
		return
	}

	err = json.Unmarshal(body, &sitesData)
	if err != nil {
		return
	}

	for _, site := range sitesData.SiteList {
		siteIDs = append(siteIDs, site.ID)
	}
	return
}

func getPartnerSites(ctx context.Context, httpClient integration.HTTPClient, partnerID string) ([]int64, error) {
	url := fmt.Sprintf(partnerSitesURLPattern, config.Config.SitesNoTokenURL, partnerID)

	var sitesList struct {
		Sites []struct {
			ID int64 `json:"siteId"`
		} `json:"outdata"`
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(transactionID.Key, transactionID.FromContext(ctx))

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logger.Log.WarnfCtx(ctx,"getPartnerSites: error while closing response body: %v", err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &sitesList); err != nil {
		return nil, err
	}

	s := make([]int64, 0)
	for _, v := range sitesList.Sites {
		s = append(s, v.ID)
	}
	return s, nil
}
