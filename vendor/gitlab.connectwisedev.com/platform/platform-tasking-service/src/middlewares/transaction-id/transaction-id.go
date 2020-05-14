package transactionID

import (
	"context"
	"net/http"

	"github.com/gocql/gocql"
)

const Key = "TransactionID"

//Middleware add Key value request context to have possibility identify operations later
func Middleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// do some stuff before
	ctx := newContextTransactionIDWithReq(r.Context(), r)
	next(rw, r.WithContext(ctx))
}

// newContextWithRequestID created new Context with inserted X-Request-ID value
func newContextTransactionIDWithReq(ctx context.Context, r *http.Request) context.Context {
	var rID = r.Header.Get(Key)
	if rID == "" {
		rID = generateTransactionID()
	}
	r.Header.Add(Key, rID)
	return context.WithValue(ctx, Key, rID)
}

// generateRequestID generates random Request ID
func generateTransactionID() string {
	uuid, err := gocql.RandomUUID()
	if err != nil {
		uuid = gocql.TimeUUID()
	}
	return uuid.String()
}

// FromContext extracts X-Request-ID value from HTTP Request' Context
func FromContext(ctx context.Context) (str string) {
	if ctx == nil {
		return ""
	}

	str, _ = ctx.Value(Key).(string)
	return
}

// NewContext returns new context that contains newly generated transaction ID
func NewContext() context.Context {
	return context.WithValue(context.Background(), Key, generateTransactionID())
}
