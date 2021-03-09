package web

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/fluxynet/goqa"
)

const (
	// ContentTypeJSON is the content type for JSON
	ContentTypeJSON = "application/json"

	// ContentTypeHTML is the content type for HTML
	ContentTypeHTML = "text/html"

	// ContentTypeEventStream used for SSE
	ContentTypeEventStream = "text/event-stream"
)

var (
	// ErrInvalidRequest means a request is either nil or not appropriate for the requested action
	ErrInvalidRequest = errors.New("request is invalid")

	// ErrResourceNotFound is when we did not get what they wanted us to get
	ErrResourceNotFound = errors.New("resource not found")

	// ErrStreamingNotSupported SSE not supported by client
	ErrStreamingNotSupported = errors.New("client does not support streaming")

	// ErrPayloadUnverified payload could not be verified wrt signature
	ErrPayloadUnverified = errors.New("payload could not be verified")
)

// Send data to the browser
func Print(w http.ResponseWriter, status int, ctype string, content []byte) {
	w.Header().Set("Content-Type", ctype)
	w.WriteHeader(status)
	w.Write(content)
}

// Json to the browser
func Json(w http.ResponseWriter, r interface{}) {
	var b, err = json.Marshal(r)
	if err == nil {
		Print(w, http.StatusOK, ContentTypeJSON, b)
	} else {
		JsonError(w, http.StatusInternalServerError, err)
	}
}

// JsonError to the browser in json format
func JsonError(w http.ResponseWriter, status int, err error) {
	var m = strings.ReplaceAll(err.Error(), `"`, `\"`)
	Print(w, status, ContentTypeJSON, []byte(`{"error":"`+m+`"}`))
}

// VerifyBody payload
func VerifyBody(hash string, b []byte, sig, token string) error {
	// todo
	return nil // ErrPayloadUnverified
}

// ReadBody from an http.Request
func ReadBody(r *http.Request) ([]byte, error) {
	if r == nil {
		return nil, ErrInvalidRequest
	}

	switch r.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		break
	default:
		return nil, ErrInvalidRequest
	}

	if r.Body == nil {
		return nil, nil
	}

	defer goqa.Closed(r.Body)
	var b, err = io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Response is a generic reply
type Response struct {
	Message string `json:"message"`
}
