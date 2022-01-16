package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/oauth2"

	oidc "github.com/coreos/go-oidc"
)

var (
	clientID     = "myclient"
	clientSecret = "BDeJPtXKIMn0TKEfyT7ChAOsJyyrclYV"
)

func main() {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "http://localhost:8080/auth/realms/myreal")
	if err != nil {
		log.Fatal(err)
	}
	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:8081/auth/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	state := "123"

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, config.AuthCodeURL(state), http.StatusFound)
	})

	http.HandleFunc("/auth/callback/", func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Query().Get("state") != state {
			http.Error(writer, "state invalido", http.StatusBadRequest)
			return
		}
		token, err := config.Exchange(ctx, request.URL.Query().Get("code"))
		if err != nil {
			http.Error(writer, "falha ao trocar o token", http.StatusInternalServerError)
			return
		}
		resp := struct {
			AccessToken *oauth2.Token //token de autorização
		}{
			token,
		}

		data, err := json.Marshal(resp)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Write(data)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
