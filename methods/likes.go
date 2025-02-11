package methods

import (
	"database/sql"
	"time"
)

type LikeMethod struct {
	DB *sql.DB
}

func (l *LikeMethod) InsertInLikePost(user, post int) (int64, error) {
	query := `INSERT INTO likepost (user_id, date, post_id) VALUES (?, ?, ?)`
	result, err := l.DB.Exec(query, user, time.Now().Format("2006-01-02 15:04:05"), post)
	if err != nil {
		return 0, err
	}
	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return postID, nil
}
func (l *LikeMethod) DeleteInLikePost(id int) error {
	query := "DELETE FROM likepost WHERE id = ?"
	_, err := l.DB.Exec(query, id)
	return err
}
func (l *LikeMethod) InsertInDislikePost(user, post int) (int64, error) {
	query := `INSERT INTO dislikepost (user_id, date, post_id) VALUES (?, ?, ?)`
	result, err := l.DB.Exec(query, user, time.Now().Format("2006-01-02 15:04:05"), post)
	if err != nil {
		return 0, err
	}
	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return postID, nil
}
func (l *LikeMethod) DeleteInDislikePost(id int) error {
	query := "DELETE FROM dislikepost WHERE id = ?"
	_, err := l.DB.Exec(query, id)
	return err
}

func (l *LikeMethod) InsertInLikeCom(user, com int) (int64, error) {
	query := `INSERT INTO likecom (user_id, date, comments_id) VALUES (?, ?, ?)`
	result, err := l.DB.Exec(query, user, time.Now().Format("2006-01-02 15:04:05"), com)
	if err != nil {
		return 0, err
	}
	comID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return comID, nil
}
func (l *LikeMethod) DeleteInLikeCom(id int) error {
	query := "DELETE FROM likecom WHERE id = ?"
	_, err := l.DB.Exec(query, id)
	return err
}
func (l *LikeMethod) InsertInDislikeCom(user, com int) (int64, error) {
	query := `INSERT INTO dislikecom (user_id, date, comments_id) VALUES (?, ?, ?)`
	result, err := l.DB.Exec(query, user, time.Now().Format("2006-01-02 15:04:05"), com)
	if err != nil {
		return 0, err
	}
	comID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return comID, nil
}
func (l *LikeMethod) DeleteInDislikeCom(id int) error {
	query := "DELETE FROM dislikecom WHERE id = ?"
	_, err := l.DB.Exec(query, id)
	return err
}
