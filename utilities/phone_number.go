package utilities

import (
	"errors"
	"strings"

	"github.com/ttacon/libphonenumber"
)

func ValidatePhoneNumber(phone_no,countrycode string) (string, error){
	num,err := libphonenumber.Parse(phone_no,countrycode)
	if err != nil {
		//number is invalid
		err_string :="Your phone number " + phone_no + " is invalid"
		return "", errors.New(err_string)
	}
	//Format phone number to E164 format
	formated_num := libphonenumber.Format(num, libphonenumber.E164)
	isvalid_no := libphonenumber.IsValidNumber(num)
	valid_no := formated_num

	if !isvalid_no {
		//number is not valid
		err_string := "Your phone number "+ phone_no + " is invalid"
		return "", errors.New(err_string)
	}
	return strings.Split(valid_no, "+")[1],nil
}