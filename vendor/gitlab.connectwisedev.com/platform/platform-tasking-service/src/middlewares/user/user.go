package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

// NewUser returns new users midlware
func NewUser(log logger.Logger) *User {
	return &User{log: log}
}

// User - represents User middleware for putting user parameters into context
type User struct {
	log logger.Logger
}

// ServeHTTP - retrieves user parameters and sets them in context
func (u *User) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token := r.Header.Get("iPlanetDirectoryProPlus")
	if len(token) == 0 {
		token = r.Header.Get("iPlanetDirectoryPro")
	}

	var isNOCAccess bool
	if realm := r.Header.Get("realm"); realm == "/activedirectory" {
		isNOCAccess = true
	}

	r.URL.Query()

	user := entities.User{
		Name:        r.Header.Get("username"),
		UID:         r.Header.Get("uid"),
		PartnerID:   mux.Vars(r)["partnerID"],
		IsNOCAccess: isNOCAccess,
		Token:       token,
	}

	if err := u.validateUser(user); err != nil {
		u.log.ErrfCtx(r.Context(), errorcode.ErrorCantValidateUserParameters, "User middleware: Can't validate user %v", err)
		common.SendBadRequest(rw, r, errorcode.ErrorCantValidateUserParameters)
		return
	}

	ctx := context.WithValue(r.Context(), config.UserKeyCTX, user)
	r = r.WithContext(ctx)
	next.ServeHTTP(rw, r)
}

func (u *User) validateUser(user entities.User) error {
	switch {
	case user.PartnerID == "":
		return fmt.Errorf("middelware: User.validateUser: PartnerID can't be empty")
	case user.UID == "":
		return fmt.Errorf("middelware: User.validateUser: UID can't be empty")
	}
	return nil
}
