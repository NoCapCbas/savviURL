package test

import (
  "net/http"
  "net/http/httptest"
  "testing"
  "urlshortener/urlshortener"
)

func TestRedirectURL(t *testing.T) {
    // Create a new URLShortener instance
    us := urlshortener.NewURLShortener()

    // Manually add a URL to the URLShortener's map
    key := us.GenerateKey(6)
    longURL := "https://google.com"
    us.UrlMap[key] = longURL

    // Create a new request to redirect
    req, _ := http.NewRequest("GET", "/"+key, nil)

    // Create a ResponseRecorder to record the response
    rr := httptest.NewRecorder()

    // Call the redirect handler function
    handler := http.HandlerFunc(us.Redirect)
    handler.ServeHTTP(rr, req)

    // Check the response status code
    if rr.Code != http.StatusFound {
        t.Errorf("Expected status code %d, got %d", http.StatusFound, rr.Code)
    }

    // Check if the Location header is set correctly
    if rr.Header().Get("Location") != longURL {
        t.Errorf("Expected Location header to be '%s', got '%s'", longURL, rr.Header().Get("Location"))
    }
}
