package main

import (
	"errors"
	"fmt"
	"github.com/harmlessprince/snippetboxapp/pkg/forms"
	"github.com/harmlessprince/snippetboxapp/pkg/models"
	"net/http"
	"strconv"
)

func ping(w http.ResponseWriter, request *http.Request) {
	w.Write([]byte("OK"))
}
func (app *application) home(w http.ResponseWriter, request *http.Request) {
	snippets, err := app.snippetModel.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := &templateData{Snippets: snippets}
	app.render(w, request, "home.page.tmpl", data)
}
func (app *application) showSnippet(w http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippetModel.Get(id)
	if errors.Is(err, models.ErrNoRecord) {
		app.notFound(w)
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := &templateData{Snippet: snippet}
	app.render(w, request, "show.page.tmpl", data)
}
func (app *application) createSnippet(w http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	title := request.PostForm.Get("title")
	content := request.PostForm.Get("content")
	expires := request.PostForm.Get("expires")

	form := forms.New(request.PostForm)

	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.MaxLength("content", 500)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, request, "create.page.tmpl", &templateData{
			Form:     form,
			FormData: request.PostForm,
		})
		return
	}

	id, err := app.snippetModel.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(request, "flash", "Snippet successfully created!")
	//Redirect user to relevant page of the snippet created
	http.Redirect(w, request, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "create.page.tmpl", &templateData{Form: forms.New(nil)})
}

func (app *application) signupUserForm(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "signup.page.tmpl", &templateData{Form: forms.New(nil)})
}

func (app *application) signupUser(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
	}
	form := forms.New(request.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRegexPattern)
	form.MinLength("password", 4)
	if !form.Valid() {
		app.render(writer, request, "signup.page.tmpl", &templateData{
			Form: form,
		})
	}
	err = app.userModel.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if errors.Is(err, models.ErrDuplicateEmail) {
		form.Errors.Add("email", "Email address already in use")
		app.render(writer, request, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(writer, err)
		return
	}
	app.session.Put(request, "flash", "Your signup was successful, Please log in")
	http.Redirect(writer, request, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "login.page.tmpl", &templateData{Form: forms.New(nil)})
}

func (app *application) loginUser(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
	}
	form := forms.New(request.PostForm)
	//form.Required("email", "password")
	//form.MatchesPattern("email", forms.EmailRegexPattern)
	//form.MinLength("password", 4)
	//if !form.Valid() {
	//	app.render(writer, request, "login.page.tmpl", &templateData{
	//		Form: form,
	//	})
	//}
	id, err := app.userModel.Authenticate(form.Get("email"), form.Get("password"))
	if errors.Is(err, models.ErrInvalidCredentials) {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(writer, request, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(writer, err)
		return
	}

	app.session.Put(request, "userID", id)

	http.Redirect(writer, request, "/snippet/create", http.StatusSeeOther)
}

func (app *application) logoutUser(writer http.ResponseWriter, request *http.Request) {
	app.session.Remove(request, "userID")
	app.session.Put(request, "flash", "You have been logged out successfully")
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}
