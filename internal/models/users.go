package models

import "database/sql"
import "errors"
import "github.com/go-sql-driver/mysql"
import "golang.org/x/crypto/bcrypt"
import "strings"
import "time"

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// interacts with the database on behalf of the user model, so thus takes a ptr and is 8b
type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Insert(name, email, password string) error {
	hpass, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `
		INSERT INTO users (name, email, hashed_password, created)
		VALUES(?, ?, ?, UTC_TIMESTAMP())
	`
	_, err = m.DB.Exec(stmt, name, email, string(hpass))
	if err != nil {
		var mySQLError *mysql.MySQLError //what?
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Exists(id int) bool {
	return false
}
