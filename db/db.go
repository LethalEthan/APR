package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var dbMutex sync.Mutex

type Member struct {
	MemberID      int    `json:"memberid"`
	FirstName     string `json:"firstname"`
	LastName      string `json:"lastname"`
	Email         string `json:"email"`
	Password      string `json:"-"`
	WellnessGoals string `json:"wellnessgoals"`
}

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./Gym_Management_DB.db") // Open the created SQLite File
	if err != nil {
		println(err.Error())
	}
	// defer db.Close() // Defer Closing the database
}

// QueryMemberInfoByID - Retrieve member information by ID
func QueryMemberInfoByID(MemberID int) (Member, error) {
	dbMutex.Lock()
	rows, err := db.Query("SELECT MemberID,FirstName,LastName,Email,Password,WellnessGoals from Member WHERE MemberID=\"" + strconv.Itoa(MemberID) + "\"")
	if err != nil {
		log.Println(err)
		return Member{}, err
	}
	dbMutex.Unlock()
	defer rows.Close()
	for rows.Next() {
		M := new(Member)
		err = rows.Scan(&M.MemberID, &M.FirstName, &M.LastName, &M.Email, &M.Password, &M.WellnessGoals)
		if err != nil {
			log.Println(err)
			return Member{}, err
		}
		fmt.Println(M.MemberID, M.FirstName, M.LastName, M.Email, M.Password, M.WellnessGoals)
		return *M, nil
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}
	return Member{}, err
}

// QueryAllMember - We get all the members from the database
func QueryAllMember() ([]Member, error) {
	dbMutex.Lock()
	rows, err := db.Query("SELECT MemberID,FirstName,LastName,Email,Password,WellnessGoals from Member")
	if err != nil {
		log.Println(err)
		return []Member{}, err
	}
	dbMutex.Unlock()
	defer rows.Close()
	var Members = make([]Member, 0, 10)
	for rows.Next() { // Go through rows and put them into Member objects
		M := new(Member)
		err = rows.Scan(&M.MemberID, &M.FirstName, &M.LastName, &M.Email, &M.Password, &M.WellnessGoals)
		if err != nil {
			log.Println(err)
			return []Member{}, err
		}

		log.Println(M.MemberID, M.FirstName, M.LastName, M.Email, M.WellnessGoals)
		Members = append(Members, *M)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return []Member{}, err
	}
	return Members, nil
}

func QueryMemberIDByToken(Token string) (int, error) {
	dbMutex.Lock()
	rows, err := db.Query("SELECT MemberID from Authentication WHERE Token=\"" + Token + "\"")
	if err != nil {
		log.Println(err)
		return 0, err
	}
	dbMutex.Unlock()
	defer rows.Close()
	for rows.Next() {
		var MemberID int
		err = rows.Scan(&MemberID)
		if err != nil {
			log.Println(err)
			return 0, err
		}
		fmt.Println(MemberID)
		return MemberID, nil
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return 0, errors.New("No MemberID found for token given!")
}

func InsertMember(M Member) error {
	dbMutex.Lock()
	_, err := db.Exec("INSERT INTO Member (FirstName,LastName,Email,Password) VALUES( \"" + M.FirstName + "\", \"" + M.LastName + "\", \"" + M.Email + "\", \"" + M.Password + "\" )")
	if err != nil {
		fmt.Println(err)
		return err
	}
	dbMutex.Unlock()
	return nil
}

func InsertToken(Token string, MemberID string) error {
	dbMutex.Lock()
	_, err := db.Exec("INSERT INTO Authentication (Token,MemberID) VALUES(\"" + Token + "\", \"" + MemberID + "\")")
	if err != nil {
		log.Println(err)
		return err
	}
	dbMutex.Unlock()
	return nil
}

func QueryMemberID(Email, Password string) (int, error) {
	dbMutex.Lock()
	rows, err := db.Query("SELECT MemberID from Member WHERE Email=\"" + Email + "\" AND Password=\"" + Password + "\"")
	if err != nil {
		log.Println(err)
		return 0, err
	}
	dbMutex.Unlock()
	defer rows.Close()
	for rows.Next() { // for rows but we only use the first one, can be change to if but since we return early the compiler emits
		var MemberID int
		err = rows.Scan(&MemberID)
		if err != nil {
			log.Println(err)
			return 0, err
		}
		fmt.Println(MemberID)
		return MemberID, nil
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return 0, errors.New("No memberID found")
}

func DeleteAllTokens() error {
	dbMutex.Lock()
	_, err := db.Exec("DELETE FROM Authentication")
	if err != nil {
		log.Println(err)
		return err
	}
	dbMutex.Unlock()
	return nil
}

func RevokeToken(token string) error {
	dbMutex.Lock()
	_, err := db.Exec("DELETE FROM Authentication WHERE Token=\"" + token + "\"")
	if err != nil {
		log.Println(err)
		return err
	}
	dbMutex.Unlock()
	return nil
}

// QueryTokenExists - Checks if token exists in DB
func QueryTokenExists(Token string) (string, error) {
	dbMutex.Lock()
	rows, err := db.Query("SELECT Token from Authentication WHERE Token=\"" + Token + "\"")
	if err != nil {
		log.Println(err)
		return "", err
	}
	dbMutex.Unlock()
	defer rows.Close()
	for rows.Next() {
		var token string
		err = rows.Scan(&token)
		if err != nil {
			log.Println(err)
			return "", err
		}
		fmt.Println(token)
		return token, nil
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return "", err
	}
	return "", errors.New("No token found!")
}
