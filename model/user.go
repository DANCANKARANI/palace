package model

import (
	"errors"

	"github.com/dancankarani/palace/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

/*
gets user using user ID
*/
type ResponseUser struct{
	ID uuid.UUID 		`json:"id"`
	FullName string 	`json:"full_name"`
	PhoneNumber string 	`json:"phone_number"`
	Email string 		`json:"email"`
	ProfilePhotoPath string	`json:"profile_photo_path"`
}

func GetOneUSer(c *fiber.Ctx)(*ResponseUser,error){
	id,err:=GetAuthUserID(c)
	if err != nil {
		return nil,errors.New("failed to get user's id:"+err.Error())
	}
	user := ResponseUser{}
	err = db.First(&User{},"id = ?",id).Scan(&user).Error
	if err != nil {
		return nil,errors.New("failed to get user details:"+err.Error())
	}
	return &user,nil
}
//gets all the users
func GetAllUsersDetails(c *fiber.Ctx)(*[]ResponseUser,error){
	response:=[]ResponseUser{}
	err := db.Model(&User{}).Scan(&response).Error
	if err != nil {
		return nil,errors.New("failed to get users:"+err.Error())
	}
	return &response,nil
}

// UpdateUser updates the user by ID and logs the changes.
func UpdateUser(c *fiber.Ctx) (*ResponseUser, error) {
    // Get the authenticated user ID
    id, err := GetAuthUserID(c)
    if err != nil {
        return nil, errors.New("failed to get user's ID: " + err.Error())
    }

    // Parse the request body into a User struct
    var body User
    if err := c.BodyParser(&body); err != nil {
        return nil, errors.New("failed to parse: " + err.Error())
    }

	//validate phone number
	if body.PhoneNumber !=""{
		_,err :=utilities.ValidatePhoneNumber(body.PhoneNumber,"KE")
		if err != nil{
			return nil, err
		}
		exist,_,err:=UserExist(c,body.PhoneNumber)
		if err != nil{
			return nil, err
		}
		if exist{
			err_str := "user with phone:"+body.PhoneNumber+" already exist"
			return nil, errors.New(err_str)
		}
	}

	//validate email
	if body.Email !=""{
		_, err := utilities.ValidateEmail(body.Email)
		if err != nil{
			return nil, err
		}
	}

	//hash password
	if body.Password != ""{
		hashed_password, err:= utilities.HashPassword(body.Password)
		if err != nil{
			return nil,err
		}
		body.Password= hashed_password
	}
    // Fetch the current user record to get old values
    oldValues := new(User)
    if err := db.First(&oldValues, "id = ?", id).Error; err != nil {
        return nil, errors.New("failed to fetch current user: " + err.Error())
    }
	response := new(ResponseUser)

    // Update the user record
    if err := db.Model(&oldValues).Updates(&body).Scan(response).Error; err != nil {
        return nil, errors.New("error in updating the user: " + err.Error())
    }


    // Audit logs

    return response, nil
}
func MapUserToResponse(user User) ResponseUser {
    return ResponseUser{
        ID:          user.ID,
        FullName:    user.FullName,
        PhoneNumber: user.PhoneNumber,
        Email:       user.Email,
    }
}
