package model

func MigrateDB(){
	db.AutoMigrate(
		&User{},
		&Product{},
		&Order{},
		&Service{},
		OrderItem{},
		Cart{},
		CartItem{},
	)
}