package methods

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserMethod struct {
	DB *sql.DB
}

func (db *UserMethod) InsertInUser(name, email, password string, uuid uuid.UUID) (int64, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO users (uuid, name, email, password, is_deleted, role_id, is_external)
	VALUES (?, ?, ?, ?, 0, 2, 0)`

	result, err := db.DB.Exec(stmt, uuid, name, email, passwordHash)
	if err != nil {
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (db *UserMethod) Authenticate(email, password string) (string, error) {
	var uuid, name string
	var passwordHash []byte

	stmt := `SELECT uuid, name, password FROM users WHERE email = ? AND is_deleted = 0`
	row := db.DB.QueryRow(stmt, email)

	err := row.Scan(&uuid, &name, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user not found")
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	return uuid, nil
}

func (db *UserMethod) EditProfile(id int, updatedField, query string, password, picture, role bool) (int64, error) {
	var result sql.Result
	var err error
	if password {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(updatedField), 12)
		if err != nil {
			return 0, err
		}
		result, err = db.DB.Exec(query, passwordHash, id)
		if err != nil {
			return 0, err
		}
	} else if picture {
		result, err = db.DB.Exec(query, []byte(updatedField), id)
		if err != nil {
			return 0, err
		}
	} else if role {
		updatedRole, _ := strconv.Atoi(updatedField)
		result, err = db.DB.Exec(query, updatedRole, id)
		if err != nil {
			return 0, err
		}
	} else {
		result, err = db.DB.Exec(query, updatedField, id)
		if err != nil {
			return 0, err
		}
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (db *UserMethod) UserFakeDeletion(id int) error {
	var uuid string
	if err := db.DB.QueryRow("SELECT uuid FROM users WHERE id = ?", id).Scan(&uuid); err != nil {
		return err
	}

	newName := "Deleted User" + strconv.Itoa(id)
	newEmail := "[Email_" + uuid + "]"

	updateQuery := `UPDATE users SET name = ?, email = ?, password = 'Bloub bloub', picture = ?, role_id = 1, is_deleted = 1 WHERE id = ?`
	result, err := db.DB.Exec(updateQuery, newName, newEmail, nil, id)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}

func (u *UserMethod) FindOrCreateUser(userUUID *string, name, email, nbr string) (int64, error) {
	var userID int64
	row := u.DB.QueryRow("SELECT id, uuid FROM users WHERE email = ? AND is_external ="+nbr, email)
	err := row.Scan(&userID, userUUID)
	if err == sql.ErrNoRows {
		newUUID, err := uuid.NewV4()
		if err != nil {
			return 0, err
		}
		*userUUID = newUUID.String()
		password := ""
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			return 0, err
		}

		res, err := u.DB.Exec("INSERT INTO users (uuid, name, email, password, is_deleted, role_id, is_external) VALUES (? ,? , ?, ?, 0, 2, "+nbr+")", *userUUID, name, email, passwordHash)
		if err != nil {
			return 0, err
		}
		userID, err = res.LastInsertId()
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}
	return userID, nil
}
