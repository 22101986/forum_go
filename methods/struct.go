package methods

type User struct {
	ID         int
	UUID       string
	Name       string
	Email      string
	Password   string
	Picture    []byte
	Role       int
	IsExternal int
}
type Post struct {
	ID        int
	Title     string
	Content   string
	Date      string
	User      *User
	IsLike    bool
	IsDislike bool
	Comments  []*Comment
	Cats      []*Categories
	Likes     []*Like
	Dislikes  []*Dislike
	LenCat    int
	Blob      []byte
}
type Comment struct {
	ID        int
	Content   string
	Date      string
	User      *User
	PostID    int
	IsLike    bool
	IsDislike bool
	Likes     []*Like
	Dislikes  []*Dislike
}
type Like struct {
	ID   int
	Date string
	User *User
}
type Dislike struct {
	ID   int
	Date string
	User *User
}
type Categories struct {
	ID   int
	Name string
}
type Notif struct {
	ID       int
	Date     string
	Type     int
	UserTo   *User
	UserFrom *User
	Post     *Post
	Comment  *Comment
}
type Activity struct { // types 1 = post, 2 = comment, 3 = like, 4 = dislike
	Post    *Post
	Comment *Comment
	Like    *Like
	Dislike *Dislike
	Type    int
}
