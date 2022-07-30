package myAuth

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
	"strings"
)

type DoAuthStruct struct {
	Code string `json:"code"`
}

//THANKS: https://stackoverflow.com/questions/48855122/keycloak-adaptor-for-golang-application !!!

func DoKcAuth(configURL string, clientID string, clientSecret string, redirectURL string, state string) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, configURL)
	if err != nil {
		log.Fatal(err)
	}
	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier := provider.Verifier(oidcConfig)

	http.HandleFunc("/oidc/auth", func(w http.ResponseWriter, r *http.Request) {
		rawAccessToken := r.Header.Get("Authorization")
		if rawAccessToken == "" {
			http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
			return
		}
		parts := strings.Split(rawAccessToken, " ")
		if len(parts) != 2 {
			w.WriteHeader(400)
			return
		}
		_, err := verifier.Verify(ctx, parts[1])
		if err != nil {
			http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
			fmt.Print(err)
			return
		}
		w.Write([]byte("valid"))
	})

	http.HandleFunc("/oidc/id", func(w http.ResponseWriter, r *http.Request) {
		var getCode DoAuthStruct
		bytes, _ := io.ReadAll(r.Body)
		json.Unmarshal(bytes, &getCode)
		oauth2Token, err := oauth2Config.Exchange(ctx, getCode.Code)
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
			return
		}
		idToken, err := verifier.Verify(ctx, rawIDToken)
		if err != nil {
			http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		resp := struct {
			OAuth2Token   *oauth2.Token
			IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
		}{oauth2Token, new(json.RawMessage)}

		if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		w.Write([]byte(resp.OAuth2Token.AccessToken))
	})
}

func ParseJWT(p string) ([]byte, error) {
	parts := strings.Split(p, ".")
	payload, _ := b64.RawURLEncoding.DecodeString(parts[1])
	return payload, nil
}
