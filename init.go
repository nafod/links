package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math/rand"
	"runtime"
	"time"
)

func Init(cfg Config) (*gorm.DB) {

	// Not needed on Go 1.5+
	runtime.GOMAXPROCS(runtime.NumCPU() + 1)

	rand.Seed(time.Now().UTC().UnixNano())

	db, err := gorm.Open("sqlite3", cfg.Database.DSN)
	if err != nil {
		log.Fatalf("[ERR] Unable to connect to database server!")
	}

	err = db.DB().Ping()
	if err != nil {
		log.Fatalf("[ERR] Unable to ping database server!")
	}

	log.Printf("[LOG] Successfully connected to database")

	return &db
}
