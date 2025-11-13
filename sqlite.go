package sqlite

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Filename = ""
)

type Userdata struct {
	ID          int
	Username    string
	Name        string
	Surname     string
	Description string
}

func openConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", Filename)
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}
	return db, nil
}

func exists(db *sql.DB, username string) (bool, error) {
	username = strings.ToLower(username)

	var id int
	row := db.QueryRow(`SELECT id FROM users WHERE username = ?`, username)
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("exists: %w", err)
	}

	return true, nil
}

func AddUser(u Userdata) error {
	u.Username = strings.ToLower(u.Username)

	db, err := openConnection()
	if err != nil {
		return fmt.Errorf("adduser: %w", err)
	}
	defer db.Close()

	exist, err := exists(db, u.Username)
	if err != nil {
		return fmt.Errorf("adduser: %w", err)
	}
	if exist {
		return fmt.Errorf("adduser: user %s already exists", u.Username)
	}

	result, err := db.Exec(`INSERT INTO users VALUES (NULL,?)`, u.Username)
	if err != nil {
		return fmt.Errorf("adduser: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("adduser: %w", err)
	}

	_, err = db.Exec(`INSERT INTO userdata VALUES (?,?,?,?)`, id, u.Name, u.Surname, u.Description)
	if err != nil {
		return fmt.Errorf("adduser: %w", err)
	}

	return nil
}

func DeleteUser(id int) error {
	db, err := openConnection()
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM userdata WHERE userid = ?`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	_, err = db.Exec(`DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}

func ListUsers() ([]Userdata, error) {
	db, err := openConnection()
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT ID, Username, Name, Surname , Description
		FROM Users, Userdata WHERE Users.ID = Userdata.UserID`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []Userdata
	for rows.Next() {
		user := Userdata{}
		err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.Name,
			&user.Surname,
			&user.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("list users: %w", err)
		}
		data = append(data, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	return data, nil
}

func UpdateUser(u Userdata) error {
	u.Username = strings.ToLower(u.Username)

	db, err := openConnection()
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	defer db.Close()

	exist, err := exists(db, u.Username)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	if !exist {
		return fmt.Errorf("user %s doesn't exist in db", u.Username)
	}

	var id int
	err = db.QueryRow(`SELECT id FROM users WHERE username = ?`, u.Username).Scan(&id)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	_, err = db.Exec(`UPDATE Userdata SET Name = ?, Surname = ?, Description = ?
		WHERE UserID = ?`, u.Name, u.Surname, u.Description, id)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}
