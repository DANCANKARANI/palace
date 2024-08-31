package model

func MigrateDB(){
	db.AutoMigrate(
		User{},
		Clothe{},
	)
}