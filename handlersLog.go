package main

import (
	"forum/methods"
	"html/template"
	"net/http"
	"regexp"

	"github.com/gofrs/uuid"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	if r.Method == http.MethodGet { // display page
		t, err := template.ParseFiles(`templates/register.html`)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error parsing file")

		t.Execute(w, nil)

	} else if r.Method == http.MethodPost { // handle registration
		err := r.ParseForm()
		ErrDiffNil(err, w, r, http.StatusBadRequest, "bad request")

		newUUID, err := uuid.NewV4() // generate an UUID
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")

		var name string
		if !VerifyContent(r.FormValue("name")) { // get content from form (comment)
			name = r.FormValue("name")
		} else {
			ErrLog(w, r, http.StatusBadRequest, "username can't be empty", `templates/register.html`)
			return
		}

		var email string
		mail := r.FormValue("email") // get string from form

		if match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, mail); !match { // verify validity of email's format
			ErrLog(w, r, http.StatusBadRequest, "invalid email format", "templates/register.html")
			return
		} else {
			email = mail
		}

		userID, err := sSquared.Users.InsertInUser(name, email, r.FormValue("password"), newUUID) // create user in database with data from form
		if err != nil {
			if err.Error() == "UNIQUE constraint failed: users.name" { // display error of already used username
				ErrLog(w, r, http.StatusBadRequest, "username already used", `templates/register.html`)
				return
			} else if err.Error() == "UNIQUE constraint failed: users.email" { // display error of already used email
				ErrLog(w, r, http.StatusBadRequest, "email already used", `templates/register.html`)
				return
			}
		}

		user := &methods.User{ // create user struc for data struct
			ID:       int(userID),
			UUID:     newUUID.String(),
			Name:     r.FormValue("name"),
			Email:    email,
			Password: r.FormValue("password"),
			Picture:  nil,
			Role:     RoleUser,
		}

		sData.AllUsers = append(sData.AllUsers, user)

		SessionGen(w, user, false) // create cookie for the session

		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	if r.Method == http.MethodGet { // display page
		t, err := template.ParseFiles(`templates/login.html`)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")

		t.Execute(w, nil)
	} else if r.Method == http.MethodPost { // handle log authentification
		err := r.ParseForm()
		ErrDiffNil(err, w, r, http.StatusBadRequest, "bad request")

		password := r.FormValue("password") // get password and email from form
		email := r.FormValue("email")

		rememberMe := r.FormValue("remember") == "remember" // get remember me option from form

		uuid, err := sSquared.Users.Authenticate(email, password) // verify credentials in database and get uuid from user
		if err != nil {
			ErrLog(w, r, http.StatusBadRequest, "invalid credentials", `templates/login.html`) // return to invalid credentials if error
			return
		}

		user := &methods.User{
			UUID: uuid,
		}
		for token := range sessions { // delete cookie if already connected
			if sessions[token] == user.UUID {
				delete(sessions, token)
				break
			}
		}
		SessionGen(w, user, rememberMe) // recreate cookie for the session

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	cookie, err := r.Cookie("session_token") // validity of cookie
	if err != nil {
		if err == http.ErrNoCookie { // invalid cookie redirect to login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ErrorHandler(w, r, http.StatusBadRequest, "gathering cookie error")
		return
	}

	sessionToken := cookie.Value // delete cookie and close session
	delete(sessions, sessionToken)
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
