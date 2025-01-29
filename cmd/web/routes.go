package main

import "github.com/julienschmidt/httprouter"
import "github.com/justinas/alice"
import "net/http"

func (app *Application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			app.notFound(w)
		},
	)

	static := http.Dir("./ui/static")
	fserver := http.FileServer(static)
	router.Handler(http.MethodGet, "/static/*fpath", http.StripPrefix("/static", fserver))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.userLogoutPost))

	midware := alice.New(app.panicHandler, app.logRequest, secureHeaders)

	return midware.Then(router)
}
