package main

import (
	"github.com/auth0/auth0-go/examples/regular-web-app/routes/callback"
	"github.com/auth0/auth0-go/examples/regular-web-app/routes/home"
	"github.com/auth0/auth0-go/examples/regular-web-app/routes/user"
	"github.com/gorilla/mux"
	"net/http"
)

func StartServer() {
	r := mux.NewRouter()

	r.HandleFunc("/", home.HomeHandler)
	r.HandleFunc("/callback", callback.CallbackHandler)
	r.HandleFunc("/user", user.UserHandler)
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
