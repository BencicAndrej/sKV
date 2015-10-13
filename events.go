package main

import (
	"encoding/json"
	"strings"

	"github.com/BencicAndrej/sKV/repository"
)

const (
	EVENT_VALUE_STORED    = "event_value_stored"
	EVENT_VALUE_RETRIEVED = "event_value_retrieved"
	EVENT_VALUE_DELETED   = "event_value_deleted"
)

type EventHandler struct {
	repo *repository.Repository
}

func (handler *EventHandler) CreatedRequestKeyHandler(kv KeyValuePair) {
	if !strings.HasPrefix(kv.Key, "/request/") {
		return
	}

	outputPath := strings.Replace(kv.Key, "request", "response", 1)

	outputBody, _ := json.Marshal(map[string]string{
		"status": "done",
	})

	handler.repo.Put(outputPath, string(outputBody))

	handler.repo.Delete(kv.Key)
}