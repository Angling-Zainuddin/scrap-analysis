package main

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDbConnection() (*gorm.DB, error) {
	dsn := "host=localhost user=crmapp password=crm@RumahSatu23! dbname=postgres port=5436 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}
