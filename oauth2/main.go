package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"html/template"
	"log"
	"net/http"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type Authorization struct {
	Access_token string
}

type GithubUser struct {
	Id   int64
	Name string
}

type GoogleUser struct {
	Id   string
	Name string
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func github_callback(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	code := r.FormValue("code")

	client_id := "7dbd0a9f61d655243969"
	client_secret, err := accessSecretVersion("projects/708733497091/secrets/github-oauth2-client-secret/versions/latest")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := "https://github.com/login/oauth/access_token"
	data := "client_id=" + client_id + "&client_secret=" + client_secret + "&code=" + code

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	authorization := Authorization{}
	json.NewDecoder(resp.Body).Decode(&authorization)
	log.Println("access_token:", authorization.Access_token)

	url = "https://api.github.com/user"

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+authorization.Access_token)
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	user := GithubUser{}
	json.NewDecoder(resp.Body).Decode(&user)
	log.Println("id:", user.Id)
	log.Println("name:", user.Name)

	http.Redirect(w, r, "/", http.StatusFound)
}

func google_callback(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	code := r.FormValue("code")

	client_id := "708733497091-l3njp9vfnni4v5misr1b2fepgbr4409t.apps.googleusercontent.com"
	client_secret, err := accessSecretVersion("projects/708733497091/secrets/google-oauth2-client-secret/versions/latest")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := "https://oauth2.googleapis.com/token"
	data := "client_id=" + client_id + "&client_secret=" + client_secret + "&code=" + code + "&grant_type=authorization_code&redirect_uri=http://localhost:8080/login/oauth2/code/google"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	authorization := Authorization{}
	json.NewDecoder(resp.Body).Decode(&authorization)
	log.Println("access_token:", authorization.Access_token)

	url = "https://www.googleapis.com/oauth2/v1/userinfo"

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+authorization.Access_token)
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	user := GoogleUser{}
	json.NewDecoder(resp.Body).Decode(&user)
	log.Println("id:", user.Id)
	log.Println("name:", user.Name)

	http.Redirect(w, r, "/", http.StatusFound)
}

// accessSecretVersion accesses the payload for the given secret version if one
// exists. The version can be a version number as a string (e.g. "5") or an
// alias (e.g. "latest").
func accessSecretVersion(name string) (string, error) {

	// name := "projects/my-project/secrets/my-secret/versions/5"
	// name := "projects/my-project/secrets/my-secret/versions/latest"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	// Verify the data checksum.
	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(result.Payload.Data, crc32c))
	if checksum != *result.Payload.DataCrc32C {
		return "", fmt.Errorf("Data corruption detected.")
	}

	return string(result.Payload.Data), nil
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login/oauth2/code/github", github_callback)
	http.HandleFunc("/login/oauth2/code/google", google_callback)

	log.Print("Listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
