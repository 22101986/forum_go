package methods

import (
	"database/sql"
	"strconv"
	"strings"
)

type PostMethod struct {
	DB *sql.DB
}

func (db *PostMethod) InsertInPosts(post *Post) (int64, error) {
	query := `INSERT INTO posts (title, content, date, user_id ) VALUES (?, ?, ?, ?)`
	result, err := db.DB.Exec(query, post.Title, post.Content, post.Date, post.User.ID)
	if err != nil {
		return 0, err
	}
	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return postID, nil
}
func (p *PostMethod) DeleteInPosts(id string) error {
	_, err := p.DB.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		return err
	}
	query := `DELETE FROM posts WHERE id = ?`
	_, err = p.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return err
}
func (db *PostMethod) InsertInRel(keyCat int, postID int64) error {
	query := `INSERT INTO catpostrel (cat_id, post_id) VALUES (?, ?)`
	_, err := db.DB.Exec(query, keyCat, int(postID))
	return err
}

func (db *PostMethod) GetPostsByCategories(catIDs []string) ([]*Post, error) {
	query := `
		SELECT p.id, p.title, p.content, p.date, u.id AS user_id, u.name, u.picture AS user_name
		FROM posts AS p
		JOIN users AS u ON p.user_id = u.id
		JOIN catpostrel AS cpr ON p.id = cpr.post_id
		WHERE cpr.cat_id IN (` + strings.Join(catIDs, ",") + `)
		GROUP BY p.id
	`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	postMap := make(map[int]*Post)
	var postIDs []int

	for rows.Next() {
		post := &Post{User: &User{}}
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Date, &post.User.ID, &post.User.Name, &post.User.Picture)
		if err != nil {
			return nil, err
		}
		postMap[post.ID] = post
		postIDs = append(postIDs, post.ID)
	}

	if len(postIDs) == 0 {
		return []*Post{}, nil
	}

	catQuery := `
		SELECT cpr.post_id, c.id, c.name
		FROM catpostrel AS cpr
		JOIN categories AS c ON cpr.cat_id = c.id
		WHERE cpr.post_id IN (` + strings.Join(IntSliceToStringSlice(postIDs), ",") + `)
	`

	catRows, err := db.DB.Query(catQuery)
	if err != nil {
		return nil, err
	}
	defer catRows.Close()

	for catRows.Next() {
		var postID, categoryID int
		var categoryName string
		err := catRows.Scan(&postID, &categoryID, &categoryName)
		if err != nil {
			return nil, err
		}
		if post, exists := postMap[postID]; exists {
			post.Cats = append(post.Cats, &Categories{ID: categoryID, Name: categoryName})
			post.LenCat = len(post.Cats) - 1
		}
	}

	likeQuery := `
		SELECT lp.post_id, lp.id, u.id AS user_id, u.name AS user_name
		FROM likepost AS lp
		JOIN users AS u ON lp.user_id = u.id
		WHERE lp.post_id IN (` + strings.Join(IntSliceToStringSlice(postIDs), ",") + `)
	`

	likeRows, err := db.DB.Query(likeQuery)
	if err != nil {
		return nil, err
	}
	defer likeRows.Close()

	for likeRows.Next() {
		var postID, likeID, userID int
		var userName string
		err := likeRows.Scan(&postID, &likeID, &userID, &userName)
		if err != nil {
			return nil, err
		}
		if post, exists := postMap[postID]; exists {
			post.Likes = append(post.Likes, &Like{ID: likeID, User: &User{ID: userID, Name: userName}})
		}
	}

	dislikeQuery := `
		SELECT dp.post_id, dp.id, u.id AS user_id, u.name AS user_name
		FROM dislikepost AS dp
		JOIN users AS u ON dp.user_id = u.id
		WHERE dp.post_id IN (` + strings.Join(IntSliceToStringSlice(postIDs), ",") + `)
	`

	dislikeRows, err := db.DB.Query(dislikeQuery)
	if err != nil {
		return nil, err
	}
	defer dislikeRows.Close()

	for dislikeRows.Next() {
		var postID, dislikeID, userID int
		var userName string
		err := dislikeRows.Scan(&postID, &dislikeID, &userID, &userName)
		if err != nil {
			return nil, err
		}
		if post, exists := postMap[postID]; exists {
			post.Dislikes = append(post.Dislikes, &Dislike{ID: dislikeID, User: &User{ID: userID, Name: userName}})
		}
	}

	posts := make([]*Post, 0, len(postMap))
	for _, post := range postMap {
		posts = append(posts, post)
	}

	return posts, nil
}

func IntSliceToStringSlice(ints []int) []string {
	strs := make([]string, len(ints))
	for i, v := range ints {
		strs[i] = strconv.Itoa(v)
	}
	return strs
}
