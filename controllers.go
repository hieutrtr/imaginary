package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"gopkg.in/h2non/bimg.v1"
	"gopkg.in/h2non/filetype.v0"

	gorilla "github.com/gorilla/mux"
)

const (
	UPLOAD int = iota
	CACHE
	PROCESS
)

var cachedAttributes = []string{
	"thumbnail",
	"watermark",
}

func indexController(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorReply(r, w, ErrNotFound, ServerOptions{})
		return
	}

	body, _ := json.Marshal(CurrentVersions)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func healthController(w http.ResponseWriter, r *http.Request) {
	health := GetHealthStats()
	body, _ := json.Marshal(health)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func imageController(o ServerOptions, operation Operation) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		getCachedImage := func() []byte {
			return nil
		}

		getImage := func() []byte {
			var imageSource = MatchSource(req)
			if imageSource == nil {
				ErrorReply(req, w, ErrMissingImageSource, o)
				return nil
			}

			buf, err := imageSource.GetImage(req)
			if err != nil {
				return nil
			}

			if len(buf) == 0 {
				ErrorReply(req, w, ErrEmptyBody, o)
				return nil
			}
			setResponseHeader(w, req)
			return buf
		}

		vars := gorilla.Vars(req)
		switch routingRequest(req) {
		case UPLOAD:
			if buf := getImage(); buf != nil {
				if err := uploadImage(req, buf); err != nil {
					ErrorReply(req, w, NewError(err.Error(), InternalError), o)
				}
				imageHandler(w, req, buf, Info, o)
			}
		case CACHE:
			if buf := getCachedImage(); buf != nil {
				imageHandler(w, req, buf, Origin, o)
			} else {
				if buf := getImage(); buf != nil {
					imageHandler(w, req, buf, operation, o)
				} else {
					var err error
					if buf, err = getDefault(vars["service"]); err != nil {
						ErrorReply(req, w, NewError(err.Error(), BadRequest), o)
					} else {
						imageHandler(w, req, buf, Origin, o)
					}
				}
			}
		case PROCESS:
			if buf := getImage(); buf != nil {
				imageHandler(w, req, buf, operation, o)
			} else {
				var err error
				if buf, err = getDefault(vars["service"]); err != nil {
					ErrorReply(req, w, NewError(err.Error(), BadRequest), o)
				} else {
					imageHandler(w, req, buf, Origin, o)
				}
			}
		}
	}
}

func getDefault(service string) ([]byte, error) {
	switch service {
	case "profile_avatar":
		return ioutil.ReadFile("./fixtures/default_avatar.png")
	default:
		return nil, NewError("Ceph: Image is not exist", InternalError)
	}
}

func routingRequest(req *http.Request) int {
	if IsUploadRequest(req) {
		return UPLOAD
	}
	attr := getCacheAttr(req.URL.Path, req.URL.RawQuery)
	if attr == DATA {
		return PROCESS
	}
	return CACHE
}

// IsUploadRequest check if request is for uploading an image
func IsUploadRequest(r *http.Request) bool {
	if r.Method == "POST" && strings.HasPrefix(r.URL.RequestURI(), "/upload/") {
		return true
	}
	return false
}

// IsPublic check if request is for uploading an image
func IsPublic(r *http.Request) bool {
	if r.Method == "GET" && strings.HasPrefix(r.URL.RequestURI(), "/public/") {
		return true
	}
	return false
}

func getCacheAttr(urlPath, rawQuery string) string {
	parts := strings.Split(urlPath, "/")
	for _, a := range cachedAttributes {
		if a == parts[len(parts)-1] {
			attr := fmt.Sprintf("%s_%s", a, rawQuery)
			return attr
		}
	}
	return DATA
}

func uploadImage(req *http.Request, buf []byte) error {
	var connection = MatchConnection(req)
	if connection == nil {
		return ErrMissingConnection
	}
	return connection.Execute(req, buf)
}

func checkSupportedMediaType(buf []byte) bool {
	// Infer the body MIME type via mimesniff algorithm
	mimeType := http.DetectContentType(buf)
	// If cannot infer the type, infer it via magic numbers
	if mimeType == "application/octet-stream" {
		kind, err := filetype.Get(buf)
		if err == nil && kind.MIME.Value != "" {
			mimeType = kind.MIME.Value
		}
	}

	// Infer text/plain responses as potential SVG image
	if strings.Contains(mimeType, "text/plain") && len(buf) > 8 {
		if bimg.IsSVGImage(buf) {
			mimeType = "image/svg+xml"
		}
	}
	return IsImageMimeTypeSupported(mimeType)
}

func setResponseHeader(w http.ResponseWriter, req *http.Request) {
	if modtime := req.Header.Get("last-modified"); modtime != "" {
		w.Header().Set("last-modified", modtime)
		h := sha256.New()
		h.Write([]byte(modtime))
		if etag := fmt.Sprintf("%x", h.Sum(nil)); etag != "" {
			w.Header().Set("Etag", etag)
		}
	}
}

func imageHandler(w http.ResponseWriter, r *http.Request, buf []byte, operation Operation, o ServerOptions) []byte {
	// Infer the body MIME type via mimesniff algorithm
	mimeType := http.DetectContentType(buf)
	// If cannot infer the type, infer it via magic numbers
	if mimeType == "application/octet-stream" {
		kind, err := filetype.Get(buf)
		if err == nil && kind.MIME.Value != "" {
			mimeType = kind.MIME.Value
		}
	}

	// Infer text/plain responses as potential SVG image
	if strings.Contains(mimeType, "text/plain") && len(buf) > 8 {
		if bimg.IsSVGImage(buf) {
			mimeType = "image/svg+xml"
		}
	}
	// Finally check if image MIME type is supported
	if IsImageMimeTypeSupported(mimeType) == false {
		ErrorReply(r, w, ErrUnsupportedMedia, o)
		return nil
	}

	opts := readParams(r.URL.Query())
	if opts.Type != "" && ImageType(opts.Type) == 0 {
		ErrorReply(r, w, ErrOutputFormat, o)
		return nil
	}

	image, err := operation.Run(buf, opts)
	if err != nil {
		ErrorReply(r, w, NewError("Error while processing the image: "+err.Error(), BadRequest), o)
		return nil
	}
	w.Header().Set("Content-Type", image.Mime)
	w.Write(image.Body)
	if image.Mime != "application/json" {
		return image.Body
	}
	return nil
}

func formController(w http.ResponseWriter, r *http.Request) {
	operations := []struct {
		name   string
		method string
		args   string
	}{
		{"Resize", "resize", "width=300&height=200&type=jpeg"},
		{"Force resize", "resize", "width=300&height=200&force=true"},
		{"Crop", "crop", "width=300&quality=95"},
		{"Extract", "extract", "top=100&left=100&areawidth=300&areaheight=150"},
		{"Enlarge", "enlarge", "width=1440&height=900&quality=95"},
		{"Rotate", "rotate", "rotate=180"},
		{"Flip", "flip", ""},
		{"Flop", "flop", ""},
		{"Thumbnail", "thumbnail", "width=100"},
		{"Zoom", "zoom", "factor=2&areawidth=300&top=80&left=80"},
		{"Color space (black&white)", "resize", "width=400&height=300&colorspace=bw"},
		{"Add watermark", "watermark", "textwidth=100&text=Hello&font=sans%2012&opacity=0.5&color=255,200,50"},
		{"Convert format", "convert", "type=png"},
		{"Image metadata", "info", ""},
		{"Upload Image", "upload", "cns=ads&cid=1269"},
	}

	html := "<html><body>"

	for _, form := range operations {
		html += fmt.Sprintf(`
    <h1>%s</h1>
    <form method="POST" action="/%s?%s" enctype="multipart/form-data">
      <input type="file" name="file" />
      <input type="submit" value="Upload" />
    </form>`, form.name, form.method, form.args)
	}

	html += "</body></html>"

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
