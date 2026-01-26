package cli

import (
	"log"

	"github.com/google/uuid"
	"github.com/joseph0x45/goutils"
	"github.com/joseph0x45/vis/db"
	"github.com/joseph0x45/vis/models"
)

func CreateUser(username, password string, conn *db.Conn) {
	if username == "" || password == "" {
		log.Println("Username and Password can not be empty")
		return
	}
	hash, err := goutils.HashPassword(password)
	if err != nil {
		log.Println(err.Error())
		return
	}
	user := &models.User{
		ID:       uuid.NewString(),
		Username: username,
		Password: hash,
	}
	if err := conn.InsertUser(user); err != nil {
		log.Println(err.Error())
		return
	}
	log.Printf("User %s created!\n", username)
}
