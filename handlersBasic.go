package main

import (
	"fmt"
	"forum/methods"
	"html/template"
	"net/http"
	"strconv"
)

// function that handles the requests related to the errors (404, 500)
func ErrorHandler(w http.ResponseWriter, r *http.Request, errorCode int, errorMessage string) {
	w.WriteHeader(errorCode) //display error message in the terminal of the navigator

	t, err := template.New(`error.html`).ParseFiles(`templates/error.html`) // parse through the files to find the error file
	if err != nil {
		fmt.Println(err)
		return
	}
	errData := struct { // structure for execute
		ErrorCode    int
		ErrorMessage string
	}{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(w, errData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// function to handle the index requests
func IndexHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	if r.URL.Path != "/" { // handle error 404
		ErrorHandler(w, r, http.StatusNotFound, "Page not found")
		return
	}

	loggedIn := LoggedInVerif(r) // verify if the cookie is setup with a session token
	DuplicateLog(loggedIn, w, r) // verify if the cookie is unique (handle double connection)

	user := &methods.User{} // get ID and Username of connected user
	if loggedIn {
		cookie, err := r.Cookie("session_token")
		ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

		userUUID := sessions[cookie.Value]

		err = sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id, is_external FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Picture, &user.Role, &user.IsExternal)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")
	}

	selectedCategories := r.URL.Query()["categories"] // getting categories for filter
	selectedCategoriesInt := []int{}
	for _, nb := range selectedCategories {
		i, _ := strconv.Atoi(nb)
		selectedCategoriesInt = append(selectedCategoriesInt, i)
	}

	var posts []*methods.Post
	var err error

	if len(selectedCategories) > 0 { // getting all posts with filter applied
		posts, err = sSquared.Posts.GetPostsByCategories(selectedCategories)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "getting posts error")
	} else {
		posts = sData.AllPosts
	}

	allPostsInv := []*methods.Post{} // invert post order
	for i := len(posts) - 1; i >= 0; i-- {
		allPostsInv = append(allPostsInv, posts[i])
	}

	for _, post := range allPostsInv { // add bool for color on likes/dislikes
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
	data := struct { // struct for execute and html
		LoggedIn    bool
		User        *methods.User
		IndexUrl    string
		AllPosts    []*methods.Post
		AllCats     []*methods.Categories
		AllSelected []int
		NbrNotif    int
	}{
		LoggedIn:    loggedIn,
		User:        user,
		IndexUrl:    r.URL.String(),
		AllPosts:    allPostsInv,
		AllCats:     sData.AllCats,
		AllSelected: selectedCategoriesInt,
		NbrNotif:    NumberOfNotif(user, sData.AllNotifs),
	}

	t, err := template.ParseFiles(`templates/index.html`)
	ErrDiffNil(err, w, r, http.StatusNotFound, "error parsing file index.html not found")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(w, data)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")
}

func ImagePostHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	picture := []byte{}
	postID, err := strconv.Atoi(r.URL.Query().Get("ID"))
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "Error getting ID")
	for _, post := range sData.AllPosts {
		if post.ID == postID {
			picture = post.Blob
			break
		}
	}
	contentType := http.DetectContentType(picture)
	w.Header().Set("Content-Type", contentType)
	w.Write(picture)
}
func ProfilePictureHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	picture := []byte{}
	userID, err := strconv.Atoi(r.URL.Query().Get("ID"))
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "Error getting ID")
	for _, user := range sData.AllUsers {
		if user.ID == userID {
			picture = user.Picture
			break
		}
	}
	contentType := http.DetectContentType(picture)
	w.Header().Set("Content-Type", contentType)
	w.Write(picture)
}
