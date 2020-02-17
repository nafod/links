package main

import (
	"log"
	"github.com/jinzhu/gorm"
	"os/exec"
)

func LinkCacher(cachechan chan string, db *gorm.DB) {
	for {
		url := <- cachechan
		log.Printf("[CACHE] Caching %s\n", url)
		err := exec.Command("wget", "--mirror", "--page-requisites", "--wait=1", "--level=3", "--adjust-extension", "--no-parent", "--convert-links", "--directory-prefix=cached", url).Run()
		if err != nil {
			log.Printf("[CACHE] Error: %s\n", err)
		}
		log.Printf("[CACHE] Successfully cached %s\n", url)
	}
}
