package app

import (
	"context"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wbL0/server_sub/cache"
	"wbL0/server_sub/config"
	"wbL0/server_sub/db"
	"wbL0/server_sub/server"
	"wbL0/server_sub/store"
	"wbL0/server_sub/subscriber"
)

type App struct {
	cfg config.Config
}

func InitApp(cfg config.Config) *App {
	app := App{}
	app.cfg = cfg
	return &app
}

func (app *App) Run() {
	db, err := db.InitDBConn(app.cfg)
	if err != nil {
		log.Println("Error while connnecting to database")
		panic(err)
	}

	cacheService := cache.CacheInit()

	storeService := store.InitStore(*cacheService, *db)

	sc := subscriber.CreateSub(*storeService)
	err = sc.Connect(app.cfg.Nats_server.Cluster_id, app.cfg.Nats_server.Client_id, app.cfg.Nats_server.Host+":"+app.cfg.Nats_server.Port)

	if err != nil {
		log.Println("Error while connecting to STAN")
	}

	sub, err := sc.SubscribeToChannel(app.cfg.Nats_server.Channel, stan.StartWithLastReceived())

	if err != nil {
		log.Println("Error while subscribing to channel")
	}

	err = storeService.RestoreCache()

	if err != nil {
		log.Println("error restoring cache: db is empty")
	}

	server := server.InitServer(*storeService, app.cfg.Http_server.Host+":"+app.cfg.Http_server.Port)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
		defer server.Stop()
		defer sub.Unsubscribe()
		defer sc.Close()
		defer db.Close()
	}()

	if err := server.Srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited")
}
