
package utilities

import (
	"errors"
	"log"
	"net/mail"
)

/*
validates email adrress
@params email
*/
func ValidateEmail(email string) (*string,error) {
	result, err := mail.ParseAddress(email)
	if err != nil {
		log.Printf("Email validation failed for: %s, Error: %s", email, err.Error())
        errStr:="your email "+email+" is not valid"
        return nil,errors.New(errStr)
	}
	return &result.Address,nil
}
