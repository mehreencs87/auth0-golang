package callback

import (
	_ "crypto/sha512"
	"encoding/json"
	"github.com/auth0/auth0-go/examples/regular-web-app/app"
	"io/ioutil"
	"net/http"
	"os"
	// "golang.org/x/oauth2"
	"github.com/auth0/golang-oauth2"
)

func CallbackHandler(w http.ResponseWriter, r *http.Request) {

	domain := os.Getenv("AUTH0_DOMAIN")

	opts, err := oauth2.New(
		oauth2.Client(os.Getenv("AUTH0_CLIENT_ID"), os.Getenv("AUTH0_CLIENT_SECRET")),
		oauth2.RedirectURL(os.Getenv("AUTH0_CALLBACK_URL")),
		oauth2.Scope("openid", "profile"),
		oauth2.Endpoint(
			"https://"+domain+"/authorize",
			"https://"+domain+"/oauth/token",
		),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	code := r.URL.Query().Get("code")

	transport, err := opts.NewTransportFromCode(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := transport.Token()

	// Getting now the userInfo
	client := http.Client{Transport: transport}
	resp, err := client.Get("https://" + domain + "/userinfo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	if err := json.Unmarshal(raw, &profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, _ := app.GlobalSessions.SessionStart(w, r)
	defer session.SessionRelease(w)

	session.Set("id_token", token.Extra("id_token"))
	session.Set("access_token", token.AccessToken)
	session.Set("profile", profile)

	// Redirect to logged in page
	http.Redirect(w, r, "/user", http.StatusMovedPermanently)

}
