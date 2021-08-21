package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"srv-1/internal/app/config"
)

func main() {
	cfg, err := config.GetConf()
	if err != nil {
		log.Fatal("Cannot read config file")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow",
		cfg.Postgres.Host,
		cfg.Postgres.User,
		cfg.Postgres.Pwd,
		cfg.Postgres.DB,
		cfg.Postgres.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println(db)

}