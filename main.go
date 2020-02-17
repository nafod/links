package main

import (
	"gopkg.in/gcfg.v1"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
)

type Config struct {
	Main struct {
		Debug  bool
		IP     string
		Domain string
	}

	Post struct {
		Password string
	}

	Database struct {
		DSN string
	}
}

type Page struct {
    Links []Link
}

func main() {

	/* Load the configuration */
	var cfg Config

	err := gcfg.ReadFileInto(&cfg, "links.conf")
	if err != nil {
		/* Couldn't read the config file */
		panic(err)
	}

	/* Create initial directories, sets GOMAXPROC, and seeds the PRNG */
	db := Init(cfg)

	if cfg.Main.Debug {
		/* db.LogMode(true) */
		db.AutoMigrate(&Link{})
	}

	defer db.Close()

	/* Start up link caching process */
	cachechan := make(chan string)
	go LinkCacher(cachechan, db)

	/* Start the check for dead links */
	go CheckDeadLinksLoop(db)

	router := httprouter.New()

	/* Main handler */
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var dblinks []Link
		db.Find(&dblinks)
		links := &Page{Links: dblinks}
	    t, err := template.ParseFiles("templates/index.html")
	    if err != nil {
	    	log.Printf("[ERR] Unable to load index page template")
	    }
	    t.Execute(w, links)
	})


	router.GET("/new", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	    t, err := template.ParseFiles("templates/new.html")
	    if err != nil {
	    	log.Printf("[ERR] Unable to load new link page template")
	    }
	    t.Execute(w, &Page{})
	})

	router.POST("/post", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if AddLinkHandler(cfg, w, r, db, cachechan) {
			// Success page
		    t, _ := template.ParseFiles("templates/success.html")
		    t.Execute(w, nil)
		} else {
			// Error
		    t, _ := template.ParseFiles("templates/error.html")
		    t.Execute(w, nil)
		}
	})

	log.Printf("[MAIN] Now listening on %s", cfg.Main.IP)
	log.Fatal(http.ListenAndServe(cfg.Main.IP, router))
}
