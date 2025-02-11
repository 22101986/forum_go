package main

import (
	"database/sql"
	"fmt"
	"forum/methods"
	"net/http"
)

// function for clear all structures
func ClearStructs(sData *StructData) {
	sData.AllCats = nil
	sData.AllUsers = nil

	for _, post := range sData.AllPosts {
		for _, comment := range post.Comments {
			comment.Likes = nil
			comment.Dislikes = nil
		}
		post.Comments = nil
		post.Cats = nil
		post.Likes = nil
		post.Dislikes = nil
	}
	sData.AllPosts = nil

	for _, notif := range sData.AllNotifs {
		notif.UserTo = nil
		notif.UserFrom = nil
		notif.Post = nil
		notif.Comment = nil
	}
	sData.AllNotifs = nil

	fmt.Println("All data structures cleared.")
}

// Function to refill all the data structure with data from database
func FillingStruct(w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	//for categories
	rowCat, err := sSquared.Posts.DB.Query("SELECT * FROM categories")
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in categories")

	for rowCat.Next() {
		newCat := &methods.Categories{}
		err := rowCat.Scan(&newCat.ID, &newCat.Name)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in categories")

		sData.AllCats = append(sData.AllCats, newCat)
	}
	err = rowCat.Err()
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in categories")
	rowCat.Close()

	// for users
	rowUser, err := sSquared.Users.DB.Query("SELECT id, uuid, name, picture, role_id, is_external FROM users")
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in users")

	for rowUser.Next() {
		newUser := &methods.User{}
		err = rowUser.Scan(&newUser.ID, &newUser.UUID, &newUser.Name, &newUser.Picture, &newUser.Role, &newUser.IsExternal)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in notifs")

		sData.AllUsers = append(sData.AllUsers, newUser)
	}
	rowUser.Close()

	// for posts
	rowPost, err := sSquared.Posts.DB.Query("SELECT * FROM posts")
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")

	for rowPost.Next() {
		newPost := &methods.Post{}
		newUser := &methods.User{}
		err = rowPost.Scan(&newPost.ID, &newPost.Title, &newPost.Content, &newPost.Date, &newUser.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")
		err = sSquared.Users.DB.QueryRow("SELECT id, uuid, name, picture, role_id FROM users WHERE id = ?", newUser.ID).Scan(&newUser.ID, &newUser.UUID, &newUser.Name, &newUser.Picture, &newUser.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")

		newPost.User = newUser
		// for comms in each post
		rowCom, err := sSquared.Comments.DB.Query("SELECT * FROM comments WHERE post_id = ?", newPost.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")

		for rowCom.Next() {
			newComment := &methods.Comment{}
			newUser2 := &methods.User{}
			err = rowCom.Scan(&newComment.ID, &newComment.Content, &newComment.Date, &newUser2.ID, &newComment.PostID)
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")
			err = sSquared.Users.DB.QueryRow("SELECT id, uuid, name, picture, role_id FROM users WHERE id = ?", newUser2.ID).Scan(&newUser2.ID, &newUser2.UUID, &newUser2.Name, &newUser2.Picture, &newUser2.Role)
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")

			newComment.User = newUser2
			//for like in each comm
			rowLike, err := sSquared.Likes.DB.Query("SELECT id, user_id, date FROM likecom WHERE comments_id = ?", newComment.ID)
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")

			for rowLike.Next() {
				newLike := &methods.Like{}
				newUser := &methods.User{}
				err := rowLike.Scan(&newLike.ID, &newUser.ID, &newLike.Date)
				ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")
				err = sSquared.Users.DB.QueryRow("SELECT id, uuid, name, picture, role_id FROM users WHERE id = ?", newUser.ID).Scan(&newUser.ID, &newUser.UUID, &newUser.Name, &newUser.Picture, &newUser.Role)
				ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")

				newLike.User = newUser
				newComment.Likes = append(newComment.Likes, newLike)
			}
			rowLike.Close()
			//for dislike in each comm
			rowDislike, err := sSquared.Likes.DB.Query("SELECT id, user_id, date FROM dislikecom WHERE comments_id = ?", newComment.ID)
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")

			for rowDislike.Next() {
				newDislike := &methods.Dislike{}
				newUser := &methods.User{}
				err = rowDislike.Scan(&newDislike.ID, &newUser.ID, &newDislike.Date)
				ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")
				err = sSquared.Users.DB.QueryRow("SELECT id, uuid, name, picture, role_id FROM users WHERE id = ?", newUser.ID).Scan(&newUser.ID, &newUser.UUID, &newUser.Name, &newUser.Picture, &newUser.Role)
				ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")

				newDislike.User = newUser
				newComment.Dislikes = append(newComment.Dislikes, newDislike)
			}
			rowDislike.Close()

			newPost.Comments = append(newPost.Comments, newComment)
		}
		rowCom.Close()
		//for categories for each post
		rowRel, err := sSquared.Posts.DB.Query("SELECT cat_id FROM catpostrel WHERE post_id = ?", newPost.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")

		for rowRel.Next() {
			newCat := &methods.Categories{}
			err = rowRel.Scan(&newCat.ID)
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")
			err = sSquared.Posts.DB.QueryRow("SELECT name FROM categories WHERE id = ?", newCat.ID).Scan(&newCat.Name)
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")

			newPost.Cats = append(newPost.Cats, newCat)
		}
		rowRel.Close()
		//for the len of categories -1 (for html) for each post
		newPost.LenCat = len(newPost.Cats) - 1
		//for like in each post
		rowLike, err := sSquared.Likes.DB.Query("SELECT id, user_id, date FROM likepost WHERE post_id = ?", newPost.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")

		for rowLike.Next() {
			newLike := &methods.Like{}
			newUser := &methods.User{}
			err = rowLike.Scan(&newLike.ID, &newUser.ID, &newLike.Date)
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")
			err = sSquared.Users.DB.QueryRow("SELECT id, uuid, name, picture, role_id FROM users WHERE id = ?", newUser.ID).Scan(&newUser.ID, &newUser.UUID, &newUser.Name, &newUser.Picture, &newUser.Role)
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")

			newLike.User = newUser
			newPost.Likes = append(newPost.Likes, newLike)
		}
		rowLike.Close()
		//for dislike in each post
		rowDislike, err := sSquared.Likes.DB.Query("SELECT id, user_id, date FROM dislikepost WHERE post_id = ?", newPost.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")

		for rowDislike.Next() {
			newDislike := &methods.Dislike{}
			newUser := &methods.User{}
			err = rowDislike.Scan(&newDislike.ID, &newUser.ID, &newDislike.Date)
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")
			err = sSquared.Users.DB.QueryRow("SELECT id, uuid, name, picture, role_id FROM users WHERE id = ?", newUser.ID).Scan(&newUser.ID, &newUser.UUID, &newUser.Name, &newUser.Picture, &newUser.Role)
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")

			newDislike.User = newUser
			newPost.Dislikes = append(newPost.Dislikes, newDislike)
		}
		rowDislike.Close()

		if err = sSquared.Blobs.DB.QueryRow("SELECT picture FROM blob WHERE post_id = ?", newPost.ID).Scan(&newPost.Blob); err != nil {
			if err != sql.ErrNoRows {
				ErrorHandler(w, r, http.StatusInternalServerError, "database error in posts")
			} else {
				newPost.Blob = nil
			}
		}
		sData.AllPosts = append(sData.AllPosts, newPost)
	}
	err = rowPost.Err()
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in posts")

	rowPost.Close()
	// for notifs
	rowNotif, err := sSquared.Notifs.DB.Query("SELECT id, date, type_id, user_id_to, user_id_from FROM notifs")
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in notifs")

	for rowNotif.Next() {
		newNotif := &methods.Notif{}
		newUserTo := &methods.User{}
		newUserFrom := &methods.User{}
		newPost := &methods.Post{}
		newComment := &methods.Comment{}
		err = rowNotif.Scan(&newNotif.ID, &newNotif.Date, &newNotif.Type, &newUserTo.ID, &newUserFrom.ID)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in notifs")
		err = sSquared.Users.DB.QueryRow("SELECT id, uuid, name, picture, role_id FROM users WHERE id = ?", newUserTo.ID).Scan(&newUserTo.ID, &newUserTo.UUID, &newUserTo.Name, &newUserTo.Picture, &newUserTo.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in notifs")
		err = sSquared.Users.DB.QueryRow("SELECT id, uuid, name, picture, role_id FROM users WHERE id = ?", newUserFrom.ID).Scan(&newUserFrom.ID, &newUserFrom.UUID, &newUserFrom.Name, &newUserFrom.Picture, &newUserFrom.Role)
		ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in notifs")

		switch newNotif.Type {
		case NotifAskModo, NotifAdminAccept, NotifAdminRefuse, NotifAcceptDeletion, NotifRefuseDeletion:
		default:
			err = sSquared.Notifs.DB.QueryRow("SELECT post_id FROM notifs WHERE id = ?", newNotif.ID).Scan(&newPost.ID)
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in notifs")
			if err = sSquared.Notifs.DB.QueryRow("SELECT comments_id FROM notifs WHERE id = ?", newNotif.ID).Scan(&newComment.ID); err != nil {
				newComment.ID = 0
			}
			if newComment.ID != 0 {
				if newNotif.Type == NotifComment {
					newComment.User = newUserFrom
				} else {
					newComment.User = newUserTo
				}

				err = sSquared.Comments.DB.QueryRow("SELECT content, date, post_id FROM comments WHERE id = ?", newComment.ID).Scan(&newComment.Content, &newComment.Date, &newComment.PostID)
				ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in notifs")
				newPost.ID = newComment.PostID
				newPost.Comments = append(newPost.Comments, newComment)
			}
			err = sSquared.Posts.DB.QueryRow("SELECT title, content, date FROM posts WHERE id = ?", newPost.ID).Scan(&newPost.Title, &newPost.Content, &newPost.Date)
			ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in notifs")
		}

		newNotif.UserTo = newUserTo
		newNotif.UserFrom = newUserFrom
		newNotif.Post = newPost
		newNotif.Comment = newComment
		sData.AllNotifs = append(sData.AllNotifs, newNotif)
	}
	err = rowNotif.Err()
	ErrDiffNil(err, w, r, http.StatusInternalServerError, "database error in notifs")
	rowNotif.Close()

	fmt.Println("All data correctly filled in data structure")
}
