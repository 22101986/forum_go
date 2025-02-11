package main

import (
	"fmt"
	"forum/methods"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func DetailPostHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	loggedIn := LoggedInVerif(r) // verify if the cookie is setup with a session token
	DuplicateLog(loggedIn, w, r) // verify if the cookie is unique (handle double connection)

	user := &methods.User{} // get ID and Username of connected user
	if loggedIn {
		cookie, err := r.Cookie("session_token")
		ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

		userUUID := sessions[cookie.Value]

		err = sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Picture, &user.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")
	}

	postID := r.URL.Query().Get("ID") // get the id of the post
	var postIDint int                 //initiate variable postIDint to store the id
	if postID == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "no ID in URL")
		return
	}
	postIDint, err := strconv.Atoi(postID) //convert the id to int and store it in the proper variable
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")
	if postIDint < 1 { //handle the case where the id's post is not existing
		ErrorHandler(w, r, http.StatusNotFound, "Page not found")
		return
	}

	var post *methods.Post //get post in the data struct
	for _, p := range sData.AllPosts {
		if p.ID == postIDint {
			post = p
			break
		}
	}

	if post == nil { //handle the case where the post is not existing
		ErrorHandler(w, r, http.StatusNotFound, "Page not found")
		return
	}

	if len(post.Likes) == 0 { //update data struct for like/dislike for html (for the post)
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
	for _, comment := range post.Comments { //update data struct for like/dislike for html (for each com)
		if len(comment.Likes) == 0 {
			comment.IsLike = false
		} else {
			for _, like := range comment.Likes {
				if like.User.ID == user.ID {
					comment.IsLike = true
					break
				} else {
					comment.IsLike = false
				}
			}
		}
		if len(comment.Dislikes) == 0 {
			comment.IsDislike = false
		} else {
			for _, dislike := range comment.Dislikes {
				if dislike.User.ID == user.ID {
					comment.IsDislike = true
					break
				} else {
					comment.IsDislike = false
				}
			}
		}
	}

	data := struct { //struct for execute and html
		ErrorMessage string
		ErrorBool    bool
		LoggedIn     bool
		User         *methods.User
		Post         *methods.Post
		NbrNotif     int
	}{
		ErrorMessage: "",
		ErrorBool:    false,
		LoggedIn:     loggedIn,
		User:         user,
		Post:         post,
		NbrNotif:     NumberOfNotif(user, sData.AllNotifs),
	}

	if r.Method == http.MethodGet {
		t, err := template.ParseFiles(`templates/detailPost.html`)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")
		t.Execute(w, data)
	} else if r.Method == http.MethodPost {
		t, err := template.ParseFiles(`templates/detailPost.html`)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")

		err = r.ParseForm()
		ErrDiffNil(err, w, r, http.StatusBadRequest, "bad request")

		var content string
		if !VerifyContent(r.FormValue("content")) { // get content from form (comment)
			content = r.FormValue("content")
		} else {
			ErrCom1(w, r, "can't post an empty comment", `templates/detailPost.html`, data.LoggedIn, data.User, data.NbrNotif, data.Post)
			return
		}

		cookie, err := r.Cookie("session_token") // get struct of connected user
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			ErrorHandler(w, r, http.StatusBadRequest, "gathering cookie error")
			return
		}

		userUUID := sessions[cookie.Value]
		user := &methods.User{}
		err = sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Picture, &user.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")

		comment := &methods.Comment{ // create a comment
			Content: content,
			Date:    time.Now().Format("2006-01-02 15:04:05"),
			User:    user,
			PostID:  post.ID,
		}
		commentID, err := sSquared.Comments.InsertInComments(comment) // Add it in database
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")
		(*comment).ID = int(commentID)
		post.Comments = append(post.Comments, comment) // Add it in the data struct

		newNotif := &methods.Notif{
			ID:       0,
			Date:     time.Now().Format("2006-01-02 15:04:05"),
			Type:     5,
			UserTo:   post.User,
			UserFrom: user,
			Post:     post,
			Comment:  comment,
		}
		newNotif.ID, _ = sSquared.Notifs.InsertInNotifs(newNotif, true)
		if newNotif.ID != 0 {
			sData.AllNotifs = append(sData.AllNotifs, newNotif)
		}

		t.Execute(w, data)

	} else {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "method not allowed")
	}
	sSquared.Notifs.DeleteAllNotifPost(postIDint, user.ID)
	ClearStructs(sData)
	FillingStruct(w, r, sSquared, sData)
}

func NewPostHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {

	loggedIn := LoggedInVerif(r) // verify if the cookie is setup with a session token
	DuplicateLog(loggedIn, w, r) // verify if the cookie is unique (handle double connection)

	user := &methods.User{} // get username of connected user
	if loggedIn {
		cookie, err := r.Cookie("session_token")
		ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

		userUUID := sessions[cookie.Value]

		err = sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Picture, &user.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")
	}

	data := struct { // struct for execute and html
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
		ErrorMessage: "",
		ErrorBool:    false,
		LoggedIn:     loggedIn,
		User:         user,
		AllCats:      sData.AllCats,
		Post:         nil,
		Title:        "",
		Content:      "",
		NbrNotif:     NumberOfNotif(user, sData.AllNotifs),
	}

	if r.Method == http.MethodGet {
		t, err := template.ParseFiles(`templates/newPost.html`)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")

		t.Execute(w, data)
	} else if r.Method == http.MethodPost {
		r.Body = http.MaxBytesReader(w, r.Body, 20<<20)
		err := r.ParseMultipartForm(5 << 20)
		if err != nil {
			ErrDiffNil(err, w, r, http.StatusBadRequest, "bad request")
			return
		}

		var image []byte
		var title string
		var content string
		if !VerifyContent(r.FormValue("title")) { // get content from form (comment)
			title = r.FormValue("title")
		} else {
			ErrPost(w, r, "can't post with an empty title", `templates/newPost.html`, data.LoggedIn, data.User, data.AllCats, title, content, data.NbrNotif, data.Post)
			return
		}
		if !VerifyContent(r.FormValue("content")) { // get content from form (comment)
			content = r.FormValue("content")
		} else {
			ErrPost(w, r, "can't post with an empty post", `templates/newPost.html`, data.LoggedIn, data.User, data.AllCats, title, content, data.NbrNotif, data.Post)
			return
		}
		if file, fileHeader, err := r.FormFile("image"); err == nil {
			image = GetImg(w, file, fileHeader)
		} else if err != http.ErrMissingFile {
			ErrorHandler(w, r, http.StatusBadRequest, "unable to get the file")
		}

		values := make(map[int]string)
		for _, cats := range sData.AllCats {
			values[cats.ID] = cats.Name
		}

		atLeastOneCat := false //verify there's at least one category for the new post
		for _, catLine := range values {
			catégorie := r.FormValue(catLine)
			if catégorie != "" {
				atLeastOneCat = true
				break
			}
		}
		if !atLeastOneCat { // send error if not the case
			ErrPost(w, r, "need at least one category", "templates/newPost.html", data.LoggedIn, data.User, data.AllCats, title, content, data.NbrNotif, data.Post)
			return
		}

		cookie, err := r.Cookie("session_token") // get struct of connected user
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			ErrorHandler(w, r, http.StatusBadRequest, "gathering cookie error")
			return
		}

		userUUID := sessions[cookie.Value]
		user := &methods.User{}
		err = sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Picture, &user.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")

		post := &methods.Post{ // create a post
			Title:   title,
			Content: content,
			Date:    time.Now().Format("2006-01-02 15:04:05"),
			User:    user,
			Blob:    image,
		}
		commentID, err := sSquared.Posts.InsertInPosts(post) // Add it in database
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")
		(*post).ID = int(commentID)
		if image != nil {
			err = sSquared.Blobs.InsertInBlob(image, post.ID)
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")
		}
		sData.AllPosts = append(sData.AllPosts, post) // Add it in data struct

		for keyCat, catLine := range values { //compare the categories checked and the ones from database
			postCat := &methods.Categories{
				ID:   keyCat,
				Name: catLine,
			}
			catégorie := r.FormValue(catLine)
			if catégorie != "" {
				err = sSquared.Posts.InsertInRel(keyCat, commentID) // add in database (rel tab)
				ErrDiffNil(err, w, r, http.StatusBadRequest, "error in insert relation")

				post.Cats = append(post.Cats, postCat) // add it in data struct
			}
		}
		post.LenCat = len(post.Cats) - 1 // add the len-1 in data struct (for html)

		ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// function for delete a post if you are author, admin or moderator
func DeletePostHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	loggedIn := LoggedInVerif(r)
	DuplicateLog(loggedIn, w, r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	id := r.URL.Query().Get("ID")
	if id == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "Missing post ID")
		return
	}
	user := &methods.User{}
	cookie, err := r.Cookie("session_token")
	ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

	userUUID := sessions[cookie.Value]
	err = sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Picture, &user.Role)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "no user")
	for i, post := range sData.AllPosts {
		if strconv.Itoa(post.ID) == id {
			if post == nil || (post.User.ID != user.ID && user.Role != RoleAdmin && user.Role != RoleModo) {
				ErrorHandler(w, r, http.StatusForbidden, "You are not allowed to delete this post")
				return
			}
			sData.AllPosts = append(sData.AllPosts[:i], sData.AllPosts[i+1:]...)
			break
		}
	}

	err = sSquared.Posts.DeleteInPosts(id)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "Failed to delete post")

	_, err = sSquared.Notifs.DB.Exec("DELETE FROM notifs WHERE post_id = ?", id)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "Failed to delete notif")

	ClearStructs(sData)
	FillingStruct(w, r, sSquared, sData)

	http.Redirect(w, r, "/", http.StatusFound)
}

// function for delete a comment if you are its author or admin
func DeleteCommentHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	loggedIn := LoggedInVerif(r)
	DuplicateLog(loggedIn, w, r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	id := r.URL.Query().Get("ID")
	if id == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "Missing post ID")
		return
	}
	user := &methods.User{}
	cookie, err := r.Cookie("session_token")
	ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

	userUUID := sessions[cookie.Value]
	err = sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Picture, &user.Role)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "no user found")

	var postID int
	for _, post := range sData.AllPosts {
		for i, comment := range post.Comments {
			if strconv.Itoa(comment.ID) == id {
				if comment == nil || (comment.User.ID != user.ID && user.Role != RoleAdmin) {
					ErrorHandler(w, r, http.StatusForbidden, "You are not allowed to delete this comment")
					return
				}
				postID = post.ID
				post.Comments = append(post.Comments[:i], post.Comments[i+1:]...)
				break
			}
		}
	}

	err = sSquared.Comments.DeleteInComments(id)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "Failed to delete comment")

	_, err = sSquared.Notifs.DB.Exec("DELETE FROM notifs WHERE comments_id = ?", id)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "Failed to delete notif")

	ClearStructs(sData)
	FillingStruct(w, r, sSquared, sData)

	http.Redirect(w, r, fmt.Sprintf("/detailPost?ID=%d", postID), http.StatusFound)
}

// function for edit a post if you are its author
func EditPostHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	loggedIn := LoggedInVerif(r)
	DuplicateLog(loggedIn, w, r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user := &methods.User{}
	cookie, err := r.Cookie("session_token")
	ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

	userUUID := sessions[cookie.Value]
	err = sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Picture, &user.Role)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "no user found")

	postID := r.URL.Query().Get("ID")
	if postID == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "no commentID")
		return
	}

	var post *methods.Post
	for _, p := range sData.AllPosts {
		if strconv.Itoa(p.ID) == postID {
			post = p
			break
		}
	}

	if post == nil || post.User.ID != user.ID {
		ErrorHandler(w, r, http.StatusForbidden, "You are not allowed to edit this comment")
		return
	}

	data := struct {
		Post      *methods.Post
		LoggedIn  bool
		User      *methods.User
		AllCats   []*methods.Categories
		Title     string
		Content   string
		ErrorBool bool
		ErrorMsg  string
		NbrNotif  int
	}{
		Post:      post,
		LoggedIn:  loggedIn,
		User:      user,
		AllCats:   sData.AllCats,
		Title:     post.Title,
		Content:   post.Content,
		ErrorBool: false,
		ErrorMsg:  "",
		NbrNotif:  NumberOfNotif(user, sData.AllNotifs),
	}

	if r.Method == http.MethodGet {
		t, err := template.ParseFiles(`templates/editPost.html`)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "Error")
		t.Execute(w, data)

	} else if r.Method == http.MethodPost {
		r.Body = http.MaxBytesReader(w, r.Body, 20<<20)
		err := r.ParseMultipartForm(5 << 20)
		ErrDiffNil(err, w, r, http.StatusBadRequest, "bad request")

		var image []byte
		var title string
		var content string
		if !VerifyContent(r.FormValue("title")) { // get content from form (comment)
			title = r.FormValue("title")
		} else {
			ErrPost(w, r, "can't post with an empty title", `templates/editPost.html`, data.LoggedIn, data.User, data.AllCats, title, content, data.NbrNotif, data.Post)
			return
		}
		if !VerifyContent(r.FormValue("content")) { // get content from form (comment)
			content = r.FormValue("content")
		} else {
			ErrPost(w, r, "can't post with an empty post", `templates/editPost.html`, data.LoggedIn, data.User, data.AllCats, title, content, data.NbrNotif, data.Post)
			return
		}
		if file, fileHeader, err := r.FormFile("image"); err == nil {
			image = GetImg(w, file, fileHeader)
		} else if err != http.ErrMissingFile {
			ErrorHandler(w, r, http.StatusBadRequest, "Unable to get the file")
		}

		query := `UPDATE posts SET title = ?, content = ? WHERE id = ?`
		_, err = sSquared.Posts.DB.Exec(query, title, content, postID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "Error updating post")

		blobQuery := `UPDATE blob SET picture = ? WHERE post_id = ?`
		_, err = sSquared.Posts.DB.Exec(blobQuery, image, postID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "Error updating post's picture")

		post.Title = title
		post.Content = content
		post.Cats = []*methods.Categories{}
		post.Blob = image

		_, err = sSquared.Users.DB.Exec(`DELETE FROM catpostrel WHERE post_id = ?`, post.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "failed to clear existing categories")

		values := make(map[int]string)
		for _, cats := range sData.AllCats {
			values[cats.ID] = cats.Name
		}

		atLeastOneCat := false //verify there's at least one category for the new post
		for _, catLine := range values {
			catégorie := r.FormValue(catLine)
			if catégorie != "" {
				atLeastOneCat = true
			}
		}
		if !atLeastOneCat { // send error if not the case
			ErrPost(w, r, "need at least one category", "templates/editPost.html", data.LoggedIn, data.User, data.AllCats, title, content, data.NbrNotif, data.Post)
			return
		}
		for keyCat, catLine := range values { //compare the categories checked and the ones from database
			postCat := &methods.Categories{
				ID:   keyCat,
				Name: catLine,
			}
			catégorie := r.FormValue(catLine)
			if catégorie != "" {
				err = sSquared.Posts.InsertInRel(keyCat, int64(post.ID)) // add in database (rel tab)
				ErrDiffNil(err, w, r, http.StatusInternalServerError, " failed the insertion in rel")
				post.Cats = append(post.Cats, postCat) // add it in data struct
			}
		}
		post.LenCat = len(post.Cats) - 1

		ClearStructs(sData)
		FillingStruct(w, r, sSquared, sData)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// function for edit a comment if you are its author
func EditCommentHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	loggedIn := LoggedInVerif(r)
	DuplicateLog(loggedIn, w, r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	cookie, err := r.Cookie("session_token")
	ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

	userUUID := sessions[cookie.Value]
	user := &methods.User{}
	err = sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Picture, &user.Role)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "no user found")

	commentID := r.URL.Query().Get("ID")
	if commentID == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "no comment ID")
		return
	}

	var comment *methods.Comment
	for _, post := range sData.AllPosts {
		for _, c := range post.Comments {
			if strconv.Itoa(c.ID) == commentID {
				comment = c
				break
			}
		}
		if comment != nil {
			break
		}
	}

	if comment == nil || comment.User.ID != user.ID {
		ErrorHandler(w, r, http.StatusForbidden, "You are not allowed to edit this comment")
		return
	}

	data := struct {
		LoggedIn  bool
		User      *methods.User
		Comment   *methods.Comment
		ErrorBool bool
		ErrorMsg  string
		NbrNotif  int
	}{
		LoggedIn:  loggedIn,
		User:      user,
		Comment:   comment,
		ErrorBool: false,
		ErrorMsg:  "",
		NbrNotif:  NumberOfNotif(user, sData.AllNotifs),
	}

	if r.Method == http.MethodGet {
		t, err := template.ParseFiles("templates/editComment.html")
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "Error parsing template")
		t.Execute(w, data)

	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		ErrDiffNil(err, w, r, http.StatusBadRequest, "bad request")

		var content string
		if !VerifyContent(r.FormValue("content")) { // get content from form (comment)
			content = r.FormValue("content")
		} else {
			ErrCom2(w, r, "can't post an empty comment", `templates/editComment.html`, data.LoggedIn, data.User, data.NbrNotif, data.Comment)
			return
		}

		_, err = sSquared.Users.DB.Exec("UPDATE comments SET content = ? WHERE id = ?", content, comment.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "Error updating comment")

		comment.Content = content

		http.Redirect(w, r, fmt.Sprintf("/detailPost?ID=%d", comment.PostID), http.StatusSeeOther)
	} else {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
