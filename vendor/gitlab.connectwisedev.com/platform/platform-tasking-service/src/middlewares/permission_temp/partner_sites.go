package permission_temp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

// Site struct that represents object with site id
type Site struct {
	ID int `json:"siteId"`
}

// outdata is the internal type for retrieving only sites
type outdataSites struct {
	Sites []Site `json:"outdata"`
}

type sites struct {
	SiteList []struct {
		ID int64 `json:"siteId"`
	} `json:"siteDetailList"`
}

const (
	userSitesKeyPattern    = "TKS_SITES_BY_PARTNER_%s_USER_%s"
	partnerSitesKeyPattern = "TKS_SITES_BY_PARTNER_%s"
	partnerSitesURLPattern = "%s/partners/%s/sites?Operation=activesites"
	userSitesURLPattern    = "%s/partner/%s/sites"
	tokenHeader            = `iPlanetDirectoryPro`
)

// PartnerSitesCheck where TransactionIDKeyCTX is added
func (md *Permission) PartnerSitesCheck(ctx context.Context, user user.User) (err error) {
	userSites, err := md.UserSites(ctx, user.PartnerID(), user.Token(), user.UID())
	if err != nil {
		return
	}

	partnerSites, err := md.PartnerSites(ctx, user.PartnerID())
	if err != nil {
		return
	}

	if len(userSites) != len(partnerSites) {
		return fmt.Errorf("number of user sites != partner sites for user %v and partner %v", user.Name(), user.PartnerID())
	}
	return
}

// UserSites gets user sites from cache
func (md *Permission) UserSites(ctx context.Context, partnerID, token string, UID string) (sites []string, err error) {
	if config.Config.AssetCacheEnabled && md.cache != nil {
		sites, err = md.GetFromCacheUserSites(partnerID, UID, md.cache)
		if err != nil {
			sites, err = md.GetUserSites(ctx, partnerID, token)
			if err != nil {
				return sites, err
			}
			md.SetToCacheUserSites(ctx, sites, partnerID, UID, md.cache)
		}
		return
	}

	sites, err = md.GetUserSites(ctx, partnerID, token)
	if err != nil {
		return sites, err
	}
	return sites, nil
}

// PartnerSites gets partner sites from cache
func (md *Permission) PartnerSites(ctx context.Context, partnerID string) (sites []string, err error) {
	if config.Config.AssetCacheEnabled && md.cache != nil {
		sites, err = md.GetFromCachePartnerSites(partnerID, md.cache)
		if err != nil {
			sites, err = md.GetPartnerSites(ctx, partnerID)
			if err != nil {
				return sites, err
			}
			md.SetToCachePartnerSites(ctx, sites, partnerID, md.cache)
		}
		return
	}

	sites, err = md.GetPartnerSites(ctx, partnerID)
	if err != nil {
		return sites, err
	}
	return sites, nil
}

// GetUserSites gets user sites from asset
func (md *Permission) GetUserSites(ctx context.Context, partnerID, token string) ([]string, error) {
	url := fmt.Sprintf(userSitesURLPattern, config.Config.SitesMsURL, partnerID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add(tokenHeader, token)

	resp, err := md.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			md.log.WarnfCtx(ctx,"Closing body err %v", err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	sitesList := sites{}
	if err = json.Unmarshal(body, &sitesList); err != nil {
		return nil, err
	}

	s := make([]string, 0)
	for _, v := range sitesList.SiteList {
		s = append(s, strconv.FormatInt(v.ID, 10))
	}
	return s, nil
}

// GetPartnerSites gets partner sites from asset
func (md *Permission) GetPartnerSites(ctx context.Context, partnerID string) ([]string, error) {
	url := fmt.Sprintf(partnerSitesURLPattern, config.Config.SitesNoTokenURL, partnerID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := md.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			md.log.ErrfCtx(ctx, errorcode.ErrorCantProcessData, "%v", err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	sitesList := outdataSites{}
	if err = json.Unmarshal(body, &sitesList); err != nil {
		return nil, err
	}

	s := make([]string, 0)
	for _, v := range sitesList.Sites {
		s = append(s, strconv.Itoa(v.ID))
	}
	return s, nil
}

// GetFromCacheUserSites gets data from cache
func (md *Permission) GetFromCacheUserSites(partnerID string, UID string, cache persistency.Cache) (sites []string, err error) {
	key := fmt.Sprintf(userSitesKeyPattern, partnerID, UID)
	keyForCache := []byte(key)

	siteBin, err := cache.Get(keyForCache)
	if err != nil {
		return sites, fmt.Errorf("error while getting data from cache. err = %v", err)
	}

	if err = json.Unmarshal(siteBin, &sites); err != nil {
		return sites, fmt.Errorf("error while unmarshaling data to site model %v", err)
	}
	return sites, nil
}

// SetToCacheUserSites sets data to cache
func (md *Permission) SetToCacheUserSites(ctx context.Context, site []string, partnerID string, UID string, cache persistency.Cache) {
	key := fmt.Sprintf(userSitesKeyPattern, partnerID, UID)
	keyForCache := []byte(key)

	sitesBytes, err := json.Marshal(site)
	if err != nil {
		md.log.WarnfCtx(ctx, "couldn't marshal sites for user with partnerID=%s and UID = %v, err:%s", partnerID, UID, err.Error())
		return
	}
	if err = cache.Set(keyForCache, sitesBytes, config.Config.PartnerSitesCacheExpiration); err != nil {
		md.log.WarnfCtx(ctx,"couldn't set site for partnerID=%s and UID=%s, err: %v", partnerID, UID, err.Error())
	}
}

// GetFromCachePartnerSites gets data from cache
func (md *Permission) GetFromCachePartnerSites(partnerID string, cache persistency.Cache) (sites []string, err error) {
	key := fmt.Sprintf(partnerSitesKeyPattern, partnerID)
	keyForCache := []byte(key)

	siteBin, err := cache.Get(keyForCache)
	if err != nil {
		return sites, fmt.Errorf("error while getting data from cache. err = %v", err)
	}

	if err = json.Unmarshal(siteBin, &sites); err != nil {
		return sites, fmt.Errorf("error while unmarshaling data to site model %v", err)
	}
	return sites, nil
}

// SetToCachePartnerSites sets data to cache
func (md *Permission) SetToCachePartnerSites(ctx context.Context, site []string, partnerID string, cache persistency.Cache) {
	key := fmt.Sprintf(partnerSitesKeyPattern, partnerID)
	keyForCache := []byte(key)

	sitesBytes, err := json.Marshal(site)
	if err != nil {
		md.log.ErrfCtx(ctx, errorcode.ErrorCache, "couldn't marshal sites for user with partnerID=%s, err:%s", partnerID, err.Error())
	}

	if err = cache.Set(keyForCache, sitesBytes, config.Config.PartnerSitesCacheExpiration); err != nil {
		md.log.ErrfCtx(ctx, errorcode.ErrorCache, "couldn't set site for partnerID=%s, err: %v", partnerID, err.Error())
	}
}
