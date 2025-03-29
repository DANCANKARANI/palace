package user

import (
	"log"
	"time"

	"github.com/dancankarani/palace/database"
	"github.com/dancankarani/palace/model"
	"github.com/dancankarani/palace/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)


var db =database.ConnectDB()
var country_code = "KE"

func CreateUserAccount(c *fiber.Ctx) error {
	db.AutoMigrate(&model.User{})
	//generating new id
	id := uuid.New()
	user:=model.User{}
	if err :=c.BodyParser(&user); err != nil {
		log.Println(err.Error())

		return utilities.ShowError(c,"failed to create account", fiber.StatusInternalServerError)
	}

	//validate email address
	_,err:=utilities.ValidateEmail(user.Email)
	if err != nil {
		return utilities.ShowError(c,err.Error(),fiber.StatusInternalServerError)
	}
	//check email existence
	emailExist,_,err := model.EmailExist(c,user.Email)
	if err != nil{
		return utilities.ShowError(c,err.Error(),fiber.StatusInternalServerError)
	}
	if emailExist{
		errStr:="user with email "+user.Email+" already exist"
		return utilities.ShowError(c,errStr,fiber.StatusConflict)
	}
	//Check if user exist
	userExist,_,err:= model.UserExist(c,user.PhoneNumber,user.UserRole)
	if err != nil{
		return utilities.ShowError(c,err.Error(),fiber.StatusInternalServerError)
	}
	if userExist{
		errStr := "user with this phone no. "+user.PhoneNumber+" already exists"
		return utilities.ShowError(c,errStr,fiber.StatusConflict)
	}
	//validate phone number
	phone,err := utilities.ValidatePhoneNumber(user.PhoneNumber,country_code)
	if err !=nil || phone ==""{
		log.Println(err.Error())
		return utilities.ShowError(c,err.Error(),fiber.StatusAccepted)
	}

	
	
	//hash password
	hashed_password, _:= utilities.HashPassword(user.Password)

	userModel := model.User{ BaseModel: model.BaseModel{
        ID: id, 
    },FirstName: user.FirstName,LastName:user.LastName,Email: user.Email,PhoneNumber: user.PhoneNumber,Password: hashed_password,ResetCode: "",UserRole: user.UserRole,}
	//create user
	userModel.CodeExpirationTime=time.Now()

	//create user model
	err = db.Create(&userModel).Error
	if err!= nil {
		log.Fatal(err.Error())
		return utilities.ShowError(c, "failed to add data to the database",fiber.StatusInternalServerError)
	}
	
	return utilities.ShowMessage(c,"account created successfully",fiber.StatusOK)
}