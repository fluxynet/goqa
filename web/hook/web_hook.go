package hook

import (
	"encoding/json"
	"net/http"

	"github.com/fluxynet/goqa"
	"github.com/fluxynet/goqa/web"
)

type Hook struct {
	Broker    goqa.Broker
	Signature string
	Token     string
}

// Receive a web hook
func (h *Hook) Receive(w http.ResponseWriter, r *http.Request) {
	var body, err = web.ReadBody(r)

	if err != nil {
		web.JsonError(w, http.StatusBadRequest, err)
		return
	}

	var hash string
	if err = web.VerifyBody(hash, body, h.Signature, h.Token); err != nil {
		web.JsonError(w, http.StatusBadRequest, err)
		return
	}

	var payload *Payload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		web.JsonError(w, http.StatusBadRequest, err)
		return
	}

	var event = CreateGithubEvent(payload)
	if event == nil {
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
