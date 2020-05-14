package asset

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	apiModel "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/asset"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	transactionID "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/memcached"
)

const (
	timeZoneOffsetKeyPrefix = "TKS_TIMEZONE_OFFSET_"
	siteIDKeyPrefix         = "TKS_SITE_ID_"
	clientIDKeyPrefix       = "TKS_CLIENT_ID_"
	machineNameKeyPrefix    = "TKS_MACHINE_NAME_"
	resourceTypeKeyPrefix   = "TKS_RESOURCE_TYPE_"
)

// Service represents a Asset Service type
type Service struct {
	mCache     memcached.Cache
	httpClient integration.HTTPClient
}

// ServiceInstance an instance of Asset MS Service
var ServiceInstance integration.Asset

// Load creates new AssetsService
func Load() {
	ServiceInstance = NewAssetsService(memcached.MemCacheInstance, defaultHTTPClient)
}

// NewAssetsService creates a new NewAssetsService Service
func NewAssetsService(memCache memcached.Cache, HTTPClient *http.Client) integration.Asset {
	if HTTPClient == nil {
		HTTPClient = defaultHTTPClient
	}
	return Service{
		mCache:     memCache,
		httpClient: HTTPClient,
	}
}

var (
	defaultHTTPClient = &http.Client{
		Timeout: time.Duration(config.Config.HTTPClientTimeoutSec) * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout:     time.Duration(config.Config.HTTPClientTimeoutSec) * time.Second,
			MaxIdleConns:        2 * config.Config.HTTPClientMaxIdleConnPerHost,
			MaxIdleConnsPerHost: config.Config.HTTPClientMaxIdleConnPerHost,
			DisableKeepAlives:   false,
		},
	}
)

// GetLocationByEndpointID get managed endpoint location by partner ID, endpoint ID.
// Retrieve timezone offset from cache if it is enabled in config otherwise obtain it from Asset MS then create location by timezone offset.
func (assetsService Service) GetLocationByEndpointID(ctx context.Context, partnerID string, endpointID gocql.UUID) (location *time.Location, err error) {
	var (
		timeZoneOffset string
		offset         time.Duration
		errCache       error
	)

	if config.Config.AssetCacheEnabled {
		var timeZoneBin *memcache.Item
		timeZoneBin, errCache = assetsService.mCache.Get(timeZoneOffsetKeyPrefix + endpointID.String())
		if errCache == nil {
			timeZoneOffset = string(timeZoneBin.Value)
		}
	}

	if !config.Config.AssetCacheEnabled || errCache != nil {
		asset, err := assetsService.getDataFromAssetMS(ctx, partnerID, endpointID)
		if err != nil {
			return time.UTC, err
		}
		timeZoneOffset = asset.System.TimeZone
	}

	offset, err = parseTimeZoneOffset(timeZoneOffset)
	if err != nil {
		logger.Log.WarnfCtx(ctx,"Couldn't parse timezone offset %s for endpoint(ID=%v), err:", timeZoneOffset, endpointID, err)
		// If received offset from AssetMS is not valid, location == nil (UTC)
		return time.UTC, nil
	}

	location = time.FixedZone("UTC"+timeZoneOffset, int(offset.Seconds()))
	return
}

// GetSiteIDByEndpointID get's sites ID by managedendpointID from asset MS
func (assetsService Service) GetSiteIDByEndpointID(ctx context.Context, partnerID string, endpointID gocql.UUID) (siteID, clientID string, err error) {
	if config.Config.AssetCacheEnabled {
		siteIDBin, errSite := assetsService.mCache.Get(siteIDKeyPrefix + endpointID.String())
		clientIDBin, errClient := assetsService.mCache.Get(clientIDKeyPrefix + endpointID.String())
		if errSite == nil && errClient == nil {
			return string(siteIDBin.Value), string(clientIDBin.Value), nil
		}
	}

	asset, err := assetsService.getDataFromAssetMS(ctx, partnerID, endpointID)
	if err != nil {
		return "", "", err
	}
	return asset.SiteID, asset.ClientID, nil
}

// GetResourceTypeByEndpointID get's endpoints resource type
func (assetsService Service) GetResourceTypeByEndpointID(ctx context.Context, partnerID string, endpointID gocql.UUID) (resType integration.ResourceType, err error) {
	if config.Config.AssetCacheEnabled {
		resTypeBin, err := assetsService.mCache.Get(resourceTypeKeyPrefix + endpointID.String())
		if err == nil {
			return integration.ResourceType(resTypeBin.Value), nil
		}
	}

	asset, err := assetsService.getDataFromAssetMS(ctx, partnerID, endpointID)
	if err != nil {
		return "", err
	}
	return integration.ResourceType(asset.EndpointType), nil
}

// GetMachineNameByEndpointID - gets machine name from asset MS
func (assetsService Service) GetMachineNameByEndpointID(ctx context.Context, partnerID string, endpointID gocql.UUID) (string, error) {
	if config.Config.AssetCacheEnabled {
		machineNameBin, err := assetsService.mCache.Get(machineNameKeyPrefix + endpointID.String())
		if err == nil {
			return string(machineNameBin.Value), nil
		}
	}
	asset, err := assetsService.getDataFromAssetMS(ctx, partnerID, endpointID)
	if err != nil {
		return "", err
	}
	return asset.System.SystemName, nil
}

// getDataFromAssetMS is used to get managed endpoint time zone offset by managed endpoint ID from AssetMS
func (assetsService Service) getDataFromAssetMS(ctx context.Context, partnerID string, endpointID gocql.UUID) (apiModel.AssetCollection, error) {
	asset := apiModel.AssetCollection{}
	assetMsURL := fmt.Sprintf("%s/partner/%s/endpoints/%v?field=system", config.Config.AssetMsURL, partnerID, endpointID)

	logger.Log.DebugfCtx(ctx, "Performing GET request by url %v", assetMsURL)

	req, err := http.NewRequest(http.MethodGet, assetMsURL, nil)
	if err != nil {
		return asset, err
	}

	req.Header.Set(transactionID.Key, transactionID.FromContext(ctx))
	resp, err := assetsService.httpClient.Do(req)
	if err != nil {
		return asset, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logger.Log.WarnfCtx(ctx,"getDataFromAssetMS: error while closing body: %v", err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return asset, err
	}

	logger.Log.DebugfCtx(ctx, "Response from GET url %v , status code '%v',  payload '%s'", assetMsURL, resp.StatusCode, string(body))
	err = json.Unmarshal(body, &asset)
	if err != nil {
		return asset, err
	}

	if config.Config.AssetCacheEnabled {
		fields := assetsService.createAssetDataMap(asset)
		assetsService.setAssetDataToCache(ctx, fields, endpointID.String())
	}

	return asset, nil
}

func (assetsService Service) createAssetDataMap(asset apiModel.AssetCollection) map[string]string {
	fields := make(map[string]string)

	fields[timeZoneOffsetKeyPrefix] = asset.System.TimeZone
	fields[siteIDKeyPrefix] = asset.SiteID
	fields[clientIDKeyPrefix] = asset.ClientID
	fields[machineNameKeyPrefix] = asset.System.SystemName
	fields[resourceTypeKeyPrefix] = asset.EndpointType

	return fields
}

func (assetsService Service) setAssetDataToCache(ctx context.Context, fields map[string]string, endpointID string) {
	for key, val := range fields {
		err := assetsService.mCache.Set(&memcache.Item{
			Key:        key + endpointID,
			Value:      []byte(val),
			Expiration: int32(time.Now().Unix() + int64(config.Config.Memcached.DefaultDataTTLSec)),
		})
		if err != nil {
			logger.Log.WarnfCtx(ctx, "Couldn't set %s for ManagedEndpointID[%s] to the Memcached:", key, endpointID)
		}
	}
	return
}

// parseTimeZoneOffset parses time zone offset and converts it in to time.Duration
func parseTimeZoneOffset(offset string) (duration time.Duration, err error) {
	if len(offset) == 0 {
		return
	}

	ok := isValidTimeZoneOffset(offset)
	if !ok {
		err = errors.New("time zone offset has wrong format")
		return
	}

	operator := string(offset[0])
	hours := offset[1:3]
	minutes := offset[3:]

	return time.ParseDuration(operator + hours + "h" + minutes + "m")
}

func isValidTimeZoneOffset(offset string) bool {
	const (
		hourOffsetMax   = 23
		minuteOffsetMax = 59
	)

	if len(offset) != 5 {
		return false
	}

	operator := string(offset[0])
	hoursStr := offset[1:3]
	minutesStr := offset[3:]

	if operator != "+" && operator != "-" {
		return false
	}

	hours, err := strconv.Atoi(hoursStr)
	if err != nil {
		return false
	}

	minutes, err := strconv.Atoi(minutesStr)
	if err != nil {
		return false
	}

	if hours > hourOffsetMax || minutes > minuteOffsetMax {
		return false
	}

	return true
}
