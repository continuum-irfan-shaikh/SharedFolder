package permission_temp

import (
	"fmt"
	"net/http"

	"strings"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
)

// roles according to task https://continuum.atlassian.net/browse/RMM-40019
const (
	RoleSuperUser        = "SuperUser"
	RoleAdmin            = "Admin"
	RolePrimarySuperUser = "PrimarySuperUser"
	RoleManager          = "Manager"
	RoleTechnician       = "Technician"
)

// Permission represents permission middleware
type Permission struct {
	log         logger.Logger
	userService user.Service
	httpClient  *http.Client
	cache       persistency.Cache
}

// NewPermissionMiddleware returns new Permission
func NewPermissionMiddleware(logger logger.Logger, userService user.Service, httpClient *http.Client, cache persistency.Cache) *Permission {
	return &Permission{
		log:         logger,
		userService: userService,
		httpClient:  httpClient,
		cache:       cache,
	}
}

// ServeHTTP where TransactionID is added
func (md *Permission) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	u := md.userService.GetUser(r, md.httpClient)
	if u.HasNOCAccess() || u.UID() == "" || u.Token() == "" {
		next(rw, r)
		return
	}

	role := strings.TrimSpace(r.Header.Get("role"))
	switch role {
	case RoleSuperUser, RoleManager, RoleAdmin, RolePrimarySuperUser, RoleTechnician:
		next(rw, r)
		return
	}

	if err := md.PartnerSitesCheck(r.Context(), u); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		md.log.WarnfCtx(r.Context(),"Forbidden request: role - %q", role)
		common.SendForbidden(rw, r, fmt.Sprintf("Forbidden request: role - %q", role))
		return
	}
	next(rw, r)
}
