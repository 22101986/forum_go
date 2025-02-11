package main

import (
	"forum/methods"
	"html/template"
	"net/http"
)

// refacto because err != nil was frustrating
func ErrDiffNil(err error, w http.ResponseWriter, r *http.Request, code int, msg string) {
	if err != nil {
		ErrorHandler(w, r, code, msg)
		return
	}
}

// Func to execute register or login template with new data (error data)
func ErrLog(w http.ResponseWriter, r *http.Request, code int, msg string, source string) {
	errData := struct {
		ErrorMessage string
	}{
		ErrorMessage: msg,
	}
	t, _ := template.ParseFiles(source)
	t.Execute(w, errData)
}

// Func to execute detailpost template with new data (error data)
func ErrCom1(w http.ResponseWriter, r *http.Request, msg string, source string, log bool, user *methods.User, nbrNotif int, post *methods.Post) {
	errData := struct {
		ErrorMessage string
		ErrorBool    bool
		LoggedIn     bool
		User         *methods.User
		Post         *methods.Post
		NbrNotif     int
	}{
		ErrorMessage: msg,
		ErrorBool:    true,
		LoggedIn:     log,
		User:         user,
		Post:         post,
		NbrNotif:     nbrNotif,
	}
	t, _ := template.ParseFiles(source)
	t.Execute(w, errData)
}

// Func to execute editcom template with new data (error data)
func ErrCom2(w http.ResponseWriter, r *http.Request, msg string, source string, log bool, user *methods.User, nbrNotif int, comment *methods.Comment) {
	errData := struct {
		ErrorMessage string
		ErrorBool    bool
		LoggedIn     bool
		User         *methods.User
		Comment      *methods.Comment
		NbrNotif     int
	}{
		ErrorMessage: msg,
		ErrorBool:    true,
		LoggedIn:     log,
		User:         user,
		Comment:      comment,
		NbrNotif:     nbrNotif,
	}
	t, _ := template.ParseFiles(source)
	t.Execute(w, errData)
}

// Func to execute newpost template with new data (error data)
func ErrPost(w http.ResponseWriter, r *http.Request, msg string, source string, log bool, user *methods.User, cats []*methods.Categories, title, content string, nbrNotif int, post *methods.Post) {
	errData := struct {
		ErrorMessage string
		ErrorBool    bool
		LoggedIn     bool
		User         *methods.User
		AllCats      []*methods.Categories
		Post         *methods.Post
		Title        string
		Content      string
		NbrNotif     int
	}{
		ErrorMessage: msg,
		ErrorBool:    true,
		LoggedIn:     log,
		User:         user,
		AllCats:      cats,
		Post:         post,
		Title:        title,
		Content:      content,
		NbrNotif:     nbrNotif,
	}
	t, _ := template.ParseFiles(source)
	t.Execute(w, errData)
}

// Func to execute profile template with new data (error data)
func ErrProfile(w http.ResponseWriter, r *http.Request, msg string, source string, log bool, user *methods.User, notif []*methods.Notif, nbrNotif int, askedModo bool) {
	errData := struct {
		ErrorMessage string
		ErrorBool    bool
		LoggedIn     bool
		User         *methods.User
		AllNotifs    []*methods.Notif
		NbrNotif     int
		AskedModo    bool
	}{
		ErrorMessage: msg,
		ErrorBool:    true,
		LoggedIn:     log,
		User:         user,
		AllNotifs:    notif,
		NbrNotif:     nbrNotif,
		AskedModo:    askedModo,
	}
	t, _ := template.ParseFiles(source)
	t.Execute(w, errData)
}
