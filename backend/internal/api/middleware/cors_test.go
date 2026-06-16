package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCorsMiddlewareAllowsCredentialedPatchPreflight(t *testing.T) {
	nextCalled := false
	handler := CorsMiddleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		nextCalled = true
	}), "http://localhost:5173")

	request := httptest.NewRequest(http.MethodOptions, "/api/posts/123", nil)
	request.Header.Set("Origin", "http://localhost:5173")
	request.Header.Set("Access-Control-Request-Method", http.MethodPatch)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusOK)
	}
	if nextCalled {
		t.Fatal("preflight request reached the next handler")
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Errorf("Access-Control-Allow-Origin = %q", got)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Errorf("Access-Control-Allow-Credentials = %q, want true", got)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Methods"); !strings.Contains(got, http.MethodPatch) {
		t.Errorf("Access-Control-Allow-Methods = %q, want PATCH", got)
	}
}
