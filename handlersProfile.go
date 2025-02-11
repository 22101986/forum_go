package main

import (
	"forum/methods"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	loggedIn := LoggedInVerif(r) // verify if the cookie is setup with a session token
	DuplicateLog(loggedIn, w, r) // verify if the cookie is unique (handle double connection)

	user := &methods.User{} // get ID and Username of connected user
	if loggedIn {
		cookie, err := r.Cookie("session_token")
		ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

		userUUID := sessions[cookie.Value]

		err = sSquared.Users.DB.QueryRow("SELECT id, name, email, picture, role_id, is_external FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Email, &user.Picture, &user.Role, &user.IsExternal)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")
	}

	data := struct { //struct for execute and html
		ErrorMessage string
		ErrorBool    bool
		LoggedIn     bool
		User         *methods.User
		AllNotifs    []*methods.Notif
		NbrNotif     int
		AskedModo    bool
	}{
		ErrorMessage: "",
		ErrorBool:    false,
		LoggedIn:     loggedIn,
		User:         user,
		AllNotifs:    sData.AllNotifs,
		NbrNotif:     NumberOfNotif(user, sData.AllNotifs),
		AskedModo:    false,
	}
	var uselessID int
	err := sSquared.Notifs.DB.QueryRow("SELECT id FROM notifs WHERE user_id_from = ? AND type_id = ? ", user.ID, NotifAskModo).Scan(&uselessID)
	if err == nil {
		data.AskedModo = true
	}

	if r.Method == http.MethodGet {
		t, err := template.ParseFiles(`templates/myProfile.html`)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")

		t.Execute(w, data)
	} else if r.Method == http.MethodPost {
		r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
		err := r.ParseMultipartForm(5 << 20)
		ErrDiffNil(err, w, r, http.StatusBadRequest, "bad request")

		var newName string
		if user.IsExternal == BasicUser {
			if !VerifyContent(r.FormValue("name")) { // get content from form (comment)
				newName = r.FormValue("name")
			} else {
				ErrProfile(w, r, "can't have an empty username", "templates/myProfile.html", loggedIn, user, sData.AllNotifs, data.NbrNotif, data.AskedModo)
				return
			}
		}

		newMail := r.FormValue("email")
		newPassword := r.FormValue("password")
		cPassword := r.FormValue("cpassword")
		var newImage []byte
		if file, fileHeader, err := r.FormFile("image"); err == nil {
			newImage = GetImg(w, file, fileHeader)
		} else if err != http.ErrMissingFile {
			ErrorHandler(w, r, http.StatusBadRequest, "Unable to get the file")
		}

		var newEmail string
		if user.IsExternal == BasicUser {
			if match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, newMail); !match { // verify validity of email's format
				ErrProfile(w, r, "invalid email format", "templates/myProfile.html", loggedIn, user, sData.AllNotifs, data.NbrNotif, data.AskedModo)
				return
			} else {
				newEmail = newMail
			}

			if newPassword != cPassword {
				ErrProfile(w, r, "passwords not identical", "templates/myProfile.html", loggedIn, user, sData.AllNotifs, data.NbrNotif, data.AskedModo)
				return
			}

			// var userID int64
			if newName != user.Name {
				_, err = sSquared.Users.EditProfile(user.ID, newName, `UPDATE users SET name = ? WHERE id = ?`, false, false, false)
				if err != nil {
					if err.Error() == "UNIQUE constraint failed: users.name" { // display error of already used username
						ErrProfile(w, r, "username already used", `templates/myProfile.html`, loggedIn, user, sData.AllNotifs, data.NbrNotif, data.AskedModo)
						return
					}
				}
			}
			if newEmail != user.Email {
				_, err = sSquared.Users.EditProfile(user.ID, newEmail, `UPDATE users SET email = ? WHERE id = ?`, false, false, false)
				if err != nil {
					if err.Error() == "UNIQUE constraint failed: users.email" { // display error of already used email
						ErrProfile(w, r, "email already used", `templates/myProfile.html`, loggedIn, user, sData.AllNotifs, data.NbrNotif, data.AskedModo)
						return
					}
				}
			}
			if newPassword != "" {
				_, err = sSquared.Users.EditProfile(user.ID, newPassword, `UPDATE users SET password = ? WHERE id = ?`, true, false, false)
				ErrDiffNil(err, w, r, http.StatusInternalServerError, "edit password unavailable")
			}
		}
		if newImage != nil {
			_, err = sSquared.Users.EditProfile(user.ID, string(newImage), `UPDATE users SET picture = ? WHERE id = ?`, false, true, false)
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "edit image unavailable")
		}

		ClearStructs(sData)
		FillingStruct(w, r, sSquared, sData)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func DeleteProfileHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	loggedIn := LoggedInVerif(r)
	DuplicateLog(loggedIn, w, r)

	user := &methods.User{} // get ID and Username of connected user
	if loggedIn {
		cookie, err := r.Cookie("session_token")
		ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

		userUUID := sessions[cookie.Value]

		err = sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id, is_external FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Picture, &user.Role, &user.IsExternal)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")
	}

	err := sSquared.Users.UserFakeDeletion(user.ID)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "error deleting account")

	ClearStructs(sData)
	FillingStruct(w, r, sSquared, sData)

	http.Redirect(w, r, "/logout", http.StatusSeeOther)
}

func PostProfileHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
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
	var myPosts []*methods.Post

	for i := len(sData.AllPosts) - 1; i >= 0; i-- { // get all post from a precise ID (here the one from the connected user) and invert them
		if sData.AllPosts[i].User.ID == user.ID {
			myPosts = append(myPosts, sData.AllPosts[i])
		}
	}

	for _, post := range myPosts { // add bool for color on likes/dislikes
		if len(post.Likes) == 0 {
			post.IsLike = false
		} else {
			for _, like := range post.Likes {
				if like.User.ID == user.ID {
					post.IsLike = true
					break
				} else {
					post.IsLike = false
				}
			}
		}
		if len(post.Dislikes) == 0 {
			post.IsDislike = false
		} else {
			for _, dislike := range post.Dislikes {
				if dislike.User.ID == user.ID {
					post.IsDislike = true
					break
				} else {
					post.IsDislike = false
				}
			}
		}
	}

	data := struct { //struct for execute and html
		LoggedIn bool
		User     *methods.User
		Posts    []*methods.Post
		NbrNotif int
	}{
		LoggedIn: loggedIn,
		User:     user,
		Posts:    myPosts,
		NbrNotif: NumberOfNotif(user, sData.AllNotifs),
	}

	t, err := template.ParseFiles(`templates/myPosts.html`)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")

	t.Execute(w, data)

}
func LikeProfileHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	loggedIn := LoggedInVerif(r) // verify if the cookie is setup with a session token
	DuplicateLog(loggedIn, w, r) // verify if the cookie is unique (handle double connection)

	user := &methods.User{} // get struct of connected user
	if loggedIn {
		cookie, err := r.Cookie("session_token")
		ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

		userUUID := sessions[cookie.Value]

		err = sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Picture, &user.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")
	}

	var myLikes []*methods.Post
	var myLikesInv []*methods.Post

	rowPost, err := sSquared.Users.DB.Query("SELECT post_id FROM likepost where user_id = ?", user.ID) // get all posts liked by the connected user
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "posts not found")
	for rowPost.Next() {
		likedPost := &methods.Post{}
		likedUser := &methods.User{}
		err = rowPost.Scan(&likedPost.ID)
		if err != nil {
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting post's ID")
			continue // Passez au suivant en cas d'erreur
		}
		if err = sSquared.Posts.DB.QueryRow("SELECT title, content, date, user_id FROM posts WHERE id = ?", likedPost.ID).Scan(&likedPost.Title, &likedPost.Content, &likedPost.Date, &likedUser.ID); err != nil {
			log.Fatal(err)
		}
		if err = sSquared.Users.DB.QueryRow("SELECT name, picture FROM users WHERE id = ?", likedUser.ID).Scan(&likedUser.Name, &likedUser.Picture); err != nil {
			log.Fatal(err)
		}
		likedPost.User = likedUser

		rowRel, err := sSquared.Posts.DB.Query("SELECT cat_id FROM catpostrel WHERE post_id = ?", likedPost.ID)
		if err != nil {
			log.Fatal(err)
		}
		for rowRel.Next() {
			newCat := &methods.Categories{}
			if err := rowRel.Scan(&newCat.ID); err != nil {
				log.Fatal(err)
			}
			if err = sSquared.Posts.DB.QueryRow("SELECT name FROM categories WHERE id = ?", newCat.ID).Scan(&newCat.Name); err != nil {
				log.Fatal(err)
			}
			likedPost.Cats = append(likedPost.Cats, newCat)
		}
		rowRel.Close()

		likedPost.LenCat = len(likedPost.Cats) - 1

		// get likes/dislikes for the posts liked
		rowLike, err := sSquared.Likes.DB.Query("SELECT id, user_id FROM likepost WHERE post_id = ?", likedPost.ID)
		if err != nil {
			log.Fatal(err)
		}
		for rowLike.Next() {
			newLike := &methods.Like{}
			newUser := &methods.User{}
			if err := rowLike.Scan(&newLike.ID, &newUser.ID); err != nil {
				log.Fatal(err)
			}
			if err = sSquared.Users.DB.QueryRow("SELECT id, name, picture FROM users WHERE id = ?", newUser.ID).Scan(&newUser.ID, &newUser.Name, &newUser.Picture); err != nil {
				log.Fatal(err)
			}
			newLike.User = newUser
			likedPost.Likes = append(likedPost.Likes, newLike)
		}
		rowLike.Close()

		rowDislike, err := sSquared.Users.DB.Query("SELECT id, user_id FROM dislikepost WHERE post_id = ?", likedPost.ID)
		if err != nil {
			log.Fatal(err)
		}
		for rowDislike.Next() {
			newDislike := &methods.Dislike{}
			newUser := &methods.User{}
			if err := rowDislike.Scan(&newDislike.ID, &newUser.ID); err != nil {
				log.Fatal(err)
			}
			if err = sSquared.Users.DB.QueryRow("SELECT id, name, picture FROM users WHERE id = ?", newUser.ID).Scan(&newUser.ID, &newUser.Name, &newUser.Picture); err != nil {
				log.Fatal(err)
			}
			newDislike.User = newUser
			likedPost.Dislikes = append(likedPost.Dislikes, newDislike)
		}
		rowDislike.Close()

		if err != nil {
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting post's data")
			continue
		}
		myLikes = append(myLikes, likedPost)
	}
	rowPost.Close()

	for i := len(myLikes) - 1; i >= 0; i-- { // Invert post list liked by the user
		myLikesInv = append(myLikesInv, myLikes[i])
	}

	for _, post := range myLikesInv { // add bool for color on likes/dislikes
		if len(post.Likes) == 0 {
			post.IsLike = false
		} else {
			for _, like := range post.Likes {
				if like.User.ID == user.ID {
					post.IsLike = true
					break
				} else {
					post.IsLike = false
				}
			}
		}
		if len(post.Dislikes) == 0 {
			post.IsDislike = false
		} else {
			for _, dislike := range post.Dislikes {
				if dislike.User.ID == user.ID {
					post.IsDislike = true
					break
				} else {
					post.IsDislike = false
				}
			}
		}
	}

	data := struct { //struct for execute and html
		LoggedIn bool
		User     *methods.User
		Posts    []*methods.Post
		NbrNotif int
	}{
		LoggedIn: loggedIn,
		User:     user,
		Posts:    myLikesInv,
		NbrNotif: NumberOfNotif(user, sData.AllNotifs),
	}
	t, err := template.ParseFiles(`templates/myLikes.html`)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")

	t.Execute(w, data)
}
