package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"realTimeForum/backend/models"
	"realTimeForum/backend/utils"
	"realTimeForum/data"
	"unicode/utf8"
)

var db = data.CreateTable()
var ErrorPost models.ErrorPost

func HandlePostData(w http.ResponseWriter, post models.Post, r *http.Request) ([]models.Posts, models.ErrorPost) {
	ErrorPost.Status = false
	if IsValidPost(post) {
		UploadImage(&post, r)
		InsertPost(w, post, r)
	}
	return GetPost(), ErrorPost
}

func InsertPost(w http.ResponseWriter, post models.Post, r *http.Request) {
	categoryJSON, err := json.Marshal(post.Category)
	if err != nil {
		fmt.Println("Erreur lors de la conversion du tableau de catégories en JSON:", err)
		GetError("InternalServerError", 500)
		return
	}

	userID, err := GetCurrentUserId(w, r)
	if err != nil {
		fmt.Println("Error retrieving user ID:", err)
		GetError("Unauthorized", 401)
		return
	}

	_, err = db.Exec("INSERT INTO Posts (user_id, title, PhotoURL, content, category) VALUES (?, ?, ?, ?, ?)", userID, post.Title, post.Image.Path, post.Content, string(categoryJSON))
	if err != nil {
		fmt.Println("Erreur lors de l'insertion du post:", err)
		GetError("InternalServerError", 500)
		return
	}
	var lastInsertID int
	err = db.QueryRow("SELECT last_insert_rowid()").Scan(&lastInsertID)
	if err != nil {
		fmt.Println("Erreur lors de la récupération de l'ID du dernier post inséré:", err)
		GetError("InternalServerError", 500)
		return
	}
}

func GetCurrentUserId(w http.ResponseWriter, r *http.Request) (int, error) {
	var userID int
	cookie := GetCookieHandler(w, r)

	userID, err := utils.GetUserIDFromDB(db, cookie)
	if err != nil {
		return 0, fmt.Errorf("failed to get user ID from session: %v", err)
	}
	return userID, nil
}

func GetCurrentUserIdModified(comment models.Comment) (int, error) {
	var userID int
	cookie := comment.Cookie

	userID, err := utils.GetUserIDFromDB(db, cookie)
	if err != nil {
		return 0, fmt.Errorf("failed to get user ID from session: %v", err)
	}
	return userID, nil
}

// func IsValidImage(post models.Post, r *http.Request) bool {
// 	if post.Image != "nophoto" {
// 		// don't delete
// 		// isImage := UploadImage(&post, r)
// 		return true
// 	}
// 	return true
// }

func UploadImage(post *models.Post, r *http.Request) bool {
	if post.Image.Status {
		dirPath := "../frontend/static/upload"
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			err := os.MkdirAll(dirPath, 0755)
			if err != nil {
				fmt.Println("Erreur lors de la création du répertoire:", err)
				return false
			}
		}

		srcFile, err := os.Open(post.Image.Path)
		if err != nil {
			fmt.Println("Erreur lors de l'ouverture du fichier source:", err)
			return false
		}
		defer srcFile.Close()

		tempFile, err := os.CreateTemp(dirPath, "upload-*"+filepath.Ext(post.Image.Name))
		if err != nil {
			fmt.Println("Erreur lors de la création du fichier temporaire:", err)
			return false
		}
		defer tempFile.Close()

		_, err = io.Copy(tempFile, srcFile)
		if err != nil {
			fmt.Println("Erreur lors de la copie du fichier:", err)
			return false
		}
		post.Image.Path = tempFile.Name()
		return true
	} else {
		return true
	}
}

func IsValidPost(post models.Post) bool {
	if post.Title == "" || post.Content == "" || len(post.Category) == 0 {
		GetError("One or more fields are empty", 400)
		return false
	}
	if utf8.RuneCountInString(post.Title) > 30 {
		GetError("too long title", 400)
		return false
	}
	return true
}

func GetPost() []models.Posts {
	allpost, err := GetAllPosts()
	if err != nil {
		fmt.Println("func GetPost allpost", err)
		GetError("InternalServerError", 500)
		return nil
	}
	return allpost
}

func GetAllPosts() ([]models.Posts, error) {
	query := `SELECT post_id, user_id, title, PhotoURL, content, category, creation_date FROM Posts ORDER BY creation_date DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allPosts []models.Posts
	for rows.Next() {
		var post models.Posts
		var categoryJSON string
		err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Photo_url, &post.Content, &categoryJSON, &post.Creation_date)
		if err != nil {
			fmt.Println("func GetAllposts rowsScan", err)
			continue
		}

		err = json.Unmarshal([]byte(categoryJSON), &post.Categories)
		if err != nil {
			fmt.Println("Erreur lors de la conversion de JSON en catégorie:", err)
			continue
		}

		author, err := GetPostAuthor(post.UserID)
		if err != nil {
			fmt.Println("func  GetAllposts Getpostauthor", err)
			continue
		}
		post.Author = author

		comment, err := GetComment(post.PostID)
		if err != nil {
			fmt.Println("func  GetAllComment-GetAllPost", err)
			continue
		}

		nbrcomment, err := GetNbrComment(post.PostID)
		if err != nil {
			fmt.Println("func  GetAllComment-GetNbrComment", err)
			continue
		}

		nbrlike, err := GetNbrLike(post.PostID)
		if err != nil {
			fmt.Println("func GetAllComment-GetNbrLike", err)
			continue
		}

		postlike, err := GetStatus(post.PostID, post.UserID)
		if err != nil {
			fmt.Println("func GetAllComment-GetStatusLike", err)
			continue
		}

		post.Like_nbr = nbrlike
		post.Comments_nbr = nbrcomment
		post.Comments = comment
		post.Like_status = postlike
		post.Formated_date = utils.FormatTimeAgo(post.Creation_date)
		allPosts = append(allPosts, post)
	}
	return allPosts, nil
}

func GetPostAuthor(userID int) (models.Author, error) {
	query := `SELECT firstname, lastname, username FROM Users WHERE user_id = ?`
	row := db.QueryRow(query, userID)
	var author models.Author
	err := row.Scan(&author.Firstname, &author.Lastname, &author.Username)
	if err != nil {
		return author, err
	}
	return author, nil
}

func GetError(msg string, code int) {
	ErrorPost.Status = true
	ErrorPost.Code = code
	ErrorPost.Message = msg
}

func GetNbrComment(postId int) (int, error) {
	var nbrComments int
	query := `SELECT COUNT(*) FROM Comments WHERE post_id = ?`
	err := db.QueryRow(query, postId).Scan(&nbrComments)
	if err != nil {
		return 0, err
	}
	return nbrComments, nil
}

func GetNbrLike(postId int) (int, error) {
	var nbrLikes int
	query := `SELECT COUNT(*) FROM LikesDislikes WHERE post_id = ? AND liked = TRUE`
	err := db.QueryRow(query, postId).Scan(&nbrLikes)
	if err != nil {
		return 0, err
	}
	return nbrLikes, nil
}

func GetStatus(postID, userID int) (bool, error) {
	query := `SELECT liked FROM LikesDislikes WHERE post_id = ? AND user_id = ?`
	var liked bool
	row := db.QueryRow(query, postID, userID)
	err := row.Scan(&liked)
	if err != nil {
		if err == sql.ErrNoRows {
			// No existing like found
			return false, nil
		}
		return false, err
	}

	return liked, nil
}
