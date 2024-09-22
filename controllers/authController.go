package controllers

import "net/http"

// register user
func Register(w http.ResponseWriter, r *http.Request) {
	CreateUser(w, r)
}
