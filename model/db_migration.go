package model

func MigrateDB(){
	db.AutoMigrate(
		User{},
		Product{},
		Order{},
		OrderItem{},
		Cart{},
		CartItem{},
	)
}