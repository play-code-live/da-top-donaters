package main

import (
	"He110/donation-report-manager/internal/business/domain/configs"
	"He110/donation-report-manager/internal/business/domain/donations"
	"He110/donation-report-manager/internal/business/domain/user"
	getConfigUseCase "He110/donation-report-manager/internal/business/use_cases/get_config_use_case"
	getTopDonatersUseCase "He110/donation-report-manager/internal/business/use_cases/get_top_donaters_use_case"
	saveConfigUseCase "He110/donation-report-manager/internal/business/use_cases/save_config_use_case"
	saveTokenUseCase "He110/donation-report-manager/internal/business/use_cases/save_token_use_case"
	donationClient "He110/donation-report-manager/internal/pkg/donation-alerts-client"
	"He110/donation-report-manager/internal/web"
	getTopDonaters "He110/donation-report-manager/internal/web/api/get_top_donaters"
	"context"
	"fmt"
	"github.com/eko/gocache/cache"
	"github.com/eko/gocache/store"
	"github.com/gorilla/mux"
	goCache "github.com/patrickmn/go-cache"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := GetConfig()
	client, err := donationClient.NewClient(cfg.DaClientId, cfg.DaClientSecret, fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		panic(err)
	}

	cacheClient := goCache.New(10*time.Minute, 30*time.Minute)
	cacheStore := store.NewGoCache(cacheClient, nil)
	cacheManager := cache.New(cacheStore)

	uStorage := user.NewStorage()
	cStorage := configs.NewStorage()

	daService := donations.NewService(client, cacheManager)

	useCaseGetConfig := getConfigUseCase.New(uStorage, cStorage)
	useCaseSaveToken := saveTokenUseCase.New(uStorage)
	useCaseSaveConfig := saveConfigUseCase.New(uStorage, cStorage, daService)
	useCaseGetTopDonaters := getTopDonatersUseCase.New(uStorage, cStorage, daService)

	container := web.NewTemplateContainer("templates/base/*.gohtml")
	if err = container.FindAndRegister("templates/pages/"); err != nil {
		panic(err)
	}

	errGroup, ctx := errgroup.WithContext(context.Background())
	app := web.NewApp(container, client, useCaseGetConfig, useCaseSaveConfig, useCaseSaveToken)

	router := mux.NewRouter()

	router.HandleFunc("/panel", app.HandlerPanel())

	router.HandleFunc("/redirect", app.HandlerRedirect())
	router.HandleFunc("/redirect/{channelId}", app.HandlerChanneledRedirect())

	router.Path("/config/{channelId}").Methods(http.MethodGet).HandlerFunc(app.HandlerGetConfig(cfg.SocketAddress))
	router.Path("/config/{channelId}").Methods(http.MethodPost).HandlerFunc(app.HandlerSaveConfig())
	router.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		if err = container.MustGet("config_anonymous").Execute(w, nil); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	router.HandleFunc("/socket/{channelId}", app.SocketBridge())

	router.HandleFunc("/api/donaters", getTopDonaters.NewTopDonatersHandler(useCaseGetTopDonaters))

	ctx, _ = signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	// #### Web Server
	errGroup.Go(func() error {
		server := http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: router,
		}
		log.Printf("webserver has started on :%d\n", cfg.Port)
		return server.ListenAndServe()
	})

	<-ctx.Done()
	log.Println("service gracefully shutdown")
}
