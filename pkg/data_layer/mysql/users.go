package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"

	"golang.org/x/crypto/bcrypt"

	datalayer "raven.net/snippetbox/pkg/data_layer"
)

type UserModel struct {
	DB *sql.DB
}

// to insert new user
func (m *UserModel) Insert(name, email, password string) error {

	//inserting a hashed (bcrypt) of plain text-passwd
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	//sql insert statement

	stmt := `INSERT INTO users (name, email, hashed_password, created_at) VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mySQLErr *mysql.MySQLError
		if errors.As(err, &mySQLErr) {
			if mySQLErr.Number == 1062 && strings.Contains(mySQLErr.Message, "users_uc_email") {
				return datalayer.ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

// to verify if mail addr exist; return userID if it exists
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPasswd []byte
	/*
	   sql statement to retrieve id, hashedpasswd from users
	   check if email exists or user is not active
	*/
	stmt := "SELECT id, hashed_password FROM users WHERE email = ? AND active = TRUE"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPasswd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, datalayer.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	//check matching plain-txt passwd and hashedpasswd
	err = bcrypt.CompareHashAndPassword(hashedPasswd, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, datalayer.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

// to fetch user with their ID
func (m *UserModel) Get(id int) (*datalayer.User, error) {
	user := &datalayer.User{}
	stmt := `SELECT id, name, email, created_at, active FROM users WHERE id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&user.ID, &user.Name, &user.Email, &user.Created_at, &user.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datalayer.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return user, nil
}
func (m *UserModel) changePassword(id int, currentPasswd, newPasswd string) error {
	u, err := m.Get(id)
	if err != nil {

	}
}
