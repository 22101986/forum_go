package methods

import (
	"database/sql"
	"fmt"
)

type NotifMethod struct {
	DB *sql.DB
}

func (db *NotifMethod) InsertInNotifs(notif *Notif, withPost bool) (int, error) {
	if notif.UserFrom.ID == notif.UserTo.ID {
		return 0, nil
	}
	var query string
	var result sql.Result
	if withPost {
		query = `INSERT INTO notifs (date, type_id, user_id_to, user_id_from, post_id) VALUES (?, ?, ?, ?, ?)`
		result, _ = db.DB.Exec(query, notif.Date, notif.Type, notif.UserTo.ID, notif.UserFrom.ID, notif.Post.ID)
	} else {
		query = `INSERT INTO notifs (date, type_id, user_id_to, user_id_from) VALUES (?, ?, ?, ?)`
		result, _ = db.DB.Exec(query, notif.Date, notif.Type, notif.UserTo.ID, notif.UserFrom.ID)
	}
	notifID, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	if notif.Comment != nil {
		_, err := db.DB.Exec("UPDATE notifs SET comments_id = ? WHERE id = ?", notif.Comment.ID, notifID)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
	}
	return int(notifID), nil
}

func (db *NotifMethod) DeleteInNotifs(id string) error {
	query := "DELETE FROM notifs WHERE id = ?"
	_, err := db.DB.Exec(query, id)
	return err
}

func (db *NotifMethod) DeleteAllNotifPost(idPost, idUser int) error {
	query := "DELETE FROM notifs WHERE user_id_to = ? AND post_id = ? AND NOT type_id in (7,8)"
	_, err := db.DB.Exec(query, idUser, idPost)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (db *NotifMethod) DeleteNotifPost(idPost, idUser int, typeNotif ...int) error {
	for _, i := range typeNotif {
		query := "DELETE FROM notifs WHERE user_id_to = ? AND post_id = ? AND type_id = ?"
		_, err := db.DB.Exec(query, idUser, idPost, i)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func (db *NotifMethod) DeleteNotifCom(idPost, idCom, idUser int, typeNotif ...int) error {
	fmt.Println(idPost, idCom, idUser, typeNotif)
	for _, i := range typeNotif {
		query := "DELETE FROM notifs WHERE post_id = ? AND comments_id = ? AND user_id_to = ? AND type_id = ?"
		_, err := db.DB.Exec(query, idPost, idCom, idUser, i)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}
