package server

import (
	"database/sql"
	"fmt"
)

type Database interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	Close() error
}

type UserDB struct {
	sqlOpenFunc sqlOpener
	db          Database
}

type sqlOpener func(driverName, dataSourceName string) (*sql.DB, error)

func NewUserDB(sqlOpenFunc sqlOpener) *UserDB {
	return &UserDB{sqlOpenFunc: sqlOpenFunc}
}

func (u *UserDB) connect(datasource, user, password, dbName string) error {
	// connect to mysql db named userdb
	var err error
	u.db, err = u.sqlOpenFunc("mysql", "user:pass@/userdb")
	if err != nil {
		return err
	}

	return nil
}

func (u *UserDB) GetName(id int) (string, error) {
	query := "select name from users where id=?"

	rows, err := u.db.Query(query, id)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var name string
	for rows.Next() {
		if err := rows.Scan(&name); err != nil {
			return "", err
		}
	}

	if name == "" {
		return "", fmt.Errorf("user does not exist. user_id=%d", id)
	}

	return name, nil
}
