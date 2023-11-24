package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/harmlessprince/snippetboxapp/pkg/models"
	"github.com/justinas/nosurf"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-XSS-Protection", "1; mode=block")
		writer.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(writer, request)
	})
}

func (app *application) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", request.RemoteAddr, request.Proto, request.Method, request.URL)
		next.ServeHTTP(writer, request)
	})
}

func (app *application) recoverFromPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				writer.Header().Set("Connection", "close")
				app.serverError(writer, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(writer, request)
	})
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if app.authenticatedUser(request) == nil {
			http.Redirect(writer, request, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(writer, request)
	})
}
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// Check if a userID value exists in the session. If this *isn't
		// present* then call the next handler in the chain as normal
		exists := app.session.Exists(request, "userID")
		if !exists {
			next.ServeHTTP(writer, request)
			return
		}
		// Fetch the details of the current user from the database. If
		// no matching record is found, remove the (invalid) userID from
		// their session and call the next handler in the chain as normal
		user, err := app.userModel.Get(app.session.GetInt(request, "userID"))
		if errors.Is(err, models.ErrNoRecord) {
			app.session.Remove(request, "userID")
			next.ServeHTTP(writer, request)
			return
		} else if err != nil {
			app.serverError(writer, err)
			return
		}
		// Otherwise, we know that the request is coming from a valid,
		// authenticated (logged in) user. We create a new copy of the
		// request with the user information added to the request context, and
		// call the next handler in the chain *using this new copy of the
		// request*.
		ctx := context.WithValue(request.Context(), contextKeyUser, user)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

// noSurf middleware function which uses a customized CSRF cookie with
// the Secure, Path and HttpOnly flags set.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}
