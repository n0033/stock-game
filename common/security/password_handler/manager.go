package password_handler

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

const DEFAULT_COST = 11

func Hash(password string) (bool, string) {

	var hashed_password, err = bcrypt.GenerateFromPassword([]byte(password), DEFAULT_COST)

	if err != nil {
		log.Fatal(err)
		return false, password
	}
	return true, string(hashed_password)
}

func Compare(password string, hashed_password string) bool {

	hash_byte := []byte(hashed_password)
	password_byte := []byte(password)

	err := bcrypt.CompareHashAndPassword(hash_byte, password_byte)

	return err == nil
}
