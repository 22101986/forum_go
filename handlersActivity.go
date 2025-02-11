package main

import (
	"forum/methods"
	"html/template"
	"net/http"
	"sort"
)

func ActivityProfileHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
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

	var myActivity []*methods.Activity

	rowPost, err := sSquared.Posts.DB.Query("SELECT id, title, content, date FROM posts WHERE user_id = ?", user.ID)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting posts' data")
	for rowPost.Next() {
		newActivity := &methods.Activity{}
		newPost := &methods.Post{}
		err := rowPost.Scan(&newPost.ID, &newPost.Title, &newPost.Content, &newPost.Date)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting posts' data")
		newPost.User = user
		newActivity.Post = newPost
		newActivity.Type = ActivityPost
		myActivity = append(myActivity, newActivity)
	}
	err = rowPost.Err()
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting posts' data")

	rowPost.Close()

	rowCom, err := sSquared.Comments.DB.Query("SELECT id, content, date, post_id FROM comments WHERE user_id = ?", user.ID)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting comments' data")
	for rowCom.Next() {
		newActivity := &methods.Activity{}
		newPost := &methods.Post{}
		newUser := &methods.User{}
		newCom := &methods.Comment{}
		rowCom.Scan(&newCom.ID, &newCom.Content, &newCom.Date, &newPost.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting comments' data")
		err = sSquared.Posts.DB.QueryRow("SELECT title, content, date, user_id FROM posts WHERE id = ?", newPost.ID).Scan(&newPost.Title, &newPost.Content, &newPost.Date, &newUser.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting comments' data")
		err = sSquared.Users.DB.QueryRow("SELECT uuid, name, picture, role_id FROM users WHERE id = ?", newUser.ID).Scan(&newUser.UUID, &newUser.Name, &newUser.Picture, &newUser.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting comments' data")

		newCom.User = user
		newPost.User = newUser
		newActivity.Comment = newCom
		newActivity.Post = newPost
		newActivity.Type = ActivityComment
		myActivity = append(myActivity, newActivity)
	}
	err = rowCom.Err()
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting comments' data")

	rowCom.Close()

	rowLikePost, err := sSquared.Likes.DB.Query("SELECT id, date, post_id FROM likepost WHERE user_id = ?", user.ID)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting likes' data")
	for rowLikePost.Next() {
		newActivity := &methods.Activity{}
		newLike := &methods.Like{}
		newPost := &methods.Post{}
		newUser := &methods.User{}
		err := rowLikePost.Scan(&newLike.ID, &newLike.Date, &newPost.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting likes' data")
		err = sSquared.Posts.DB.QueryRow("SELECT title, content, date, user_id FROM posts WHERE id = ?", newPost.ID).Scan(&newPost.Title, &newPost.Content, &newPost.Date, &newUser.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting likes' data")
		err = sSquared.Users.DB.QueryRow("SELECT uuid, name, picture, role_id FROM users WHERE id = ?", newUser.ID).Scan(&newUser.UUID, &newUser.Name, &newUser.Picture, &newUser.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting likes' data")
		newPost.User = newUser
		newActivity.Post = newPost
		newActivity.Like = newLike
		newActivity.Type = ActivityLike
		myActivity = append(myActivity, newActivity)
	}
	err = rowLikePost.Err()
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting likes' data")

	rowLikePost.Close()

	rowDislikePost, err := sSquared.Likes.DB.Query("SELECT id, date, post_id FROM dislikepost WHERE user_id = ?", user.ID)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting dislikes' data")
	for rowDislikePost.Next() {
		newActivity := &methods.Activity{}
		newDislike := &methods.Dislike{}
		newPost := &methods.Post{}
		newUser := &methods.User{}
		err := rowDislikePost.Scan(&newDislike.ID, &newDislike.Date, &newPost.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting dislikes' data")
		err = sSquared.Posts.DB.QueryRow("SELECT title, content, date, user_id FROM posts WHERE id = ?", newPost.ID).Scan(&newPost.Title, &newPost.Content, &newPost.Date, &newUser.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting dislikes' data")
		err = sSquared.Users.DB.QueryRow("SELECT uuid, name, picture, role_id FROM users WHERE id = ?", newUser.ID).Scan(&newUser.UUID, &newUser.Name, &newUser.Picture, &newUser.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting dislikes' data")
		newPost.User = newUser
		newActivity.Post = newPost
		newActivity.Dislike = newDislike
		newActivity.Type = ActivityDislike
		myActivity = append(myActivity, newActivity)
	}
	err = rowDislikePost.Err()
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting dislikes' data")

	rowDislikePost.Close()

	rowLikeCom, err := sSquared.Likes.DB.Query("SELECT id, date, comments_id FROM likecom WHERE user_id = ?", user.ID)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting likes' data")
	for rowLikeCom.Next() {
		newActivity := &methods.Activity{}
		newLike := &methods.Like{}
		newPost := &methods.Post{User: &methods.User{}}
		newCom := &methods.Comment{User: &methods.User{}}
		err := rowLikeCom.Scan(&newLike.ID, &newLike.Date, &newCom.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting likes' data")
		err = sSquared.Posts.DB.QueryRow("SELECT content, date, user_id, post_id FROM comments WHERE id = ?", newCom.ID).Scan(&newCom.Content, &newCom.Date, &newCom.User.ID, &newCom.PostID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting likes' data")
		err = sSquared.Posts.DB.QueryRow("SELECT title, content, date, user_id FROM posts WHERE id = ?", newCom.PostID).Scan(&newPost.Title, &newPost.Content, &newPost.Date, &newPost.User.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting likes' data")
		err = sSquared.Users.DB.QueryRow("SELECT uuid, name, picture, role_id FROM users WHERE id = ?", newCom.User.ID).Scan(&newCom.User.UUID, &newCom.User.Name, &newPost.User.Picture, &newCom.User.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting likes' data")
		err = sSquared.Users.DB.QueryRow("SELECT uuid, name, picture, role_id FROM users WHERE id = ?", newPost.User.ID).Scan(&newPost.User.UUID, &newPost.User.Name, &newPost.User.Picture, &newPost.User.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting likes' data")

		newActivity.Comment = newCom
		newActivity.Post = newPost
		newActivity.Like = newLike
		newActivity.Type = ActivityLike
		myActivity = append(myActivity, newActivity)
	}
	err = rowLikeCom.Err()
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting likes' data")

	rowLikeCom.Close()

	rowDislikeCom, err := sSquared.Likes.DB.Query("SELECT id, date, comments_id FROM dislikecom WHERE user_id = ?", user.ID)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting dislikes' data")
	for rowDislikeCom.Next() {
		newActivity := &methods.Activity{}
		newDislike := &methods.Dislike{}
		newPost := &methods.Post{User: &methods.User{}}
		newCom := &methods.Comment{User: &methods.User{}}
		err := rowDislikeCom.Scan(&newDislike.ID, &newDislike.Date, &newCom.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting dislikes' data")
		err = sSquared.Posts.DB.QueryRow("SELECT content, date, user_id, post_id FROM comments WHERE id = ?", newCom.ID).Scan(&newCom.Content, &newCom.Date, &newCom.User.ID, &newCom.PostID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting dislikes' data")
		err = sSquared.Posts.DB.QueryRow("SELECT title, content, date, user_id FROM posts WHERE id = ?", newCom.PostID).Scan(&newPost.Title, &newPost.Content, &newPost.Date, &newPost.User.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting dislikes' data")
		err = sSquared.Users.DB.QueryRow("SELECT uuid, name, picture, role_id FROM users WHERE id = ?", newCom.User.ID).Scan(&newCom.User.UUID, &newCom.User.Name, &newPost.User.Picture, &newCom.User.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting dislikes' data")
		err = sSquared.Users.DB.QueryRow("SELECT uuid, name, picture, role_id FROM users WHERE id = ?", newPost.User.ID).Scan(&newPost.User.UUID, &newPost.User.Name, &newPost.User.Picture, &newPost.User.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting dislikes' data")

		newActivity.Comment = newCom
		newActivity.Post = newPost
		newActivity.Dislike = newDislike
		newActivity.Type = ActivityDislike
		myActivity = append(myActivity, newActivity)
	}
	err = rowDislikeCom.Err()
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "error getting posts' data")

	rowDislikeCom.Close()

	sort.Slice(myActivity, func(i, j int) bool {
		return GetDate(myActivity[j]).Before(GetDate(myActivity[i]))
	})

	data := struct { //struct for execute and html
		LoggedIn bool
		User     *methods.User
		Activity []*methods.Activity
		NbrNotif int
	}{
		LoggedIn: loggedIn,
		User:     user,
		Activity: myActivity,
		NbrNotif: NumberOfNotif(user, sData.AllNotifs),
	}

	t, err := template.ParseFiles(`templates/myActivity.html`)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")

	t.Execute(w, data)
}
