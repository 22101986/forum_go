package methods

import (
	"database/sql"
)

type CommentMethod struct {
	DB *sql.DB
}

func (p *CommentMethod) InsertInComments(comment *Comment) (int64, error) {
	query := `INSERT INTO comments (content, Date, user_id, post_id ) VALUES (?, ?, ?, ?)`
	result, err := p.DB.Exec(query, comment.Content, comment.Date, comment.User.ID, comment.PostID)
	if err != nil {
		return 0, err
	}
	commentID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return commentID, nil
}

func (p *CommentMethod) DeleteInComments(id string) error {
	_, err := p.DB.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		return err
	}
	query := `DELETE FROM comments WHERE id = ?`
	_, err = p.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return err
}
