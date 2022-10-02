package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yashvardhan-kukreja/vaultie-talkie/internal/target"
)

func TestWebhookTargetSuccess(t *testing.T) {
	m := mockWebhookServerOpts{
		port:                     8000,
		path:                     "/webhook",
		returnResponseStatusCode: 200,
		returnResponseBody:       map[string]interface{}{"success": true, "message": "Well, done!"},
	}
	assert.NoError(t, m.setup())
	assert.NotNil(t, m.mockServer)

	wt := WebhookTarget{
		Url: fmt.Sprintf("http://localhost:8000%s", m.path),
	}

	oldKeyStore := target.KeyStore(map[string]interface{}{
		"foo": "bar",
	})
	newKeyStore := target.KeyStore(map[string]interface{}{
		"foo": "bar",
		"a":   "b",
	})

	assert.NoError(t, wt.Execute(oldKeyStore, newKeyStore))
	assert.NoError(t, m.teardown())
}

func TestWebhookTargetFailure(t *testing.T) {
	m := mockWebhookServerOpts{
		port:                     8000,
		path:                     "/webhook-2",
		returnResponseStatusCode: 401,
		returnResponseBody:       map[string]interface{}{"success": true, "message": "Well, done!"},
	}
	assert.NoError(t, m.setup())
	assert.NotNil(t, m.mockServer)

	wt := WebhookTarget{
		Url: fmt.Sprintf("http://localhost:8000%s", m.path),
	}

	oldKeyStore := target.KeyStore(map[string]interface{}{
		"foo": "bar",
	})
	newKeyStore := target.KeyStore(map[string]interface{}{
		"foo": "bar",
		"a":   "b",
	})

	expectedResp, err := json.Marshal(m.returnResponseBody)
	assert.NoError(t, err)

	assert.EqualError(t, wt.Execute(oldKeyStore, newKeyStore), fmt.Sprintf("request ended up with a client/server-side error with status code '%d': %s", m.returnResponseStatusCode, string(expectedResp)))
	assert.Nil(t, m.teardown())
}

type mockWebhookServerOpts struct {
	port                     int
	path                     string
	returnResponseStatusCode int
	returnResponseBody       map[string]interface{}
	mockServer               *http.Server
	expectedBearerToken      string
}

func (m *mockWebhookServerOpts) setup() error {
	http.HandleFunc(m.path, func(w http.ResponseWriter, r *http.Request) {
		if m.expectedBearerToken != "" {
			receivedToken := r.Header.Get("Authorization")
			if fmt.Sprintf("Bearer %s", m.expectedBearerToken) != receivedToken {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("authorization bearer token didn't match the expected one"))
				return
			}
		}
		resp, err := json.Marshal(m.returnResponseBody)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed to marshal the response body"))
			return
		}
		w.WriteHeader(m.returnResponseStatusCode)
		w.Write(resp)
	})

	m.mockServer = &http.Server{Addr: fmt.Sprintf(":%d", m.port)}
	go func(sv *http.Server) {
		sv.ListenAndServe()
	}(m.mockServer)
	return nil
}

func (m *mockWebhookServerOpts) teardown() error {
	if m.mockServer != nil {
		return m.mockServer.Shutdown(context.Background())
	}
	return nil
}
