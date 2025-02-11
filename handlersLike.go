package main

import (
	"forum/methods"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func LikePostHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	postID := r.URL.Query().Get("ID") // check what post is liked
	if postID == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "no ID in URL")
		return
	}
	postIDint, err := strconv.Atoi(postID)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")
	if postIDint < 1 {
		ErrorHandler(w, r, http.StatusNotFound, "Page not found")
		return
	}

	var spePost *methods.Post //get post in the data struct
	for _, post := range sData.AllPosts {
		if post.ID == postIDint {
			spePost = post
			break
		}
	}

	cookie, err := r.Cookie("session_token") // getting connected user
	ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

	userUUID := sessions[cookie.Value]
	var user methods.User
	err = sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Picture, &user.Role)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")

	var likeID, dislikeID int
	speLike := methods.Like{}
	newNotif := &methods.Notif{
		ID:       0,
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Type:     1,
		UserTo:   spePost.User,
		UserFrom: &user,
		Post:     spePost,
	}

	// Check if the user has already disliked the post
	err = sSquared.Likes.DB.QueryRow("SELECT id FROM dislikepost WHERE user_id = ? AND post_id = ?", user.ID, postIDint).Scan(&dislikeID)
	if err == nil {
		sSquared.Likes.DeleteInDislikePost(dislikeID)                                  // Remove dislike if is exists (in the database)
		sSquared.Notifs.DeleteNotifPost(spePost.ID, spePost.User.ID, NotifDislikePost) // Delete notification if exist (type 2 -> dislike post)
		var tempTab []*methods.Dislike
		var tempNot []*methods.Notif
		for _, dislike := range spePost.Dislikes { // In the data struct
			if dislike.ID != dislikeID {
				tempTab = append(tempTab, dislike)
			}
		}
		spePost.Dislikes = tempTab
		for _, notif := range sData.AllNotifs {
			if notif.UserTo.ID != spePost.User.ID && notif.Post.ID != spePost.ID && notif.Type != NotifDislikePost {
				tempNot = append(tempNot, notif)
			}
		}
		sData.AllNotifs = tempNot
	}

	// Check if the user has already liked the post
	err = sSquared.Likes.DB.QueryRow("SELECT id FROM likepost WHERE user_id = ? AND post_id = ?", user.ID, postIDint).Scan(&likeID)
	if err == nil {
		sSquared.Likes.DeleteInLikePost(likeID) // Remove like if it exists (in the database)
		sSquared.Notifs.DeleteNotifPost(spePost.ID, spePost.User.ID, NotifLikePost)
		var tempTab []*methods.Like
		var tempNot []*methods.Notif
		for _, like := range spePost.Likes { // In the data struct
			if like.ID != likeID {
				tempTab = append(tempTab, like)
			}
		}
		spePost.Likes = tempTab
		for _, notif := range sData.AllNotifs {
			if notif.UserTo.ID != spePost.User.ID && notif.Post.ID != spePost.ID && notif.Type != NotifDislikePost {
				tempNot = append(tempNot, notif)
			}
		}
		sData.AllNotifs = tempNot
	} else {
		newLikeID, _ := sSquared.Likes.InsertInLikePost(user.ID, postIDint) // Add like if not already liked
		newNotif.ID, _ = sSquared.Notifs.InsertInNotifs(newNotif, true)
		speLike.ID, speLike.User = int(newLikeID), &user
		spePost.Likes = append(spePost.Likes, &speLike)
		if newNotif.ID != 0 {
			sData.AllNotifs = append(sData.AllNotifs, newNotif)
		}
	}

	url := r.URL.Query().Get("source") // get the url where like happened
	if url == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "no URL")
		return
	}
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func DislikePostHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	postID := r.URL.Query().Get("ID") // check what post is disliked
	if postID == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "no ID in URL")
		return
	}
	postIDint, err := strconv.Atoi(postID)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")
	if postIDint < 1 {
		ErrorHandler(w, r, http.StatusNotFound, "Page not found")
		return
	}

	var spePost *methods.Post //get post in the data struct
	for _, post := range sData.AllPosts {
		if post.ID == postIDint {
			spePost = post
			break
		}
	}

	cookie, err := r.Cookie("session_token") // getting connected user
	ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

	userUUID := sessions[cookie.Value]
	var user methods.User
	err = sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Picture, &user.Role)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")

	var likeID, dislikeID int
	speDislike := methods.Dislike{}
	newNotif := &methods.Notif{
		ID:       0,
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Type:     NotifDislikePost,
		UserTo:   spePost.User,
		UserFrom: &user,
		Post:     spePost,
	}

	// Check if the user has already liked the post
	err = sSquared.Likes.DB.QueryRow("SELECT id FROM likepost WHERE user_id = ? AND post_id = ?", user.ID, postIDint).Scan(&likeID)
	if err == nil {
		sSquared.Likes.DeleteInLikePost(likeID) // Remove like if it exists (in the database)
		sSquared.Notifs.DeleteNotifPost(spePost.ID, spePost.User.ID, NotifLikePost)
		var tempTab []*methods.Like
		var tempNot []*methods.Notif
		for _, like := range spePost.Likes { // In the data struct
			if like.ID != likeID {
				tempTab = append(tempTab, like)
			}
		}
		spePost.Likes = tempTab
		for _, notif := range sData.AllNotifs {
			if notif.UserTo.ID != spePost.User.ID && notif.Post.ID != spePost.ID && notif.Type != NotifLikePost {
				tempNot = append(tempNot, notif)
			}
		}
		sData.AllNotifs = tempNot

	}

	// Check if the user has already disliked the post
	err = sSquared.Likes.DB.QueryRow("SELECT id FROM dislikepost WHERE user_id = ? AND post_id = ?", user.ID, postIDint).Scan(&dislikeID)
	if err == nil {
		sSquared.Likes.DeleteInDislikePost(dislikeID) // Remove dislike if it exists (in the database)
		sSquared.Notifs.DeleteNotifPost(spePost.ID, spePost.User.ID, NotifDislikePost)
		var tempTab []*methods.Dislike
		var tempNot []*methods.Notif
		for _, dislike := range spePost.Dislikes { // In the data struct
			if dislike.ID != dislikeID {
				tempTab = append(tempTab, dislike)
			}
		}
		spePost.Dislikes = tempTab
		for _, notif := range sData.AllNotifs {
			if notif.UserTo.ID != spePost.User.ID && notif.Post.ID != spePost.ID && notif.Type != NotifDislikePost {
				tempNot = append(tempNot, notif)
			}
		}
		sData.AllNotifs = tempNot

	} else {
		newDislikeID, _ := sSquared.Likes.InsertInDislikePost(user.ID, postIDint) // Add dislike if not already disliked
		newNotif.ID, _ = sSquared.Notifs.InsertInNotifs(newNotif, true)
		speDislike.ID, speDislike.User = int(newDislikeID), &user
		spePost.Dislikes = append(spePost.Dislikes, &speDislike)
		if newNotif.ID != 0 {
			sData.AllNotifs = append(sData.AllNotifs, newNotif)
		}

	}

	url := r.URL.Query().Get("source") // get the url where like happened
	if url == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "no URL")
		return
	}
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func LikeComHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	comID := r.URL.Query().Get("ID") // check what com is disliked
	if comID == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "no ID in URL")
		return
	}
	comIDint, err := strconv.Atoi(comID)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")
	if comIDint < 1 {
		ErrorHandler(w, r, http.StatusNotFound, "Page not found")
		return
	}

	url := r.URL.Query().Get("source") // get the url where like happened
	if url == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "no URL")
		return
	}
	_, postID, _ := strings.Cut(url, "=") // check what post is associated to the com
	postIDint, err := strconv.Atoi(postID)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")
	var spePost *methods.Post //get post in the data struct
	for _, post := range sData.AllPosts {
		if post.ID == postIDint {
			spePost = post
			break
		}
	}
	var speCom *methods.Comment //get com in the data struct
	for _, com := range spePost.Comments {
		if com.ID == comIDint {
			speCom = com
			break
		}
	}

	cookie, err := r.Cookie("session_token") // getting connected user
	ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

	userUUID := sessions[cookie.Value]
	var user methods.User
	err = sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Picture, &user.Role)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")

	var likeID, dislikeID int
	speLike := methods.Like{}

	newNotif := &methods.Notif{
		ID:       0,
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Type:     NotifLikeCom,
		UserTo:   speCom.User,
		UserFrom: &user,
		Post:     spePost,
		Comment:  speCom,
	}

	// Check if the user has already disliked the com
	err = sSquared.Likes.DB.QueryRow("SELECT id FROM dislikecom WHERE user_id = ? AND comments_id = ?", user.ID, comIDint).Scan(&dislikeID)
	if err == nil {
		sSquared.Likes.DeleteInDislikeCom(dislikeID) // Remove dislike if is exists (in the database)
		sSquared.Notifs.DeleteNotifCom(spePost.ID, speCom.ID, speCom.User.ID, NotifDislikeCom)
		var tempTab []*methods.Dislike
		var tempNot []*methods.Notif
		for _, dislike := range speCom.Dislikes { // In the data struct
			if dislike.ID != dislikeID {
				tempTab = append(tempTab, dislike)
			}
		}
		speCom.Dislikes = tempTab
		for _, notif := range sData.AllNotifs {
			if notif.UserTo.ID != spePost.User.ID && notif.Post.ID != spePost.ID && notif.Type != NotifDislikeCom {
				tempNot = append(tempNot, notif)
			}
		}
		sData.AllNotifs = tempNot
	}

	// Check if the user has already liked the com
	err = sSquared.Likes.DB.QueryRow("SELECT id FROM likecom WHERE user_id = ? AND comments_id = ?", user.ID, comIDint).Scan(&likeID)
	if err == nil {
		sSquared.Likes.DeleteInLikeCom(likeID) // Remove like if it exists	(in the database)
		sSquared.Notifs.DeleteNotifCom(spePost.ID, speCom.ID, speCom.User.ID, NotifLikeCom)
		var tempTab []*methods.Like
		var tempNot []*methods.Notif
		for _, like := range speCom.Likes { // In the data struct
			if like.ID != likeID {
				tempTab = append(tempTab, like)
			}
		}
		speCom.Likes = tempTab
		for _, notif := range sData.AllNotifs {
			if notif.UserTo.ID != spePost.User.ID && notif.Post.ID != spePost.ID && notif.Type != NotifDislikeCom {
				tempNot = append(tempNot, notif)
			}
		}
		sData.AllNotifs = tempNot
	} else {
		newLikeID, _ := sSquared.Likes.InsertInLikeCom(user.ID, comIDint) // Add like if not already liked
		speLike.ID, speLike.User = int(newLikeID), &user
		speCom.Likes = append(speCom.Likes, &speLike)

		newNotif.ID, _ = sSquared.Notifs.InsertInNotifs(newNotif, true)
		if newNotif.ID != 0 {
			sData.AllNotifs = append(sData.AllNotifs, newNotif)
		}
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func DislikeComHandler(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	comID := r.URL.Query().Get("ID") // check what com is disliked
	if comID == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "no ID in URL")
		return
	}
	comIDint, err := strconv.Atoi(comID)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")
	if comIDint < 1 {
		ErrorHandler(w, r, http.StatusNotFound, "Page not found")
		return
	}

	url := r.URL.Query().Get("source") // get the url where like happened
	if url == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "no URL")
		return
	}
	_, postID, _ := strings.Cut(url, "=") // check what post is associated to the com
	postIDint, err := strconv.Atoi(postID)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "internal server error")
	var spePost *methods.Post //get post in the data struct
	for _, post := range sData.AllPosts {
		if post.ID == postIDint {
			spePost = post
			break
		}
	}
	var speCom *methods.Comment //get com in the data struct
	for _, com := range spePost.Comments {
		if com.ID == comIDint {
			speCom = com
			break
		}
	}

	cookie, err := r.Cookie("session_token") // getting connected user
	ErrDiffNil(err, w, r, http.StatusBadRequest, "gathering cookie error")

	userUUID := sessions[cookie.Value]
	var user methods.User
	err = sSquared.Users.DB.QueryRow("SELECT id, name, picture, role_id FROM users WHERE uuid = ?", userUUID).Scan(&user.ID, &user.Name, &user.Picture, &user.Role)
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "user not found")

	var likeID, dislikeID int
	speDislike := methods.Dislike{}

	newNotif := &methods.Notif{
		ID:       0,
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Type:     NotifDislikeCom,
		UserTo:   speCom.User,
		UserFrom: &user,
		Post:     spePost,
		Comment:  speCom,
	}

	// Check if the user has already liked the com
	err = sSquared.Likes.DB.QueryRow("SELECT id FROM likecom WHERE user_id = ? AND comments_id = ?", user.ID, comIDint).Scan(&likeID)
	if err == nil {
		sSquared.Likes.DeleteInLikeCom(likeID) // Remove like if it exists	(in the database)
		sSquared.Notifs.DeleteNotifCom(spePost.ID, speCom.ID, speCom.User.ID, NotifLikeCom)
		var tempTab []*methods.Like
		var tempNot []*methods.Notif
		for _, like := range speCom.Likes { // In the data struct
			if like.ID != likeID {
				tempTab = append(tempTab, like)
			}
		}
		speCom.Likes = tempTab
		for _, notif := range sData.AllNotifs {
			if notif.UserTo.ID != spePost.User.ID && notif.Post.ID != spePost.ID && notif.Type != NotifLikeCom {
				tempNot = append(tempNot, notif)
			}
		}
		sData.AllNotifs = tempNot
	}

	// Check if the user has already disliked the com
	err = sSquared.Likes.DB.QueryRow("SELECT id FROM dislikecom WHERE user_id = ? AND comments_id = ?", user.ID, comIDint).Scan(&dislikeID)
	if err == nil {
		sSquared.Likes.DeleteInDislikeCom(dislikeID) // Remove dislike if it exists	(in the database)
		sSquared.Notifs.DeleteNotifCom(spePost.ID, speCom.ID, speCom.User.ID, NotifDislikeCom)
		var tempTab []*methods.Dislike
		var tempNot []*methods.Notif
		for _, dislike := range speCom.Dislikes { // In the data struct
			if dislike.ID != dislikeID {
				tempTab = append(tempTab, dislike)
			}
		}
		speCom.Dislikes = tempTab
		for _, notif := range sData.AllNotifs {
			if notif.UserTo.ID != spePost.User.ID && notif.Post.ID != spePost.ID && notif.Type != NotifDislikeCom {
				tempNot = append(tempNot, notif)
			}
		}
		sData.AllNotifs = tempNot
	} else {
		newDislikeID, _ := sSquared.Likes.InsertInDislikeCom(user.ID, comIDint) // Add dislike if not already disliked
		speDislike.ID, speDislike.User = int(newDislikeID), &user
		speCom.Dislikes = append(speCom.Dislikes, &speDislike)

		newNotif.ID, _ = sSquared.Notifs.InsertInNotifs(newNotif, true)
		if newNotif.ID != 0 {
			sData.AllNotifs = append(sData.AllNotifs, newNotif)
		}
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}
