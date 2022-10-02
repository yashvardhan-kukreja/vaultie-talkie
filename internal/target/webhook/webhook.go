package webhook

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/yashvardhan-kukreja/vaultie-talkie/internal/target"
	"io"
	"net/http"
	"time"
)

type WebhookTarget struct {
	Url         string
	BearerToken string
}

func (w *WebhookTarget) Args() {
	flag.StringVar(&w.Url, "webhook-url", "", "Webhook URL which against which a POST request is triggered in case of vault-key store changes")
	flag.StringVar(&w.Url, "webhook-access-token", "", "Access token used to authn/authz vaultie-talkie against the Webhook URL")
}

func (w WebhookTarget) Execute(oldKeyStore, newKeyStore target.KeyStore) error {
	reqPayload := map[string]target.KeyStore{
		"old_key_store": oldKeyStore,
		"new_key_store": newKeyStore,
	}
	reqBody := new(bytes.Buffer)
	if err := json.NewEncoder(reqBody).Encode(reqPayload); err != nil {
		return fmt.Errorf("error occurred while marshalling the JSON of the webhook payload %+v: %w", reqPayload, err)
	}

	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("POST", w.Url, reqBody)
	if err != nil {
		return fmt.Errorf("error occurred while bootstrapping a new request for the webhook target: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if w.BearerToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.BearerToken))
	}

	log.Debugf("Triggering a webhook request at POST %s | Body: %s", w.Url, reqBody)

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error occurred while making the request to the webhook target: %w", err)
	}
	if resp.StatusCode >= 400 {
		log.Debug("Response was found to be of the status code", resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("request ended up with a client/server-side error with status code '%d': error occurred while reading the response body: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("request ended up with a client/server-side error with status code '%d': %s", resp.StatusCode, string(respBody))
	}
	log.Debug("Webhook request found to be delivered with status code", resp.StatusCode)
	return nil

}
