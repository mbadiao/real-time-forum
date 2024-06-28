package models

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type RegisterType struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Age       string `json:"age"`
	Gender    string `json:"gender"`
}

type LoginType struct {
	EmailOrUsername string `json:"emailOrUsername"`
	Password        string `json:"password"`
}

type UserLogin struct {
	EmailOrUsername string
	Password        string
}

type UserRegister struct {
	Username  string
	Age       int
	Gender    string
	FirstName string
	LastName  string
	Email     string
	Password  string
	Id        int
}

type Post struct {
	Type     string   `json:"type"`
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Category []string `json:"category"`
	Image    Image    `json:"image"`
}

type Image struct {
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	Type   string `json:"type"`
	Status bool   `json:"status"`
	Path   string `json:"path"`
}

type Like struct {
	PostID  int    `json:"postid"`
	UserID  int    `json:"userid"`
}


type Author struct {
	Firstname string
	Lastname  string
	Username  string
}

type Posts struct {
	PostID         int
	UserID         int
	Author         Author
	Title          string
	Content        string
	Photo_url      string
	Creation_date  time.Time
	Formated_date  string
	Categories     []string
	Like_nbr       int
	Dislike_nbr    int
	Comments_nbr   int
	Comments       []Comments
	Like_status    bool
	Dislike_status bool
}

type Comments struct {
	CommentID     int
	UserID        int
	PostID        int
	Content       string
	Author        Author
	Creation_date time.Time
	Formated_date string
}

type Comment struct {
	PostID  int    `json:"postid"`
	UserID  int    `json:"userid"`
	Content string `json:"content"`
	Cookie string `json:"cookie"`
}

type ErrorPost struct {
	Status  bool
	Code    int
	Message string
}

type ServicePost struct {
	Posts  []Posts
	Errors ErrorPost
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

type RequestBody struct {
    CookieValue string `json:"cookieValue"`
}

type ResponseBody struct {
    ProcessedValue string `json:"processedValue"`
	Firstname string `json:"firstname"`
	Username string `json:"username"`
}