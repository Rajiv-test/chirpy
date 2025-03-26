package auth

import (

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string,error){

	hashedPass,err  := bcrypt.GenerateFromPassword([]byte(password),10)
	if err != nil {
		return "", err
	}
	return string(hashedPass),err
}


func CheckPasswordHash(hash,password string) error{
	err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	if err != nil {
		return err
	}
	return nil
}