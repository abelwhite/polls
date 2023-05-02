// Filename: cmd/web/routes.go
package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// ROUTES: 10
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	dynamicMiddleware := alice.New(app.sessionManager.LoadAndSave, noSurf)
	//we wrap

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/about", app.about)
	router.Handler(http.MethodGet, "/poll/success", dynamicMiddleware.ThenFunc(app.pollSuccessShow))
	router.Handler(http.MethodGet, "/poll/create", dynamicMiddleware.ThenFunc(app.pollCreateShow)) //wrap poll.create with dynamic middlware
	router.HandlerFunc(http.MethodPost, "/poll/create", app.pollCreateSubmit)
	router.HandlerFunc(http.MethodGet, "/options/create", app.optionsCreateShow)
	router.HandlerFunc(http.MethodPost, "/options/create", app.optionsCreateSubmit)

	router.Handler(http.MethodGet, "/user/signup", dynamicMiddleware.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamicMiddleware.ThenFunc(app.userSignupSubmit))
	router.Handler(http.MethodGet, "/user/login", dynamicMiddleware.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamicMiddleware.ThenFunc(app.userLoginSubmit))
	router.Handler(http.MethodPost, "/user/logout", dynamicMiddleware.ThenFunc(app.userLogoutSubmit))

	//protect the /poll/reply route
	protected := dynamicMiddleware.Append(app.requireAuthenticationMiddleware)

	router.Handler(http.MethodGet, "/poll/reply", protected.ThenFunc(app.pollReplyShow))
	router.Handler(http.MethodPost, "/poll/reply", protected.ThenFunc(app.pollReplySubmit))
	//returns to the router to our middleware interesting in between go server and mux
	//Client -> Goserver ->Middleware -> Router -> Handler

	//tidy up the middle wear
	standardMiddleware := alice.New(app.RecoverPanicMiddleware, app.logRequestMiddleware, securityHeadersMiddleware)

	return standardMiddleware.Then(router)
}
