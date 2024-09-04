package service

import (
	"net/http"
	"tg-backend/server/util"
)

func HealthCheck() util.HttpHandler {
	return func(r *http.Request) util.HandleResult {
		payload := struct {
			Alive bool `json:"alive"`
		}{true}

		return util.HandleResult{
			Payload: payload,
			Type:    util.ResponseTypeJson,
		}
	}
}
