package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetRequest(t *testing.T) {
	var capturedReq *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key", "i123")
	resp, err := client.Do("GET", "/api/v1/athlete/{athleteId}/activities", map[string]string{"foo": "bar"}, nil)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	defer resp.Body.Close()

	if capturedReq.Method != "GET" {
		t.Errorf("Method = %q, want GET", capturedReq.Method)
	}

	if !strings.Contains(capturedReq.URL.Path, "i123") {
		t.Errorf("Path %q should contain athleteId i123", capturedReq.URL.Path)
	}

	if capturedReq.URL.Query().Get("foo") != "bar" {
		t.Errorf("Query param foo = %q, want bar", capturedReq.URL.Query().Get("foo"))
	}

	user, pass, ok := capturedReq.BasicAuth()
	if !ok {
		t.Fatal("BasicAuth not present")
	}
	if user != "API_KEY" {
		t.Errorf("BasicAuth user = %q, want API_KEY", user)
	}
	if pass != "test-key" {
		t.Errorf("BasicAuth pass = %q, want test-key", pass)
	}
}

func TestPostWithJSONBody(t *testing.T) {
	var capturedReq *http.Request
	var capturedBody map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &capturedBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key", "i123")
	payload, _ := json.Marshal(map[string]string{"name": "Morning Run"})
	resp, err := client.Do("POST", "/api/v1/athlete/{athleteId}/activities", nil, payload)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	defer resp.Body.Close()

	if capturedReq.Method != "POST" {
		t.Errorf("Method = %q, want POST", capturedReq.Method)
	}

	if capturedReq.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", capturedReq.Header.Get("Content-Type"))
	}

	if capturedBody["name"] != "Morning Run" {
		t.Errorf("body name = %q, want 'Morning Run'", capturedBody["name"])
	}
}

func TestPathSubstitution(t *testing.T) {
	var capturedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key", "i123")
	resp, err := client.Do("GET", "/api/v1/athlete/{athleteId}/activities/{activityId}", map[string]string{"activityId": "abc42"}, nil)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	defer resp.Body.Close()

	if !strings.Contains(capturedPath, "i123") {
		t.Errorf("Path %q should contain athleteId i123", capturedPath)
	}
	if !strings.Contains(capturedPath, "abc42") {
		t.Errorf("Path %q should contain activityId abc42", capturedPath)
	}

	// activityId should be consumed (not in query string)
	// We need to make sure there's no query string with activityId
	// (we can't check query string from path alone; just check path is correct)
	if strings.Contains(capturedPath, "activityId") {
		t.Errorf("Path %q should not contain literal 'activityId' placeholder", capturedPath)
	}
}

func TestPathSubstitutionNotInQuery(t *testing.T) {
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key", "i123")
	resp, err := client.Do("GET", "/api/v1/athlete/{athleteId}/activities/{activityId}", map[string]string{"activityId": "abc42"}, nil)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	defer resp.Body.Close()

	if strings.Contains(capturedQuery, "activityId") {
		t.Errorf("Query %q should not contain activityId (should be consumed by path substitution)", capturedQuery)
	}
}

func TestAPIErrorFormatting(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("activity not found"))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key", "i123")
	_, err := client.Do("GET", "/api/v1/athlete/{athleteId}/activities/999", nil, nil)
	if err == nil {
		t.Fatal("Do() error = nil, want error for 404")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("error type = %T, want *APIError", err)
	}

	if apiErr.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
	}

	errMsg := apiErr.Error()
	if !strings.Contains(errMsg, "404") {
		t.Errorf("Error() = %q, should contain '404'", errMsg)
	}
}

func TestDoRejectsPathInjection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("server should not be reached for injection attempt: %s", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key", "i123")
	cases := []string{"a/../b", "a?b", "a#frag", "a&x=1", "a%2e"}
	for _, bad := range cases {
		_, err := client.Do("GET", "/api/v1/athlete/{athleteId}/activities/{activityId}", map[string]string{"activityId": bad}, nil)
		if err == nil {
			t.Errorf("Do() with activityId=%q: expected error, got nil", bad)
		}
	}
}

func TestDryRunReturnsString(t *testing.T) {
	client := NewClient("https://example.com", "test-key", "i123")
	out, err := client.DryRun("GET", "/api/v1/athlete/{athleteId}/activities/{activityId}", map[string]string{"activityId": "abc42"}, nil)
	if err != nil {
		t.Fatalf("DryRun() error = %v", err)
	}
	if !strings.Contains(out, "abc42") || !strings.Contains(out, "i123") {
		t.Errorf("DryRun() = %q, want it to contain abc42 and i123", out)
	}
}

func TestDryRunErrorsOnBadPathValue(t *testing.T) {
	client := NewClient("https://example.com", "test-key", "i123")
	_, err := client.DryRun("GET", "/api/v1/athlete/{athleteId}/activities/{activityId}", map[string]string{"activityId": "a/../b"}, nil)
	if err == nil {
		t.Fatal("DryRun() error = nil, want error for injection path value")
	}
}

func TestDeleteRequest(t *testing.T) {
	var capturedMethod string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key", "i123")
	_, err := client.Do("DELETE", "/api/v1/athlete/{athleteId}/activities/123", nil, nil)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	if capturedMethod != "DELETE" {
		t.Errorf("Method = %q, want DELETE", capturedMethod)
	}
}
