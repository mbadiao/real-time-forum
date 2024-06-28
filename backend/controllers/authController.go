package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"realTimeForum/backend/models"
	"realTimeForum/backend/utils"
	"realTimeForum/data"
	"strconv"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func CreateCookie(w http.ResponseWriter) http.Cookie {
	Tokens, _ := uuid.NewV4()
	now := time.Now()
	expires := now.Add(time.Hour * 1)
	cookie := http.Cookie{
		Name:     "ForumCookie",
		Value:    Tokens.String(),
		Expires:  expires,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	return cookie
}

func GetCookieHandler(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("ForumCookie")
	if err != nil {
		return ""
	}
	return (cookie.Value)
}


func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// db := data.CreateTable()
	utils.FileService("index.html", w, "login")
	// CookieHandler(w, r, db)
}

var rgUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Permettre les connexions WebSocket depuis toutes les origines
		return true
	},
}

var lgUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Permettre les connexions WebSocket depuis toutes les origines
		return true
	},
}

func InsertRegistredUser(user models.UserRegister) error {
	db := data.CreateTable()
	// Prépare la requête d'insertion
	stmt, err := db.Prepare(`
        INSERT INTO Users (username, age, gender, firstname, lastname, email, password_hash, user_status)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Hacher le mot de passe de l'utilisateur
	hashedPassword, err := models.HashPassword(user.Password)
	if err != nil {
		return err
	}

	// Exécute la requête avec les valeurs de l'utilisateur
	_, err = stmt.Exec(
		user.Username,
		user.Age,
		user.Gender,
		user.FirstName,
		user.LastName,
		user.Email,
		hashedPassword,
		true,
	)
	return err
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := rgUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error", err)
		return
	}
	defer conn.Close()
	var registerData models.RegisterType
	for {
		messageType, message, err := conn.ReadMessage()

		if err != nil {
			log.Println("err", err)
			break
		}

		err = json.Unmarshal(message, &registerData)

		if err != nil {
			log.Println("111", err)
			return
		}

		userData := models.UserRegister{
			Username:  registerData.Username,
			Gender:    registerData.Gender,
			FirstName: registerData.Firstname,
			LastName:  registerData.Lastname,
			Email:     registerData.Email,
			Password:  registerData.Password,
		}
		userData.Age, err = strconv.Atoi(registerData.Age)
		if err != nil {
			fmt.Println("error sur la fonction atoi", err)
			return
		}

		err = InsertRegistredUser(userData)
		if err != nil {

			fmt.Println("error lors de l'insertion", err)
			errorType := utils.AnalyzeError(err)

			if errorType == "username" {

				response := map[string]string{"errorRegister":"Username already Exist"}
				responseJSON, _ := json.Marshal(response)
				if err := conn.WriteMessage(messageType, responseJSON); err != nil {
					log.Println("Write error:", err)
				}

			} else if errorType == "email" {
			
				response := map[string]string{"errorRegister":"Email already Exist"}
				responseJSON, _ := json.Marshal(response)
				if err := conn.WriteMessage(messageType, responseJSON); err != nil {
					log.Println("Write error:", err)
				}
			} else {
				fmt.Println("other error man")
			}
			return
		}

		_, err, cookie := InitSession(registerData.Username, w)
		if err != nil {
			fmt.Println("impossible d'initier la base de donnees")
			return
		}

		response := map[string]string{"status": "success", "cookie": cookie, "username":userData.Username, "firstname":userData.FirstName}
		responseJSON, _ := json.Marshal(response)
		if err := conn.WriteMessage(messageType, responseJSON); err != nil {
			log.Println("Write error:", err)
		}
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := lgUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("errorr", err)
		return
	}
	defer conn.Close()

	var loginData models.LoginType
	for {
		messageType, message, err := conn.ReadMessage()

		if err != nil {
			log.Println("errr", err)
			break
		}

		err = json.Unmarshal(message, &loginData)

		if err != nil {
			log.Println("111", err)
		}

		userData := models.UserLogin{
			EmailOrUsername: loginData.EmailOrUsername,
			Password:        loginData.Password,
		}


		isPresent, err, user := CheckExistenceOfUser(userData.EmailOrUsername, userData.Password)
		if err != nil {
			fmt.Println("the user does not exist", err)

			response := map[string]string{"errorLogin":"invalid identifants or password"}
				responseJSON, _ := json.Marshal(response)
				if err := conn.WriteMessage(messageType, responseJSON); err != nil {
					log.Println("Write error:", err)
				}

			return
		}

		if isPresent {
			_, err, cookie := InitSession(userData.EmailOrUsername, w)

			response := map[string]string{"status": "success", "cookie": cookie , "username":user.Username, "firstname":user.FirstName}
			responseJSON, _ := json.Marshal(response)
			if err := conn.WriteMessage(messageType, responseJSON); err != nil {
				log.Println("Write error:", err)
			}

			if err != nil {
				fmt.Println("impossible d'initier la base de donnees")
				return
			}
		}
	}
}

func CheckExistenceOfUser(emailOrUsername string, password string) (bool, error, models.UserRegister) {
	db := data.CreateTable()

	var user models.UserRegister

	query := `select user_id, username, email, password_hash, firstname from Users where email = ? or username = ?`

	err := db.QueryRow(query, emailOrUsername, emailOrUsername).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.FirstName)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("no row found")
			return false, err, user
		} else {
			fmt.Println("autre erreur", err)
			return false, err, user
		}

	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, err, user
	}
	return true, nil, user
}

func InitSession(emailOrUsername string, w http.ResponseWriter) (bool, error, string) {
	db := data.CreateTable()

	var id int

	query := `SELECT user_id FROM Users WHERE username = ? OR email= ?`

	err := db.QueryRow(query, emailOrUsername, emailOrUsername).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("no row found")
			return false, err, ""
		} else {
			fmt.Println("autre erreur", err)
			return false, err, ""
		}
	}

	_, err = CheckAndDeleteSession(id)
	if err != nil {
		fmt.Println("error while checking and deleting the user id")
		return false, err, ""
	}

	cookie := CreateCookie(w).Value

	err = InsertSession(id, cookie)

	if err != nil {
		fmt.Println("error lors de l'insertion de l'utilisation de insertSession")
		return false, err, ""
	}

	return true, nil, cookie
}

func CheckAndDeleteSession(userID int) (bool, error) {
	// Connexion à la base de données SQLite
	db := data.CreateTable()
	defer db.Close()

	var exists bool

	// Préparer la requête pour vérifier l'existence du user_id
	checkQuery := `SELECT EXISTS(SELECT 1 FROM Sessions WHERE user_id = ?)`
	err := db.QueryRow(checkQuery, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error querying the database: %w", err)
	}

	if !exists {
		// Le user_id n'existe pas dans la table session
		fmt.Println("user_id not found in session table")
		return false, nil
	}

	// Préparer et exécuter la requête pour supprimer le user_id
	deleteQuery := `DELETE FROM Sessions WHERE user_id = ?`
	stmtDelete, err := db.Prepare(deleteQuery)
	if err != nil {
		return false, fmt.Errorf("error preparing delete query: %w", err)
	}
	defer stmtDelete.Close()

	_, err = stmtDelete.Exec(userID)
	if err != nil {
		return false, fmt.Errorf("error deleting from session table: %w", err)
	}

	fmt.Println("user_id successfully deleted from session table")
	return true, nil
}


func CheckSession(cookieValue string) (bool, int, error) {
    db := data.CreateTable()

    var exists bool
    var userID int

    query := "SELECT user_id FROM sessions WHERE cookie_value = ? LIMIT 1"
    err := db.QueryRow(query, cookieValue).Scan(&userID)
    if err != nil {
        if err == sql.ErrNoRows {
            return false, 0, nil
        }
        return false, 0, err
    }

    exists = true
    return exists, userID, nil
}



func InsertSession(userID int, cookieValue string) error {
	// Connexion à la base de données SQLite
	db := data.CreateTable()
	defer db.Close()

	// Préparer la requête pour insérer un enregistrement dans la table Sessions
	insertQuery := `INSERT INTO Sessions (user_id, cookie_value) VALUES (?, ?)`
	stmtInsert, err := db.Prepare(insertQuery)
	if err != nil {
		return fmt.Errorf("error preparing insert query: %w", err)
	}
	defer stmtInsert.Close()

	// Exécuter la requête d'insertion
	_, err = stmtInsert.Exec(userID, cookieValue)
	if err != nil {
		return fmt.Errorf("error inserting into Sessions table: %w", err)
	}

	fmt.Println("Session successfully inserted")
	return nil
}

func IsRightSession(w http.ResponseWriter, r *http.Request){
	var reqBody models.RequestBody
	err:= json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		fmt.Println("error decoding", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	var resBody models.ResponseBody
	isPresent, userid, err:=CheckSession(reqBody.CookieValue)

	if err!=nil {
		fmt.Println("error why testing if the sessions is the right")
		return
	}

	if isPresent{
		username,firstname,_:= GetUsernameAndFirstname(userid)
		resBody.ProcessedValue="ok"
		resBody.Username=username
		resBody.Firstname=firstname
	}else{
		resBody.ProcessedValue="ko"
	}

	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resBody)
}

func GetUsernameAndFirstname(userID int) (string, string, error) {
    db := data.CreateTable()

    var username string
    var firstname string

    query := "SELECT username, firstname FROM Users WHERE user_id = ? LIMIT 1"
    err := db.QueryRow(query, userID).Scan(&username, &firstname)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", "", nil
        }
        return "", "", err
    }

    return username, firstname, nil
}
