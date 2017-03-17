package main

import (
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	gorilla "github.com/gorilla/mux"
)

type ServerOptions struct {
	Port              int
	Burst             int
	Concurrency       int
	HttpCacheTtl      int
	HttpReadTimeout   int
	HttpWriteTimeout  int
	CORS              bool
	Gzip              bool
	AuthForwarding    bool
	EnableURLSource   bool
	EnablePlaceholder bool
	Address           string
	PathPrefix        string
	ApiKey            string
	Mount             string
	CertFile          string
	KeyFile           string
	Authorization     string
	Placeholder       string
	PlaceholderImage  []byte
	AlloweOrigins     []*url.URL
	MaxAllowedSize    int
	EnableCeph        bool
	CephConfig        string
	EnableFriendly    bool
	EnableSafeRoute   bool
	SafeKey           string
	EnableTracking    bool
	UseCephBlock      bool
	CephBlockURL      string
	EnableS3          bool
}

func Server(o ServerOptions) error {
	addr := o.Address + ":" + strconv.Itoa(o.Port)
	handler := NewLog(NewServerMux(o), os.Stdout)

	server := &http.Server{
		Addr:           addr,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    time.Duration(o.HttpReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(o.HttpWriteTimeout) * time.Second,
	}

	return listenAndServe(server, o)
}

func listenAndServe(s *http.Server, o ServerOptions) error {
	if o.CertFile != "" && o.KeyFile != "" {
		return s.ListenAndServeTLS(o.CertFile, o.KeyFile)
	}
	return s.ListenAndServe()
}

func joinImageRoute(o ServerOptions, route string) string {
	var middleRoute string
	if o.EnableSafeRoute {
		middleRoute = "/{safehash}"
	}
	if o.EnableCeph {
		middleRoute = middleRoute + "/{service}/{oid}"
	}
	return path.Join(o.PathPrefix, middleRoute, route)
}

func join(o ServerOptions, route string) string {
	return path.Join(o.PathPrefix, route)
}

func joinPublic(o ServerOptions, route string) string {
	return path.Join("/public", route)
}

// NewServerMux creates a new HTTP server route multiplexer.
func NewServerMux(o ServerOptions) http.Handler {
	// mux := http.NewServeMux()
	mux := gorilla.NewRouter()
	mux.Handle(joinPublic(o, "/"), Middleware(indexController, o))
	mux.Handle(joinPublic(o, "/form"), Middleware(formController, o))
	mux.Handle(joinPublic(o, "/health"), Middleware(healthController, o))

	image := ImageMiddleware(o)
	mux.Handle(join(o, "/upload/{service}/{oid}"), image(Info))
	mux.Handle(joinImageRoute(o, "/"), image(Origin))
	mux.Handle(joinImageRoute(o, "/resize"), image(Resize))
	mux.Handle(joinImageRoute(o, "/enlarge"), image(Enlarge))
	mux.Handle(joinImageRoute(o, "/extract"), image(Extract))
	mux.Handle(joinImageRoute(o, "/crop"), image(Crop))
	mux.Handle(joinImageRoute(o, "/rotate"), image(Rotate))
	mux.Handle(joinImageRoute(o, "/flip"), image(Flip))
	mux.Handle(joinImageRoute(o, "/flop"), image(Flop))
	mux.Handle(joinImageRoute(o, "/thumbnail"), image(Thumbnail))
	mux.Handle(joinImageRoute(o, "/zoom"), image(Zoom))
	mux.Handle(joinImageRoute(o, "/convert"), image(Convert))
	mux.Handle(joinImageRoute(o, "/watermark"), image(Watermark))
	mux.Handle(joinImageRoute(o, "/info"), image(Info))

	mux.Handle("/friendly/{service}/{op}/{id}", FriendlyImageMiddleware(o))
	return mux
}
