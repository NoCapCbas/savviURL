package test

import (
  "bytes"
  "encoding/json"
  "net/http"
  "net/http/httptest"
  "testing"
  "urlshortener/urlshortener"
)

func TestShortenURL(t *testing.T) {
    // Create a new URLShortener instance
    us := urlshortener.NewURLShortener()

    // Create a new request with a JSON payload
    reqBody := map[string]string{"URL": "https://google.com"}
    reqBytes, _ := json.Marshal(reqBody)
    req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(reqBytes))
    req.Header.Set("Content-Type", "application/json")

    // Create a ResponseRecorder to record the response
    rr := httptest.NewRecorder()

    // Call the shortenURL handler function
    handler := http.HandlerFunc(us.ShortenURL)
    handler.ServeHTTP(rr, req)

    // Check the response status code
    if rr.Code != http.StatusOK {
        t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
    }

    // Decode the response body
    var respBody map[string]string
    json.NewDecoder(rr.Body).Decode(&respBody)

    // Check if the short_url is present in the response
    if _, ok := respBody["short_url"]; !ok {
        t.Error("Expected 'short_url' field in response, got none")
    }
}
