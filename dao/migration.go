package dao

func migration() {
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate()
	if err != nil {
		panic(err)
	}
}
