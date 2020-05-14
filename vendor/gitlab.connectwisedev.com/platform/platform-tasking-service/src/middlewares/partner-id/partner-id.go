package partnerID

import (
	"context"
	"net/http"
	"regexp"
)

// To prevent golint complaint "should not use basic type string as key in context.WithValue"
type key string

const (
	// PartnerIDKey contains a key name which the partner ID is stored in the URL & Context under
	PartnerIDKey key    = "partnerID"
	pattern      string = `/partners/(.*?)/`
)

var re = regexp.MustCompile(pattern)

// Middleware transfers PartnerID from URL to Context
func Middleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	id := re.FindStringSubmatch(r.RequestURI)
	if len(id) == 2 {
		ctx := context.WithValue(r.Context(), PartnerIDKey, id[1])
		next(w, r.WithContext(ctx))
		return
	}
	next(w, r)
}
