//filename: cmd/web/middleware.go

package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("X-Frame-Options", "deny")

			next.ServeHTTP(w, r)
		})
}

//the func( w http.ResponseWrite..) is a regular function that is being passed
//Intercepts a response0

func (app *application) logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//when the request comes to me
		start := time.Now()
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr,
			r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
		// when the response comes to me
		app.infoLog.Printf("Request took %v", time.Since(start))
	})
}

func (app *application) RecoverPanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("connection", "close")
				trace := fmt.Sprintf("%s\n%s", err, debug.Stack()) //we make a string and then debug the stack
				app.errorLog.Output(2, trace)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
