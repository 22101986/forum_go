package main

import (
	"database/sql"
	"forum/methods"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

const (
	RoleDrowned         = 1
	RoleUser            = 2
	RoleModo            = 3
	RoleAdmin           = 4
	NotifLikePost       = 1
	NotifDislikePost    = 2
	NotifLikeCom        = 3
	NotifDislikeCom     = 4
	NotifComment        = 5
	NotifAskModo        = 6
	NotifReportPost     = 7
	NotifReportCom      = 8
	NotifAdminAccept    = 9
	NotifAdminRefuse    = 10
	NotifAcceptDeletion = 11
	NotifRefuseDeletion = 12
	BasicUser           = 0
	GoogleUser          = 1
	DiscordUser         = 2
	GithubUser          = 3
	ActivityPost        = 1
	ActivityComment     = 2
	ActivityLike        = 3
	ActivityDislike     = 4
)

// Create the struct for methods
func CreateStructSquare(db *sql.DB) StructSquare {
	sSquared := StructSquare{
		Users: &methods.UserMethod{
			DB: db,
		},
		Posts: &methods.PostMethod{
			DB: db,
		},
		Comments: &methods.CommentMethod{
			DB: db,
		},
		Likes: &methods.LikeMethod{
			DB: db,
		},
		Notifs: &methods.NotifMethod{
			DB: db,
		},
		Blobs: &methods.BlobMethod{
			DB: db,
		},
	}
	return sSquared
}

// Create the struct containing all the data
func CreateStructData() StructData {
	sData := StructData{}
	return sData
}

func GetDate(activity *methods.Activity) time.Time {
	const layout = "2006-01-02 15:04:05"

	var dateStr string
	switch activity.Type {
	case ActivityPost:
		if activity.Post == nil {
			return time.Time{}
		}
		dateStr = activity.Post.Date
	case ActivityComment:
		if activity.Comment == nil {
			return time.Time{}
		}
		dateStr = activity.Comment.Date
	case ActivityLike:
		if activity.Like == nil {
			return time.Time{}
		}
		dateStr = activity.Like.Date
	case ActivityDislike:
		if activity.Dislike == nil {
			return time.Time{}
		}
		dateStr = activity.Dislike.Date
	}

	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}
	}
	return date
}

func GetImg(w http.ResponseWriter, file multipart.File, fileHeader *multipart.FileHeader) []byte {
	imageData := make([]byte, fileHeader.Size)
	_, err := file.Read(imageData)
	if err != nil && err != io.EOF {
		http.Error(w, "Erreur lors de la lecture du fichier", http.StatusInternalServerError)
	}
	return imageData
}

func NumberOfNotif(user *methods.User, allNotifs []*methods.Notif) int {
	tabNotifs := []*methods.Notif{}
	for _, notif := range allNotifs {
		if notif.UserTo.ID == user.ID {
			tabNotifs = append(tabNotifs, notif)
		}
	}
	return len(tabNotifs)
}

func VerifyContent(content string) bool {
	if content == "" {
		return true
	}
	if content[0] == ' ' || content[0] == '\n' || content[0] == '\r' {
		return VerifyContent(content[1:])
	}
	return false
}
