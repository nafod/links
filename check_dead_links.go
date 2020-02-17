package main

import (
	"log"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

func CheckDeadLinks(db *gorm.DB) {
	var links []Link
	db.Find(&links)

	active_urls := 0
	notfound := 0

	log.Printf("[CHECK] Checking for dead links")

	for _, elem := range links {

		/* If it's already dead, we just count it and move on */
		if elem.Active == false {
			notfound = notfound + 1
			continue
		}

		resp, err := http.Get(elem.URL)
		if err != nil {
			// There was some non-HTTP code related error, let's just log it and move on
			log.Printf("[CHECK] Unable to fetch URL %+v\n", err)
			continue
		}

		defer resp.Body.Close()

		// Check the status code
		if resp.StatusCode == 200 {
			/*
				We don't need to do anything right now. Later we'll check the actual contents
				of the page to see how much has changed and redownload if neccesary.
			*/
			active_urls = active_urls + 1
		} else if resp.StatusCode == 404 {
			/* Page has 404'd */
			notfound = notfound + 1
			elem.Active = false
			db.Save(&elem)
		} else {
			/* Some other response code */
			log.Printf("[CHECK] Response code %d from URL %s\n", resp.StatusCode, elem.URL)
			active_urls = active_urls + 1
		}
	}

	log.Printf("[CHECK] Active: %d | Dead: %d | Total: %d\n", active_urls, notfound, active_urls + notfound)
}

func CheckDeadLinksLoop(db *gorm.DB) {
	log.Printf("[CHECK] CheckDeadLinks thread started")
	for {
		time.Sleep(12 * time.Hour)
		CheckDeadLinks(db)
	}
}

