package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/abelwhite/poll/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	flash := app.sessionManager.PopString(r.Context(), "flash")
	//render
	data := &templateData{ //putting flash into template data
		Flash: flash,
	}
	RenderTemplate(w, "home.page.tmpl", data)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "about.page.tmpl", nil)

}

func (app *application) pollReplyShow(w http.ResponseWriter, r *http.Request) {
	question, err := app.questions.Get()
	if err != nil {
		return
	}
	data := &templateData{
		Question: question,
	}
	RenderTemplate(w, "poll.page.tmpl", data)
}

func (app *application) pollReplySubmit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	response := r.PostForm.Get("emotion")
	question_id := r.PostForm.Get("id")
	identifier, err := strconv.ParseInt(question_id, 10, 64)
	if err != nil {
		return
	}
	_, err = app.responses.Insert(response, identifier)
	if err != nil {
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	//do a redirect
	app.sessionManager.Put(r.Context(), "flash", "poll submitted successfully!") //store it in flash
	http.Redirect(w, r, "/poll/success", http.StatusSeeOther)
}

func (app *application) pollSuccessShow(w http.ResponseWriter, r *http.Request) {
	// remove the entry from the session manager
	flash := app.sessionManager.PopString(r.Context(), "flash")
	//render
	data := &templateData{ //putting flash into template data
		Flash: flash,
	}
	RenderTemplate(w, "success.page.tmpl", data)
}

func (app *application) pollCreateShow(w http.ResponseWriter, r *http.Request) {
	//display the input box
	RenderTemplate(w, "poll.create.page.tmpl", nil)
}

func (app *application) pollCreateSubmit(w http.ResponseWriter, r *http.Request) {
	// add the question to the datastore
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	question := r.PostForm.Get("new_question")
	_, err = app.questions.Insert(question)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) optionsCreateShow(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "options.create.page.tmpl", nil)
}

func (app *application) optionsCreateSubmit(w http.ResponseWriter, r *http.Request) {
	// get the four options
	r.ParseForm()
	option_1 := r.PostForm.Get("option_1")
	option_2 := r.PostForm.Get("option_2")
	option_3 := r.PostForm.Get("option_3")
	option_4 := r.PostForm.Get("option_4")
	// save the options
	_, err := app.options.Insert(option_1, option_2, option_3, option_4)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	flash := app.sessionManager.PopString(r.Context(), "flash")
	//render
	data := &templateData{ //putting flash into template data
		Flash: flash,
	}
	RenderTemplate(w, "signup.page.tmpl", data)
}

func (app *application) userSignupSubmit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.PostForm.Get("name") //"name" is the name of the form
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	//lets write the data to the table
	err := app.users.Insert(name, email, password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			RenderTemplate(w, "signup.page.tmpl", nil)
		}
		app.sessionManager.Put(r.Context(), "flash", "Signup Was successful")
		//http.Redirect(w, r, "/users/login", http.StatusSeeOther)
	}
}
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
func (app *application) userLoginSubmit(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
func (app *application) userLogoutSubmit(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
