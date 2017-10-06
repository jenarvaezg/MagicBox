package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jenarvaezg/magicbox/models"
	"github.com/jenarvaezg/magicbox/utils"
)

// RequireJSONMiddleware is a struct that has a ServeHTTP method
type RequireJSONMiddleware struct {
}

//RequireBoxMiddleware is a middleware that ensures a url's id parameter is a valid ID related to a Box document
type RequireBoxMiddleware struct {
}

//RequireUserMiddleware is a middleware that ensures a url's id parameter is a valid ID related to a User document
type RequireUserMiddleware struct {
}

// NewRequireJSONMiddleware returns a RequireJSONMiddleware
func NewRequireJSONMiddleware() *RequireJSONMiddleware {
	return &RequireJSONMiddleware{}
}

// NewRequireBoxMiddleware returns a RequireBoxMiddleware
func NewRequireBoxMiddleware() *RequireBoxMiddleware {
	return &RequireBoxMiddleware{}
}

// NewRequireUserMiddleware returns a RequireUserMiddleware
func NewRequireUserMiddleware() *RequireUserMiddleware {
	return &RequireUserMiddleware{}
}

/*
RequireJSONMiddleware's handler, which asserts that POST and PUT methods include content-type header
and is set to application/json
*/
func (l *RequireJSONMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	methodNeedsJSON := func(method string) bool {
		return method == "POST" || method == "PUT"
	}
	if methodNeedsJSON(r.Method) && r.Header.Get("content-type") != "application/json" {
		utils.ResponseError(w, "Expected content-type to be application/json", http.StatusBadRequest)
	} else {
		next(w, r)
	}
}

func getBox(r *http.Request) (models.Box, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	return models.GetBoxByID(id)
}

/*
RequireBoxMiddleware's handler, which asserts that url's id parameter is a valid ID and is related to a Box
document in the database
*/
func (l *RequireBoxMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	box, err := getBox(r)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), utils.ContextKeyBox, box))

	next(w, r)
}

func getUser(r *http.Request) (models.User, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	return models.GetUserByID(id)
}

/*
RequireUserMiddleware's handler, which asserts that url's id parameter is a valid ID and is related to a User
document in the database
*/
func (l *RequireUserMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	user, err := getUser(r)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), utils.ContextKeyUser, user))

	next(w, r)
}
