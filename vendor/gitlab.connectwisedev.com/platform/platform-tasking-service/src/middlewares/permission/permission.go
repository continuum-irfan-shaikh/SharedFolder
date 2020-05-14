package permission

//go:generate mockgen -destination=./user_uc_mock_test.go -package=permission -source=./permission.go

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	en "gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

type userUC interface {
	Sites(ctx context.Context) ([]string, error)
	Endpoints(ctx context.Context, sites []string) ([]string, error)
}

// NewPermission - returns a new instance of Permission middleware
func NewPermission(permissionUC userUC, log logger.Logger) *Permission {
	return &Permission{
		userUC: permissionUC,
		log:    log,
	}
}

// Permission - represents User middleware for putting Permission filters in context
type Permission struct {
	userUC userUC
	log    logger.Logger
}

// ServeHTTP - retrieves User Permission filters and sets them in context
func (p *Permission) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := r.Context()
	user, err := p.user(ctx)
	if err != nil {
		p.log.ErrfCtx(r.Context(), errorcode.ErrorCantGetUser, err.Error())
		common.SendBadRequest(rw, r, errorcode.ErrorCantGetUser)
		return
	}

	if user.IsNOCAccess {
		next.ServeHTTP(rw, r)
		return
	}

	sites, err := p.userUC.Sites(ctx)
	if err != nil {
		p.log.ErrfCtx(r.Context(), errorcode.ErrorCantGetUserSites, err.Error())
		common.SendBadRequest(rw, r, errorcode.ErrorCantGetUserSites)
		return
	}

	//ed - user endpoints
	ed, err := p.userUC.Endpoints(ctx, sites)
	if err != nil {
		p.log.ErrfCtx(r.Context(), errorcode.ErrorCantGetUserEndPointsBySites, err.Error())
		common.SendBadRequest(rw, r, errorcode.ErrorCantGetUserEndPointsBySites)
		return
	}

	ctx = context.WithValue(ctx, config.UserEndPointsKeyCTX, ed)
	r = r.WithContext(ctx)

	next.ServeHTTP(rw, r)
}

// user - get user from context
func (p *Permission) user(ctx context.Context) (u en.User, err error) {
	user, ok := ctx.Value(config.UserKeyCTX).(en.User)
	if !ok {
		return u, fmt.Errorf("PermissionMD.user: can't get user")
	}
	return user, nil
}
