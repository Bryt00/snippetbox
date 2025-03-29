package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleWare := alice.New(app.session.Enable, noSurf, app.authenticate)
	mux := pat.New()
	mux.Get("/", dynamicMiddleWare.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleWare.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleWare.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleWare.ThenFunc(app.showSnippet))

	//auth routes
	mux.Get("/user/signup", dynamicMiddleWare.ThenFunc(app.signUpUsrForm))
	mux.Post("/user/signup", dynamicMiddleWare.ThenFunc(app.signUpUsr))
	mux.Get("/user/login", dynamicMiddleWare.ThenFunc(app.loginUsrForm))
	mux.Post("/user/login", dynamicMiddleWare.ThenFunc(app.loginUsr))
	mux.Post("/user/logout", dynamicMiddleWare.Append(app.requireAuthentication).ThenFunc(app.logoutUser))
	mux.Get("/about", dynamicMiddleWare.ThenFunc(app.aboutPage))
	mux.Get("/user/profile", dynamicMiddleWare.Append(app.requireAuthentication).ThenFunc(app.userProfile))
	mux.Get("/user/changepassword", dynamicMiddleWare.Append(app.requireAuthentication).ThenFunc(app.changePasswdForm))
	mux.Post("/user/changepassword", dynamicMiddleWare.Append(app.requireAuthentication).ThenFunc(app.changePasswd))

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}
