package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/daaku/go.httpgzip"
	gorilla "github.com/gorilla/mux"
	"github.com/rs/cors"
	"gopkg.in/h2non/bimg.v1"
	"gopkg.in/throttled/throttled.v2"
	"gopkg.in/throttled/throttled.v2/store/memstore"
)

func Middleware(fn func(http.ResponseWriter, *http.Request), o ServerOptions) http.Handler {
	next := http.Handler(http.HandlerFunc(fn))

	if o.Concurrency > 0 {
		next = throttle(next, o)
	}
	if o.Gzip {
		next = httpgzip.NewHandler(next)
	}
	if o.CORS {
		next = cors.Default().Handler(next)
	}
	if o.ApiKey != "" {
		next = authorizeClient(next, o)
	}
	if o.HttpCacheTtl >= 0 {
		next = setCacheHeaders(next, o.HttpCacheTtl)
	}

	if o.EnableSafeRoute {
		next = checkSafeKey(next, o)
	}

	return validate(defaultHeaders(next), o)
}

func checkSafeKey(next http.Handler, o ServerOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := gorilla.Vars(r)
		safeHash := vars["safehash"]
		route := strings.Replace(r.URL.RequestURI(), "/"+safeHash, "", 1)
		if safeHash == "" || safeHash != hashRoute([]byte(route), []byte(o.SafeKey)) {
			ErrorReply(r, w, ErrSafeHash, o)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// hashRoute to sha1 with safeKey then encoding as hex
func hashRoute(route []byte, safeKey []byte) string {
	mac := hmac.New(sha1.New, safeKey)
	mac.Write(route)
	return hex.EncodeToString(mac.Sum(nil))
}

func ImageMiddleware(o ServerOptions) func(Operation) http.Handler {
	return func(fn Operation) http.Handler {
		return validateImage(Middleware(imageController(o, Operation(fn)), o), o)
	}
}

func CephMiddleware(o ServerOptions) func(Operation) http.Handler {
	return func(fn Operation) http.Handler {
		return validateImage(Middleware(cephController(o, Operation(fn)), o), o)
	}
}

func FriendlyImageMiddleware(o ServerOptions) http.Handler {
	return validateImage(friendlyRoute(o), o)
}

func friendlyRoute(o ServerOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if o.EnableFriendly {
			vars := gorilla.Vars(r)
			sq := ServiceQueryMap[vars["service"]]               // TODO : handle error
			r.URL.RawQuery = sq.getQuery(vars["op"], vars["id"]) // TODO : handle error
			Middleware(imageController(o, sq.getOperation(vars["op"])), o).ServeHTTP(w, r)
		} else {
			ErrorReply(r, w, ErrFriendlyNotAllowed, o)
		}
	})
}

func throttleError(err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "throttle error: "+err.Error(), http.StatusInternalServerError)
	})
}

func throttle(next http.Handler, o ServerOptions) http.Handler {
	store, err := memstore.New(65536)
	if err != nil {
		return throttleError(err)
	}

	quota := throttled.RateQuota{throttled.PerSec(o.Concurrency), o.Burst}
	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		return throttleError(err)
	}

	httpRateLimiter := throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{Method: true},
	}

	return httpRateLimiter.RateLimit(next)
}

func validate(next http.Handler, o ServerOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "POST" {
			ErrorReply(r, w, ErrMethodNotAllowed, o)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateImage(next http.Handler, o ServerOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if r.Method == "GET" && isPublicPath(path) {
			next.ServeHTTP(w, r)
			return
		}

		if r.Method == "GET" && o.Mount == "" && o.EnableURLSource == false {
			ErrorReply(r, w, ErrMethodNotAllowed, o)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func authorizeClient(next http.Handler, o ServerOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("API-Key")
		if key == "" {
			key = r.URL.Query().Get("key")
		}

		if key != o.ApiKey {
			ErrorReply(r, w, ErrInvalidApiKey, o)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func defaultHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", fmt.Sprintf("imaginary %s (bimg %s)", Version, bimg.Version))
		next.ServeHTTP(w, r)
	})
}

func setCacheHeaders(next http.Handler, ttl int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer next.ServeHTTP(w, r)

		if r.Method != "GET" || isPublicPath(r.URL.Path) {
			return
		}

		ttlDiff := time.Duration(ttl) * time.Second
		expires := time.Now().Add(ttlDiff)

		w.Header().Add("Expires", strings.Replace(expires.Format(time.RFC1123), "UTC", "GMT", -1))
		w.Header().Add("Cache-Control", getCacheControl(ttl))
	})
}

func getCacheControl(ttl int) string {
	if ttl == 0 {
		return "private, no-cache, no-store, must-revalidate"
	}
	return fmt.Sprintf("public, s-maxage=%d, max-age=%d, no-transform", ttl, ttl)
}

func isPublicPath(path string) bool {
	return path == "/" || path == "/health" || path == "/form"
}
