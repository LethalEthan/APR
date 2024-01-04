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
	DateOfBirth   string `json:"dob"`
}

type PersonalLog struct {
	LogID      int    `json:"logid"`
	MemberID   int    `json:"memberid"`
	LogContent string `json:"logcontent"`
	LogRoutine string `json:"logroutine"`
	LogDate    string `json:"logdate"`
}

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./Gym_Management_DB.db") // Open the created SQLite File
	if err != nil {
		println(err.Error())
	}
}

func CloseDB() {
	dbMutex.Lock()
	db.Close()
	dbMutex.Unlock()
}

// QueryMemberInfoByID - Retrieve member information by ID
func QueryMemberInfoByID(MemberID int) (Member, error) {
	dbMutex.Lock()
	rows, err := db.Query("SELECT MemberID,FirstName,LastName,Email,Password,WellnessGoals,DOB from Member WHERE MemberID=\"" + strconv.Itoa(MemberID) + "\"")
	if err != nil {
		log.Println(err)
		return Member{}, err
	}
	dbMutex.Unlock()
	defer rows.Close()
	for rows.Next() {
		M := new(Member)
		err = rows.Scan(&M.MemberID, &M.FirstName, &M.LastName, &M.Email, &M.Password, &M.WellnessGoals, &M.DateOfBirth)
		if err != nil {
			log.Println(err)
			return Member{}, err
		}
		fmt.Println(M.MemberID, M.FirstName, M.LastName, M.Email, M.Password, M.WellnessGoals, M.DateOfBirth)
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
	rows, err := db.Query("SELECT MemberID,FirstName,LastName,Email,Password,WellnessGoals,DOB from Member")
	if err != nil {
		log.Println(err)
		return []Member{}, err
	}
	dbMutex.Unlock()
	defer rows.Close()
	var Members = make([]Member, 0, 10)
	for rows.Next() { // Go through rows and put them into Member objects
		M := new(Member)
		err = rows.Scan(&M.MemberID, &M.FirstName, &M.LastName, &M.Email, &M.Password, &M.WellnessGoals, &M.DateOfBirth)
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
	_, err := db.Exec("INSERT INTO Member (FirstName,LastName,Email,Password,DOB) VALUES( \"" + M.FirstName + "\", \"" + M.LastName + "\", \"" + M.Email + "\", \"" + M.Password + "\", \"" + M.DateOfBirth + "\")")
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
func QueryTokenExists(Token string) (string, bool, error) {
	dbMutex.Lock()
	rows, err := db.Query("SELECT Token from Authentication WHERE Token=\"" + Token + "\"")
	if err != nil {
		log.Println(err)
		return "", false, err
	}
	dbMutex.Unlock()
	defer rows.Close()
	for rows.Next() {
		var token string
		err = rows.Scan(&token)
		if err != nil {
			log.Println(err)
			return "", false, err
		}
		fmt.Println(token)
		return token, true, nil
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return "", false, err
	}
	return "", false, errors.New("No token found!")
}

func QueryPersonalLogs() ([]PersonalLog, error) {
	dbMutex.Lock()
	rows, err := db.Query("SELECT LogID, MemberID, LogContent, LogRoutines, LogDate from PersonalLog")
	if err != nil {
		log.Println(err)
		return []PersonalLog{}, err
	}
	dbMutex.Unlock()
	defer rows.Close()
	var PersonalLogs = make([]PersonalLog, 0, 10)
	for rows.Next() {
		var PL PersonalLog
		err = rows.Scan(&PL.LogID, &PL.MemberID, &PL.LogContent, &PL.LogRoutine, &PL.LogDate)
		if err != nil {
			log.Println(err)
			return []PersonalLog{}, err
		}
		fmt.Println(PL.LogID, PL.MemberID, PL.LogContent, PL.LogRoutine, PL.LogDate)
		PersonalLogs = append(PersonalLogs, PL)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return []PersonalLog{}, err
	}
	return PersonalLogs, err
}

func QueryPersonalLogsByID(MemberID int) ([]PersonalLog, error) {
	dbMutex.Lock()
	rows, err := db.Query("SELECT LogID, MemberID, LogContent, LogRoutines, LogDate from PersonalLog WHERE MemberID=" + "\"" + strconv.Itoa(MemberID) + "\"")
	if err != nil {
		log.Println(err)
		return []PersonalLog{}, err
	}
	dbMutex.Unlock()
	defer rows.Close()
	var PersonalLogs = make([]PersonalLog, 0, 10)
	for rows.Next() {
		var PL PersonalLog
		err = rows.Scan(&PL.LogID, &PL.MemberID, &PL.LogContent, &PL.LogRoutine, &PL.LogDate)
		if err != nil {
			log.Println(err)
			return []PersonalLog{}, err
		}
		fmt.Println(PL.LogID, PL.MemberID, PL.LogContent, PL.LogRoutine, PL.LogDate)
		PersonalLogs = append(PersonalLogs, PL)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return []PersonalLog{}, err
	}
	return PersonalLogs, err
}

func UpdateMember(M Member) error {
	dbMutex.Lock()
	_, err := db.Exec("UPDATE Member SET FirstName=\"" + M.FirstName + "\",LastName=\"" + M.LastName + "\",Email=\"" + M.Email + "\",Password=\"" + M.Password + "\",WellnessGoals=\"" + M.WellnessGoals + "\",DOB=\"" + M.DateOfBirth + "\"" + " WHERE MemberID=" + strconv.Itoa(M.MemberID))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
