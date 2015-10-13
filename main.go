package main

import (
	"log"
	"net/http"

	"github.com/BencicAndrej/EventBus"
	"github.com/BencicAndrej/sKV/repository"

	"github.com/boltdb/bolt"
)

type KeyValuePair struct {
	Key, Value string
}

func main() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)

	eventHandler := &EventHandler{repo: repo}

	eventBus := EventBus.New()
	eventBus.Subscribe(EVENT_VALUE_STORED, eventHandler.CreatedRequestKeyHandler)

	c := &KvController{Repo: repo, EventBus: eventBus}

	http.HandleFunc("/", func(wr http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			c.Get(wr, r)
		case "PUT":
			c.Put(wr, r)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
