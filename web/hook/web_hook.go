package hook

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fluxynet/goqa"
	"github.com/fluxynet/goqa/web"
)

const githubHeaderSignature = "X-Hub-Signature"

var (
	errIncompleteRequest = errors.New("request incomplete")
)

type Hook struct {
	Broker goqa.Broker
	// SigKey used in hash
	SigKey string
}

// Receive a web hook
func (h *Hook) Receive(w http.ResponseWriter, r *http.Request) {
	var (
		body, err = web.ReadBody(r)
		signature = r.Header.Get(githubHeaderSignature)
	)

	if err != nil || len(body) == 0 || signature == "" {
		web.JsonError(w, http.StatusBadRequest, errIncompleteRequest)
		return
	}

	if err = web.VerifyBody(body, signature, h.SigKey); err != nil {
		web.JsonError(w, http.StatusBadRequest, err)
		return
	}

	var payload *Payload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		web.JsonError(w, http.StatusBadRequest, errIncompleteRequest)
		return
	}

	var event = CreateGithubEvent(payload)
	if len(event.Coverage) == 0 {
		web.Json(w, web.Response{Message: "web hook was not very interesting"})
		return
	}

	err = h.Broker.Publish(r.Context(), event)
	if err != nil {
		web.JsonError(w, http.StatusInternalServerError, err)
		return
	}

	web.Json(w, web.Response{Message: "web hook well received"})
}
