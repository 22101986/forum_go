package methods

import (
	"database/sql"
)

type BlobMethod struct {
	DB *sql.DB
}

func (b *BlobMethod) InsertInBlob(image []byte, id int) error {
	query := "INSERT INTO blob (picture, post_id) VALUES (?, ?)"
	_, err := b.DB.Exec(query, image, id)
	if err != nil {
		return err
	}
	return nil
}
