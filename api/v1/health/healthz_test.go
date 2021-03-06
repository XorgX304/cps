package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rapid7/cps/watchers/v1/consul"
	"github.com/rapid7/cps/watchers/v1/s3"
)

func TestGetHealthz(t *testing.T) {
	consul.Up = false
	s3.Up = false

	req, err := http.NewRequest("GET", "/v1/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetHealthz)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("Status code is wrong when unhealthy: expected %v got %v", status, http.StatusServiceUnavailable)
	}

	expectedJSON := `{"status":"down","consul":false,"s3":false}`
	assert.Equal(t, expectedJSON, rr.Body.String())

	consul.Up = true
	s3.Up = true

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(GetHealthz)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code is wrong when services are healthyish: expected %v got %v", status, http.StatusOK)
	}

	assert.NotNil(t, rr.Body.String())

	expectedJSON = `{"status":"up","consul":true,"s3":true}`
	assert.Equal(t, expectedJSON, rr.Body.String())
}
