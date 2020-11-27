package webdav

import (
	"context"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/webdav"

	wdc "webdav/config"
)

// Config is the configuration of a WebDAV instance.
type Config struct {
	wdc.DavServer
	handler *webdav.Handler
}

func checkPassword(saved, input string) bool {
	if strings.HasPrefix(saved, "{bcrypt}") {
		savedPassword := strings.TrimPrefix(saved, "{bcrypt}")
		return bcrypt.CompareHashAndPassword([]byte(savedPassword), []byte(input)) == nil
	}

	return saved == input
}

// ServeHTTP determines if the request is for this plugin, and if all prerequisites are met.
func (c *Config) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

	// 未提供用户
	username, password, ok := r.BasicAuth()
	if !ok {
		log.Println("未提供用户:", r.RemoteAddr)
		http.Error(w, "Not authorized", 401)
		return
	}

	if c.Auth {
		if username != c.User || !checkPassword(c.Pass, password) {
			http.Error(w, "Not authorized", 401)
			log.Printf("%s: %s(%s) not authorized.\n", r.RemoteAddr, username, password)
		}
	}

	if r.Method == "HEAD" {
		w = newResponseWriterNoBody(w)
	}

	// If this request modified the files and the user doesn't have permission
	// to do so, return forbidden.
	if (r.Method == "PUT" || r.Method == "POST" || r.Method == "MKCOL" ||
		r.Method == "DELETE" || r.Method == "COPY" || r.Method == "MOVE") &&
		!c.Modify {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Excerpt from RFC4918, section 9.4:
	//
	// 		GET, when applied to a collection, may return the contents of an
	//		"index.html" resource, a human-readable view of the contents of
	//		the collection, or something else altogether.
	//
	// Get, when applied to collection, will return the same as PROPFIND method.
	if r.Method == "GET" {
		info, err := c.handler.FileSystem.Stat(context.TODO(), r.URL.Path)
		if err == nil && info.IsDir() {
			r.Method = "PROPFIND"
			if r.Header.Get("Depth") == "" {
				r.Header.Add("Depth", "1")
			}
		}
	}

	log.Println(r.RemoteAddr, r.URL.Path)
	c.handler.ServeHTTP(w, r)
}

// responseWriterNoBody is a wrapper used to suprress the body of the response
// to a request. Mainly used for HEAD requests.
type responseWriterNoBody struct {
	http.ResponseWriter
}

// newResponseWriterNoBody creates a new responseWriterNoBody.
func newResponseWriterNoBody(w http.ResponseWriter) *responseWriterNoBody {
	return &responseWriterNoBody{w}
}

// Header executes the Header method from the http.ResponseWriter.
func (w responseWriterNoBody) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Write suprresses the body.
func (w responseWriterNoBody) Write(_ []byte) (int, error) {
	return 0, nil
}

// WriteHeader writes the header to the http.ResponseWriter.
func (w responseWriterNoBody) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}
