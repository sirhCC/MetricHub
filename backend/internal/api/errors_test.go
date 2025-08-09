package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
)

// helper to capture response for an error
func captureError(t *testing.T, err error) (int, map[string]any) {
	t.Helper()
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	respondError(c, err, "test error", nil)
	res := recorder.Result()
	defer res.Body.Close()
	var body map[string]any
	_ = json.NewDecoder(res.Body).Decode(&body)
	return res.StatusCode, body
}

func TestErrorMapping(t *testing.T) {
	cases := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{"not found", ErrNotFound, http.StatusNotFound, "not_found"},
		{"validation", ErrValidation, http.StatusBadRequest, "validation_error"},
		{"unauthorized", ErrUnauthorized, http.StatusUnauthorized, "unauthorized"},
		{"forbidden", ErrForbidden, http.StatusForbidden, "forbidden"},
		{"conflict", ErrConflict, http.StatusConflict, "conflict"},
		{"timeout", ErrTimeout, http.StatusGatewayTimeout, "timeout"},
		{"internal", ErrInternal, http.StatusInternalServerError, "internal_error"},
		{"unknown falls back", errUnknownType{}, http.StatusInternalServerError, "internal_error"},
	}

	for _, tc := range cases {
		status, body := captureError(t, tc.err)
		if status != tc.wantStatus {
			to := body["error"].(map[string]any)
			// include details for easier debugging
			if status != tc.wantStatus {
				// redundant but explicit for clarity
				// no t.Fatalf to allow seeing all failures
				t.Errorf("%s: unexpected status got=%d want=%d body=%v", tc.name, status, tc.wantStatus, body)
			}
			if to["code"] != tc.wantCode {
				t.Errorf("%s: unexpected code got=%v want=%s", tc.name, to["code"], tc.wantCode)
			}
			continue
		}
		errPayload, _ := body["error"].(map[string]any)
		if errPayload["code"] != tc.wantCode {
			to := errPayload["code"]
			t.Errorf("%s: code mismatch got=%v want=%s", tc.name, to, tc.wantCode)
		}
	}
}

// errUnknownType is a sentinel to test default mapping path
// implements error interface only

type errUnknownType struct{}

func (e errUnknownType) Error() string { return "unknown" }
