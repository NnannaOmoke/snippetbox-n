package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Print(trace)
	statusText := http.StatusText(http.StatusInternalServerError)
	http.Error(w, statusText, http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, status int) {
	statusText := http.StatusText(status)
	http.Error(w, statusText, status)
}

func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *Application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("The template {%s} does not exist.", page)
		app.serverError(w, err)
		return
	}

	var buf bytes.Buffer

	err := ts.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)

}

func (app Application) newTemplateData(r *http.Request) templateData {
	return templateData{CurrentYear: time.Now().Year(), Flash: app.sessionManager.PopString(r.Context(), "flash")}
}
