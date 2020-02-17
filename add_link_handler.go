package main

import (
	"github.com/jinzhu/gorm"
	"net/http"
)

func AddLinkHandler(cfg Config, w http.ResponseWriter, r *http.Request, db *gorm.DB, cachechan chan string) bool {

	/* Parse the form and check the URL arguments */
	r.ParseForm()

	title := r.FormValue("title")
	url := r.FormValue("url")
	description := r.FormValue("description")
	password := r.FormValue("password")

	if len(title) == 0 || len(url) == 0 || len(description) == 0 || len(password) == 0 {
		/* Return an error page if we're missing a field */
		return false
	}

	if password != cfg.Post.Password {
		/* Bad password, so we should show an error page */
		return false
	}

	link := Link{Title: title, URL: url, Description: description, Active: true}
	db.Create(&link)

	/* Let's queue the link for caching */
	cachechan <- url

	return true
}
