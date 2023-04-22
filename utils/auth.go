package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/mitchellh/go-homedir"
)

const (
    baseAuthURL = "https://stackoverflow.com/oauth/dialog"
    clientId = "26062"
	key = "w1BFZmzoMKahE3t5WYlEBA(("
    scope = "no_expiry"
    redirectUri = "http://localhost:6789/sotui-callback"
)

var token string

func GetToken() string {
	if token != "" {
		return token
	}

	dir, err := homedir.Dir(); if err != nil {
		panic("Unable to get home directory")
	}

	filePath := dir + "/.sotui/token"

	access_token, err := ioutil.ReadFile(filePath); if err != nil {
		return ""
	} else {
		token = string(access_token)
	}

	return token
}

func SetToken(access_token string) {
	dir, err := homedir.Dir(); if err != nil {
		panic("Unable to get home directory")
	}

	filePath := dir + "/.sotui/token"

	if _, err := os.Stat(filePath); os.IsNotExist(err) { 
		os.Mkdir(dir + "/.sotui" , 0777)
	} //TODO: More restrictive permissions

	err = os.WriteFile(filePath, []byte(access_token), 0777); if err != nil {
		panic(err)
	}
	token = access_token
}

func Oauth2() {
	m := http.NewServeMux()
    s := http.Server{Addr: ":6789", Handler: m}

	m.HandleFunc("/sotui-callback", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		html := `
		<script>
		document.write("Authenticating... , make sure javascript is enabled")
		window.location.href = window.location.href.replace("#", "?")
		</script>
		`
		if r.URL.RawQuery == "" {
			w.Write([]byte(html))
		} else{
			token := r.URL.Query().Get("access_token")

			if token == "" {
				w.Write([]byte("Authentication failed!"))
				panic("Unable to get token")
			} else {
				w.Write([]byte("Authentication successful!"))
				SetToken(token)
				go s.Shutdown(context.Background())
			}
		}

    })

	s.ListenAndServe()
}

func GetAuthURL() string {
	return fmt.Sprintf("%s?client_id=%s&scope=%s&redirect_uri=%s", baseAuthURL, clientId, scope, redirectUri)
}