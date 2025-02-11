package main

import (
	"forum/methods"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

func ModoHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	loggedIn := LoggedInVerif(r)
	DuplicateLog(loggedIn, w, r)

	user := &methods.User{} // get ID and Username of connected user
	if loggedIn {
		cookie, err := r.Cookie("session_token")
		ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

		userUUID := sessions[cookie.Value]

		err = sSquared.Users.DB.QueryRow("SELECT id, uuid, name, picture, role_id FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.UUID, &user.Name, &user.Picture, &user.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")
	}
	userAdmin := &methods.User{}
	err := sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE role_id = ?", RoleAdmin).Scan(&userAdmin.ID, &userAdmin.Name, &userAdmin.Picture, &userAdmin.Role)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")

	// verifier qu'une notif existe déjà.
	newNotif := &methods.Notif{
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Type:     NotifAskModo,
		UserTo:   userAdmin,
		UserFrom: user,
		Post:     nil,
	}
	newNotif.ID, _ = sSquared.Notifs.InsertInNotifs(newNotif, false)
	if newNotif.ID != 0 {
		sData.AllNotifs = append(sData.AllNotifs, newNotif)
	}
	http.Redirect(w, r, "/myProfile", http.StatusSeeOther)
}

func ResponseHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) { // A modif pour les reponses de report
	loggedIn := LoggedInVerif(r)
	DuplicateLog(loggedIn, w, r)

	user := &methods.User{} // get ID and Username of connected user
	if loggedIn {
		cookie, err := r.Cookie("session_token")
		ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

		userUUID := sessions[cookie.Value]

		err = sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Picture, &user.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")
	}

	userFromUUID := r.URL.Query().Get("user")
	notifID := r.URL.Query().Get("notif")
	result := r.URL.Query().Get("result")
	postBool, _ := strconv.ParseBool(r.URL.Query().Get("postBool"))
	comBool, _ := strconv.ParseBool(r.URL.Query().Get("comBool"))

	userFrom := &methods.User{}
	err := sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE uuid = ?", userFromUUID).Scan(&userFrom.ID, &userFrom.Name, &userFrom.Picture, &userFrom.Role)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")

	oldNotif := &methods.Notif{Post: &methods.Post{}, Comment: &methods.Comment{}}
	oldNotif.ID, _ = strconv.Atoi(notifID)
	var sourceIDStr string
	err = sSquared.Notifs.DB.QueryRow("SELECT type_id FROM notifs WHERE id = ?", oldNotif.ID).Scan(&oldNotif.Type)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "notif typeless ?")
	if oldNotif.Type == NotifReportPost || oldNotif.Type == NotifReportCom {
		err = sSquared.Notifs.DB.QueryRow("SELECT post_id FROM notifs WHERE id = ?", oldNotif.ID).Scan(&oldNotif.Post.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "notif postless ?")
		sourceIDStr = strconv.Itoa(oldNotif.Post.ID)
	}
	if oldNotif.Type == NotifReportCom {
		err = sSquared.Notifs.DB.QueryRow("SELECT comments_id FROM notifs WHERE id = ?", oldNotif.ID).Scan(&oldNotif.Comment.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "notif comless ?")
		sourceIDStr = strconv.Itoa(oldNotif.Comment.ID)
	}

	newNotif := &methods.Notif{
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Type:     NotifAdminAccept,
		UserTo:   userFrom,
		UserFrom: user,
	}

	switch {
	case postBool && !comBool:
		if result == "accept" {
			newNotif.Type = NotifAcceptDeletion
			sSquared.Notifs.InsertInNotifs(newNotif, false)
			err = sSquared.Notifs.DeleteInNotifs(notifID)
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "error during notif deletion")
			ClearStructs(sData)
			FillingStruct(w, r, sSquared, sData)
			http.Redirect(w, r, "/deletePost?ID="+sourceIDStr, http.StatusSeeOther)
		} else {
			newNotif.Type = NotifRefuseDeletion
			sSquared.Notifs.InsertInNotifs(newNotif, false)
			err = sSquared.Notifs.DeleteInNotifs(notifID)
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "error during notif deletion")
			ClearStructs(sData)
			FillingStruct(w, r, sSquared, sData)
			http.Redirect(w, r, "/notifications", http.StatusSeeOther)
		}
	case !postBool && comBool:
		if result == "accept" {
			newNotif.Type = NotifAcceptDeletion
			sSquared.Notifs.InsertInNotifs(newNotif, false)
			err = sSquared.Notifs.DeleteInNotifs(notifID)
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "error during notif deletion")
			ClearStructs(sData)
			FillingStruct(w, r, sSquared, sData)
			http.Redirect(w, r, "/deleteComment?ID="+sourceIDStr, http.StatusSeeOther)
		} else {
			newNotif.Type = NotifRefuseDeletion
			sSquared.Notifs.InsertInNotifs(newNotif, false)
			err = sSquared.Notifs.DeleteInNotifs(notifID)
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "error during notif deletion")
			ClearStructs(sData)
			FillingStruct(w, r, sSquared, sData)
			http.Redirect(w, r, "/notifications", http.StatusSeeOther)
		}
	default:
		if userFrom.Role < RoleModo {
			if result == "accept" { //Change the user role to 3
				_, err = sSquared.Users.EditProfile(userFrom.ID, "3", `UPDATE users SET role_id = ? WHERE id = ?`, false, false, true)
				ErrDiffNil(err, w, r, http.StatusInternalServerError, "edit role error")

				newNotif.ID, _ = sSquared.Notifs.InsertInNotifs(newNotif, false)
				if newNotif.ID != 0 {
					sData.AllNotifs = append(sData.AllNotifs, newNotif)
				}
			} else if result == "refuse" {
				newNotif.Type = NotifAdminRefuse
				newNotif.ID, _ = sSquared.Notifs.InsertInNotifs(newNotif, false)
				if newNotif.ID != 0 {
					sData.AllNotifs = append(sData.AllNotifs, newNotif)
				}
			}
			err = sSquared.Notifs.DeleteInNotifs(notifID)
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "error during notif deletion")

			ClearStructs(sData)
			FillingStruct(w, r, sSquared, sData)
		}
		http.Redirect(w, r, "/notifications", http.StatusSeeOther)
	}
}

func ReportHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	loggedIn := LoggedInVerif(r)
	DuplicateLog(loggedIn, w, r)

	user := &methods.User{} // get ID and Username of connected user
	if loggedIn {
		cookie, err := r.Cookie("session_token")
		ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

		userUUID := sessions[cookie.Value]

		err = sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Picture, &user.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")
	}
	userAdmin := &methods.User{}
	err := sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE role_id = ?", RoleAdmin).Scan(&userAdmin.ID, &userAdmin.Name, &userAdmin.Picture, &userAdmin.Role)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")
	newNotif := &methods.Notif{
		ID:       0,
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Type:     NotifReportPost,
		UserTo:   userAdmin,
		UserFrom: user,
		Post:     &methods.Post{},
		Comment:  &methods.Comment{},
	}
	postIDStr := r.URL.Query().Get("postID")
	comIDStr := r.URL.Query().Get("comID")
	postID, _ := strconv.Atoi(postIDStr)
	if err = sSquared.Posts.DB.QueryRow("SELECT id, title, content, date FROM posts WHERE id = ?", postID).Scan(&newNotif.Post.ID, &newNotif.Post.Title, &newNotif.Post.Content, &newNotif.Post.Date); err != nil {
		log.Fatal(err)
	}
	if comIDStr != "" {
		comID, _ := strconv.Atoi(comIDStr)
		if err = sSquared.Posts.DB.QueryRow("SELECT id, content, date FROM comments WHERE id = ?", comID).Scan(&newNotif.Comment.ID, &newNotif.Comment.Content, &newNotif.Comment.Date); err != nil {
			log.Fatal(err)
		}
		newNotif.Type = NotifReportCom
	} else {
		newNotif.Comment = nil
	}

	newNotif.ID, _ = sSquared.Notifs.InsertInNotifs(newNotif, true)
	if newNotif.ID != 0 {
		sData.AllNotifs = append(sData.AllNotifs, newNotif)
	}

	http.Redirect(w, r, "/detailPost?ID="+postIDStr, http.StatusSeeOther)
}

func AdminHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	loggedIn := LoggedInVerif(r)
	DuplicateLog(loggedIn, w, r)

	user := &methods.User{} // get ID and Username of connected user
	if loggedIn {
		cookie, err := r.Cookie("session_token")
		ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

		userUUID := sessions[cookie.Value]

		err = sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Picture, &user.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")
	}

	var modosList []*methods.User

	rowUser, err := sSquared.Users.DB.Query("SELECT id, uuid, name, picture FROM users WHERE role_id = 3")
	if err != nil {
		log.Fatal(err)
	}
	for rowUser.Next() {
		newModo := &methods.User{}
		if err := rowUser.Scan(&newModo.ID, &newModo.UUID, &newModo.Name, &newModo.Picture); err != nil {
			log.Fatal(err)
		}
		modosList = append(modosList, newModo)
	}
	rowUser.Close()

	data := struct { //struct for execute and html
		LoggedIn bool
		User     *methods.User
		Modos    []*methods.User
		NbrNotif int
	}{
		LoggedIn: loggedIn,
		User:     user,
		Modos:    modosList,
		NbrNotif: NumberOfNotif(user, sData.AllNotifs),
	}

	t, err := template.ParseFiles(`templates/admin.html`)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "error parsing file")

	t.Execute(w, data)
}

func DemoteHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	loggedIn := LoggedInVerif(r)
	DuplicateLog(loggedIn, w, r)

	user := &methods.User{} // get ID and Username of connected user
	if loggedIn {
		cookie, err := r.Cookie("session_token")
		ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

		userUUID := sessions[cookie.Value]

		err = sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Picture, &user.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "admin not found")
	}

	userID, _ := strconv.Atoi(r.URL.Query().Get("user"))

	modo := &methods.User{}
	err := sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE id = ?", userID).Scan(&modo.ID, &modo.Name, &modo.Picture, &modo.Role)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "modo not found")

	if user.Role == RoleAdmin && modo.Role == RoleModo {
		_, err = sSquared.Users.EditProfile(modo.ID, "2", `UPDATE users SET role_id = ? WHERE id = ?`, false, false, true)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "edit role error")

		ClearStructs(sData)
		FillingStruct(w, r, sSquared, sData)
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
