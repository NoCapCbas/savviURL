package main

import (
  "log"
  "net/http"
  "urlshortener/urlshortener"
  "urlshortener/database"
)

func main() {
  database.Connect()
  database.AutoMigrate()
  defer database.CloseDB()

  us := urlshortener.NewURLShortener()

  http.HandleFunc("/savvi-url/shorten", us.ShortenURL)
  http.HandleFunc("/savvi-url/", us.Redirect)

  log.Println("Server started at http://localhost:8083")
  if err := http.ListenAndServe(":8083", nil); err != nil {
    log.Fatalf("Failed to start server: %v", err)
  }
}

