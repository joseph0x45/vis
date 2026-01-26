package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/joseph0x45/vis/models"
)

func (c *Conn) InsertUser(user *models.User) error {
	const query = `
    insert into users (
      id, username, password
    )
    values (
      :id, :username, :password
    )
    on conflict do nothing;
  `
	if _, err := c.db.NamedExec(query, user); err != nil {
		return fmt.Errorf("Error while inserting user: %w", err)
	}
	return nil
}

func (c *Conn) GetUserBy(by, value string) (*models.User, error) {
	user := &models.User{}
	var query string
	if by == "id" {
		query = "select * from users where id=?"
	} else {
		query = "select * from users where username=?"
	}
	if err := c.db.Get(user, query, value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Error while getting user: %w", err)
	}
	return user, nil
}
