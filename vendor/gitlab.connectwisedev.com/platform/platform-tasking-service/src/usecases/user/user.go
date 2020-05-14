package user

//go:generate mockgen -destination=./user_repo_mock_test.go -package=usecases -source=./user.go

import (
	"context"
	"fmt"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	en "gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
)

type userRepoDB interface {
	Endpoints(ctx context.Context, partnerID string, siteIDs []string) ([]en.Endpoints, error)
	SaveEndpoints(ctx context.Context, ep []en.Endpoints)
}

type userRepoHTTP interface {
	Endpoints(siteIDs []string, partnerID string) ([]en.Endpoints, error)
}

type siteRepoHTTP interface {
	Sites(partnerID, token string) ([]string, error)
}

// NewUser - creates a new instance of a User use case
func NewUser(d userRepoDB, u userRepoHTTP, s siteRepoHTTP) *user {
	return &user{
		repoDB:       d,
		userRepoHTTP: u,
		siteRepoHTTP: s,
	}
}

type user struct {
	repoDB       userRepoDB
	userRepoHTTP userRepoHTTP
	siteRepoHTTP siteRepoHTTP
}

// Sites - returns user sites by partnerID
func (u *user) Sites(ctx context.Context) ([]string, error) {
	user, err := u.user(ctx)
	if err != nil {
		return nil, err
	}

	//us - user sites
	us, err := u.siteRepoHTTP.Sites(user.PartnerID, user.Token)
	if err != nil {
		msg := "UserUC.Endpoints: can't get user sites from siteRepoHTTP: %v"
		return nil, fmt.Errorf(msg, err)
	}

	return us, nil
}

// Endpoints - returns user End Points by Endpoints from DB or Asset service
func (u *user) Endpoints(ctx context.Context, sites []string) (ep []string, err error) {
	user, err := u.user(ctx)
	if err != nil {
		return nil, err
	}

	endpoints, err := u.repoDB.Endpoints(ctx, user.PartnerID, sites)
	if err != nil {
		msg := "UserUC.Endpoints: can't get endPoints by sites from repoDB: %v"
		return nil, fmt.Errorf(msg, err)
	}

	if len(sites) == len(endpoints) {
		for _, v := range endpoints {
			ep = append(ep, v.Endpoints...)
		}
		return ep, nil
	}

	endpoints, err = u.EndpointsFromAsset(ctx, sites)
	if err != nil {
		return nil, err
	}

	go u.SaveEndpoints(ctx, endpoints)

	for _, v := range endpoints {
		ep = append(ep, v.Endpoints...)
	}

	return ep, nil
}

// EndpointsFromAsset - returns user End Points by Endpoints from Asset service
func (u *user) EndpointsFromAsset(ctx context.Context, sites []string) ([]en.Endpoints, error) {
	user, err := u.user(ctx)
	if err != nil {
		return nil, err
	}

	ep, err := u.userRepoHTTP.Endpoints(sites, user.PartnerID)
	if err != nil {
		msg := "UserUC.EndpointsFromAsset: can't get endPoints by sites from userRepoHTTP: %v"
		return nil, fmt.Errorf(msg, err)
	}
	return ep, nil
}

// SaveEndpoints - save endpoints by partnerID and sites
func (u *user) SaveEndpoints(ctx context.Context, ep []en.Endpoints) {
	u.repoDB.SaveEndpoints(ctx, ep)
}

// user - get user from context, validate user's parameters
func (u *user) user(ctx context.Context) (user en.User, err error) {
	user, ok := ctx.Value(config.UserKeyCTX).(en.User)
	if !ok {
		return user, fmt.Errorf("UserUC.user: can't get user")
	}
	return user, nil
}
