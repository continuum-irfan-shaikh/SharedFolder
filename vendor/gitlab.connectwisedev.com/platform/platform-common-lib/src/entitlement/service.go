package entitlement

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	apiModel "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/entitlement"
	"github.com/coocood/freecache"
)

// Service represents an Entitlement Service type
type Service struct {
	cache           *freecache.Cache
	httpClient      *http.Client
	url             string
	cacheDataTTLSec int
}

// NewEntitlementService creates a new Entitlement Service
func NewEntitlementService(httpClient *http.Client, entitlementMsURL string, cacheDataTTLSec, cacheSize int) Service {
	return Service{
		cache:           freecache.NewCache(cacheSize),
		httpClient:      httpClient,
		url:             entitlementMsURL,
		cacheDataTTLSec: cacheDataTTLSec,
	}
}

// GetPartnerFeatures retrieve features for Partner from Entitlement MS or from local cache
func (es Service) GetPartnerFeatures(partnerID string) (features []apiModel.Feature, err error) {
	var featuresBin []byte
	partnerIDBin := []byte(partnerID)

	featuresBin, err = es.cache.Get(partnerIDBin)
	if err != nil {
		resp, entitlementErr := es.httpClient.Get(es.url + "/partners/" + partnerID + "/features")
		if entitlementErr != nil {
			return features, entitlementErr
		}
		defer resp.Body.Close()

		featuresBin, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return features, err
		}

		err = es.cache.Set(partnerIDBin, featuresBin, es.cacheDataTTLSec)

		if err != nil {
			return features, err
		}
	}

	err = json.Unmarshal(featuresBin, &features)

	return features, err
}

// IsPartnerAuthorized checks if the Partner has enabled feature in the Entitlement Service
func (es Service) IsPartnerAuthorized(partnerID, featureName string) bool {
	features, err := es.GetPartnerFeatures(partnerID)
	if err != nil {
		return false
	}

	for _, feature := range features {
		if feature.Name == featureName {
			return true
		}
	}
	return false
}
