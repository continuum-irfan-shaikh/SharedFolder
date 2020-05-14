package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	cweb "gitlab.connectwisedev.com/platform/platform-common-lib/src/web"
)

//RequestContextKey Type to handle request context keys
type requestContextKey string

const (
	//RequestContextMappingKey is a key use in context request to store the value
	RequestContextMappingKey requestContextKey = "RequestContextMappingKey"
)

// AuthMiddleware is a middleware it will do some pre processing of the request
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxWithUser := context.WithValue(r.Context(), RequestContextMappingKey, "endpointmap")
		next.ServeHTTP(w, r.WithContext(ctxWithUser))
	})
}

// HandlerFuncWhenMiddleware will calll when middleware execution will complet
func HandlerFuncWhenMiddleware(w http.ResponseWriter, r *http.Request) {
	mapping := r.Context().Value(RequestContextMappingKey)
	fmt.Printf("Value from context%+v", mapping)
}

// HandlerFuncWhenwithoutmiddlware it will call when request arrived
func HandlerFuncWhenwithoutmiddlware(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handlerFuncWhenwithoutmiddlware invoked")
}

// Registering route with middleware
func addRouteWithMiddleware(r1 cweb.Router) {
	r1.AddHandle("/path-to-api", AuthMiddleware(HandlerFuncWhenMiddleware), http.MethodGet)
}

// Registering route without middleware
func addRouteWithoutMiddleware(r2 cweb.Router) {
	r2.AddFunc("/apth-to-api", HandlerFuncWhenwithoutmiddlware, http.MethodGet)
}

type testMiddleware struct {
	next http.Handler
}

func (h *testMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Middleware function invoked")
	content := []byte("test")
	w.Write(content)
	h.next.ServeHTTP(w, r)
}

func (h *testMiddleware) dummyMiddleware(handler http.Handler) http.Handler {
	return &testMiddleware{next: handler}
}

// Add middleware to router
func addMiddlewareToRouter(r cweb.Router) {
	tm := &testMiddleware{}
	//add variadic middleware
	r.Use(tm.dummyMiddleware, tm.dummyMiddleware, tm.dummyMiddleware)
}

// getting Server Config Object
var getServerConfig = func() *cweb.ServerConfig {
	return &cweb.ServerConfig{ListenURL: ":8080"}
}

// getting Server Object
var createServer = func(cfg *cweb.ServerConfig) cweb.HTTPServer {
	return cweb.Create(cfg)
}

func main() {
	// getting Server Object
	cfg := getServerConfig()
	// getting Server Config Object
	h := createServer(cfg)
	// getting associated router
	r1 := h.GetRouter()
	addRouteWithMiddleware(r1)
	addRouteWithoutMiddleware(r1)
	addMiddlewareToRouter(r1)
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	log.Fatal(h.Start(ctx))
}
