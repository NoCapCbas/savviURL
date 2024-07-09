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

  http.HandleFunc("/shorten", us.ShortenURL)
  http.HandleFunc("/", us.Redirect)

  log.Println("Server started at http://localhost:8083")
  if err := http.ListenAndServe(":8083", nil); err != nil {
    log.Fatalf("Failed to start server: %v", err)
  }
}

