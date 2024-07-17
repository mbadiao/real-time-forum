package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"realTimeForum/backend/models"
)

func HandleLikedata(w http.ResponseWriter, like models.Like, r *http.Request) ([]models.Posts, models.ErrorPost) {
	GetStatusLike(like)
	return GetPost(), ErrorPost
}

func GetStatusLike(like models.Like) {
	actualStatus, err := GetExistingLike(like.PostID, like.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No existing like record, treat it as if the user has not liked the post
			err = InsertLike(like, true) // Insert a new like with 'liked' status
			if err != nil {
				fmt.Println("error insert like:", err)
				GetError("InternalServerError", 500)
				return
			}
			return
		} else {
			fmt.Println("error Get Statuslike:", err)
			GetError("InternalServerError", 500)
			return
		}
	}
	
	// Toggle the like status
	newLikeStatus := !actualStatus

	err = UpdateLike(like, newLikeStatus)
	if err != nil {
		fmt.Println("error update like:", err)
		GetError("InternalServerError", 500)
		return
	}
}

func GetExistingLike(postID, userID int) (bool, error) {
	query := `SELECT liked FROM LikesDislikes WHERE post_id = ? AND user_id = ?`
	var liked bool
	row := db.QueryRow(query, postID, userID)
	err := row.Scan(&liked)
	if err != nil {
		if err == sql.ErrNoRows {
			// No existing like found
			return false, sql.ErrNoRows
		}
		return false, err
	}

	return liked, nil
}

func InsertLike(like models.Like, status bool) error {
	query := `INSERT INTO LikesDislikes (post_id, user_id, liked) VALUES (?, ?, ?)`
	_, err := db.Exec(query, like.PostID, like.UserID, status)
	return err
}

func UpdateLike(like models.Like, status bool) error {
	query := `UPDATE LikesDislikes SET liked = ? WHERE post_id = ? AND user_id = ?`
	_, err := db.Exec(query, status, like.PostID, like.UserID)
	return err
}
