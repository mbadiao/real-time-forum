package controllers

import (
	"fmt"
	"net/http"
	"realTimeForum/backend/models"
	"realTimeForum/backend/utils"
	"unicode/utf8"
)

func HandleCommentData(w http.ResponseWriter, r *http.Request, comment models.Comment) ([]models.Posts, models.ErrorPost) {
	// ErrorPost.Status = false
	if IsValidComment(comment) {
		InsertComment(w, r, comment)
	}

	return GetPost(), ErrorPost
}

func InsertComment(w http.ResponseWriter, r *http.Request, comment models.Comment) {
	userID, err := GetCurrentUserIdModified(comment)
	if err != nil {
		fmt.Println("Error retrieving user ID:", err)
		GetError("Unauthorized", 401)
		return
	}

	_, err = db.Exec("INSERT INTO Comments (user_id, post_id, content) VALUES (?, ?, ?)", userID, comment.PostID, comment.Content)
	if err != nil {
		fmt.Println("Error inserting comment:", err)
		GetError("InternalServerError", 500)
		return
	}
}

func IsValidComment(comment models.Comment) bool {
	if utf8.RuneCountInString(comment.Content) > 500 {
		GetError("InvalidComment", 400)
		return false
	}
	return true
}

func GetComment(postID int) ([]models.Comments, error) {
	query := `SELECT comment_id, post_id, user_id, content, creation_date FROM Comments WHERE post_id = ? ORDER BY creation_date DESC`
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var allComments []models.Comments
	for rows.Next() {
		var comments models.Comments
		err := rows.Scan(&comments.CommentID, &comments.PostID, &comments.UserID, &comments.Content, &comments.Creation_date)
		if err != nil {
			fmt.Println("func GetallComments rowsScan", err)
			continue
		}

		author, err := GetPostAuthor(comments.UserID)
		if err != nil {
			fmt.Println("func  GetallComments Getpostauthor", err)
			continue
		}

		comments.Author = author
		comments.Formated_date = utils.FormatTimeAgo(comments.Creation_date)
		allComments = append(allComments, comments)
	}
	return allComments, nil
}

