package main

import (
	"time"
)

type Link struct {
	ID           int
	Title     	 string
	URL     	 string
	Description  string
	Active		 bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}
