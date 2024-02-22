package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	flag "github.com/spf13/pflag"
	"net/http"
	"os"
	"os/signal"
	"service-token/internal/feedbackGorush"
	"service-token/internal/health"
	"service-token/internal/news"
	"service-token/internal/notification"
	userToken "service-token/internal/userToken"
	"service-token/pkg/api"
	"service-token/pkg/config"
	"strconv"
	"time"
)

const file = "config/service.yaml"

var configFile string

func main() {
	flag.Parse()
	err, conf := config.Load(configFile)
	if err != nil {
		fmt.Println("Read config error")
		os.Exit(1)
	}

	repo, err := userToken.NewRepository(conf.Mongo)
	if err != nil {
		fmt.Println("could not initialize token repository")
		os.Exit(1)
	}
	userTokenService := userToken.New(repo)
	userTokenHandler := api.UserTokenHandler{
		S: userTokenService,
	}
	health.InitHealthSystem(conf.Healthy, repo)

	nService := notification.New(userTokenService)
	newsService := news.New(nService)
	newsHandler := api.NewsHandler{
		S: newsService,
	}

	feedbackGorushService := feedbackGorush.New(userTokenService)
	feedbackGorushHandler := api.FeedbackGorushHandler{
		S: feedbackGorushService,
	}

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Mount("/userToken", userTokenHandler.Routes())
		r.Mount("/feedback", feedbackGorushHandler.Routes())
		r.Mount("/news", newsHandler.Routes())
		r.Mount("/health", health.Routes())
	})
	var srv *http.Server
	srv = &http.Server{
		Addr:         "0.0.0.0:" + strconv.Itoa(conf.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
	go func() {
		fmt.Printf("starting http server on address: %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("error starting server: %s\n", err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("waiting for clients")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_ = srv.Shutdown(ctx)

	fmt.Println("finished")

	os.Exit(0)
}

func init() {
	flag.StringVarP(&configFile, "config", "c", file, "this is the path and filename to the config file")
}
