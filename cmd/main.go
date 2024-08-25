package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var tmpl *template.Template
	// Load templates
	tmpl, err = template.ParseFiles(
		filepath.Join("ui", "templates", "base.html"),
		filepath.Join("ui", "templates", "login.html"),
	)
	if err != nil {
		log.Println("Error parsing templates:", err)
		os.Exit(1)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println("Error executing template:", err)
		os.Exit(1)
	}
}

func googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	googleOauthConfig := &oauth2.Config{
		RedirectURL:  "https://urls.savvilabs.co/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	token, err := googleOauthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Println("Error exchanging token:", err)
		os.Exit(1)
	}

	fmt.Println(token)
}

func googleSignInHandler(w http.ResponseWriter, r *http.Request) {
	// Google OAuth2 configuration
	googleOauthConfig := &oauth2.Config{
		RedirectURL:  "https://urls.savvilabs.co/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	url := googleOauthConfig.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func main() {
	var err error
	// serve static files
	fs := http.FileServer(http.Dir("ui/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/auth/google/callback", googleCallbackHandler)
	http.HandleFunc("/auth/google/signin", googleSignInHandler)

	log.Println("Server is running on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error starting server:", err)
	}
	os.Exit(0)
}
