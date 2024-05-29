package urlshortener

import (
  "encoding/json"
  "math/rand"
  "net/http"
  "sync"
  "time"
  "fmt"
)

type URLShortener struct {
  sync.Mutex
  UrlMap map[string]string
}

func NewURLShortener() *URLShortener {
  return &URLShortener{
    UrlMap: make(map[string]string),
  }
}

func (us *URLShortener) GenerateKey(n int) string {
  const chrs = "abcdefghijklmnopqrstuvwxyz0123456789"
  b := make([]rune, n)
  for i := range b {
    b[i] = rune(chrs[rand.Intn(len(chrs))])
  }
  return string(b)
}

func (us *URLShortener) ShortenURL(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Pinged ShortenURL Handler...")
  var req struct {
    URL string `json:"URL"`
  }
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    http.Error(w, "Invalid request payload", http.StatusBadRequest)
    return
  }
  
  // Debugging
  fmt.Printf("Request Body: %+v\n", req)

  us.Lock()
  defer us.Unlock()

  key := us.GenerateKey(6)
  us.UrlMap[key] = req.URL

  resp := map[string]string{"short_url": "http://localhost:8083/" + key}
  if err := json.NewEncoder(w).Encode(resp); err != nil {
    http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    return
  }
}

func (us *URLShortener) Redirect(w http.ResponseWriter, r *http.Request) {
  key := r.URL.Path[1:]

  us.Lock()
  longURL, exists := us.UrlMap[key]
  us.Unlock()

  if !exists {
    http.NotFound(w, r)
    return
  }

  http.Redirect(w, r, longURL, http.StatusFound)
}

func init() {
  rand.Seed(time.Now().UnixNano())
}

