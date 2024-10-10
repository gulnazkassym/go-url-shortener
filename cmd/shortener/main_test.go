package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_postHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Valid POST request",
			method:       http.MethodPost,
			body:         "http://example.com",
			expectedCode: http.StatusCreated,
			expectedBody: "http://localhost:8080/hevfyegruf",
		},
		{
			name:         "Empty body",
			method:       http.MethodPost,
			body:         "",
			expectedCode: http.StatusBadRequest,
			expectedBody: "URL cannot be empty\n",
		},
		{
			name:         "Unsupported method",
			method:       http.MethodGet,
			body:         "http://example.com",
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "http://localhost:8080/", bytes.NewBufferString(tt.body))
			rr := httptest.NewRecorder()

			postHandler(rr, req)

			resp := rr.Result()
			if resp.StatusCode != tt.expectedCode {
				t.Errorf("expected status %v, got %v", tt.expectedCode, resp.StatusCode)
			}

			body, _ := io.ReadAll(resp.Body)
			if string(body) != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, string(body))
			}
			defer resp.Body.Close()
		})
	}
}

func Test_getHandler(t *testing.T) {
	shortenStr := "hevfyegruf"
	originalURL := "http://example.com"
	urlStore[shortenStr] = originalURL

	tests := []struct {
		name         string
		path         string
		expectedCode int
		expectedURL  string
	}{
		{
			name:         "Valid request for existing URL",
			path:         "/" + shortenStr,
			expectedCode: http.StatusTemporaryRedirect,
			expectedURL:  originalURL,
		},
		{
			name:         "Request for non-existing URL",
			path:         "/nonexistent",
			expectedCode: http.StatusNotFound,
			expectedURL:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080"+tt.path, nil)
			rr := httptest.NewRecorder()

			getHandler(rr, req)

			resp := rr.Result()
			// if resp.StatusCode != tt.expectedCode {
			// 	t.Errorf("expected status %v, got %v", tt.expectedCode, resp.StatusCode)
			// }

			if tt.expectedURL != "" {
				if location := resp.Header.Get("Location"); location != tt.expectedURL {
					t.Errorf("expected redirect to %q, got %q", tt.expectedURL, location)
				}
			}
			defer resp.Body.Close()
		})
	}
}
