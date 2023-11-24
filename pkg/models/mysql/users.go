package mysql

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/harmlessprince/snippetboxapp/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserModel struct {
	DB *sql.DB
}

func (userModel *UserModel) Insert(name string, email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil
	}
	statement := `INSERT INTO users (name, email, hashed_password, created) VALUES(?, ?, ?, UTC_TIMESTAMP())`
	_, err = userModel.DB.Exec(statement, name, email, string(hashedPassword))
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "") {
				return models.ErrDuplicateEmail
			}
		}
	}
	return err
}

func (userModel *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	row := userModel.DB.QueryRow("SELECT  id, hashed_password FROM users WHERE email = ?", email)
	err := row.Scan(&id, &hashedPassword)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, nil
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, nil
	}
	return id, nil

}

func (userModel *UserModel) Get(id int) (*models.User, error) {
	user := &models.User{}
	statement := `SELECT id, name, email, created FROM users WHERE id = ?`
	err := userModel.DB.QueryRow(statement, id).Scan(&user.ID, &user.Name, &user.Email, &user.Created)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return user, nil
}
