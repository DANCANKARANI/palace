package model

func MigrateDB(){
	db.AutoMigrate(
		&User{},
		&Rating{},
		&Product{},
		&Order{},
		&Service{},
		&OrderItem{},
		&Cart{},
		&CartItem{},
		&Payment{},
	)
}