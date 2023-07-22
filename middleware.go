package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Jiraffe7/imgram/user"
)

var errInvalidAuth = errors.New("invalid authentication")

// UserAuth is a middleware that authenticates the user.
func UserAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: actual auth
		username, _, ok := r.BasicAuth()
		if !ok {
			respondError(w, http.StatusUnauthorized, errInvalidAuth)
			return
		}

		// Assume username is userID.
		//TODO: actual impl
		userID, err := strconv.ParseUint(username, 10, 64)
		if err != nil {
			respondError(w, http.StatusUnauthorized, errInvalidAuth)
			return
		}

		u := user.User{
			ID: userID,
		}
		ctx := user.NewContext(r.Context(), &u)

		newR := r.WithContext(ctx)

		h.ServeHTTP(w, newR)
	})
}
