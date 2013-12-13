// Spirytus is a toolkit for constructing web servers in conjunction with
// Go's net/http without resorting to fancy tricks like reflection or
// dependency injection.
//
// Spirytus can be hard to swallow and too much of it will make you blind.
package spirytus

import (
	"encoding/json"
	"net/http"
)

// JSONResponse writes a JSON-encoded response with the provided status code to the ResponseWriter.
// If the value cannot be encoded an error is returned and nothing is written to the writer.
func JSONResponse(w http.ResponseWriter, code int, value interface{}) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(v)
	return nil
}

// JSONRequest reads the body of req in to v using a JSON decoder.
func JSONRequest(req *http.Request, v interface{}) error {
	dec := json.NewDecoder(req.Body)
	return dec.Decode(v)
}

// A resource describes an HTTP endpoint that can respond to a set of methods.
// It is a regular http.Handler so can be used with any router.
type Resource struct {
	allow   string
	methods []methodHandler
}

type methodHandler struct {
	method  string
	handler http.Handler
}

// Handle instructs the resource to handle the given method with a handler.
func (r *Resource) Handle(method string, handler http.Handler) {
	h := methodHandler{method, handler}
	for i, m := range r.methods {
		if m.method == method {
			r.methods[i] = h
			return
		}
	}
	r.methods = append(r.methods, h)
	if r.allow != "" {
		method = ", " + method
	}
	r.allow += method
}

func (r *Resource) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r == nil || len(r.methods) == 0 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if req.Method == "OPTIONS" {
		r.serveOptions(w, req)
		return
	}

	for _, m := range r.methods {
		if req.Method == m.method {
			m.handler.ServeHTTP(w, req)
			return
		}
	}
	w.Header().Set("Allow", r.allow)
	http.Error(w, "Not allowed", http.StatusMethodNotAllowed)
}

func (r *Resource) serveOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", r.allow)
	w.WriteHeader(http.StatusOK)
	return
}

func (r *Resource) allowOrigin(req *http.Request) bool {
	return true
}

func (r *Resource) serveCORS(w http.ResponseWriter, req *http.Request) {
	// If the origin is not allowed, continue as normal.
	if origin := req.Header.Get("Origin"); origin != "" && r.allowOrigin(req) {
		r.serveCORS(w, req)
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	if req.Method == "OPTIONS" {

	} else {
	}
}
