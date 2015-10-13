package main

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/net/html"

	"github.com/BencicAndrej/EventBus"
	"github.com/BencicAndrej/sKV/repository"
)

type Action func(rw http.ResponseWriter, r *http.Request) error

type BaseController struct{}

func (c *BaseController) Action(a Action) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err := a(rw, r); err != nil {
			http.Error(rw, err.Error(), 500)
		}
	})
}

type KvController struct {
	BaseController
	Repo     *repository.Repository
	EventBus *EventBus.EventBus
}

func (c *KvController) Get(rw http.ResponseWriter, r *http.Request) error {
	record := c.Repo.Get(html.EscapeString(r.URL.Path))

	if len(record) == 0 {
		rw.WriteHeader(http.StatusNotFound)
		return nil
	}

	rw.Write([]byte(record))

	c.EventBus.Publish(EVENT_VALUE_RETRIEVED, KeyValuePair{Key: html.EscapeString(r.URL.Path), Value: record})

	return nil
}

func (c *KvController) Put(rw http.ResponseWriter, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = c.Repo.Put(html.EscapeString(r.URL.Path), string(body))
	if err != nil {
		return err
	}

	rw.WriteHeader(http.StatusCreated)

	c.EventBus.Publish(EVENT_VALUE_STORED, KeyValuePair{Key: html.EscapeString(r.URL.Path), Value: string(body)})

	return nil
}
