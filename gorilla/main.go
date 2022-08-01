package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var (
	key   = []byte(securecookie.GenerateRandomKey(32))
	store = sessions.NewCookieStore(key)
	//store = sessions.NewFilesystemStore("./", key)
)

func index(w http.ResponseWriter, r *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := store.Get(r, "go-session")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	fmt.Fprintln(w, "Hello, logged in user!")
}

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)

	case "POST":
		// TODO: Validate username and password!
		if r.FormValue("username") != "" && r.FormValue("password") != "" {
			session, _ := store.Get(r, "go-session")

			// Set user as authenticated
			session.Values["authenticated"] = true
			err := session.Save(r, w)
			if err != nil {
				http.Error(w, "Failed to save session", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "go-session")

	session.Values["authenticated"] = false
	session.Options.MaxAge = -1  // will delete the session if using FilesystemStore

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Logged out!")
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	log.Print("Listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
