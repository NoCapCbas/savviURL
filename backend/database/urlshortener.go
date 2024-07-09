package database

import (
  "log"
)

func CreateURLMapping(key, url string) error {

  _, err := db.Exec(`
  INSERT INTO url_mappings (short_key, long_url)
  VALUES ($1, $2)
  `, key, url)
  if err != nil {
    log.Printf("Error inserting url mapping: %v", err)
    return err
  }
  return nil
}

func GetURLMapping(key string) (string, error) {
  var longURL string
  err := db.QueryRow(`
   SELECT long_url FROM url_mappings WHERE short_key = $1;
  `, key).Scan(&longURL)
  if err != nil {
    log.Printf("Error querying database for short key: %v", err)
    return "", err
  }
  return longURL, nil
}
