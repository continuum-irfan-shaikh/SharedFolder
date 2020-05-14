package user

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

const defaultMsg = `failed on unexpected value of result "%v"`

func TestNewUser(t *testing.T) {
	RegisterTestingT(t)

	expected := &user{}
	actual := NewUser(nil, nil, nil)
	Ω(actual).To(Equal(expected), fmt.Sprintf(defaultMsg, expected))
}

func TestUser_Sites(t *testing.T) {
	// Register gomega
	RegisterTestingT(t)
	mockCtrl := gomock.NewController(t)

	user := entities.User{
		PartnerID: "PartnerID",
		Token:     "Token",
	}

	userCTX := context.Background()
	userCTX = context.WithValue(userCTX, config.UserKeyCTX, user)

	type expected struct {
		result []string
		err    string
	}

	tc := []struct {
		name     string
		payload  func() (context.Context, siteRepoHTTP)
		expected expected
	}{
		{
			name: "error: can't get user from the context",
			payload: func() (i context.Context, db siteRepoHTTP) {
				return context.Background(), nil
			},
			expected: expected{
				err: "UserUC.user: can't get user",
			},
		},
		{
			name: "error: can't get sites for user",
			payload: func() (i context.Context, db siteRepoHTTP) {
				repoMock := NewMocksiteRepoHTTP(mockCtrl)
				err := errors.New("can't get user sites")
				repoMock.EXPECT().Sites("PartnerID", "Token").Return(nil, err).Times(1)
				return userCTX, repoMock
			},
			expected: expected{
				err: "UserUC.Endpoints: can't get user sites from siteRepoHTTP: can't get user sites",
			},
		},
		{
			name: "success: returns user sites",
			payload: func() (i context.Context, db siteRepoHTTP) {
				repoMock := NewMocksiteRepoHTTP(mockCtrl)
				sites := []string{"1", "2", "3"}
				repoMock.EXPECT().Sites("PartnerID", "Token").Return(sites, nil).Times(1)
				return userCTX, repoMock
			},
			expected: expected{
				result: []string{"1", "2", "3"},
			},
		},
	}

	for _, test := range tc {
		ctx, repo := test.payload()
		user := NewUser(nil, nil, repo)

		actual, err := user.Sites(ctx)

		if test.expected.err != "" {
			Ω(err.Error()).To(Equal(test.expected.err), fmt.Sprintf(defaultMsg, test.name))
			continue
		}
		Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, test.name))
		Ω(actual).To(Equal(test.expected.result), fmt.Sprintf(defaultMsg, test.name))
	}
}

func TestMockuserRepoDB_Endpoints(t *testing.T) {
	// Register gomega
	RegisterTestingT(t)
	mockCtrl := gomock.NewController(t)

	user := entities.User{
		PartnerID: "PartnerID",
	}

	userCTX := context.Background()
	userCTX = context.WithValue(userCTX, config.UserKeyCTX, user)

	siteIDs := []string{"1", "2", "3"}

	type expected struct {
		result []string
		err    string
	}

	type payload struct {
		ctx       context.Context
		repoDB    func() userRepoDB
		repoAsset func() userRepoHTTP
	}

	tc := []struct {
		name     string
		payload  payload
		expected expected
	}{
		{
			name: "error: can't get user from the context",
			payload: payload{
				ctx:       context.Background(),
				repoDB:    func() userRepoDB { return nil },
				repoAsset: func() userRepoHTTP { return nil },
			},
			expected: expected{
				err: "UserUC.user: can't get user",
			},
		},
		{
			name: "error: can't get endpoints from the DB repo",
			payload: payload{
				ctx: userCTX,
				repoDB: func() userRepoDB {
					repoMock := NewMockuserRepoDB(mockCtrl)
					err := errors.New("can't get endpoints")
					repoMock.EXPECT().Endpoints(userCTX, "PartnerID", siteIDs).Return(nil, err).Times(1)
					return repoMock
				},
				repoAsset: func() userRepoHTTP { return nil },
			},
			expected: expected{
				err: "UserUC.Endpoints: can't get endPoints by sites from repoDB: can't get endpoints",
			},
		},
		{
			name: "success: returns endpoints from repo DB",
			payload: payload{
				ctx: userCTX,
				repoDB: func() userRepoDB {
					result := []entities.Endpoints{
						{
							SiteID:    "1",
							Endpoints: []string{"1"},
						},
						{
							SiteID:    "2",
							Endpoints: []string{"2"},
						},
						{
							SiteID:    "3",
							Endpoints: []string{"3"},
						},
					}
					repoMock := NewMockuserRepoDB(mockCtrl)
					repoMock.EXPECT().Endpoints(userCTX, "PartnerID", siteIDs).Return(result, nil).Times(1)
					return repoMock
				},
				repoAsset: func() userRepoHTTP { return nil },
			},
			expected: expected{
				result: []string{"1", "2", "3"},
			},
		},
		{
			name: "error: can't get endpoints from the asset repo",
			payload: payload{
				ctx: userCTX,
				repoDB: func() userRepoDB {
					result := []entities.Endpoints{}
					repoMock := NewMockuserRepoDB(mockCtrl)
					repoMock.EXPECT().Endpoints(userCTX, "PartnerID", siteIDs).Return(result, nil).Times(1)
					return repoMock
				},
				repoAsset: func() userRepoHTTP {
					err := errors.New("cat't get endpoints")
					repoMock := NewMockuserRepoHTTP(mockCtrl)
					repoMock.EXPECT().Endpoints(siteIDs, "PartnerID").Return(nil, err).Times(1)
					return repoMock
				},
			},
			expected: expected{
				err: "UserUC.EndpointsFromAsset: can't get endPoints by sites from userRepoHTTP: cat't get endpoints",
			},
		},
		{
			name: "success: returns endpoints from asset DB",
			payload: payload{
				ctx: userCTX,
				repoDB: func() userRepoDB {
					result := []entities.Endpoints{}
					repoMock := NewMockuserRepoDB(mockCtrl)
					repoMock.EXPECT().Endpoints(userCTX, "PartnerID", siteIDs).Return(result, nil).Times(1)
					repoMock.EXPECT().SaveEndpoints(userCTX, gomock.Any()).Times(1)
					return repoMock
				},
				repoAsset: func() userRepoHTTP {
					result := []entities.Endpoints{
						{
							SiteID:    "1",
							Endpoints: []string{"1"},
						},
					}
					repoMock := NewMockuserRepoHTTP(mockCtrl)
					repoMock.EXPECT().Endpoints(siteIDs, "PartnerID").Return(result, nil).Times(1)
					return repoMock
				},
			},
			expected: expected{
				result: []string{"1"},
			},
		},
	}

	for _, test := range tc {
		user := NewUser(test.payload.repoDB(), test.payload.repoAsset(), nil)

		actual, err := user.Endpoints(test.payload.ctx, siteIDs)

		if test.expected.err != "" {
			Ω(err.Error()).To(Equal(test.expected.err), fmt.Sprintf(defaultMsg, test.name))
			continue
		}
		Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, test.name))
		Ω(actual).To(Equal(test.expected.result), fmt.Sprintf(defaultMsg, test.name))
	}
}

func TestUser_EndpointsFromAsset(t *testing.T) {
	// Register gomega
	RegisterTestingT(t)
	siteIDs := []string{"1", "2", "3"}

	type expected struct {
		result []string
		err    string
	}

	type payload struct {
		ctx       context.Context
		repoAsset func() userRepoHTTP
	}

	tc := []struct {
		name     string
		payload  payload
		expected expected
	}{
		{
			name: "error: can't get user from the context",
			payload: payload{
				ctx:       context.Background(),
				repoAsset: func() userRepoHTTP { return nil },
			},
			expected: expected{
				err: "UserUC.user: can't get user",
			},
		},
	}

	for _, test := range tc {
		user := NewUser(nil, test.payload.repoAsset(), nil)

		actual, err := user.EndpointsFromAsset(test.payload.ctx, siteIDs)

		if test.expected.err != "" {
			Ω(err.Error()).To(Equal(test.expected.err), fmt.Sprintf(defaultMsg, test.name))
			continue
		}
		Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, test.name))
		Ω(actual).To(Equal(test.expected.result), fmt.Sprintf(defaultMsg, test.name))
	}
}

func TestUser_SaveEndpoints(t *testing.T) {
	ctx := context.Background()
	ep := make([]entities.Endpoints, 1)

	mockCtrl := gomock.NewController(t)
	repoMock := NewMockuserRepoDB(mockCtrl)
	repoMock.EXPECT().SaveEndpoints(ctx, ep)

	user := NewUser(repoMock, nil, nil)
	user.SaveEndpoints(ctx, ep)
}
