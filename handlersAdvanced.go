package main

import (
	"forum/methods"
	"html/template"
	"net/http"
)

func NotifHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	loggedIn := LoggedInVerif(r) // verify if the cookie is setup with a session token
	DuplicateLog(loggedIn, w, r) // verify if the cookie is unique (handle double connection)

	user := &methods.User{} // get struct of connected user
	if loggedIn {
		cookie, err := r.Cookie("session_token")
		ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

		userUUID := sessions[cookie.Value]

		err = sSquared.Users.DB.QueryRow("SELECT id, uuid, name, email, picture, role_id FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.UUID, &user.Name, &user.Email, &user.Picture, &user.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")
	}

	tabNotifs := []*methods.Notif{}

	for _, notif := range sData.AllNotifs {
		if notif.UserTo.ID == user.ID || (notif.Type == NotifAskModo && user.Role == RoleAdmin) {
			tabNotifs = append(tabNotifs, notif)
		}
	}
	data := struct {
		LoggedIn  bool
		User      *methods.User
		AllNotifs []*methods.Notif
	}{
		LoggedIn:  loggedIn,
		User:      user,
		AllNotifs: tabNotifs,
	}

	sSquared.Notifs.DB.Exec("DELETE FROM notifs WHERE type_id IN (9, 10, 11, 12) AND user_id_to = ?", user.ID)

	t, err := template.ParseFiles(`templates/notifications.html`)
	ErrDiffNil(err, w, r, http.StatusNotFound, "index.html not found")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(w, data)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")

	/*ClearStructs(sData)
	FillingStruct(sSquared, sData)*/
	tabTemp := []*methods.Notif{}
	for _, notif := range sData.AllNotifs {
		if (notif.Type != 9 && notif.Type != 10 && notif.Type != 11 && notif.Type != 12) || notif.UserTo.ID != user.ID {
			tabTemp = append(tabTemp, notif)
		}
	}
	sData.AllNotifs = tabTemp
}
