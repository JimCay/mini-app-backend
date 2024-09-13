package server

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"io"
	"net/http"
	"tg-backend/bot"
	"tg-backend/config"
	"tg-backend/db"
	_ "tg-backend/docs"
	"tg-backend/pkg/log"
	"tg-backend/server/middleware"
	"tg-backend/server/service"
	"tg-backend/server/util"
	"time"
)

type Server struct {
	httpServer *http.Server
	tgBot      *bot.TelegramBot
	conf       *config.Config
}

func NewServer(storage db.Storage, config *config.Config, writer io.Writer) *Server {
	serviceManager := service.NewServiceManager(storage, config)
	httpServ := NewHttpServer(serviceManager, config, writer)
	myBot := bot.NewTelegramBot(config.TgConf)
	return &Server{
		httpServer: httpServ,
		tgBot:      myBot,
		conf:       config,
	}
}

// NewHttpServer
// @title tg backend API
// @version 1.0
// @description tg backend API Document
// @BasePath /
func NewHttpServer(sm *service.ServiceManager, config *config.Config, writer io.Writer) *http.Server {
	r := mux.NewRouter()
	if config.TgConf.Swagger {
		r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	}
	r.HandleFunc("/health", util.ResponseWrapper(service.HealthCheck())).Methods("GET")
	apiRouter := r.PathPrefix("/api").Subrouter()
	authMiddleware := middleware.NewTelegramAuthMiddleware(config.TgConf.TelegramBotToken, config.TgConf.Expire)
	apiRouter.Use(authMiddleware)

	handlerRouter(apiRouter, "/user/login", "POST", service.LoginHandler(sm.User, sm.Point))
	handlerRouter(apiRouter, "/user/invite", "GET", service.InviteHandler(sm.User))
	handlerRouter(apiRouter, "/user/friends", "GET", service.FriendHandler(sm.User))

	handlerRouter(apiRouter, "/point/query", "GET", service.GetPointHandler(sm.Point))
	handlerRouter(apiRouter, "/point/update", "POST", service.UpdatePointHandler(sm.Point))
	handlerRouter(apiRouter, "/point/rank", "GET", service.GetRankHandler(sm.Point))

	handlerRouter(apiRouter, "/task/get", "GET", service.GetTasksHandler(sm.Task))
	handlerRouter(apiRouter, "/task/check", "POST", service.TaskCheckHandler(sm.Task))

	handler := handlers.CombinedLoggingHandler(writer, r)

	server := &http.Server{
		Addr:              ":" + config.TgConf.Port,
		Handler:           handlerCors(handler),
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Second * 30,
	}
	return server
}

func handlerRouter(router *mux.Router, path, method string, handler util.HttpHandler) {
	router.HandleFunc(path, util.ResponseWrapper(handler)).Methods(method)
}

func handlerCors(h http.Handler) http.Handler {
	originsOk := handlers.AllowedOrigins([]string{"*"})
	headers := handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT"})
	return handlers.CORS(originsOk, headers, methods)(h)
}

func (s *Server) Start() {

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered. Error:\n", r)
			}
		}()
		if s.conf.TgConf.Ssl {
			log.Info("start https server")
			if err := s.httpServer.ListenAndServeTLS("./cert.pem", "key.pem"); err != nil {
				log.Error("httpServer error ", err)
			}
		} else {
			log.Info("start http server")
			if err := s.httpServer.ListenAndServe(); err != nil {
				log.Error("httpServer error ", err)
			}
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
