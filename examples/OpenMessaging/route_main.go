package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/gildas/go-purecloud"
)

func mainRouteHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.Must(logger.FromContext(r.Context()))
	config := core.Must(ConfigFromContext(r.Context())).(*Config)

	log.Infof("Request Headers: %#v", r.Header)
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Failed to read the request body", err)
		core.RespondWithError(w, http.StatusBadRequest, err)
		return
	}
	log.Infof("Received body (%d bytes): %s", len(body), string(body))

	signature := r.Header.Get("X-Hub-Signature-256")

	if len(signature) == 0 {
		log.Errorf("Request is missing [X-Hub-Signature-256] header")
		core.RespondWithError(w, http.StatusForbidden, err)
		return
	}

	crypto := hmac.New(sha256.New, []byte(config.IntegrationWebhookToken))
	_, _ = crypto.Write(body)
	signed := base64.StdEncoding.EncodeToString(crypto.Sum(nil))
	log.Debugf("Signed by Integration Token: %s", signed)
	if strings.TrimPrefix(signature, "sha256=") != signed {
		log.Errorf("Expected Signature (%s) does not match Token Signature (%s), rejecting", signature, signed)
		core.RespondWithError(w, http.StatusForbidden, errors.ArgumentInvalid.With("signature", signature))
		return
	}

	message := &purecloud.OpenMessage{}
	if err = json.Unmarshal(body, &message); err != nil {
		log.Errorf("Failed to unmarshal message", err)
		core.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	log.Record("message", message).Infof("Received From GCloud: %s", message.Text)
	core.RespondWithJSON(w, http.StatusOK, struct{}{})
}
