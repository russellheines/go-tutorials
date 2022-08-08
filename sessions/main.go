package main

import (
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

type User struct {
	ID   string
	Name string
}

func index(w http.ResponseWriter, r *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := store.Get(r, "go-session")

	var user *User
	if _, ok := session.Values["userid"]; ok {
		user = &User{
			Name: session.Values["name"].(string),
		}
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, user)
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "go-session")

	session.Values["userid"] = 101
	session.Values["name"] = "Gorilla User"

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "go-session")

	delete(session.Values, "userid")
	delete(session.Values, "name")

	session.Options.MaxAge = -1 // if using FilesystemStore, this will also delete the session on the server side

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
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