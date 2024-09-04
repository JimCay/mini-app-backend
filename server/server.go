package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"tg-backend/bot"
	"tg-backend/config"
	"tg-backend/db"
	"tg-backend/pkg/log"
	"tg-backend/server/middleware"
	"tg-backend/server/service"
	"tg-backend/server/util"
	"time"
)

type Server struct {
	httpServer *http.Server
	tgBot      *bot.TelegramBot
}

func NewServer(storage *db.Storage, config *config.Config) *Server {
	serviceManager := service.NewServiceManager(storage, config)
	httpServ := NewHttpServer(serviceManager, config)
	myBot := bot.NewTelegramBot(config.TgConf)
	return &Server{
		httpServer: httpServ,
		tgBot:      myBot,
	}
}

func NewHttpServer(serviceManager *service.ServiceManager, config *config.Config) *http.Server {

	r := mux.NewRouter()
	r.HandleFunc("/health", util.ResponseWrapper(service.HealthCheck())).Methods("GET")
	apiRouter := r.PathPrefix("/api").Subrouter()
	authMiddleware := middleware.NewTelegramAuthMiddleware(config.TgConf.TelegramBotToken)
	apiRouter.Use(authMiddleware)

	apiRouter.HandleFunc("/user/login", util.ResponseWrapper(
		service.LoginHandler(serviceManager.User))).Methods("GET")

	server := &http.Server{
		Addr:              ":" + config.TgConf.Port,
		Handler:           r,
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Second * 30,
	}
	return server
}

func (s *Server) Start() {
	log.Info("start http server")
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered. Error:\n", r)
			}
		}()
		if err := s.httpServer.ListenAndServeTLS("./cert.pem", "key.pem"); err != nil {
			log.Error("httpServer error ", err)
		}
	}()
	log.Info("start tg bot")
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("tg bot recovered. Panic:\n", r)
			}
		}()
		s.tgBot.Start()
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return nil

}
